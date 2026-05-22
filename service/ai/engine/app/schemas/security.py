from pydantic import BaseModel, Field


class SecurityLogSearchRequest(BaseModel):
    user_id: str = "demo-user"
    query: str = ""
    ip: str = ""
    username: str = ""
    start_time: str = ""
    end_time: str = ""
    page: int = 1
    page_size: int = 20


class SecurityLogItem(BaseModel):
    id: str
    timestamp: str = ""
    host: str = ""
    ip: str = ""
    username: str = ""
    message: str


class SecurityLogSearchResponse(BaseModel):
    total: int
    list: list[SecurityLogItem]
    mode: str
    message: str


class SecurityAnalyzeRequest(BaseModel):
    user_id: str = "demo-user"
    raw_logs: list[str] = Field(default_factory=list)


class SecurityEvent(BaseModel):
    event_type: str
    severity: str
    ip: str = ""
    username: str = ""
    summary: str


class SecurityAnalyzeResponse(BaseModel):
    events: list[SecurityEvent]
    alerts: list[str]
    message: str