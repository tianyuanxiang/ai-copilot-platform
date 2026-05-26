import json

import httpx
import pytest
from fastapi.testclient import TestClient

from app.core.config import Settings, load_settings
from app.main import app
from app.services import rag


def test_mock_embed_returns_vectors_and_token_usage():
    client = TestClient(app)

    response = client.post("/v1/embed", json={"texts": ["hello world", "中文 文本"]})

    assert response.status_code == 200
    body = response.json()
    assert body["mode"] == "mock-embedding"
    assert body["model"] == "text-embedding-v4"
    assert body["dimension"] == 1024
    assert len(body["vectors"]) == 2
    assert all(len(vector) == 1024 for vector in body["vectors"])
    assert body["token_counts"] == [2, 2]
    assert body["total_tokens"] == 4


def test_embed_rejects_blank_text():
    client = TestClient(app)

    response = client.post("/v1/embed", json={"texts": ["valid", "   "]})

    assert response.status_code == 400
    assert "texts[1]" in response.json()["detail"]


@pytest.mark.asyncio
async def test_dashscope_embed_batches_requests_and_reads_usage():
    seen_batches: list[list[str]] = []

    def handler(request: httpx.Request) -> httpx.Response:
        assert request.url == "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings"
        assert request.headers["authorization"] == "Bearer test-key"
        payload = request.read()
        body = json.loads(payload)
        seen_batches.append(body["input"])
        assert body["model"] == "text-embedding-v4"
        assert body["dimensions"] == 1024
        assert body["encoding_format"] == "float"
        assert "text_type" not in body

        offset = sum(len(batch) for batch in seen_batches[:-1])
        data = [
            {
                "object": "embedding",
                "index": offset + index,
                "embedding": [float(offset + index)] * 1024,
            }
            for index, _ in enumerate(body["input"])
        ]
        return httpx.Response(
            200,
            json={
                "object": "list",
                "data": data,
                "model": "text-embedding-v4",
                "usage": {
                    "prompt_tokens": len(body["input"]) * 3,
                    "total_tokens": len(body["input"]) * 3,
                },
            },
        )

    settings = Settings(
        embedding_provider="dashscope",
        embedding_api_key="test-key",
        embedding_model="text-embedding-v4",
        embedding_dimension=1024,
        embedding_batch_size=10,
    )
    transport = httpx.MockTransport(handler)

    result = await rag.embed_texts(
        [f"text {index}" for index in range(11)],
        input_type="query",
        settings=settings,
        transport=transport,
    )

    assert [len(batch) for batch in seen_batches] == [10, 1]
    assert result.mode == "dashscope-openai-compatible-embedding"
    assert result.model == "text-embedding-v4"
    assert result.dimension == 1024
    assert len(result.vectors) == 11
    assert result.token_counts == [3] * 11
    assert result.total_tokens == 33


def test_embedding_config_loads_yaml_section(tmp_path):
    config_file = tmp_path / "config.yaml"
    config_file.write_text(
        """
embedding:
  provider: dashscope
  model: text-embedding-v4
  api_key: config-key
  dimension: 768
  batch_size: 5
  timeout_seconds: 12
""",
        encoding="utf-8",
    )

    settings = load_settings(config_file)

    assert settings.embedding_provider == "dashscope"
    assert settings.embedding_api_key == "config-key"
    assert settings.embedding_dimension == 768
    assert settings.embedding_batch_size == 5
    assert settings.embedding_timeout_seconds == 12
