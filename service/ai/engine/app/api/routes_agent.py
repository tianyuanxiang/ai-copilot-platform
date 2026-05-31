"""HTTP and SSE routes for the Wind ReAct Agent."""

import json

from fastapi import APIRouter, Request
from fastapi.responses import StreamingResponse

from app.schemas.agent import AgentResumeRequest, AgentRunRequest, AgentRunResponse, AgentStreamEvent
from app.services.agent import resume_agent, run_agent, stream_agent
from app.services.agent_runtime import AgentRuntime


router = APIRouter(tags=["agent"])


@router.post("/agent/run", response_model=AgentRunResponse)
async def agent_run(payload: AgentRunRequest, request: Request) -> AgentRunResponse:
    """Run the same Agent graph as SSE and aggregate its final state."""

    return await run_agent(
        _runtime(request),
        user_id=payload.user_id,
        conversation_id=payload.conversation_id,
        user_input=payload.input,
    )


@router.post("/agent/stream")
async def agent_stream(payload: AgentRunRequest, request: Request) -> StreamingResponse:
    """Start one Agent turn and stream ReAct events as SSE."""

    async def event_generator():
        disconnected = False
        async for event in stream_agent(
            _runtime(request),
            user_id=payload.user_id,
            conversation_id=payload.conversation_id,
            user_input=payload.input,
        ):
            if await request.is_disconnected():
                disconnected = True
                break
            yield _sse(event)
        if not disconnected:
            yield "data: [DONE]\n\n"

    return _streaming_response(event_generator())


@router.post("/agent/resume/stream")
async def agent_resume_stream(payload: AgentResumeRequest, request: Request) -> StreamingResponse:
    """Resume a clarification or draft-approval interrupt and continue SSE."""

    async def event_generator():
        disconnected = False
        async for event in resume_agent(
            _runtime(request),
            user_id=payload.user_id,
            conversation_id=payload.conversation_id,
            action=payload.action,
            content=payload.content,
        ):
            if await request.is_disconnected():
                disconnected = True
                break
            yield _sse(event)
        if not disconnected:
            yield "data: [DONE]\n\n"

    return _streaming_response(event_generator())


def _runtime(request: Request) -> AgentRuntime:
    runtime = getattr(request.app.state, "agent_runtime", None)
    if runtime is None:
        raise RuntimeError("Wind Agent runtime is not initialized")
    return runtime


def _sse(event: AgentStreamEvent) -> str:
    return "data: " + json.dumps(event.model_dump(exclude_none=True), ensure_ascii=False) + "\n\n"


def _streaming_response(events) -> StreamingResponse:
    return StreamingResponse(
        events,
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "X-Accel-Buffering": "no",
        },
    )
