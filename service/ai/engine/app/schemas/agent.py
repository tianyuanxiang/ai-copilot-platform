"""Public request, response, and SSE models for the Wind ReAct Agent."""

from typing import Literal

from pydantic import BaseModel, Field, model_validator


class Citation(BaseModel):
    """A knowledge-base citation returned by the controlled Go RAG tool."""

    document_id: int = 0
    chunk_id: int = 0
    title: str = ""
    snippet: str = ""
    score: float = 0.0


class ToolCall(BaseModel):
    """A single controlled Go tool invocation visible to callers."""

    tool_call_id: int = 0
    tool_name: str
    status: str
    arguments_json: str = "{}"
    result_json: str = "{}"
    message: str = ""
    latency_ms: int = 0


class WindDraftRef(BaseModel):
    """A lightweight reference to a generated draft stored by Go."""

    draft_type: str = ""
    draft_id: int = 0
    title: str = ""


class AgentRunRequest(BaseModel):
    """Starts a new Agent turn.

    conversation_id may be omitted for the first turn. The engine returns the
    generated id in the first SSE event and reuses it for resume requests.
    """

    user_id: int = Field(gt=0)
    conversation_id: str = ""
    input: str = Field(min_length=1, max_length=10000)


class AgentResumeRequest(BaseModel):
    """Resumes an interrupted clarification or draft-approval node."""

    user_id: int = Field(gt=0)
    conversation_id: str = Field(min_length=1, max_length=128)
    action: Literal["approve", "reject", "clarify"]
    content: str = Field(default="", max_length=10000)

    @model_validator(mode="after")
    def validate_clarification_content(self) -> "AgentResumeRequest":
        if self.action == "clarify" and not self.content.strip():
            raise ValueError("content is required when action=clarify")
        return self


class AgentStreamEvent(BaseModel):
    """SSE payload aligned with Go WindAgentStreamEvent where possible."""

    type: str
    trace_id: str = ""
    conversation_id: str = ""
    content: str = ""
    tool_call: ToolCall | None = None
    tool_calls: list[ToolCall] = Field(default_factory=list)
    citations: list[Citation] = Field(default_factory=list)
    draft: WindDraftRef | None = None
    error_msg: str = ""


class AgentRunResponse(BaseModel):
    """Aggregated non-streaming result for diagnostics and internal callers."""

    answer: str
    tool_calls: list[ToolCall] = Field(default_factory=list)
    trace_id: str
    conversation_id: str
    citations: list[Citation] = Field(default_factory=list)
    draft: WindDraftRef | None = None
    status: str = "done"
