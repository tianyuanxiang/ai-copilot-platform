import uuid

from app.schemas.agent import AgentRunResponse, ToolCall


ALLOWED_TOOLS = {
    "search_maintenance_sop",
    "query_sensor_timeseries",
    "query_alarm_events",
    "get_turbine_metadata",
    "compare_sensor_trend",
    "generate_health_report",
    "create_maintenance_ticket_draft",
}


async def run_agent(user_id: str, user_input: str) -> AgentRunResponse:
    tool_calls = [
        ToolCall(name=name, status="reserved", message="allowlisted read-only or controlled tool")
        for name in sorted(ALLOWED_TOOLS)
    ]
    return AgentRunResponse(
        answer="Wind O&M tool calling scaffold is ready. Go-side whitelist, evidence assembly, and audit logging remain the control plane.",
        tool_calls=tool_calls,
        trace_id=str(uuid.uuid4()),
    )
