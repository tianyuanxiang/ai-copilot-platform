"""LangGraph ReAct state machine for the Wind O&M Copilot.

Python decides which controlled tool to call and whether evidence is sufficient.
All operational facts, RAG retrieval, persistence, and audit logging stay in Go.
"""

from __future__ import annotations

import json
import logging
import uuid
from collections.abc import Awaitable, Callable
from typing import Any, Protocol, TypedDict

from langgraph.graph import END, START, StateGraph

from app.clients.wind_tool_rpc import WindToolExecuteRequest, WindToolExecuteResult
from app.schemas.agent import AgentStreamEvent, Citation, ToolCall, WindDraftRef
from app.services import llm
from app.services.wind_agent_tools import (
    AGENT_TOOL_SCHEMAS,
    ALL_PLANNER_ACTIONS,
    DRAFT_TOOLS,
    GO_TOOLS,
    PLANNER_SYSTEM_PROMPT,
)


logger = logging.getLogger(__name__)

MAX_TOOL_RESULT_CHARS = 20000
MAX_EVIDENCE_ITEMS = 12
MAX_CITATIONS = 20


class ToolExecutor(Protocol):
    """Minimal protocol implemented by the async Go RPC client and test fakes."""
    async def execute(self, request: WindToolExecuteRequest) -> WindToolExecuteResult: ...


Planner = Callable[[list[dict[str, str]], list[dict[str, Any]], str], Awaitable[dict[str, Any]]]
class WindAgentState(TypedDict, total=False):
    """Serializable state persisted by the LangGraph checkpointer."""

    user_id: int
    agent_session_id: str
    trace_id: str
    input: str
    messages: list[dict[str, str]]
    step: int
    max_steps: int
    next_action: dict[str, Any]
    route: str
    tool_calls: list[dict[str, Any]]
    evidence: list[Any]
    citations: list[dict[str, Any]]
    draft: dict[str, Any]
    failures: dict[str, int]
    answer: str
    events: list[dict[str, Any]]
    resume_action: dict[str, str]


def initial_agent_state(
    *,
    user_id: int,
    agent_session_id: str,
    user_input: str,
    max_steps: int,
) -> WindAgentState:
    """Create a fresh turn while keeping agent_session_id stable across turns."""

    return {
        "user_id": user_id,
        "agent_session_id": agent_session_id,
        "trace_id": f"wind-agent-{uuid.uuid4()}",
        "input": user_input.strip(),
        "messages": [
            {"role": "system", "content": PLANNER_SYSTEM_PROMPT},
            {"role": "user", "content": user_input.strip()},
        ],
        "step": 0,
        "max_steps": max_steps,
        "next_action": {},
        "route": "",
        "tool_calls": [],
        "evidence": [],
        "citations": [],
        "draft": {},
        "failures": {},
        "answer": "",
        "events": [],
        "resume_action": {},
    }


def build_wind_agent_graph(
    tool_executor: ToolExecutor,
    *,
    checkpointer: Any,
    planner: Planner = llm.plan_agent_action,
):
    """Compile the Wind ReAct graph with injected infrastructure dependencies."""

    async def prepare_turn(state: WindAgentState) -> WindAgentState:
        # 会判断是不是恢复之前等待确认的状态。
        if state.get("resume_action"):
            return _resume_waiting_state(state, state["resume_action"])
        return {**state, "route": "plan", "events": [_event_payload(state, "start", content="Wind Agent 已开始分析。")]}

    async def plan_next_action(state: WindAgentState) -> WindAgentState:
        if state.get("step", 0) >= state.get("max_steps", 6):
            return {**state, "route": "degraded"}
        try:
            # 让大模型根据当前对话和工具列表，决定下一步要干什么。
            action = await planner(state.get("messages", []), AGENT_TOOL_SCHEMAS, state["trace_id"])
        except Exception as exc:
            logger.exception("wind agent planner failed trace_id=%s", state.get("trace_id"))
            return _with_failure(state, f"planner:{exc}", route="degraded")

        name = str(action.get("name") or "").strip()
        arguments = action.get("arguments")
        if name not in ALL_PLANNER_ACTIONS or not isinstance(arguments, dict):
            return _with_failure(state, f"planner returned unsupported action: {name}", route="degraded")

        content = "正在整理最终回答。" if name == "finish_answer" else f"Agent 选择下一步动作：{name}"
        route = _route_planned_action(name)
        return {
            **state,
            "next_action": {"name": name, "arguments": arguments},
            "route": route,
            "events": [_event_payload(state, "intent", content=content)],
        }

    def approval_gate(state: WindAgentState) -> WindAgentState:
        """Persist a confirmation wait state before any draft tool is executed.

        LangGraph's native async interrupt() relies on Python 3.11 context
        propagation. This service still supports Python 3.10, so the graph
        stores an equivalent durable wait state in the configured checkpointer.
        """

        action = state["next_action"]
        pending = ToolCall(
            tool_name=action["name"],
            status="confirmation_required",
            arguments_json=_json_dumps(action.get("arguments", {})),
            message="草稿工具将在确认后执行。",
        )
        return {
            **state,
            "route": "waiting_approval",
            "events": [
                _event_payload(
                    state,
                    "confirmation_required",
                    content="该操作会生成并保存草稿，请确认是否继续。",
                    tool_call=pending,
                )
            ],
        }

    def clarification_gate(state: WindAgentState) -> WindAgentState:
        action = state["next_action"]
        question = str(action.get("arguments", {}).get("question") or "请补充继续分析所需的信息。")
        return {
            **state,
            "route": "waiting_clarification",
            "events": [_event_payload(state, "clarification_required", content=question)],
        }

    async def announce_tool_start(state: WindAgentState) -> WindAgentState:
        action = state["next_action"]
        pending = ToolCall(
            tool_name=action["name"],
            status="running",
            arguments_json=_json_dumps(action.get("arguments", {})),   # 构造请求方法的参数
        )
        return {**state, "events": [_event_payload(state, "tool_start", tool_call=pending)]}

    async def execute_tool(state: WindAgentState) -> WindAgentState:
        action = state["next_action"]
        tool_name = action["name"]
        arguments_json = _json_dumps(action.get("arguments", {}))
        step = state.get("step", 0) + 1
        try:
            result = await tool_executor.execute(
                WindToolExecuteRequest(
                    user_id=state["user_id"],
                    trace_id=state["trace_id"],
                    conversation_id=state["agent_session_id"],
                    tool_name=tool_name,
                    arguments_json=arguments_json,
                    step=step,
                )
            )
            tool_call = result.tool_call.model_copy(update={"result_json": _clip_json(result.tool_call.result_json)})
        except Exception as exc:   # 如果工具调用失败，构造失败结果
            logger.exception("Go tool execution failed trace_id=%s tool=%s", state.get("trace_id"), tool_name)
            result = WindToolExecuteResult(
                tool_call=ToolCall(
                    tool_name=tool_name,
                    status="failed",
                    arguments_json=arguments_json,
                    message=str(exc),
                ),
                evidence_json="",
                citations=[],
            )
            tool_call = result.tool_call

        # 把这次工具调用追加到历史记录里
        tool_calls = [*state.get("tool_calls", []), tool_call.model_dump()]
        evidence = _append_evidence(state.get("evidence", []), result.evidence_json) # 把工具返回的 evidence_json 追加到证据列表
        citations = _merge_citations(state.get("citations", []), result.citations)
        draft = _draft_from_result(tool_call.result_json) or state.get("draft", {}) # 如果本次工具结果里包含草稿信息，就取出来放到 state["draft"]；否则沿用原来的 draft。
        messages = _append_message(   # 把工具结果写进对话上下文，让下一次 planner 能看到。
            state,
            "assistant",
            f"工具结果 {tool_name}: status={tool_call.status}; message={tool_call.message}; "
            f"result={tool_call.result_json}",
        )
        events = [_event_payload(state, "tool_result", tool_call=tool_call)] # 生成事件，给前端流式展示
        if result.evidence_json:
            events.append(_event_payload(state, "evidence", content=_clip_json(result.evidence_json)))

        failure_key = f"{tool_name}:{arguments_json}"
        failures = dict(state.get("failures", {}))
        if tool_call.status == "success":
            failures.pop(failure_key, None)
        else:
            failures[failure_key] = failures.get(failure_key, 0) + 1

        route = "plan"
        # 如果工具被拒绝执行、同一个工具同一组参数失败 2 次、工具调用步数超过最大限制
        if tool_call.status == "denied" or failures.get(failure_key, 0) >= 2 or step >= state.get("max_steps", 6):
            route = "degraded"
        return {
            **state,
            "step": step,
            "route": route,
            "tool_calls": tool_calls,
            "evidence": evidence,
            "citations": citations,
            "draft": draft,
            "messages": messages,
            "failures": failures,
            "events": events,
        }

    async def assess_evidence(state: WindAgentState) -> WindAgentState:
        arguments = state.get("next_action", {}).get("arguments", {})
        requirement = str(arguments.get("evidenceRequirement") or "none")
        if requirement == "none":
            return {**state, "route": "answer", "events": []}
        if requirement == "sop" and state.get("citations"):
            return {**state, "route": "answer", "events": []}
        if requirement == "operational" and state.get("evidence"):
            return {**state, "route": "answer", "events": []}

        question = "当前证据不足。请补充风场、风机、设备类型或希望分析的时间范围。"
        return {
            **state,
            "next_action": {
                "name": "request_clarification",
                "arguments": {"question": question, "reason": f"missing {requirement} evidence"},
            },
            "route": "clarify",
            "events": [],
        }

    async def build_degraded_answer(state: WindAgentState) -> WindAgentState:
        failures = list(state.get("failures", {}).keys())
        answer = "当前无法完成完整分析。"
        if state.get("evidence"):
            answer += " 已保留可用证据，请结合工具结果人工复核。"
        if failures:
            answer += f" 失败原因：{failures[-1]}"
        finished = {**state, "answer": answer}
        return {
            **finished,
            "route": "done",
            "events": [
                _event_payload(state, "token", content=answer),
                build_done_event(finished).model_dump(exclude_none=True),
            ],
        }

    graph = StateGraph(WindAgentState)
    graph.add_node("prepare_turn", prepare_turn)
    graph.add_node("plan_next_action", plan_next_action)
    graph.add_node("approval_gate", approval_gate)
    graph.add_node("clarification_gate", clarification_gate)
    graph.add_node("announce_tool_start", announce_tool_start)
    graph.add_node("execute_tool", execute_tool)
    graph.add_node("assess_evidence", assess_evidence)
    graph.add_node("build_degraded_answer", build_degraded_answer)

    graph.add_edge(START, "prepare_turn")
    graph.add_conditional_edges(
        "prepare_turn",
        lambda state: state["route"],
        {
            "plan": "plan_next_action",
            "execute": "announce_tool_start",
            "degraded": "build_degraded_answer",
            "waiting_clarification": END,
        },
    )
    graph.add_conditional_edges(
        "plan_next_action",
        lambda state: state["route"],  # 看 state["route"] 的值。如果是 execute，就去执行工具；如果是 approval，就等待用户确认；如果是 clarify，就让用户补充信息；如果是 degraded，就生成降级回答。
        {
            "execute": "announce_tool_start",
            "approval": "approval_gate",
            "clarify": "clarification_gate",    #   需要用户补充信息
            "assess": "assess_evidence",
            "degraded": "build_degraded_answer",
        },
    )
    graph.add_conditional_edges(
        "approval_gate",
        lambda state: state["route"],
        {"waiting_approval": END},
    )
    graph.add_edge("announce_tool_start", "execute_tool")
    # 需要用户确认时暂停，返回一个事件confirmation_required，然后走到 END，等待用户确认。
    graph.add_conditional_edges(
        "clarification_gate",
        lambda state: state["route"],
        {"waiting_clarification": END},
    )
    graph.add_conditional_edges(
        "execute_tool",
        lambda state: state["route"],
        {"plan": "plan_next_action", "degraded": "build_degraded_answer"},
    )
    graph.add_conditional_edges(
        "assess_evidence",
        lambda state: state["route"],
        {"answer": END, "clarify": "clarification_gate"},
    )
    graph.add_edge("build_degraded_answer", END)

    return graph.compile(checkpointer=checkpointer)


def _route_planned_action(name: str) -> str:
    if name in DRAFT_TOOLS:
        return "approval"
    if name in GO_TOOLS:
        return "execute"
    if name == "request_clarification":
        return "clarify"
    return "assess"


def _resume_waiting_state(state: WindAgentState, resume_action: dict[str, str]) -> WindAgentState:
    # state["route"] = waiting_clarification
    action = str(resume_action.get("action") or "").strip()
    if state.get("route") == "waiting_approval":
        if action == "approve":
            return {**state, "route": "execute", "resume_action": {}, "events": []}

        if action == "reject":
            messages = _append_message(
                state,
                "user",
                "用户拒绝执行草稿工具，请不要创建草稿。"
            )
            question = (
                "已取消该草稿操作，未执行任何写入。"
                "你希望我接下来怎么处理？可以选择："
                "1）仅基于当前证据给出分析建议；"
                "2）重新补充条件后再分析；"
                "3）结束本次任务。"
            )
            return {
                **state,
                "messages": messages,
                "next_action": {
                    "name": "request_clarification",
                    "arguments": {
                        "question": question,
                        "reason": "user_rejected_draft_tool",
                    },
                },
                "route": "waiting_clarification",
                "resume_action": {},
                "events": [
                    _event_payload(
                        state,
                        "clarification_required",
                        content=question,
                    )
                ],
            }
    if state.get("route") == "waiting_clarification" and action == "clarify":
        content = str(resume_action.get("content") or "").strip()
        if content:
            messages = _append_message(state, "user", f"用户补充信息：{content}")
            return {**state, "messages": messages, "next_action": {}, "route": "plan", "resume_action": {},
                    "events": []}

    failed = _with_failure(state, "恢复动作与当前等待节点不匹配，无法继续分析。", route="degraded")
    return {**failed, "resume_action": {}}



def _append_message(state: WindAgentState, role: str, content: str) -> list[dict[str, str]]:
    return [*state.get("messages", []), {"role": role, "content": content[:MAX_TOOL_RESULT_CHARS]}]


def _with_failure(state: WindAgentState, message: str, *, route: str) -> WindAgentState:
    failures = dict(state.get("failures", {}))
    failures[message] = failures.get(message, 0) + 1
    return {**state, "failures": failures, "route": route, "events": []}


def _append_evidence(current: list[Any], raw: str) -> list[Any]:
    if not raw.strip():
        return current
    try:
        value = json.loads(raw)
    except json.JSONDecodeError:
        value = {"raw": _clip_json(raw)}
    return [*current, _compact_value(value)][-MAX_EVIDENCE_ITEMS:]


def _merge_citations(current: list[dict[str, Any]], additions: list[Citation]) -> list[dict[str, Any]]:
    merged = [*current]
    seen = {(item.get("document_id", 0), item.get("chunk_id", 0)) for item in merged}
    for citation in additions:
        key = (citation.document_id, citation.chunk_id)
        if key in seen:
            continue
        seen.add(key)
        merged.append(citation.model_dump())
    return merged[:MAX_CITATIONS]


def _draft_from_result(raw: str) -> dict[str, Any]:
    try:
        value = json.loads(raw)
    except json.JSONDecodeError:
        return {}
    if not isinstance(value, dict) or not value.get("draftId"):
        return {}
    return WindDraftRef(
        draft_type=str(value.get("draftType") or ""),
        draft_id=int(value.get("draftId") or 0),
        title=str(value.get("title") or ""),
    ).model_dump()


def _compact_value(value: Any) -> Any:
    serialized = _json_dumps(value)
    if len(serialized) <= MAX_TOOL_RESULT_CHARS:
        return value
    return {"truncated": True, "original_chars": len(serialized), "preview": serialized[:1000]}


def _clip_json(raw: str) -> str:
    if len(raw) <= MAX_TOOL_RESULT_CHARS:
        return raw
    return _json_dumps({"truncated": True, "original_chars": len(raw), "preview": raw[:1000]})


def _json_dumps(value: Any) -> str:
    return json.dumps(value, ensure_ascii=False, separators=(",", ":"))


def build_answer_prompt(state: WindAgentState) -> str:
    action = state.get("next_action", {}).get("arguments", {})
    context = {
        "question": state.get("input", ""),
        "answer_focus": action.get("answerFocus", ""),
        "tool_calls": state.get("tool_calls", []),
        "evidence": state.get("evidence", []),
        "citations": state.get("citations", []),
        "draft": state.get("draft", {}),
    }
    return (
        "你是风机混塔智能运维 AI Copilot。仅根据下面事实回答，不得编造原因、数据或 SOP。"
        "事实不足时明确说明。涉及现场动作时提醒人工复核。引用 SOP 时标明文档标题。\n\n"
        f"{_json_dumps(context)}"
    )


def _event_payload(
    state: WindAgentState,
    event_type: str,
    *,
    content: str = "",
    tool_call: ToolCall | None = None,
) -> dict[str, Any]:
    return AgentStreamEvent(
        type=event_type,
        trace_id=state.get("trace_id", ""),
        agent_session_id=state.get("agent_session_id", ""),
        content=content,
        tool_call=tool_call,
    ).model_dump(exclude_none=True)


def build_done_event(state: WindAgentState) -> AgentStreamEvent:
    """Build the public terminal event after answer generation or degradation."""

    return AgentStreamEvent(
        type="done",
        trace_id=state.get("trace_id", ""),
        agent_session_id=state.get("agent_session_id", ""),
        content=state.get("answer", ""),
        tool_calls=[ToolCall.model_validate(item) for item in state.get("tool_calls", [])],
        citations=[Citation.model_validate(item) for item in state.get("citations", [])],
        draft=WindDraftRef.model_validate(state["draft"]) if state.get("draft") else None,
    )
