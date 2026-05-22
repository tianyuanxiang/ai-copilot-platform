from fastapi import APIRouter

from app.schemas.security import (
    SecurityAnalyzeRequest,
    SecurityAnalyzeResponse,
    SecurityLogSearchRequest,
    SecurityLogSearchResponse,
)
from app.services import security

router = APIRouter(tags=["security"])


@router.post("/logs/search", response_model=SecurityLogSearchResponse)
async def search_logs(payload: SecurityLogSearchRequest) -> SecurityLogSearchResponse:
    return await security.search_logs(payload.query, payload.ip, payload.username, payload.page, payload.page_size)


@router.post("/events/analyze", response_model=SecurityAnalyzeResponse)
async def analyze_events(payload: SecurityAnalyzeRequest) -> SecurityAnalyzeResponse:
    return await security.analyze_events(payload.raw_logs)