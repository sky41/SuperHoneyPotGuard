-- 提交目的：修复 HFish API 调用问题
-- 提交内容：分析 HFish API 错误原因并提供解决方案
-- 提交时间：2026-01-18

-- ========================================
-- HFish API 错误分析
-- ========================================

-- 错误信息：
-- response_code: 1004
-- verbose_msg: "请求数据非法"

-- 可能的原因：
-- 1. API Key 不正确或已过期
-- 2. 端点 URL 不正确
-- 3. 请求方法不正确
-- 4. 请求参数格式不正确

-- ========================================
-- 解决方案
-- ========================================

-- 方案 1：验证 API Key
-- 确认 API Key 是否正确：kGZJpRLpvBCTgbXTfEjaOfwupNLXZsYuLMXUvmboyuagOVZAoUlXJbzbuyGUkdRT
-- 检查 API Key 是否已过期
-- 联系 HFish 管理员获取新的 API Key

-- 方案 2：验证端点 URL
-- 检查端点 URL 是否正确
-- 参考 HFish API 文档确认正确的端点

-- 方案 3：验证请求方法
-- 确认每个端点使用的 HTTP 方法是否正确
-- GET: /api/v1/hfish/sys_info
-- POST: /api/v1/attack/ip
-- POST: /api/v1/attack/detail
-- POST: /api/v1/attack/account

-- 方案 4：验证请求参数
-- 确认请求参数格式是否正确
-- 检查是否需要特定的请求头

-- 方案 5：测试 API 连接
-- 使用 curl 命令测试 API 连接
-- 检查网络连接是否正常
-- 检查防火墙设置

-- ========================================
-- 测试命令
-- ========================================

-- 测试获取系统信息
-- curl "https://115.190.62.202:4433/api/v1/hfish/sys_info?api_key=kGZJpRLpvBCTgbXTfEjaOfwupNLXZsYuLMXUvmboyuagOVZAoUlXJbzbuyGUkdRT"

-- 测试获取攻击 IP
-- curl -X POST "https://115.190.62.202:4433/api/v1/attack/ip?api_key=kGZJpRLpvBCTgbXTfEjaOfwupNLXZsYuLMXUvmboyuagOVZAoUlXJbzbuyGUkdRT"

-- ========================================
-- 调试建议
-- ========================================

-- 1. 查看后端日志输出
-- 2. 检查 HFish API 服务器状态
-- 3. 检查网络连接
-- 4. 验证 API Key 是否有效
-- 5. 联系 HFish 技术支持

-- ========================================
-- 代码修改建议
-- ========================================

-- 1. 添加请求超时设置
-- 2. 添加重试机制
-- 3. 添加更详细的错误处理
-- 4. 添加 API 响应日志记录
