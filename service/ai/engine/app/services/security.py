from app.core.config import get_settings
from app.schemas.security import (
    SecurityAnalyzeResponse,
    SecurityEvent,
    SecurityLogItem,
    SecurityLogSearchResponse,
)


async def search_logs(query: str, ip: str, username: str, page: int, page_size: int) -> SecurityLogSearchResponse:
    settings = get_settings()
    if not settings.elasticsearch_url:
        return SecurityLogSearchResponse(
            total=0,
            list=[],
            mode="degraded",
            message="ELASTICSEARCH_URL is empty; security log full-text search is reserved",
        )

    return SecurityLogSearchResponse(
        total=1,
        list=[
            SecurityLogItem(
                id="mock-log-1",
                timestamp="",
                host="demo-host",
                ip=ip,
                username=username,
                message=f"Mock ES hit for query: {query}",
            )
        ],
        mode="mock-es",
        message=f"Elasticsearch index reserved: {settings.security_logs_index}",
    )


async def analyze_events(raw_logs: list[str]) -> SecurityAnalyzeResponse:
    events: list[SecurityEvent] = []
    for raw in raw_logs:
        lower = raw.lower()
        if "failed password" in lower or "authentication failure" in lower:
            events.append(
                SecurityEvent(
                    event_type="ssh_login_failed",
                    severity="medium",
                    summary=raw,
                )
            )

    alerts: list[str] = []
    if len(events) >= 5:
        alerts.append("possible brute-force SSH login attempts")

    return SecurityAnalyzeResponse(
        events=events,
        alerts=alerts,
        message="rule skeleton completed; IP geo and threat intelligence are reserved",
    )