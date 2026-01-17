-- 用户数据库表结构设计
-- 创建时间: 2026-01-17
-- 说明: 用户管理功能相关数据库表

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码(加密后)',
  `email` VARCHAR(100) UNIQUE COMMENT '邮箱',
  `phone` VARCHAR(20) COMMENT '手机号',
  `real_name` VARCHAR(50) COMMENT '真实姓名',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `last_login_time` DATETIME COMMENT '最后登录时间',
  `last_login_ip` VARCHAR(50) COMMENT '最后登录IP',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` INT COMMENT '创建人ID',
  `updated_by` INT COMMENT '更新人ID',
  INDEX `idx_username` (`username`),
  INDEX `idx_email` (`email`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
  `role_name` VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
  `role_code` VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
  `description` VARCHAR(200) COMMENT '角色描述',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` INT COMMENT '创建人ID',
  `updated_by` INT COMMENT '更新人ID',
  INDEX `idx_role_code` (`role_code`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 权限表
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '权限ID',
  `permission_name` VARCHAR(50) NOT NULL COMMENT '权限名称',
  `permission_code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
  `permission_type` VARCHAR(20) NOT NULL COMMENT '权限类型: menu-菜单, button-按钮, api-接口',
  `parent_id` INT DEFAULT 0 COMMENT '父权限ID',
  `path` VARCHAR(200) COMMENT '路由路径',
  `component` VARCHAR(200) COMMENT '组件路径',
  `icon` VARCHAR(50) COMMENT '图标',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `description` VARCHAR(200) COMMENT '权限描述',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_permission_code` (`permission_code`),
  INDEX `idx_parent_id` (`parent_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
  `user_id` INT NOT NULL COMMENT '用户ID',
  `role_id` INT NOT NULL COMMENT '角色ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` INT COMMENT '创建人ID',
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_role_id` (`role_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
  `role_id` INT NOT NULL COMMENT '角色ID',
  `permission_id` INT NOT NULL COMMENT '权限ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` INT COMMENT '创建人ID',
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  INDEX `idx_role_id` (`role_id`),
  INDEX `idx_permission_id` (`permission_id`),
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
  `user_id` INT COMMENT '用户ID',
  `username` VARCHAR(50) COMMENT '用户名',
  `operation` VARCHAR(100) NOT NULL COMMENT '操作类型',
  `method` VARCHAR(10) COMMENT '请求方法',
  `url` VARCHAR(500) COMMENT '请求URL',
  `ip` VARCHAR(50) COMMENT 'IP地址',
  `location` VARCHAR(100) COMMENT '地理位置',
  `params` TEXT COMMENT '请求参数',
  `result` TEXT COMMENT '返回结果',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-失败, 1-成功',
  `error_msg` VARCHAR(500) COMMENT '错误信息',
  `execute_time` INT COMMENT '执行时间(ms)',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_username` (`username`),
  INDEX `idx_operation` (`operation`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 插入初始数据
-- 插入默认管理员角色
INSERT INTO `roles` (`role_name`, `role_code`, `description`, `status`) VALUES
('超级管理员', 'ADMIN', '系统超级管理员，拥有所有权限', 1),
('普通用户', 'USER', '普通用户，拥有基础权限', 1),
('审计员', 'AUDITOR', '审计员，只能查看日志和报表', 1);

-- 插入基础权限
INSERT INTO `permissions` (`permission_name`, `permission_code`, `permission_type`, `parent_id`, `path`, `component`, `icon`, `sort_order`, `description`, `status`) VALUES
('用户管理', 'user:manage', 'menu', 0, '/user', 'UserManage', 'UserOutlined', 1, '用户管理菜单', 1),
('用户列表', 'user:list', 'api', 1, '/api/user/list', '', '', 1, '获取用户列表接口', 1),
('用户添加', 'user:add', 'button', 1, '', '', '', 2, '添加用户按钮', 1),
('用户编辑', 'user:edit', 'button', 1, '', '', '', 3, '编辑用户按钮', 1),
('用户删除', 'user:delete', 'button', 1, '', '', '', 4, '删除用户按钮', 1),
('角色管理', 'role:manage', 'menu', 0, '/role', 'RoleManage', 'TeamOutlined', 2, '角色管理菜单', 1),
('角色列表', 'role:list', 'api', 6, '/api/role/list', '', '', 1, '获取角色列表接口', 1),
('权限管理', 'permission:manage', 'menu', 0, '/permission', 'PermissionManage', 'SafetyOutlined', 3, '权限管理菜单', 1),
('日志管理', 'log:manage', 'menu', 0, '/log', 'LogManage', 'FileTextOutlined', 4, '日志管理菜单', 1);

-- 插入默认管理员用户 (密码: admin123, 需要使用BCrypt加密)
-- 注意: 实际使用时需要使用密码加密工具生成加密后的密码
INSERT INTO `users` (`username`, `password`, `email`, `real_name`, `status`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@example.com', '系统管理员', 1);

-- 为管理员分配角色
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- 为超级管理员角色分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions`;
