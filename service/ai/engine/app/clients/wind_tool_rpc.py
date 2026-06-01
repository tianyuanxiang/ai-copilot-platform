"""Async gRPC adapter for the controlled Go Wind Agent tool runtime."""

from __future__ import annotations

from dataclasses import dataclass

import grpc

from app.generated import ai_pb2, ai_pb2_grpc
from app.schemas.agent import Citation, ToolCall


@dataclass(slots=True)
class WindToolExecuteRequest:
    """Input required by Go AiWindAgentService.ExecuteTool."""

    user_id: int
    trace_id: str
    conversation_id: str
    tool_name: str
    arguments_json: str
    step: int


@dataclass(slots=True)
class WindToolExecuteResult:
    """Normalized result returned by the Go tool runtime."""

    tool_call: ToolCall
    evidence_json: str
    citations: list[Citation]


class WindToolRPCError(RuntimeError):
    """Raised when Python cannot reach the Go tool runtime."""

# Python Agent调 Go 工具服务的客户端

class WindToolRPCClient:
    """Calls Go ExecuteTool through one reusable async gRPC channel."""
    """Python Agent 调 Go 后端工具的 RPC 客户端。"""
    """ 理解用户问题 决定下一步调用什么工具  组织最终回答"""
    def __init__(
        self,
        target: str,
        timeout_seconds: float,
        *,
        channel: grpc.aio.Channel | None = None,
        stub: ai_pb2_grpc.AiWindAgentServiceStub | None = None,
    ) -> None:
        self._channel = channel or grpc.aio.insecure_channel(target)
        self._owns_channel = channel is None
        self._stub = stub or ai_pb2_grpc.AiWindAgentServiceStub(self._channel)
        self._timeout_seconds = timeout_seconds

    async def execute(self, request: WindToolExecuteRequest) -> WindToolExecuteResult:
        """Execute one allowlisted tool and map protobuf fields to Python models."""

        payload = ai_pb2.WindToolExecuteReq(
            user_id=request.user_id,
            trace_id=request.trace_id,
            conversation_id=request.conversation_id,
            tool_name=request.tool_name,
            arguments_json=request.arguments_json,
            step=request.step,
        )
        try:
            response = await self._stub.ExecuteTool(payload, timeout=self._timeout_seconds)
        except grpc.aio.AioRpcError as exc:
            raise WindToolRPCError(f"Go ExecuteTool RPC failed: {exc.code().name}: {exc.details()}") from exc

        citations = [
            Citation(
                document_id=item.document_id,
                chunk_id=item.chunk_id,
                title=item.title,
                snippet=item.snippet,
                score=item.score,
            )
            for item in response.citations
        ]
        return WindToolExecuteResult(
            tool_call=ToolCall(
                tool_call_id=response.tool_call_id,
                tool_name=response.tool_name,
                status=response.status,
                arguments_json=request.arguments_json,
                result_json=response.result_json or "{}",
                message=response.message,
                latency_ms=response.latency_ms,
            ),
            evidence_json=response.evidence_json,
            citations=citations,
        )

    async def close(self) -> None:
        """Close the owned gRPC channel during FastAPI shutdown."""

        if self._owns_channel:
            await self._channel.close()
