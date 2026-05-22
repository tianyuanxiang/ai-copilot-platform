import json
import uuid

from fastapi import APIRouter, Request
from fastapi.responses import StreamingResponse

from app.schemas.chat import ChatStreamRequest, StreamEvent
from app.services.llm import stream_chat

router = APIRouter(tags=["chat"])


def _sse(event: StreamEvent) -> str:
    payload = event.model_dump()
    return "data: " + json.dumps(payload, ensure_ascii=False) + "\n\n"


@router.post("/chat/stream")
async def chat_stream(payload: ChatStreamRequest, request: Request) -> StreamingResponse:
    trace_id = str(uuid.uuid4())

    async def event_generator():
        async for event in stream_chat(payload, trace_id=trace_id):
            if await request.is_disconnected():
                break
            yield _sse(event)
        yield "data: [DONE]\n\n"

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "X-Accel-Buffering": "no",
        },
    )