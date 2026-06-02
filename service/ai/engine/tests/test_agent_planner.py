import httpx
import pytest

from app.core.config import Settings
from app.services import llm


async def test_plan_agent_action_disables_thinking_when_tool_choice_is_required(monkeypatch):
    class FakeResponse:
        def raise_for_status(self):
            return None

        def json(self):
            return {
                "choices": [
                    {
                        "message": {
                            "tool_calls": [
                                {
                                    "function": {
                                        "name": "finish_answer",
                                        "arguments": '{"evidenceRequirement":"none"}',
                                    }
                                }
                            ]
                        }
                    }
                ]
            }

    class FakeAsyncClient:
        def __init__(self, timeout):
            self.timeout = timeout

        async def __aenter__(self):
            return self

        async def __aexit__(self, *_):
            return None

        async def post(self, url, headers, json):
            assert json["tool_choice"] == "required"
            assert json["thinking"] == {"type": "disabled"}
            return FakeResponse()

    monkeypatch.setattr(
        llm,
        "get_settings",
        lambda: Settings(mock_llm=False, llm_provider="deepseek", llm_api_key="test-key"),
    )
    monkeypatch.setattr(llm.httpx, "AsyncClient", FakeAsyncClient)

    action = await llm.plan_agent_action([{"role": "user", "content": "test"}], [], "trace-1")

    assert action == {"name": "finish_answer", "arguments": {"evidenceRequirement": "none"}}


async def test_plan_agent_action_reports_provider_error_body(monkeypatch):
    class FakeAsyncClient:
        def __init__(self, timeout):
            self.timeout = timeout

        async def __aenter__(self):
            return self

        async def __aexit__(self, *_):
            return None

        async def post(self, url, headers, json):
            request = httpx.Request("POST", url)
            return httpx.Response(400, request=request, text='{"error":{"message":"invalid tool_choice"}}')

    monkeypatch.setattr(
        llm,
        "get_settings",
        lambda: Settings(mock_llm=False, llm_provider="deepseek", llm_api_key="test-key"),
    )
    monkeypatch.setattr(llm.httpx, "AsyncClient", FakeAsyncClient)

    with pytest.raises(RuntimeError, match="invalid tool_choice"):
        await llm.plan_agent_action([{"role": "user", "content": "test"}], [], "trace-1")
