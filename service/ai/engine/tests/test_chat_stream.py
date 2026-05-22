from fastapi.testclient import TestClient

from app.main import app


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
