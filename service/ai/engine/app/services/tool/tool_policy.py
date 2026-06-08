from dataclasses import dataclass

from app.services.tool.tool_errors import ToolErrorType

@dataclass(frozen=True, slots=True)
class ToolPolicy:
    timeout_ms: int
    max_retries: int = 0
    retryable_errors: frozenset[ToolErrorType] = frozenset()
    fallback: str = "degraded"

DEFAULT_TOOL_POLICY = ToolPolicy(
    timeout_ms=2000,
    max_retries=0,
)

TOOL_POLICIES: dict[str, ToolPolicy] = {
    "get_turbine_metadata": ToolPolicy(
        timeout_ms=1600,
        max_retries=1,
        retryable_errors=frozenset({
            ToolErrorType.RPC_TIMEOUT,
            ToolErrorType.RPC_UNAVAILABLE,
        }),
        fallback="clarify",
    ),
    "search_maintenance_sop": ToolPolicy(
        timeout_ms=6000,
        max_retries=1,
        retryable_errors=frozenset({
            ToolErrorType.RPC_TIMEOUT,
            ToolErrorType.RPC_UNAVAILABLE,
            ToolErrorType.TOOL_FAILED,
        }),
        fallback="answer_without_sop",
    ),
    "query_alarm_events": ToolPolicy(
        timeout_ms=4000,
        max_retries=1,
        retryable_errors=frozenset({
            ToolErrorType.RPC_TIMEOUT,
            ToolErrorType.RPC_UNAVAILABLE,
            ToolErrorType.TOOL_FAILED,
        }),
        fallback="partial_answer",
    ),
    "query_sensor_timeseries": ToolPolicy(
        timeout_ms=6000,
        max_retries=1,
        retryable_errors=frozenset({
            ToolErrorType.RPC_TIMEOUT,
            ToolErrorType.RPC_UNAVAILABLE,
        }),
        fallback="partial_answer",
    ),
    "compare_sensor_trend": ToolPolicy(
        timeout_ms=4000,
        max_retries=0,
        retryable_errors=frozenset(),
        fallback="raw_data_answer",
    ),
}


def get_tool_policy(tool_name: str) -> ToolPolicy:
    return TOOL_POLICIES.get(tool_name, DEFAULT_TOOL_POLICY)
