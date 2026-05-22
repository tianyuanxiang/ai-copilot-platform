import base64
import zipfile
from io import BytesIO

from fastapi.testclient import TestClient

from app.main import app


def test_parse_markdown_returns_structural_parent_and_child_chunks():
    client = TestClient(app)
    markdown = """
# 部署手册

## PostgreSQL

PostgreSQL 需要创建 pgvector 扩展，并配置连接池。

```yaml
postgres:
  dsn: postgres://user:pass@127.0.0.1:5432/app
```

## Elasticsearch

kb_chunks_index 必须写入 kb_type、owner_user_id、domain_id、visibility 和 uploaded_by。
"""

    response = client.post(
        "/v1/parse",
        json={"file_name": "deploy.md", "file_type": "md", "content": markdown},
    )

    assert response.status_code == 200
    body = response.json()
    assert body["metadata"]["parser"] == "markdown-structure"
    assert body["metadata"]["parent_count"] >= 1
    assert body["metadata"]["child_count"] >= 1
    assert "PostgreSQL" in body["text"]
    assert any("```yaml" in parent["content"] for parent in body["parents"])
    assert all(parent["children"] for parent in body["parents"])


def test_parse_docx_extracts_paragraphs_and_tables():
    client = TestClient(app)
    content = base64.b64encode(_minimal_docx_bytes()).decode("ascii")

    response = client.post(
        "/v1/parse",
        json={"file_name": "knowledge.docx", "content": content},
    )

    assert response.status_code == 200
    body = response.json()
    assert body["metadata"]["file_type"] == "docx"
    assert body["metadata"]["parser"] == "docx-xml"
    assert "知识库说明" in body["text"]
    assert "字段 | 含义" in body["text"]
    assert body["parents"][0]["children"]


def test_parse_docx_rejects_non_base64_content():
    client = TestClient(app)

    response = client.post(
        "/v1/parse",
        json={"file_name": "broken.docx", "content": "not a docx"},
    )

    assert response.status_code == 400
    assert "base64" in response.json()["detail"]


def _minimal_docx_bytes() -> bytes:
    document_xml = """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p>
      <w:r><w:t>知识库说明</w:t></w:r>
    </w:p>
    <w:p>
      <w:r><w:t>这是一段从 docx 抽取出来的正文。</w:t></w:r>
    </w:p>
    <w:tbl>
      <w:tr>
        <w:tc><w:p><w:r><w:t>字段</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>含义</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:t>uploaded_by</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>上传人</w:t></w:r></w:p></w:tc>
      </w:tr>
    </w:tbl>
  </w:body>
</w:document>
"""
    buffer = BytesIO()
    with zipfile.ZipFile(buffer, "w") as package:
        package.writestr("word/document.xml", document_xml)
    return buffer.getvalue()
