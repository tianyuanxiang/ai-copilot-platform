import json

from fastapi.testclient import TestClient
from langgraph.checkpoint.memory import InMemorySaver

from app.core.config import Settings
from app.main import app
from app.services import agent_runtime
from app.services.agent_runtime import AgentRuntime


def parse_sse(body: str) -> list[dict]:
    events = []
    for block in body.split("\n\n"):
        if not block.startswith("data: "):
            continue
        payload = block.removeprefix("data: ")
        if payload != "[DONE]":
            events.append(json.loads(payload))
    return events


def test_agent_stream_returns_start_tokens_and_done():
    with TestClient(app) as client:
        with client.stream("POST", "/v1/agent/stream", json={"user_id": 123, "input": "hello"}) as response:
            body = "".join(response.iter_text())

    events = parse_sse(body)
    assert response.status_code == 200
    assert events[0]["type"] == "start"
    assert events[0]["agent_session_id"]
    assert any(item["type"] == "token" for item in events)
    assert events[-1]["type"] == "done"
    assert "data: [DONE]" in body


def test_agent_run_aggregates_same_graph():
    with TestClient(app) as client:
        response = client.post("/v1/agent/run", json={"user_id": 123, "input": "hello"})

    assert response.status_code == 200
    assert response.json()["status"] == "done"
    assert response.json()["agent_session_id"]
    assert response.json()["answer"]


def test_agent_resume_without_waiting_node_returns_error_event():
    with TestClient(app) as client:
        with client.stream(
            "POST",
            "/v1/agent/resume/stream",
            json={"user_id": 123, "agent_session_id": "missing", "action": "approve"},
        ) as response:
            body = "".join(response.iter_text())

    events = parse_sse(body)
    assert response.status_code == 200
    assert events[0]["type"] == "error"
    assert "data: [DONE]" in body


class FakeToolClient:
    async def close(self):
        return None


class FakePostgresContext:
    def __init__(self, dsn: str):
        self.dsn = dsn
        self.checkpointer = InMemorySaver()
        self.closed = False

    async def __aenter__(self):
        return self.checkpointer

    async def __aexit__(self, *_):
        self.closed = True


async def test_postgres_runtime_uses_configured_checkpoint_dsn_without_startup_ddl(monkeypatch):
    contexts = []

    def fake_from_conn_string(dsn: str):
        context = FakePostgresContext(dsn)
        contexts.append(context)
        return context

    monkeypatch.setattr(agent_runtime.AsyncPostgresSaver, "from_conn_string", fake_from_conn_string)
    runtime = AgentRuntime(
        Settings(
            agent_checkpoint_backend="postgres",
            agent_checkpoint_dsn="postgresql://checkpoint.example/wind_agent",
        ),
        tool_client=FakeToolClient(),
    )

    await runtime.start()
    await runtime.close()

    assert contexts[0].dsn == "postgresql://checkpoint.example/wind_agent"
    assert contexts[0].closed is True
