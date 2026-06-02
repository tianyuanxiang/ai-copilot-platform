"""Runtime lifecycle for the Wind ReAct graph, gRPC client, and checkpoints."""

from __future__ import annotations

from contextlib import AbstractAsyncContextManager
from typing import Any

from langgraph.checkpoint.memory import InMemorySaver
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver

from app.clients.wind_tool_rpc import WindToolRPCClient
from app.core.config import Settings
from app.schemas.agent import AgentStreamEvent
from app.services import llm
from app.services.wind_agent_graph import build_answer_prompt, build_done_event, build_wind_agent_graph, initial_agent_state


class AgentRuntime:
    """Owns long-lived Agent infrastructure and exposes stream helpers."""

    def __init__(
        self,
        settings: Settings,
        *,
        tool_client: Any | None = None,
        checkpointer: Any | None = None,
        planner: Any = llm.plan_agent_action,
        answer_streamer: Any = llm.stream_agent_answer,
    ) -> None:
        self.settings = settings
        # 如果外部传了 tool_client，就用外部传进来的；如果没传，就默认创建一个 WindToolRPCClient。
        if tool_client is None:
            self.tool_client = WindToolRPCClient(
                target=settings.agent_go_rpc_target,
                timeout_seconds=settings.agent_go_rpc_timeout_seconds,
            )
        else:
            self.tool_client = tool_client
        self.checkpointer = checkpointer
        self.planner = planner
        self.answer_streamer = answer_streamer
        self.graph: Any | None = None
        self._checkpointer_context: AbstractAsyncContextManager[Any] | None = None

    async def start(self) -> None:
        """Connect checkpoint storage and compile the graph.

        PostgreSQL schema setup is intentionally not performed here. Production
        deployments run scripts/setup_agent_checkpoints.py once before startup.
        """

        if self.checkpointer is None:
            if self.settings.agent_checkpoint_backend == "memory":
                self.checkpointer = InMemorySaver()
            else:
                dsn = self.settings.agent_checkpoint_dsn or self.settings.postgres_dsn
                if not dsn:
                    raise RuntimeError("agent checkpoint DSN is empty")
                self._checkpointer_context = AsyncPostgresSaver.from_conn_string(dsn)
                self.checkpointer = await self._checkpointer_context.__aenter__()
        #搭建LangGraph 流程图。
        self.graph = build_wind_agent_graph(self.tool_client, checkpointer=self.checkpointer, planner=self.planner)

    async def close(self) -> None:
        """Release gRPC and PostgreSQL resources during FastAPI shutdown."""

        if hasattr(self.tool_client, "close"):
            await self.tool_client.close()
        if self._checkpointer_context is not None:
            await self._checkpointer_context.__aexit__(None, None, None)

    # 查询 checkpoint
    async def stream_new_turn(self, *, user_id: int, agent_session_id: str, user_input: str):
        """Start one Agent turn and stream custom events."""

        self._require_started()
        config = self._config(agent_session_id)
        snapshot = await self.graph.aget_state(config)   # 读取该会话以前保存的状态
        if _is_waiting(snapshot):
            yield AgentStreamEvent(
                type="error",
                agent_session_id=agent_session_id,
                error_msg="当前会话仍在等待澄清或人工确认，请调用 /v1/agent/resume/stream。",
            )
            return
        state = initial_agent_state(
            user_id=user_id,
            agent_session_id=agent_session_id,
            user_input=user_input,
            max_steps=self.settings.agent_max_steps,
        )
        async for event in self._stream_graph(state, config):
            yield event

    async def stream_resume(self, *, user_id: int, agent_session_id: str, action: str, content: str):
        """Resume the currently interrupted node for a conversation."""

        self._require_started()
        config = self._config(agent_session_id)
        snapshot = await self.graph.aget_state(config)
        values = snapshot.values or {}
        if not _is_waiting(snapshot):
            yield AgentStreamEvent(
                type="error",
                agent_session_id=agent_session_id,
                error_msg="当前会话没有等待恢复的 Agent 节点。",
            )
            return
        if int(values.get("user_id", 0)) != user_id:
            yield AgentStreamEvent(
                type="error",
                agent_session_id=agent_session_id,
                error_msg="当前用户无权恢复该会话。",
            )
            return
        resumed_state = {
            **values,
            "resume_action": {"action": action, "content": content},
            "events": [],
        }
        async for event in self._stream_graph(resumed_state, config):
            yield event

    # 真正运行 LangGraph
    async def _stream_graph(self, graph_input: Any, config: dict[str, Any]):

        # stream_mode="updates":每执行完一个节点，LangGraph 就返回该节点产生的新状态。
        async for update in self.graph.astream(graph_input, config=config, stream_mode="updates"):  # 正式启动状态机,执行节点 prepare_turn
            for node_update in update.values():
                if not isinstance(node_update, dict):
                    continue
                for raw_event in node_update.get("events", []):
                    yield AgentStreamEvent.model_validate(raw_event)
        snapshot = await self.graph.aget_state(config)   # 跑着跑着停了，主动去数据库/内存中抓取该任务当前的完整快照
        pending = _pending_interrupts(snapshot)
        for payload in pending:
            yield AgentStreamEvent.model_validate(payload)
        values = snapshot.values or {}
        if not pending and values.get("route") == "answer" and not values.get("answer"):
            parts: list[str] = []
            try:
                # 手动调用了一个外部的 LLM 流式接口
                async for token in self.answer_streamer(
                    build_answer_prompt(values),
                    user_id=values["user_id"],
                    conversation_id=values["agent_session_id"],
                    trace_id=values["trace_id"],
                ):
                    parts.append(token)
                    yield AgentStreamEvent(
                        type="token",
                        trace_id=values["trace_id"],
                        agent_session_id=values["agent_session_id"],
                        content=token,
                    )
            except Exception as exc:
                parts = [f"最终答案生成失败：{exc}"]
                yield AgentStreamEvent(
                    type="token",
                    trace_id=values["trace_id"],
                    agent_session_id=values["agent_session_id"],
                    content=parts[0],
                )
            answer = "".join(parts).strip()
            values = {**values, "answer": answer, "route": "done", "events": []}
            # The graph has already exited through the assess_evidence -> END
            # branch. Persist only the generated answer here: feeding a new
            # route value back into that completed conditional branch would
            # make LangGraph try to route "done" as if it were a graph edge.
            await self.graph.aupdate_state(config, {"answer": answer, "events": []})
            yield build_done_event(values)

    def _config(self, agent_session_id: str) -> dict[str, Any]:
        return {"configurable": {"thread_id": agent_session_id}, "recursion_limit": 48}

    def _require_started(self) -> None:
        if self.graph is None:
            raise RuntimeError("AgentRuntime.start() must be called before use")


def _pending_interrupts(snapshot: Any) -> list[dict[str, Any]]:
    payloads: list[dict[str, Any]] = []
    for task in getattr(snapshot, "tasks", ()) or ():
        for item in getattr(task, "interrupts", ()) or ():
            if isinstance(item.value, dict):
                payloads.append(item.value)
    return payloads


def _is_waiting(snapshot: Any) -> bool:
    values = snapshot.values or {}
    return values.get("route") in {"waiting_approval", "waiting_clarification"} or bool(_pending_interrupts(snapshot))
