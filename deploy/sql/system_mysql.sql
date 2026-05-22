-- PostgreSQL 14 compatible schema and data converted from the original MySQL dump.
SET client_encoding = 'UTF8';

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS casbin_rule;
CREATE TABLE casbin_rule (
  id BIGSERIAL NOT NULL,
  ptype VARCHAR(100) NULL DEFAULT NULL,
  v0 VARCHAR(100) NULL DEFAULT NULL,
  v1 VARCHAR(100) NULL DEFAULT NULL,
  v2 VARCHAR(100) NULL DEFAULT NULL,
  v3 VARCHAR(100) NULL DEFAULT NULL,
  v4 VARCHAR(100) NULL DEFAULT NULL,
  v5 VARCHAR(100) NULL DEFAULT NULL,
  PRIMARY KEY (id)
);

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO casbin_rule VALUES (1, 'p', 'admin', '/api/*', '*', NULL, NULL, NULL);

-- ----------------------------
-- Table structure for sys_api
-- ----------------------------
DROP TABLE IF EXISTS sys_api;
CREATE TABLE sys_api (
  id BIGSERIAL NOT NULL,
  api_path VARCHAR(256) NOT NULL,
  method VARCHAR(16) NOT NULL,
  api_group VARCHAR(64) NOT NULL DEFAULT '',
  description VARCHAR(256) NULL DEFAULT '',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  api_name VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_api_uk_path_method ON sys_api (api_path, method);
CREATE INDEX sys_api_idx_group ON sys_api (api_group);

-- ----------------------------
-- Records of sys_api
-- ----------------------------
INSERT INTO sys_api VALUES (1, '/api/auth/login', 'POST', '认证管理', '用户登录获取Token', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '用户登录');
INSERT INTO sys_api VALUES (2, '/api/auth/refresh', 'POST', '认证管理', '刷新访问令牌', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '刷新令牌');
INSERT INTO sys_api VALUES (3, '/api/auth/changePassword', 'POST', '认证管理', '登录用户修改密码', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '修改密码');
INSERT INTO sys_api VALUES (4, '/api/auth/logout', 'POST', '认证管理', '退出登录并注销Token', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '退出登录');
INSERT INTO sys_api VALUES (5, '/api/auth/userInfo', 'GET', '认证管理', '获取当前登录用户信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '获取用户信息');
INSERT INTO sys_api VALUES (6, '/api/system/user', 'POST', '用户管理', '创建系统用户', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增用户');
INSERT INTO sys_api VALUES (7, '/api/system/user', 'GET', '用户管理', '分页查询用户列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '用户列表');
INSERT INTO sys_api VALUES (8, '/api/system/user/:id', 'PUT', '用户管理', '修改用户基本信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑用户');
INSERT INTO sys_api VALUES (9, '/api/system/user/:id', 'DELETE', '用户管理', '软删除用户', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除用户');
INSERT INTO sys_api VALUES (10, '/api/system/user/:id', 'GET', '用户管理', '查询单个用户信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '用户详情');
INSERT INTO sys_api VALUES (11, '/api/system/user/:id/reset-password', 'POST', '用户管理', '重置用户密码', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '重置密码');
INSERT INTO sys_api VALUES (12, '/api/system/role', 'POST', '角色管理', '创建系统角色', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增角色');
INSERT INTO sys_api VALUES (13, '/api/system/role', 'GET', '角色管理', '分页查询角色列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '角色列表');
INSERT INTO sys_api VALUES (14, '/api/system/role/:id', 'PUT', '角色管理', '修改角色信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑角色');
INSERT INTO sys_api VALUES (15, '/api/system/role/:id', 'DELETE', '角色管理', '删除角色', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除角色');
INSERT INTO sys_api VALUES (16, '/api/system/role/:id', 'GET', '角色管理', '查询单个角色信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '角色详情');
INSERT INTO sys_api VALUES (17, '/api/system/role/:id/permissions', 'PUT', '角色管理', '分配角色的菜单和接口权限', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '分配权限');
INSERT INTO sys_api VALUES (18, '/api/system/role/all', 'GET', '角色管理', '查询所有角色（下拉选择用）', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '全部角色');
INSERT INTO sys_api VALUES (19, '/api/system/menu', 'POST', '菜单管理', '创建菜单或按钮', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增菜单');
INSERT INTO sys_api VALUES (20, '/api/system/menu/tree', 'GET', '菜单管理', '获取完整菜单树（管理页用）', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '菜单树');
INSERT INTO sys_api VALUES (21, '/api/system/menu/current', 'GET', '菜单管理', '获取当前登录用户的菜单树', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '当前用户菜单');
INSERT INTO sys_api VALUES (22, '/api/system/menu/:id', 'PUT', '菜单管理', '修改菜单信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑菜单');
INSERT INTO sys_api VALUES (23, '/api/system/menu/:id', 'DELETE', '菜单管理', '软删除菜单', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除菜单');
INSERT INTO sys_api VALUES (24, '/api/system/api', 'POST', '接口管理', '创建系统接口', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增接口');
INSERT INTO sys_api VALUES (25, '/api/system/api', 'GET', '接口管理', '分页查询接口列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '接口列表');
INSERT INTO sys_api VALUES (26, '/api/system/api/:id', 'PUT', '接口管理', '修改接口信息', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑接口');
INSERT INTO sys_api VALUES (27, '/api/system/api/:id', 'DELETE', '接口管理', '删除接口', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除接口');
INSERT INTO sys_api VALUES (28, '/api/system/api/all', 'GET', '接口管理', '查询所有接口（分配权限用）', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '全部接口');
INSERT INTO sys_api VALUES (29, '/api/system/dict/type', 'POST', '字典管理', '创建字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增字典类型');
INSERT INTO sys_api VALUES (30, '/api/system/dict/type', 'GET', '字典管理', '分页查询字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '字典类型列表');
INSERT INTO sys_api VALUES (31, '/api/system/dict/type/:id', 'PUT', '字典管理', '修改字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑字典类型');
INSERT INTO sys_api VALUES (32, '/api/system/dict/type/:id', 'DELETE', '字典管理', '删除字典类型', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除字典类型');
INSERT INTO sys_api VALUES (33, '/api/system/dict/data', 'POST', '字典管理', '创建字典数据', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '新增字典数据');
INSERT INTO sys_api VALUES (34, '/api/system/dict/data/:dictType', 'GET', '字典管理', '根据字典类型查询数据列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '字典数据列表');
INSERT INTO sys_api VALUES (35, '/api/system/dict/data/:id', 'PUT', '字典管理', '修改字典数据', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '编辑字典数据');
INSERT INTO sys_api VALUES (36, '/api/system/dict/data/:id', 'DELETE', '字典管理', '删除字典数据', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除字典数据');
INSERT INTO sys_api VALUES (37, '/api/system/file/upload', 'POST', '文件管理', '上传文件到服务器', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '上传文件');
INSERT INTO sys_api VALUES (38, '/api/system/file', 'GET', '文件管理', '分页查询文件列表', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '文件列表');
INSERT INTO sys_api VALUES (39, '/api/system/file/:id', 'DELETE', '文件管理', '删除文件', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '删除文件');
INSERT INTO sys_api VALUES (40, '/api/system/log/login', 'GET', '日志管理', '分页查询登录日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '登录日志列表');
INSERT INTO sys_api VALUES (41, '/api/system/log/login/clear', 'DELETE', '日志管理', '清除全部登录日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '清空登录日志');
INSERT INTO sys_api VALUES (42, '/api/system/log/oper', 'GET', '日志管理', '分页查询操作日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '操作日志列表');
INSERT INTO sys_api VALUES (43, '/api/system/log/oper/clear', 'DELETE', '日志管理', '清除全部操作日志', '2026-04-30 17:35:41', '2026-04-30 17:35:41', NULL, '清空操作日志');
INSERT INTO sys_api VALUES (44, '/api/nothing', 'GET', 'none', '只更新了名称', '2026-04-30 17:44:00', '2026-05-03 10:40:27', '2026-05-03 10:40:26', '不存在');
INSERT INTO sys_api VALUES (45, '/api/test/create', 'POST', 'test', '', '2026-04-30 17:45:06', '2026-05-03 10:41:42', '2026-05-03 10:41:43', '测试接口-创建');
INSERT INTO sys_api VALUES (46, '/api/test/:id', 'GET', 'test', '根据ID获取详情', '2026-04-30 17:45:34', '2026-04-30 17:45:34', NULL, '测试接口-详情');
INSERT INTO sys_api VALUES (47, '/api/test/:id', 'PUT', 'test', '', '2026-04-30 17:45:54', '2026-04-30 17:45:54', NULL, '测试接口-更新');
INSERT INTO sys_api VALUES (48, '/api/system/casbin/rules', 'GET', '权限管理', '只读查询Casbin策略列表', '2026-05-13 00:00:00', '2026-05-13 00:00:00', NULL, 'Casbin策略列表');

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS sys_dict_data;
CREATE TABLE sys_dict_data (
  id BIGSERIAL NOT NULL,
  type_id BIGINT NOT NULL,
  label VARCHAR(128) NOT NULL,
  dict_value VARCHAR(128) NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(512) NOT NULL DEFAULT '',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE INDEX sys_dict_data_idx_type_id ON sys_dict_data (type_id);
CREATE INDEX sys_dict_data_idx_status ON sys_dict_data (status);

-- ----------------------------
-- Records of sys_dict_data
-- ----------------------------
INSERT INTO sys_dict_data VALUES (1, 1, '启用', '1', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (2, 1, '禁用', '0', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (3, 2, '目录', '0', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (4, 2, '菜单', '1', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (5, 2, '按钮', '2', 3, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (6, 3, '显示', '1', 1, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (7, 3, '隐藏', '0', 2, 1, '', '2026-04-07 15:28:49', '2026-04-07 15:28:49', NULL);
INSERT INTO sys_dict_data VALUES (14, 1, '只改标签', 'test_value', 10, 1, '更新后备注', '2026-05-03 11:03:59', '2026-05-03 16:00:12', '2026-05-03 16:00:12');
INSERT INTO sys_dict_data VALUES (15, 6, '简单标签', 'simple', 0, 1, '', '2026-05-03 14:54:03', '2026-05-03 16:07:52', '2026-05-03 16:07:53');
INSERT INTO sys_dict_data VALUES (17, 6, '禁用标签', 'disabled_value', 99, 0, '', '2026-05-03 15:12:24', '2026-05-03 16:07:52', '2026-05-03 16:07:53');

-- ----------------------------
-- Table structure for sys_dict_type
-- ----------------------------
DROP TABLE IF EXISTS sys_dict_type;
CREATE TABLE sys_dict_type (
  id BIGSERIAL NOT NULL,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(512) NOT NULL DEFAULT '',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_dict_type_uk_code ON sys_dict_type (code);

-- ----------------------------
-- Records of sys_dict_type
-- ----------------------------
INSERT INTO sys_dict_type VALUES (1, '系统状态', 'sys_status', 1, '通用状态：启用/禁用', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_dict_type VALUES (2, '菜单类型', 'sys_menu_type', 1, '菜单类型：目录/菜单/按钮', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_dict_type VALUES (3, '是否显示', 'sys_visible', 1, '菜单是否显示', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_dict_type VALUES (6, '更新后字典名2', 'test_dict_type_updated222', 1, '更新后备注2', '2026-05-03 10:53:40', '2026-05-03 16:07:52', '2026-05-03 16:07:53');
INSERT INTO sys_dict_type VALUES (7, '最简字典', 'test_simple', 1, '', '2026-05-03 10:56:56', '2026-05-03 16:08:56', '2026-05-03 16:08:56');
INSERT INTO sys_dict_type VALUES (8, '不存在', 'not_exist', 1, '创建时即禁用', '2026-05-03 10:57:18', '2026-05-03 15:50:00', NULL);

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS sys_file;
CREATE TABLE sys_file (
  id BIGSERIAL NOT NULL,
  filename VARCHAR(256) NOT NULL,
  origin_name VARCHAR(256) NOT NULL DEFAULT '',
  file_path VARCHAR(512) NOT NULL DEFAULT '',
  file_url VARCHAR(512) NOT NULL DEFAULT '',
  file_size BIGINT NOT NULL DEFAULT 0,
  file_type VARCHAR(64) NOT NULL DEFAULT '',
  uploader_id BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE INDEX sys_file_idx_uploader_id ON sys_file (uploader_id);
CREATE INDEX sys_file_idx_created_at ON sys_file (created_at);

-- ----------------------------
-- Records of sys_file
-- ----------------------------

-- ----------------------------
-- Table structure for sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS sys_login_log;
CREATE TABLE sys_login_log (
  id BIGSERIAL NOT NULL,
  user_id BIGINT NOT NULL DEFAULT 0,
  username VARCHAR(64) NOT NULL DEFAULT '',
  ip VARCHAR(64) NOT NULL DEFAULT '',
  location VARCHAR(128) NOT NULL DEFAULT '',
  browser VARCHAR(128) NOT NULL DEFAULT '',
  os VARCHAR(64) NOT NULL DEFAULT '',
  status SMALLINT NOT NULL DEFAULT 1,
  msg VARCHAR(256) NOT NULL DEFAULT '',
  login_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE INDEX sys_login_log_idx_user_id ON sys_login_log (user_id);
CREATE INDEX sys_login_log_idx_username ON sys_login_log (username);
CREATE INDEX sys_login_log_idx_login_time ON sys_login_log (login_time);

-- ----------------------------
-- Records of sys_login_log
-- ----------------------------
INSERT INTO sys_login_log VALUES (1, 1, 'admin', '127.0.0.1', '', '其他', '其他', 1, '登录成功', '2026-04-28 22:35:48', NULL);
INSERT INTO sys_login_log VALUES (2, 1, 'admin', '127.0.0.1', '', 'Postman', '其他', 1, '登录成功', '2026-04-29 07:36:27', NULL);
INSERT INTO sys_login_log VALUES (3, 1, 'admin', '127.0.0.1', '', 'Postman', '其他', 1, '登录成功', '2026-04-29 07:42:16', NULL);
INSERT INTO sys_login_log VALUES (4, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-04-29 09:31:44', NULL);
INSERT INTO sys_login_log VALUES (5, 1, 'admin', '127.0.0.1', '', 'Postman', '其他', 1, '登录成功', '2026-04-30 07:04:42', NULL);
INSERT INTO sys_login_log VALUES (6, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-04-30 09:58:29', NULL);
INSERT INTO sys_login_log VALUES (7, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-04-30 17:43:16', NULL);
INSERT INTO sys_login_log VALUES (8, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-05-03 10:23:35', NULL);
INSERT INTO sys_login_log VALUES (9, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-05-05 09:29:38', NULL);
INSERT INTO sys_login_log VALUES (10, 1, 'admin', '[::1]:52934', '', '', '', 1, '登录成功', '2026-05-06 15:07:33', NULL);
INSERT INTO sys_login_log VALUES (11, 1, 'admin', '[::1]:52572', '', '', '', 1, '登录成功', '2026-05-06 16:03:31', NULL);
INSERT INTO sys_login_log VALUES (12, 1, 'admin', '[::1]:65434', '', '', '', 0, '密码错误', '2026-05-06 16:03:56', NULL);
INSERT INTO sys_login_log VALUES (13, 1, 'admin', '[::1]:62171', '', '', '', 0, '密码错误', '2026-05-06 16:06:15', NULL);
INSERT INTO sys_login_log VALUES (14, 1, 'admin', '[::1]:64787', '', '', '', 0, '密码错误', '2026-05-06 16:08:19', NULL);
INSERT INTO sys_login_log VALUES (15, 1, 'admin', '[::1]:53437', '', '', '', 0, '密码错误', '2026-05-06 16:18:55', NULL);
INSERT INTO sys_login_log VALUES (16, 1, 'admin', '[::1]:55671', '', '', '', 0, '密码错误', '2026-05-06 16:21:53', NULL);
INSERT INTO sys_login_log VALUES (17, 1, 'admin', '[::1]:65364', '', '', '', 0, '密码错误', '2026-05-06 16:24:43', NULL);
INSERT INTO sys_login_log VALUES (18, 1, 'admin', '[::1]:53073', '', '', '', 0, '密码错误', '2026-05-06 16:28:58', NULL);
INSERT INTO sys_login_log VALUES (19, 1, 'admin', '[::1]:53406', '', '', '', 0, '密码错误', '2026-05-06 16:34:01', NULL);
INSERT INTO sys_login_log VALUES (20, 1, 'admin', '[::1]:61567', '', '', '', 0, '密码错误', '2026-05-06 16:36:17', NULL);
INSERT INTO sys_login_log VALUES (21, 1, 'admin', '[::1]:55601', '', '', '', 1, '登录成功', '2026-05-07 10:19:26', NULL);
INSERT INTO sys_login_log VALUES (22, 1, 'admin', '[::1]:55601', '', '', '', 0, '密码错误', '2026-05-07 10:19:40', NULL);
INSERT INTO sys_login_log VALUES (23, 1, 'admin', '[::1]:54234', '', '', '', 1, '登录成功', '2026-05-07 10:49:05', NULL);
INSERT INTO sys_login_log VALUES (24, 1, 'admin', '[::1]:58537', '', '', '', 1, '登录成功', '2026-05-07 10:57:59', NULL);
INSERT INTO sys_login_log VALUES (25, 1, 'admin', '[::1]:63271', '', '', '', 1, '登录成功', '2026-05-07 14:52:35', NULL);
INSERT INTO sys_login_log VALUES (26, 1, 'admin', '[::1]:55467', '', '', '', 0, '密码错误', '2026-05-07 14:57:59', NULL);
INSERT INTO sys_login_log VALUES (27, 1, 'admin', '[::1]:55467', '', '', '', 1, '登录成功', '2026-05-07 14:58:04', NULL);
INSERT INTO sys_login_log VALUES (28, 1, 'admin', '[::1]:53366', '', '', '', 1, '登录成功', '2026-05-07 15:45:22', NULL);
INSERT INTO sys_login_log VALUES (29, 1, 'admin', '[::1]:64388', '', '', '', 1, '登录成功', '2026-05-07 19:14:51', NULL);
INSERT INTO sys_login_log VALUES (30, 1, 'admin', '[::1]:52524', '', '', '', 1, '登录成功', '2026-05-08 09:47:48', NULL);
INSERT INTO sys_login_log VALUES (31, 1, 'admin', '[::1]', '', 'Postman', '其他', 1, '登录成功', '2026-05-08 10:16:55', NULL);
INSERT INTO sys_login_log VALUES (32, 1, 'admin', '[::1]:53513', '', '', '', 1, '登录成功', '2026-05-08 11:14:33', NULL);

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS sys_menu;
CREATE TABLE sys_menu (
  id BIGSERIAL NOT NULL,
  parent_id BIGINT NOT NULL DEFAULT 0,
  name VARCHAR(64) NOT NULL,
  menu_path VARCHAR(256) NOT NULL DEFAULT '',
  component VARCHAR(256) NOT NULL DEFAULT '',
  icon VARCHAR(128) NOT NULL DEFAULT '',
  menu_type SMALLINT NOT NULL DEFAULT 0,
  permission VARCHAR(128) NOT NULL DEFAULT '',
  sort INTEGER NOT NULL DEFAULT 0,
  visible SMALLINT NOT NULL DEFAULT 1,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(512) NOT NULL DEFAULT '',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE INDEX sys_menu_idx_parent_id ON sys_menu (parent_id);
CREATE INDEX sys_menu_idx_type ON sys_menu (menu_type);

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO sys_menu VALUES (1, 0, '系统管理', '/system', 'Layout', 'setting', 0, '', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (2, 1, '用户管理', '/system/user', 'system/user/index', 'user', 1, 'system:user:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (3, 1, '角色管理', '/system/role', 'system/role/index', 'peoples', 1, 'system:role:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (4, 1, '菜单管理', '/system/menu', 'system/menu/index', 'tree-table', 1, 'system:menu:list', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (5, 1, '接口管理', '/system/api', 'system/api/index', 'api', 1, 'system:api:list', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (6, 1, '字典管理', '/system/dict', 'system/dict/index', 'dict', 1, 'system:dict:list', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (7, 0, '日志管理', '/log', 'Layout', 'log', 0, '', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (8, 7, '登录日志', '/log/login', 'log/login/index', 'logininfor', 1, 'system:loginlog:list', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (9, 7, '操作日志', '/log/oper', 'log/oper/index', 'form', 1, 'system:operlog:list', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (20, 2, '新增用户', '', '', '', 2, 'system:user:create', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (21, 2, '编辑用户', '', '', '', 2, 'system:user:update', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (22, 2, '删除用户', '', '', '', 2, 'system:user:delete', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (23, 3, '新增角色', '', '', '', 2, 'system:role:create', 1, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (24, 3, '编辑角色', '', '', '', 2, 'system:role:update', 2, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (25, 3, '删除角色', '', '', '', 2, 'system:role:delete', 3, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (26, 3, '分配菜单', '', '', '', 2, 'system:role:assignmenu', 4, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (27, 3, '分配接口', '', '', '', 2, 'system:role:assignapi', 5, 1, 1, '', '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_menu VALUES (33, 2, '重置密码', '', '', '', 2, 'system:user:resetpwd', 4, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (34, 4, '新增菜单', '', '', '', 2, 'system:menu:create', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (35, 4, '编辑菜单', '', '', '', 2, 'system:menu:update', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (36, 4, '删除菜单', '', '', '', 2, 'system:menu:delete', 3, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (37, 5, '新增接口', '', '', '', 2, 'system:api:create', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (38, 5, '编辑接口', '', '', '', 2, 'system:api:update', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (39, 5, '删除接口', '', '', '', 2, 'system:api:delete', 3, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (40, 6, '新增字典类型', '', '', '', 2, 'system:dict:type:create', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (41, 6, '编辑字典类型', '', '', '', 2, 'system:dict:type:update', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (42, 6, '删除字典类型', '', '', '', 2, 'system:dict:type:delete', 3, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (43, 6, '新增字典数据', '', '', '', 2, 'system:dict:data:create', 4, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (44, 6, '编辑字典数据', '', '', '', 2, 'system:dict:data:update', 5, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (45, 6, '删除字典数据', '', '', '', 2, 'system:dict:data:delete', 6, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (46, 8, '清空日志', '', '', '', 2, 'system:loginlog:clear', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (47, 9, '清空日志', '', '', '', 2, 'system:operlog:clear', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (48, 1, '文件管理', '/system/file', 'system/file/index', 'upload', 1, 'system:file:list', 6, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (49, 48, '上传文件', '', '', '', 2, 'system:file:upload', 1, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);
INSERT INTO sys_menu VALUES (50, 48, '删除文件', '', '', '', 2, 'system:file:delete', 2, 1, 1, '', '2026-04-30 17:29:16', '2026-04-30 17:29:16', NULL);

-- ----------------------------
-- Table structure for sys_oper_log
-- ----------------------------
DROP TABLE IF EXISTS sys_oper_log;
CREATE TABLE sys_oper_log (
  id BIGSERIAL NOT NULL,
  title VARCHAR(64) NOT NULL DEFAULT '',
  business_type INTEGER NOT NULL DEFAULT 0,
  method VARCHAR(256) NOT NULL DEFAULT '',
  request_method VARCHAR(16) NOT NULL DEFAULT '',
  operator_type INTEGER NOT NULL DEFAULT 0,
  operator_name VARCHAR(64) NOT NULL DEFAULT '',
  operator_id BIGINT NOT NULL DEFAULT 0,
  dept_name VARCHAR(64) NOT NULL DEFAULT '',
  oper_url VARCHAR(512) NOT NULL DEFAULT '',
  oper_ip VARCHAR(64) NOT NULL DEFAULT '',
  oper_param TEXT NULL,
  json_result TEXT NULL,
  status SMALLINT NOT NULL DEFAULT 1,
  error_msg VARCHAR(2000) NOT NULL DEFAULT '',
  oper_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE INDEX sys_oper_log_idx_operator_id ON sys_oper_log (operator_id);
CREATE INDEX sys_oper_log_idx_oper_time ON sys_oper_log (oper_time);
CREATE INDEX sys_oper_log_idx_status ON sys_oper_log (status);

-- ----------------------------
-- Records of sys_oper_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS sys_role;
CREATE TABLE sys_role (
  id BIGSERIAL NOT NULL,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(512) NOT NULL DEFAULT '',
  sort INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_role_uk_code ON sys_role (code);
CREATE INDEX sys_role_idx_status ON sys_role (status);
CREATE INDEX sys_role_idx_deleted_at ON sys_role (deleted_at);

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO sys_role VALUES (1, '超级管理员', 'admin', 1, '系统内置超级管理员角色，拥有所有权限', 1, '2026-04-07 15:28:48', '2026-04-07 15:28:48', NULL);
INSERT INTO sys_role VALUES (6, '管理员', 'manager', 1, '能管人、管角色、管接口、但不碰业务数据', 2, '2026-04-29 11:32:53', '2026-04-29 11:32:53', NULL);
INSERT INTO sys_role VALUES (7, '操作者', 'operator', 1, '能录入业务数据、查看状态，不能管理系统', 3, '2026-04-29 14:35:57', '2026-04-29 16:24:47', NULL);
INSERT INTO sys_role VALUES (8, '查看者', 'viewer', 1, '只能看，什么都不能改', 4, '2026-04-29 14:36:36', '2026-04-29 17:28:48', NULL);

-- ----------------------------
-- Table structure for sys_role_api
-- ----------------------------
DROP TABLE IF EXISTS sys_role_api;
CREATE TABLE sys_role_api (
  id BIGSERIAL NOT NULL,
  role_id BIGINT NOT NULL,
  api_id BIGINT NOT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_role_api_uk_role_api ON sys_role_api (role_id, api_id);
CREATE INDEX sys_role_api_idx_api_id ON sys_role_api (api_id);

-- ----------------------------
-- Records of sys_role_api
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS sys_role_menu;
CREATE TABLE sys_role_menu (
  id BIGSERIAL NOT NULL,
  role_id BIGINT NOT NULL,
  menu_id BIGINT NOT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_role_menu_uk_role_menu ON sys_role_menu (role_id, menu_id);
CREATE INDEX sys_role_menu_idx_menu_id ON sys_role_menu (menu_id);

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO sys_role_menu VALUES (1, 1, 1);
INSERT INTO sys_role_menu VALUES (2, 1, 2);
INSERT INTO sys_role_menu VALUES (8, 6, 1);
INSERT INTO sys_role_menu VALUES (13, 6, 2);
INSERT INTO sys_role_menu VALUES (9, 6, 3);
INSERT INTO sys_role_menu VALUES (18, 6, 4);
INSERT INTO sys_role_menu VALUES (3, 6, 5);
INSERT INTO sys_role_menu VALUES (4, 6, 6);
INSERT INTO sys_role_menu VALUES (14, 6, 7);
INSERT INTO sys_role_menu VALUES (5, 6, 8);
INSERT INTO sys_role_menu VALUES (10, 6, 9);
INSERT INTO sys_role_menu VALUES (11, 6, 20);
INSERT INTO sys_role_menu VALUES (6, 6, 21);
INSERT INTO sys_role_menu VALUES (19, 6, 22);
INSERT INTO sys_role_menu VALUES (15, 6, 23);
INSERT INTO sys_role_menu VALUES (7, 6, 24);
INSERT INTO sys_role_menu VALUES (12, 6, 25);
INSERT INTO sys_role_menu VALUES (16, 6, 26);
INSERT INTO sys_role_menu VALUES (17, 6, 27);
INSERT INTO sys_role_menu VALUES (39, 8, 1);
INSERT INTO sys_role_menu VALUES (42, 8, 2);
INSERT INTO sys_role_menu VALUES (40, 8, 6);
INSERT INTO sys_role_menu VALUES (41, 8, 7);
INSERT INTO sys_role_menu VALUES (44, 8, 8);
INSERT INTO sys_role_menu VALUES (43, 8, 20);
INSERT INTO sys_role_menu VALUES (38, 8, 42);

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS sys_user;
CREATE TABLE sys_user (
  id BIGSERIAL NOT NULL,
  username VARCHAR(64) NOT NULL,
  password VARCHAR(128) NOT NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  email VARCHAR(128) NOT NULL DEFAULT '',
  phone VARCHAR(20) NOT NULL DEFAULT '',
  avatar VARCHAR(512) NOT NULL DEFAULT '',
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(512) NOT NULL DEFAULT '',
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITHOUT TIME ZONE NULL DEFAULT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_user_uk_username ON sys_user (username);
CREATE INDEX sys_user_idx_status ON sys_user (status);
CREATE INDEX sys_user_idx_deleted_at ON sys_user (deleted_at);

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO sys_user VALUES (1, 'admin', '$2a$10$vh6VEcx58MphNdTaShd...ErM.V1CkWwYOt5wo47qy0soZ6SKZqKK', '系统管理员', 'admin@example.com', '18288888888', '', 1, '系统内置超级管理员账号', '2026-04-07 15:28:48', '2026-05-07 14:57:52', NULL);
INSERT INTO sys_user VALUES (2, 'testuser01', '$2a$10$1ZXON13uBALtKlyzAB7UOeoS0UaQnrqNsoauedo5Q4dFTy1tbZbAu', '只改昵称', 'updated@example.com', '13900139001', 'https://example.com/new-avatar.png', 1, '更新备注', '2026-04-30 07:05:13', '2026-05-05 11:40:47', NULL);
INSERT INTO sys_user VALUES (3, 'testuser02', '$2a$10$pRRoEYtW/fiFGwRvVXtneu3SHKuxaDwMkVoof0HPTIHvE.WQ53gkK', '测试用户02', 'test01@example.com', '13800138002', 'https://example.com/avatar.png', 1, '测试创建的用户2', '2026-04-30 07:12:36', '2026-04-30 07:12:36', NULL);
INSERT INTO sys_user VALUES (4, 'testuser03', '$2a$10$m1cUikaSJFBVHMfUa4XkX.leDwvgDpIb86KaMhu2MbILx13EjOifq', '测试用户03', 'test01@example.com', '13800138003', 'https://example.com/avatar.png', 1, '测试创建的用户2', '2026-04-30 07:20:06', '2026-05-05 13:24:32', '2026-05-05 13:24:34');
INSERT INTO sys_user VALUES (9, 'tianyx', '$2a$10$gFiCU3xpbadRwvsoPCXz2.i3U7AEHKJJblE67JewjbCNuEEBVuf/C', '管理员01', '1787817883@qq.com', '1888888888', 'https://example.com/avatar.png', 1, '', '2026-05-05 09:56:55', '2026-05-05 09:56:55', NULL);

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS sys_user_role;
CREATE TABLE sys_user_role (
  id BIGSERIAL NOT NULL,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX sys_user_role_uk_user_role ON sys_user_role (user_id, role_id);
CREATE INDEX sys_user_role_idx_role_id ON sys_user_role (role_id);

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO sys_user_role VALUES (3, 1, 1);
INSERT INTO sys_user_role VALUES (4, 1, 8);
INSERT INTO sys_user_role VALUES (8, 9, 6);
INSERT INTO sys_user_role VALUES (9, 9, 7);


-- Reset BIGSERIAL sequences to match the original MySQL next values.
SELECT setval(pg_get_serial_sequence('casbin_rule', 'id')::regclass, 57, false);
SELECT setval(pg_get_serial_sequence('sys_api', 'id')::regclass, 49, false);
SELECT setval(pg_get_serial_sequence('sys_dict_data', 'id')::regclass, 18, false);
SELECT setval(pg_get_serial_sequence('sys_dict_type', 'id')::regclass, 9, false);
SELECT setval(pg_get_serial_sequence('sys_file', 'id')::regclass, 1, false);
SELECT setval(pg_get_serial_sequence('sys_login_log', 'id')::regclass, 33, false);
SELECT setval(pg_get_serial_sequence('sys_menu', 'id')::regclass, 51, false);
SELECT setval(pg_get_serial_sequence('sys_oper_log', 'id')::regclass, 1, false);
SELECT setval(pg_get_serial_sequence('sys_role', 'id')::regclass, 9, false);
SELECT setval(pg_get_serial_sequence('sys_role_api', 'id')::regclass, 8, false);
SELECT setval(pg_get_serial_sequence('sys_role_menu', 'id')::regclass, 45, false);
SELECT setval(pg_get_serial_sequence('sys_user', 'id')::regclass, 10, false);
SELECT setval(pg_get_serial_sequence('sys_user_role', 'id')::regclass, 13, false);
