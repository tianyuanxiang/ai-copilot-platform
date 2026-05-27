import asyncio
import json
import logging
import uuid
from collections.abc import AsyncIterator

import httpx

from app.core.config import get_settings
from app.schemas.chat import ChatStreamRequest, StreamEvent


logger = logging.getLogger(__name__)


async def stream_chat(request: ChatStreamRequest, trace_id: str | None = None) -> AsyncIterator[StreamEvent]:
    current_trace_id = trace_id or str(uuid.uuid4())
    settings = get_settings()
    if settings.mock_llm or settings.llm_provider == "mock":
        async for event in _mock_stream_chat(current_trace_id):
            yield event
        return

    if settings.llm_provider == "deepseek":
        async for event in _stream_deepseek_chat(request, current_trace_id):
            yield event
        return

    yield StreamEvent(
        type="error",
        content=f"unsupported llm provider: {settings.llm_provider}",
        trace_id=current_trace_id,
    )
    yield StreamEvent(type="done", content="", trace_id=current_trace_id)


async def _mock_stream_chat(trace_id: str) -> AsyncIterator[StreamEvent]:
    text = (
        "ASGI streaming is ready. "
        "This mock answer is emitted token by token, "
        "so you can replace this function with a real LLM streaming client later."
    )

    for token in text.split(" "):
        await asyncio.sleep(0.04)
        yield StreamEvent(type="token", content=token + " ", trace_id=trace_id)
    yield StreamEvent(type="done", content="", trace_id=trace_id)


async def _stream_deepseek_chat(request: ChatStreamRequest, trace_id: str) -> AsyncIterator[StreamEvent]:
    settings = get_settings()
    if not settings.llm_api_key:
        yield StreamEvent(type="error", content="llm api key is empty", trace_id=trace_id)
        yield StreamEvent(type="done", content="", trace_id=trace_id)
        return

    body = {
        "model": settings.llm_model,
        "stream": True,
        "messages": _build_messages(request),
    }
    headers = {
        "Authorization": f"Bearer {settings.llm_api_key}",
        "Content-Type": "application/json",
        "Accept": "text/event-stream",
    }

    try:
        timeout = httpx.Timeout(float(settings.llm_timeout_seconds))
        async with httpx.AsyncClient(timeout=timeout) as client:
            async with client.stream(
                "POST",
                settings.deepseek_chat_url,
                headers=headers,
                json=body,
            ) as response:
                if response.status_code < 200 or response.status_code >= 300:
                    content = await response.aread()
                    yield StreamEvent(
                        type="error",
                        content=f"deepseek http {response.status_code}: {content.decode('utf-8', errors='ignore')}",
                        trace_id=trace_id,
                    )
                    yield StreamEvent(type="done", content="", trace_id=trace_id)
                    return

                async for line in response.aiter_lines():
                    token = _parse_deepseek_sse_line(line)
                    if token is None:
                        continue
                    yield StreamEvent(type="token", content=token, trace_id=trace_id)
    except Exception as exc:
        logger.exception("deepseek stream chat failed")
        yield StreamEvent(type="error", content=f"deepseek stream failed: {exc}", trace_id=trace_id)
        yield StreamEvent(type="done", content="", trace_id=trace_id)
        return

    yield StreamEvent(type="done", content="", trace_id=trace_id)


def _build_messages(request: ChatStreamRequest) -> list[dict[str, str]]:
    messages: list[dict[str, str]] = []
    for item in request.history:
        role = _normalize_role(item.role)
        content = item.content.strip()
        if not role or not content:
            continue
        messages.append({"role": role, "content": content})
    messages.append({"role": "user", "content": request.question})
    return messages


def _normalize_role(role: str) -> str:
    normalized = role.strip().lower()
    if normalized in {"system", "user", "assistant"}:
        return normalized
    return ""


def _parse_deepseek_sse_line(line: str) -> str | None:
    line = line.strip()
    if not line or not line.startswith("data:"):
        return None

    data = line.removeprefix("data:").strip()
    if not data or data == "[DONE]":
        return None

    payload = json.loads(data)
    choices = payload.get("choices") or []
    if not choices:
        return None
    delta = choices[0].get("delta") or {}
    content = delta.get("content")
    if not content:
        return None
    return str(content)
