from pydantic import BaseModel, Field


class ParseRequest(BaseModel):
    file_name: str
    file_type: str = ""
    content: str


class ChildChunk(BaseModel):
    chunk_index: int
    content: str
    token_count: int


class ParentChunk(BaseModel):
    parent_index: int
    content: str
    token_count: int
    children: list[ChildChunk] = Field(default_factory=list)


""" 返回示例: 
{
    "text": "完整解析后的文本...",
    "metadata": {
      "file_name": "readme.md",
      "file_type": "md",
      "parser": "markdown-structure",
      "parent_count": 2,
      "child_count": 4
    },
    "parents": [
      {
        "parent_index": 0,
        "content": "第一个父块的文本内容，大约800-1200字...",
        "token_count": 250,
        "children": [
          {
            "chunk_index": 0,
            "content": "第一个子块，大约200-400字...",
            "token_count": 80
          },
        ]
      },
      {
        "parent_index": 1,
        "content": "第二个父块的文本内容...",
        "token_count": 230,
        "children": [
          {
            "chunk_index": 2,
            "content": "第三个子块...",
            "token_count": 75
          },
        ]
      }
    ]
  }
"""
class ParseResponse(BaseModel):
    text: str
    metadata: dict = Field(default_factory=dict) # dict 是"键值对"，相当于 Go 的 map[string]any;Field(default_factory=dict):如果不传这个字段，默认值是一个空字典 {}
    parents: list[ParentChunk] = Field(default_factory=list) # ParentChunk 类型的数组，可为空


class EmbedRequest(BaseModel):
    texts: list[str]


class EmbedResponse(BaseModel):
    vectors: list[list[float]]
    token_counts: list[int]
    total_tokens: int
    model: str
    dimension: int
    mode: str


class RetrieveRequest(BaseModel):
    user_id: str
    kb_id: str
    query: str
    top_k: int = 8


class RetrievedChunk(BaseModel):
    chunk_id: str
    document_id: str
    title: str
    content: str
    score: float
    source: str


class RetrieveResponse(BaseModel):
    chunks: list[RetrievedChunk]
    mode: str
    message: str


class RerankRequest(BaseModel):
    query: str
    chunks: list[RetrievedChunk]
    top_k: int = 5


class RerankResponse(BaseModel):
    chunks: list[RetrievedChunk]
    mode: str
