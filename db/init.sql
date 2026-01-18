-- SuperHoneyPotGuard 数据库初始化脚本
-- 提交目的：整合所有数据库表结构和初始数据到一个文件中，避免多文件导致数据库异常
-- 提交内容：创建完整的数据库初始化脚本，包含所有表结构、索引、初始数据和权限配置
-- 提交时间：2026-01-18

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `superhoneypotguard` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `superhoneypotguard`;

-- ========================================
-- 1. 用户表
-- ========================================
DROP TABLE IF EXISTS `user_roles`;
DROP TABLE IF EXISTS `role_permissions`;
DROP TABLE IF EXISTS `operation_logs`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `roles`;
DROP TABLE IF EXISTS `permissions`;

CREATE TABLE `users` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
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
  `created_by` BIGINT COMMENT '创建人ID',
  `updated_by` BIGINT COMMENT '更新人ID',
  INDEX `idx_username` (`username`),
  INDEX `idx_email` (`email`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ========================================
-- 2. 角色表
-- ========================================
CREATE TABLE `roles` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
  `role_name` VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
  `role_code` VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
  `description` VARCHAR(200) COMMENT '角色描述',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` BIGINT COMMENT '创建人ID',
  `updated_by` BIGINT COMMENT '更新人ID',
  INDEX `idx_role_code` (`role_code`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ========================================
-- 3. 权限表
-- ========================================
CREATE TABLE `permissions` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '权限ID',
  `permission_name` VARCHAR(50) NOT NULL COMMENT '权限名称',
  `permission_code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
  `permission_type` VARCHAR(20) NOT NULL COMMENT '权限类型: menu-菜单, button-按钮, api-接口',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父权限ID',
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

-- ========================================
-- 4. 用户角色关联表
-- ========================================
CREATE TABLE `user_roles` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` BIGINT COMMENT '创建人ID',
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_role_id` (`role_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- ========================================
-- 5. 角色权限关联表
-- ========================================
CREATE TABLE `role_permissions` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '关联ID',
  `role_id` BIGINT NOT NULL COMMENT '角色ID',
  `permission_id` BIGINT NOT NULL COMMENT '权限ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` BIGINT COMMENT '创建人ID',
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  INDEX `idx_role_id` (`role_id`),
  INDEX `idx_permission_id` (`permission_id`),
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- ========================================
-- 6. 操作日志表
-- ========================================
CREATE TABLE `operation_logs` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
  `user_id` BIGINT COMMENT '用户ID',
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
  INDEX `idx_status` (`status`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- ========================================
-- 7. 初始数据
-- ========================================

-- 插入默认管理员用户
INSERT INTO `users` (`username`, `password`, `email`, `real_name`, `status`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@example.com', '系统管理员', 1);

-- 插入默认角色
INSERT INTO `roles` (`role_name`, `role_code`, `description`, `status`) VALUES
('超级管理员', 'admin', '拥有系统所有权限', 1),
('普通用户', 'user', '普通用户权限', 1);

-- 插入默认权限
INSERT INTO `permissions` (`permission_name`, `permission_code`, `permission_type`, `parent_id`, `path`, `component`, `icon`, `sort_order`, `description`, `status`) VALUES
('首页', 'dashboard:view', 'menu', 0, '/dashboard', 'Dashboard', 'DashboardOutlined', 0, '首页仪表盘', 1),
('系统管理', 'system', 'menu', 0, '/system', NULL, 'SettingOutlined', 1, '系统管理', 1),
('用户管理', 'user:manage', 'menu', 2, '/system/user', 'UserManage', 'UserOutlined', 1, '用户管理', 1),
('角色管理', 'role:manage', 'menu', 2, '/system/role', 'RoleManage', 'TeamOutlined', 2, '角色管理', 1),
('权限管理', 'permission:manage', 'menu', 2, '/system/permission', 'PermissionManage', 'SafetyOutlined', 3, '权限管理', 1),
('操作日志', 'log:manage', 'menu', 2, '/system/log', 'LogManage', 'FileTextOutlined', 4, '操作日志', 1),
('查看日志', 'log:view', 'button', 0, NULL, NULL, NULL, 1, '查看日志详情', 1),
('删除日志', 'log:delete', 'button', 0, NULL, NULL, NULL, 1, '删除日志', 1),
('清空日志', 'log:clear', 'button', 0, NULL, NULL, NULL, 1, '清空所有日志', 1);

-- 为超级管理员角色分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions` WHERE status = 1;

-- 为管理员分配超级管理员角色
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- ========================================
-- 8. 验证数据
-- ========================================

-- 验证用户和权限
SELECT
    u.id as user_id,
    u.username,
    r.id as role_id,
    r.role_name,
    p.id as permission_id,
    p.permission_code,
    p.permission_name
FROM users u
INNER JOIN user_roles ur ON u.id = ur.user_id
INNER JOIN roles r ON ur.role_id = r.id
INNER JOIN role_permissions rp ON r.id = rp.role_id
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE u.username = 'admin'
ORDER BY p.permission_code;

-- ========================================
-- 数据库初始化完成
-- ========================================
