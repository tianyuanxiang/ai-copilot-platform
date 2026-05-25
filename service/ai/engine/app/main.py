from fastapi import FastAPI

from app.api.routes_agent import router as agent_router
from app.api.routes_chat import router as chat_router
from app.api.routes_health import router as health_router
from app.api.routes_knowledge import router as knowledge_router
from app.api.routes_wind import router as wind_router
from app.core.config import get_settings


def create_app() -> FastAPI:
    settings = get_settings()
    app = FastAPI(title=settings.app_name, version="0.1.0")
    app.include_router(health_router, prefix="/v1")
    app.include_router(knowledge_router, prefix="/v1")
    app.include_router(chat_router, prefix="/v1")
    app.include_router(wind_router, prefix="/v1/wind")
    app.include_router(agent_router, prefix="/v1")
    return app


app = create_app()
