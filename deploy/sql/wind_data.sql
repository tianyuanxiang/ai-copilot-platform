/*
 Navicat Premium Data Transfer

 Source Server         : 59.110.219.98-5432
 Source Server Type    : PostgreSQL
 Source Server Version : 140020
 Source Host           : 59.110.219.98:5432
 Source Catalog        : fuyu
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140020
 File Encoding         : 65001

 Date: 24/05/2026 22:26:19
*/


-- ----------------------------
-- Sequence structure for alarm_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."alarm_type_id_seq";
CREATE SEQUENCE "public"."alarm_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for camera_record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."camera_record_id_seq";
CREATE SEQUENCE "public"."camera_record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for device_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."device_id_seq";
CREATE SEQUENCE "public"."device_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for device_meta_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."device_meta_id_seq";
CREATE SEQUENCE "public"."device_meta_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for device_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."device_type_id_seq";
CREATE SEQUENCE "public"."device_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for farm_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."farm_id_seq";
CREATE SEQUENCE "public"."farm_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for map_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."map_id_seq";
CREATE SEQUENCE "public"."map_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for models_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."models_id_seq";
CREATE SEQUENCE "public"."models_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."record_id_seq";
CREATE SEQUENCE "public"."record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for remote_switch_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."remote_switch_id_seq";
CREATE SEQUENCE "public"."remote_switch_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for structure_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."structure_id_seq";
CREATE SEQUENCE "public"."structure_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for switch_asset_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."switch_asset_id_seq";
CREATE SEQUENCE "public"."switch_asset_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_depart_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_depart_id_seq";
CREATE SEQUENCE "public"."sys_depart_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_login_log_record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_login_log_record_id_seq";
CREATE SEQUENCE "public"."sys_login_log_record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_menu_menu_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_menu_menu_id_seq";
CREATE SEQUENCE "public"."sys_menu_menu_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_operate_log_record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_operate_log_record_id_seq";
CREATE SEQUENCE "public"."sys_operate_log_record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_role_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_role_id_seq";
CREATE SEQUENCE "public"."sys_role_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_users_uid_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_users_uid_seq";
CREATE SEQUENCE "public"."sys_users_uid_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for tower_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."tower_id_seq";
CREATE SEQUENCE "public"."tower_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for alarm_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."alarm_type";
CREATE TABLE "public"."alarm_type" (
  "id" int4 NOT NULL DEFAULT nextval('alarm_type_id_seq'::regclass),
  "alarm_code" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "alarm_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."alarm_type"."id" IS '类型ID';
COMMENT ON COLUMN "public"."alarm_type"."alarm_code" IS '告警类型编号';
COMMENT ON COLUMN "public"."alarm_type"."alarm_name" IS '告警类型名称';
COMMENT ON COLUMN "public"."alarm_type"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."alarm_type"."updated_at" IS '记录修改时间';
COMMENT ON COLUMN "public"."alarm_type"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."alarm_type"."remark" IS '备注';
COMMENT ON TABLE "public"."alarm_type" IS '告警类型表';

-- ----------------------------
-- Records of alarm_type
-- ----------------------------
INSERT INTO "public"."alarm_type" VALUES (1, '10000', '设备离线报警', '2025-07-21 16:29:14', '2025-07-21 16:29:14', 0, '1');
INSERT INTO "public"."alarm_type" VALUES (2, '10001', '索力报警', '2025-07-21 16:29:19', '2025-07-21 16:29:19', 0, '2');
INSERT INTO "public"."alarm_type" VALUES (3, '10002', '缝隙报警', '2025-07-22 14:51:00', '2025-07-22 14:51:00', 0, '');
INSERT INTO "public"."alarm_type" VALUES (4, '10003', '液位报警', '2025-07-22 14:51:15', '2025-07-22 14:51:15', 0, '');
INSERT INTO "public"."alarm_type" VALUES (5, '10004', '沉降报警', '2025-07-22 14:51:24', '2025-07-22 14:51:24', 0, '');
INSERT INTO "public"."alarm_type" VALUES (6, '10086', '塔身基频报警', '2025-07-22 14:51:33', '2025-07-22 14:51:33', 0, '');

-- ----------------------------
-- Table structure for camera_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."camera_record";
CREATE TABLE "public"."camera_record" (
  "record_id" int4 NOT NULL DEFAULT nextval('camera_record_id_seq'::regclass),
  "record_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "capture_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "image_url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."camera_record"."record_id" IS '网络摄像机抓拍记录id';
COMMENT ON COLUMN "public"."camera_record"."record_name" IS '网络摄像机抓拍记录名称';
COMMENT ON COLUMN "public"."camera_record"."capture_time" IS '抓拍时间';
COMMENT ON COLUMN "public"."camera_record"."farm_code" IS '所属风场编号，例: FY';
COMMENT ON COLUMN "public"."camera_record"."tower_code" IS '风机编号，如 “04” 4号风机';
COMMENT ON COLUMN "public"."camera_record"."device_code" IS '传感器编号, 从00开始(通道编号)';
COMMENT ON COLUMN "public"."camera_record"."image_url" IS '图片地址';
COMMENT ON COLUMN "public"."camera_record"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON TABLE "public"."camera_record" IS '网络摄像机抓拍记录表';

-- ----------------------------
-- Records of camera_record
-- ----------------------------
INSERT INTO "public"."camera_record" VALUES (822, 'FY_01_IPC_01_20250905180000000', '2025-09-05 18:00:00', 'FY', '01', '01', 'http://59.110.219.98:9081/5,f6069d2549.jpg', 0);
INSERT INTO "public"."camera_record" VALUES (823, 'FY_03_IPC_01_20250905180002000', '2025-09-05 18:00:02', 'FY', '03', '01', 'http://59.110.219.98:9081/2,f7313de26f.jpg', 0);





-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS "public"."device";
CREATE TABLE "public"."device" (
  "device_id" int4 NOT NULL DEFAULT nextval('device_id_seq'::regclass),
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
COMMENT ON COLUMN "public"."device"."device_id" IS '传感器实例ID';
COMMENT ON COLUMN "public"."device"."device_code" IS '传感器编号, 从00开始(通道编号)';
COMMENT ON COLUMN "public"."device"."device_type_id" IS '传感器类型表的主键id';
COMMENT ON COLUMN "public"."device"."device_type_code" IS '传感器类型编号，如: HLS';
COMMENT ON COLUMN "public"."device"."device_type_name" IS '传感器类型名称，如: 静力水准仪';
COMMENT ON COLUMN "public"."device"."tower_id" IS '所属风机id';
COMMENT ON COLUMN "public"."device"."structure_id" IS '风机结构类型id';
COMMENT ON COLUMN "public"."device"."structure_code" IS '结构编号，例如: BL';
COMMENT ON COLUMN "public"."device"."structure_name" IS '结构名称，例如: 叶片';
COMMENT ON COLUMN "public"."device"."frequency" IS '采样频率，单位为HZ';
COMMENT ON COLUMN "public"."device"."height" IS '传感器高度，单位为m，精度为小数点后2位。';
COMMENT ON COLUMN "public"."device"."longitude" IS '传感器经度';
COMMENT ON COLUMN "public"."device"."latitude" IS '传感器纬度';
COMMENT ON COLUMN "public"."device"."status" IS '状态：在线0 离线 1';
COMMENT ON COLUMN "public"."device"."install_date" IS '安装日期';
COMMENT ON COLUMN "public"."device"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."device"."updated_at" IS '记录修改时间';
COMMENT ON COLUMN "public"."device"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."device"."remark" IS '备注';
COMMENT ON COLUMN "public"."device"."rtsp_url" IS '摄像机视频流地址';
COMMENT ON COLUMN "public"."device"."video_url" IS '摄像机前端播放地址';
COMMENT ON COLUMN "public"."device"."switch_ip" IS '所属的远程开关ip';
COMMENT ON COLUMN "public"."device"."location_code" IS '监测位置下的逻辑编号';
COMMENT ON COLUMN "public"."device"."install_image" IS '传感器安装图片';
COMMENT ON TABLE "public"."device" IS '传感器实例表';

-- ----------------------------
-- Records of device
-- ----------------------------
INSERT INTO "public"."device" VALUES (62, '01', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:09:44', '2025-08-28 10:12:47', 0, '-', '', '', '', '04', '00', '');
INSERT INTO "public"."device" VALUES (63, '02', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:10:18', '2025-08-27 15:25:04', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (177, '06', 12, 'ACCY', '垂直主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:34:14', '2025-10-24 16:34:14', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (183, '12', 12, 'ACCY', '垂直主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:42:01', '2025-10-24 16:42:01', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (64, '03', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:10:35', '2025-08-27 15:25:21', 0, '-', '', '', '', '02', '02', '');
INSERT INTO "public"."device" VALUES (65, '04', 6, 'JMT', '测缝计', 1, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:10:51', '2025-08-27 15:25:36', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."device" VALUES (88, '01', 9, 'IPC', '网络摄像机', 5, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 16:01:17', '2025-08-27 15:55:55', 0, '-', 'rtsp://admin:Hd135246@10.184.10.39:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (84, '01', 9, 'IPC', '网络摄像机', 1, 9, 'FD_00', '基础地基基坑', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 15:40:33', '2025-08-28 10:23:01', 0, '-', 'rtsp://admin:Hd135246@10.184.10.14:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (85, '01', 9, 'IPC', '网络摄像机', 2, 9, 'FD_00', '基础地基基坑', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 15:44:54', '2025-08-27 15:31:09', 0, '-', 'rtsp://admin:Hd135246@10.184.10.17:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (87, '01', 9, 'IPC', '网络摄像机', 4, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 16:00:02', '2025-08-27 15:52:31', 0, '-', 'rtsp://admin:Hd135246@10.184.10.34:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (86, '01', 9, 'IPC', '网络摄像机', 3, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 15:57:32', '2025-08-27 15:47:52', 0, '-', 'rtsp://admin:Hd135246@10.184.10.29:554/h264/ch1/main/av_stream', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (187, '16', 12, 'ACCY', '垂直主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:44:39', '2025-10-24 16:44:39', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (67, '02', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:11:24', '2025-08-27 15:40:14', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (68, '03', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:11:47', '2025-08-27 15:40:22', 0, '-', '', '', '', '02', '02', '');
INSERT INTO "public"."device" VALUES (69, '04', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:12:02', '2025-08-27 15:40:30', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."device" VALUES (180, '09', 4, 'ACCX', '主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:40:29', '2025-10-24 16:40:29', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (91, '03', 4, 'ACCX', '主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 10:45:52', '2025-09-04 10:45:52', 0, '-', '', '', '', '01', '03', '');
INSERT INTO "public"."device" VALUES (191, '20', 12, 'ACCY', '垂直主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:58:58', '2025-10-24 16:58:58', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (89, '02', 4, 'ACCX', '主风向加速度计', 1, 16, 'TW_20', '塔架_20米', '50', '20', 'x', '', 0, '', '2025-09-04 10:41:22', '2025-09-04 10:41:22', 1, '-', '', '', '', '', '01', '');
INSERT INTO "public"."device" VALUES (35, '14', 4, 'ACCX', '主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:48', '2025-08-27 15:41:24', 0, '2', '', '', '', '02', '12', '');
INSERT INTO "public"."device" VALUES (39, '06', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:30:34', '2025-07-31 14:30:34', 0, '2', '', '', '', '04', '04', '');
INSERT INTO "public"."device" VALUES (52, '03', 4, 'ACCX', '主风向加速度计', 4, 1, 'BL', '叶片', '51', '80.12', '143.648', '45.856', 0, '', '2025-08-04 11:15:14', '2006-01-02 15:04:05', 1, '2', '', '', '', '', '02', '');
INSERT INTO "public"."device" VALUES (80, '01', 10, 'ULS', '超声波液位计', 5, 1, 'BL', '叶片', '', '', '', '', 0, '', '2025-08-07 10:15:52', '2025-08-07 10:15:52', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."device" VALUES (79, '01', 10, 'ULS', '超声波液位计', 4, 1, 'BL', '叶片', '', '', '', '', 0, '', '2025-08-07 10:15:30', '2025-08-07 10:15:30', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."device" VALUES (131, '01', 9, 'IPC', '网络摄像机', 12, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:25:01', '2025-09-18 17:25:01', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (128, '01', 9, 'IPC', '网络摄像机', 9, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:22:09', '2025-09-18 17:22:09', 0, '-', 'rtmp://172.16.90.70/live/33', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (1, '01', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:20:04', '2025-08-27 15:33:27', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (2, '02', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:20:31', '2025-08-27 15:33:47', 0, '2', '', '', '', '02', '01', '');
INSERT INTO "public"."device" VALUES (3, '03', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:20:50', '2025-08-27 15:34:00', 0, '2', '', '', '', '03', '02', '');
INSERT INTO "public"."device" VALUES (17, '01', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:25:01', '2025-08-27 15:53:14', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (138, '02', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:47:51', '2025-10-24 15:47:51', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (70, '01', 1, 'HLS', '静力水准仪', 1, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:12:53', '2025-10-09 17:14:15', 0, '-', '', '', '', '01', '00', 'http://59.110.219.98:81/device_image/accel_1.png');
INSERT INTO "public"."device" VALUES (134, '02', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:44:10', '2025-10-24 15:44:10', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (155, '02', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:02:04', '2025-10-24 16:02:04', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (135, '03', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:44:39', '2025-10-24 15:44:39', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (184, '13', 4, 'ACCX', '主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:43:33', '2025-10-24 16:43:33', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (188, '17', 4, 'ACCX', '主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:57:53', '2025-10-24 16:57:53', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (141, '05', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:49:34', '2025-10-24 15:49:34', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (145, '09', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:52:18', '2025-10-24 15:52:18', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (149, '01', 10, 'ULS', '超声波液位计', 11, 9, 'FD_00', '基础地基基坑', '50', '00', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:55:04', '2025-10-24 15:55:04', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (152, '03', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:58:06', '2025-10-24 15:58:06', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (37, '01', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:30:13', '2025-08-27 15:41:49', 0, '2', '', '', '', '02', '00', '');
INSERT INTO "public"."device" VALUES (44, '02', 4, 'ACCX', '主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:31:48', '2025-07-31 14:31:48', 0, '2', '', '', '', '02', '01', '');
INSERT INTO "public"."device" VALUES (158, '05', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:03:45', '2025-10-24 16:03:45', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (42, '06', 4, 'ACCX', '主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:31:19', '2025-08-27 15:50:47', 0, '2', '', '', '', '01', '05', '');
INSERT INTO "public"."device" VALUES (29, '18', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:32', '2025-08-27 15:39:00', 0, '2', '', '', '', '02', '17', '');
INSERT INTO "public"."device" VALUES (30, '19', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '123.648', '41.856', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:48', '2025-08-27 15:39:10', 0, '2', '', '', '', '01', '19', '');
INSERT INTO "public"."device" VALUES (81, '01', 3, 'GNSS', 'GNSS', 3, 17, 'NA_156', '机舱顶', 'x', 'x', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:16:33', '2025-12-10 09:49:31', 0, '-', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (53, '01', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:05:37', '2025-08-27 15:29:41', 0, '-', '', '', '', '02', '00', '');
INSERT INTO "public"."device" VALUES (54, '02', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:05:54', '2025-08-27 15:30:00', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (55, '03', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:06:17', '2025-08-27 15:30:17', 0, '-', '', '', '', '04', '02', '');
INSERT INTO "public"."device" VALUES (56, '04', 5, 'ATS', '锚索计', 2, 10, 'CT_06', '索力6米', '50', '20', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:07:00', '2025-08-27 15:30:28', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."device" VALUES (60, '04', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:08:27', '2025-08-27 15:46:43', 0, '-', '', '', '', '03', '03', '');
INSERT INTO "public"."device" VALUES (83, '04', 1, 'HLS', '静力水准仪', 2, 1, 'BL', '叶片', '22', '22', '22', '22', 0, '2025-08-13T16:00:00.000Z', '2025-08-07 11:34:17', '2025-08-07 11:34:17', 1, '-', '', '', '', '', '03', '');
INSERT INTO "public"."device" VALUES (74, '02', 1, 'HLS', '静力水准仪', 3, 8, 'TW_00', '塔架塔底平台', '50', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:43', '2025-08-27 15:47:24', 0, '-', '', '', '', '03', '01', '');
INSERT INTO "public"."device" VALUES (75, '03', 1, 'HLS', '静力水准仪', 3, 8, 'TW_00', '塔架塔底平台', '50', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:14:00', '2025-08-27 15:47:33', 0, '-', '', '', '', '02', '02', '');
INSERT INTO "public"."device" VALUES (49, '02', 2, 'INSX', '主风向倾角传感器', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:34:06', '2025-08-27 15:50:11', 0, '2', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (47, '01', 2, 'INSX', '主风向倾角传感器', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:33:24', '2025-08-27 15:38:13', 0, '2', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (120, '01', 13, 'INSY', '垂直主风向倾角传感器', 4, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'x', 'x', 0, '', '2025-09-12 09:57:11', '2025-09-12 09:57:11', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (118, '02', 13, 'INSY', '垂直主风向倾角传感器', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-09-11T16:00:00.000Z', '2025-09-12 09:50:02', '2025-09-12 15:56:17', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (119, '01', 13, 'INSY', '垂直主风向倾角传感器', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '', '2025-09-12 09:52:47', '2025-09-12 09:52:47', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (121, '02', 13, 'INSY', '垂直主风向倾角传感器', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '', '2025-09-12 09:58:16', '2025-09-12 09:58:16', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (161, '08', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:06:37', '2025-10-24 16:06:37', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (164, '11', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:08:08', '2025-10-24 16:08:08', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (167, '01', 8, 'WPR', '测风雷达', 12, 17, 'NA_156', '机舱顶', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:12:23', '2025-10-24 16:13:01', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (170, '01', 2, 'INSX', '主风向倾角传感器', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:17:49', '2025-10-24 16:17:49', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (173, '02', 12, 'ACCY', '垂直主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:30:36', '2025-10-24 16:30:36', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (31, '22', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:58', '2025-08-27 15:39:24', 0, '2', '', '', '', '04', '20', '');
INSERT INTO "public"."device" VALUES (32, '23', 4, 'ACCX', '主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:10', '2025-08-27 15:39:36', 0, '2', '', '', '', '03', '23', '');
INSERT INTO "public"."device" VALUES (66, '01', 6, 'JMT', '测缝计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:11:04', '2025-08-27 15:40:03', 0, '-', '', '', '', '04', '00', '');
INSERT INTO "public"."device" VALUES (9, '09', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:05', '2025-08-27 15:35:48', 0, '2', '', '', '', '01', '08', '');
INSERT INTO "public"."device" VALUES (10, '10', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:08', '2025-08-27 15:35:59', 0, '2', '', '', '', '02', '09', '');
INSERT INTO "public"."device" VALUES (11, '11', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:18', '2025-08-27 15:36:09', 0, '2', '', '', '', '03', '10', '');
INSERT INTO "public"."device" VALUES (175, '04', 12, 'ACCY', '垂直主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:32:31', '2025-10-24 16:32:31', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (181, '10', 12, 'ACCY', '垂直主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:40:53', '2025-10-24 16:40:53', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (185, '14', 12, 'ACCY', '垂直主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:43:54', '2025-10-24 16:43:54', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (57, '01', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:07:34', '2025-08-27 15:46:16', 0, '-', '', '', '', '02', '00', '');
INSERT INTO "public"."device" VALUES (58, '02', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:07:54', '2025-08-27 15:46:24', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (59, '03', 5, 'ATS', '锚索计', 3, 10, 'CT_06', '索力6米', '50', '20', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:08:09', '2025-08-27 15:46:33', 0, '-', '', '', '', '04', '02', '');
INSERT INTO "public"."device" VALUES (109, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '35', '', '', 0, '', '2025-09-10 16:01:23', '2025-09-10 16:01:23', 1, '-', '', '', '', '35', '02', '');
INSERT INTO "public"."device" VALUES (110, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '40', '', '', 0, '', '2025-09-10 16:03:25', '2025-09-10 16:03:25', 1, '-', '', '', '', '40', '03', '');
INSERT INTO "public"."device" VALUES (73, '01', 1, 'HLS', '静力水准仪', 3, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:26', '2025-08-27 15:47:13', 0, '-', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (46, '02', 2, 'INSX', '主风向倾角传感器', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:33:02', '2025-09-12 15:56:07', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (12, '12', 7, 'STM', '应变计', 3, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:23:21', '2025-08-27 15:36:21', 0, '2', '', '', '', '04', '11', '');
INSERT INTO "public"."device" VALUES (23, '09', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:26:34', '2025-08-27 15:55:13', 0, '2', '', '', '', '01', '08', '');
INSERT INTO "public"."device" VALUES (24, '10', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '', '156', '', '', 0, '', '2025-07-31 14:26:52', '2025-07-31 14:26:52', 0, '2', '', '', '', '02', '09', '');
INSERT INTO "public"."device" VALUES (25, '11', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '', '156', '', '', 0, '', '2025-07-31 14:27:00', '2025-07-31 14:27:00', 0, '2', '', '', '', '03', '10', '');
INSERT INTO "public"."device" VALUES (174, '03', 4, 'ACCX', '主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:32:04', '2025-10-24 16:32:04', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (178, '07', 4, 'ACCX', '主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:34:51', '2025-10-24 16:34:51', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (93, '07', 4, 'ACCX', '主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 10:48:32', '2025-09-04 10:48:32', 0, '-', '', '', '', '02', '07', '');
INSERT INTO "public"."device" VALUES (8, '08', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:22:50', '2025-08-27 15:35:32', 0, '2', '', '', '', '04', '07', '');
INSERT INTO "public"."device" VALUES (16, '05', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:24:30', '2025-08-27 15:54:35', 0, '2', '', '', '', '01', '04', '');
INSERT INTO "public"."device" VALUES (18, '06', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '', '156', '', '', 0, '', '2025-07-31 14:25:24', '2025-07-31 14:25:24', 0, '2', '', '', '', '02', '05', '');
INSERT INTO "public"."device" VALUES (19, '07', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '', '156', '', '', 0, '', '2025-07-31 14:25:37', '2025-07-31 14:25:37', 0, '2', '', '', '', '03', '06', '');
INSERT INTO "public"."device" VALUES (22, '08', 7, 'STM', '应变计', 5, 11, 'BL02_156', '叶片02', '', '156', '', '', 0, '', '2025-07-31 14:26:22', '2025-07-31 14:26:22', 0, '2', '', '', '', '04', '07', '');
INSERT INTO "public"."device" VALUES (5, '05', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:21:19', '2025-08-27 15:34:27', 0, '2', '', '', '', '01', '04', '');
INSERT INTO "public"."device" VALUES (111, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '50', '', '', 0, '', '2025-09-10 16:04:24', '2025-09-10 16:04:24', 1, '-', '', '', '', '50', '04', '');
INSERT INTO "public"."device" VALUES (6, '06', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:21:34', '2025-08-27 15:34:47', 0, '2', '', '', '', '02', '05', '');
INSERT INTO "public"."device" VALUES (7, '07', 7, 'STM', '应变计', 3, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:22:24', '2025-08-27 15:35:08', 0, '2', '', '', '', '03', '06', '');
INSERT INTO "public"."device" VALUES (78, '01', 10, 'ULS', '超声波液位计', 3, 9, 'FD_00', '基础地基基坑', '0', '0', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:15:15', '2025-08-27 15:48:10', 0, '-', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (112, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '60', '', '', 0, '', '2025-09-10 16:05:02', '2025-09-10 16:05:02', 1, '-', '', '', '', '60', '05', '');
INSERT INTO "public"."device" VALUES (113, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '70', '', '', 0, '', '2025-09-10 16:05:31', '2025-09-10 16:05:31', 1, '-', '', '', '', '70', '06', '');
INSERT INTO "public"."device" VALUES (114, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '80', '', '', 0, '', '2025-09-10 16:06:18', '2025-09-10 16:06:18', 1, '-', '', '', '', '80', '07', '');
INSERT INTO "public"."device" VALUES (115, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '90', '', '', 0, '', '2025-09-10 16:06:56', '2025-09-10 16:06:56', 1, '-', '', '', '', '90', '08', '');
INSERT INTO "public"."device" VALUES (116, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '100', '', '', 0, '', '2025-09-10 16:07:34', '2025-09-10 16:07:34', 1, '-', '', '', '', '100', '09', '');
INSERT INTO "public"."device" VALUES (117, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '', '110', '', '', 0, '', '2025-09-10 16:08:23', '2025-09-10 16:08:23', 1, '-', '', '', '', '110', '10', '');
INSERT INTO "public"."device" VALUES (61, '01', 8, 'WPR', '测风雷达', 3, 17, 'NA_156', '机舱顶', '10s', '156', 'x', 'x', 0, '', '2025-08-07 10:09:23', '2025-08-07 10:09:23', 0, '-', '', '', '', '01', '01', '');
INSERT INTO "public"."device" VALUES (129, '01', 9, 'IPC', '网络摄像机', 10, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:24:18', '2025-09-18 17:24:18', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (13, '02', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '', '156', '', '', 0, '', '2025-07-31 14:23:51', '2025-07-31 14:23:51', 0, '2', '', '', '', '02', '01', '');
INSERT INTO "public"."device" VALUES (132, '01', 9, 'IPC', '网络摄像机', 13, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:25:22', '2025-09-18 17:25:22', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (14, '03', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '', '156', '', '', 0, '', '2025-07-31 14:24:20', '2025-07-31 14:24:20', 0, '2', '', '', '', '03', '02', '');
INSERT INTO "public"."device" VALUES (15, '04', 7, 'STM', '应变计', 5, 6, 'BL01_156', '叶片01', '', '156', '', '', 0, '', '2025-07-31 14:24:26', '2025-07-31 14:24:26', 0, '2', '', '', '', '04', '03', '');
INSERT INTO "public"."device" VALUES (4, '04', 7, 'STM', '应变计', 3, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:21:05', '2025-08-27 15:34:09', 0, '2', '', '', '', '04', '03', '');
INSERT INTO "public"."device" VALUES (139, '03', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:48:19', '2025-10-24 15:48:19', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (71, '02', 1, 'HLS', '静力水准仪', 1, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:02', '2025-10-09 09:20:46', 0, '-', '', '', '', '03', '01', 'http://172.16.90.70:8081/device_image/login_bg1.jpeg');
INSERT INTO "public"."device" VALUES (72, '03', 1, 'HLS', '静力水准仪', 1, 8, 'TW_00', '塔架塔底平台', '50', '0', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-08-07 10:13:14', '2025-10-09 09:21:02', 0, '-', '', '', '', '02', '02', 'http://172.16.90.70:8081/device_image/login_bg2.jpeg');
INSERT INTO "public"."device" VALUES (136, '04', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:45:05', '2025-10-24 15:45:05', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (142, '06', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:50:05', '2025-10-24 15:50:05', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (143, '07', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:50:34', '2025-10-24 15:50:34', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (146, '10', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:52:53', '2025-10-24 15:52:53', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (102, '17', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:18:17', '2025-09-04 11:18:17', 0, '-', '', '', '', '02', '16', '');
INSERT INTO "public"."device" VALUES (150, '01', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:55:44', '2025-10-24 15:56:23', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (153, '04', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:58:28', '2025-10-24 15:58:28', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (159, '06', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:05:42', '2025-10-24 16:05:42', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (162, '09', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:07:14', '2025-10-24 16:07:14', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (165, '12', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:08:31', '2025-10-24 16:08:31', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (168, '02', 2, 'INSX', '主风向倾角传感器', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:15:24', '2025-10-24 16:16:49', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (171, '01', 13, 'INSY', '垂直主风向倾角传感器', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'x', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:18:13', '2025-10-24 16:18:13', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (156, '03', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:02:30', '2025-10-24 16:02:30', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (26, '12', 7, 'STM', '应变计', 5, 12, 'BL03_156', '叶片03', '', '156', '', '', 0, '', '2025-07-31 14:27:13', '2025-07-31 14:27:13', 0, '2', '', '', '', '04', '11', '');
INSERT INTO "public"."device" VALUES (76, '01', 10, 'ULS', '超声波液位计', 1, 2, 'TW', '塔架', '', '0', '', '', 0, '', '2025-08-07 10:14:41', '2025-08-07 10:14:41', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."device" VALUES (77, '01', 10, 'ULS', '超声波液位计', 2, 1, 'BL', '叶片', '', '', '', '', 0, '', '2025-08-07 10:14:51', '2025-08-07 10:14:51', 1, '-', '', '', '', '', '00', '');
INSERT INTO "public"."device" VALUES (94, '02', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:02:00', '2025-09-04 11:02:00', 0, '-', '', '', '', '02', '01', 'http://127.0.0.1:81/device_image\FY\04\TW_20\ACCX\02/accel_1.png');
INSERT INTO "public"."device" VALUES (130, '01', 9, 'IPC', '网络摄像机', 11, 9, 'FD_00', '基础地基基坑', '2h', '00', 'xx', 'xx', 0, '2025-09-17T16:00:00.000Z', '2025-09-18 17:24:39', '2025-09-18 17:24:39', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (98, '10', 4, 'ACCX', '主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-08T16:00:00.000Z', '2025-09-04 11:10:51', '2025-10-09 16:50:34', 0, '-', '', '', '', '02', '09', 'http://172.16.90.70:8081/device_image/accel_1.png');
INSERT INTO "public"."device" VALUES (172, '01', 4, 'ACCX', '主风向加速度计', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:29:05', '2025-10-24 16:29:05', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (176, '05', 4, 'ACCX', '主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:33:37', '2025-10-24 16:33:37', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (182, '11', 4, 'ACCX', '主风向加速度计', 12, 14, 'TW_80', '塔架_80米', '50', '80', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:41:28', '2025-10-24 16:41:28', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (186, '15', 4, 'ACCX', '主风向加速度计', 12, 15, 'TW_50', '塔架_50米', '50', '50', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:44:20', '2025-10-24 16:44:20', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (144, '08', 7, 'STM', '应变计', 10, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:51:20', '2025-10-24 15:51:20', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (147, '11', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:53:29', '2025-10-24 15:53:29', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (148, '12', 7, 'STM', '应变计', 10, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:53:49', '2025-10-24 15:53:49', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (190, '19', 4, 'ACCX', '主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:58:39', '2025-10-24 16:58:39', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (28, '27', 4, 'ACCX', '主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:16', '2025-08-27 15:37:23', 0, '2', '', '', '', '02', '27', '');
INSERT INTO "public"."device" VALUES (151, '02', 5, 'ATS', '锚索计', 11, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:56:14', '2025-10-24 15:56:14', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (160, '07', 7, 'STM', '应变计', 13, 11, 'BL02_156', '叶片02', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:06:10', '2025-10-24 16:06:10', 0, '-', '', '', '', '03', NULL, '');
INSERT INTO "public"."device" VALUES (163, '10', 7, 'STM', '应变计', 13, 12, 'BL03_156', '叶片03', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:07:46', '2025-10-24 16:07:46', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (166, '01', 3, 'GNSS', 'GNSS', 12, 17, 'NA_156', '机舱顶', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:11:31', '2025-10-24 16:12:56', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (27, '26', 4, 'ACCX', '主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '22', '156', '22', '22', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:28:07', '2025-08-27 15:37:02', 0, '2', '', '', '', '01', '25', '');
INSERT INTO "public"."device" VALUES (169, '02', 13, 'INSY', '垂直主风向倾角传感器', 12, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:15:55', '2025-10-24 16:16:54', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (179, '08', 12, 'ACCY', '垂直主风向加速度计', 12, 7, 'TW_114', '塔架钢混转接', '50', '114', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:35:12', '2025-10-24 16:35:12', 0, '-', '', '', '', '02', NULL, '');
INSERT INTO "public"."device" VALUES (189, '18', 12, 'ACCY', '垂直主风向加速度计', 12, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:58:16', '2025-10-24 16:58:16', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (99, '11', 4, 'ACCX', '主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', '', '', 0, '', '2025-09-04 11:11:40', '2025-09-04 11:11:40', 0, '-', '', '', '', '01', '11', '');
INSERT INTO "public"."device" VALUES (36, '15', 4, 'ACCX', '主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:30:03', '2025-08-27 15:41:32', 0, '2', '', '', '', '01', '15', '');
INSERT INTO "public"."device" VALUES (137, '01', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:47:21', '2025-10-24 15:47:21', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (108, '28', 12, 'ACCY', '垂直主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', '', '', 0, '', '2025-09-04 11:28:44', '2025-09-04 11:28:44', 0, '-', '', '', '', '02', '26', '');
INSERT INTO "public"."device" VALUES (107, '25', 12, 'ACCY', '垂直主风向加速度计', 3, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', '', '', 0, '', '2025-09-04 11:27:43', '2025-09-04 11:27:43', 0, '-', '', '', '', '01', '24', '');
INSERT INTO "public"."device" VALUES (92, '08', 12, 'ACCY', '垂直主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 10:47:57', '2025-09-04 10:47:57', 0, '-', '', '', '', '02', '06', '');
INSERT INTO "public"."device" VALUES (103, '20', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:19:43', '2025-09-04 11:19:43', 0, '-', '', '', '', '01', '18', '');
INSERT INTO "public"."device" VALUES (43, '01', 12, 'ACCY', '垂直主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:31:35', '2025-07-31 14:31:35', 0, '2', '', '', '', '02', '00', '');
INSERT INTO "public"."device" VALUES (104, '21', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:20:42', '2025-09-04 11:20:42', 0, '-', '', '', '', '04', '21', '');
INSERT INTO "public"."device" VALUES (105, '24', 12, 'ACCY', '垂直主风向加速度计', 3, 7, 'TW_114', '塔架钢混转接', '50', '114', '', '', 0, '', '2025-09-04 11:21:16', '2025-09-04 11:21:16', 0, '-', '', '', '', '03', '22', '');
INSERT INTO "public"."device" VALUES (33, '09', 12, 'ACCY', '垂直主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:20', '2025-08-27 15:41:03', 0, '2', '', '', '', '02', '08', '');
INSERT INTO "public"."device" VALUES (90, '04', 12, 'ACCY', '垂直主风向加速度计', 4, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 10:45:07', '2025-09-04 10:45:07', 0, '-', '', '', '', '01', '02', '');
INSERT INTO "public"."device" VALUES (34, '12', 12, 'ACCY', '垂直主风向加速度计', 3, 14, 'TW_80', '塔架_80米', '50', '80', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:29:33', '2025-08-27 15:41:13', 0, '2', '', '', '', '01', '10', '');
INSERT INTO "public"."device" VALUES (100, '13', 12, 'ACCY', '垂直主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', '', '', 0, '', '2025-09-04 11:14:13', '2025-09-04 11:14:13', 0, '-', '', '', '', '02', '13', '');
INSERT INTO "public"."device" VALUES (101, '16', 12, 'ACCY', '垂直主风向加速度计', 3, 15, 'TW_50', '塔架_50米', '50', '50', '', '', 0, '', '2025-09-04 11:15:17', '2025-09-04 11:15:17', 0, '-', '', '', '', '01', '14', '');
INSERT INTO "public"."device" VALUES (140, '04', 7, 'STM', '应变计', 10, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:48:42', '2025-10-24 15:48:42', 0, '-', '', '', '', '04', NULL, '');
INSERT INTO "public"."device" VALUES (41, '05', 12, 'ACCY', '垂直主风向加速度计', 4, 7, 'TW_114', '塔架钢混转接', '50', '114', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:31:01', '2025-08-27 15:50:39', 0, '2', '', '', '', '01', '04', '');
INSERT INTO "public"."device" VALUES (38, '04', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:30:26', '2025-07-31 14:30:26', 0, '2', '', '', '', '01', '02', '');
INSERT INTO "public"."device" VALUES (96, '05', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:04:41', '2025-09-04 11:04:41', 0, '-', '', '', '', '04', '05', '');
INSERT INTO "public"."device" VALUES (40, '08', 12, 'ACCY', '垂直主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-07-31 14:30:45', '2025-07-31 14:30:45', 0, '2', '', '', '', '03', '06', '');
INSERT INTO "public"."device" VALUES (95, '03', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:03:23', '2025-09-04 11:03:23', 0, '-', '', '', '', '01', '03', '');
INSERT INTO "public"."device" VALUES (97, '07', 4, 'ACCX', '主风向加速度计', 3, 16, 'TW_20', '塔架_20米', '50', '20', '', '', 0, '', '2025-09-04 11:05:56', '2025-09-04 11:05:56', 0, '-', '', '', '', '03', '07', '');
INSERT INTO "public"."device" VALUES (133, '01', 5, 'ATS', '锚索计', 9, 16, 'TW_20', '塔架_20米', '50', '20', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 15:43:17', '2025-12-01 18:45:04', 0, '1234', '', '', '', '01', NULL, 'http://59.110.219.98:81/device_image/制作电脑桌面背景.png');
INSERT INTO "public"."device" VALUES (48, '01', 2, 'INSX', '主风向倾角传感器', 4, 13, 'TW_156', '塔架_钢段顶平台', '50', '156', 'x', 'x', 0, '2025-08-26T16:00:00.000Z', '2025-07-31 14:33:41', '2025-08-27 15:50:00', 0, '2', '', '', '', '01', '00', '');
INSERT INTO "public"."device" VALUES (154, '01', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:01:40', '2025-10-24 16:01:40', 0, '-', '', '', '', '01', NULL, '');
INSERT INTO "public"."device" VALUES (157, '04', 7, 'STM', '应变计', 13, 6, 'BL01_156', '叶片01', '50', '156', 'xx', 'xx', 0, '2025-10-23T16:00:00.000Z', '2025-10-24 16:02:55', '2025-10-24 16:02:55', 0, '-', '', '', '', '04', NULL, '');

-- ----------------------------
-- Table structure for device_meta
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_meta";
CREATE TABLE "public"."device_meta" (
  "meta_id" int4 NOT NULL DEFAULT nextval('device_meta_id_seq'::regclass),
  "device_type_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "column_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "display_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "unit" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "ord" int4 NOT NULL DEFAULT 1,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_type_id" int4 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."device_meta"."meta_id" IS '传感器元数据id';
COMMENT ON COLUMN "public"."device_meta"."device_type_code" IS '传感器类型编号';
COMMENT ON COLUMN "public"."device_meta"."column_name" IS '数据库中的列名';
COMMENT ON COLUMN "public"."device_meta"."display_name" IS '前端显示名';
COMMENT ON COLUMN "public"."device_meta"."unit" IS '单位';
COMMENT ON COLUMN "public"."device_meta"."ord" IS '字段显示顺序';
COMMENT ON COLUMN "public"."device_meta"."remark" IS '备注';
COMMENT ON COLUMN "public"."device_meta"."device_type_id" IS '传感器类型id';
COMMENT ON TABLE "public"."device_meta" IS '传感器元数据映射表';

-- ----------------------------
-- Records of device_meta
-- ----------------------------
INSERT INTO "public"."device_meta" VALUES (6, 'INSX', 'x', '倾角', '°', 1, '', 2, true);
INSERT INTO "public"."device_meta" VALUES (7, 'INSY', 'y', '倾角', '°', 2, '', 13, true);
INSERT INTO "public"."device_meta" VALUES (2, 'ACCY', 'accel', '加速度', 'mg', 2, '', 12, true);
INSERT INTO "public"."device_meta" VALUES (1, 'ACCX', 'accel', '加速度', 'mg', 1, '', 4, true);
INSERT INTO "public"."device_meta" VALUES (5, 'STM', 'strain', '动静态应变', 'με', 1, '', 7, true);
INSERT INTO "public"."device_meta" VALUES (8, 'ATS', 'tension', '张拉力', 'kN', 1, '', 5, true);
INSERT INTO "public"."device_meta" VALUES (25, 'WPR', 'ti', '湍流强度', ' ', 5, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (28, 'WPR', 'v_sheer', '垂直风切变', ' ', 8, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (29, 'WPR', 'h_sheer', '水平风切变', ' ', 9, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (17, 'GNSS', 'vertical_offset', '竖向偏移量', 'mm', 6, '-', 3, true);
INSERT INTO "public"."device_meta" VALUES (18, 'GNSS', 'horizontal_offset', '横向偏移量', 'mm', 7, '-', 3, true);
INSERT INTO "public"."device_meta" VALUES (19, 'GNSS', 'ordinate_offset', '纵向偏移量', 'mm', 8, '-', 3, true);
INSERT INTO "public"."device_meta" VALUES (23, 'WPR', 'veer', '垂直风向变化率', 'deg/m', 3, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (31, 'WPR', 'direction_high', '上平面处风向', '°', 11, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (9, 'JMT', 'joint', '缝隙', 'mm', 1, '', 6, true);
INSERT INTO "public"."device_meta" VALUES (10, 'HLS', 'settlement', '沉降', 'mm', 1, '', 1, true);
INSERT INTO "public"."device_meta" VALUES (20, 'IPC', 'image_url', '图像名称', 'px', 1, '', 9, true);
INSERT INTO "public"."device_meta" VALUES (12, 'GNSS', 'longitude', '经度', '°', 1, '', 3, true);
INSERT INTO "public"."device_meta" VALUES (13, 'GNSS', 'latitude', '纬度', '°', 2, '', 3, true);
INSERT INTO "public"."device_meta" VALUES (11, 'ULS', 'height', '液位高度', 'mm', 1, '', 10, true);
INSERT INTO "public"."device_meta" VALUES (21, 'WPR', 'd', '测量距离', 'm', 1, '', 8, true);
INSERT INTO "public"."device_meta" VALUES (22, 'WPR', 'rws', '视向风速', 'm/s', 2, '', 8, true);
INSERT INTO "public"."device_meta" VALUES (24, 'WPR', 'raws', '轴向投影风速', 'm/s', 4, '', 8, true);
INSERT INTO "public"."device_meta" VALUES (26, 'WPR', 'hw_shub', '轮毂高度处风速', 'm/s', 6, '', 8, true);
INSERT INTO "public"."device_meta" VALUES (30, 'WPR', 'hw_shigh', '上平面处风速', 'm/s', 10, '', 8, true);
INSERT INTO "public"."device_meta" VALUES (32, 'WPR', 'hw_slow', '下平面处风速', 'm/s', 12, '', 8, true);
INSERT INTO "public"."device_meta" VALUES (33, 'WPR', 'direction_low', '下平面处风向', '°', 13, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (27, 'WPR', 'direction_hub', '轮毂高度处风向', '°', 7, '-', 8, true);
INSERT INTO "public"."device_meta" VALUES (14, 'GNSS', 'vertical', '竖向坐标', 'm', 3, '-', 3, true);
INSERT INTO "public"."device_meta" VALUES (15, 'GNSS', 'horizontal', '横向坐标', 'm', 4, '-', 3, true);
INSERT INTO "public"."device_meta" VALUES (16, 'GNSS', 'ordinate', '纵向坐标', 'm', 4, '-', 3, true);

-- ----------------------------
-- Table structure for device_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."device_type";
CREATE TABLE "public"."device_type" (
  "device_type_id" int4 NOT NULL DEFAULT nextval('device_type_id_seq'::regclass),
  "device_type_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "device_type_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "sample_unit" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying
)
;
COMMENT ON COLUMN "public"."device_type"."device_type_id" IS '记录ID';
COMMENT ON COLUMN "public"."device_type"."device_type_code" IS '传感器类型编号，如: HLS';
COMMENT ON COLUMN "public"."device_type"."device_type_name" IS '传感器类型名称，如: 静力水准仪';
COMMENT ON COLUMN "public"."device_type"."sample_unit" IS '采样单位';
COMMENT ON COLUMN "public"."device_type"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."device_type"."updated_at" IS '记录修改时间';
COMMENT ON COLUMN "public"."device_type"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."device_type"."remark" IS '备注';
COMMENT ON TABLE "public"."device_type" IS '传感器类型表';

-- ----------------------------
-- Records of device_type
-- ----------------------------
INSERT INTO "public"."device_type" VALUES (2, 'INSX', '主风向倾角传感器', '°', '2025-07-31 14:15:49', '2025-09-17 15:10:05', 0, '与当地水平面夹角（绝对值）');
INSERT INTO "public"."device_type" VALUES (13, 'INSY', '垂直主风向倾角传感器', '°', '2025-09-12 09:46:01', '2025-09-17 15:10:22', 0, '与当地水平面夹角（绝对值）');
INSERT INTO "public"."device_type" VALUES (6, 'JMT', '测缝计', 'mm', '2025-07-31 14:17:04', '2025-09-17 15:10:43', 0, '混凝土裂缝张开宽度');
INSERT INTO "public"."device_type" VALUES (1, 'HLS', '静力水准仪', 'mm', '2025-07-31 13:43:32', '2025-09-17 15:09:29', 0, '沉降量（绝对值）');
INSERT INTO "public"."device_type" VALUES (10, 'ULS', '超声波液位计', 'mm', '2025-07-31 14:17:51', '2025-09-17 15:16:35', 0, '基坑水位深度');
INSERT INTO "public"."device_type" VALUES (11, 'HLS2222', '超声', '', '2025-08-04 10:35:19', '2006-01-02 15:04:05', 1, '-');
INSERT INTO "public"."device_type" VALUES (5, 'ATS', '锚索计', 'kN', '2025-07-31 14:16:52', '2025-09-17 15:09:03', 0, '索力张拉力');
INSERT INTO "public"."device_type" VALUES (7, 'STM', '应变计', 'με', '2025-07-31 14:17:15', '2025-09-17 15:13:19', 0, '叶根测点应变值（绝对值）');
INSERT INTO "public"."device_type" VALUES (12, 'ACCY', '垂直主风向加速度计', 'mg', '2025-09-04 10:31:24', '2025-09-17 15:08:44', 0, '采样50Hz-200Hz');
INSERT INTO "public"."device_type" VALUES (3, 'GNSS', 'GNSS', '', '2025-07-31 14:16:04', '2025-07-31 14:16:04', 0, '-');
INSERT INTO "public"."device_type" VALUES (9, 'IPC', '网络摄像机', '', '2025-07-31 14:17:38', '2025-07-31 14:17:38', 0, '-');
INSERT INTO "public"."device_type" VALUES (8, 'WPR', '测风雷达', 'm/s', '2025-07-31 14:17:26', '2025-08-19 14:05:50', 0, '-');
INSERT INTO "public"."device_type" VALUES (4, 'ACCX', '主风向加速度计', 'mg', '2025-07-31 14:16:36', '2025-11-24 15:37:21', 0, '采样50Hz-200Hz');

-- ----------------------------
-- Table structure for farms
-- ----------------------------
DROP TABLE IF EXISTS "public"."farms";
CREATE TABLE "public"."farms" (
  "farm_id" int4 NOT NULL DEFAULT nextval('farm_id_seq'::regclass),
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "province" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "location" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "latitude" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "longitude" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying
)
;
COMMENT ON COLUMN "public"."farms"."farm_id" IS '记录ID';
COMMENT ON COLUMN "public"."farms"."farm_code" IS '风场编号，例: FY';
COMMENT ON COLUMN "public"."farms"."farm_name" IS '风场名称，例: 扶余：市级单位缩写';
COMMENT ON COLUMN "public"."farms"."province" IS '省份缩写，例: JL(吉林)';
COMMENT ON COLUMN "public"."farms"."location" IS '地理位置';
COMMENT ON COLUMN "public"."farms"."latitude" IS '纬度';
COMMENT ON COLUMN "public"."farms"."longitude" IS '经度';
COMMENT ON COLUMN "public"."farms"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."farms"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."farms"."remark" IS '备注';
COMMENT ON TABLE "public"."farms" IS '风场表';

-- ----------------------------
-- Records of farms
-- ----------------------------
INSERT INTO "public"."farms" VALUES (1, 'FY', '华电万兴风电场', 'JL', '松原市扶余市', '45.227703', '125.482944', '2025-07-31 14:05:12', 0, '-');
INSERT INTO "public"."farms" VALUES (2, 'YS', '华电龙岗风电场', 'JL', '长春市榆树市', '44.803879', '126.439690', '2025-07-31 14:06:12', 0, '-');

-- ----------------------------
-- Table structure for model_device_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."model_device_map";
CREATE TABLE "public"."model_device_map" (
  "map_id" int4 NOT NULL DEFAULT nextval('map_id_seq'::regclass),
  "model_id" int4 NOT NULL DEFAULT 1,
  "device_type_id" int4 NOT NULL DEFAULT 1,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."model_device_map"."map_id" IS '映射表id';
COMMENT ON COLUMN "public"."model_device_map"."model_id" IS '模型id';
COMMENT ON COLUMN "public"."model_device_map"."device_type_id" IS '设备类型id';
COMMENT ON COLUMN "public"."model_device_map"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."model_device_map"."remark" IS '备注';
COMMENT ON COLUMN "public"."model_device_map"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON TABLE "public"."model_device_map" IS '模型设备映射表';

-- ----------------------------
-- Records of model_device_map
-- ----------------------------
INSERT INTO "public"."model_device_map" VALUES (3, 2, 2, '2025-10-16 13:47:44', '', 0);
INSERT INTO "public"."model_device_map" VALUES (4, 2, 13, '2025-10-16 13:48:33', '', 0);
INSERT INTO "public"."model_device_map" VALUES (5, 3, 2, '2025-10-16 13:48:44', '', 0);
INSERT INTO "public"."model_device_map" VALUES (6, 3, 13, '2025-10-16 13:48:50', '', 0);
INSERT INTO "public"."model_device_map" VALUES (8, 5, 1, '2025-10-16 13:49:26', '', 0);
INSERT INTO "public"."model_device_map" VALUES (9, 6, 10, '2025-10-16 13:49:53', '', 0);
INSERT INTO "public"."model_device_map" VALUES (14, 94, 5, '2025-10-17 15:32:26', '', 1);
INSERT INTO "public"."model_device_map" VALUES (15, 94, 13, '2025-10-17 15:32:26', '', 1);
INSERT INTO "public"."model_device_map" VALUES (16, 172, 12, '2025-10-17 17:20:55', '', 1);
INSERT INTO "public"."model_device_map" VALUES (17, 172, 4, '2025-10-17 17:20:55', '', 1);
INSERT INTO "public"."model_device_map" VALUES (1, 1, 2, '2025-10-16 13:46:12', '', 0);
INSERT INTO "public"."model_device_map" VALUES (2, 1, 13, '2025-10-16 13:46:51', '', 0);
INSERT INTO "public"."model_device_map" VALUES (18, 4, 6, '2025-11-28 11:34:47', '', 0);

-- ----------------------------
-- Table structure for model_records
-- ----------------------------
DROP TABLE IF EXISTS "public"."model_records";
CREATE TABLE "public"."model_records" (
  "record_id" int4 NOT NULL DEFAULT nextval('models_id_seq'::regclass),
  "run_id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "model_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "model_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "display_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "image_url" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "gif_url" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "txt_url" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "started_at" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "completed_at" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "duration_sec" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "data_start_at" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "data_end_at" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."model_records"."record_id" IS '模型id';
COMMENT ON COLUMN "public"."model_records"."run_id" IS '算法任务的唯一值id';
COMMENT ON COLUMN "public"."model_records"."model_code" IS '模型对应的编号';
COMMENT ON COLUMN "public"."model_records"."model_name" IS '模型名称';
COMMENT ON COLUMN "public"."model_records"."display_name" IS '前端显示的模型名称';
COMMENT ON COLUMN "public"."model_records"."farm_code" IS '风场名称';
COMMENT ON COLUMN "public"."model_records"."tower_code" IS '风机名称';
COMMENT ON COLUMN "public"."model_records"."image_url" IS '结果静图url';
COMMENT ON COLUMN "public"."model_records"."gif_url" IS '结果动图url';
COMMENT ON COLUMN "public"."model_records"."txt_url" IS '输出日志url';
COMMENT ON COLUMN "public"."model_records"."started_at" IS '开始时间';
COMMENT ON COLUMN "public"."model_records"."completed_at" IS '完成时间';
COMMENT ON COLUMN "public"."model_records"."duration_sec" IS '执行时间';
COMMENT ON COLUMN "public"."model_records"."remark" IS '备注';
COMMENT ON COLUMN "public"."model_records"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."model_records"."data_start_at" IS '数据开始时间';
COMMENT ON COLUMN "public"."model_records"."data_end_at" IS '数据结束时间';
COMMENT ON TABLE "public"."model_records" IS '算法模型解析记录表';

-- ----------------------------
-- Records of model_records
-- ----------------------------
INSERT INTO "public"."model_records" VALUES (192, '20251017173643690', '1', 'tilt_relative', '钢塔相对倾角', 'FY', '04', 'http://59.110.219.98:9081/11,0ae91f8df473.png', 'http://59.110.219.98:9081/8,0aea4af7a403.gif', 'http://59.110.219.98:9081/9,0ae889ca64ea.txt', '2025-10-17 17:36:46', '2025-10-17 17:37:06', '19.769', '', 0, '', '');
INSERT INTO "public"."model_records" VALUES (193, '20251017173643690', '1', 'tilt_top', '塔顶倾角', 'FY', '04', 'http://59.110.219.98:9081/11,0aec043978ba.png', 'http://59.110.219.98:9081/12,0aed556eb6f5.gif', 'http://59.110.219.98:9081/14,0aeb270adeaa.txt', '2025-10-17 17:37:06', '2025-10-17 17:37:24', '17.511', '', 0, '', '');

-- ----------------------------
-- Table structure for models
-- ----------------------------
DROP TABLE IF EXISTS "public"."models";
CREATE TABLE "public"."models" (
  "model_id" int4 NOT NULL DEFAULT nextval('models_id_seq'::regclass),
  "model_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "model_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "display_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."models"."model_id" IS '模型id';
COMMENT ON COLUMN "public"."models"."model_code" IS '模型对应的编号';
COMMENT ON COLUMN "public"."models"."model_name" IS '模型名称';
COMMENT ON COLUMN "public"."models"."display_name" IS '前端显示的模型名称';
COMMENT ON COLUMN "public"."models"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."models"."remark" IS '备注';
COMMENT ON COLUMN "public"."models"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON TABLE "public"."models" IS '算法模型表';

-- ----------------------------
-- Records of models
-- ----------------------------
INSERT INTO "public"."models" VALUES (5, '5', 'static_level', '静力水准仪', '2025-10-14 10:25:21', '-', 0);
INSERT INTO "public"."models" VALUES (2, '1', 'tilt_top', '塔顶倾角', '2025-10-14 10:24:49', '-', 0);
INSERT INTO "public"."models" VALUES (3, '1', 'tilt_joint', '钢混连接处倾角', '2025-10-14 10:25:01', '-', 0);
INSERT INTO "public"."models" VALUES (94, 'Acc', 'acc', '加速度', '2025-10-17 11:04:13', '', 1);
INSERT INTO "public"."models" VALUES (172, '', 'sss', 'sss', '2025-10-17 17:02:29', 'sss...', 1);
INSERT INTO "public"."models" VALUES (6, '6', 'ultrasonic_level', '超声波液位计', '2025-10-14 10:25:29', '-', 0);
INSERT INTO "public"."models" VALUES (7, '1', 'tilt_relative', '钢塔相对倾角', '2025-11-27 08:51:20', '', 1);
INSERT INTO "public"."models" VALUES (1, '1', 'tilt_relative', '钢塔相对倾角', '2025-10-14 10:24:36', '-', 0);
INSERT INTO "public"."models" VALUES (4, '4', 'joint_meter', '测缝计', '2025-10-14 10:25:10', '-', 0);

-- ----------------------------
-- Table structure for structure_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."structure_type";
CREATE TABLE "public"."structure_type" (
  "structure_id" int4 NOT NULL DEFAULT nextval('structure_id_seq'::regclass),
  "structure_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "structure_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying
)
;
COMMENT ON COLUMN "public"."structure_type"."structure_id" IS '结构类型id';
COMMENT ON COLUMN "public"."structure_type"."structure_code" IS '结构编号，例如: BL';
COMMENT ON COLUMN "public"."structure_type"."structure_name" IS '结构名称，例如: 叶片';
COMMENT ON COLUMN "public"."structure_type"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."structure_type"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."structure_type"."remark" IS '备注';
COMMENT ON TABLE "public"."structure_type" IS '风机结构类型表';

-- ----------------------------
-- Records of structure_type
-- ----------------------------
INSERT INTO "public"."structure_type" VALUES (7, 'TW_114', '塔架钢混转接', '2025-08-27 15:08:53', 0, '1,2.9,0');
INSERT INTO "public"."structure_type" VALUES (8, 'TW_00', '塔架塔底平台', '2025-08-27 15:09:39', 0, '-1,1.1,0');
INSERT INTO "public"."structure_type" VALUES (9, 'FD_00', '基础地基基坑', '2025-08-27 15:10:46', 0, '1,1.1,0');
INSERT INTO "public"."structure_type" VALUES (10, 'CT_06', '索力6米', '2025-08-27 15:11:25', 0, '-1,1.6,0');
INSERT INTO "public"."structure_type" VALUES (14, 'TW_80', '塔架_80米', '2025-08-27 15:20:30', 0, '-1,2.5,0');
INSERT INTO "public"."structure_type" VALUES (15, 'TW_50', '塔架_50米', '2025-08-27 15:20:55', 0, '1,2.2,0');
INSERT INTO "public"."structure_type" VALUES (16, 'TW_20', '塔架_20米', '2025-08-27 15:21:21', 0, '1,1.6,0');
INSERT INTO "public"."structure_type" VALUES (13, 'TW_156', '塔架_钢段顶平台', '2025-08-27 15:19:16', 0, '1,3.7,0');
INSERT INTO "public"."structure_type" VALUES (12, 'BL03_156', '叶片03', '2025-08-27 15:18:28', 0, '-1,3.7,0');
INSERT INTO "public"."structure_type" VALUES (11, 'BL02_156', '叶片02', '2025-08-27 15:17:53', 0, '-1,4.2,0');
INSERT INTO "public"."structure_type" VALUES (17, 'NA_156', '机舱顶', '2025-09-12 15:38:08', 0, '0,4.2,0');
INSERT INTO "public"."structure_type" VALUES (2, 'TW', '塔架', '2025-07-31 14:13:16', 1, '-');
INSERT INTO "public"."structure_type" VALUES (3, 'FD', '基础', '2025-07-31 14:14:32', 1, '-');
INSERT INTO "public"."structure_type" VALUES (4, 'CT', '索力', '2025-07-31 14:14:45', 1, '-');
INSERT INTO "public"."structure_type" VALUES (1, 'BL', '叶片', '2025-07-31 14:13:04', 1, '-');
INSERT INTO "public"."structure_type" VALUES (6, 'BL01_156', '叶片01', '2025-08-27 11:14:41', 0, '1,4.2,0');

-- ----------------------------
-- Table structure for switch_asset
-- ----------------------------
DROP TABLE IF EXISTS "public"."switch_asset";
CREATE TABLE "public"."switch_asset" (
  "asset_id" int4 NOT NULL DEFAULT nextval('switch_asset_id_seq'::regclass),
  "asset_code" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "asset_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "asset_type" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "switch_id" int4 NOT NULL DEFAULT 0,
  "asset_ip" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."switch_asset"."asset_id" IS '远程开关控制资产id';
COMMENT ON COLUMN "public"."switch_asset"."asset_code" IS '远程开关控制资产编号';
COMMENT ON COLUMN "public"."switch_asset"."asset_name" IS '远程开关控制资产名称';
COMMENT ON COLUMN "public"."switch_asset"."asset_type" IS '远程开关控制资产类型';
COMMENT ON COLUMN "public"."switch_asset"."switch_id" IS '远程开关id';
COMMENT ON COLUMN "public"."switch_asset"."asset_ip" IS '资产ip';
COMMENT ON COLUMN "public"."switch_asset"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."switch_asset"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."switch_asset"."remark" IS '备注';
COMMENT ON TABLE "public"."switch_asset" IS '远程开关资产表';

-- ----------------------------
-- Records of switch_asset
-- ----------------------------
INSERT INTO "public"."switch_asset" VALUES (1, '1', 'QA-J4Q4W8', '1', 1, '10.184.10.11', '2025-08-11 15:17:02', 0, '1');
INSERT INTO "public"."switch_asset" VALUES (2, '2', '串口模块', '1', 1, '10.184.10.13', '2025-08-11 15:17:20', 0, '');
INSERT INTO "public"."switch_asset" VALUES (3, '3', '网络摄像机', '1', 1, '10.184.10.14', '2025-08-11 15:19:49', 0, '');

-- ----------------------------
-- Table structure for switch_remote
-- ----------------------------
DROP TABLE IF EXISTS "public"."switch_remote";
CREATE TABLE "public"."switch_remote" (
  "switch_id" int4 NOT NULL DEFAULT nextval('remote_switch_id_seq'::regclass),
  "switch_code" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "switch_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "tower_id" int4 NOT NULL DEFAULT 0,
  "switch_ip" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "port" int4 NOT NULL DEFAULT 1031,
  "protocol" varchar(16) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'tcp'::character varying,
  "cmd_restart" varchar(16) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '!RST\r'::character varying,
  "enabled" int2 NOT NULL DEFAULT 0,
  "last_seen_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."switch_remote"."switch_id" IS '远程开关id';
COMMENT ON COLUMN "public"."switch_remote"."switch_code" IS '远程开关编号';
COMMENT ON COLUMN "public"."switch_remote"."switch_name" IS '远程开关名称';
COMMENT ON COLUMN "public"."switch_remote"."tower_id" IS '风机id';
COMMENT ON COLUMN "public"."switch_remote"."switch_ip" IS '远程开关ip';
COMMENT ON COLUMN "public"."switch_remote"."port" IS '远程开关端口';
COMMENT ON COLUMN "public"."switch_remote"."protocol" IS '远程开关协议';
COMMENT ON COLUMN "public"."switch_remote"."cmd_restart" IS '远程开关重启命令';
COMMENT ON COLUMN "public"."switch_remote"."enabled" IS '远程开关是否可用  0 可用 1 故障';
COMMENT ON COLUMN "public"."switch_remote"."last_seen_at" IS '上次执行时间';
COMMENT ON COLUMN "public"."switch_remote"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."switch_remote"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."switch_remote"."remark" IS '备注';
COMMENT ON TABLE "public"."switch_remote" IS '远程开关管理表';

-- ----------------------------
-- Records of switch_remote
-- ----------------------------
INSERT INTO "public"."switch_remote" VALUES (7, '', '7#远程开关1', 5, '127.0.0.1', 1037, 'tcp', '!RST\r', 1, '2025-08-21 14:53:28', '2025-08-11 15:15:35', 0, '');
INSERT INTO "public"."switch_remote" VALUES (4, '', '4#远程开关2', 3, '127.0.0.1', 1034, 'tcp', '!RST\r', 1, '2025-08-21 14:53:37', '2025-08-11 15:13:22', 0, '');
INSERT INTO "public"."switch_remote" VALUES (6, '', '5#远程开关2', 4, '127.0.0.1', 1036, 'tcp', '!RST\r', 1, '2025-08-21 14:53:55', '2025-08-11 15:14:55', 0, '');
INSERT INTO "public"."switch_remote" VALUES (5, '', '5#远程开关1', 4, '127.0.0.1', 1035, 'tcp', '!RST\r', 1, '2025-08-21 14:54:04', '2025-08-11 15:14:38', 0, '');
INSERT INTO "public"."switch_remote" VALUES (3, '', '4#远程开关1', 3, '127.0.0.1', 1033, 'tcp', '!RST\r', 1, '2025-08-21 14:56:20', '2025-08-11 15:12:59', 0, '');
INSERT INTO "public"."switch_remote" VALUES (8, '', '7#远程开关2', 5, '127.0.0.1', 1038, 'tcp', '!RST\r', 0, '2025-08-11 15:15:54', '2025-08-11 15:15:54', 0, '');
INSERT INTO "public"."switch_remote" VALUES (2, '', '3#远程开关1', 2, '127.0.0.1', 1032, 'tcp', '!RST\r', 1, '2025-08-21 15:03:24', '2025-08-11 15:12:13', 0, '');
INSERT INTO "public"."switch_remote" VALUES (1, '', '1#远程开关1', 1, '127.0.0.1', 1031, 'tcp', '!RST\r', 1, '2025-08-21 15:56:25', '2025-08-11 15:11:18', 0, '');

-- ----------------------------
-- Table structure for sys_departments
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_departments";
CREATE TABLE "public"."sys_departments" (
  "depart_id" int4 NOT NULL DEFAULT nextval('sys_depart_id_seq'::regclass),
  "depart_name" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "create_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."sys_departments"."depart_id" IS '部门ID';
COMMENT ON COLUMN "public"."sys_departments"."depart_name" IS '部门名称';
COMMENT ON COLUMN "public"."sys_departments"."create_time" IS '部门创建时间';
COMMENT ON COLUMN "public"."sys_departments"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_departments"."is_delete" IS '此部门是否已删除,1已删除 0未删除';
COMMENT ON TABLE "public"."sys_departments" IS '部门信息表';

-- ----------------------------
-- Records of sys_departments
-- ----------------------------
INSERT INTO "public"."sys_departments" VALUES (64, '产品设计4', '2025-07-09 11:31:40', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (65, 'test2', '2025-07-09 15:11:24', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (1, '系统研发部', '2025-07-09 09:17:26', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (21, '应用研发部', '2025-07-09 10:59:40', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (22, '智慧城市bg', '2025-07-09 11:00:28', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (23, '智慧教育bg', '2025-07-09 11:00:45', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (41, '页面设计部门', '2025-07-09 11:12:36', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (61, 'test2', '2025-07-09 11:18:03', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (62, '产品设计2', '2025-07-09 11:19:05', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (63, '产品设计3', '2025-07-09 11:20:47', '-', 1);
INSERT INTO "public"."sys_departments" VALUES (67, '研发部', '2025-07-11 14:14:19', '-', 0);
INSERT INTO "public"."sys_departments" VALUES (66, '财务部', '2025-07-11 14:13:15', '-', 0);

-- ----------------------------
-- Table structure for sys_login_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_login_logs";
CREATE TABLE "public"."sys_login_logs" (
  "id" int4 NOT NULL DEFAULT nextval('sys_login_log_record_id_seq'::regclass),
  "user_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "status" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "ip_addr" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "login_location" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "browser" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "os" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "platform" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "login_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "msg" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "dept_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."sys_login_logs"."id" IS '记录ID';
COMMENT ON COLUMN "public"."sys_login_logs"."user_name" IS '用户名';
COMMENT ON COLUMN "public"."sys_login_logs"."status" IS '状态';
COMMENT ON COLUMN "public"."sys_login_logs"."ip_addr" IS '登录ip地址';
COMMENT ON COLUMN "public"."sys_login_logs"."login_location" IS '登录位置,归属地';
COMMENT ON COLUMN "public"."sys_login_logs"."browser" IS '浏览器';
COMMENT ON COLUMN "public"."sys_login_logs"."os" IS '操作系统';
COMMENT ON COLUMN "public"."sys_login_logs"."platform" IS '固件';
COMMENT ON COLUMN "public"."sys_login_logs"."login_time" IS '登录时间';
COMMENT ON COLUMN "public"."sys_login_logs"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_login_logs"."msg" IS '信息';
COMMENT ON COLUMN "public"."sys_login_logs"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_login_logs"."updated_at" IS '最后更新时间';
COMMENT ON COLUMN "public"."sys_login_logs"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."sys_login_logs"."dept_name" IS '部门名称';
COMMENT ON TABLE "public"."sys_login_logs" IS '登录日志表';

-- ----------------------------
-- Records of sys_login_logs
-- ----------------------------
INSERT INTO "public"."sys_login_logs" VALUES (1, '管理员', '2', '127.0.0.1:56473', '192.168.188.1', 'PostmanRuntime 7.45.0', '', '', '2025-08-04 10:20:12', 'PostmanRuntime/7.45.0', '登录成功', '2025-08-04 10:20:12', '2025-08-04 10:20:12', 0, '系统研发部');
INSERT INTO "public"."sys_login_logs" VALUES (2, '管理员', '2', '172.22.0.2:52120', '172.16.90.70', 'Chrome 138.0.0.0', 'Windows 10', 'Windows', '2025-08-04 14:02:56', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36', '登录成功', '2025-08-04 14:02:56', '2025-08-04 14:02:56', 0, '系统研发部');
INSERT INTO "public"."sys_login_logs" VALUES (3, '管理员', '2', '172.22.0.2:36196', '172.16.90.70', 'Chrome 138.0.0.0', 'Windows 10', 'Windows', '2025-08-05 09:20:12', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36', '登录成功', '2025-08-05 09:20:12', '2025-08-05 09:20:12', 0, '系统研发部');
-- ----------------------------
-- Table structure for sys_menus
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_menus";
CREATE TABLE "public"."sys_menus" (
  "menu_id" int4 NOT NULL DEFAULT nextval('sys_menu_menu_id_seq'::regclass),
  "fid" int8 NOT NULL DEFAULT 0,
  "route" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "cn_name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "order" int8 NOT NULL DEFAULT 9999,
  "status" int2 NOT NULL DEFAULT 0,
  "normal_auth" int2 NOT NULL DEFAULT 0,
  "en_name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "create_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "remark" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "path" varchar(600) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "perms" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "icon" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "menu_type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."sys_menus"."menu_id" IS '菜单表的唯一标识';
COMMENT ON COLUMN "public"."sys_menus"."fid" IS '父页面的id,如果是一级目录,则为0';
COMMENT ON COLUMN "public"."sys_menus"."route" IS '路由';
COMMENT ON COLUMN "public"."sys_menus"."cn_name" IS '菜单节点的中文名';
COMMENT ON COLUMN "public"."sys_menus"."order" IS '桌面图标排序';
COMMENT ON COLUMN "public"."sys_menus"."status" IS '页面是否展示 0 展示 1 不展示';
COMMENT ON COLUMN "public"."sys_menus"."normal_auth" IS '普通用户是否有权限 1无权限  0 有权限';
COMMENT ON COLUMN "public"."sys_menus"."en_name" IS '菜单节点的英文名';
COMMENT ON COLUMN "public"."sys_menus"."create_time" IS '菜单创建时间';
COMMENT ON COLUMN "public"."sys_menus"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_menus"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."sys_menus"."path" IS '菜单链接地址';
COMMENT ON COLUMN "public"."sys_menus"."perms" IS '权限字段';
COMMENT ON COLUMN "public"."sys_menus"."icon" IS '图标';
COMMENT ON COLUMN "public"."sys_menus"."menu_type" IS '菜单类型';
COMMENT ON TABLE "public"."sys_menus" IS '菜单信息表';

-- ----------------------------
-- Records of sys_menus
-- ----------------------------
INSERT INTO "public"."sys_menus" VALUES (4, 1, '/', '测试biyan', 9999, 0, 0, '', '2025-04-16 09:10:25', '-', 1, '/testbiyan', '', '864654', '');
INSERT INTO "public"."sys_menus" VALUES (10, 7, '', '组织管理', 9999, 0, 0, '', '2025-04-18 11:27:41', '-', 0, '/settings/department', '', 'menu-setting', '');
INSERT INTO "public"."sys_menus" VALUES (11, 7, '', '权限管理', 9999, 0, 0, '', '2025-04-18 11:27:58', '-', 0, '/settings/permission', '', 'menu-setting', '');
INSERT INTO "public"."sys_menus" VALUES (12, 1, '', '日志管理', 9999, 0, 0, '', '2025-04-18 11:28:16', '-', 0, '/logs', '', 'menu-logs', '');
INSERT INTO "public"."sys_menus" VALUES (13, 12, '', '登录日志', 9999, 0, 0, '', '2025-04-18 11:28:36', '-', 0, '/logs/login', '', 'menu-logs', '');
INSERT INTO "public"."sys_menus" VALUES (14, 12, '', '操作日志', 9999, 0, 0, '', '2025-04-18 11:28:51', '-', 0, '/logs/operation', '', 'menu-logs', '');
INSERT INTO "public"."sys_menus" VALUES (1, 0, '', '风电数据分析', 1, 0, 1, '', '2025-04-16 09:05:35', '-', 0, '/', '', '/', '');
INSERT INTO "public"."sys_menus" VALUES (5, 1, '/dashboard', '工作台', 9999, 0, 0, '', '2025-04-18 11:26:00', '-', 0, '/dashboard', '', 'menu-dashboard', '');
INSERT INTO "public"."sys_menus" VALUES (7, 1, '/system', '系统设置', 9999, 0, 0, '', '2025-04-18 11:26:39', '-', 0, '/setting', '', 'menu-setting', '');
INSERT INTO "public"."sys_menus" VALUES (8, 7, '/menu', '菜单管理', 9999, 0, 0, '', '2025-04-18 11:26:59', '-', 0, '/settings/menu', '', 'menu-setting', '');
INSERT INTO "public"."sys_menus" VALUES (9, 7, '/user', '人员管理', 9999, 0, 0, '', '2025-04-18 11:27:23', '-', 0, '/settings/user', '', 'menu-setting', '');
INSERT INTO "public"."sys_menus" VALUES (3, 1, '/test', '测试', 9999, 0, 0, '', '2025-04-16 09:09:30', '-', 1, '/test', '', '5345', '');
INSERT INTO "public"."sys_menus" VALUES (2, 1, '', '风场信息', 9999, 0, 0, '', '2025-11-28 09:08:44', '-', 0, '/wind', '', 'menu-type', '');
INSERT INTO "public"."sys_menus" VALUES (15, 1, '', '设备管理', 9999, 0, 0, '', '2025-11-28 09:16:46', '-', 0, '/equipment', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (16, 2, '', '风场管理', 9999, 0, 0, '', '2025-11-28 09:17:36', '-', 0, '/wind/station', '', 'menu-type', '');
INSERT INTO "public"."sys_menus" VALUES (17, 2, '', '风机列表', 9999, 0, 0, '', '2025-11-28 09:17:55', '-', 0, '/wind/list', '', 'menu-type', '');
INSERT INTO "public"."sys_menus" VALUES (18, 15, '', '结构类型', 9999, 0, 0, '', '2025-11-28 09:18:19', '-', 0, '/equipment/structure', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (19, 15, '', '设备分类', 9999, 0, 0, '', '2025-11-28 09:18:37', '-', 0, '/equipment/classify', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (20, 15, '', '设备列表', 9999, 0, 0, '', '2025-11-28 09:18:56', '-', 0, '/equipment/sensor', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (21, 1, '', '远程开关', 9999, 0, 0, '', '2025-11-28 09:19:25', '-', 0, '/switch', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (22, 21, '', '开关列表', 9999, 0, 0, '', '2025-11-28 09:19:44', '-', 0, '/switch/list', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (23, 21, '', '资产列表', 9999, 0, 0, '', '2025-11-28 09:20:04', '-', 0, '/switch/asset', '', 'menu-equipment', '');
INSERT INTO "public"."sys_menus" VALUES (6, 1, '', '识别记录', 9999, 0, 0, '', '2025-04-18 11:26:24', '-', 1, '/record', '', 'menu-record', '');
INSERT INTO "public"."sys_menus" VALUES (24, 1, '', '数据管理', 9999, 0, 0, '', '2025-11-28 09:20:49', '-', 0, '/datamanage', '', 'menu-record', '');
INSERT INTO "public"."sys_menus" VALUES (25, 24, '', '数据查询', 9999, 0, 0, '', '2025-11-28 09:21:10', '-', 0, '/datamanage/list', '', 'menu-record', '');
INSERT INTO "public"."sys_menus" VALUES (26, 24, '', '数据分析', 9999, 0, 0, '', '2025-11-28 09:21:27', '-', 0, '/datamanage/statistics', '', 'menu-record', '');

-- ----------------------------
-- Table structure for sys_operate_logs
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_operate_logs";
CREATE TABLE "public"."sys_operate_logs" (
  "id" int4 NOT NULL DEFAULT nextval('sys_operate_log_record_id_seq'::regclass),
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "business_type" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "method" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "request_method" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "operator_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "dept_name" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "operator_url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "operator_ip" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "operator_location" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "status" int2 NOT NULL DEFAULT 0,
  "operator_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "latency_time" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "user_agent" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."sys_operate_logs"."id" IS '记录ID';
COMMENT ON COLUMN "public"."sys_operate_logs"."title" IS '操作模块';
COMMENT ON COLUMN "public"."sys_operate_logs"."business_type" IS '操作类型';
COMMENT ON COLUMN "public"."sys_operate_logs"."method" IS '函数';
COMMENT ON COLUMN "public"."sys_operate_logs"."request_method" IS '请求方式 GET POST PUT DELETE';
COMMENT ON COLUMN "public"."sys_operate_logs"."operator_name" IS '操作者';
COMMENT ON COLUMN "public"."sys_operate_logs"."dept_name" IS '部门名称';
COMMENT ON COLUMN "public"."sys_operate_logs"."operator_url" IS '访问地址';
COMMENT ON COLUMN "public"."sys_operate_logs"."operator_ip" IS '客户端ip';
COMMENT ON COLUMN "public"."sys_operate_logs"."operator_location" IS '访问位置';
COMMENT ON COLUMN "public"."sys_operate_logs"."status" IS '操作状态 0:正常 1:关闭';
COMMENT ON COLUMN "public"."sys_operate_logs"."operator_time" IS '操作时间';
COMMENT ON COLUMN "public"."sys_operate_logs"."latency_time" IS '耗时';
COMMENT ON COLUMN "public"."sys_operate_logs"."user_agent" IS 'ua';
COMMENT ON COLUMN "public"."sys_operate_logs"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_operate_logs"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."sys_operate_logs"."remark" IS '备注';
COMMENT ON TABLE "public"."sys_operate_logs" IS '操作日志表';

-- ----------------------------
-- Records of sys_operate_logs
-- ----------------------------
INSERT INTO "public"."sys_operate_logs" VALUES (20, '设备管理', '新增', '新增设备', 'POST', '管理员', '系统研发部', '/api/device/device/addDevice', '127.0.0.1', '192.168.188.1', 0, '2025-08-04 11:13:13', '76.3581ms', 'PostmanRuntime/7.45.0', '2025-08-04 11:13:14', 0, '');
INSERT INTO "public"."sys_operate_logs" VALUES (21, '设备管理', '新增', '新增设备', 'POST', '管理员', '系统研发部', '/api/device/device/addDevice', '127.0.0.1', '192.168.188.1', 0, '2025-08-04 11:15:13', '76.9707ms', 'PostmanRuntime/7.45.0', '2025-08-04 11:15:14', 0, '');
INSERT INTO "public"."sys_operate_logs" VALUES (22, '设备管理', '更新', '更新设备实例', 'POST', '管理员', '系统研发部', '/api/device/device/editDevice', '127.0.0.1', '192.168.188.1', 0, '2025-08-04 11:17:26', '39.6784ms', 'PostmanRuntime/7.45.0', '2025-08-04 11:17:27', 0, '');

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_role";
CREATE TABLE "public"."sys_role" (
  "id" int4 NOT NULL DEFAULT nextval('sys_role_id_seq'::regclass),
  "role_id" int4 NOT NULL DEFAULT 2,
  "role_name" varchar(200) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "role_auth" json NOT NULL,
  "detail" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "role_create_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "alter_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."sys_role"."id" IS '角色表的唯一标识';
COMMENT ON COLUMN "public"."sys_role"."role_id" IS '角色id 默认为 2 普通用户';
COMMENT ON COLUMN "public"."sys_role"."role_name" IS '角色名称';
COMMENT ON COLUMN "public"."sys_role"."role_auth" IS '角色权限';
COMMENT ON COLUMN "public"."sys_role"."detail" IS '角色描述';
COMMENT ON COLUMN "public"."sys_role"."role_create_time" IS '角色创建时间';
COMMENT ON COLUMN "public"."sys_role"."alter_time" IS '角色修改时间';
COMMENT ON COLUMN "public"."sys_role"."is_delete" IS '此角色是否已删除,1 已删除 0 未删除';
COMMENT ON TABLE "public"."sys_role" IS '角色信息表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO "public"."sys_role" VALUES (5, 4, '测试3', '{"1": 1, "3": 1, "4": 0}', '', '2025-04-16 09:44:24', '2025-04-16 09:44:24', 1);
INSERT INTO "public"."sys_role" VALUES (1, 2, '普通用户', '{"1": 1, "5": 1, "6": 1, "7": 0, "8": 0, "9": 0, "10": 0, "11": 0, "12": 0, "13": 0, "14": 0}', '默认权限', '2025-04-11 06:36:33', '2025-04-11 06:36:33', 0);
INSERT INTO "public"."sys_role" VALUES (4, 1, '管理员', '{"1": 1, "5": 1, "6": 1, "7": 1, "8": 1, "9": 1, "10": 1, "11": 1, "12": 1, "13": 1, "14": 1}', '', '2025-04-16 09:43:31', '2025-04-16 09:43:31', 0);
INSERT INTO "public"."sys_role" VALUES (2, 1, 'UI测试3', '{"1":1,"10":0,"11":0,"12":0,"13":0,"14":0,"3":1,"5":0,"6":0,"7":0,"8":0,"9":0}', '', '2025-07-10 09:43:46', '2025-07-10 09:43:46', 0);

-- ----------------------------
-- Table structure for sys_users
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_users";
CREATE TABLE "public"."sys_users" (
  "uid" int4 NOT NULL DEFAULT nextval('sys_users_uid_seq'::regclass),
  "account" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "password" char(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '123456'::bpchar,
  "user_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "depart_id" int8 NOT NULL DEFAULT 0,
  "depart_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "is_admin" int2 NOT NULL DEFAULT 0,
  "user_create_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "last_login_time" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "user_role_id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '2'::character varying,
  "user_role_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '普通用户'::character varying,
  "sex" varchar(10) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '1'::character varying,
  "phone" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "email" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying
)
;
COMMENT ON COLUMN "public"."sys_users"."uid" IS '用户ID';
COMMENT ON COLUMN "public"."sys_users"."account" IS '用户名或账号';
COMMENT ON COLUMN "public"."sys_users"."password" IS '密码';
COMMENT ON COLUMN "public"."sys_users"."user_name" IS '用户名';
COMMENT ON COLUMN "public"."sys_users"."depart_id" IS '部门id';
COMMENT ON COLUMN "public"."sys_users"."depart_name" IS '部门名称';
COMMENT ON COLUMN "public"."sys_users"."is_admin" IS '是否为超管, 1 超管 0 不是超管';
COMMENT ON COLUMN "public"."sys_users"."user_create_time" IS '用户注册时间';
COMMENT ON COLUMN "public"."sys_users"."last_login_time" IS '用户最后登录时间';
COMMENT ON COLUMN "public"."sys_users"."is_delete" IS '此用户是否已删除,1已删除';
COMMENT ON COLUMN "public"."sys_users"."user_role_id" IS '角色id,如果为0则无任何权限,2 普通用户 1 超管';
COMMENT ON COLUMN "public"."sys_users"."user_role_name" IS '用户角色名';
COMMENT ON COLUMN "public"."sys_users"."sex" IS '用户性别 1 男 2 女';
COMMENT ON COLUMN "public"."sys_users"."phone" IS '用户手机号';
COMMENT ON COLUMN "public"."sys_users"."email" IS '用户邮箱号';
COMMENT ON COLUMN "public"."sys_users"."remark" IS '备注';
COMMENT ON TABLE "public"."sys_users" IS '用户信息表';

-- ----------------------------
-- Records of sys_users
-- ----------------------------
INSERT INTO "public"."sys_users" VALUES (22, 'root', '123456                          ', 'tianyx', 1, '系统研发部', 0, '2025-07-10 09:29:35', '2025-07-10 09:29:35', 0, '1', '管理员', '2', '2354365476', '', '-');
INSERT INTO "public"."sys_users" VALUES (21, 'tian', '123456                          ', 'tianyuanxiang', 1, '系统研发部', 0, '2025-07-09 16:14:51', '2025-07-09 16:14:51', 1, '1', '普通用户', '2', '2354365476', '', '-');
INSERT INTO "public"."sys_users" VALUES (1, 'admin', 'Jl147258.                       ', '管理员', 1, '', 1, '2025-07-09 09:17:05', '2025-07-09 09:17:05', 0, '2', 'UI测试3', '1', '16666666666', '', '-');

-- ----------------------------
-- Table structure for towers
-- ----------------------------
DROP TABLE IF EXISTS "public"."towers";
CREATE TABLE "public"."towers" (
  "tower_id" int4 NOT NULL DEFAULT nextval('tower_id_seq'::regclass),
  "tower_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '01'::character varying,
  "farm_id" int4 NOT NULL DEFAULT 1,
  "farm_code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "farm_name" varchar(32) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "is_delete" int2 NOT NULL DEFAULT 0,
  "remark" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '-'::character varying,
  "longitude" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "latitude" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "models" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."towers"."tower_id" IS '记录ID';
COMMENT ON COLUMN "public"."towers"."tower_code" IS '风机名称，如 “04” 4号风机';
COMMENT ON COLUMN "public"."towers"."farm_id" IS '对应所属的风场id';
COMMENT ON COLUMN "public"."towers"."farm_code" IS '所属风场编号，例: FY';
COMMENT ON COLUMN "public"."towers"."farm_name" IS '所属风场名称，例: 扶余：市级单位缩写';
COMMENT ON COLUMN "public"."towers"."created_at" IS '记录创建时间';
COMMENT ON COLUMN "public"."towers"."is_delete" IS '是否删除 0 未删除 1 已删除';
COMMENT ON COLUMN "public"."towers"."remark" IS '备注';
COMMENT ON COLUMN "public"."towers"."longitude" IS '风机经度';
COMMENT ON COLUMN "public"."towers"."latitude" IS '风机纬度';
COMMENT ON COLUMN "public"."towers"."models" IS '风机绑定的模型编号';
COMMENT ON TABLE "public"."towers" IS '风机表';

-- ----------------------------
-- Records of towers
-- ----------------------------
INSERT INTO "public"."towers" VALUES (5, '07', 1, 'FY', '华电万兴风电场', '2025-07-31 14:11:17', 0, '-', '125.642331', '45.227637', '');
INSERT INTO "public"."towers" VALUES (2, '03', 1, 'FY', '华电万兴风电场', '2025-07-31 14:10:40', 0, '-', '125.642831', '45.227637', '');
INSERT INTO "public"."towers" VALUES (9, '01', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:20:27', 0, '-', '126.532532', '44.844675', '');
INSERT INTO "public"."towers" VALUES (10, '02', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:21:00', 0, '-', '126.532532', '44.844675', '');
INSERT INTO "public"."towers" VALUES (11, '03', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:21:24', 0, '-', '126.532532', '44.844675', '');
INSERT INTO "public"."towers" VALUES (13, '06', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:22:06', 0, '-', '126.532532', '44.844675', '');
INSERT INTO "public"."towers" VALUES (7, '02', 2, 'YS', '华电龙岗风电场', '2025-08-04 10:27:21', 1, '-', '', '', '');
INSERT INTO "public"."towers" VALUES (8, '09', 2, 'YS', '华电龙岗风电场', '2025-08-22 14:55:50', 1, '-', '126.532532', '44.844675', '');
INSERT INTO "public"."towers" VALUES (12, '04', 2, 'YS', '华电龙岗风电场', '2025-09-15 09:21:44', 0, '-', '126.532532', '44.844675', '');
INSERT INTO "public"."towers" VALUES (3, '04', 1, 'FY', '华电万兴风电场', '2025-07-31 14:10:53', 0, '-', '125.642931', '45.227637', '1,2,3,4,6,5');
INSERT INTO "public"."towers" VALUES (4, '05', 1, 'FY', '华电万兴风电场', '2025-07-31 14:11:05', 0, '-', '125.642431', '45.227637', '2,3,1');
INSERT INTO "public"."towers" VALUES (1, '01', 1, 'FY', '华电万兴风电场', '2025-07-31 13:58:44', 0, '-', '125.642731', '45.227637', '4,5');

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."alarm_type_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."camera_record_id_seq"', 24918, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."device_id_seq"', 191, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."device_meta_id_seq"', 35, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."device_type_id_seq"', 13, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."farm_id_seq"', 5, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."map_id_seq"', 18, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."models_id_seq"', 8192, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."record_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."remote_switch_id_seq"', 8, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."structure_id_seq"', 17, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."switch_asset_id_seq"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sys_depart_id_seq"', 2, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sys_login_log_record_id_seq"', 226, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sys_menu_menu_id_seq"', 26, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sys_operate_log_record_id_seq"', 472, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sys_role_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."sys_users_uid_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."tower_id_seq"', 13, true);

-- ----------------------------
-- Primary Key structure for table alarm_type
-- ----------------------------
ALTER TABLE "public"."alarm_type" ADD CONSTRAINT "alarm_type_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table camera_record
-- ----------------------------
ALTER TABLE "public"."camera_record" ADD CONSTRAINT "camera_record_pkey" PRIMARY KEY ("record_id");

-- ----------------------------
-- Primary Key structure for table camera_record_test
-- ----------------------------
ALTER TABLE "public"."camera_record_test" ADD CONSTRAINT "camera_record_test_pkey" PRIMARY KEY ("record_id");

-- ----------------------------
-- Primary Key structure for table device
-- ----------------------------
ALTER TABLE "public"."device" ADD CONSTRAINT "device_pkey" PRIMARY KEY ("device_id");

-- ----------------------------
-- Primary Key structure for table device_meta
-- ----------------------------
ALTER TABLE "public"."device_meta" ADD CONSTRAINT "device_meta_pkey" PRIMARY KEY ("meta_id");

-- ----------------------------
-- Primary Key structure for table device_type
-- ----------------------------
ALTER TABLE "public"."device_type" ADD CONSTRAINT "device_type_pkey" PRIMARY KEY ("device_type_id");

-- ----------------------------
-- Primary Key structure for table farms
-- ----------------------------
ALTER TABLE "public"."farms" ADD CONSTRAINT "farms_pkey" PRIMARY KEY ("farm_id");

-- ----------------------------
-- Primary Key structure for table model_device_map
-- ----------------------------
ALTER TABLE "public"."model_device_map" ADD CONSTRAINT "model_device_map_pkey" PRIMARY KEY ("map_id");

-- ----------------------------
-- Primary Key structure for table model_records
-- ----------------------------
ALTER TABLE "public"."model_records" ADD CONSTRAINT "model_records_pkey" PRIMARY KEY ("record_id");

-- ----------------------------
-- Primary Key structure for table models
-- ----------------------------
ALTER TABLE "public"."models" ADD CONSTRAINT "models_pkey" PRIMARY KEY ("model_id");

-- ----------------------------
-- Primary Key structure for table structure_type
-- ----------------------------
ALTER TABLE "public"."structure_type" ADD CONSTRAINT "structure_type_pkey" PRIMARY KEY ("structure_id");

-- ----------------------------
-- Primary Key structure for table switch_asset
-- ----------------------------
ALTER TABLE "public"."switch_asset" ADD CONSTRAINT "switch_asset_pkey" PRIMARY KEY ("asset_id");

-- ----------------------------
-- Primary Key structure for table switch_remote
-- ----------------------------
ALTER TABLE "public"."switch_remote" ADD CONSTRAINT "remote_switch_pkey" PRIMARY KEY ("switch_id");

-- ----------------------------
-- Primary Key structure for table sys_departments
-- ----------------------------
ALTER TABLE "public"."sys_departments" ADD CONSTRAINT "sys_depart_pkey" PRIMARY KEY ("depart_id");

-- ----------------------------
-- Primary Key structure for table sys_login_logs
-- ----------------------------
ALTER TABLE "public"."sys_login_logs" ADD CONSTRAINT "sys_login_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sys_menus
-- ----------------------------
ALTER TABLE "public"."sys_menus" ADD CONSTRAINT "sys_menu_pkey" PRIMARY KEY ("menu_id");

-- ----------------------------
-- Primary Key structure for table sys_operate_logs
-- ----------------------------
ALTER TABLE "public"."sys_operate_logs" ADD CONSTRAINT "sys_operate_logs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table sys_role
-- ----------------------------
ALTER TABLE "public"."sys_role" ADD CONSTRAINT "sys_role_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_users
-- ----------------------------
CREATE INDEX "idx_depart_id" ON "public"."sys_users" USING btree (
  "depart_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_users
-- ----------------------------
ALTER TABLE "public"."sys_users" ADD CONSTRAINT "sys_users_pkey" PRIMARY KEY ("uid");

-- ----------------------------
-- Primary Key structure for table towers
-- ----------------------------
ALTER TABLE "public"."towers" ADD CONSTRAINT "towers_pkey" PRIMARY KEY ("tower_id");
