/*
 Navicat Premium Data Transfer

 Source Server         : 云服务器-8.140.204.47
 Source Server Type    : PostgreSQL
 Source Server Version : 140017
 Source Host           : 8.140.204.47:5432
 Source Catalog        : system
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140017
 File Encoding         : 65001

 Date: 17/05/2026 12:31:30
*/


-- ----------------------------
-- Sequence structure for casbin_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."casbin_rule_id_seq";
CREATE SEQUENCE "public"."casbin_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_api_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_api_id_seq";
CREATE SEQUENCE "public"."sys_api_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_dict_data_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_dict_data_id_seq";
CREATE SEQUENCE "public"."sys_dict_data_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_dict_type_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_dict_type_id_seq";
CREATE SEQUENCE "public"."sys_dict_type_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_file_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_file_id_seq";
CREATE SEQUENCE "public"."sys_file_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_login_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_login_log_id_seq";
CREATE SEQUENCE "public"."sys_login_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_menu_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_menu_id_seq";
CREATE SEQUENCE "public"."sys_menu_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_oper_log_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_oper_log_id_seq";
CREATE SEQUENCE "public"."sys_oper_log_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_role_api_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_role_api_id_seq";
CREATE SEQUENCE "public"."sys_role_api_id_seq" 
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
-- Sequence structure for sys_role_menu_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_role_menu_id_seq";
CREATE SEQUENCE "public"."sys_role_menu_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_user_id_seq";
CREATE SEQUENCE "public"."sys_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for sys_user_role_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sys_user_role_id_seq";
CREATE SEQUENCE "public"."sys_user_role_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."casbin_rule";
CREATE TABLE "public"."casbin_rule" (
  "id" int8 NOT NULL DEFAULT nextval('casbin_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default",
  "v0" varchar(100) COLLATE "pg_catalog"."default",
  "v1" varchar(100) COLLATE "pg_catalog"."default",
  "v2" varchar(100) COLLATE "pg_catalog"."default",
  "v3" varchar(100) COLLATE "pg_catalog"."default",
  "v4" varchar(100) COLLATE "pg_catalog"."default",
  "v5" varchar(100) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."casbin_rule"."id" IS '主键';
COMMENT ON TABLE "public"."casbin_rule" IS 'Casbin权限规则表';

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO "public"."casbin_rule" VALUES (1, 'p', 'admin', '/api/*', '*', NULL, NULL, NULL);
INSERT INTO "public"."casbin_rule" VALUES (19, 'p', 'manager', '/api/auth/login', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (20, 'p', 'manager', '/api/auth/refresh', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (21, 'p', 'manager', '/api/auth/changePassword', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (22, 'p', 'manager', '/api/auth/logout', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (23, 'p', 'manager', '/api/auth/userInfo', 'GET', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (24, 'p', 'manager', '/api/system/user', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (25, 'p', 'manager', '/api/system/user', 'GET', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (26, 'p', 'manager', '/api/system/user/:id', 'PUT', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (27, 'p', 'manager', '/api/system/user/:id', 'GET', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (28, 'p', 'manager', '/api/system/user/:id/reset-password', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (29, 'p', 'manager', '/api/system/role', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (30, 'p', 'manager', '/api/system/role', 'GET', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (31, 'p', 'manager', '/api/system/role/:id', 'PUT', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (32, 'p', 'manager', '/api/system/role/:id', 'GET', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (33, 'p', 'manager', '/api/system/role/:id/permissions', 'PUT', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (34, 'p', 'manager', '/api/system/role/all', 'GET', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (35, 'p', 'manager', '/api/system/menu', 'POST', '', '', '');
INSERT INTO "public"."casbin_rule" VALUES (36, 'p', 'manager', '/api/system/menu/tree', 'GET', '', '', '');

-- ----------------------------
-- Table structure for sys_api
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_api";
CREATE TABLE "public"."sys_api" (
  "id" int8 NOT NULL DEFAULT nextval('sys_api_id_seq'::regclass),
  "api_path" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
  "method" varchar(16) COLLATE "pg_catalog"."default" NOT NULL,
  "api_group" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "description" varchar(256) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6),
  "api_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."sys_api"."id" IS '接口ID';
COMMENT ON COLUMN "public"."sys_api"."api_path" IS '接口路径，如：/api/system/user';
COMMENT ON COLUMN "public"."sys_api"."method" IS 'HTTP方法：GET/POST/PUT/DELETE';
COMMENT ON COLUMN "public"."sys_api"."api_group" IS '接口分组，如：系统管理';
COMMENT ON COLUMN "public"."sys_api"."description" IS '接口描述';
COMMENT ON COLUMN "public"."sys_api"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_api"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sys_api"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON COLUMN "public"."sys_api"."api_name" IS '接口中文名称';
COMMENT ON TABLE "public"."sys_api" IS '系统接口权限表';

-- ----------------------------
-- Records of sys_api
-- ----------------------------
INSERT INTO "public"."sys_api" VALUES (1, '/api/auth/login', 'POST', '认证管理', '用户登录获取Token', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '用户登录');
INSERT INTO "public"."sys_api" VALUES (2, '/api/auth/refresh', 'POST', '认证管理', '刷新访问令牌', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '刷新令牌');
INSERT INTO "public"."sys_api" VALUES (3, '/api/auth/changePassword', 'POST', '认证管理', '登录用户修改密码', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '修改密码');
INSERT INTO "public"."sys_api" VALUES (4, '/api/auth/logout', 'POST', '认证管理', '退出登录并注销Token', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '退出登录');
INSERT INTO "public"."sys_api" VALUES (5, '/api/auth/userInfo', 'GET', '认证管理', '获取当前登录用户信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '获取用户信息');
INSERT INTO "public"."sys_api" VALUES (6, '/api/system/user', 'POST', '用户管理', '创建系统用户', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增用户');
INSERT INTO "public"."sys_api" VALUES (7, '/api/system/user', 'GET', '用户管理', '分页查询用户列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '用户列表');
INSERT INTO "public"."sys_api" VALUES (8, '/api/system/user/:id', 'PUT', '用户管理', '修改用户基本信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑用户');
INSERT INTO "public"."sys_api" VALUES (9, '/api/system/user/:id', 'DELETE', '用户管理', '软删除用户', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除用户');
INSERT INTO "public"."sys_api" VALUES (10, '/api/system/user/:id', 'GET', '用户管理', '查询单个用户信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '用户详情');
INSERT INTO "public"."sys_api" VALUES (11, '/api/system/user/:id/reset-password', 'POST', '用户管理', '重置用户密码', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '重置密码');
INSERT INTO "public"."sys_api" VALUES (12, '/api/system/role', 'POST', '角色管理', '创建系统角色', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增角色');
INSERT INTO "public"."sys_api" VALUES (13, '/api/system/role', 'GET', '角色管理', '分页查询角色列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '角色列表');
INSERT INTO "public"."sys_api" VALUES (14, '/api/system/role/:id', 'PUT', '角色管理', '修改角色信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑角色');
INSERT INTO "public"."sys_api" VALUES (15, '/api/system/role/:id', 'DELETE', '角色管理', '删除角色', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除角色');
INSERT INTO "public"."sys_api" VALUES (16, '/api/system/role/:id', 'GET', '角色管理', '查询单个角色信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '角色详情');
INSERT INTO "public"."sys_api" VALUES (17, '/api/system/role/:id/permissions', 'PUT', '角色管理', '分配角色的菜单和接口权限', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '分配权限');
INSERT INTO "public"."sys_api" VALUES (18, '/api/system/role/all', 'GET', '角色管理', '查询所有角色（下拉选择用）', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '全部角色');
INSERT INTO "public"."sys_api" VALUES (19, '/api/system/menu', 'POST', '菜单管理', '创建菜单或按钮', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增菜单');
INSERT INTO "public"."sys_api" VALUES (20, '/api/system/menu/tree', 'GET', '菜单管理', '获取完整菜单树（管理页用）', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '菜单树');
INSERT INTO "public"."sys_api" VALUES (21, '/api/system/menu/current', 'GET', '菜单管理', '获取当前登录用户的菜单树', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '当前用户菜单');
INSERT INTO "public"."sys_api" VALUES (22, '/api/system/menu/:id', 'PUT', '菜单管理', '修改菜单信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑菜单');
INSERT INTO "public"."sys_api" VALUES (23, '/api/system/menu/:id', 'DELETE', '菜单管理', '软删除菜单', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除菜单');
INSERT INTO "public"."sys_api" VALUES (24, '/api/system/api', 'POST', '接口管理', '创建系统接口', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增接口');
INSERT INTO "public"."sys_api" VALUES (25, '/api/system/api', 'GET', '接口管理', '分页查询接口列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '接口列表');
INSERT INTO "public"."sys_api" VALUES (26, '/api/system/api/:id', 'PUT', '接口管理', '修改接口信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑接口');
INSERT INTO "public"."sys_api" VALUES (27, '/api/system/api/:id', 'DELETE', '接口管理', '删除接口', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除接口');
INSERT INTO "public"."sys_api" VALUES (28, '/api/system/api/all', 'GET', '接口管理', '查询所有接口（分配权限用）', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '全部接口');
INSERT INTO "public"."sys_api" VALUES (29, '/api/system/dict/type', 'POST', '字典管理', '创建字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增字典类型');
INSERT INTO "public"."sys_api" VALUES (30, '/api/system/dict/type', 'GET', '字典管理', '分页查询字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '字典类型列表');
INSERT INTO "public"."sys_api" VALUES (31, '/api/system/dict/type/:id', 'PUT', '字典管理', '修改字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑字典类型');
INSERT INTO "public"."sys_api" VALUES (32, '/api/system/dict/type/:id', 'DELETE', '字典管理', '删除字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除字典类型');
INSERT INTO "public"."sys_api" VALUES (33, '/api/system/dict/data', 'POST', '字典管理', '创建字典数据', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增字典数据');
INSERT INTO "public"."sys_api" VALUES (34, '/api/system/dict/data/:dictType', 'GET', '字典管理', '根据字典类型查询数据列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '字典数据列表');
INSERT INTO "public"."sys_api" VALUES (35, '/api/system/dict/data/:id', 'PUT', '字典管理', '修改字典数据', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑字典数据');
INSERT INTO "public"."sys_api" VALUES (36, '/api/system/dict/data/:id', 'DELETE', '字典管理', '删除字典数据', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除字典数据');
INSERT INTO "public"."sys_api" VALUES (37, '/api/system/file/upload', 'POST', '文件管理', '上传文件到服务器', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '上传文件');
INSERT INTO "public"."sys_api" VALUES (38, '/api/system/file', 'GET', '文件管理', '分页查询文件列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '文件列表');
INSERT INTO "public"."sys_api" VALUES (39, '/api/system/file/:id', 'DELETE', '文件管理', '删除文件', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除文件');
INSERT INTO "public"."sys_api" VALUES (40, '/api/system/log/login', 'GET', '日志管理', '分页查询登录日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '登录日志列表');
INSERT INTO "public"."sys_api" VALUES (41, '/api/system/log/login/clear', 'DELETE', '日志管理', '清除全部登录日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '清空登录日志');
INSERT INTO "public"."sys_api" VALUES (42, '/api/system/log/oper', 'GET', '日志管理', '分页查询操作日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '操作日志列表');
INSERT INTO "public"."sys_api" VALUES (43, '/api/system/log/oper/clear', 'DELETE', '日志管理', '清除全部操作日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '清空操作日志');
INSERT INTO "public"."sys_api" VALUES (44, '/api/nothing', 'GET', 'none', '只更新了名称', '2026-04-30 17:44:00', '2026-05-03 10:40:27', '2026-05-03 10:40:26', '不存在');
INSERT INTO "public"."sys_api" VALUES (45, '/api/test/create', 'POST', 'test', '', '2026-04-30 17:45:06', '2026-05-03 10:41:42', '2026-05-03 10:41:43', '测试接口-创建');
INSERT INTO "public"."sys_api" VALUES (46, '/api/test/:id', 'GET', 'test', '根据ID获取详情', '2026-04-30 17:45:34', '2026-04-30 17:45:34', NULL, '测试接口-详情');
INSERT INTO "public"."sys_api" VALUES (47, '/api/test/:id', 'PUT', 'test', '', '2026-04-30 17:45:54', '2026-04-30 17:45:54', NULL, '测试接口-更新');

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_dict_data";
CREATE TABLE "public"."sys_dict_data" (
  "id" int8 NOT NULL DEFAULT nextval('sys_dict_data_id_seq'::regclass),
  "type_id" int8 NOT NULL,
  "label" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "dict_value" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "sort" int4 NOT NULL DEFAULT 0,
  "status" int2 NOT NULL DEFAULT 1,
  "remark" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_dict_data"."id" IS '字典数据ID';
COMMENT ON COLUMN "public"."sys_dict_data"."type_id" IS '字典类型ID';
COMMENT ON COLUMN "public"."sys_dict_data"."label" IS '字典标签（显示名称）';
COMMENT ON COLUMN "public"."sys_dict_data"."dict_value" IS '字典键值';
COMMENT ON COLUMN "public"."sys_dict_data"."sort" IS '排序值';
COMMENT ON COLUMN "public"."sys_dict_data"."status" IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN "public"."sys_dict_data"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_dict_data"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_dict_data"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sys_dict_data"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_dict_data" IS '字典数据表';

-- ----------------------------
-- Records of sys_dict_data
-- ----------------------------
INSERT INTO "public"."sys_dict_data" VALUES (1, 1, '启用', '1', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (2, 1, '禁用', '0', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (3, 2, '目录', '0', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (4, 2, '菜单', '1', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (5, 2, '按钮', '2', 3, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (6, 3, '显示', '1', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (7, 3, '隐藏', '0', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO "public"."sys_dict_data" VALUES (14, 1, '只改标签', 'test_value', 10, 1, '更新后备注', '2026-05-03 11:03:59', '2026-05-03 16:00:12', '2026-05-03 16:00:12');
INSERT INTO "public"."sys_dict_data" VALUES (15, 6, '简单标签', 'simple', 0, 1, '', '2026-05-03 14:54:03', '2026-05-03 16:07:52', '2026-05-03 16:07:53');
INSERT INTO "public"."sys_dict_data" VALUES (17, 6, '禁用标签', 'disabled_value', 99, 0, '', '2026-05-03 15:12:24', '2026-05-03 16:07:52', '2026-05-03 16:07:53');

-- ----------------------------
-- Table structure for sys_dict_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_dict_type";
CREATE TABLE "public"."sys_dict_type" (
  "id" int8 NOT NULL DEFAULT nextval('sys_dict_type_id_seq'::regclass),
  "name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "code" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "status" int2 NOT NULL DEFAULT 1,
  "remark" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_dict_type"."id" IS '字典类型ID';
COMMENT ON COLUMN "public"."sys_dict_type"."name" IS '字典类型名称';
COMMENT ON COLUMN "public"."sys_dict_type"."code" IS '字典类型编码';
COMMENT ON COLUMN "public"."sys_dict_type"."status" IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN "public"."sys_dict_type"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_dict_type"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_dict_type"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sys_dict_type"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_dict_type" IS '字典类型表';

-- ----------------------------
-- Records of sys_dict_type
-- ----------------------------
INSERT INTO "public"."sys_dict_type" VALUES (1, '系统状态', 'sys_status', 1, '通用状态：启用/禁用', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_dict_type" VALUES (2, '菜单类型', 'sys_menu_type', 1, '菜单类型：目录/菜单/按钮', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_dict_type" VALUES (3, '是否显示', 'sys_visible', 1, '菜单是否显示', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_dict_type" VALUES (6, '更新后字典名2', 'test_dict_type_updated222', 1, '更新后备注2', '2026-05-03 10:53:40', '2026-05-03 16:07:52', '2026-05-03 16:07:53');
INSERT INTO "public"."sys_dict_type" VALUES (7, '最简字典', 'test_simple', 1, '', '2026-05-03 10:56:56', '2026-05-03 16:08:56', '2026-05-03 16:08:56');
INSERT INTO "public"."sys_dict_type" VALUES (8, '不存在', 'not_exist', 1, '创建时即禁用', '2026-05-03 10:57:18', '2026-05-03 15:50:00', NULL);

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_file";
CREATE TABLE "public"."sys_file" (
  "id" int8 NOT NULL DEFAULT nextval('sys_file_id_seq'::regclass),
  "filename" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
  "origin_name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "file_path" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "file_url" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "file_size" int8 NOT NULL DEFAULT 0,
  "file_type" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "uploader_id" int8 NOT NULL DEFAULT 0,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_file"."id" IS '文件ID';
COMMENT ON COLUMN "public"."sys_file"."filename" IS '存储文件名（UUID生成）';
COMMENT ON COLUMN "public"."sys_file"."origin_name" IS '原始文件名';
COMMENT ON COLUMN "public"."sys_file"."file_path" IS '文件存储路径';
COMMENT ON COLUMN "public"."sys_file"."file_url" IS '文件访问URL';
COMMENT ON COLUMN "public"."sys_file"."file_size" IS '文件大小（字节）';
COMMENT ON COLUMN "public"."sys_file"."file_type" IS '文件类型/MIME类型';
COMMENT ON COLUMN "public"."sys_file"."uploader_id" IS '上传人用户ID';
COMMENT ON COLUMN "public"."sys_file"."created_at" IS '上传时间';
COMMENT ON COLUMN "public"."sys_file"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_file" IS '文件记录表';

-- ----------------------------
-- Records of sys_file
-- ----------------------------
INSERT INTO "public"."sys_file" VALUES (1, '63e68410-4338-40e8-8734-e4086a9d0ed4.png', 'code.png', 'D:\GoProject\project_new\blade\backend\gateway\uploads\2026\05\63e68410-4338-40e8-8734-e4086a9d0ed4.png', '/D:/GoProject/project_new/blade/backend/gateway/uploads/2026/05/63e68410-4338-40e8-8734-e4086a9d0ed4.png', 45599, 'image/png', 1, '2026-05-14 10:22:04.926635', NULL);

-- ----------------------------
-- Table structure for sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_login_log";
CREATE TABLE "public"."sys_login_log" (
  "id" int8 NOT NULL DEFAULT nextval('sys_login_log_id_seq'::regclass),
  "user_id" int8 NOT NULL DEFAULT 0,
  "username" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "ip" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "location" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "browser" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "os" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "status" int2 NOT NULL DEFAULT 1,
  "msg" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "login_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_login_log"."id" IS '日志ID';
COMMENT ON COLUMN "public"."sys_login_log"."user_id" IS '用户ID';
COMMENT ON COLUMN "public"."sys_login_log"."username" IS '用户名';
COMMENT ON COLUMN "public"."sys_login_log"."ip" IS '登录IP';
COMMENT ON COLUMN "public"."sys_login_log"."location" IS '登录地点';
COMMENT ON COLUMN "public"."sys_login_log"."browser" IS '浏览器';
COMMENT ON COLUMN "public"."sys_login_log"."os" IS '操作系统';
COMMENT ON COLUMN "public"."sys_login_log"."status" IS '登录状态：1=成功，0=失败';
COMMENT ON COLUMN "public"."sys_login_log"."msg" IS '提示消息';
COMMENT ON COLUMN "public"."sys_login_log"."login_time" IS '登录时间';
COMMENT ON COLUMN "public"."sys_login_log"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_login_log" IS '登录日志表';

-- ----------------------------
-- Records of sys_login_log
-- ----------------------------
INSERT INTO "public"."sys_login_log" VALUES (1, 1, 'admin', '127.0.0.1', '', '其他', '其他', 1, '登录成功', '2026-04-28 22:35:48', NULL);
INSERT INTO "public"."sys_login_log" VALUES (2, 1, 'admin', '127.0.0.1', '', 'Postman', '其他', 1, '登录成功', '2026-04-29 07:36:27', NULL);
INSERT INTO "public"."sys_login_log" VALUES (3, 1, 'admin', '127.0.0.1', '', 'Postman', '其他', 1, '登录成功', '2026-04-29 07:42:16', NULL);
INSERT INTO "public"."sys_login_log" VALUES (4, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-04-29 09:31:44', NULL);
INSERT INTO "public"."sys_login_log" VALUES (5, 1, 'admin', '127.0.0.1', '', 'Postman', '其他', 1, '登录成功', '2026-04-30 07:04:42', NULL);
INSERT INTO "public"."sys_login_log" VALUES (6, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-04-30 09:58:29', NULL);
INSERT INTO "public"."sys_login_log" VALUES (7, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-04-30 17:43:16', NULL);
INSERT INTO "public"."sys_login_log" VALUES (8, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-05-03 10:23:35', NULL);
INSERT INTO "public"."sys_login_log" VALUES (9, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-05-05 09:29:38', NULL);
INSERT INTO "public"."sys_login_log" VALUES (10, 1, 'admin', '[::1]:52934', '', '', '', 1, '登录成功', '2026-05-06 15:07:33', NULL);
INSERT INTO "public"."sys_login_log" VALUES (11, 1, 'admin', '[::1]:52572', '', '', '', 1, '登录成功', '2026-05-06 16:03:31', NULL);
INSERT INTO "public"."sys_login_log" VALUES (12, 1, 'admin', '[::1]:65434', '', '', '', 0, '密码错误', '2026-05-06 16:03:56', NULL);
INSERT INTO "public"."sys_login_log" VALUES (13, 1, 'admin', '[::1]:62171', '', '', '', 0, '密码错误', '2026-05-06 16:06:15', NULL);
INSERT INTO "public"."sys_login_log" VALUES (14, 1, 'admin', '[::1]:64787', '', '', '', 0, '密码错误', '2026-05-06 16:08:19', NULL);
INSERT INTO "public"."sys_login_log" VALUES (15, 1, 'admin', '[::1]:53437', '', '', '', 0, '密码错误', '2026-05-06 16:18:55', NULL);
INSERT INTO "public"."sys_login_log" VALUES (16, 1, 'admin', '[::1]:55671', '', '', '', 0, '密码错误', '2026-05-06 16:21:53', NULL);
INSERT INTO "public"."sys_login_log" VALUES (17, 1, 'admin', '[::1]:65364', '', '', '', 0, '密码错误', '2026-05-06 16:24:43', NULL);
INSERT INTO "public"."sys_login_log" VALUES (18, 1, 'admin', '[::1]:53073', '', '', '', 0, '密码错误', '2026-05-06 16:28:58', NULL);
INSERT INTO "public"."sys_login_log" VALUES (19, 1, 'admin', '[::1]:53406', '', '', '', 0, '密码错误', '2026-05-06 16:34:01', NULL);
INSERT INTO "public"."sys_login_log" VALUES (20, 1, 'admin', '[::1]:61567', '', '', '', 0, '密码错误', '2026-05-06 16:36:17', NULL);
INSERT INTO "public"."sys_login_log" VALUES (21, 1, 'admin', '[::1]:55601', '', '', '', 1, '登录成功', '2026-05-07 10:19:26', NULL);
INSERT INTO "public"."sys_login_log" VALUES (22, 1, 'admin', '[::1]:55601', '', '', '', 0, '密码错误', '2026-05-07 10:19:40', NULL);
INSERT INTO "public"."sys_login_log" VALUES (23, 1, 'admin', '[::1]:54234', '', '', '', 1, '登录成功', '2026-05-07 10:49:05', NULL);
INSERT INTO "public"."sys_login_log" VALUES (24, 1, 'admin', '[::1]:58537', '', '', '', 1, '登录成功', '2026-05-07 10:57:59', NULL);
INSERT INTO "public"."sys_login_log" VALUES (25, 1, 'admin', '[::1]:63271', '', '', '', 1, '登录成功', '2026-05-07 14:52:35', NULL);
INSERT INTO "public"."sys_login_log" VALUES (26, 1, 'admin', '[::1]:55467', '', '', '', 0, '密码错误', '2026-05-07 14:57:59', NULL);
INSERT INTO "public"."sys_login_log" VALUES (27, 1, 'admin', '[::1]:55467', '', '', '', 1, '登录成功', '2026-05-07 14:58:04', NULL);
INSERT INTO "public"."sys_login_log" VALUES (28, 1, 'admin', '[::1]:53366', '', '', '', 1, '登录成功', '2026-05-07 15:45:22', NULL);
INSERT INTO "public"."sys_login_log" VALUES (29, 1, 'admin', '[::1]:64388', '', '', '', 1, '登录成功', '2026-05-07 19:14:51', NULL);
INSERT INTO "public"."sys_login_log" VALUES (30, 1, 'admin', '[::1]:52524', '', '', '', 1, '登录成功', '2026-05-08 09:47:48', NULL);
INSERT INTO "public"."sys_login_log" VALUES (31, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-05-08 10:16:55', NULL);
INSERT INTO "public"."sys_login_log" VALUES (32, 1, 'admin', '[::1]:53513', '', '', '', 1, '登录成功', '2026-05-08 11:14:33', NULL);
INSERT INTO "public"."sys_login_log" VALUES (33, 1, 'admin', '[::1]:50290', '', '', '', 1, '登录成功', '2026-05-09 10:25:20.287368', NULL);
INSERT INTO "public"."sys_login_log" VALUES (34, 1, 'admin', '127.0.0.1:52788', '', '', '', 0, '密码错误', '2026-05-09 17:26:07.985826', NULL);
INSERT INTO "public"."sys_login_log" VALUES (35, 1, 'admin', '127.0.0.1:51128', '', '', '', 0, '密码错误', '2026-05-09 17:26:11.260637', NULL);
INSERT INTO "public"."sys_login_log" VALUES (36, 1, 'admin', '127.0.0.1:51758', '', '', '', 0, '密码错误', '2026-05-09 17:26:30.885793', NULL);
INSERT INTO "public"."sys_login_log" VALUES (37, 1, 'admin', '127.0.0.1:55853', '', '', '', 1, '登录成功', '2026-05-09 17:26:38.905005', NULL);
INSERT INTO "public"."sys_login_log" VALUES (38, 1, 'admin', '127.0.0.1:49672', '', '', '', 0, '密码错误', '2026-05-09 17:27:43.480278', NULL);
INSERT INTO "public"."sys_login_log" VALUES (39, 1, 'admin', '127.0.0.1:61982', '', '', '', 1, '登录成功', '2026-05-09 17:27:50.269777', NULL);
INSERT INTO "public"."sys_login_log" VALUES (40, 1, 'admin', '127.0.0.1:49870', '', '', '', 0, '密码错误', '2026-05-09 19:14:17.016531', NULL);
INSERT INTO "public"."sys_login_log" VALUES (41, 1, 'admin', '127.0.0.1:58367', '', '', '', 1, '登录成功', '2026-05-09 19:14:35.024186', NULL);
INSERT INTO "public"."sys_login_log" VALUES (42, 1, 'admin', '127.0.0.1:56665', '', '', '', 1, '登录成功', '2026-05-09 19:40:40.89324', NULL);
INSERT INTO "public"."sys_login_log" VALUES (43, 1, 'admin', '127.0.0.1:54580', '', '', '', 1, '登录成功', '2026-05-09 19:42:16.5761', NULL);
INSERT INTO "public"."sys_login_log" VALUES (44, 1, 'admin', '127.0.0.1:61638', '', '', '', 1, '登录成功', '2026-05-09 19:52:48.623521', NULL);
INSERT INTO "public"."sys_login_log" VALUES (45, 1, 'admin', '127.0.0.1:49699', '', '', '', 1, '登录成功', '2026-05-11 09:10:34.304927', NULL);
INSERT INTO "public"."sys_login_log" VALUES (46, 1, 'admin', '127.0.0.1:61461', '', '', '', 1, '登录成功', '2026-05-13 09:02:02.983039', NULL);
INSERT INTO "public"."sys_login_log" VALUES (47, 1, 'admin', '127.0.0.1:56826', '', '', '', 1, '登录成功', '2026-05-14 10:21:42.907902', NULL);
INSERT INTO "public"."sys_login_log" VALUES (48, 1, 'admin', '127.0.0.1:53354', '', '', '', 0, '密码错误', '2026-05-15 17:15:09.945006', NULL);
INSERT INTO "public"."sys_login_log" VALUES (49, 1, 'admin', '127.0.0.1:62723', '', '', '', 1, '登录成功', '2026-05-15 17:15:14.428377', NULL);

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_menu";
CREATE TABLE "public"."sys_menu" (
  "id" int8 NOT NULL DEFAULT nextval('sys_menu_id_seq'::regclass),
  "parent_id" int8 NOT NULL DEFAULT 0,
  "name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "menu_path" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "component" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "icon" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "menu_type" int2 NOT NULL DEFAULT 0,
  "permission" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "sort" int4 NOT NULL DEFAULT 0,
  "visible" int2 NOT NULL DEFAULT 1,
  "status" int2 NOT NULL DEFAULT 1,
  "remark" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_menu"."id" IS '菜单ID';
COMMENT ON COLUMN "public"."sys_menu"."parent_id" IS '父菜单ID，0表示顶级';
COMMENT ON COLUMN "public"."sys_menu"."name" IS '菜单名称';
COMMENT ON COLUMN "public"."sys_menu"."menu_path" IS '路由路径';
COMMENT ON COLUMN "public"."sys_menu"."component" IS '前端组件路径';
COMMENT ON COLUMN "public"."sys_menu"."icon" IS '菜单图标';
COMMENT ON COLUMN "public"."sys_menu"."menu_type" IS '菜单类型：0=目录，1=菜单，2=按钮';
COMMENT ON COLUMN "public"."sys_menu"."permission" IS '权限标识符，如：system:user:list';
COMMENT ON COLUMN "public"."sys_menu"."sort" IS '排序值，越小越靠前';
COMMENT ON COLUMN "public"."sys_menu"."visible" IS '是否可见：1=可见，0=隐藏';
COMMENT ON COLUMN "public"."sys_menu"."status" IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN "public"."sys_menu"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_menu"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_menu"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sys_menu"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_menu" IS '系统菜单表';

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO "public"."sys_menu" VALUES (1, 0, '系统管理', '/system', 'Layout', 'setting', 0, '', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (2, 1, '用户管理', '/system/user', 'system/user/index', 'user', 1, 'system:user:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (3, 1, '角色管理', '/system/role', 'system/role/index', 'peoples', 1, 'system:role:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (4, 1, '菜单管理', '/system/menu', 'system/menu/index', 'tree-table', 1, 'system:menu:list', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (5, 1, '接口管理', '/system/api', 'system/api/index', 'api', 1, 'system:api:list', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (6, 1, '字典管理', '/system/dict', 'system/dict/index', 'dict', 1, 'system:dict:list', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (7, 0, '日志管理', '/log', 'Layout', 'log', 0, '', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (8, 7, '登录日志', '/log/login', 'log/login/index', 'logininfor', 1, 'system:loginlog:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (9, 7, '操作日志', '/log/oper', 'log/oper/index', 'form', 1, 'system:operlog:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (20, 2, '新增用户', '', '', '', 2, 'system:user:create', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (21, 2, '编辑用户', '', '', '', 2, 'system:user:update', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (22, 2, '删除用户', '', '', '', 2, 'system:user:delete', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (23, 3, '新增角色', '', '', '', 2, 'system:role:create', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (24, 3, '编辑角色', '', '', '', 2, 'system:role:update', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (25, 3, '删除角色', '', '', '', 2, 'system:role:delete', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (26, 3, '分配菜单', '', '', '', 2, 'system:role:assignmenu', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (27, 3, '分配接口', '', '', '', 2, 'system:role:assignapi', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_menu" VALUES (33, 2, '重置密码', '', '', '', 2, 'system:user:resetpwd', 4, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (34, 4, '新增菜单', '', '', '', 2, 'system:menu:create', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (35, 4, '编辑菜单', '', '', '', 2, 'system:menu:update', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (36, 4, '删除菜单', '', '', '', 2, 'system:menu:delete', 3, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (37, 5, '新增接口', '', '', '', 2, 'system:api:create', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (38, 5, '编辑接口', '', '', '', 2, 'system:api:update', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (39, 5, '删除接口', '', '', '', 2, 'system:api:delete', 3, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (40, 6, '新增字典类型', '', '', '', 2, 'system:dict:type:create', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (41, 6, '编辑字典类型', '', '', '', 2, 'system:dict:type:update', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (42, 6, '删除字典类型', '', '', '', 2, 'system:dict:type:delete', 3, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (43, 6, '新增字典数据', '', '', '', 2, 'system:dict:data:create', 4, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (44, 6, '编辑字典数据', '', '', '', 2, 'system:dict:data:update', 5, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (45, 6, '删除字典数据', '', '', '', 2, 'system:dict:data:delete', 6, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (46, 8, '清空日志', '', '', '', 2, 'system:loginlog:clear', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (47, 9, '清空日志', '', '', '', 2, 'system:operlog:clear', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (48, 1, '文件管理', '/system/file', 'system/file/index', 'upload', 1, 'system:file:list', 6, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (49, 48, '上传文件', '', '', '', 2, 'system:file:upload', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO "public"."sys_menu" VALUES (50, 48, '删除文件', '', '', '', 2, 'system:file:delete', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);

-- ----------------------------
-- Table structure for sys_oper_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_oper_log";
CREATE TABLE "public"."sys_oper_log" (
  "id" int8 NOT NULL DEFAULT nextval('sys_oper_log_id_seq'::regclass),
  "title" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "business_type" int4 NOT NULL DEFAULT 0,
  "method" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "request_method" varchar(16) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "operator_type" int4 NOT NULL DEFAULT 0,
  "operator_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "operator_id" int8 NOT NULL DEFAULT 0,
  "dept_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "oper_url" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "oper_ip" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "oper_param" text COLLATE "pg_catalog"."default",
  "json_result" text COLLATE "pg_catalog"."default",
  "status" int2 NOT NULL DEFAULT 1,
  "error_msg" varchar(2000) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "oper_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_oper_log"."id" IS '操作日志ID';
COMMENT ON COLUMN "public"."sys_oper_log"."title" IS '操作模块标题';
COMMENT ON COLUMN "public"."sys_oper_log"."business_type" IS '业务类型：0=其他，1=新增，2=修改，3=删除，4=查询';
COMMENT ON COLUMN "public"."sys_oper_log"."method" IS '方法名称';
COMMENT ON COLUMN "public"."sys_oper_log"."request_method" IS '请求方式：GET/POST/PUT/DELETE';
COMMENT ON COLUMN "public"."sys_oper_log"."operator_type" IS '操作人类型：0=其他，1=后台用户，2=手机端用户';
COMMENT ON COLUMN "public"."sys_oper_log"."operator_name" IS '操作人员名称';
COMMENT ON COLUMN "public"."sys_oper_log"."operator_id" IS '操作人员ID';
COMMENT ON COLUMN "public"."sys_oper_log"."dept_name" IS '部门名称';
COMMENT ON COLUMN "public"."sys_oper_log"."oper_url" IS '请求URL';
COMMENT ON COLUMN "public"."sys_oper_log"."oper_ip" IS '操作IP';
COMMENT ON COLUMN "public"."sys_oper_log"."oper_param" IS '请求参数（JSON）';
COMMENT ON COLUMN "public"."sys_oper_log"."json_result" IS '返回结果（JSON）';
COMMENT ON COLUMN "public"."sys_oper_log"."status" IS '操作状态：1=成功，0=失败';
COMMENT ON COLUMN "public"."sys_oper_log"."error_msg" IS '错误消息';
COMMENT ON COLUMN "public"."sys_oper_log"."oper_time" IS '操作时间';
COMMENT ON COLUMN "public"."sys_oper_log"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_oper_log" IS '操作日志表';

-- ----------------------------
-- Records of sys_oper_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_role";
CREATE TABLE "public"."sys_role" (
  "id" int8 NOT NULL DEFAULT nextval('sys_role_id_seq'::regclass),
  "name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "code" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "status" int2 NOT NULL DEFAULT 1,
  "remark" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "sort" int4 NOT NULL DEFAULT 0,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_role"."id" IS '角色ID';
COMMENT ON COLUMN "public"."sys_role"."name" IS '角色名称';
COMMENT ON COLUMN "public"."sys_role"."code" IS '角色编码（casbin中使用）';
COMMENT ON COLUMN "public"."sys_role"."status" IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN "public"."sys_role"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_role"."sort" IS '排序值，越小越靠前';
COMMENT ON COLUMN "public"."sys_role"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_role"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sys_role"."deleted_at" IS '软删除时间, 空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_role" IS '系统角色表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO "public"."sys_role" VALUES (1, '超级管理员', 'admin', 1, '系统内置超级管理员角色，拥有所有权限', 1, '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO "public"."sys_role" VALUES (6, '管理员', 'manager', 1, '能管人、管角色、管接口、但不碰业务数据', 2, '2026-04-29 11:32:53', '2026-04-29 11:32:53', NULL);
INSERT INTO "public"."sys_role" VALUES (7, '操作者', 'operator', 1, '能录入业务数据、查看状态，不能管理系统', 3, '2026-04-29 14:35:57', '2026-04-29 16:24:47', NULL);
INSERT INTO "public"."sys_role" VALUES (8, '查看者', 'viewer', 1, '只能看，什么都不能改', 4, '2026-04-29 14:36:36', '2026-04-29 17:28:48', NULL);

-- ----------------------------
-- Table structure for sys_role_api
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_role_api";
CREATE TABLE "public"."sys_role_api" (
  "id" int8 NOT NULL DEFAULT nextval('sys_role_api_id_seq'::regclass),
  "role_id" int8 NOT NULL,
  "api_id" int8 NOT NULL
)
;
COMMENT ON COLUMN "public"."sys_role_api"."id" IS '主键ID';
COMMENT ON COLUMN "public"."sys_role_api"."role_id" IS '角色ID';
COMMENT ON COLUMN "public"."sys_role_api"."api_id" IS '接口ID';
COMMENT ON TABLE "public"."sys_role_api" IS '角色接口关联表';

-- ----------------------------
-- Records of sys_role_api
-- ----------------------------
INSERT INTO "public"."sys_role_api" VALUES (8, 1, 1);
INSERT INTO "public"."sys_role_api" VALUES (9, 1, 2);
INSERT INTO "public"."sys_role_api" VALUES (10, 1, 3);
INSERT INTO "public"."sys_role_api" VALUES (11, 1, 4);
INSERT INTO "public"."sys_role_api" VALUES (12, 1, 5);
INSERT INTO "public"."sys_role_api" VALUES (13, 1, 6);
INSERT INTO "public"."sys_role_api" VALUES (14, 1, 7);
INSERT INTO "public"."sys_role_api" VALUES (15, 1, 8);
INSERT INTO "public"."sys_role_api" VALUES (16, 1, 9);
INSERT INTO "public"."sys_role_api" VALUES (17, 1, 10);
INSERT INTO "public"."sys_role_api" VALUES (18, 1, 11);
INSERT INTO "public"."sys_role_api" VALUES (19, 1, 12);
INSERT INTO "public"."sys_role_api" VALUES (20, 1, 13);
INSERT INTO "public"."sys_role_api" VALUES (21, 1, 14);
INSERT INTO "public"."sys_role_api" VALUES (22, 1, 15);
INSERT INTO "public"."sys_role_api" VALUES (23, 1, 16);
INSERT INTO "public"."sys_role_api" VALUES (24, 1, 17);
INSERT INTO "public"."sys_role_api" VALUES (25, 1, 18);
INSERT INTO "public"."sys_role_api" VALUES (26, 1, 19);
INSERT INTO "public"."sys_role_api" VALUES (27, 1, 20);
INSERT INTO "public"."sys_role_api" VALUES (28, 1, 21);
INSERT INTO "public"."sys_role_api" VALUES (29, 1, 22);
INSERT INTO "public"."sys_role_api" VALUES (30, 1, 23);
INSERT INTO "public"."sys_role_api" VALUES (31, 1, 24);
INSERT INTO "public"."sys_role_api" VALUES (32, 1, 25);
INSERT INTO "public"."sys_role_api" VALUES (33, 1, 26);
INSERT INTO "public"."sys_role_api" VALUES (34, 1, 27);
INSERT INTO "public"."sys_role_api" VALUES (35, 1, 28);
INSERT INTO "public"."sys_role_api" VALUES (36, 1, 29);
INSERT INTO "public"."sys_role_api" VALUES (37, 1, 30);
INSERT INTO "public"."sys_role_api" VALUES (38, 1, 31);
INSERT INTO "public"."sys_role_api" VALUES (39, 1, 32);
INSERT INTO "public"."sys_role_api" VALUES (40, 1, 33);
INSERT INTO "public"."sys_role_api" VALUES (41, 1, 34);
INSERT INTO "public"."sys_role_api" VALUES (42, 1, 35);
INSERT INTO "public"."sys_role_api" VALUES (43, 1, 36);
INSERT INTO "public"."sys_role_api" VALUES (44, 1, 37);
INSERT INTO "public"."sys_role_api" VALUES (45, 1, 38);
INSERT INTO "public"."sys_role_api" VALUES (46, 1, 39);
INSERT INTO "public"."sys_role_api" VALUES (47, 1, 40);
INSERT INTO "public"."sys_role_api" VALUES (48, 1, 41);
INSERT INTO "public"."sys_role_api" VALUES (49, 1, 42);
INSERT INTO "public"."sys_role_api" VALUES (50, 1, 43);
INSERT INTO "public"."sys_role_api" VALUES (51, 1, 46);
INSERT INTO "public"."sys_role_api" VALUES (52, 1, 47);
INSERT INTO "public"."sys_role_api" VALUES (70, 6, 1);
INSERT INTO "public"."sys_role_api" VALUES (71, 6, 2);
INSERT INTO "public"."sys_role_api" VALUES (72, 6, 3);
INSERT INTO "public"."sys_role_api" VALUES (73, 6, 4);
INSERT INTO "public"."sys_role_api" VALUES (74, 6, 6);
INSERT INTO "public"."sys_role_api" VALUES (75, 6, 5);
INSERT INTO "public"."sys_role_api" VALUES (76, 6, 7);
INSERT INTO "public"."sys_role_api" VALUES (77, 6, 10);
INSERT INTO "public"."sys_role_api" VALUES (78, 6, 8);
INSERT INTO "public"."sys_role_api" VALUES (79, 6, 11);
INSERT INTO "public"."sys_role_api" VALUES (80, 6, 12);
INSERT INTO "public"."sys_role_api" VALUES (81, 6, 13);
INSERT INTO "public"."sys_role_api" VALUES (82, 6, 14);
INSERT INTO "public"."sys_role_api" VALUES (83, 6, 16);
INSERT INTO "public"."sys_role_api" VALUES (84, 6, 17);
INSERT INTO "public"."sys_role_api" VALUES (85, 6, 18);
INSERT INTO "public"."sys_role_api" VALUES (86, 6, 19);
INSERT INTO "public"."sys_role_api" VALUES (87, 6, 20);

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_role_menu";
CREATE TABLE "public"."sys_role_menu" (
  "id" int8 NOT NULL DEFAULT nextval('sys_role_menu_id_seq'::regclass),
  "role_id" int8 NOT NULL,
  "menu_id" int8 NOT NULL
)
;
COMMENT ON COLUMN "public"."sys_role_menu"."id" IS '主键ID';
COMMENT ON COLUMN "public"."sys_role_menu"."role_id" IS '角色ID';
COMMENT ON COLUMN "public"."sys_role_menu"."menu_id" IS '菜单ID';
COMMENT ON TABLE "public"."sys_role_menu" IS '角色菜单关联表';

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO "public"."sys_role_menu" VALUES (1, 1, 1);
INSERT INTO "public"."sys_role_menu" VALUES (2, 1, 2);
INSERT INTO "public"."sys_role_menu" VALUES (39, 8, 1);
INSERT INTO "public"."sys_role_menu" VALUES (42, 8, 2);
INSERT INTO "public"."sys_role_menu" VALUES (40, 8, 6);
INSERT INTO "public"."sys_role_menu" VALUES (41, 8, 7);
INSERT INTO "public"."sys_role_menu" VALUES (44, 8, 8);
INSERT INTO "public"."sys_role_menu" VALUES (43, 8, 20);
INSERT INTO "public"."sys_role_menu" VALUES (38, 8, 42);
INSERT INTO "public"."sys_role_menu" VALUES (45, 1, 3);
INSERT INTO "public"."sys_role_menu" VALUES (46, 1, 4);
INSERT INTO "public"."sys_role_menu" VALUES (47, 1, 5);
INSERT INTO "public"."sys_role_menu" VALUES (48, 1, 6);
INSERT INTO "public"."sys_role_menu" VALUES (49, 1, 7);
INSERT INTO "public"."sys_role_menu" VALUES (50, 1, 8);
INSERT INTO "public"."sys_role_menu" VALUES (51, 1, 9);
INSERT INTO "public"."sys_role_menu" VALUES (52, 1, 20);
INSERT INTO "public"."sys_role_menu" VALUES (53, 1, 21);
INSERT INTO "public"."sys_role_menu" VALUES (54, 1, 22);
INSERT INTO "public"."sys_role_menu" VALUES (55, 1, 23);
INSERT INTO "public"."sys_role_menu" VALUES (56, 1, 24);
INSERT INTO "public"."sys_role_menu" VALUES (57, 1, 25);
INSERT INTO "public"."sys_role_menu" VALUES (58, 1, 26);
INSERT INTO "public"."sys_role_menu" VALUES (59, 1, 27);
INSERT INTO "public"."sys_role_menu" VALUES (60, 1, 33);
INSERT INTO "public"."sys_role_menu" VALUES (61, 1, 34);
INSERT INTO "public"."sys_role_menu" VALUES (62, 1, 35);
INSERT INTO "public"."sys_role_menu" VALUES (63, 1, 36);
INSERT INTO "public"."sys_role_menu" VALUES (64, 1, 37);
INSERT INTO "public"."sys_role_menu" VALUES (65, 1, 38);
INSERT INTO "public"."sys_role_menu" VALUES (66, 1, 39);
INSERT INTO "public"."sys_role_menu" VALUES (67, 1, 40);
INSERT INTO "public"."sys_role_menu" VALUES (68, 1, 41);
INSERT INTO "public"."sys_role_menu" VALUES (69, 1, 42);
INSERT INTO "public"."sys_role_menu" VALUES (70, 1, 43);
INSERT INTO "public"."sys_role_menu" VALUES (71, 1, 44);
INSERT INTO "public"."sys_role_menu" VALUES (72, 1, 45);
INSERT INTO "public"."sys_role_menu" VALUES (73, 1, 46);
INSERT INTO "public"."sys_role_menu" VALUES (74, 1, 47);
INSERT INTO "public"."sys_role_menu" VALUES (75, 1, 48);
INSERT INTO "public"."sys_role_menu" VALUES (76, 1, 49);
INSERT INTO "public"."sys_role_menu" VALUES (77, 1, 50);
INSERT INTO "public"."sys_role_menu" VALUES (96, 6, 5);
INSERT INTO "public"."sys_role_menu" VALUES (97, 6, 21);
INSERT INTO "public"."sys_role_menu" VALUES (98, 6, 3);
INSERT INTO "public"."sys_role_menu" VALUES (99, 6, 26);
INSERT INTO "public"."sys_role_menu" VALUES (100, 6, 6);
INSERT INTO "public"."sys_role_menu" VALUES (101, 6, 48);
INSERT INTO "public"."sys_role_menu" VALUES (102, 6, 8);
INSERT INTO "public"."sys_role_menu" VALUES (103, 6, 9);
INSERT INTO "public"."sys_role_menu" VALUES (104, 6, 2);
INSERT INTO "public"."sys_role_menu" VALUES (105, 6, 20);
INSERT INTO "public"."sys_role_menu" VALUES (106, 6, 24);
INSERT INTO "public"."sys_role_menu" VALUES (107, 6, 7);
INSERT INTO "public"."sys_role_menu" VALUES (108, 6, 1);
INSERT INTO "public"."sys_role_menu" VALUES (109, 6, 22);
INSERT INTO "public"."sys_role_menu" VALUES (110, 6, 23);
INSERT INTO "public"."sys_role_menu" VALUES (111, 6, 25);
INSERT INTO "public"."sys_role_menu" VALUES (112, 6, 27);
INSERT INTO "public"."sys_role_menu" VALUES (113, 6, 4);

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_user";
CREATE TABLE "public"."sys_user" (
  "id" int8 NOT NULL DEFAULT nextval('sys_user_id_seq'::regclass),
  "username" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "nickname" varchar(64) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "email" varchar(128) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "phone" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "avatar" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "status" int2 NOT NULL DEFAULT 1,
  "remark" varchar(512) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "created_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sys_user"."id" IS '用户ID';
COMMENT ON COLUMN "public"."sys_user"."username" IS '登录用户名';
COMMENT ON COLUMN "public"."sys_user"."password" IS '登录密码（bcrypt加密）';
COMMENT ON COLUMN "public"."sys_user"."nickname" IS '昵称/显示名';
COMMENT ON COLUMN "public"."sys_user"."email" IS '邮箱';
COMMENT ON COLUMN "public"."sys_user"."phone" IS '手机号';
COMMENT ON COLUMN "public"."sys_user"."avatar" IS '头像URL';
COMMENT ON COLUMN "public"."sys_user"."status" IS '状态：1=启用，0=禁用';
COMMENT ON COLUMN "public"."sys_user"."remark" IS '备注';
COMMENT ON COLUMN "public"."sys_user"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sys_user"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sys_user"."deleted_at" IS '软删除时间,空-未删除，非空为已删除时间';
COMMENT ON TABLE "public"."sys_user" IS '系统用户表';

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO "public"."sys_user" VALUES (1, 'admin', '$2a$10$vh6VEcx58MphNdTaShd...ErM.V1CkWwYOt5wo47qy0soZ6SKZqKK', '系统管理员', 'admin@example.com', '18288888888', '', 1, '系统内置超级管理员账号', '2026-04-07 15:28:48', '2026-05-07 14:57:52', NULL);
INSERT INTO "public"."sys_user" VALUES (2, 'testuser01', '$2a$10$1ZXON13uBALtKlyzAB7UOeoS0UaQnrqNsoauedo5Q4dFTy1tbZbAu', '只改昵称', 'updated@example.com', '13900139001', 'https://example.com/new-avatar.png', 1, '更新备注', '2026-04-30 07:05:13', '2026-05-05 11:40:47', NULL);
INSERT INTO "public"."sys_user" VALUES (3, 'testuser02', '$2a$10$pRRoEYtW/fiFGwRvVXtneu3SHKuxaDwMkVoof0HPTIHvE.WQ53gkK', '测试用户02', 'test01@example.com', '13800138002', 'https://example.com/avatar.png', 1, '测试创建的用户2', '2026-04-30 07:12:36', '2026-04-30 07:12:36', NULL);
INSERT INTO "public"."sys_user" VALUES (4, 'testuser03', '$2a$10$m1cUikaSJFBVHMfUa4XkX.leDwvgDpIb86KaMhu2MbILx13EjOifq', '测试用户03', 'test01@example.com', '13800138003', 'https://example.com/avatar.png', 1, '测试创建的用户2', '2026-04-30 07:20:06', '2026-05-05 13:24:32', '2026-05-05 13:24:34');
INSERT INTO "public"."sys_user" VALUES (9, 'tianyx', '$2a$10$gFiCU3xpbadRwvsoPCXz2.i3U7AEHKJJblE67JewjbCNuEEBVuf/C', '管理员01', '1787817883@qq.com', '1888888888', 'https://example.com/avatar.png', 1, '', '2026-05-05 09:56:55', '2026-05-05 09:56:55', NULL);

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."sys_user_role";
CREATE TABLE "public"."sys_user_role" (
  "id" int8 NOT NULL DEFAULT nextval('sys_user_role_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "role_id" int8 NOT NULL
)
;
COMMENT ON COLUMN "public"."sys_user_role"."id" IS '主键ID';
COMMENT ON COLUMN "public"."sys_user_role"."user_id" IS '用户ID';
COMMENT ON COLUMN "public"."sys_user_role"."role_id" IS '角色ID';
COMMENT ON TABLE "public"."sys_user_role" IS '用户角色关联表';

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO "public"."sys_user_role" VALUES (3, 1, 1);
INSERT INTO "public"."sys_user_role" VALUES (4, 1, 8);
INSERT INTO "public"."sys_user_role" VALUES (8, 9, 6);
INSERT INTO "public"."sys_user_role" VALUES (9, 9, 7);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."casbin_rule_id_seq"
OWNED BY "public"."casbin_rule"."id";
SELECT setval('"public"."casbin_rule_id_seq"', 36, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_api_id_seq"
OWNED BY "public"."sys_api"."id";
SELECT setval('"public"."sys_api_id_seq"', 48, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_dict_data_id_seq"
OWNED BY "public"."sys_dict_data"."id";
SELECT setval('"public"."sys_dict_data_id_seq"', 18, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_dict_type_id_seq"
OWNED BY "public"."sys_dict_type"."id";
SELECT setval('"public"."sys_dict_type_id_seq"', 9, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_file_id_seq"
OWNED BY "public"."sys_file"."id";
SELECT setval('"public"."sys_file_id_seq"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_login_log_id_seq"
OWNED BY "public"."sys_login_log"."id";
SELECT setval('"public"."sys_login_log_id_seq"', 49, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_menu_id_seq"
OWNED BY "public"."sys_menu"."id";
SELECT setval('"public"."sys_menu_id_seq"', 51, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_oper_log_id_seq"
OWNED BY "public"."sys_oper_log"."id";
SELECT setval('"public"."sys_oper_log_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_role_api_id_seq"
OWNED BY "public"."sys_role_api"."id";
SELECT setval('"public"."sys_role_api_id_seq"', 87, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_role_id_seq"
OWNED BY "public"."sys_role"."id";
SELECT setval('"public"."sys_role_id_seq"', 9, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_role_menu_id_seq"
OWNED BY "public"."sys_role_menu"."id";
SELECT setval('"public"."sys_role_menu_id_seq"', 113, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_user_id_seq"
OWNED BY "public"."sys_user"."id";
SELECT setval('"public"."sys_user_id_seq"', 10, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sys_user_role_id_seq"
OWNED BY "public"."sys_user_role"."id";
SELECT setval('"public"."sys_user_role_id_seq"', 9, true);

-- ----------------------------
-- Indexes structure for table casbin_rule
-- ----------------------------
CREATE UNIQUE INDEX "idx_casbin_rule" ON "public"."casbin_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table casbin_rule
-- ----------------------------
ALTER TABLE "public"."casbin_rule" ADD CONSTRAINT "casbin_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_api
-- ----------------------------
CREATE INDEX "sys_api_idx_group" ON "public"."sys_api" USING btree (
  "api_group" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sys_api_uk_path_method" ON "public"."sys_api" USING btree (
  "api_path" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "method" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_api
-- ----------------------------
ALTER TABLE "public"."sys_api" ADD CONSTRAINT "sys_api_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_dict_data
-- ----------------------------
CREATE INDEX "sys_dict_data_idx_status" ON "public"."sys_dict_data" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE INDEX "sys_dict_data_idx_type_id" ON "public"."sys_dict_data" USING btree (
  "type_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_dict_data
-- ----------------------------
ALTER TABLE "public"."sys_dict_data" ADD CONSTRAINT "sys_dict_data_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_dict_type
-- ----------------------------
CREATE UNIQUE INDEX "sys_dict_type_uk_code" ON "public"."sys_dict_type" USING btree (
  "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_dict_type
-- ----------------------------
ALTER TABLE "public"."sys_dict_type" ADD CONSTRAINT "sys_dict_type_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_file
-- ----------------------------
CREATE INDEX "sys_file_idx_created_at" ON "public"."sys_file" USING btree (
  "created_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "sys_file_idx_uploader_id" ON "public"."sys_file" USING btree (
  "uploader_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_file
-- ----------------------------
ALTER TABLE "public"."sys_file" ADD CONSTRAINT "sys_file_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_login_log
-- ----------------------------
CREATE INDEX "sys_login_log_idx_login_time" ON "public"."sys_login_log" USING btree (
  "login_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "sys_login_log_idx_user_id" ON "public"."sys_login_log" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "sys_login_log_idx_username" ON "public"."sys_login_log" USING btree (
  "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_login_log
-- ----------------------------
ALTER TABLE "public"."sys_login_log" ADD CONSTRAINT "sys_login_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_menu
-- ----------------------------
CREATE INDEX "sys_menu_idx_parent_id" ON "public"."sys_menu" USING btree (
  "parent_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "sys_menu_idx_type" ON "public"."sys_menu" USING btree (
  "menu_type" "pg_catalog"."int2_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_menu
-- ----------------------------
ALTER TABLE "public"."sys_menu" ADD CONSTRAINT "sys_menu_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_oper_log
-- ----------------------------
CREATE INDEX "sys_oper_log_idx_oper_time" ON "public"."sys_oper_log" USING btree (
  "oper_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "sys_oper_log_idx_operator_id" ON "public"."sys_oper_log" USING btree (
  "operator_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "sys_oper_log_idx_status" ON "public"."sys_oper_log" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_oper_log
-- ----------------------------
ALTER TABLE "public"."sys_oper_log" ADD CONSTRAINT "sys_oper_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_role
-- ----------------------------
CREATE INDEX "sys_role_idx_deleted_at" ON "public"."sys_role" USING btree (
  "deleted_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "sys_role_idx_status" ON "public"."sys_role" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sys_role_uk_code" ON "public"."sys_role" USING btree (
  "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_role
-- ----------------------------
ALTER TABLE "public"."sys_role" ADD CONSTRAINT "sys_role_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_role_api
-- ----------------------------
CREATE INDEX "sys_role_api_idx_api_id" ON "public"."sys_role_api" USING btree (
  "api_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sys_role_api_uk_role_api" ON "public"."sys_role_api" USING btree (
  "role_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "api_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_role_api
-- ----------------------------
ALTER TABLE "public"."sys_role_api" ADD CONSTRAINT "sys_role_api_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_role_menu
-- ----------------------------
CREATE INDEX "sys_role_menu_idx_menu_id" ON "public"."sys_role_menu" USING btree (
  "menu_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sys_role_menu_uk_role_menu" ON "public"."sys_role_menu" USING btree (
  "role_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "menu_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_role_menu
-- ----------------------------
ALTER TABLE "public"."sys_role_menu" ADD CONSTRAINT "sys_role_menu_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_user
-- ----------------------------
CREATE INDEX "sys_user_idx_deleted_at" ON "public"."sys_user" USING btree (
  "deleted_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "sys_user_idx_status" ON "public"."sys_user" USING btree (
  "status" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sys_user_uk_username" ON "public"."sys_user" USING btree (
  "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_user
-- ----------------------------
ALTER TABLE "public"."sys_user" ADD CONSTRAINT "sys_user_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sys_user_role
-- ----------------------------
CREATE INDEX "sys_user_role_idx_role_id" ON "public"."sys_user_role" USING btree (
  "role_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "sys_user_role_uk_user_role" ON "public"."sys_user_role" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "role_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sys_user_role
-- ----------------------------
ALTER TABLE "public"."sys_user_role" ADD CONSTRAINT "sys_user_role_pkey" PRIMARY KEY ("id");
