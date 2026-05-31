from app.clients.wind_tool_rpc import WindToolExecuteRequest, WindToolRPCClient
from app.generated import ai_pb2


class FakeChannel:
    async def close(self):
        return None


class FakeStub:
    def __init__(self):
        self.payload = None
        self.timeout = None

    async def ExecuteTool(self, payload, *, timeout):
        self.payload = payload
        self.timeout = timeout
        return ai_pb2.WindToolExecuteResp(
            tool_call_id=17,
            tool_name="search_maintenance_sop",
            status="success",
            result_json='{"chunks":[]}',
            evidence_json='{"chunks":[]}',
            message="ok",
            latency_ms=12,
            citations=[
                ai_pb2.Citation(
                    document_id=3,
                    chunk_id=9,
                    title="齿轮箱维护 SOP",
                    snippet="检查油温和振动趋势",
                    score=0.91,
                )
            ],
        )


async def test_execute_tool_maps_request_and_response():
    stub = FakeStub()
    client = WindToolRPCClient("unused:0", 8, channel=FakeChannel(), stub=stub)

    result = await client.execute(
        WindToolExecuteRequest(
            user_id=123,
            trace_id="trace-1",
            conversation_id="conversation-1",
            tool_name="search_maintenance_sop",
            arguments_json='{"query":"齿轮箱油温"}',
            step=2,
        )
    )

    assert stub.timeout == 8
    assert stub.payload.user_id == 123
    assert stub.payload.trace_id == "trace-1"
    assert stub.payload.conversation_id == "conversation-1"
    assert stub.payload.tool_name == "search_maintenance_sop"
    assert stub.payload.arguments_json == '{"query":"齿轮箱油温"}'
    assert stub.payload.step == 2
    assert result.tool_call.tool_call_id == 17
    assert result.tool_call.status == "success"
    assert result.evidence_json == '{"chunks":[]}'
    assert [(item.document_id, item.chunk_id) for item in result.citations] == [(3, 9)]

