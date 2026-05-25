from typing import Any

from pydantic import BaseModel, Field


class WindEvidenceRequest(BaseModel):
    user_id: str = "demo-user"
    farm_code: str = Field(default="", max_length=64)
    tower_code: str = Field(default="", max_length=64)
    question: str = Field(default="", max_length=4000)
    evidence: list[dict[str, Any]] = Field(default_factory=list)
    evidence_json: str = ""


class WindHealthReportDraftRequest(WindEvidenceRequest):
    report_type: str = "single_turbine_health"
    start_time: str = ""
    end_time: str = ""


class WindTicketDraftRequest(WindEvidenceRequest):
    alarm_code: str = ""
    priority: str = "normal"


class WindScaffoldResponse(BaseModel):
    title: str
    status: str
    summary: str
    evidence_count: int
    message: str
    metrics: dict[str, Any] = Field(default_factory=dict)
    sections: list[dict[str, Any]] = Field(default_factory=list)
    recommendations: list[str] = Field(default_factory=list)
    todo: list[str]
