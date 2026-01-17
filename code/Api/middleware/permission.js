const pool = require('../config/database');
const logger = require('../utils/logger');

const checkPermission = (permissionCode) => {
  return async (req, res, next) => {
    try {
      const userId = req.user.id;

      const [rows] = await pool.query(`
        SELECT DISTINCT p.permission_code
        FROM users u
        INNER JOIN user_roles ur ON u.id = ur.user_id
        INNER JOIN roles r ON ur.role_id = r.id
        INNER JOIN role_permissions rp ON r.id = rp.role_id
        INNER JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = ? AND p.permission_code = ? AND p.status = 1 AND r.status = 1 AND u.status = 1
      `, [userId, permissionCode]);

      if (rows.length === 0) {
        return res.status(403).json({
          success: false,
          message: '权限不足'
        });
      }

      next();
    } catch (error) {
      logger.error('权限验证中间件错误:', error);
      return res.status(500).json({
        success: false,
        message: '权限验证失败'
      });
    }
  };
};

module.exports = checkPermission;
