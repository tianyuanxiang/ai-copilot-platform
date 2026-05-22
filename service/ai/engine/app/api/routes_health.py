from fastapi import APIRouter

from app.core.config import get_settings

router = APIRouter(tags=["health"])


@router.get("/health")

async def health() -> dict:
    settings = get_settings()
    return {
        "status": "ok",
        "mode": settings.app_env,
        "asgi": True,
        "streaming": "sse",
    }