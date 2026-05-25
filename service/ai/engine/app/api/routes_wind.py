from fastapi import APIRouter

from app.schemas.wind import (
    WindEvidenceRequest,
    WindHealthReportDraftRequest,
    WindScaffoldResponse,
    WindTicketDraftRequest,
)
from app.services import wind

router = APIRouter(tags=["wind"])


@router.post("/summary/timeseries", response_model=WindScaffoldResponse)
async def wind_timeseries_summary(payload: WindEvidenceRequest) -> WindScaffoldResponse:
    return await wind.summarize_timeseries(payload)


@router.post("/summary/alarm", response_model=WindScaffoldResponse)
async def wind_alarm_summary(payload: WindEvidenceRequest) -> WindScaffoldResponse:
    return await wind.summarize_alarm(payload)


@router.post("/reports/health/draft", response_model=WindScaffoldResponse)
async def wind_health_report_draft(payload: WindHealthReportDraftRequest) -> WindScaffoldResponse:
    return await wind.draft_health_report(payload)


@router.post("/tickets/draft", response_model=WindScaffoldResponse)
async def wind_ticket_draft(payload: WindTicketDraftRequest) -> WindScaffoldResponse:
    return await wind.draft_ticket(payload)
