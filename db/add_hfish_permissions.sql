-- 添加 HFish 数据查看和封禁权限
-- 提交目的：添加 HFish 数据查看和 IP 封禁权限到权限表
-- 提交内容：插入 hfish:view 和 hfish:block 权限，并分配给超级管理员角色
-- 提交时间：2026-01-18

-- 1. 插入 HFish 数据查看权限
INSERT INTO `permissions` (`permission_name`, `permission_code`, `permission_type`, `parent_id`, `path`, `component`, `icon`, `sort_order`, `description`, `status`)
VALUES
('HFish 数据', 'hfish:view', 'menu', 0, '/hfish', 'HFishData', 'SecurityScanOutlined', 5, 'HFish 蜜罐数据查看', 1);

-- 2. 插入 HFish IP 封禁权限
INSERT INTO `permissions` (`permission_name`, `permission_code`, `permission_type`, `parent_id`, `path`, `component`, `icon`, `sort_order`, `description`, `status`)
VALUES
('封禁 IP', 'hfish:block', 'button', 0, NULL, NULL, NULL, 1, '手动封禁 IP 地址', 1);

-- 3. 为超级管理员角色分配所有 HFish 权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions` WHERE permission_code LIKE 'hfish:%' AND status = 1;

-- 4. 验证结果
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
