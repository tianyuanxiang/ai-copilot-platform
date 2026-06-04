/*
 Navicat Premium Data Transfer

 Source Server         : 云服务器-8.140.204.47
 Source Server Type    : PostgreSQL
 Source Server Version : 140023
 Source Host           : 8.140.204.47:5432
 Source Catalog        : ai_copilot
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140023
 File Encoding         : 65001

 Date: 04/06/2026 14:16:17
*/


-- ----------------------------
-- Type structure for halfvec
-- ----------------------------
DROP TYPE IF EXISTS "public"."halfvec";
CREATE TYPE "public"."halfvec" (
  INPUT = "public"."halfvec_in",
  OUTPUT = "public"."halfvec_out",
  RECEIVE = "public"."halfvec_recv",
  SEND = "public"."halfvec_send",
  TYPMOD_IN = "public"."halfvec_typmod_in",
  INTERNALLENGTH = VARIABLE,
  STORAGE = external,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."halfvec" OWNER TO "postgres";

-- ----------------------------
-- Type structure for sparsevec
-- ----------------------------
DROP TYPE IF EXISTS "public"."sparsevec";
CREATE TYPE "public"."sparsevec" (
  INPUT = "public"."sparsevec_in",
  OUTPUT = "public"."sparsevec_out",
  RECEIVE = "public"."sparsevec_recv",
  SEND = "public"."sparsevec_send",
  TYPMOD_IN = "public"."sparsevec_typmod_in",
  INTERNALLENGTH = VARIABLE,
  STORAGE = external,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."sparsevec" OWNER TO "postgres";

-- ----------------------------
-- Type structure for vector
-- ----------------------------
DROP TYPE IF EXISTS "public"."vector";
CREATE TYPE "public"."vector" (
  INPUT = "public"."vector_in",
  OUTPUT = "public"."vector_out",
  RECEIVE = "public"."vector_recv",
  SEND = "public"."vector_send",
  TYPMOD_IN = "public"."vector_typmod_in",
  INTERNALLENGTH = VARIABLE,
  STORAGE = external,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."vector" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for ai_alarm_analysis_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_alarm_analysis_id_seq";
CREATE SEQUENCE "public"."ai_alarm_analysis_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_conversation_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_conversation_id_seq";
CREATE SEQUENCE "public"."ai_conversation_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_document_chunk_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_document_chunk_id_seq";
CREATE SEQUENCE "public"."ai_document_chunk_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_document_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_document_id_seq";
CREATE SEQUENCE "public"."ai_document_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_document_parent_chunk_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_document_parent_chunk_id_seq";
CREATE SEQUENCE "public"."ai_document_parent_chunk_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_health_report_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_health_report_id_seq";
CREATE SEQUENCE "public"."ai_health_report_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_kb_domain_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_kb_domain_id_seq";
CREATE SEQUENCE "public"."ai_kb_domain_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_kb_member_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_kb_member_id_seq";
CREATE SEQUENCE "public"."ai_kb_member_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_knowledge_base_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_knowledge_base_id_seq";
CREATE SEQUENCE "public"."ai_knowledge_base_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_llm_call_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_llm_call_log_id_seq";
CREATE SEQUENCE "public"."ai_llm_call_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_maintenance_ticket_draft_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_maintenance_ticket_draft_id_seq";
CREATE SEQUENCE "public"."ai_maintenance_ticket_draft_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_message_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_message_id_seq";
CREATE SEQUENCE "public"."ai_message_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_model_config_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_model_config_id_seq";
CREATE SEQUENCE "public"."ai_model_config_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for ai_tool_call_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ai_tool_call_log_id_seq";
CREATE SEQUENCE "public"."ai_tool_call_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for ai_alarm_analysis
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_alarm_analysis";
CREATE TABLE "public"."ai_alarm_analysis" (
  "id" int8 NOT NULL DEFAULT nextval('ai_alarm_analysis_id_seq'::regclass),
  "user_id" int8 NOT NULL DEFAULT 0,
  "trace_id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "alarm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "content" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "evidence" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'draft'::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."ai_alarm_analysis" IS 'AI 告警分析草稿表，一期只保存脚手架结果';

-- ----------------------------
-- Table structure for ai_conversation
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_conversation";
CREATE TABLE "public"."ai_conversation" (
  "id" int8 NOT NULL DEFAULT nextval('ai_conversation_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "kb_id" int8,
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "conversation_summary" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "deleted_at" timestamptz(6),
  "deleted_by" int8
)
;
COMMENT ON COLUMN "public"."ai_conversation"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_conversation"."user_id" IS '会话所属用户ID';
COMMENT ON COLUMN "public"."ai_conversation"."kb_id" IS '关联知识库ID，可空（不在知识库上下文中的对话）';
COMMENT ON COLUMN "public"."ai_conversation"."title" IS '会话标题';
COMMENT ON COLUMN "public"."ai_conversation"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."ai_conversation"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."ai_conversation"."conversation_summary" IS '会话长期上下文滚动摘要';
COMMENT ON COLUMN "public"."ai_conversation"."deleted_at" IS '软删除时间';
COMMENT ON COLUMN "public"."ai_conversation"."deleted_by" IS '软删除操作人ID';
COMMENT ON TABLE "public"."ai_conversation" IS '对话会话表，保存一次问答会话';

-- ----------------------------
-- Table structure for ai_document
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_document";
CREATE TABLE "public"."ai_document" (
  "id" int8 NOT NULL DEFAULT nextval('ai_document_id_seq'::regclass),
  "kb_id" int8 NOT NULL,
  "uploaded_by" int8 NOT NULL,
  "file_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "file_type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "content_hash" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'pending'::character varying,
  "error_msg" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_document"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_document"."kb_id" IS '所属知识库ID，外键关联ai_knowledge_base';
COMMENT ON COLUMN "public"."ai_document"."uploaded_by" IS '上传人用户ID';
COMMENT ON COLUMN "public"."ai_document"."file_name" IS '文件名';
COMMENT ON COLUMN "public"."ai_document"."file_type" IS '文件类型，如md、pdf、docx等';
COMMENT ON COLUMN "public"."ai_document"."content_hash" IS '文档内容哈希，用于去重校验，同一kb_id下不可重复';
COMMENT ON COLUMN "public"."ai_document"."status" IS '文档处理状态，pending/parsing/chunking/indexing/ready/failed';
COMMENT ON COLUMN "public"."ai_document"."error_msg" IS '处理失败时的错误信息';
COMMENT ON COLUMN "public"."ai_document"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."ai_document"."updated_at" IS '更新时间';
COMMENT ON TABLE "public"."ai_document" IS '文档表，保存文档元信息和入库状态，文档归属知识库';


-- ----------------------------
-- Table structure for ai_document_chunk
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_document_chunk";
CREATE TABLE "public"."ai_document_chunk" (
  "id" int8 NOT NULL DEFAULT nextval('ai_document_chunk_id_seq'::regclass),
  "document_id" int8 NOT NULL,
  "parent_chunk_id" int8 NOT NULL,
  "chunk_index" int4 NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "embedding" "public"."vector",
  "token_count" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_document_chunk"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_document_chunk"."document_id" IS '所属文档ID，外键关联ai_document，级联删除';
COMMENT ON COLUMN "public"."ai_document_chunk"."parent_chunk_id" IS '所属父块ID，外键关联ai_document_parent_chunk，级联删除';
COMMENT ON COLUMN "public"."ai_document_chunk"."chunk_index" IS '子块在父块中的索引顺序';
COMMENT ON COLUMN "public"."ai_document_chunk"."content" IS '子块文本内容，200-400中文字符，用于向量检索召回';
COMMENT ON COLUMN "public"."ai_document_chunk"."embedding" IS '子块向量，维度取决于embedding模型（当前1536维），用于pgvector语义检索';
COMMENT ON COLUMN "public"."ai_document_chunk"."token_count" IS '子块token数量';
COMMENT ON COLUMN "public"."ai_document_chunk"."created_at" IS '创建时间';
COMMENT ON TABLE "public"."ai_document_chunk" IS '文档子块表，保存子块文本和向量，负责检索召回';


-- ----------------------------
-- Table structure for ai_document_parent_chunk
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_document_parent_chunk";
CREATE TABLE "public"."ai_document_parent_chunk" (
  "id" int8 NOT NULL DEFAULT nextval('ai_document_parent_chunk_id_seq'::regclass),
  "document_id" int8 NOT NULL,
  "parent_index" int4 NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "token_count" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_document_parent_chunk"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_document_parent_chunk"."document_id" IS '所属文档ID，外键关联ai_document，级联删除';
COMMENT ON COLUMN "public"."ai_document_parent_chunk"."parent_index" IS '父块在文档中的顺序索引';
COMMENT ON COLUMN "public"."ai_document_parent_chunk"."content" IS '父块文本内容，800-1200中文字符，用于给LLM提供完整上下文';
COMMENT ON COLUMN "public"."ai_document_parent_chunk"."token_count" IS '父块token数量';
COMMENT ON COLUMN "public"."ai_document_parent_chunk"."created_at" IS '创建时间';
COMMENT ON TABLE "public"."ai_document_parent_chunk" IS '文档父块表，保存父块文本，负责给LLM提供完整上下文';

-- ----------------------------
-- Table structure for ai_health_report
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_health_report";
CREATE TABLE "public"."ai_health_report" (
  "id" int8 NOT NULL DEFAULT nextval('ai_health_report_id_seq'::regclass),
  "user_id" int8 NOT NULL DEFAULT 0,
  "trace_id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "report_type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'health'::character varying,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "start_time" timestamp(0),
  "end_time" timestamp(0),
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "content" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "evidence" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'draft'::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."ai_health_report" IS 'AI 健康报告草稿表，后续扩展日报、周报、单机报告和故障复盘';

-- ----------------------------
-- Records of ai_health_report
-- ----------------------------
INSERT INTO "public"."ai_health_report" VALUES (1, 1, 'wind-1780039862330142900', 'health', 'FY', '4', '2026-05-16 00:55:00', '2026-05-29 01:00:00', 'FY 4 健康报告草稿', '{"message":"报告由 LangGraph 规则模板生成，事实完全来自 Go 侧 evidence。；已执行可选 LLM 润色","metrics":{"alarm":{"alarm_count":16652,"by_alarm_code_top":[{"count":2842,"key":"10004"},{"count":2818,"key":"10086"},{"count":2807,"key":"10001"},{"count":2757,"key":"10002"},{"count":2742,"key":"10000"},{"count":2686,"key":"10003"}],"by_time_bucket":[{"bucket_start":"2026-05-16 00:00:00.000","count":7309},{"bucket_start":"2026-05-17 00:00:00.000","count":6239},{"bucket_start":"2026-05-28 00:00:00.000","count":2781},{"bucket_start":"2026-05-29 00:00:00.000","count":323}],"granularity":"1d","level_counts":{"1":4107,"2":4163,"3":4192,"4":4190},"peak_bucket":{"bucket_start":"2026-05-16 00:00:00.000","count":7309},"returned_records":27,"risk":"critical","status_counts":{"0":16652},"tower_counts":{"4":16652},"truncated":true},"timeseries":{}},"recommendations":["保存到 ai_health_report 后允许人工编辑","后续补齐日报、周报、单机报告和故障复盘模板"],"sections":[{"content":"报告类型：single_turbine_health；时间范围：2026-05-16 00:55:00 至 2026-05-29 01:00:00。","name":"一、概览"},{"content":"识别到 0 条测点记录，数值字段 0 个。","name":"二、测点趋势"},{"content":"识别到 16652 条告警，风险等级建议为 critical。","name":"三、告警情况"},{"content":"当前为规则模板草稿，提交前需由运维人员结合现场情况复核。","name":"四、风险建议"},{"content":"在线率、缺测率、阈值配置、SOP 引用和历史故障案例。","name":"五、待补证据"}],"status":"draft","summary":"问题描述：健康报告显示告警数量极高（总计16652条），风险评级为critical，且数据仅返回27条记录（truncated=true），告警峰值出现在2026-05-16（7309条）及2026-05-17（6239条），告警码10004、10086、10001、10002、10000、10003各自出现约2700-2800次，各等级告警（1~4级）数量均超过4100条。  \n建议步骤：优先排查2026-05-16至17日的告警高峰原因，针对告警码10004、10086等前六位高频码进行根因分析，检查对应测点与部件状态，并补充被截断的完整告警记录以评估整体风险。  \n证据引用：告警总数16652，级别分布（1级4107、2级4163、3级4192、4级4190），峰值时段（2026-05-16 7309条，2026-05-17 6239条），告警码top6及对应计数（2842、2818、2807、2757、2742、2686），风险等级critical，数据截断标记true。","todo":["接入健康评分","接入 SOP 引用","接入 LLM 润色"]}', '[{"total": 16652, "where": "ts >= ''2026-05-16 00:55:00'' AND ts < ''2026-05-29 01:00:00'' and tower_id=4", "source": "tdengine.alarm", "last_ts": "2026-05-29T00:57:29.984+08:00", "sampled": 27, "by_level": {"1": 4107, "2": 4163, "3": 4192, "4": 4190}, "database": "fuyu", "end_time": "2026-05-29 01:00:00", "first_ts": "2026-05-16T00:56:09.106+08:00", "scaffold": false, "by_status": {"0": 16652}, "farm_code": "FY", "truncated": true, "alarm_code": "", "start_time": "2026-05-16 00:55:00", "tower_code": "4", "granularity": "1d", "peak_bucket": {"count": 7309, "bucket_start": "2026-05-16 00:00:00.000"}, "by_tower_top": [{"key": "4", "count": 16652}], "by_time_bucket": [{"count": 7309, "bucket_start": "2026-05-16 00:00:00.000"}, {"count": 6239, "bucket_start": "2026-05-17 00:00:00.000"}, {"count": 2781, "bucket_start": "2026-05-28 00:00:00.000"}, {"count": 323, "bucket_start": "2026-05-29 00:00:00.000"}], "by_alarm_code_top": [{"key": "10004", "count": 2842}, {"key": "10086", "count": 2818}, {"key": "10001", "count": 2807}, {"key": "10002", "count": 2757}, {"key": "10000", "count": 2742}, {"key": "10003", "count": 2686}]}]', 'draft', '2026-05-29 15:31:14');

-- ----------------------------
-- Table structure for ai_kb_domain
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_kb_domain";
CREATE TABLE "public"."ai_kb_domain" (
  "id" int8 NOT NULL DEFAULT nextval('ai_kb_domain_id_seq'::regclass),
  "name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "code" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "sort" int4 NOT NULL DEFAULT 0,
  "status" int2 NOT NULL DEFAULT 1,
  "created_by" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_kb_domain"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_kb_domain"."name" IS '领域名称，如"后端开发"、"运维部署"';
COMMENT ON COLUMN "public"."ai_kb_domain"."code" IS '领域编码，唯一标识，如"backend"、"ops"';
COMMENT ON COLUMN "public"."ai_kb_domain"."description" IS '领域描述';
COMMENT ON COLUMN "public"."ai_kb_domain"."sort" IS '排序权重，值越小越靠前';
COMMENT ON COLUMN "public"."ai_kb_domain"."status" IS '状态，1=启用，0=禁用';
COMMENT ON COLUMN "public"."ai_kb_domain"."created_by" IS '创建人用户ID';
COMMENT ON COLUMN "public"."ai_kb_domain"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."ai_kb_domain"."updated_at" IS '更新时间';
COMMENT ON TABLE "public"."ai_kb_domain" IS '公共知识库领域表，管理公共知识库的分类领域';

-- ----------------------------
-- Records of ai_kb_domain
-- ----------------------------
INSERT INTO "public"."ai_kb_domain" VALUES (1, '后端开发', 'backend', 'Go、Go-zero、数据库、接口设计', 10, 1, 1, '2026-05-19 16:40:12.460062+08', '2026-05-19 16:40:12.460062+08');
INSERT INTO "public"."ai_kb_domain" VALUES (2, '运维部署', 'ops', 'Linux、Docker、服务部署、排障', 20, 1, 1, '2026-05-19 16:40:12.460062+08', '2026-05-19 16:40:12.460062+08');
INSERT INTO "public"."ai_kb_domain" VALUES (3, 'AI 应用开发', 'ai_app', 'RAG、Agent、LLMOps、模型调用', 30, 1, 1, '2026-05-19 16:40:12.460062+08', '2026-05-19 16:40:12.460062+08');
INSERT INTO "public"."ai_kb_domain" VALUES (4, '风电混塔智能运维', 'wmtd', 'AI智能运维', 1, 1, 1, '2026-06-02 22:46:47.539368+08', '2026-06-02 22:46:47.539368+08');

-- ----------------------------
-- Table structure for ai_kb_member
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_kb_member";
CREATE TABLE "public"."ai_kb_member" (
  "id" int8 NOT NULL DEFAULT nextval('ai_kb_member_id_seq'::regclass),
  "kb_id" int8 NOT NULL,
  "user_id" int8 NOT NULL,
  "role" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "created_by" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_kb_member"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_kb_member"."kb_id" IS '知识库ID，外键关联ai_knowledge_base，级联删除';
COMMENT ON COLUMN "public"."ai_kb_member"."user_id" IS '成员用户ID';
COMMENT ON COLUMN "public"."ai_kb_member"."role" IS '成员角色，viewer=可查看，editor=可编辑上传，manager=可管理成员';
COMMENT ON COLUMN "public"."ai_kb_member"."created_by" IS '添加该成员的用户ID';
COMMENT ON COLUMN "public"."ai_kb_member"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."ai_kb_member"."updated_at" IS '更新时间';
COMMENT ON TABLE "public"."ai_kb_member" IS '知识库成员表，保存公共知识库团队成员权限，个人知识库不使用此表';

-- ----------------------------
-- Records of ai_kb_member
-- ----------------------------
INSERT INTO "public"."ai_kb_member" VALUES (1, 7, 1, 'editor', 1, '2026-05-26 08:43:34.887005+08', '2026-05-26 08:43:34.887005+08');
INSERT INTO "public"."ai_kb_member" VALUES (2, 7, 2, 'viewer', 1, '2026-05-26 08:43:47.813171+08', '2026-05-26 08:43:47.813171+08');
INSERT INTO "public"."ai_kb_member" VALUES (4, 8, 1, 'manager', 1, '2026-06-03 09:57:36.540692+08', '2026-06-03 09:57:36.540692+08');

-- ----------------------------
-- Table structure for ai_knowledge_base
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_knowledge_base";
CREATE TABLE "public"."ai_knowledge_base" (
  "id" int8 NOT NULL DEFAULT nextval('ai_knowledge_base_id_seq'::regclass),
  "kb_type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "owner_user_id" int8,
  "domain_id" int8,
  "name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "visibility" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'private'::character varying,
  "status" int2 NOT NULL DEFAULT 1,
  "created_by" int8 NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_knowledge_base"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_knowledge_base"."kb_type" IS '知识库类型，personal=个人知识库，public=公共知识库';
COMMENT ON COLUMN "public"."ai_knowledge_base"."owner_user_id" IS '所有者用户ID，个人知识库必填，公共知识库为空';
COMMENT ON COLUMN "public"."ai_knowledge_base"."domain_id" IS '所属领域ID，公共知识库必填，个人知识库为空';
COMMENT ON COLUMN "public"."ai_knowledge_base"."name" IS '知识库名称';
COMMENT ON COLUMN "public"."ai_knowledge_base"."description" IS '知识库描述';
COMMENT ON COLUMN "public"."ai_knowledge_base"."visibility" IS '可见性，private=私有，public=公开可读';
COMMENT ON COLUMN "public"."ai_knowledge_base"."status" IS '状态，1=启用，0=禁用';
COMMENT ON COLUMN "public"."ai_knowledge_base"."created_by" IS '创建人用户ID';
COMMENT ON COLUMN "public"."ai_knowledge_base"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."ai_knowledge_base"."updated_at" IS '更新时间';
COMMENT ON TABLE "public"."ai_knowledge_base" IS '知识库主表，统一保存个人知识库和公共知识库';

-- ----------------------------
-- Records of ai_knowledge_base
-- ----------------------------
INSERT INTO "public"."ai_knowledge_base" VALUES (1, 'personal', 1, NULL, '我的技术知识库', '个人 Markdown 和部署文档', 'private', 1, 1, '2026-05-21 10:32:19.86099+08', '2026-05-21 10:32:19.86099+08');
INSERT INTO "public"."ai_knowledge_base" VALUES (7, 'public', NULL, 3, 'AI应用开发', 'AI应用开发相关测试文档', 'public', 1, 1, '2026-05-26 08:42:54.580897+08', '2026-05-26 08:42:54.580897+08');
INSERT INTO "public"."ai_knowledge_base" VALUES (8, 'public', NULL, 4, '风电AI运维', '', 'public', 1, 1, '2026-06-02 22:47:56.856718+08', '2026-06-02 22:47:56.856718+08');

-- ----------------------------
-- Table structure for ai_llm_call_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."ai_llm_call_log";
CREATE TABLE "public"."ai_llm_call_log" (
  "id" int8 NOT NULL DEFAULT nextval('ai_llm_call_log_id_seq'::regclass),
  "user_id" int8,
  "trace_id" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "provider" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "model" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "prompt" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "prompt_tokens" int4 NOT NULL DEFAULT 0,
  "completion_tokens" int4 NOT NULL DEFAULT 0,
  "latency_ms" int4 NOT NULL DEFAULT 0,
  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "error_msg" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
COMMENT ON COLUMN "public"."ai_llm_call_log"."id" IS '主键ID，自增';
COMMENT ON COLUMN "public"."ai_llm_call_log"."user_id" IS '调用用户ID，可空（系统调用时为空）';
COMMENT ON COLUMN "public"."ai_llm_call_log"."trace_id" IS '链路追踪ID，用于关联整个请求链路';
COMMENT ON COLUMN "public"."ai_llm_call_log"."provider" IS '模型提供商';
COMMENT ON COLUMN "public"."ai_llm_call_log"."model" IS '模型名称';
COMMENT ON COLUMN "public"."ai_llm_call_log"."prompt" IS '输入提示文本';
COMMENT ON COLUMN "public"."ai_llm_call_log"."prompt_tokens" IS '输入token数量';
COMMENT ON COLUMN "public"."ai_llm_call_log"."completion_tokens" IS '输出token数量';
COMMENT ON COLUMN "public"."ai_llm_call_log"."latency_ms" IS '调用延迟（毫秒）';
COMMENT ON COLUMN "public"."ai_llm_call_log"."status" IS '调用状态，success/failed';
COMMENT ON COLUMN "public"."ai_llm_call_log"."error_msg" IS '失败时的错误信息';
COMMENT ON COLUMN "public"."ai_llm_call_log"."created_at" IS '创建时间';
COMMENT ON TABLE "public"."ai_llm_call_log" IS 'LLM调用日志表，记录每次LLM调用，方便排查成本、延迟、失败';



DROP TABLE IF EXISTS "public"."wind_camera_record";
CREATE TABLE "public"."wind_camera_record" (
  "record_id" int4 NOT NULL,
  "record_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "capture_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "image_url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0
)
;
COMMENT ON TABLE "public"."wind_camera_record" IS '网络摄像机抓拍记录表';

-- ----------------------------
-- Records of wind_camera_record
-- ----------------------------

-- ----------------------------
-- Table structure for wind_device
-- ----------------------------
DROP TABLE IF EXISTS "public"."wind_device";
CREATE TABLE "public"."wind_device" (
  "device_id" int4 NOT NULL,
  "device_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '00'::character varying,
  "device_type_id" int4 NOT NULL DEFAULT 1,
  "device_type_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_type_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_id" int4 NOT NULL DEFAULT 1,
  "structure_id" int4 NOT NULL DEFAULT 1,
  "structure_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "structure_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "frequency" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "height" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "longitude" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "latitude" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "status" int4 NOT NULL DEFAULT 0,
  "install_date" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "rtsp_url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "video_url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "switch_ip" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "location_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_code_old" varchar(32) COLLATE "pg_catalog"."default",
  "install_image" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."wind_device"."device_id" IS '传感器实例ID';
COMMENT ON COLUMN "public"."wind_device"."device_code" IS '传感器编号, 从00开始(通道编号)';
COMMENT ON COLUMN "public"."wind_device"."device_type_id" IS '传感器类型表的主键id';
COMMENT ON COLUMN "public"."wind_device"."device_type_code" IS '传感器类型编号，如: HLS';
COMMENT ON COLUMN "public"."wind_device"."device_type_name" IS '传感器类型名称，如: 静力水准仪';
COMMENT ON COLUMN "public"."wind_device"."tower_id" IS '所属风机id';
COMMENT ON COLUMN "public"."wind_device"."structure_id" IS '风机结构类型id';
COMMENT ON COLUMN "public"."wind_device"."structure_code" IS '结构编号，例如: BL';
COMMENT ON COLUMN "public"."wind_device"."structure_name" IS '结构名称，例如: 叶片';
COMMENT ON COLUMN "public"."wind_device"."frequency" IS '采样频率，单位为HZ';
COMMENT ON COLUMN "public"."wind_device"."height" IS '传感器高度，单位为m，精度为小数点后2位。';
COMMENT ON COLUMN "public"."wind_device"."longitude" IS '传感器经度';
COMMENT ON COLUMN "public"."wind_device"."latitude" IS '传感器纬度';
COMMENT ON COLUMN "public"."wind_device"."status" IS '状态：在线0 离线 1';
COMMENT ON COLUMN "public"."wind_device"."install_date" IS '安装日期';
COMMENT ON COLUMN "public"."wind_device"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."wind_device"."updated_at" IS '记录修改时间';
COMMENT ON COLUMN "public"."wind_device"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."wind_device"."remark" IS '备注';
COMMENT ON COLUMN "public"."wind_device"."rtsp_url" IS '摄像机视频流地址';
COMMENT ON COLUMN "public"."wind_device"."video_url" IS '摄像机前端播放地址';
COMMENT ON COLUMN "public"."wind_device"."switch_ip" IS '所属的远程开关ip';
COMMENT ON COLUMN "public"."wind_device"."location_code" IS '监测位置下的逻辑编号';
COMMENT ON COLUMN "public"."wind_device"."install_image" IS '传感器安装图片';
COMMENT ON TABLE "public"."wind_device" IS '传感器实例表';

-- ----------------------------
-- Records of wind_device
-- ----------------------------
INSERT INTO "public"."wind_device" VALUES (62, '01', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:09:44', '2025-08-28 10:12:47', 0, '-', '', '', '', '04', '00', '');
INSERT INTO "public"."wind_device" VALUES (63, '02', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:10:18', '2025-08-27 15:25:04', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (177, '06', 12, 'ACCY', '垂直主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:34:14', '2025-10-24 16:34:14', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (183, '12', 12, 'ACCY', '垂直主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:42:01', '2025-10-24 16:42:01', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (64, '03', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:10:35', '2025-08-27 15:25:21', 0, '-', '', '', '', '02', '02', '');
INSERT INTO "public"."wind_device" VALUES (65, '04', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:10:51', '2025-08-27 15:25:36', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."wind_device" VALUES (88, '01', 9, 'IPC', '网络摄像机', 5, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 16:01:17', '2025-08-27 15:55:55', 0, '-', 'rtsp://admin:Hd135246@10.184.10.39:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (84, '01', 9, 'IPC', '网络摄像机', 1, 9, 'FD_00', '基础地基基坑', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 15:40:33', '2025-08-28 10:23:01', 0, '-', 'rtsp://admin:Hd135246@10.184.10.14:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (85, '01', 9, 'IPC', '网络摄像机', 2, 9, 'FD_00', '基础地基基坑', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 15:44:54', '2025-08-27 15:31:09', 0, '-', 'rtsp://admin:Hd135246@10.184.10.17:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (87, '01', 9, 'IPC', '网络摄像机', 4, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 16:00:02', '2025-08-27 15:52:31', 0, '-', 'rtsp://admin:Hd135246@10.184.10.34:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (86, '01', 9, 'IPC', '网络摄像机', 3, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 15:57:32', '2025-08-27 15:47:52', 0, '-', 'rtsp://admin:Hd135246@10.184.10.29:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (187, '16', 12, 'ACCY', '垂直主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:44:39', '2025-10-24 16:44:39', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (67, '02', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:11:24', '2025-08-27 15:40:14', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (68, '03', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:11:47', '2025-08-27 15:40:22', 0, '-', '', '', '', '02', '02', '');
INSERT INTO "public"."wind_device" VALUES (69, '04', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:12:02', '2025-08-27 15:40:30', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."wind_device" VALUES (180, '09', 4, 'ACCX', '主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:40:29', '2025-10-24 16:40:29', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (91, '03', 4, 'ACCX', '主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 10:45:52', '2025-09-04 10:45:52', 0, '-', '', '', '', '01', '03', '');
INSERT INTO "public"."wind_device" VALUES (191, '20', 12, 'ACCY', '垂直主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:58:58', '2025-10-24 16:58:58', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (89, '02', 4, 'ACCX', '主风向加速度计', 1, 16, 'TW_20', '塔架_20米', '50', '20', 'x', '', 0, '', '2025-09-04 10:41:22', '2025-09-04 10:41:22', 1, '-', '', '', '', '', '01', '');
INSERT INTO "public"."wind_device" VALUES (35, '14', 4, 'ACCX', '主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:48', '2025-08-27 15:41:24', 0, '2', '', '', '', '02', '12', '');
INSERT INTO "public"."wind_device" VALUES (39, '06', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:30:34', '2025-07-31 14:30:34', 0, '2', '', '', '', '04', '04', '');
INSERT INTO "public"."wind_device" VALUES (52, '03', 4, 'ACCX', '主风向加速度计', 4, 1, 'BL', '叶片', '51', '80.12', '143.648', '45.856', 0, '', '2025-08-04 11:15:14', '2006-01-02 15:04:05', 1, '2', '', '', '', '', '02', '');
INSERT INTO "public"."wind_device" VALUES (80, '01', 10, 'ULS', '超声波液位计', 5, 1, 'BL', '叶片', '', '', '', '', 0, '', '2025-08-07 10:15:52', '2025-08-07 10:15:52', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."wind_device" VALUES (79, '01', 10, 'ULS', '超声波液位计', 4, 1, 'BL', '叶片', '', '', '', '', 0, '', '2025-08-07 10:15:30', '2025-08-07 10:15:30', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."wind_device" VALUES (131, '01', 9, 'IPC', '网络摄像机', 12, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:25:01', '2025-09-18 17:25:01', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (128, '01', 9, 'IPC', '网络摄像机', 9, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:22:09', '2025-09-18 17:22:09', 0, '-', 'rtmp://172.16.90.70/live/33', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (1, '01', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:20:04', '2025-08-27 15:33:27', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (2, '02', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:20:31', '2025-08-27 15:33:47', 0, '2', '', '', '', '02', '01', '');
INSERT INTO "public"."wind_device" VALUES (3, '03', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:20:50', '2025-08-27 15:34:00', 0, '2', '', '', '', '03', '02', '');
INSERT INTO "public"."wind_device" VALUES (17, '01', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:25:01', '2025-08-27 15:53:14', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (138, '02', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:47:51', '2025-10-24 15:47:51', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (70, '01', 1, 'HLS', '静力水准仪', 1, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:12:53', '2025-10-09 17:14:15', 0, '-', '', '', '', '01', '00', 'http://59.110.219.98:81/device_image/accel_1.png');
INSERT INTO "public"."wind_device" VALUES (134, '02', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:44:10', '2025-10-24 15:44:10', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (155, '02', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:02:04', '2025-10-24 16:02:04', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (135, '03', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:44:39', '2025-10-24 15:44:39', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (184, '13', 4, 'ACCX', '主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:43:33', '2025-10-24 16:43:33', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (188, '17', 4, 'ACCX', '主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:57:53', '2025-10-24 16:57:53', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (141, '05', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:49:34', '2025-10-24 15:49:34', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (145, '09', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:52:18', '2025-10-24 15:52:18', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (149, '01', 10, 'ULS', '超声波液位计', 11, 9, 'FD_00', '基础地基基坑', '50', '00', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:55:04', '2025-10-24 15:55:04', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (152, '03', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:58:06', '2025-10-24 15:58:06', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (37, '01', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:30:13', '2025-08-27 15:41:49', 0, '2', '', '', '', '02', '00', '');
INSERT INTO "public"."wind_device" VALUES (44, '02', 4, 'ACCX', '主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:31:48', '2025-07-31 14:31:48', 0, '2', '', '', '', '02', '01', '');
INSERT INTO "public"."wind_device" VALUES (158, '05', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:03:45', '2025-10-24 16:03:45', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (42, '06', 4, 'ACCX', '主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:31:19', '2025-08-27 15:50:47', 0, '2', '', '', '', '01', '05', '');
INSERT INTO "public"."wind_device" VALUES (29, '18', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:32', '2025-08-27 15:39:00', 0, '2', '', '', '', '02', '17', '');
INSERT INTO "public"."wind_device" VALUES (30, '19', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '123.648', '41.856', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:48', '2025-08-27 15:39:10', 0, '2', '', '', '', '01', '19', '');
INSERT INTO "public"."wind_device" VALUES (81, '01', 3, 'GNSS', 'GNSS', 3, 17, 'NA_156', '机舱顶', 'x', 'x', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:16:33', '2025-12-10 09:49:31', 0, '-', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (53, '01', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:05:37', '2025-08-27 15:29:41', 0, '-', '', '', '', '02', '00', '');
INSERT INTO "public"."wind_device" VALUES (54, '02', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:05:54', '2025-08-27 15:30:00', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (55, '03', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:06:17', '2025-08-27 15:30:17', 0, '-', '', '', '', '04', '02', '');
INSERT INTO "public"."wind_device" VALUES (56, '04', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:07:00', '2025-08-27 15:30:28', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."wind_device" VALUES (60, '04', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:08:27', '2025-08-27 15:46:43', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."wind_device" VALUES (83, '04', 1, 'HLS', '静力水准仪', 2, 1, 'BL', '叶片', '22', '22', '22', '22', 0, '2025-08-13T16:00:00.000Z', '2025-08-07 11:34:17', '2025-08-07 11:34:17', 1, '-', '', '', '', '', '03', '');
INSERT INTO "public"."wind_device" VALUES (74, '02', 1, 'HLS', '静力水准仪', 3, 8, 'TW_00', '塔架塔底平台', '50', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:43', '2025-08-27 15:47:24', 0, '-', '', '', '', '03', '01', '');
INSERT INTO "public"."wind_device" VALUES (75, '03', 1, 'HLS', '静力水准仪', 3, 8, 'TW_00', '塔架塔底平台', '50', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:14:00', '2025-08-27 15:47:33', 0, '-', '', '', '', '02', '02', '');
INSERT INTO "public"."wind_device" VALUES (49, '02', 2, 'INSX', '主风向倾角传感器', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:34:06', '2025-08-27 15:50:11', 0, '2', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (47, '01', 2, 'INSX', '主风向倾角传感器', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:33:24', '2025-08-27 15:38:13', 0, '2', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (120, '01', 13, 'INSY', '垂直主风向倾角传感器', 4, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'x', 'x', 0, '', '2025-09-12 09:57:11', '2025-09-12 09:57:11', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (118, '02', 13, 'INSY', '垂直主风向倾角传感器', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-09-11T16:00:00.000Z', '2025-09-12 09:50:02', '2025-09-12 15:56:17', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (119, '01', 13, 'INSY', '垂直主风向倾角传感器', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '', '2025-09-12 09:52:47', '2025-09-12 09:52:47', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (121, '02', 13, 'INSY', '垂直主风向倾角传感器', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '', '2025-09-12 09:58:16', '2025-09-12 09:58:16', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (161, '08', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:06:37', '2025-10-24 16:06:37', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (164, '11', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:08:08', '2025-10-24 16:08:08', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (167, '01', 8, 'WPR', '测风雷达', 12, 17, 'NA_156', '机舱顶', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:12:23', '2025-10-24 16:13:01', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (170, '01', 2, 'INSX', '主风向倾角传感器', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:17:49', '2025-10-24 16:17:49', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (173, '02', 12, 'ACCY', '垂直主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:30:36', '2025-10-24 16:30:36', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (31, '22', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:58', '2025-08-27 15:39:24', 0, '2', '', '', '', '04', '20', '');
INSERT INTO "public"."wind_device" VALUES (32, '23', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:10', '2025-08-27 15:39:36', 0, '2', '', '', '', '03', '23', '');
INSERT INTO "public"."wind_device" VALUES (66, '01', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:11:04', '2025-08-27 15:40:03', 0, '-', '', '', '', '04', '00', '');
INSERT INTO "public"."wind_device" VALUES (9, '09', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:05', '2025-08-27 15:35:48', 0, '2', '', '', '', '01', '08', '');
INSERT INTO "public"."wind_device" VALUES (10, '10', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:08', '2025-08-27 15:35:59', 0, '2', '', '', '', '02', '09', '');
INSERT INTO "public"."wind_device" VALUES (11, '11', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:18', '2025-08-27 15:36:09', 0, '2', '', '', '', '03', '10', '');
INSERT INTO "public"."wind_device" VALUES (175, '04', 12, 'ACCY', '垂直主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:32:31', '2025-10-24 16:32:31', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (181, '10', 12, 'ACCY', '垂直主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:40:53', '2025-10-24 16:40:53', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (185, '14', 12, 'ACCY', '垂直主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:43:54', '2025-10-24 16:43:54', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (57, '01', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:07:34', '2025-08-27 15:46:16', 0, '-', '', '', '', '02', '00', '');
INSERT INTO "public"."wind_device" VALUES (58, '02', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:07:54', '2025-08-27 15:46:24', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (59, '03', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:08:09', '2025-08-27 15:46:33', 0, '-', '', '', '', '04', '02', '');
INSERT INTO "public"."wind_device" VALUES (109, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '35', '', '', 0, '', '2025-09-10 16:01:23', '2025-09-10 16:01:23', 1, '-', '', '', '', '35', '02', '');
INSERT INTO "public"."wind_device" VALUES (110, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '40', '', '', 0, '', '2025-09-10 16:03:25', '2025-09-10 16:03:25', 1, '-', '', '', '', '40', '03', '');
INSERT INTO "public"."wind_device" VALUES (73, '01', 1, 'HLS', '静力水准仪', 3, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:26', '2025-08-27 15:47:13', 0, '-', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (46, '02', 2, 'INSX', '主风向倾角传感器', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:33:02', '2025-09-12 15:56:07', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (12, '12', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:21', '2025-08-27 15:36:21', 0, '2', '', '', '', '04', '11', '');
INSERT INTO "public"."wind_device" VALUES (23, '09', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:26:34', '2025-08-27 15:55:13', 0, '2', '', '', '', '01', '08', '');
INSERT INTO "public"."wind_device" VALUES (24, '10', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '', '156', '', '', 0, '', '2025-07-31 14:26:52', '2025-07-31 14:26:52', 0, '2', '', '', '', '02', '09', '');
INSERT INTO "public"."wind_device" VALUES (25, '11', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '', '156', '', '', 0, '', '2025-07-31 14:27:00', '2025-07-31 14:27:00', 0, '2', '', '', '', '03', '10', '');
INSERT INTO "public"."wind_device" VALUES (174, '03', 4, 'ACCX', '主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:32:04', '2025-10-24 16:32:04', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (178, '07', 4, 'ACCX', '主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:34:51', '2025-10-24 16:34:51', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (93, '07', 4, 'ACCX', '主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 10:48:32', '2025-09-04 10:48:32', 0, '-', '', '', '', '02', '07', '');
INSERT INTO "public"."wind_device" VALUES (8, '08', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:22:50', '2025-08-27 15:35:32', 0, '2', '', '', '', '04', '07', '');
INSERT INTO "public"."wind_device" VALUES (16, '05', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:24:30', '2025-08-27 15:54:35', 0, '2', '', '', '', '01', '04', '');
INSERT INTO "public"."wind_device" VALUES (18, '06', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '', '156', '', '', 0, '', '2025-07-31 14:25:24', '2025-07-31 14:25:24', 0, '2', '', '', '', '02', '05', '');
INSERT INTO "public"."wind_device" VALUES (19, '07', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '', '156', '', '', 0, '', '2025-07-31 14:25:37', '2025-07-31 14:25:37', 0, '2', '', '', '', '03', '06', '');
INSERT INTO "public"."wind_device" VALUES (22, '08', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '', '156', '', '', 0, '', '2025-07-31 14:26:22', '2025-07-31 14:26:22', 0, '2', '', '', '', '04', '07', '');
INSERT INTO "public"."wind_device" VALUES (5, '05', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:21:19', '2025-08-27 15:34:27', 0, '2', '', '', '', '01', '04', '');
INSERT INTO "public"."wind_device" VALUES (111, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '50', '', '', 0, '', '2025-09-10 16:04:24', '2025-09-10 16:04:24', 1, '-', '', '', '', '50', '04', '');
INSERT INTO "public"."wind_device" VALUES (6, '06', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:21:34', '2025-08-27 15:34:47', 0, '2', '', '', '', '02', '05', '');
INSERT INTO "public"."wind_device" VALUES (7, '07', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:22:24', '2025-08-27 15:35:08', 0, '2', '', '', '', '03', '06', '');
INSERT INTO "public"."wind_device" VALUES (78, '01', 10, 'ULS', '超声波液位计', 3, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:15:15', '2025-08-27 15:48:10', 0, '-', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (112, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '60', '', '', 0, '', '2025-09-10 16:05:02', '2025-09-10 16:05:02', 1, '-', '', '', '', '60', '05', '');
INSERT INTO "public"."wind_device" VALUES (113, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '70', '', '', 0, '', '2025-09-10 16:05:31', '2025-09-10 16:05:31', 1, '-', '', '', '', '70', '06', '');
INSERT INTO "public"."wind_device" VALUES (114, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '80', '', '', 0, '', '2025-09-10 16:06:18', '2025-09-10 16:06:18', 1, '-', '', '', '', '80', '07', '');
INSERT INTO "public"."wind_device" VALUES (115, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '90', '', '', 0, '', '2025-09-10 16:06:56', '2025-09-10 16:06:56', 1, '-', '', '', '', '90', '08', '');
INSERT INTO "public"."wind_device" VALUES (116, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '100', '', '', 0, '', '2025-09-10 16:07:34', '2025-09-10 16:07:34', 1, '-', '', '', '', '100', '09', '');
INSERT INTO "public"."wind_device" VALUES (117, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '110', '', '', 0, '', '2025-09-10 16:08:23', '2025-09-10 16:08:23', 1, '-', '', '', '', '110', '10', '');
INSERT INTO "public"."wind_device" VALUES (61, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '10s', '156', 'x', 'x', 0, '', '2025-08-07 10:09:23', '2025-08-07 10:09:23', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."wind_device" VALUES (129, '01', 9, 'IPC', '网络摄像机', 10, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:24:18', '2025-09-18 17:24:18', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (13, '02', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '', '156', '', '', 0, '', '2025-07-31 14:23:51', '2025-07-31 14:23:51', 0, '2', '', '', '', '02', '01', '');
INSERT INTO "public"."wind_device" VALUES (132, '01', 9, 'IPC', '网络摄像机', 13, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:25:22', '2025-09-18 17:25:22', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (14, '03', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '', '156', '', '', 0, '', '2025-07-31 14:24:20', '2025-07-31 14:24:20', 0, '2', '', '', '', '03', '02', '');
INSERT INTO "public"."wind_device" VALUES (15, '04', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '', '156', '', '', 0, '', '2025-07-31 14:24:26', '2025-07-31 14:24:26', 0, '2', '', '', '', '04', '03', '');
INSERT INTO "public"."wind_device" VALUES (4, '04', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:21:05', '2025-08-27 15:34:09', 0, '2', '', '', '', '04', '03', '');
INSERT INTO "public"."wind_device" VALUES (139, '03', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:48:19', '2025-10-24 15:48:19', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (71, '02', 1, 'HLS', '静力水准仪', 1, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:02', '2025-10-09 09:20:46', 0, '-', '', '', '', '03', '01', 'http://172.16.90.70:8081/device_image/login_bg1.jpeg');
INSERT INTO "public"."wind_device" VALUES (72, '03', 1, 'HLS', '静力水准仪', 1, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:14', '2025-10-09 09:21:02', 0, '-', '', '', '', '02', '02', 'http://172.16.90.70:8081/device_image/login_bg2.jpeg');
INSERT INTO "public"."wind_device" VALUES (136, '04', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:45:05', '2025-10-24 15:45:05', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (142, '06', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:50:05', '2025-10-24 15:50:05', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (143, '07', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:50:34', '2025-10-24 15:50:34', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (146, '10', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:52:53', '2025-10-24 15:52:53', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (102, '17', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:18:17', '2025-09-04 11:18:17', 0, '-', '', '', '', '02', '16', '');
INSERT INTO "public"."wind_device" VALUES (150, '01', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:55:44', '2025-10-24 15:56:23', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (153, '04', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:58:28', '2025-10-24 15:58:28', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (159, '06', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:05:42', '2025-10-24 16:05:42', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (162, '09', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:07:14', '2025-10-24 16:07:14', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (165, '12', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:08:31', '2025-10-24 16:08:31', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (168, '02', 2, 'INSX', '主风向倾角传感器', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:15:24', '2025-10-24 16:16:49', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (171, '01', 13, 'INSY', '垂直主风向倾角传感器', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'x', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:18:13', '2025-10-24 16:18:13', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (156, '03', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:02:30', '2025-10-24 16:02:30', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (26, '12', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '', '156', '', '', 0, '', '2025-07-31 14:27:13', '2025-07-31 14:27:13', 0, '2', '', '', '', '04', '11', '');
INSERT INTO "public"."wind_device" VALUES (76, '01', 10, 'ULS', '超声波液位计', 1, 2, 'TW', '塔架', '', '0', '', '', 0, '', '2025-08-07 10:14:41', '2025-08-07 10:14:41', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."wind_device" VALUES (77, '01', 10, 'ULS', '超声波液位计', 2, 1, 'BL', '叶片', '', '', '', '', 0, '', '2025-08-07 10:14:51', '2025-08-07 10:14:51', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."wind_device" VALUES (94, '02', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:02:00', '2025-09-04 11:02:00', 0, '-', '', '', '', '02', '01', 'http://127.0.0.1:81/device_image\FY\04\TW_20\ACCX\02/accel_1.png');
INSERT INTO "public"."wind_device" VALUES (130, '01', 9, 'IPC', '网络摄像机', 11, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:24:39', '2025-09-18 17:24:39', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (98, '10', 4, 'ACCX', '主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-08T16:00:00.000Z', '2025-09-04 11:10:51', '2025-10-09 16:50:34', 0, '-', '', '', '', '02', '09', 'http://172.16.90.70:8081/device_image/accel_1.png');
INSERT INTO "public"."wind_device" VALUES (172, '01', 4, 'ACCX', '主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:29:05', '2025-10-24 16:29:05', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (176, '05', 4, 'ACCX', '主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:33:37', '2025-10-24 16:33:37', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (182, '11', 4, 'ACCX', '主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:41:28', '2025-10-24 16:41:28', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (186, '15', 4, 'ACCX', '主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:44:20', '2025-10-24 16:44:20', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (144, '08', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:51:20', '2025-10-24 15:51:20', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (147, '11', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:53:29', '2025-10-24 15:53:29', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (148, '12', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:53:49', '2025-10-24 15:53:49', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (190, '19', 4, 'ACCX', '主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:58:39', '2025-10-24 16:58:39', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (28, '27', 4, 'ACCX', '主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:16', '2025-08-27 15:37:23', 0, '2', '', '', '', '02', '27', '');
INSERT INTO "public"."wind_device" VALUES (151, '02', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:56:14', '2025-10-24 15:56:14', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (160, '07', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:06:10', '2025-10-24 16:06:10', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."wind_device" VALUES (163, '10', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:07:46', '2025-10-24 16:07:46', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (166, '01', 3, 'GNSS', 'GNSS', 12, 17, 'NA_156', '机舱顶', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:11:31', '2025-10-24 16:12:56', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (27, '26', 4, 'ACCX', '主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '22', '156', '22', '22', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:07', '2025-08-27 15:37:02', 0, '2', '', '', '', '01', '25', '');
INSERT INTO "public"."wind_device" VALUES (169, '02', 13, 'INSY', '垂直主风向倾角传感器', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:15:55', '2025-10-24 16:16:54', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (179, '08', 12, 'ACCY', '垂直主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:35:12', '2025-10-24 16:35:12', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."wind_device" VALUES (189, '18', 12, 'ACCY', '垂直主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:58:16', '2025-10-24 16:58:16', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (99, '11', 4, 'ACCX', '主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', '', '', 0, '', '2025-09-04 11:11:40', '2025-09-04 11:11:40', 0, '-', '', '', '', '01', '11', '');
INSERT INTO "public"."wind_device" VALUES (36, '15', 4, 'ACCX', '主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:30:03', '2025-08-27 15:41:32', 0, '2', '', '', '', '01', '15', '');
INSERT INTO "public"."wind_device" VALUES (137, '01', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:47:21', '2025-10-24 15:47:21', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (108, '28', 12, 'ACCY', '垂直主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', '', '', 0, '', '2025-09-04 11:28:44', '2025-09-04 11:28:44', 0, '-', '', '', '', '02', '26', '');
INSERT INTO "public"."wind_device" VALUES (107, '25', 12, 'ACCY', '垂直主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', '', '', 0, '', '2025-09-04 11:27:43', '2025-09-04 11:27:43', 0, '-', '', '', '', '01', '24', '');
INSERT INTO "public"."wind_device" VALUES (92, '08', 12, 'ACCY', '垂直主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 10:47:57', '2025-09-04 10:47:57', 0, '-', '', '', '', '02', '06', '');
INSERT INTO "public"."wind_device" VALUES (103, '20', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:19:43', '2025-09-04 11:19:43', 0, '-', '', '', '', '01', '18', '');
INSERT INTO "public"."wind_device" VALUES (43, '01', 12, 'ACCY', '垂直主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:31:35', '2025-07-31 14:31:35', 0, '2', '', '', '', '02', '00', '');
INSERT INTO "public"."wind_device" VALUES (104, '21', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:20:42', '2025-09-04 11:20:42', 0, '-', '', '', '', '04', '21', '');
INSERT INTO "public"."wind_device" VALUES (105, '24', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:21:16', '2025-09-04 11:21:16', 0, '-', '', '', '', '03', '22', '');
INSERT INTO "public"."wind_device" VALUES (33, '09', 12, 'ACCY', '垂直主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:20', '2025-08-27 15:41:03', 0, '2', '', '', '', '02', '08', '');
INSERT INTO "public"."wind_device" VALUES (90, '04', 12, 'ACCY', '垂直主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 10:45:07', '2025-09-04 10:45:07', 0, '-', '', '', '', '01', '02', '');
INSERT INTO "public"."wind_device" VALUES (34, '12', 12, 'ACCY', '垂直主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:33', '2025-08-27 15:41:13', 0, '2', '', '', '', '01', '10', '');
INSERT INTO "public"."wind_device" VALUES (100, '13', 12, 'ACCY', '垂直主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', '', '', 0, '', '2025-09-04 11:14:13', '2025-09-04 11:14:13', 0, '-', '', '', '', '02', '13', '');
INSERT INTO "public"."wind_device" VALUES (101, '16', 12, 'ACCY', '垂直主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', '', '', 0, '', '2025-09-04 11:15:17', '2025-09-04 11:15:17', 0, '-', '', '', '', '01', '14', '');
INSERT INTO "public"."wind_device" VALUES (140, '04', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:48:42', '2025-10-24 15:48:42', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."wind_device" VALUES (41, '05', 12, 'ACCY', '垂直主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:31:01', '2025-08-27 15:50:39', 0, '2', '', '', '', '01', '04', '');
INSERT INTO "public"."wind_device" VALUES (38, '04', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:30:26', '2025-07-31 14:30:26', 0, '2', '', '', '', '01', '02', '');
INSERT INTO "public"."wind_device" VALUES (96, '05', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:04:41', '2025-09-04 11:04:41', 0, '-', '', '', '', '04', '05', '');
INSERT INTO "public"."wind_device" VALUES (40, '08', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:30:45', '2025-07-31 14:30:45', 0, '2', '', '', '', '03', '06', '');
INSERT INTO "public"."wind_device" VALUES (95, '03', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:03:23', '2025-09-04 11:03:23', 0, '-', '', '', '', '01', '03', '');
INSERT INTO "public"."wind_device" VALUES (97, '07', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:05:56', '2025-09-04 11:05:56', 0, '-', '', '', '', '03', '07', '');
INSERT INTO "public"."wind_device" VALUES (133, '01', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:43:17', '2025-12-01 18:45:04', 0, '1234', '', '', '', '01', NULL, 'http://59.110.219.98:81/device_image/制作电脑桌面背景.png');
INSERT INTO "public"."wind_device" VALUES (48, '01', 2, 'INSX', '主风向倾角传感器', 4, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:33:41', '2025-08-27 15:50:00', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."wind_device" VALUES (154, '01', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:01:40', '2025-10-24 16:01:40', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."wind_device" VALUES (157, '04', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:02:55', '2025-10-24 16:02:55', 0, '-', '', '', '', '04', NULL, '');

-- ----------------------------
-- Table structure for wind_device_meta
-- ----------------------------
DROP TABLE IF EXISTS "public"."wind_device_meta";
CREATE TABLE "public"."wind_device_meta" (
  "meta_id" int4 NOT NULL,
  "device_type_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "column_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "display_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "unit" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "ord" int4 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_type_id" int4 NOT NULL DEFAULT 0,
  "ai_enabled" bool NOT NULL DEFAULT true,
  "threshold_config" jsonb NOT NULL DEFAULT '{}'::jsonb
)
;
COMMENT ON COLUMN "public"."wind_device_meta"."device_type_code" IS '传感器类型编号';
COMMENT ON COLUMN "public"."wind_device_meta"."column_name" IS 'TDengine 字段名，例如 x、y、accel、strain、settlement';
COMMENT ON COLUMN "public"."wind_device_meta"."display_name" IS '前端显示名';
COMMENT ON COLUMN "public"."wind_device_meta"."unit" IS '单位';
COMMENT ON COLUMN "public"."wind_device_meta"."ord" IS '字段显示顺序';
COMMENT ON COLUMN "public"."wind_device_meta"."remark" IS '备注';
COMMENT ON COLUMN "public"."wind_device_meta"."device_type_id" IS '传感器类型id';
COMMENT ON TABLE "public"."wind_device_meta" IS '传感器元数据映射表';

-- ----------------------------
-- Records of wind_device_meta
-- ----------------------------
INSERT INTO "public"."wind_device_meta" VALUES (6, 'INSX', 'x', '倾角', '°', 1, '', 2, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (7, 'INSY', 'y', '倾角', '°', 2, '', 13, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (2, 'ACCY', 'accel', '加速度', 'mg', 2, '', 12, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (1, 'ACCX', 'accel', '加速度', 'mg', 1, '', 4, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (5, 'STM', 'strain', '动静态应变', 'με', 1, '', 7, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (8, 'ATS', 'tension', '张拉力', 'kN', 1, '', 5, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (25, 'WPR', 'ti', '湍流强度', ' ', 5, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (28, 'WPR', 'v_sheer', '垂直风切变', ' ', 8, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (29, 'WPR', 'h_sheer', '水平风切变', ' ', 9, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (17, 'GNSS', 'vertical_offset', '竖向偏移量', 'mm', 6, '-', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (18, 'GNSS', 'horizontal_offset', '横向偏移量', 'mm', 7, '-', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (19, 'GNSS', 'ordinate_offset', '纵向偏移量', 'mm', 8, '-', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (23, 'WPR', 'veer', '垂直风向变化率', 'deg/m', 3, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (31, 'WPR', 'direction_high', '上平面处风向', '°', 11, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (9, 'JMT', 'joint', '缝隙', 'mm', 1, '', 6, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (10, 'HLS', 'settlement', '沉降', 'mm', 1, '', 1, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (20, 'IPC', 'image_url', '图像名称', 'px', 1, '', 9, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (12, 'GNSS', 'longitude', '经度', '°', 1, '', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (13, 'GNSS', 'latitude', '纬度', '°', 2, '', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (11, 'ULS', 'height', '液位高度', 'mm', 1, '', 10, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (21, 'WPR', 'd', '测量距离', 'm', 1, '', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (22, 'WPR', 'rws', '视向风速', 'm/s', 2, '', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (24, 'WPR', 'raws', '轴向投影风速', 'm/s', 4, '', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (26, 'WPR', 'hw_shub', '轮毂高度处风速', 'm/s', 6, '', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (30, 'WPR', 'hw_shigh', '上平面处风速', 'm/s', 10, '', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (32, 'WPR', 'hw_slow', '下平面处风速', 'm/s', 12, '', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (33, 'WPR', 'direction_low', '下平面处风向', '°', 13, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (27, 'WPR', 'direction_hub', '轮毂高度处风向', '°', 7, '-', 8, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (14, 'GNSS', 'vertical', '竖向坐标', 'm', 3, '-', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (15, 'GNSS', 'horizontal', '横向坐标', 'm', 4, '-', 3, 't', '{}');
INSERT INTO "public"."wind_device_meta" VALUES (16, 'GNSS', 'ordinate', '纵向坐标', 'm', 4, '-', 3, 't', '{}');

-- ----------------------------
-- Table structure for wind_device_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."wind_device_type";
CREATE TABLE "public"."wind_device_type" (
  "device_type_id" int4 NOT NULL,
  "device_type_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_type_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "sample_unit" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "threshold_config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "ai_enabled" bool NOT NULL DEFAULT true,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."wind_device_type"."threshold_config" IS 'AI 分析阈值配置，JSON 格式';
COMMENT ON TABLE "public"."wind_device_type" IS '传感器类型表，补充 TDengine 超级表映射';

-- ----------------------------
-- Records of wind_device_type
-- ----------------------------
INSERT INTO "public"."wind_device_type" VALUES (2, 'INSX', '主风向倾角传感器', '°', '{}', 't', '2025-07-31 14:15:49', '2025-09-17 15:10:05', 0, '与当地水平面夹角（绝对值）');
INSERT INTO "public"."wind_device_type" VALUES (13, 'INSY', '垂直主风向倾角传感器', '°', '{}', 't', '2025-09-12 09:46:01', '2025-09-17 15:10:22', 0, '与当地水平面夹角（绝对值）');
INSERT INTO "public"."wind_device_type" VALUES (6, 'JMT', '测缝计', 'mm', '{}', 't', '2025-07-31 14:17:04', '2025-09-17 15:10:43', 0, '混凝土裂缝张开宽度');
INSERT INTO "public"."wind_device_type" VALUES (1, 'HLS', '静力水准仪', 'mm', '{}', 't', '2025-07-31 13:43:32', '2025-09-17 15:09:29', 0, '沉降量（绝对值）');
INSERT INTO "public"."wind_device_type" VALUES (10, 'ULS', '超声波液位计', 'mm', '{}', 't', '2025-07-31 14:17:51', '2025-09-17 15:16:35', 0, '基坑水位深度');
INSERT INTO "public"."wind_device_type" VALUES (11, 'HLS2222', '超声', '', '{}', 't', '2025-08-04 10:35:19', '2006-01-02 15:04:05', 1, '-');
INSERT INTO "public"."wind_device_type" VALUES (5, 'ATS', '锚索计', 'kN', '{}', 't', '2025-07-31 14:16:52', '2025-09-17 15:09:03', 0, '索力张拉力');
INSERT INTO "public"."wind_device_type" VALUES (7, 'STM', '应变计', 'με', '{}', 't', '2025-07-31 14:17:15', '2025-09-17 15:13:19', 0, '叶根测点应变值（绝对值）');
INSERT INTO "public"."wind_device_type" VALUES (12, 'ACCY', '垂直主风向加速度计', 'mg', '{}', 't', '2025-09-04 10:31:24', '2025-09-17 15:08:44', 0, '采样50Hz-200Hz');
INSERT INTO "public"."wind_device_type" VALUES (3, 'GNSS', 'GNSS', '', '{}', 't', '2025-07-31 14:16:04', '2025-07-31 14:16:04', 0, '-');
INSERT INTO "public"."wind_device_type" VALUES (9, 'IPC', '网络摄像机', '', '{}', 't', '2025-07-31 14:17:38', '2025-07-31 14:17:38', 0, '-');
INSERT INTO "public"."wind_device_type" VALUES (8, 'WPR', '测风雷达', 'm/s', '{}', 't', '2025-07-31 14:17:26', '2025-08-19 14:05:50', 0, '-');
INSERT INTO "public"."wind_device_type" VALUES (4, 'ACCX', '主风向加速度计', 'mg', '{}', 't', '2025-07-31 14:16:36', '2025-11-24 15:37:21', 0, '采样50Hz-200Hz');

-- ----------------------------
-- Table structure for wind_farm
-- ----------------------------
DROP TABLE IF EXISTS "public"."wind_farm";
CREATE TABLE "public"."wind_farm" (
  "farm_id" int4 NOT NULL,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "province" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "location" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "latitude" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "longitude" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "td_database" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "ai_enabled" bool NOT NULL DEFAULT true
)
;
COMMENT ON COLUMN "public"."wind_farm"."farm_code" IS '风场编号，例如 FY、YS';
COMMENT ON COLUMN "public"."wind_farm"."td_database" IS '对应 TDengine 数据库名，例如 fuyu、yushu';
COMMENT ON COLUMN "public"."wind_farm"."ai_enabled" IS '是否允许 AI Copilot 查询该风场数据';
COMMENT ON TABLE "public"."wind_farm" IS '风场表，补充 TDengine 数据库名和 AI 开关';

-- ----------------------------
-- Records of wind_farm
-- ----------------------------
INSERT INTO "public"."wind_farm" VALUES (1, 'FY', '华电万兴风电场', 'JL', '松原市扶余市', '45.227703', '125.482944', '2025-07-31 14:05:12', 0, '-', 'fuyu', 't');
INSERT INTO "public"."wind_farm" VALUES (2, 'YS', '华电龙岗风电场', 'JL', '长春市榆树市', '44.803879', '126.439690', '2025-07-31 14:06:12', 0, '-', 'yushu', 't');

-- ----------------------------
-- Table structure for wind_structure_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."wind_structure_type";
CREATE TABLE "public"."wind_structure_type" (
  "structure_id" int4 NOT NULL,
  "structure_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "structure_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying
)
;
COMMENT ON TABLE "public"."wind_structure_type" IS '风机结构类型表';

-- ----------------------------
-- Records of wind_structure_type
-- ----------------------------
INSERT INTO "public"."wind_structure_type" VALUES (7, 'TW_114', '塔架钢混转接', '2025-08-27 15:08:53', 0, '1,2.9,0');
INSERT INTO "public"."wind_structure_type" VALUES (8, 'TW_00', '塔架塔底平台', '2025-08-27 15:09:39', 0, '-1,1.1,0');
INSERT INTO "public"."wind_structure_type" VALUES (9, 'FD_00', '基础地基基坑', '2025-08-27 15:10:46', 0, '1,1.1,0');
INSERT INTO "public"."wind_structure_type" VALUES (10, 'CT_06', '索力6米', '2025-08-27 15:11:25', 0, '-1,1.6,0');
INSERT INTO "public"."wind_structure_type" VALUES (14, 'TW_80', '塔架_80米', '2025-08-27 15:20:30', 0, '-1,2.5,0');
INSERT INTO "public"."wind_structure_type" VALUES (15, 'TW_50', '塔架_50米', '2025-08-27 15:20:55', 0, '1,2.2,0');
INSERT INTO "public"."wind_structure_type" VALUES (16, 'TW_20', '塔架_20米', '2025-08-27 15:21:21', 0, '1,1.6,0');
INSERT INTO "public"."wind_structure_type" VALUES (13, 'TW_156', '塔架_钢段顶平台', '2025-08-27 15:19:16', 0, '1,3.7,0');
INSERT INTO "public"."wind_structure_type" VALUES (12, 'BL03_156', '叶片03', '2025-08-27 15:18:28', 0, '-1,3.7,0');
INSERT INTO "public"."wind_structure_type" VALUES (11, 'BL02_156', '叶片02', '2025-08-27 15:17:53', 0, '-1,4.2,0');
INSERT INTO "public"."wind_structure_type" VALUES (17, 'NA_156', '机舱顶', '2025-09-12 15:38:08', 0, '0,4.2,0');
INSERT INTO "public"."wind_structure_type" VALUES (2, 'TW', '塔架', '2025-07-31 14:13:16', 1, '-');
INSERT INTO "public"."wind_structure_type" VALUES (3, 'FD', '基础', '2025-07-31 14:14:32', 1, '-');
INSERT INTO "public"."wind_structure_type" VALUES (4, 'CT', '索力', '2025-07-31 14:14:45', 1, '-');
INSERT INTO "public"."wind_structure_type" VALUES (1, 'BL', '叶片', '2025-07-31 14:13:04', 1, '-');
INSERT INTO "public"."wind_structure_type" VALUES (6, 'BL01_156', '叶片01', '2025-08-27 11:14:41', 0, '1,4.2,0');

-- ----------------------------
-- Table structure for wind_tower
-- ----------------------------
DROP TABLE IF EXISTS "public"."wind_tower";
CREATE TABLE "public"."wind_tower" (
  "tower_id" int4 NOT NULL,
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_id" int4 NOT NULL DEFAULT 0,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "longitude" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "latitude" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "models" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "risk_level" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'normal'::character varying,
  "ai_enabled" bool NOT NULL DEFAULT true
)
;
COMMENT ON COLUMN "public"."wind_tower"."tower_code" IS '风机编号，如 04 表示 4 号风机';
COMMENT ON COLUMN "public"."wind_tower"."risk_level" IS 'AI 风险等级：normal/warning/high/critical';
COMMENT ON TABLE "public"."wind_tower" IS '风机表，补充 AI 风险等级';

-- ----------------------------
-- Records of wind_tower
-- ----------------------------
INSERT INTO "public"."wind_tower" VALUES (5, '07', 1, 'FY', '华电万兴风电场', '2025-07-31 14:11:17', 0, '-', '125.642331', '45.227637', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (2, '03', 1, 'FY', '华电万兴风电场', '2025-07-31 14:10:40', 0, '-', '125.642831', '45.227637', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (9, '01', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:20:27', 0, '-', '126.532532', '44.844675', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (10, '02', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:21:00', 0, '-', '126.532532', '44.844675', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (11, '03', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:21:24', 0, '-', '126.532532', '44.844675', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (13, '06', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:22:06', 0, '-', '126.532532', '44.844675', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (7, '02', 2, 'YS', '华电龙岗风电场', '2025-08-04 10:27:21', 1, '-', '', '', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (8, '09', 2, 'YS', '华电龙岗风电场', '2025-08-22 14:55:50', 1, '-', '126.532532', '44.844675', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (12, '04', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:21:44', 0, '-', '126.532532', '44.844675', '', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (3, '04', 1, 'FY', '华电万兴风电场', '2025-07-31 14:10:53', 0, '-', '125.642931', '45.227637', '1,2,3,4,6,5', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (4, '05', 1, 'FY', '华电万兴风电场', '2025-07-31 14:11:05', 0, '-', '125.642431', '45.227637', '2,3,1', 'normal', 't');
INSERT INTO "public"."wind_tower" VALUES (1, '01', 1, 'FY', '华电万兴风电场', '2025-07-31 13:58:44', 0, '-', '125.642731', '45.227637', '4,5', 'normal', 't');

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_numeric, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_numeric, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_int4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_int4, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_float4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_float4, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_halfvec"(_float8, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_halfvec"(_float8, int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'array_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_int4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_int4, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_numeric, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_numeric, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_float8, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_float8, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_sparsevec"(_float4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_sparsevec"(_float4, int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'array_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_float4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_float4, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_numeric, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_numeric, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_float8, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_float8, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for array_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."array_to_vector"(_int4, int4, bool);
CREATE OR REPLACE FUNCTION "public"."array_to_vector"(_int4, int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'array_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for binary_quantize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."binary_quantize"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."binary_quantize"("public"."halfvec")
  RETURNS "pg_catalog"."bit" AS '$libdir/vector', 'halfvec_binary_quantize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for binary_quantize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."binary_quantize"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."binary_quantize"("public"."vector")
  RETURNS "pg_catalog"."bit" AS '$libdir/vector', 'binary_quantize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for cosine_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."cosine_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."cosine_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_cosine_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for cosine_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."cosine_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."cosine_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_cosine_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for cosine_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."cosine_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."cosine_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'cosine_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec"("public"."halfvec", int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_accum
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_accum"(_float8, "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_accum"(_float8, "public"."halfvec")
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'halfvec_accum'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_add
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_add"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_add"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_add'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_avg
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_avg"(_float8);
CREATE OR REPLACE FUNCTION "public"."halfvec_avg"(_float8)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_avg'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_cmp"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_cmp"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'halfvec_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_combine
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_combine"(_float8, _float8);
CREATE OR REPLACE FUNCTION "public"."halfvec_combine"(_float8, _float8)
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'vector_combine'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_concat
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_concat"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_concat"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_concat'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_eq
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_eq"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_eq"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_eq'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_ge"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_ge"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_gt"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_gt"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_in"(cstring, oid, int4);
CREATE OR REPLACE FUNCTION "public"."halfvec_in"(cstring, oid, int4)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_l2_squared_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_l2_squared_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_l2_squared_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l2_squared_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_le"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_le"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_lt"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_lt"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_mul
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_mul"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_mul"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_mul'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_ne
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_ne"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_ne"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'halfvec_ne'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_negative_inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_negative_inner_product"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_negative_inner_product"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_negative_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_out"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_out"("public"."halfvec")
  RETURNS "pg_catalog"."cstring" AS '$libdir/vector', 'halfvec_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_recv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_recv"(internal, oid, int4);
CREATE OR REPLACE FUNCTION "public"."halfvec_recv"(internal, oid, int4)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_recv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_send
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_send"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_send"("public"."halfvec")
  RETURNS "pg_catalog"."bytea" AS '$libdir/vector', 'halfvec_send'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_spherical_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_spherical_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_spherical_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_spherical_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_sub
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_sub"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."halfvec_sub"("public"."halfvec", "public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_sub'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_to_float4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_to_float4"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec_to_float4"("public"."halfvec", int4, bool)
  RETURNS "pg_catalog"."_float4" AS '$libdir/vector', 'halfvec_to_float4'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_to_sparsevec"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec_to_sparsevec"("public"."halfvec", int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'halfvec_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_to_vector"("public"."halfvec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."halfvec_to_vector"("public"."halfvec", int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'halfvec_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for halfvec_typmod_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."halfvec_typmod_in"(_cstring);
CREATE OR REPLACE FUNCTION "public"."halfvec_typmod_in"(_cstring)
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'halfvec_typmod_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for hamming_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hamming_distance"(bit, bit);
CREATE OR REPLACE FUNCTION "public"."hamming_distance"(bit, bit)
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'hamming_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for hnsw_bit_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnsw_bit_support"(internal);
CREATE OR REPLACE FUNCTION "public"."hnsw_bit_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'hnsw_bit_support'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for hnsw_halfvec_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnsw_halfvec_support"(internal);
CREATE OR REPLACE FUNCTION "public"."hnsw_halfvec_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'hnsw_halfvec_support'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for hnsw_sparsevec_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnsw_sparsevec_support"(internal);
CREATE OR REPLACE FUNCTION "public"."hnsw_sparsevec_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'hnsw_sparsevec_support'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for hnswhandler
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."hnswhandler"(internal);
CREATE OR REPLACE FUNCTION "public"."hnswhandler"(internal)
  RETURNS "pg_catalog"."index_am_handler" AS '$libdir/vector', 'hnswhandler'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."inner_product"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."inner_product"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."inner_product"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."inner_product"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."inner_product"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."inner_product"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for ivfflat_bit_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."ivfflat_bit_support"(internal);
CREATE OR REPLACE FUNCTION "public"."ivfflat_bit_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'ivfflat_bit_support'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for ivfflat_halfvec_support
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."ivfflat_halfvec_support"(internal);
CREATE OR REPLACE FUNCTION "public"."ivfflat_halfvec_support"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/vector', 'ivfflat_halfvec_support'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for ivfflathandler
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."ivfflathandler"(internal);
CREATE OR REPLACE FUNCTION "public"."ivfflathandler"(internal)
  RETURNS "pg_catalog"."index_am_handler" AS '$libdir/vector', 'ivfflathandler'
  LANGUAGE c VOLATILE
  COST 1;

-- ----------------------------
-- Function structure for jaccard_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."jaccard_distance"(bit, bit);
CREATE OR REPLACE FUNCTION "public"."jaccard_distance"(bit, bit)
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'jaccard_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l1_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l1_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l1_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l1_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l1_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l1_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l1_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l1_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l1_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l1_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."l1_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'l1_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l2_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l2_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_distance"("public"."halfvec", "public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l2_distance"("public"."halfvec", "public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l2_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."l2_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'l2_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_norm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_norm"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l2_norm"("public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l2_norm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_norm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_norm"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l2_norm"("public"."halfvec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'halfvec_l2_norm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_normalize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_normalize"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."l2_normalize"("public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'l2_normalize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_normalize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_normalize"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."l2_normalize"("public"."sparsevec")
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec_l2_normalize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for l2_normalize
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."l2_normalize"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."l2_normalize"("public"."halfvec")
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_l2_normalize'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec"("public"."sparsevec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."sparsevec"("public"."sparsevec", int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_cmp"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_cmp"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'sparsevec_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_eq
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_eq"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_eq"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_eq'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_ge"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_ge"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_gt"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_gt"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_in"(cstring, oid, int4);
CREATE OR REPLACE FUNCTION "public"."sparsevec_in"(cstring, oid, int4)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_l2_squared_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_l2_squared_distance"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_l2_squared_distance"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_l2_squared_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_le"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_le"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_lt"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_lt"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_ne
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_ne"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_ne"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'sparsevec_ne'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_negative_inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_negative_inner_product"("public"."sparsevec", "public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_negative_inner_product"("public"."sparsevec", "public"."sparsevec")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'sparsevec_negative_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_out"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_out"("public"."sparsevec")
  RETURNS "pg_catalog"."cstring" AS '$libdir/vector', 'sparsevec_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_recv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_recv"(internal, oid, int4);
CREATE OR REPLACE FUNCTION "public"."sparsevec_recv"(internal, oid, int4)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'sparsevec_recv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_send
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_send"("public"."sparsevec");
CREATE OR REPLACE FUNCTION "public"."sparsevec_send"("public"."sparsevec")
  RETURNS "pg_catalog"."bytea" AS '$libdir/vector', 'sparsevec_send'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_to_halfvec"("public"."sparsevec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."sparsevec_to_halfvec"("public"."sparsevec", int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'sparsevec_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_to_vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_to_vector"("public"."sparsevec", int4, bool);
CREATE OR REPLACE FUNCTION "public"."sparsevec_to_vector"("public"."sparsevec", int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'sparsevec_to_vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for sparsevec_typmod_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."sparsevec_typmod_in"(_cstring);
CREATE OR REPLACE FUNCTION "public"."sparsevec_typmod_in"(_cstring)
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'sparsevec_typmod_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for subvector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."subvector"("public"."halfvec", int4, int4);
CREATE OR REPLACE FUNCTION "public"."subvector"("public"."halfvec", int4, int4)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'halfvec_subvector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for subvector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."subvector"("public"."vector", int4, int4);
CREATE OR REPLACE FUNCTION "public"."subvector"("public"."vector", int4, int4)
  RETURNS "public"."vector" AS '$libdir/vector', 'subvector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector"("public"."vector", int4, bool)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_accum
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_accum"(_float8, "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_accum"(_float8, "public"."vector")
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'vector_accum'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_add
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_add"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_add"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_add'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_avg
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_avg"(_float8);
CREATE OR REPLACE FUNCTION "public"."vector_avg"(_float8)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_avg'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_cmp
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_cmp"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_cmp"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'vector_cmp'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_combine
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_combine"(_float8, _float8);
CREATE OR REPLACE FUNCTION "public"."vector_combine"(_float8, _float8)
  RETURNS "pg_catalog"."_float8" AS '$libdir/vector', 'vector_combine'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_concat
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_concat"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_concat"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_concat'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_dims
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_dims"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_dims"("public"."vector")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'vector_dims'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_dims
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_dims"("public"."halfvec");
CREATE OR REPLACE FUNCTION "public"."vector_dims"("public"."halfvec")
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'halfvec_vector_dims'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_eq
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_eq"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_eq"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_eq'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_ge
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_ge"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_ge"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_ge'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_gt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_gt"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_gt"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_gt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_in"(cstring, oid, int4);
CREATE OR REPLACE FUNCTION "public"."vector_in"(cstring, oid, int4)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_l2_squared_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_l2_squared_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_l2_squared_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_l2_squared_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_le
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_le"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_le"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_le'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_lt
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_lt"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_lt"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_lt'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_mul
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_mul"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_mul"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_mul'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_ne
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_ne"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_ne"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."bool" AS '$libdir/vector', 'vector_ne'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_negative_inner_product
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_negative_inner_product"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_negative_inner_product"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_negative_inner_product'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_norm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_norm"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_norm"("public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_norm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_out"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_out"("public"."vector")
  RETURNS "pg_catalog"."cstring" AS '$libdir/vector', 'vector_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_recv
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_recv"(internal, oid, int4);
CREATE OR REPLACE FUNCTION "public"."vector_recv"(internal, oid, int4)
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_recv'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_send
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_send"("public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_send"("public"."vector")
  RETURNS "pg_catalog"."bytea" AS '$libdir/vector', 'vector_send'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_spherical_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_spherical_distance"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_spherical_distance"("public"."vector", "public"."vector")
  RETURNS "pg_catalog"."float8" AS '$libdir/vector', 'vector_spherical_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_sub
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_sub"("public"."vector", "public"."vector");
CREATE OR REPLACE FUNCTION "public"."vector_sub"("public"."vector", "public"."vector")
  RETURNS "public"."vector" AS '$libdir/vector', 'vector_sub'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_to_float4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_to_float4"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector_to_float4"("public"."vector", int4, bool)
  RETURNS "pg_catalog"."_float4" AS '$libdir/vector', 'vector_to_float4'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_to_halfvec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_to_halfvec"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector_to_halfvec"("public"."vector", int4, bool)
  RETURNS "public"."halfvec" AS '$libdir/vector', 'vector_to_halfvec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_to_sparsevec
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_to_sparsevec"("public"."vector", int4, bool);
CREATE OR REPLACE FUNCTION "public"."vector_to_sparsevec"("public"."vector", int4, bool)
  RETURNS "public"."sparsevec" AS '$libdir/vector', 'vector_to_sparsevec'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for vector_typmod_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."vector_typmod_in"(_cstring);
CREATE OR REPLACE FUNCTION "public"."vector_typmod_in"(_cstring)
  RETURNS "pg_catalog"."int4" AS '$libdir/vector', 'vector_typmod_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_alarm_analysis_id_seq"
OWNED BY "public"."ai_alarm_analysis"."id";
SELECT setval('"public"."ai_alarm_analysis_id_seq"', 7, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_conversation_id_seq"
OWNED BY "public"."ai_conversation"."id";
SELECT setval('"public"."ai_conversation_id_seq"', 12, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_document_chunk_id_seq"
OWNED BY "public"."ai_document_chunk"."id";
SELECT setval('"public"."ai_document_chunk_id_seq"', 492, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_document_id_seq"
OWNED BY "public"."ai_document"."id";
SELECT setval('"public"."ai_document_id_seq"', 19, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_document_parent_chunk_id_seq"
OWNED BY "public"."ai_document_parent_chunk"."id";
SELECT setval('"public"."ai_document_parent_chunk_id_seq"', 220, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_health_report_id_seq"
OWNED BY "public"."ai_health_report"."id";
SELECT setval('"public"."ai_health_report_id_seq"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_kb_domain_id_seq"
OWNED BY "public"."ai_kb_domain"."id";
SELECT setval('"public"."ai_kb_domain_id_seq"', 4, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_kb_member_id_seq"
OWNED BY "public"."ai_kb_member"."id";
SELECT setval('"public"."ai_kb_member_id_seq"', 4, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_knowledge_base_id_seq"
OWNED BY "public"."ai_knowledge_base"."id";
SELECT setval('"public"."ai_knowledge_base_id_seq"', 8, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_llm_call_log_id_seq"
OWNED BY "public"."ai_llm_call_log"."id";
SELECT setval('"public"."ai_llm_call_log_id_seq"', 36, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_maintenance_ticket_draft_id_seq"
OWNED BY "public"."ai_maintenance_ticket_draft"."id";
SELECT setval('"public"."ai_maintenance_ticket_draft_id_seq"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_message_id_seq"
OWNED BY "public"."ai_message"."id";
SELECT setval('"public"."ai_message_id_seq"', 58, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_model_config_id_seq"
OWNED BY "public"."ai_model_config"."id";
SELECT setval('"public"."ai_model_config_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."ai_tool_call_log_id_seq"
OWNED BY "public"."ai_tool_call_log"."id";
SELECT setval('"public"."ai_tool_call_log_id_seq"', 139, true);

-- ----------------------------
-- Indexes structure for table ai_alarm_analysis
-- ----------------------------
CREATE INDEX "idx_ai_alarm_analysis_trace" ON "public"."ai_alarm_analysis" USING btree (
  "trace_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_alarm_analysis
-- ----------------------------
ALTER TABLE "public"."ai_alarm_analysis" ADD CONSTRAINT "ai_alarm_analysis_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_conversation
-- ----------------------------
CREATE INDEX "idx_ai_conversation_user" ON "public"."ai_conversation" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_ai_conversation_user_deleted_updated" ON "public"."ai_conversation" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST,
  "updated_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Primary Key structure for table ai_conversation
-- ----------------------------
ALTER TABLE "public"."ai_conversation" ADD CONSTRAINT "ai_conversation_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_document
-- ----------------------------
CREATE INDEX "idx_ai_document_kb" ON "public"."ai_document" USING btree (
  "kb_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_ai_document_uploaded_by" ON "public"."ai_document" USING btree (
  "uploaded_by" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_ai_document_hash" ON "public"."ai_document" USING btree (
  "kb_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "content_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_document
-- ----------------------------
ALTER TABLE "public"."ai_document" ADD CONSTRAINT "ai_document_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_document_chunk
-- ----------------------------
CREATE INDEX "idx_ai_chunk_document" ON "public"."ai_document_chunk" USING btree (
  "document_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_chunk_embedding" ON "public"."ai_document_chunk" (
  "embedding" "public"."vector_cosine_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_chunk_parent" ON "public"."ai_document_chunk" USING btree (
  "parent_chunk_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_document_chunk
-- ----------------------------
ALTER TABLE "public"."ai_document_chunk" ADD CONSTRAINT "ai_document_chunk_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_document_parent_chunk
-- ----------------------------
CREATE INDEX "idx_ai_parent_chunk_document" ON "public"."ai_document_parent_chunk" USING btree (
  "document_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_document_parent_chunk
-- ----------------------------
ALTER TABLE "public"."ai_document_parent_chunk" ADD CONSTRAINT "ai_document_parent_chunk_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_health_report
-- ----------------------------
CREATE INDEX "idx_ai_health_report_user_time" ON "public"."ai_health_report" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_health_report
-- ----------------------------
ALTER TABLE "public"."ai_health_report" ADD CONSTRAINT "ai_health_report_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_kb_domain
-- ----------------------------
CREATE INDEX "idx_ai_kb_domain_status_sort" ON "public"."ai_kb_domain" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST,
  "sort" "pg_catalog"."int4_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_ai_kb_domain_code" ON "public"."ai_kb_domain" USING btree (
  "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_kb_domain
-- ----------------------------
ALTER TABLE "public"."ai_kb_domain" ADD CONSTRAINT "ai_kb_domain_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_kb_member
-- ----------------------------
CREATE INDEX "idx_ai_kb_member_user" ON "public"."ai_kb_member" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uk_ai_kb_member_user" ON "public"."ai_kb_member" USING btree (
  "kb_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Checks structure for table ai_kb_member
-- ----------------------------
ALTER TABLE "public"."ai_kb_member" ADD CONSTRAINT "ck_ai_kb_member_role" CHECK (role::text = ANY (ARRAY['viewer'::character varying, 'editor'::character varying, 'manager'::character varying]::text[]));

-- ----------------------------
-- Primary Key structure for table ai_kb_member
-- ----------------------------
ALTER TABLE "public"."ai_kb_member" ADD CONSTRAINT "ai_kb_member_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_knowledge_base
-- ----------------------------
CREATE INDEX "idx_ai_kb_owner" ON "public"."ai_knowledge_base" USING btree (
  "owner_user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_kb_type_domain" ON "public"."ai_knowledge_base" USING btree (
  "kb_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "domain_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_kb_visibility" ON "public"."ai_knowledge_base" USING btree (
  "visibility" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Checks structure for table ai_knowledge_base
-- ----------------------------
ALTER TABLE "public"."ai_knowledge_base" ADD CONSTRAINT "ck_ai_kb_type" CHECK (kb_type::text = ANY (ARRAY['personal'::character varying, 'public'::character varying]::text[]));
ALTER TABLE "public"."ai_knowledge_base" ADD CONSTRAINT "ck_ai_kb_visibility" CHECK (visibility::text = ANY (ARRAY['private'::character varying, 'public'::character varying]::text[]));
ALTER TABLE "public"."ai_knowledge_base" ADD CONSTRAINT "ck_ai_kb_personal_owner" CHECK (kb_type::text = 'personal'::text AND owner_user_id IS NOT NULL OR kb_type::text = 'public'::text);
ALTER TABLE "public"."ai_knowledge_base" ADD CONSTRAINT "ck_ai_kb_public_domain" CHECK (kb_type::text = 'public'::text AND domain_id IS NOT NULL OR kb_type::text = 'personal'::text);

-- ----------------------------
-- Primary Key structure for table ai_knowledge_base
-- ----------------------------
ALTER TABLE "public"."ai_knowledge_base" ADD CONSTRAINT "ai_knowledge_base_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_llm_call_log
-- ----------------------------
CREATE INDEX "idx_ai_llm_call_trace" ON "public"."ai_llm_call_log" USING btree (
  "trace_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_llm_call_user_time" ON "public"."ai_llm_call_log" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Primary Key structure for table ai_llm_call_log
-- ----------------------------
ALTER TABLE "public"."ai_llm_call_log" ADD CONSTRAINT "ai_llm_call_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_maintenance_ticket_draft
-- ----------------------------
CREATE INDEX "idx_ai_ticket_draft_trace" ON "public"."ai_maintenance_ticket_draft" USING btree (
  "trace_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_maintenance_ticket_draft
-- ----------------------------
ALTER TABLE "public"."ai_maintenance_ticket_draft" ADD CONSTRAINT "ai_maintenance_ticket_draft_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_message
-- ----------------------------
CREATE INDEX "idx_ai_message_conversation" ON "public"."ai_message" USING btree (
  "conversation_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_ai_message_conversation_deleted_created" ON "public"."ai_message" USING btree (
  "conversation_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_message
-- ----------------------------
ALTER TABLE "public"."ai_message" ADD CONSTRAINT "ai_message_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_model_config
-- ----------------------------
CREATE INDEX "idx_ai_model_config_type_enabled" ON "public"."ai_model_config" USING btree (
  "model_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "enabled" "pg_catalog"."bool_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_model_config
-- ----------------------------
ALTER TABLE "public"."ai_model_config" ADD CONSTRAINT "ai_model_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table ai_tool_call_log
-- ----------------------------
CREATE INDEX "idx_ai_tool_call_trace" ON "public"."ai_tool_call_log" USING btree (
  "trace_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table ai_tool_call_log
-- ----------------------------
ALTER TABLE "public"."ai_tool_call_log" ADD CONSTRAINT "ai_tool_call_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table checkpoint_blobs
-- ----------------------------
CREATE INDEX "checkpoint_blobs_thread_id_idx" ON "public"."checkpoint_blobs" USING btree (
  "thread_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table checkpoint_blobs
-- ----------------------------
ALTER TABLE "public"."checkpoint_blobs" ADD CONSTRAINT "checkpoint_blobs_pkey" PRIMARY KEY ("thread_id", "checkpoint_ns", "channel", "version");

-- ----------------------------
-- Primary Key structure for table checkpoint_migrations
-- ----------------------------
ALTER TABLE "public"."checkpoint_migrations" ADD CONSTRAINT "checkpoint_migrations_pkey" PRIMARY KEY ("v");

-- ----------------------------
-- Indexes structure for table checkpoint_writes
-- ----------------------------
CREATE INDEX "checkpoint_writes_thread_id_idx" ON "public"."checkpoint_writes" USING btree (
  "thread_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table checkpoint_writes
-- ----------------------------
ALTER TABLE "public"."checkpoint_writes" ADD CONSTRAINT "checkpoint_writes_pkey" PRIMARY KEY ("thread_id", "checkpoint_ns", "checkpoint_id", "task_id", "idx");

-- ----------------------------
-- Indexes structure for table checkpoints
-- ----------------------------
CREATE INDEX "checkpoints_thread_id_idx" ON "public"."checkpoints" USING btree (
  "thread_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table checkpoints
-- ----------------------------
ALTER TABLE "public"."checkpoints" ADD CONSTRAINT "checkpoints_pkey" PRIMARY KEY ("thread_id", "checkpoint_ns", "checkpoint_id");

-- ----------------------------
-- Primary Key structure for table wind_camera_record
-- ----------------------------
ALTER TABLE "public"."wind_camera_record" ADD CONSTRAINT "wind_camera_record_pkey" PRIMARY KEY ("record_id");

-- ----------------------------
-- Indexes structure for table wind_device
-- ----------------------------
CREATE INDEX "idx_wind_device_tower_type" ON "public"."wind_device" USING btree (
  "tower_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "device_type_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table wind_device
-- ----------------------------
ALTER TABLE "public"."wind_device" ADD CONSTRAINT "wind_device_pkey" PRIMARY KEY ("device_id");

-- ----------------------------
-- Indexes structure for table wind_device_meta
-- ----------------------------
CREATE INDEX "idx_wind_device_meta_type" ON "public"."wind_device_meta" USING btree (
  "device_type_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "ord" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table wind_device_meta
-- ----------------------------
ALTER TABLE "public"."wind_device_meta" ADD CONSTRAINT "wind_device_meta_pkey" PRIMARY KEY ("meta_id");

-- ----------------------------
-- Primary Key structure for table wind_device_type
-- ----------------------------
ALTER TABLE "public"."wind_device_type" ADD CONSTRAINT "wind_device_type_pkey" PRIMARY KEY ("device_type_id");

-- ----------------------------
-- Indexes structure for table wind_farm
-- ----------------------------
CREATE INDEX "idx_wind_farm_code" ON "public"."wind_farm" USING btree (
  "farm_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table wind_farm
-- ----------------------------
ALTER TABLE "public"."wind_farm" ADD CONSTRAINT "wind_farm_pkey" PRIMARY KEY ("farm_id");

-- ----------------------------
-- Primary Key structure for table wind_structure_type
-- ----------------------------
ALTER TABLE "public"."wind_structure_type" ADD CONSTRAINT "wind_structure_type_pkey" PRIMARY KEY ("structure_id");

-- ----------------------------
-- Indexes structure for table wind_tower
-- ----------------------------
CREATE INDEX "idx_wind_tower_farm_code" ON "public"."wind_tower" USING btree (
  "farm_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "tower_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table wind_tower
-- ----------------------------
ALTER TABLE "public"."wind_tower" ADD CONSTRAINT "wind_tower_pkey" PRIMARY KEY ("tower_id");

-- ----------------------------
-- Foreign Keys structure for table ai_document
-- ----------------------------
ALTER TABLE "public"."ai_document" ADD CONSTRAINT "ai_document_kb_id_fkey" FOREIGN KEY ("kb_id") REFERENCES "public"."ai_knowledge_base" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table ai_document_chunk
-- ----------------------------
ALTER TABLE "public"."ai_document_chunk" ADD CONSTRAINT "ai_document_chunk_document_id_fkey" FOREIGN KEY ("document_id") REFERENCES "public"."ai_document" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."ai_document_chunk" ADD CONSTRAINT "ai_document_chunk_parent_chunk_id_fkey" FOREIGN KEY ("parent_chunk_id") REFERENCES "public"."ai_document_parent_chunk" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table ai_document_parent_chunk
-- ----------------------------
ALTER TABLE "public"."ai_document_parent_chunk" ADD CONSTRAINT "ai_document_parent_chunk_document_id_fkey" FOREIGN KEY ("document_id") REFERENCES "public"."ai_document" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table ai_kb_member
-- ----------------------------
ALTER TABLE "public"."ai_kb_member" ADD CONSTRAINT "ai_kb_member_kb_id_fkey" FOREIGN KEY ("kb_id") REFERENCES "public"."ai_knowledge_base" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table ai_knowledge_base
-- ----------------------------
ALTER TABLE "public"."ai_knowledge_base" ADD CONSTRAINT "ai_knowledge_base_domain_id_fkey" FOREIGN KEY ("domain_id") REFERENCES "public"."ai_kb_domain" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table ai_message
-- ----------------------------
ALTER TABLE "public"."ai_message" ADD CONSTRAINT "ai_message_conversation_id_fkey" FOREIGN KEY ("conversation_id") REFERENCES "public"."ai_conversation" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
