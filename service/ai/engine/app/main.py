import logging
import time
import uuid
from contextlib import asynccontextmanager

from fastapi import FastAPI
from starlette.requests import Request

from app.api.routes_agent import router as agent_router
from app.api.routes_chat import router as chat_router
from app.api.routes_health import router as health_router
from app.api.routes_knowledge import router as knowledge_router
from app.api.routes_wind import router as wind_router
from app.core.config import get_settings
from app.core.logging import setup_logging
from app.services.agent_runtime import AgentRuntime


logger = logging.getLogger("app.request")


def create_app() -> FastAPI:
    settings = get_settings()
    setup_logging()

    @asynccontextmanager
    async def lifespan(app: FastAPI):
        runtime = AgentRuntime(settings)
        await runtime.start()
        app.state.agent_runtime = runtime
        try:
            yield
        finally:
            await runtime.close()

    app = FastAPI(title=settings.app_name, version="0.1.0", lifespan=lifespan)

    @app.middleware("http")
    async def request_logging_middleware(request: Request, call_next):
        request_id = request.headers.get("x-request-id") or str(uuid.uuid4())
        request.state.request_id = request_id
        started_at = time.perf_counter()
        logger.info(
            "request.start request_id=%s method=%s path=%s client=%s",
            request_id,
            request.method,
            request.url.path,
            request.client.host if request.client else "-",
        )
        try:
            response = await call_next(request)
        except Exception:
            duration_ms = (time.perf_counter() - started_at) * 1000
            logger.exception(
                "request.error request_id=%s method=%s path=%s duration_ms=%.2f",
                request_id,
                request.method,
                request.url.path,
                duration_ms,
            )
            raise
        duration_ms = (time.perf_counter() - started_at) * 1000
        response.headers["x-request-id"] = request_id
        logger.info(
            "request.end request_id=%s method=%s path=%s status=%s duration_ms=%.2f",
            request_id,
            request.method,
            request.url.path,
            response.status_code,
            duration_ms,
        )
        return response

    app.include_router(health_router, prefix="/v1")
    app.include_router(knowledge_router, prefix="/v1")
    app.include_router(chat_router, prefix="/v1")
    app.include_router(wind_router, prefix="/v1/wind")
    app.include_router(agent_router, prefix="/v1")
    return app


app = create_app()
