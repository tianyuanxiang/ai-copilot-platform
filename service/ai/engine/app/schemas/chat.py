from pydantic import BaseModel, Field


class ChatMessage(BaseModel):
    role: str
    content: str


class Citation(BaseModel):
    document_id: str = ""
    chunk_id: str = ""
    title: str = ""
    snippet: str = ""
    score: float = 0.0


class ChatStreamRequest(BaseModel):
    user_id: str = Field(default="demo-user")
    kb_id: str = ""
    conversation_id: str = ""
    question: str = Field(min_length=1, max_length=4000)
    history: list[ChatMessage] = Field(default_factory=list)


class StreamEvent(BaseModel):
    type: str
    content: str = ""
    trace_id: str
    citations: list[Citation] = Field(default_factory=list)