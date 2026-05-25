-- 风机混塔智能运维 AI Copilot 一期建表脚本
-- 说明：
-- 1. 业务基础表来自 deploy/sql/wind_data.sql 的原始表结构，统一改名为 wind_*。
-- 2. 知识库相关表继续沿用 ai_copilot.sql 中的 ai_knowledge_base、ai_document、ai_document_parent_chunk、ai_document_chunk 等表。
-- 3. 本脚本只新增风机业务元数据表和 AI 应用脚手架表，不修改 RabbitMQ/RocketMQ/TDengine 采集链路。

CREATE TABLE IF NOT EXISTS "public"."wind_farm" (
  "farm_id" int4 NOT NULL,
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "farm_name" varchar(128) NOT NULL DEFAULT '',
  "province" varchar(32) NOT NULL DEFAULT '',
  "location" varchar(255) NOT NULL DEFAULT '',
  "latitude" varchar(64) NOT NULL DEFAULT '',
  "longitude" varchar(64) NOT NULL DEFAULT '',
  "td_database" varchar(64) NOT NULL DEFAULT '',
  "ai_enabled" bool NOT NULL DEFAULT true,
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) NOT NULL DEFAULT '-',
  CONSTRAINT "wind_farm_pkey" PRIMARY KEY ("farm_id")
);
COMMENT ON TABLE "public"."wind_farm" IS '风场表，复用原 farms 表字段，补充 TDengine 数据库名和 AI 开关';
COMMENT ON COLUMN "public"."wind_farm"."farm_code" IS '风场编号，例如 FY、YS';
COMMENT ON COLUMN "public"."wind_farm"."td_database" IS '对应 TDengine 数据库名，例如 fuyu、yushu';
COMMENT ON COLUMN "public"."wind_farm"."ai_enabled" IS '是否允许 AI Copilot 查询该风场数据';

CREATE TABLE IF NOT EXISTS "public"."wind_turbine" (
  "tower_id" int4 NOT NULL,
  "tower_code" varchar(32) NOT NULL DEFAULT '',
  "farm_id" int4 NOT NULL DEFAULT 0,
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "farm_name" varchar(128) NOT NULL DEFAULT '',
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) NOT NULL DEFAULT '-',
  "longitude" varchar(64) NOT NULL DEFAULT '',
  "latitude" varchar(64) NOT NULL DEFAULT '',
  "models" varchar(255) NOT NULL DEFAULT '',
  "risk_level" varchar(32) NOT NULL DEFAULT 'normal',
  "ai_enabled" bool NOT NULL DEFAULT true,
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_turbine_pkey" PRIMARY KEY ("tower_id")
);
COMMENT ON TABLE "public"."wind_turbine" IS '风机表，复用原 towers 表字段，补充 AI 风险等级';
COMMENT ON COLUMN "public"."wind_turbine"."tower_code" IS '风机编号，如 04 表示 4 号风机';
COMMENT ON COLUMN "public"."wind_turbine"."risk_level" IS 'AI 风险等级：normal/warning/high/critical';

CREATE TABLE IF NOT EXISTS "public"."wind_device_type" (
  "device_type_id" int4 NOT NULL,
  "device_type_code" varchar(32) NOT NULL DEFAULT '',
  "device_type_name" varchar(128) NOT NULL DEFAULT '',
  "sample_unit" varchar(32) NOT NULL DEFAULT '',
  "td_stable" varchar(64) NOT NULL DEFAULT '',
  "threshold_config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "ai_enabled" bool NOT NULL DEFAULT true,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) NOT NULL DEFAULT '',
  CONSTRAINT "wind_device_type_pkey" PRIMARY KEY ("device_type_id")
);
COMMENT ON TABLE "public"."wind_device_type" IS '传感器类型表，复用原 device_type 表字段，补充 TDengine 超级表映射';
COMMENT ON COLUMN "public"."wind_device_type"."td_stable" IS 'TDengine 超级表名，例如 inclinometer、accel、strain';
COMMENT ON COLUMN "public"."wind_device_type"."threshold_config" IS 'AI 分析阈值配置，JSON 格式';

CREATE TABLE IF NOT EXISTS "public"."wind_structure_type" (
  "structure_id" int4 NOT NULL,
  "structure_code" varchar(32) NOT NULL DEFAULT '',
  "structure_name" varchar(32) NOT NULL DEFAULT '',
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) NOT NULL DEFAULT '-',
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_structure_type_pkey" PRIMARY KEY ("structure_id")
);
COMMENT ON TABLE "public"."wind_structure_type" IS '风机结构类型表，复用原 structure_type 表字段';

CREATE TABLE IF NOT EXISTS "public"."wind_device" (
  "device_id" int4 NOT NULL,
  "device_code" varchar(32) NOT NULL DEFAULT '',
  "device_type_id" int4 NOT NULL DEFAULT 1,
  "device_type_code" varchar(32) NOT NULL DEFAULT '',
  "device_type_name" varchar(128) NOT NULL DEFAULT '',
  "tower_id" int4 NOT NULL DEFAULT 0,
  "structure_id" int4 NOT NULL DEFAULT 0,
  "structure_code" varchar(32) NOT NULL DEFAULT '',
  "structure_name" varchar(128) NOT NULL DEFAULT '',
  "frequency" varchar(32) NOT NULL DEFAULT '',
  "height" numeric(10,2) NOT NULL DEFAULT 0,
  "longitude" varchar(64) NOT NULL DEFAULT '',
  "latitude" varchar(64) NOT NULL DEFAULT '',
  "status" int2 NOT NULL DEFAULT 0,
  "install_date" date NULL,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) NOT NULL DEFAULT '',
  "rtsp_url" varchar(500) NOT NULL DEFAULT '',
  "video_url" varchar(500) NOT NULL DEFAULT '',
  "switch_ip" varchar(64) NOT NULL DEFAULT '',
  "location_code" varchar(32) NOT NULL DEFAULT '',
  "install_image" varchar(500) NOT NULL DEFAULT '',
  "td_stable" varchar(64) NOT NULL DEFAULT '',
  "threshold_config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "ai_enabled" bool NOT NULL DEFAULT true,
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_device_pkey" PRIMARY KEY ("device_id")
);
COMMENT ON TABLE "public"."wind_device" IS '传感器实例表，复用原 device 表字段，补充 AI 查询映射';
COMMENT ON COLUMN "public"."wind_device"."device_code" IS '传感器编号，从 00 开始，对应 TDengine device_channel 时需按旧系统规则转换';
COMMENT ON COLUMN "public"."wind_device"."td_stable" IS '可覆盖设备类型的 TDengine 超级表名';

CREATE TABLE IF NOT EXISTS "public"."wind_device_meta" (
  "meta_id" int4 NOT NULL,
  "device_type_code" varchar(32) NOT NULL DEFAULT '',
  "column_name" varchar(64) NOT NULL DEFAULT '',
  "display_name" varchar(128) NOT NULL DEFAULT '',
  "unit" varchar(32) NOT NULL DEFAULT '',
  "ord" int4 NOT NULL DEFAULT 0,
  "remark" varchar(500) NOT NULL DEFAULT '',
  "device_type_id" int4 NOT NULL DEFAULT 0,
  "ai_enabled" bool NOT NULL DEFAULT true,
  "threshold_config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_device_meta_pkey" PRIMARY KEY ("meta_id")
);
COMMENT ON TABLE "public"."wind_device_meta" IS '传感器字段元数据表，复用原 device_meta 表字段，作为 TDengine 字段白名单';
COMMENT ON COLUMN "public"."wind_device_meta"."column_name" IS 'TDengine 字段名，例如 x、y、accel、strain、settlement';

CREATE TABLE IF NOT EXISTS "public"."wind_camera_record" (
  "record_id" int4 NOT NULL,
  "record_name" varchar(128) NOT NULL DEFAULT '',
  "capture_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "tower_code" varchar(32) NOT NULL DEFAULT '',
  "device_code" varchar(32) NOT NULL DEFAULT '',
  "image_url" varchar(255) NOT NULL DEFAULT '',
  "is_delete" int2 NOT NULL DEFAULT 0,
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_camera_record_pkey" PRIMARY KEY ("record_id")
);
COMMENT ON TABLE "public"."wind_camera_record" IS '网络摄像机抓拍记录表，复用原 camera_record 表字段';

CREATE TABLE IF NOT EXISTS "public"."wind_model" (
  "model_id" int4 NOT NULL,
  "model_code" varchar(32) NOT NULL DEFAULT '',
  "model_name" varchar(32) NOT NULL DEFAULT '',
  "display_name" varchar(128) NOT NULL DEFAULT '',
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(500) NOT NULL DEFAULT '-',
  "is_delete" int2 NOT NULL DEFAULT 0,
  "ai_enabled" bool NOT NULL DEFAULT true,
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_model_pkey" PRIMARY KEY ("model_id")
);
COMMENT ON TABLE "public"."wind_model" IS '算法模型表，复用原 models 表字段';

CREATE TABLE IF NOT EXISTS "public"."wind_model_device_map" (
  "map_id" int4 NOT NULL,
  "model_id" int4 NOT NULL DEFAULT 1,
  "device_type_id" int4 NOT NULL DEFAULT 1,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(500) NOT NULL DEFAULT '',
  "is_delete" int2 NOT NULL DEFAULT 0,
  CONSTRAINT "wind_model_device_map_pkey" PRIMARY KEY ("map_id")
);
COMMENT ON TABLE "public"."wind_model_device_map" IS '模型设备映射表，复用原 model_device_map 表字段';

CREATE TABLE IF NOT EXISTS "public"."wind_model_record" (
  "record_id" int4 NOT NULL,
  "run_id" varchar(64) NOT NULL DEFAULT '',
  "model_code" varchar(32) NOT NULL DEFAULT '',
  "model_name" varchar(32) NOT NULL DEFAULT '',
  "display_name" varchar(128) NOT NULL DEFAULT '',
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "tower_code" varchar(32) NOT NULL DEFAULT '',
  "image_url" varchar(500) NOT NULL DEFAULT '',
  "gif_url" varchar(500) NOT NULL DEFAULT '',
  "txt_url" varchar(500) NOT NULL DEFAULT '',
  "started_at" timestamp(0) NULL,
  "completed_at" timestamp(0) NULL,
  "duration_sec" varchar(32) NOT NULL DEFAULT '',
  "remark" varchar(500) NOT NULL DEFAULT '',
  "is_delete" int2 NOT NULL DEFAULT 0,
  "data_start_at" varchar(64) NOT NULL DEFAULT '',
  "data_end_at" varchar(64) NOT NULL DEFAULT '',
  "extra_json" jsonb NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT "wind_model_record_pkey" PRIMARY KEY ("record_id")
);
COMMENT ON TABLE "public"."wind_model_record" IS '算法模型解析记录表，复用原 model_records 表字段';

CREATE TABLE IF NOT EXISTS "public"."ai_alarm_analysis" (
  "id" bigserial PRIMARY KEY,
  "user_id" int8 NOT NULL DEFAULT 0,
  "trace_id" varchar(64) NOT NULL DEFAULT '',
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "tower_code" varchar(32) NOT NULL DEFAULT '',
  "alarm_code" varchar(32) NOT NULL DEFAULT '',
  "title" varchar(255) NOT NULL DEFAULT '',
  "content" text NOT NULL DEFAULT '',
  "evidence" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "status" varchar(32) NOT NULL DEFAULT 'draft',
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE "public"."ai_alarm_analysis" IS 'AI 告警分析草稿表，一期只保存脚手架结果';

CREATE TABLE IF NOT EXISTS "public"."ai_health_report" (
  "id" bigserial PRIMARY KEY,
  "user_id" int8 NOT NULL DEFAULT 0,
  "trace_id" varchar(64) NOT NULL DEFAULT '',
  "report_type" varchar(32) NOT NULL DEFAULT 'health',
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "tower_code" varchar(32) NOT NULL DEFAULT '',
  "start_time" timestamp(0) NULL,
  "end_time" timestamp(0) NULL,
  "title" varchar(255) NOT NULL DEFAULT '',
  "content" text NOT NULL DEFAULT '',
  "evidence" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "status" varchar(32) NOT NULL DEFAULT 'draft',
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE "public"."ai_health_report" IS 'AI 健康报告草稿表，后续扩展日报、周报、单机报告和故障复盘';

CREATE TABLE IF NOT EXISTS "public"."ai_maintenance_ticket_draft" (
  "id" bigserial PRIMARY KEY,
  "user_id" int8 NOT NULL DEFAULT 0,
  "trace_id" varchar(64) NOT NULL DEFAULT '',
  "farm_code" varchar(32) NOT NULL DEFAULT '',
  "tower_code" varchar(32) NOT NULL DEFAULT '',
  "alarm_code" varchar(32) NOT NULL DEFAULT '',
  "title" varchar(255) NOT NULL DEFAULT '',
  "content" text NOT NULL DEFAULT '',
  "evidence" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "status" varchar(32) NOT NULL DEFAULT 'draft',
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE "public"."ai_maintenance_ticket_draft" IS 'AI 维修工单草稿表，一期不自动提交工单';

CREATE INDEX IF NOT EXISTS "idx_wind_farm_code" ON "public"."wind_farm" ("farm_code");
CREATE INDEX IF NOT EXISTS "idx_wind_turbine_farm_code" ON "public"."wind_turbine" ("farm_code", "tower_code");
CREATE INDEX IF NOT EXISTS "idx_wind_device_tower_type" ON "public"."wind_device" ("tower_id", "device_type_code");
CREATE INDEX IF NOT EXISTS "idx_wind_device_meta_type" ON "public"."wind_device_meta" ("device_type_code", "ord");
CREATE INDEX IF NOT EXISTS "idx_ai_alarm_analysis_trace" ON "public"."ai_alarm_analysis" ("trace_id");
CREATE INDEX IF NOT EXISTS "idx_ai_health_report_user_time" ON "public"."ai_health_report" ("user_id", "created_at");
CREATE INDEX IF NOT EXISTS "idx_ai_ticket_draft_trace" ON "public"."ai_maintenance_ticket_draft" ("trace_id");

-- 常用 TDengine 映射初始化。实际数据可按 wind_data.sql 的 INSERT 继续迁移。
INSERT INTO "public"."wind_device_type"
  ("device_type_id", "device_type_code", "device_type_name", "sample_unit", "td_stable", "remark")
VALUES
  (1, 'HLS', '静力水准仪', 'mm', 'hydrostatic', '沉降量'),
  (2, 'INSX', '主风向倾角传感器', '°', 'inclinometer', '倾角 x'),
  (3, 'GNSS', 'GNSS', '', 'gnss', '位移坐标'),
  (4, 'ACCX', '主风向加速度计', 'mg', 'accel', '加速度 x'),
  (5, 'ATS', '锚索计', 'kN', 'tension', '张拉力'),
  (6, 'JMT', '测缝计', 'mm', 'joint_meter', '缝隙'),
  (7, 'STM', '应变计', 'με', 'strain', '动静态应变'),
  (8, 'WPR', '测风雷达', 'm/s', 'radar', '激光测风雷达'),
  (9, 'IPC', '网络摄像机', '', 'camera_record', '摄像机'),
  (10, 'ULS', '超声波液位计', 'mm', 'ultrasonic_level', '液位'),
  (12, 'ACCY', '垂直主风向加速度计', 'mg', 'accel', '加速度 y'),
  (13, 'INSY', '垂直主风向倾角传感器', '°', 'inclinometer', '倾角 y')
ON CONFLICT ("device_type_id") DO NOTHING;
