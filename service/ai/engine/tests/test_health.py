from fastapi.testclient import TestClient

from app.main import app


def test_health_marks_asgi_streaming():
    client = TestClient(app)
    response = client.get("/v1/health")

    assert response.status_code == 200
    assert response.json()["asgi"] is True
    assert response.json()["streaming"] == "sse"