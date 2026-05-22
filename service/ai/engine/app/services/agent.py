import uuid

from app.schemas.agent import AgentRunResponse, ToolCall


ALLOWED_TOOLS = {
    "search_knowledge_base",
    "query_security_events",
    "query_security_alerts",
    "get_daily_report",
    "generate_daily_report",
}


async def run_agent(user_id: str, user_input: str) -> AgentRunResponse:
    tool_calls = [
        ToolCall(name=name, status="reserved", message="allowlisted read-only or controlled tool")
        for name in sorted(ALLOWED_TOOLS)
    ]
    return AgentRunResponse(
        answer="Controlled tool calling skeleton is ready. The LLM will suggest, and Go-side permission checks will decide.",
        tool_calls=tool_calls,
        trace_id=str(uuid.uuid4()),
    )