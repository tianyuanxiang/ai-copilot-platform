from fastapi import APIRouter, HTTPException

from app.schemas.knowledge import (
    EmbedRequest,
    EmbedResponse,
    ParseRequest,
    ParseResponse,
    RerankRequest,
    RerankResponse,
    RetrieveRequest,
    RetrieveResponse,
)
from app.services import rag

router = APIRouter(tags=["knowledge"])


@router.post("/parse", response_model=ParseResponse)
async def parse(payload: ParseRequest) -> ParseResponse:
    try:
        # async/await 就是异步编程，相当于 Go 的 goroutine + channel 的简化版
        # await = "等这个异步操作完成再往下走"
        return await rag.parse_document(payload.file_name, payload.content, payload.file_type)
    except rag.ParseInputError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc


@router.post("/embed", response_model=EmbedResponse)
async def embed(payload: EmbedRequest) -> EmbedResponse:
    try:
        return await rag.embed_texts(payload.texts)
    except rag.EmbeddingInputError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc
    except rag.EmbeddingProviderError as exc:
        raise HTTPException(status_code=502, detail=str(exc)) from exc


@router.post("/retrieve", response_model=RetrieveResponse)
async def retrieve(payload: RetrieveRequest) -> RetrieveResponse:
    return await rag.retrieve(payload.user_id, payload.kb_id, payload.query, payload.top_k)


@router.post("/rerank", response_model=RerankResponse)
async def rerank(payload: RerankRequest) -> RerankResponse:
    return await rag.rerank(payload.query, payload.chunks, payload.top_k)
