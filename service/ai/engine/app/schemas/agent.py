from pydantic import BaseModel, Field


class AgentRunRequest(BaseModel):
    user_id: str = "demo-user"
    conversation_id: str = ""
    input: str = Field(min_length=1, max_length=4000)


class ToolCall(BaseModel):
    name: str
    status: str
    message: str


class AgentRunResponse(BaseModel):
    answer: str
    tool_calls: list[ToolCall]
    trace_id: str