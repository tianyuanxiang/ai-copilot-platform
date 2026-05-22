/*
 Navicat Premium Data Transfer

 Source Server         : 云服务器-8.140.204.47
 Source Server Type    : PostgreSQL
 Source Server Version : 140017
 Source Host           : 8.140.204.47:5432
 Source Catalog        : blade
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140017
 File Encoding         : 65001

 Date: 17/05/2026 12:31:04
*/


-- ----------------------------
-- Sequence structure for blade_input_file_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blade_input_file_id_seq";
CREATE SEQUENCE "public"."blade_input_file_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for blade_job_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blade_job_id_seq";
CREATE SEQUENCE "public"."blade_job_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for blade_job_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blade_job_log_id_seq";
CREATE SEQUENCE "public"."blade_job_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for blade_job_result_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blade_job_result_id_seq";
CREATE SEQUENCE "public"."blade_job_result_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for blade_parameter_set_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blade_parameter_set_id_seq";
CREATE SEQUENCE "public"."blade_parameter_set_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for blade_project_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blade_project_id_seq";
CREATE SEQUENCE "public"."blade_project_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for blade_input_file
-- ----------------------------
DROP TABLE IF EXISTS "public"."blade_input_file";
CREATE TABLE "public"."blade_input_file" (
  "id" int8 NOT NULL DEFAULT nextval('blade_input_file_id_seq'::regclass),
  "project_id" int8 NOT NULL,
  "file_type" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "version_no" int8 NOT NULL DEFAULT 1,
  "file_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "stored_path" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "file_size" int8 NOT NULL DEFAULT 0,
  "is_current" int8 NOT NULL DEFAULT 1,
  "validation_status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'not_checked'::character varying,
  "validation_message" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "status" int8 NOT NULL DEFAULT 1,
  "uploaded_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;

-- ----------------------------
-- Records of blade_input_file
-- ----------------------------
INSERT INTO "public"."blade_input_file" VALUES (6, 1, 'airfoil', 1, 'Airfoil_data.xls', 'uploads\1\airfoil\1778634701188931900_Airfoil_data.xls', 93184, 0, 'passed', '文件上传并校验通过', 0, '2026-05-13 09:11:41.268867', '2026-05-13 09:11:41.268867', '2026-05-13 11:20:13.625613', NULL);
INSERT INTO "public"."blade_input_file" VALUES (7, 1, 'airfoil', 2, 'Airfoil_data.xls', 'D:\GoProject\project_new\blade\backend\gateway\uploads\1\airfoil\1778642413474366400_Airfoil_data.xls', 93184, 1, 'passed', '文件上传并校验通过', 1, '2026-05-13 11:20:13.625613', '2026-05-13 11:20:13.625613', '2026-05-13 11:20:13.625613', NULL);
INSERT INTO "public"."blade_input_file" VALUES (5, 1, 'aerodynamic_config', 1, 'Aerodynamic_configuration.xls', 'uploads\1\aerodynamic_config\1778634690263510300_Aerodynamic_configuration.xls', 26624, 0, 'passed', '文件上传并校验通过', 0, '2026-05-13 09:11:30.34346', '2026-05-13 09:11:30.34346', '2026-05-13 11:20:36.325903', NULL);
INSERT INTO "public"."blade_input_file" VALUES (8, 1, 'aerodynamic_config', 2, 'Aerodynamic_configuration.xls', 'D:\GoProject\project_new\blade\backend\gateway\uploads\1\aerodynamic_config\1778642436252020200_Aerodynamic_configuration.xls', 26624, 1, 'passed', '文件上传并校验通过', 1, '2026-05-13 11:20:36.325903', '2026-05-13 11:20:36.325903', '2026-05-13 11:20:36.325903', NULL);
INSERT INTO "public"."blade_input_file" VALUES (1, 1, 'materials', 1, 'Materials.xlsx', 'uploads\1\materials\1778634405046089500_Materials.xlsx', 9356, 0, 'passed', '文件上传并校验通过', 0, '2026-05-13 09:06:45.1859', '2026-05-13 09:06:45.1859', '2026-05-13 11:20:45.069845', NULL);
INSERT INTO "public"."blade_input_file" VALUES (9, 1, 'materials', 2, 'Materials.xlsx', 'D:\GoProject\project_new\blade\backend\gateway\uploads\1\materials\1778642444996490700_Materials.xlsx', 9356, 1, 'passed', '文件上传并校验通过', 1, '2026-05-13 11:20:45.069845', '2026-05-13 11:20:45.069845', '2026-05-13 11:20:45.069845', NULL);
INSERT INTO "public"."blade_input_file" VALUES (3, 1, 'positioning_line_reference', 1, 'Positioning_line_reference.xlsx', 'uploads\1\positioning_line_reference\1778634507365860000_Positioning_line_reference.xlsx', 13367, 0, 'passed', '文件上传并校验通过', 0, '2026-05-13 09:08:27.457426', '2026-05-13 09:08:27.457426', '2026-05-13 11:20:55.949026', NULL);
INSERT INTO "public"."blade_input_file" VALUES (10, 1, 'positioning_line_reference', 2, 'Positioning_line_reference.xlsx', 'D:\GoProject\project_new\blade\backend\gateway\uploads\1\positioning_line_reference\1778642455878306500_Positioning_line_reference.xlsx', 13367, 1, 'passed', '文件上传并校验通过', 1, '2026-05-13 11:20:55.949026', '2026-05-13 11:20:55.949026', '2026-05-13 11:20:55.949026', NULL);
INSERT INTO "public"."blade_input_file" VALUES (2, 1, 'positioning_line', 1, 'Positioning_line.xlsx', 'uploads\1\positioning_line\1778634452404677500_Positioning_line.xlsx', 22999, 0, 'passed', '文件上传并校验通过', 0, '2026-05-13 09:07:32.496585', '2026-05-13 09:07:32.496585', '2026-05-13 11:21:06.313467', NULL);
INSERT INTO "public"."blade_input_file" VALUES (11, 1, 'positioning_line', 2, 'Positioning_line.xlsx', 'D:\GoProject\project_new\blade\backend\gateway\uploads\1\positioning_line\1778642466225008000_Positioning_line.xlsx', 22999, 1, 'passed', '文件上传并校验通过', 1, '2026-05-13 11:21:06.313467', '2026-05-13 11:21:06.313467', '2026-05-13 11:21:06.313467', NULL);
INSERT INTO "public"."blade_input_file" VALUES (4, 1, 'layup_input', 1, 'Layup_input.xlsx', 'uploads\1\layup_input\1778634524176698300_Layup_input.xlsx', 33131, 0, 'passed', '文件上传并校验通过', 0, '2026-05-13 09:08:44.275736', '2026-05-13 09:08:44.275736', '2026-05-13 11:21:19.445019', NULL);
INSERT INTO "public"."blade_input_file" VALUES (12, 1, 'layup_input', 2, 'Layup_input.xlsx', 'D:\GoProject\project_new\blade\backend\gateway\uploads\1\layup_input\1778642479376958400_Layup_input.xlsx', 33131, 1, 'passed', '文件上传并校验通过', 1, '2026-05-13 11:21:19.445019', '2026-05-13 11:21:19.445019', '2026-05-13 11:21:19.445019', NULL);

-- ----------------------------
-- Table structure for blade_job
-- ----------------------------
DROP TABLE IF EXISTS "public"."blade_job";
CREATE TABLE "public"."blade_job" (
  "id" int8 NOT NULL DEFAULT nextval('blade_job_id_seq'::regclass),
  "job_no" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "project_id" int8 NOT NULL,
  "parameter_set_id" int8 NOT NULL DEFAULT 0,
  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "current_stage" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "error_code" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "error_message" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "started_at" timestamp(6),
  "finished_at" timestamp(6),
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6),
  "job_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."blade_job"."id" IS '任务主键ID';
COMMENT ON COLUMN "public"."blade_job"."job_no" IS '任务编号，未删除任务内唯一';
COMMENT ON COLUMN "public"."blade_job"."project_id" IS '所属项目ID，对应 blade_project.id';
COMMENT ON COLUMN "public"."blade_job"."parameter_set_id" IS '使用的参数集ID，对应 blade_parameter_set.id';
COMMENT ON COLUMN "public"."blade_job"."status" IS '任务状态，例如 pending、running、succeeded、failed、canceled';
COMMENT ON COLUMN "public"."blade_job"."current_stage" IS '当前任务阶段';
COMMENT ON COLUMN "public"."blade_job"."error_code" IS '失败错误码';
COMMENT ON COLUMN "public"."blade_job"."error_message" IS '失败错误信息';
COMMENT ON COLUMN "public"."blade_job"."started_at" IS '任务开始时间';
COMMENT ON COLUMN "public"."blade_job"."finished_at" IS '任务结束时间';
COMMENT ON COLUMN "public"."blade_job"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."blade_job"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."blade_job"."deleted_at" IS '软删除时间，NULL 表示未删除';
COMMENT ON COLUMN "public"."blade_job"."job_name" IS '任务名称，未删除任务内唯一';
COMMENT ON TABLE "public"."blade_job" IS '叶片建模任务表';

-- ----------------------------
-- Records of blade_job
-- ----------------------------
INSERT INTO "public"."blade_job" VALUES (10, 'BLADE-20260515-0001', 1, 2, 'running', 'call_python_engine', '', '', '2026-05-15 17:18:26.253725', NULL, '2026-05-15 17:18:25.786578', '2026-05-15 17:18:26.253725', NULL, '');

-- ----------------------------
-- Table structure for blade_job_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."blade_job_log";
CREATE TABLE "public"."blade_job_log" (
  "id" int8 NOT NULL DEFAULT nextval('blade_job_log_id_seq'::regclass),
  "job_id" int8 NOT NULL,
  "seq_no" int8 NOT NULL DEFAULT 0,
  "stage" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "level" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'info'::character varying,
  "message" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "log_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."blade_job_log"."id" IS '日志主键ID';
COMMENT ON COLUMN "public"."blade_job_log"."job_id" IS '所属任务ID，对应 blade_job.id';
COMMENT ON COLUMN "public"."blade_job_log"."seq_no" IS '任务内日志序号';
COMMENT ON COLUMN "public"."blade_job_log"."stage" IS '日志所属任务阶段';
COMMENT ON COLUMN "public"."blade_job_log"."level" IS '日志级别，例如 info、warn、error';
COMMENT ON COLUMN "public"."blade_job_log"."message" IS '日志内容';
COMMENT ON COLUMN "public"."blade_job_log"."log_time" IS '日志发生时间';
COMMENT ON COLUMN "public"."blade_job_log"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."blade_job_log"."deleted_at" IS '软删除时间，NULL 表示未删除';
COMMENT ON TABLE "public"."blade_job_log" IS '叶片建模任务日志表';

-- ----------------------------
-- Records of blade_job_log
-- ----------------------------
INSERT INTO "public"."blade_job_log" VALUES (1, 1, 1, 'prepare', 'info', '任务已创建，等待执行', '2026-05-11 15:00:02.831549', '2026-05-11 15:00:02.831549', NULL);
INSERT INTO "public"."blade_job_log" VALUES (2, 1, 2, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-11 15:00:03.354905', '2026-05-11 15:00:03.354905', NULL);
INSERT INTO "public"."blade_job_log" VALUES (3, 1, 3, 'prepare', 'info', '任务已重新提交，进入排队状态', '2026-05-11 15:01:50.182983', '2026-05-11 15:01:50.182983', NULL);
INSERT INTO "public"."blade_job_log" VALUES (4, 1, 4, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-11 15:01:50.342518', '2026-05-11 15:01:50.342518', NULL);
INSERT INTO "public"."blade_job_log" VALUES (5, 2, 1, 'prepare', 'info', '任务已创建，等待执行', '2026-05-13 09:12:26.766677', '2026-05-13 09:12:26.766677', '2026-05-13 11:22:02.489129');
INSERT INTO "public"."blade_job_log" VALUES (6, 2, 2, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-13 09:12:27.212753', '2026-05-13 09:12:27.212753', '2026-05-13 11:22:02.489129');
INSERT INTO "public"."blade_job_log" VALUES (7, 3, 1, 'prepare', 'info', 'task created, waiting execution...', '2026-05-13 11:23:36.962432', '2026-05-13 11:23:36.962432', NULL);
INSERT INTO "public"."blade_job_log" VALUES (8, 4, 1, 'prepare', 'info', 'task created, waiting execution...', '2026-05-13 11:26:43.464818', '2026-05-13 11:26:43.464818', '2026-05-13 13:40:31.569739');
INSERT INTO "public"."blade_job_log" VALUES (9, 5, 1, 'prepare', 'info', 'task created, waiting execution...', '2026-05-13 13:42:04.32564', '2026-05-13 13:42:04.32564', NULL);
INSERT INTO "public"."blade_job_log" VALUES (10, 5, 2, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-13 13:42:04.823754', '2026-05-13 13:42:04.823754', NULL);
INSERT INTO "public"."blade_job_log" VALUES (11, 5, 3, 'validate', 'info', 'validate success', '2026-05-13 13:43:27.43485', '2026-05-13 13:43:27.43485', NULL);
INSERT INTO "public"."blade_job_log" VALUES (12, 5, 4, 'geometry', 'info', 'geometry success', '2026-05-13 13:43:27.546627', '2026-05-13 13:43:27.546627', NULL);
INSERT INTO "public"."blade_job_log" VALUES (13, 5, 5, 'mesh', 'info', 'mesh success', '2026-05-13 13:43:27.609925', '2026-05-13 13:43:27.609925', NULL);
INSERT INTO "public"."blade_job_log" VALUES (14, 5, 6, 'assembly', 'info', 'assembly success', '2026-05-13 13:43:27.676861', '2026-05-13 13:43:27.676861', NULL);
INSERT INTO "public"."blade_job_log" VALUES (15, 5, 7, 'export', 'info', 'export success', '2026-05-13 13:43:27.830737', '2026-05-13 13:43:27.830737', NULL);
INSERT INTO "public"."blade_job_log" VALUES (16, 5, 8, 'visualization', 'info', 'visualization success', '2026-05-13 13:43:27.895684', '2026-05-13 13:43:27.895684', NULL);
INSERT INTO "public"."blade_job_log" VALUES (17, 3, 2, 'prepare', 'info', '任务已取消', '2026-05-13 13:49:25.672295', '2026-05-13 13:49:25.672295', NULL);
INSERT INTO "public"."blade_job_log" VALUES (18, 3, 3, 'prepare', 'info', '任务已重新提交，进入排队状态', '2026-05-13 13:49:30.423538', '2026-05-13 13:49:30.423538', NULL);
INSERT INTO "public"."blade_job_log" VALUES (19, 3, 4, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-13 13:49:30.56165', '2026-05-13 13:49:30.56165', NULL);
INSERT INTO "public"."blade_job_log" VALUES (20, 3, 5, 'validate', 'info', 'validate success', '2026-05-13 13:51:10.500337', '2026-05-13 13:51:10.500337', NULL);
INSERT INTO "public"."blade_job_log" VALUES (21, 3, 6, 'geometry', 'info', 'geometry success', '2026-05-13 13:51:10.580956', '2026-05-13 13:51:10.580956', NULL);
INSERT INTO "public"."blade_job_log" VALUES (22, 3, 7, 'mesh', 'info', 'mesh success', '2026-05-13 13:51:10.659981', '2026-05-13 13:51:10.659981', NULL);
INSERT INTO "public"."blade_job_log" VALUES (23, 3, 8, 'assembly', 'info', 'assembly success', '2026-05-13 13:51:10.727022', '2026-05-13 13:51:10.727022', NULL);
INSERT INTO "public"."blade_job_log" VALUES (24, 3, 9, 'export', 'info', 'export success', '2026-05-13 13:51:10.792766', '2026-05-13 13:51:10.792766', NULL);
INSERT INTO "public"."blade_job_log" VALUES (25, 3, 10, 'visualization', 'info', 'visualization success', '2026-05-13 13:51:10.859539', '2026-05-13 13:51:10.859539', NULL);
INSERT INTO "public"."blade_job_log" VALUES (26, 6, 1, 'prepare', 'info', 'task created, waiting execution...', '2026-05-13 17:29:57.557424', '2026-05-13 17:29:57.557424', NULL);
INSERT INTO "public"."blade_job_log" VALUES (27, 6, 2, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-13 17:29:58.081601', '2026-05-13 17:29:58.081601', NULL);
INSERT INTO "public"."blade_job_log" VALUES (28, 6, 3, 'validate', 'info', 'validate success', '2026-05-13 17:31:55.960736', '2026-05-13 17:31:55.960736', NULL);
INSERT INTO "public"."blade_job_log" VALUES (29, 6, 4, 'geometry', 'info', 'geometry success', '2026-05-13 17:31:56.090217', '2026-05-13 17:31:56.090217', NULL);
INSERT INTO "public"."blade_job_log" VALUES (30, 6, 5, 'mesh', 'info', 'mesh success', '2026-05-13 17:31:56.200232', '2026-05-13 17:31:56.200232', NULL);
INSERT INTO "public"."blade_job_log" VALUES (31, 6, 6, 'assembly', 'info', 'assembly success', '2026-05-13 17:31:56.260924', '2026-05-13 17:31:56.260924', NULL);
INSERT INTO "public"."blade_job_log" VALUES (32, 6, 7, 'export', 'info', 'export success', '2026-05-13 17:31:56.336827', '2026-05-13 17:31:56.336827', NULL);
INSERT INTO "public"."blade_job_log" VALUES (33, 6, 8, 'visualization', 'info', 'visualization success', '2026-05-13 17:31:56.404531', '2026-05-13 17:31:56.404531', NULL);
INSERT INTO "public"."blade_job_log" VALUES (34, 10, 1, 'prepare', 'info', 'task created, waiting execution...', '2026-05-15 17:18:25.862507', '2026-05-15 17:18:25.862507', NULL);
INSERT INTO "public"."blade_job_log" VALUES (35, 10, 2, 'call_python_engine', 'info', 'Python 引擎服务调用中', '2026-05-15 17:18:26.320789', '2026-05-15 17:18:26.320789', NULL);
INSERT INTO "public"."blade_job_log" VALUES (36, 10, 3, 'queued', 'info', 'engine job accepted', '2026-05-15 17:18:26.564888', '2026-05-15 17:18:26.564888', NULL);
INSERT INTO "public"."blade_job_log" VALUES (37, 10, 4, 'call_python_engine', 'info', 'Python engine accepted job, waiting for callbacks', '2026-05-15 17:18:26.632704', '2026-05-15 17:18:26.632704', NULL);

-- ----------------------------
-- Table structure for blade_job_result
-- ----------------------------
DROP TABLE IF EXISTS "public"."blade_job_result";
CREATE TABLE "public"."blade_job_result" (
  "id" int8 NOT NULL DEFAULT nextval('blade_job_result_id_seq'::regclass),
  "job_id" int8 NOT NULL,
  "inp_file_path" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "section_property_path" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "summary_json_path" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "node_count" int8 NOT NULL DEFAULT 0,
  "element_count" int8 NOT NULL DEFAULT 0,
  "material_count" int8 NOT NULL DEFAULT 0,
  "element_set_count" int8 NOT NULL DEFAULT 0,
  "generated_at" timestamp(6),
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."blade_job_result"."id" IS '结果主键ID';
COMMENT ON COLUMN "public"."blade_job_result"."job_id" IS '所属任务ID，对应 blade_job.id，未删除结果内唯一';
COMMENT ON COLUMN "public"."blade_job_result"."inp_file_path" IS '生成的 ABAQUS INP 文件路径';
COMMENT ON COLUMN "public"."blade_job_result"."section_property_path" IS '截面属性文件路径';
COMMENT ON COLUMN "public"."blade_job_result"."summary_json_path" IS '任务摘要 JSON 文件路径';
COMMENT ON COLUMN "public"."blade_job_result"."node_count" IS '节点数量';
COMMENT ON COLUMN "public"."blade_job_result"."element_count" IS '单元数量';
COMMENT ON COLUMN "public"."blade_job_result"."material_count" IS '材料数量';
COMMENT ON COLUMN "public"."blade_job_result"."element_set_count" IS '单元集合数量';
COMMENT ON COLUMN "public"."blade_job_result"."generated_at" IS '结果生成时间';
COMMENT ON COLUMN "public"."blade_job_result"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."blade_job_result"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."blade_job_result"."deleted_at" IS '软删除时间，NULL 表示未删除';
COMMENT ON TABLE "public"."blade_job_result" IS '叶片建模任务结果表';

-- ----------------------------
-- Records of blade_job_result
-- ----------------------------
INSERT INTO "public"."blade_job_result" VALUES (1, 5, 'blade_jobs\5\output\blade.inp', 'blade_jobs\5\output\section_property.txt', 'blade_jobs\5\output\summary.json', 213285, 214812, 9, 345, '2026-05-13 13:43:28.045882', '2026-05-13 13:43:28.045882', '2026-05-13 13:43:28.045882', NULL);
INSERT INTO "public"."blade_job_result" VALUES (2, 3, 'blade_jobs\3\output\blade.inp', 'blade_jobs\3\output\section_property.txt', 'blade_jobs\3\output\summary.json', 213285, 214812, 9, 345, '2026-05-13 13:51:10.994851', '2026-05-13 13:51:10.994851', '2026-05-13 13:51:10.994851', NULL);
INSERT INTO "public"."blade_job_result" VALUES (3, 6, 'blade_jobs\6\output\blade.inp', 'blade_jobs\6\output\section_property.txt', 'blade_jobs\6\output\summary.json', 213285, 214812, 9, 345, '2026-05-13 17:31:56.580432', '2026-05-13 17:31:56.580432', '2026-05-13 17:31:56.580432', NULL);

-- ----------------------------
-- Table structure for blade_parameter_set
-- ----------------------------
DROP TABLE IF EXISTS "public"."blade_parameter_set";
CREATE TABLE "public"."blade_parameter_set" (
  "id" int8 NOT NULL DEFAULT nextval('blade_parameter_set_id_seq'::regclass),
  "project_id" int8 NOT NULL,
  "name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "is_template" int8 NOT NULL DEFAULT 0,
  "element_size" float8 NOT NULL DEFAULT 0,
  "tolerance" float8 NOT NULL DEFAULT 0,
  "check_overlap" int8 NOT NULL DEFAULT 0,
  "extra_config" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."blade_parameter_set"."id" IS '参数集主键ID';
COMMENT ON COLUMN "public"."blade_parameter_set"."project_id" IS '所属项目ID，对应 blade_project.id';
COMMENT ON COLUMN "public"."blade_parameter_set"."name" IS '参数集名称，同一项目下未删除参数集内唯一';
COMMENT ON COLUMN "public"."blade_parameter_set"."is_template" IS '是否模板参数集，0 否，1 是';
COMMENT ON COLUMN "public"."blade_parameter_set"."element_size" IS '网格单元尺寸';
COMMENT ON COLUMN "public"."blade_parameter_set"."tolerance" IS '计算容差';
COMMENT ON COLUMN "public"."blade_parameter_set"."check_overlap" IS '是否检查铺层重叠，0 否，1 是';
COMMENT ON COLUMN "public"."blade_parameter_set"."extra_config" IS '扩展配置，JSON 字符串';
COMMENT ON COLUMN "public"."blade_parameter_set"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."blade_parameter_set"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."blade_parameter_set"."deleted_at" IS '软删除时间，NULL 表示未删除';
COMMENT ON TABLE "public"."blade_parameter_set" IS '叶片参数集表';

-- ----------------------------
-- Records of blade_parameter_set
-- ----------------------------
INSERT INTO "public"."blade_parameter_set" VALUES (1, 2, 'test_a', 0, 50, 1e-06, 1, '', '2026-05-11 14:58:06.360267', '2026-05-11 14:58:06.360267', NULL);
INSERT INTO "public"."blade_parameter_set" VALUES (2, 1, 'test', 0, 50, 1e-06, 1, '', '2026-05-13 09:12:04.264011', '2026-05-13 09:12:04.264011', NULL);

-- ----------------------------
-- Table structure for blade_project
-- ----------------------------
DROP TABLE IF EXISTS "public"."blade_project";
CREATE TABLE "public"."blade_project" (
  "id" int8 NOT NULL DEFAULT nextval('blade_project_id_seq'::regclass),
  "project_code" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "project_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "blade_model" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "description" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'active'::character varying,
  "latest_job_status" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."blade_project"."id" IS '项目主键ID';
COMMENT ON COLUMN "public"."blade_project"."project_code" IS '项目编码，未删除项目内唯一';
COMMENT ON COLUMN "public"."blade_project"."project_name" IS '项目名称';
COMMENT ON COLUMN "public"."blade_project"."blade_model" IS '叶片型号';
COMMENT ON COLUMN "public"."blade_project"."description" IS '项目描述';
COMMENT ON COLUMN "public"."blade_project"."status" IS '项目状态，例如 active、inactive';
COMMENT ON COLUMN "public"."blade_project"."latest_job_status" IS '项目最近一次建模任务状态';
COMMENT ON COLUMN "public"."blade_project"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."blade_project"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."blade_project"."deleted_at" IS '软删除时间，NULL 表示未删除';
COMMENT ON TABLE "public"."blade_project" IS '叶片项目表';

-- ----------------------------
-- Records of blade_project
-- ----------------------------
INSERT INTO "public"."blade_project" VALUES (2, 'bbb', 'bbb', 'bbb', '', 'active', 'created', '2026-05-11 14:37:17.888874', '2026-05-11 15:00:02.984014', NULL);
INSERT INTO "public"."blade_project" VALUES (1, 'aaa', 'aaa', 'aaa', '', 'active', 'created', '2026-05-11 14:36:36.679893', '2026-05-15 17:18:25.987207', NULL);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blade_input_file_id_seq"
OWNED BY "public"."blade_input_file"."id";
SELECT setval('"public"."blade_input_file_id_seq"', 12, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blade_job_id_seq"
OWNED BY "public"."blade_job"."id";
SELECT setval('"public"."blade_job_id_seq"', 10, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blade_job_log_id_seq"
OWNED BY "public"."blade_job_log"."id";
SELECT setval('"public"."blade_job_log_id_seq"', 37, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blade_job_result_id_seq"
OWNED BY "public"."blade_job_result"."id";
SELECT setval('"public"."blade_job_result_id_seq"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blade_parameter_set_id_seq"
OWNED BY "public"."blade_parameter_set"."id";
SELECT setval('"public"."blade_parameter_set_id_seq"', 2, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blade_project_id_seq"
OWNED BY "public"."blade_project"."id";
SELECT setval('"public"."blade_project_id_seq"', 2, true);

-- ----------------------------
-- Indexes structure for table blade_input_file
-- ----------------------------
CREATE INDEX "idx_blade_input_file_current_active" ON "public"."blade_input_file" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "file_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "is_current" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL AND status <> 0;
CREATE INDEX "idx_blade_input_file_project_active" ON "public"."blade_input_file" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "file_type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "status" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Primary Key structure for table blade_input_file
-- ----------------------------
ALTER TABLE "public"."blade_input_file" ADD CONSTRAINT "blade_input_file_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table blade_job
-- ----------------------------
CREATE INDEX "idx_blade_job_parameter_set_id_active" ON "public"."blade_job" USING btree (
  "parameter_set_id" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_job_parameter_set_id_active" IS '未删除任务的参数集查询索引';
CREATE INDEX "idx_blade_job_project_created_at_active" ON "public"."blade_job" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamp_ops" DESC NULLS FIRST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_job_project_created_at_active" IS '未删除任务的项目内创建时间排序索引';
CREATE INDEX "idx_blade_job_project_status_active" ON "public"."blade_job" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_job_project_status_active" IS '未删除任务的项目状态查询索引';
CREATE UNIQUE INDEX "ux_blade_job_job_name_active" ON "public"."blade_job" USING btree (
  "job_name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX "ux_blade_job_job_no_active" ON "public"."blade_job" USING btree (
  "job_no" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."ux_blade_job_job_no_active" IS '未删除任务的任务编号唯一索引';

-- ----------------------------
-- Primary Key structure for table blade_job
-- ----------------------------
ALTER TABLE "public"."blade_job" ADD CONSTRAINT "blade_job_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table blade_job_log
-- ----------------------------
CREATE INDEX "idx_blade_job_log_job_seq_active" ON "public"."blade_job_log" USING btree (
  "job_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "seq_no" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_job_log_job_seq_active" IS '未删除任务日志的任务序号查询索引';
CREATE INDEX "idx_blade_job_log_job_time_active" ON "public"."blade_job_log" USING btree (
  "job_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "log_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_job_log_job_time_active" IS '未删除任务日志的任务时间查询索引';

-- ----------------------------
-- Primary Key structure for table blade_job_log
-- ----------------------------
ALTER TABLE "public"."blade_job_log" ADD CONSTRAINT "blade_job_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table blade_job_result
-- ----------------------------
CREATE UNIQUE INDEX "ux_blade_job_result_job_id_active" ON "public"."blade_job_result" USING btree (
  "job_id" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."ux_blade_job_result_job_id_active" IS '未删除任务结果的任务 ID 唯一索引';

-- ----------------------------
-- Primary Key structure for table blade_job_result
-- ----------------------------
ALTER TABLE "public"."blade_job_result" ADD CONSTRAINT "blade_job_result_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table blade_parameter_set
-- ----------------------------
CREATE INDEX "idx_blade_parameter_set_project_id_active" ON "public"."blade_parameter_set" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_parameter_set_project_id_active" IS '未删除参数集的项目查询索引';
CREATE UNIQUE INDEX "ux_blade_parameter_set_project_name_active" ON "public"."blade_parameter_set" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."ux_blade_parameter_set_project_name_active" IS '未删除参数集的项目内名称唯一索引';

-- ----------------------------
-- Primary Key structure for table blade_parameter_set
-- ----------------------------
ALTER TABLE "public"."blade_parameter_set" ADD CONSTRAINT "blade_parameter_set_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table blade_project
-- ----------------------------
CREATE INDEX "idx_blade_project_status_active" ON "public"."blade_project" USING btree (
  "status" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_project_status_active" IS '未删除项目的状态查询索引';
CREATE INDEX "idx_blade_project_updated_at_active" ON "public"."blade_project" USING btree (
  "updated_at" "pg_catalog"."timestamp_ops" DESC NULLS FIRST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."idx_blade_project_updated_at_active" IS '未删除项目的更新时间排序索引';
CREATE UNIQUE INDEX "ux_blade_project_project_code_active" ON "public"."blade_project" USING btree (
  "project_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;
COMMENT ON INDEX "public"."ux_blade_project_project_code_active" IS '未删除项目的项目编码唯一索引';

-- ----------------------------
-- Primary Key structure for table blade_project
-- ----------------------------
ALTER TABLE "public"."blade_project" ADD CONSTRAINT "blade_project_pkey" PRIMARY KEY ("id");
