"""FastAPI-facing facade for the Wind ReAct Agent runtime."""

from __future__ import annotations

import uuid
from collections.abc import AsyncIterator

from app.schemas.agent import AgentRunResponse, AgentStreamEvent
from app.services.agent_runtime import AgentRuntime
from app.services.tool.wind_agent_tools import GO_TOOLS


ALLOWED_TOOLS = GO_TOOLS


async def stream_agent(
    runtime: AgentRuntime,
    *,
    user_id: int,
    agent_session_id: str,
    user_input: str,
) -> AsyncIterator[AgentStreamEvent]:
    """Start a fresh Agent turn and expose LangGraph custom events."""

    resolved_agent_session_id = agent_session_id.strip() or str(uuid.uuid4())
    # 准备langgraph状态机输入
    async for event in runtime.stream_new_turn(
        user_id=user_id,
        agent_session_id=resolved_agent_session_id,
        user_input=user_input,
    ):
        yield event


async def resume_agent(
    runtime: AgentRuntime,
    *,
    user_id: int,
    agent_session_id: str,
    action: str,
    content: str,
) -> AsyncIterator[AgentStreamEvent]:
    """Resume a clarification or draft confirmation interrupt."""

    async for event in runtime.stream_resume(
        user_id=user_id,
        agent_session_id=agent_session_id.strip(),
        action=action,
        content=content,
    ):
        yield event


async def run_agent(
    runtime: AgentRuntime,
    *,
    user_id: int,
    agent_session_id: str,
    user_input: str,
) -> AgentRunResponse:
    """Aggregate the streaming Agent path for diagnostics and internal callers."""

    answer = ""
    trace_id = ""
    resolved_agent_session_id = agent_session_id.strip()
    tool_calls = []
    citations = []
    draft = None
    status = "done"
    async for event in stream_agent(
        runtime,
        user_id=user_id,
        agent_session_id=resolved_agent_session_id,
        user_input=user_input,
    ):
        trace_id = event.trace_id or trace_id
        resolved_agent_session_id = event.agent_session_id or resolved_agent_session_id
        if event.type == "done":
            answer = event.content
            tool_calls = event.tool_calls
            citations = event.citations
            draft = event.draft
        elif event.type in {"confirmation_required", "clarification_required"}:
            answer = event.content
            status = event.type
        elif event.type == "error":
            answer = event.error_msg or event.content
            status = "error"
    return AgentRunResponse(
        answer=answer,
        tool_calls=tool_calls,
        trace_id=trace_id,
        agent_session_id=resolved_agent_session_id,
        citations=citations,
        draft=draft,
        status=status,
    )
