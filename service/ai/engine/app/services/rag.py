import base64
import math
import re
import zipfile
from io import BytesIO
from pathlib import Path
from typing import Any
from xml.etree import ElementTree

import httpx

from app.core.config import Settings, get_settings
from app.schemas.knowledge import (
    ChildChunk,
    EmbedResponse,
    ParentChunk,
    ParseResponse,
    RerankResponse,
    RetrieveResponse,
    RetrievedChunk,
)

DASHSCOPE_EMBEDDINGS_URL = "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings"
MAX_EMBEDDING_TOKENS = 8192
PARENT_TARGET_CHARS = 1000
PARENT_MAX_CHARS = 1200
CHILD_TARGET_CHARS = 350
CHILD_MAX_CHARS = 420
CHILD_OVERLAP_CHARS = 50

WORD_NAMESPACE = "{http://schemas.openxmlformats.org/wordprocessingml/2006/main}"
TEXT_FORMAT_MARKS = "\ufeff\u200b\u200c\u200d"


class ParseInputError(ValueError):
    pass


class EmbeddingInputError(ValueError):
    pass


class EmbeddingProviderError(RuntimeError):
    pass

# content 存的是文件内容的文本
async def parse_document(file_name: str, content: str, file_type: str = "") -> ParseResponse:
    document_type = _detect_document_type(file_name, file_type)
    # Go 把 DOCX 文件做 base64 编码后传过来
    if document_type == "docx":
        text = _extract_docx_text(_decode_docx_content(content))
        parser = "docx-xml"
    elif document_type in {"md", "markdown", "txt"}:
        text = _normalize_text(content)
        parser = "markdown-structure" if document_type in {"md", "markdown"} else "plain-text"
    else:
        text = _normalize_text(content)
        parser = "plain-text"

    parents = chunk_document(text)
    return ParseResponse(
        text=text,
        metadata={
            "file_name": file_name,
            "file_type": document_type,
            "parser": parser,
            "parent_count": len(parents),
            "child_count": sum(len(parent.children) for parent in parents),
        },
        parents=parents,
    )


def chunk_document(text: str) -> list[ParentChunk]:
    normalized = _normalize_text(text)
    if not normalized:
        return []

    blocks = _split_structural_blocks(normalized)
    parent_texts = _merge_blocks(blocks, PARENT_TARGET_CHARS, PARENT_MAX_CHARS)

    # 准备装所有父块的列表
    parents: list[ParentChunk] = []
    for parent_index, parent_text in enumerate(parent_texts):
        # 1. 把父块文本再切成子块(200-400字)
        child_texts = _split_child_chunks(parent_text)
        parents.append(
            ParentChunk(
                parent_index=parent_index,
                content=parent_text,
                token_count=_estimate_chunk_tokens(parent_text),
                # 列表推导式
                # [要生成的元素 for 临时变量 in 可遍历对象]
                children=[
                    ChildChunk(
                        chunk_index=child_index,
                        content=child_text,
                        token_count=_estimate_chunk_tokens(child_text),
                    )
                    for child_index, child_text in enumerate(child_texts)
                ],
            )
        )
    return parents


def _detect_document_type(file_name: str, file_type: str) -> str:
    clean_type = file_type.strip().lower().lstrip(".")
    if clean_type:
        return clean_type

    suffix = Path(file_name).suffix.lower().lstrip(".")
    if suffix:
        return suffix
    return "txt"


def _normalize_text(text: str) -> str:
    normalized = text.replace("\r\n", "\n").replace("\r", "\n")
    normalized = normalized.lstrip(TEXT_FORMAT_MARKS)
    normalized = re.sub(rf"\n[{TEXT_FORMAT_MARKS}]+", "\n", normalized)
    normalized = re.sub(r"[ \t]+\n", "\n", normalized)
    normalized = re.sub(r"\n{3,}", "\n\n", normalized)
    return normalized.strip()


def _decode_docx_content(content: str) -> bytes:
    raw = content.strip()
    if raw.startswith("data:"):
        _, _, raw = raw.partition(",")

    try:
        data = base64.b64decode(raw, validate=True) # # base64 解码回二进制
    except ValueError as exc:
        raise ParseInputError("docx content must be base64 encoded") from exc

    if not zipfile.is_zipfile(BytesIO(data)):
        raise ParseInputError("docx content is not a valid .docx zip package")
    return data


def _extract_docx_text(data: bytes) -> str:
    try:
        with zipfile.ZipFile(BytesIO(data)) as package:
            document_xml = package.read("word/document.xml")
    except KeyError as exc:
        raise ParseInputError("docx is missing word/document.xml") from exc
    except zipfile.BadZipFile as exc:
        raise ParseInputError("docx content is not a readable zip package") from exc

    root = ElementTree.fromstring(document_xml)
    body = root.find(f"{WORD_NAMESPACE}body")
    if body is None:
        return ""

    parts: list[str] = []
    for child in body:
        if child.tag == f"{WORD_NAMESPACE}p":
            paragraph = _docx_paragraph_text(child)
            if paragraph:
                parts.append(paragraph)
        elif child.tag == f"{WORD_NAMESPACE}tbl":
            table = _docx_table_text(child)
            if table:
                parts.append(table)

    return _normalize_text("\n\n".join(parts))


def _docx_paragraph_text(paragraph: ElementTree.Element) -> str:
    fragments: list[str] = []
    for node in paragraph.iter():
        if node.tag == f"{WORD_NAMESPACE}t" and node.text:
            fragments.append(node.text)
        elif node.tag == f"{WORD_NAMESPACE}tab":
            fragments.append("\t")
        elif node.tag in {f"{WORD_NAMESPACE}br", f"{WORD_NAMESPACE}cr"}:
            fragments.append("\n")
    return "".join(fragments).strip()


def _docx_table_text(table: ElementTree.Element) -> str:
    rows: list[str] = []
    for row in table.findall(f".//{WORD_NAMESPACE}tr"):
        cells: list[str] = []
        for cell in row.findall(f"{WORD_NAMESPACE}tc"):
            paragraphs = [
                _docx_paragraph_text(paragraph)
                for paragraph in cell.findall(f"{WORD_NAMESPACE}p")
            ]
            cell_text = " ".join(item for item in paragraphs if item).strip()
            cells.append(cell_text)
        if any(cells):
            rows.append(" | ".join(cells))
    return "\n".join(rows).strip()


def _split_structural_blocks(text: str) -> list[str]:
    blocks: list[str] = []      # 切好的小块，最终输出
    paragraph: list[str] = []   # 当前正在攒的普通段落
    code_block: list[str] = []  # 当前正在攒的代码块
    in_code_block = False       # 是否在代码块里

    def flush_paragraph() -> None:
        if paragraph:
            blocks.append("\n".join(paragraph).strip())
            paragraph.clear()

    for line in text.split("\n"):
        stripped = line.strip() # 去掉一行前后的空格。
        # 判断当前行是不是代码块的开始或结束。
        is_fence = stripped.startswith("```") or stripped.startswith("~~~")

        if is_fence:
            if in_code_block:
                code_block.append(line)
                blocks.append("\n".join(code_block).strip())
                code_block.clear()
                in_code_block = False
            else:
                flush_paragraph()
                code_block.append(line)
                in_code_block = True
            continue

        if in_code_block:
            code_block.append(line)
            continue
        # 处理标题行
        # 1. 先把前面的普通段落收尾
        # 2. 把标题单独作为一个 block
        # 3. 继续处理下一行
        if _is_markdown_heading(stripped):
            flush_paragraph()
            blocks.append(stripped)
            continue

        # 处理空行
        if not stripped:
            flush_paragraph()
            continue

        # 处理普通文本暂时放进 paragraph。后面遇到空行、标题、代码块时，才会通过 flush_paragraph() 放进 blocks
        paragraph.append(line.rstrip())

    if in_code_block and code_block:
        blocks.append("\n".join(code_block).strip())
    flush_paragraph()

    expanded: list[str] = []

    for block in blocks:
        expanded.extend(_split_oversized_block(block, PARENT_MAX_CHARS))  # 如果某个 block 超过 PARENT_MAX_CHARS，就继续把它切小。
    return [block for block in expanded if block]


# 一般只要以 # 开头，就认为是标题：
def _is_markdown_heading(line: str) -> bool:
    return bool(re.match(r"^#{1,6}\s+\S+", line))


def _split_oversized_block(block: str, max_chars: int) -> list[str]:
    if len(block) <= max_chars:
        return [block]

    if block.startswith("```") or block.startswith("~~~"):
        return _split_by_window(block, max_chars, 0)

    sentences = re.split(r"(?<=[。！？.!?])\s+", block)
    chunks: list[str] = []
    current = ""
    for sentence in sentences:
        if not sentence:
            continue
        if len(sentence) > max_chars:
            if current:
                chunks.append(current.strip())
                current = ""
            chunks.extend(_split_by_window(sentence, max_chars, 0))
            continue
        candidate = f"{current} {sentence}".strip() if current else sentence
        if len(candidate) > max_chars and current:
            chunks.append(current.strip())
            current = sentence
        else:
            current = candidate
    if current:
        chunks.append(current.strip())
    return chunks


def _merge_blocks(blocks: list[str], target_chars: int, max_chars: int) -> list[str]:
    merged: list[str] = []
    current: list[str] = []

    for block in blocks:
        candidate = _join_blocks([*current, block])
        if current and len(candidate) > max_chars: # 当前父块已经够大了，再塞 block 就超长
            merged.append(_join_blocks(current))   # 所以先把 current 合并后放进 merged
            current = [block]                      # 然后用当前 block 开启新的 current
            continue
        # 如果当前父块已经达到目标长度了，并且现在遇到一个新的 Markdown 标题，那就把前面的内容收尾，从这个标题开始新建父块。
        if current and len(_join_blocks(current)) >= target_chars and _is_markdown_heading(block):
            merged.append(_join_blocks(current))
            current = [block]
            continue
        current.append(block)

    if current:
        merged.append(_join_blocks(current)) # 默认情况：把 block 放进当前父块
    return merged


# 接收一个父块 parent_text，再把这个父块切成多个更小的子块，并且给相邻子块加一点重叠内容。

def _split_child_chunks(parent_text: str) -> list[str]:
    blocks = _split_structural_blocks(parent_text)
    child_blocks: list[str] = []
    for block in blocks:
        # extend 把列表里的元素逐个放进去
        child_blocks.extend(_split_oversized_block(block, CHILD_MAX_CHARS)) # 如果某个结构块太长，就继续切小，保证单个 block 不超过 CHILD_MAX_CHARS。
    blocks = child_blocks

    child_texts = _merge_blocks(blocks, CHILD_TARGET_CHARS, CHILD_MAX_CHARS)
    if len(child_texts) <= 1:
        return child_texts

    # 从第二个子块开始，把前一个子块的尾巴复制一点到当前子块前面。
    overlapped: list[str] = []
    for index, child_text in enumerate(child_texts):
        if index == 0:
            overlapped.append(child_text)
            continue

        prefix = _tail_text(child_texts[index - 1], CHILD_OVERLAP_CHARS)
        if prefix and not child_text.startswith(prefix):
            overlapped.append(f"{prefix}\n\n{child_text}")
        else:
            overlapped.append(child_text)
    return overlapped


# 把多个 block 用两个换行符拼成一个大字符串。

# 遍历 blocks
# 去掉每个 block 前后的空格
# 过滤掉空字符串

def _join_blocks(blocks: list[str]) -> str:
    return "\n\n".join(block.strip() for block in blocks if block.strip()).strip()


def _tail_text(text: str, size: int) -> str:
    compact = re.sub(r"\s+", " ", text).strip()
    if len(compact) <= size:
        return compact
    return compact[-size:].strip()


def _split_by_window(text: str, size: int, overlap: int) -> list[str]:
    if size <= overlap:
        raise ValueError("size must be greater than overlap")

    chunks: list[str] = []
    start = 0
    while start < len(text):
        end = min(len(text), start + size)
        chunks.append(text[start:end].strip())
        if end == len(text):
            break
        start = end - overlap
    return [chunk for chunk in chunks if chunk]


def _estimate_chunk_tokens(text: str) -> int:
    return max(1, math.ceil(len(text) / 4))


async def embed_texts(
    texts: list[str],
    settings: Settings | None = None,
    transport: httpx.AsyncBaseTransport | None = None,
) -> EmbedResponse:
    settings = settings or get_settings()
    normalized = _validate_embedding_texts(texts)

    if settings.embedding_provider == "dashscope":
        return await _embed_dashscope(normalized, settings, transport)

    token_counts = [_estimate_tokens(text) for text in normalized]
    return EmbedResponse(
        vectors=[_mock_vector(text, settings.embedding_dimension) for text in normalized],
        token_counts=token_counts,
        total_tokens=sum(token_counts),
        model=settings.embedding_model,
        dimension=settings.embedding_dimension,
        mode="mock-embedding",
    )


def _validate_embedding_texts(texts: list[str]) -> list[str]:
    if not texts:
        raise EmbeddingInputError("texts must contain at least one item")

    normalized: list[str] = []
    for index, text in enumerate(texts):
        clean_text = text.strip()
        if not clean_text:
            raise EmbeddingInputError(f"texts[{index}] must not be blank")
        token_count = _estimate_tokens(clean_text)
        if token_count > MAX_EMBEDDING_TOKENS:
            raise EmbeddingInputError(f"texts[{index}] exceeds {MAX_EMBEDDING_TOKENS} tokens")
        normalized.append(clean_text)
    return normalized


def _estimate_tokens(text: str) -> int:
    return max(1, len(text.split()))


def _mock_vector(text: str, dimension: int) -> list[float]:
    seed = sum(ord(char) for char in text) or 1
    return [float(((seed + idx) % 997) / 997) for idx in range(dimension)]


async def _embed_dashscope(
    texts: list[str],
    settings: Settings,
    transport: httpx.AsyncBaseTransport | None,
) -> EmbedResponse:
    if not settings.embedding_api_key:
        raise EmbeddingProviderError("embedding api key is empty; set embedding.api_key or DASHSCOPE_API_KEY")

    batch_size = max(1, min(settings.embedding_batch_size, 10))
    vectors: list[list[float]] = []
    token_counts: list[int] = []
    total_tokens = 0

    async with httpx.AsyncClient(timeout=settings.embedding_timeout_seconds, transport=transport) as client:
        for start in range(0, len(texts), batch_size):
            batch = texts[start : start + batch_size]
            payload = {
                "model": settings.embedding_model,
                "input": batch,
                "dimensions": settings.embedding_dimension,
                "encoding_format": "float",
            }
            response = await client.post(
                DASHSCOPE_EMBEDDINGS_URL,
                headers={"Authorization": f"Bearer {settings.embedding_api_key}"},
                json=payload,
            )
            if response.status_code < 200 or response.status_code >= 300:
                raise EmbeddingProviderError(
                    f"dashscope embedding returned {response.status_code}: {response.text[:4096]}"
                )

            body = response.json()
            batch_vectors = _extract_dashscope_vectors(body, settings.embedding_dimension)
            if len(batch_vectors) != len(batch):
                raise EmbeddingProviderError(
                    f"dashscope embedding returned {len(batch_vectors)} vectors for {len(batch)} texts"
                )
            vectors.extend(batch_vectors)

            usage_total = _extract_usage_total(body)
            if usage_total > 0:
                batch_token_counts = _allocate_usage_tokens(usage_total, len(batch))
            else:
                batch_token_counts = [_estimate_tokens(text) for text in batch]
            token_counts.extend(batch_token_counts)
            total_tokens += sum(batch_token_counts)

    return EmbedResponse(
        vectors=vectors,
        token_counts=token_counts,
        total_tokens=total_tokens,
        model=settings.embedding_model,
        dimension=settings.embedding_dimension,
        mode="dashscope-embedding",
    )


def _extract_dashscope_vectors(body: dict[str, Any], expected_dimension: int) -> list[list[float]]:
    data = body.get("data")
    if not isinstance(data, list):
        raise EmbeddingProviderError("dashscope embedding response missing data list")

    vectors: list[list[float]] = []
    for item in data:
        if not isinstance(item, dict):
            raise EmbeddingProviderError("dashscope embedding response contains invalid data item")
        embedding = item.get("embedding")
        if not isinstance(embedding, list):
            raise EmbeddingProviderError("dashscope embedding response missing embedding vector")
        if len(embedding) != expected_dimension:
            raise EmbeddingProviderError(
                f"dashscope embedding dimension is {len(embedding)}, expected {expected_dimension}"
            )
        vectors.append([float(value) for value in embedding])
    return vectors


def _extract_usage_total(body: dict[str, Any]) -> int:
    usage = body.get("usage")
    if not isinstance(usage, dict):
        return 0
    value = usage.get("total_tokens") or usage.get("prompt_tokens") or 0
    try:
        return int(value)
    except (TypeError, ValueError):
        return 0


def _allocate_usage_tokens(total_tokens: int, item_count: int) -> list[int]:
    base = total_tokens // item_count
    remainder = total_tokens % item_count
    return [base + (1 if index < remainder else 0) for index in range(item_count)]


async def retrieve(user_id: str, kb_id: str, query: str, top_k: int) -> RetrieveResponse:
    settings = get_settings()
    mode = "hybrid-degraded"
    message = "pgvector and Elasticsearch are not connected yet; returning mock chunks"
    if settings.elasticsearch_url:
        message = f"Elasticsearch BM25 channel reserved at {settings.elasticsearch_url}"

    chunks = [
        RetrievedChunk(
            chunk_id="chunk-demo-1",
            document_id="doc-demo-1",
            title="ASGI Engine Skeleton",
            content=f"Mock retrieval result for query: {query}",
            score=0.5,
            source="mock",
        )
    ][:top_k]
    return RetrieveResponse(chunks=chunks, mode=mode, message=message)


async def rerank(query: str, chunks: list[RetrievedChunk], top_k: int) -> RerankResponse:
    ranked = sorted(chunks, key=lambda item: item.score, reverse=True)[:top_k]
    return RerankResponse(chunks=ranked, mode="mock-rerank")
