from fastapi import APIRouter

from app.schemas.agent import AgentRunRequest, AgentRunResponse
from app.services.agent import run_agent

router = APIRouter(tags=["agent"])


@router.post("/agent/run", response_model=AgentRunResponse)
async def agent_run(payload: AgentRunRequest) -> AgentRunResponse:
    return await run_agent(payload.user_id, payload.input)