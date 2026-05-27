from fastapi.testclient import TestClient

from app.core.config import get_settings
from app.main import app
from app.schemas.chat import ChatStreamRequest
from app.services import llm


def test_chat_stream_returns_sse_tokens():
    client = TestClient(app)
    with client.stream("POST", "/v1/chat/stream", json={"question": "hello"}) as response:
        assert response.status_code == 200
        body = "".join(response.iter_text())

    assert "data:" in body
    assert "[DONE]" in body
    assert '"content": "ASGI "' in body
    assert '"content": "streaming "' in body
    assert '"content": "ready. "' in body


def test_chat_stream_returns_error_when_deepseek_key_is_empty(tmp_path, monkeypatch):
    config_path = tmp_path / "config.yaml"
    config_path.write_text(
        """
app:
  name: ai-copilot-engine-test
  env: dev
llm:
  mock: false
  provider: deepseek
  model: deepseek-v4-pro
  deepseek_chat_url: "https://api.deepseek.com/chat/completions"
  api_key: ""
  timeout_seconds: 60
postgres:
  dsn: ""
elasticsearch:
  url: ""
  kb_chunks_index: kb_chunks_index
  security_logs_index: security_logs_index
embedding:
  provider: mock
  model: text-embedding-v4
  api_key: ""
  dimension: 1024
  batch_size: 10
  timeout_seconds: 30
rerank:
  provider: mock
  model: qwen3-rerank
  api_key: ""
  top_n: 5
  timeout_seconds: 30
""",
        encoding="utf-8",
    )
    monkeypatch.setenv("APP_CONFIG", str(config_path))
    get_settings.cache_clear()

    try:
        client = TestClient(app)
        with client.stream("POST", "/v1/chat/stream", json={"question": "hello"}) as response:
            assert response.status_code == 200
            body = "".join(response.iter_text())
    finally:
        get_settings.cache_clear()

    assert "data:" in body
    assert '"type": "error"' in body
    assert '"content": "llm api key is empty"' in body
    assert "[DONE]" in body


async def test_stream_chat_parses_deepseek_sse(tmp_path, monkeypatch):
    config_path = tmp_path / "config.yaml"
    config_path.write_text(
        """
app:
  name: ai-copilot-engine-test
  env: dev
llm:
  mock: false
  provider: deepseek
  model: deepseek-v4-pro
  deepseek_chat_url: "https://api.deepseek.com/chat/completions"
  api_key: "test-key"
  timeout_seconds: 60
postgres:
  dsn: ""
elasticsearch:
  url: ""
  kb_chunks_index: kb_chunks_index
  security_logs_index: security_logs_index
embedding:
  provider: mock
  model: text-embedding-v4
  api_key: ""
  dimension: 1024
  batch_size: 10
  timeout_seconds: 30
rerank:
  provider: mock
  model: qwen3-rerank
  api_key: ""
  top_n: 5
  timeout_seconds: 30
""",
        encoding="utf-8",
    )
    monkeypatch.setenv("APP_CONFIG", str(config_path))
    get_settings.cache_clear()

    class FakeResponse:
        status_code = 200

        async def __aenter__(self):
            return self

        async def __aexit__(self, exc_type, exc, tb):
            return False

        async def aiter_lines(self):
            yield ""
            yield 'data: {"choices":[{"delta":{"content":"你"}}]}'
            yield 'data: {"choices":[{"delta":{"content":"好"}}]}'
            yield "data: [DONE]"

    class FakeAsyncClient:
        def __init__(self, timeout):
            self.timeout = timeout

        async def __aenter__(self):
            return self

        async def __aexit__(self, exc_type, exc, tb):
            return False

        def stream(self, method, url, headers, json):
            assert method == "POST"
            assert url == "https://api.deepseek.com/chat/completions"
            assert headers["Authorization"] == "Bearer test-key"
            assert json["model"] == "deepseek-v4-pro"
            assert json["stream"] is True
            assert json["messages"][-1] == {"role": "user", "content": "hello"}
            return FakeResponse()

    monkeypatch.setattr(llm.httpx, "AsyncClient", FakeAsyncClient)

    try:
        events = [
            event
            async for event in llm.stream_chat(
                ChatStreamRequest(question="hello"),
                trace_id="trace-1",
            )
        ]
    finally:
        get_settings.cache_clear()

    assert [(event.type, event.content, event.trace_id) for event in events] == [
        ("token", "你", "trace-1"),
        ("token", "好", "trace-1"),
        ("done", "", "trace-1"),
    ]
