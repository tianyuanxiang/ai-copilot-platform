import asyncio
import uuid
from collections.abc import AsyncIterator

from app.schemas.chat import ChatStreamRequest, StreamEvent


async def stream_chat(request: ChatStreamRequest, trace_id: str | None = None) -> AsyncIterator[StreamEvent]:
    current_trace_id = trace_id or str(uuid.uuid4())
    text = (
        "ASGI streaming is ready. "
        "This mock answer is emitted token by token, "
        "so you can replace this function with a real LLM streaming client later."
    )

    for token in text.split(" "):
        # 异步等待0.04秒
        await asyncio.sleep(0.04)
        yield StreamEvent(type="token", content=token + " ", trace_id=current_trace_id)
    # 循环结束后，最后再发一个done事件
    yield StreamEvent(type="done", content="", trace_id=current_trace_id)