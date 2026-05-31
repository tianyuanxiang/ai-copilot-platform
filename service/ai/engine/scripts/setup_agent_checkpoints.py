"""Create LangGraph PostgreSQL checkpoint tables for the Wind Agent.

This is an explicit deployment step. The FastAPI process intentionally does
not run DDL during startup.
"""

import asyncio

from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver

from app.core.config import get_settings


async def setup() -> None:
    settings = get_settings()
    dsn = settings.agent_checkpoint_dsn or settings.postgres_dsn
    if not dsn:
        raise RuntimeError("agent checkpoint DSN is empty")
    async with AsyncPostgresSaver.from_conn_string(dsn) as checkpointer:
        await checkpointer.setup()


if __name__ == "__main__":
    asyncio.run(setup())
