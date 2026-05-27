import asyncio
import json
import logging
import time
import uuid
from collections.abc import AsyncIterator

import httpx

from app.core.config import get_settings
from app.schemas.chat import ChatStreamRequest, StreamEvent


logger = logging.getLogger(__name__)


async def stream_chat(request: ChatStreamRequest, trace_id: str | None = None) -> AsyncIterator[StreamEvent]:
    current_trace_id = trace_id or str(uuid.uuid4())
    settings = get_settings()
    started_at = time.perf_counter()
    logger.info(
        "chat.stream开始：trace_id=%s provider=%s model=%s user_id=%s kb_id=%s conversation_id=%s question_chars=%s history_count=%s",
        current_trace_id,
        settings.llm_provider,
        settings.llm_model,
        request.user_id,
        request.kb_id,
        request.conversation_id,
        len(request.question),
        len(request.history),
    )
    if settings.mock_llm or settings.llm_provider == "mock":
        async for event in _log_stream_events(_mock_stream_chat(current_trace_id), current_trace_id, started_at):
            yield event
        return

    if settings.llm_provider == "deepseek":
        # 从下层函数一个一个取出 StreamEvent  原封不动地往上层传
        async for event in _log_stream_events(_stream_deepseek_chat(request, current_trace_id), current_trace_id, started_at):
            yield event
        return

    async for event in _log_stream_events(
        _static_events(
            StreamEvent(
                type="error",
                content=f"unsupported llm provider: {settings.llm_provider}",
                trace_id=current_trace_id,
            ),
            StreamEvent(type="done", content="", trace_id=current_trace_id),
        ),
        current_trace_id,
        started_at,
    ):
        yield event


async def _log_stream_events(
    events: AsyncIterator[StreamEvent],
    trace_id: str,
    started_at: float,
) -> AsyncIterator[StreamEvent]:
    token_count = 0
    answer_parts: list[str] = []
    async for event in events:
        if event.type == "token":
            token_count += 1
            answer_parts.append(event.content)
            logger.debug(
                "llm.token trace_id=%s token_index=%s token_chars=%s",
                trace_id,
                token_count,
                len(event.content),
            )
        elif event.type == "error":
            logger.error("chat.stream.error trace_id=%s content=%s", trace_id, event.content)
        yield event
    logger.info(
        "chat.stream.done trace_id=%s token_count=%s answer_chars=%s duration_ms=%.2f",
        trace_id,
        token_count,
        len("".join(answer_parts)),
        (time.perf_counter() - started_at) * 1000,
    )


async def _static_events(*events: StreamEvent) -> AsyncIterator[StreamEvent]:
    for event in events:
        yield event


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

    # 组装请求
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
    logger.info(
        "llm.request trace_id=%s provider=deepseek model=%s stream=%s messages=%s api_key=%s url=%s",
        trace_id,
        settings.llm_model,
        body["stream"],
        len(body["messages"]),
        "present" if settings.llm_api_key else "empty",
        settings.deepseek_chat_url,
    )

    try:
        timeout = httpx.Timeout(float(settings.llm_timeout_seconds))
        async with httpx.AsyncClient(timeout=timeout) as client:
            async with client.stream(   # 发送 HTTP 流式请求
                "POST",
                settings.deepseek_chat_url,
                headers=headers,
                json=body,
            ) as response:
                if response.status_code < 200 or response.status_code >= 300:
                    content = await response.aread()
                    logger.error(
                        "llm.provider_error trace_id=%s provider=deepseek status=%s response=%s",
                        trace_id,
                        response.status_code,
                        content.decode("utf-8", errors="ignore")[:512],
                    )
                    yield StreamEvent(
                        type="error",
                        content=f"deepseek http {response.status_code}: {content.decode('utf-8', errors='ignore')}",
                        trace_id=trace_id,
                    )
                    yield StreamEvent(type="done", content="", trace_id=trace_id)
                    return

                async for line in response.aiter_lines():
                    token = _parse_deepseek_sse_line(line) # 解析每行
                    if token is None:
                        continue
                    yield StreamEvent(type="token", content=token, trace_id=trace_id)
    except Exception as exc:
        logger.exception("deepseek stream chat failed trace_id=%s", trace_id)
        yield StreamEvent(type="error", content=f"deepseek stream failed: {exc}", trace_id=trace_id)
        yield StreamEvent(type="done", content="", trace_id=trace_id)
        return

    # 全部读完，发结束信号
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
