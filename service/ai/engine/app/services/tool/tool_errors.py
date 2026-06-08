from enum import Enum


class ToolErrorType(str, Enum):
    INVALID_ARGUMENTS = "invalid_arguments"
    PERMISSION_DENIED = "permission_denied"
    RPC_TIMEOUT = "rpc_timeout"
    RPC_UNAVAILABLE = "rpc_unavailable"
    TOOL_FAILED = "tool_failed"
    TOOL_PARTIAL = "tool_partial"
    LLM_PLAN_ERROR = "llm_plan_error"
    UNKNOWN = "unknown"


def classify_tool_error(status: str, message: str) -> ToolErrorType:
    s = (status or "").lower()
    m = (message or "").lower()

    if s in {"invalid_arguments", "bad_request", "invalid"}:
        return ToolErrorType.INVALID_ARGUMENTS

    if s in {"denied", "permission_denied", "forbidden"}:
        return ToolErrorType.PERMISSION_DENIED

    if s in {"timeout", "deadline_exceeded"} or "timeout" in m or "deadline" in m:
        return ToolErrorType.RPC_TIMEOUT

    if s in {"unavailable", "rpc_unavailable"} or "connection refused" in m or "unavailable" in m:
        return ToolErrorType.RPC_UNAVAILABLE

    if s == "partial":
        return ToolErrorType.TOOL_PARTIAL

    if s == "failed":
        return ToolErrorType.TOOL_FAILED

    return ToolErrorType.UNKNOWN


def is_retryable_error(error_type: ToolErrorType) -> bool:
    """
    只有临时性错误才重试。
    参数错误、权限错误不重试。
    """
    return error_type in {
        ToolErrorType.RPC_TIMEOUT,
        ToolErrorType.RPC_UNAVAILABLE,
        ToolErrorType.TOOL_FAILED,
    }