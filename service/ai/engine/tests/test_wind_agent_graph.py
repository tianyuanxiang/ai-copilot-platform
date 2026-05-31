from collections.abc import Iterable

from langgraph.checkpoint.memory import InMemorySaver

from app.clients.wind_tool_rpc import WindToolExecuteResult
from app.schemas.agent import Citation, ToolCall
from app.services.wind_agent_graph import build_wind_agent_graph, initial_agent_state


class FakeToolClient:
    def __init__(self, *, status: str = "success", citations: list[Citation] | None = None):
        self.status = status
        self.citations = citations or []
        self.calls = []

    async def execute(self, request):
        self.calls.append(request)
        evidence_json = '{"source":"go","count":1}' if self.status == "success" else ""
        return WindToolExecuteResult(
            tool_call=ToolCall(
                tool_name=request.tool_name,
                status=self.status,
                arguments_json=request.arguments_json,
                result_json="{}",
                message="ok" if self.status == "success" else "failed",
            ),
            evidence_json=evidence_json,
            citations=self.citations,
        )


def scripted_planner(actions: Iterable[dict]):
    iterator = iter(actions)

    async def planner(*_):
        return next(iterator)

    return planner


async def collect(graph, graph_input, conversation_id: str):
    config = {"configurable": {"thread_id": conversation_id}}
    events = []
    async for update in graph.astream(graph_input, config=config, stream_mode="updates"):
        for node_update in update.values():
            events.extend(node_update.get("events", []))
    return events, await graph.aget_state(config)


def new_state(conversation_id: str, *, max_steps: int = 6):
    return initial_agent_state(
        user_id=123,
        conversation_id=conversation_id,
        user_input="分析 FY 风场告警并查询 SOP",
        max_steps=max_steps,
    )


async def test_read_only_sop_tool_calls_go_and_deduplicates_citations():
    citation = Citation(document_id=3, chunk_id=9, title="SOP")
    tool = FakeToolClient(citations=[citation, citation])
    graph = build_wind_agent_graph(
        tool,
        checkpointer=InMemorySaver(),
        planner=scripted_planner(
            [
                {"name": "search_maintenance_sop", "arguments": {"query": "齿轮箱油温"}},
                {"name": "finish_answer", "arguments": {"evidenceRequirement": "sop"}},
            ]
        ),
    )

    events, snapshot = await collect(graph, new_state("sop"), "sop")

    assert [item.type if hasattr(item, "type") else item["type"] for item in events] == [
        "start",
        "intent",
        "tool_start",
        "tool_result",
        "evidence",
        "intent",
    ]
    assert [item.tool_name for item in tool.calls] == ["search_maintenance_sop"]
    assert snapshot.values["route"] == "answer"
    assert len(snapshot.values["citations"]) == 1


async def test_draft_tool_waits_for_approval_then_executes_once():
    tool = FakeToolClient()
    graph = build_wind_agent_graph(
        tool,
        checkpointer=InMemorySaver(),
        planner=scripted_planner(
            [
                {"name": "create_maintenance_ticket_draft", "arguments": {"farmCode": "FY"}},
                {"name": "finish_answer", "arguments": {"evidenceRequirement": "operational"}},
            ]
        ),
    )

    events, snapshot = await collect(graph, new_state("approval"), "approval")

    assert events[-1]["type"] == "confirmation_required"
    assert snapshot.values["route"] == "waiting_approval"
    assert tool.calls == []

    resumed = {**snapshot.values, "resume_action": {"action": "approve", "content": ""}, "events": []}
    _, snapshot = await collect(graph, resumed, "approval")

    assert [item.tool_name for item in tool.calls] == ["create_maintenance_ticket_draft"]
    assert snapshot.values["route"] == "answer"


async def test_reject_does_not_execute_draft_tool():
    tool = FakeToolClient()
    graph = build_wind_agent_graph(
        tool,
        checkpointer=InMemorySaver(),
        planner=scripted_planner(
            [
                {"name": "generate_health_report", "arguments": {"farmCode": "FY"}},
                {"name": "finish_answer", "arguments": {"evidenceRequirement": "none"}},
            ]
        ),
    )

    _, snapshot = await collect(graph, new_state("reject"), "reject")
    resumed = {**snapshot.values, "resume_action": {"action": "reject", "content": ""}, "events": []}
    _, snapshot = await collect(graph, resumed, "reject")

    assert tool.calls == []
    assert snapshot.values["route"] == "answer"


async def test_clarification_resume_continues_planning():
    tool = FakeToolClient()
    graph = build_wind_agent_graph(
        tool,
        checkpointer=InMemorySaver(),
        planner=scripted_planner(
            [
                {"name": "request_clarification", "arguments": {"question": "请补充风场编码。"}},
                {"name": "finish_answer", "arguments": {"evidenceRequirement": "none"}},
            ]
        ),
    )

    events, snapshot = await collect(graph, new_state("clarify"), "clarify")

    assert events[-1]["type"] == "clarification_required"
    assert snapshot.values["route"] == "waiting_clarification"

    resumed = {**snapshot.values, "resume_action": {"action": "clarify", "content": "风场编码为 FY"}, "events": []}
    _, snapshot = await collect(graph, resumed, "clarify")

    assert snapshot.values["route"] == "answer"
    assert snapshot.values["messages"][-1]["content"] == "用户补充信息：风场编码为 FY"


async def test_same_tool_failure_twice_builds_degraded_answer():
    tool = FakeToolClient(status="failed")
    graph = build_wind_agent_graph(
        tool,
        checkpointer=InMemorySaver(),
        planner=scripted_planner(
            [
                {"name": "query_alarm_events", "arguments": {"farmCode": "FY"}},
                {"name": "query_alarm_events", "arguments": {"farmCode": "FY"}},
            ]
        ),
    )

    events, snapshot = await collect(graph, new_state("failed"), "failed")

    assert len(tool.calls) == 2
    assert snapshot.values["route"] == "done"
    assert events[-1]["type"] == "done"


async def test_max_steps_builds_degraded_answer():
    tool = FakeToolClient()
    graph = build_wind_agent_graph(
        tool,
        checkpointer=InMemorySaver(),
        planner=scripted_planner([{"name": "query_alarm_events", "arguments": {"farmCode": "FY"}}]),
    )

    _, snapshot = await collect(graph, new_state("max-steps", max_steps=1), "max-steps")

    assert len(tool.calls) == 1
    assert snapshot.values["route"] == "done"

