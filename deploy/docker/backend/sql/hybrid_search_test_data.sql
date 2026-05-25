-- Hybrid Search standard test data.
-- Run after the normal schema is installed and pgvector extension is enabled.
-- The vector values use deterministic 1024-dimensional mock vectors for local validation.

CREATE EXTENSION IF NOT EXISTS vector;

DELETE FROM ai_document_chunk WHERE document_id BETWEEN 301 AND 308;
DELETE FROM ai_document_parent_chunk WHERE document_id BETWEEN 301 AND 308;
DELETE FROM ai_document WHERE id BETWEEN 301 AND 308;
DELETE FROM ai_kb_member WHERE kb_id BETWEEN 201 AND 208;
DELETE FROM ai_knowledge_base WHERE id BETWEEN 201 AND 208;
DELETE FROM ai_kb_domain WHERE id IN (101, 102);

INSERT INTO ai_kb_domain (id, name, description, status, created_by, created_at, updated_at)
VALUES
  (101, '后端研发', 'go-zero、网关、微服务相关知识库领域', 1, 1001, NOW(), NOW()),
  (102, '运维部署', 'etcd、Docker、部署配置相关知识库领域', 1, 1001, NOW(), NOW());

INSERT INTO ai_knowledge_base
  (id, kb_type, owner_user_id, domain_id, name, description, visibility, status, created_by, created_at, updated_at)
VALUES
  (201, 'personal', 1001, NULL, '1001 的个人部署知识库', '当前用户自己的私有知识库', 'private', 1, 1001, NOW(), NOW()),
  (202, 'personal', 1003, NULL, '1003 的个人知识库', '当前用户无权访问', 'private', 1, 1003, NOW(), NOW()),
  (203, 'public', NULL, 101, '公开后端研发知识库', '所有人可读', 'public', 1, 1001, NOW(), NOW()),
  (204, 'public', NULL, 102, '私有运维部署知识库 viewer', '1001 是 viewer', 'private', 1, 1001, NOW(), NOW()),
  (205, 'public', NULL, 102, '私有运维部署知识库 no access', '1001 不是成员', 'private', 1, 1002, NOW(), NOW()),
  (206, 'public', NULL, 101, '私有后端研发知识库 editor', '1001 是 editor', 'private', 1, 1001, NOW(), NOW()),
  (207, 'public', NULL, 101, '私有后端研发知识库 manager', '1001 是 manager', 'private', 1, 1001, NOW(), NOW()),
  (208, 'public', NULL, 101, '禁用知识库', 'status 禁用，不应被检索', 'public', 0, 1001, NOW(), NOW());

INSERT INTO ai_kb_member (kb_id, user_id, role, created_at, updated_at)
VALUES
  (204, 1001, 'viewer', NOW(), NOW()),
  (206, 1001, 'editor', NOW(), NOW()),
  (207, 1001, 'manager', NOW(), NOW()),
  (205, 1002, 'viewer', NOW(), NOW());

INSERT INTO ai_document
  (id, kb_id, uploaded_by, file_name, file_type, content_hash, status, error_msg, created_at, updated_at)
VALUES
  (301, 201, 1001, 'go-zero部署文档.md', 'md', 'hybrid-doc-301', 'ready', '', NOW(), NOW()),
  (302, 203, 1001, 'etcd配置说明.md', 'md', 'hybrid-doc-302', 'ready', '', NOW(), NOW()),
  (303, 204, 1001, 'Docker部署手册.md', 'md', 'hybrid-doc-303', 'ready', '', NOW(), NOW()),
  (304, 206, 1001, '微服务网关说明.md', 'md', 'hybrid-doc-304', 'ready', '', NOW(), NOW()),
  (305, 207, 1001, 'RAG检索说明.md', 'md', 'hybrid-doc-305', 'ready', '', NOW(), NOW()),
  (306, 202, 1003, '无权限个人文档.md', 'md', 'hybrid-doc-306', 'ready', '', NOW(), NOW()),
  (307, 201, 1001, '失败文档.md', 'md', 'hybrid-doc-307', 'failed', 'test failed document', NOW(), NOW()),
  (308, 208, 1001, '禁用知识库文档.md', 'md', 'hybrid-doc-308', 'ready', '', NOW(), NOW());

INSERT INTO ai_document_parent_chunk
  (id, document_id, parent_index, content, token_count, created_at)
VALUES
  (401, 301, 0, 'go-zero 服务部署包括配置 gateway、rpc、etcd 注册中心和 systemd 启动流程。', 32, NOW()),
  (402, 302, 0, 'etcd 配置说明包括 Hosts、Key、服务发现、注册中心地址和 go-zero yaml 配置位置。', 34, NOW()),
  (403, 303, 0, 'Docker 部署手册说明 Docker Compose、容器网络、环境变量和后端服务启动顺序。', 32, NOW()),
  (404, 304, 0, '微服务网关说明介绍 gateway 路由、AuthMiddleware、CasbinMiddleware 和 RPC 调用。', 34, NOW()),
  (405, 305, 0, 'RAG 检索说明介绍 pgvector 语义检索、Elasticsearch BM25 和 RRF 融合排序。', 36, NOW()),
  (406, 306, 0, '这是无权限用户 1003 的私有文档，当前用户 1001 不应该检索到。', 30, NOW()),
  (407, 307, 0, '这是 failed 文档，即使属于 1001，也不应该出现在检索结果中。', 28, NOW()),
  (408, 308, 0, '这是禁用知识库中的文档，即使文档 ready，也不应该被跨库检索返回。', 30, NOW());

INSERT INTO ai_document_chunk
  (id, document_id, parent_chunk_id, chunk_index, content, embedding, token_count, created_at)
VALUES
  (501, 301, 401, 0, 'go-zero 怎么部署：先准备 etcd，然后启动 ai.rpc，再启动 gateway-api，并检查 etcd 服务注册。', ('[' || array_to_string(array_fill(0.10::float8, ARRAY[1024]), ',') || ']')::vector, 36, NOW()),
  (502, 301, 401, 1, 'go-zero 部署配置文件通常在 etc/gateway-api.yaml 和 service/ai/rpc/etc/ai.yaml。', ('[' || array_to_string(array_fill(0.11::float8, ARRAY[1024]), ',') || ']')::vector, 30, NOW()),
  (503, 302, 402, 0, 'etcd 配置在哪里：go-zero 的 Etcd Hosts 和 Key 写在 yaml 配置文件中。', ('[' || array_to_string(array_fill(0.20::float8, ARRAY[1024]), ',') || ']')::vector, 32, NOW()),
  (504, 302, 402, 1, '服务发现依赖 etcd，RPC 服务启动后会按照 Key 注册到指定 Hosts。', ('[' || array_to_string(array_fill(0.21::float8, ARRAY[1024]), ',') || ']')::vector, 30, NOW()),
  (505, 303, 403, 0, 'Docker Compose 启动服务时，需要先启动 PostgreSQL、Elasticsearch、etcd，再启动业务服务。', ('[' || array_to_string(array_fill(0.30::float8, ARRAY[1024]), ',') || ']')::vector, 34, NOW()),
  (506, 304, 404, 0, '微服务网关负责 HTTP 路由转 RPC 调用，鉴权通过 AuthMiddleware 和 CasbinMiddleware 完成。', ('[' || array_to_string(array_fill(0.40::float8, ARRAY[1024]), ',') || ']')::vector, 34, NOW()),
  (507, 305, 405, 0, 'RAG 检索使用 pgvector 做语义召回，Elasticsearch 做 BM25 关键词召回，最后用 RRF 融合。', ('[' || array_to_string(array_fill(0.50::float8, ARRAY[1024]), ',') || ']')::vector, 38, NOW()),
  (508, 306, 406, 0, '无权限 chunk：这个片段包含 go-zero 怎么部署，但当前用户 1001 不应该看到。', ('[' || array_to_string(array_fill(0.10::float8, ARRAY[1024]), ',') || ']')::vector, 34, NOW()),
  (509, 307, 407, 0, 'failed chunk：这个片段包含 etcd 配置在哪里，但文档状态不是 ready。', ('[' || array_to_string(array_fill(0.20::float8, ARRAY[1024]), ',') || ']')::vector, 32, NOW()),
  (510, 308, 408, 0, 'disabled kb chunk：这个片段包含 RAG 检索，但知识库 status 禁用。', ('[' || array_to_string(array_fill(0.50::float8, ARRAY[1024]), ',') || ']')::vector, 32, NOW());
