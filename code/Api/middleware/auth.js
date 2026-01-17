const { verifyToken } = require('../utils/jwt');
const logger = require('../utils/logger');

const authMiddleware = async (req, res, next) => {
  try {
    const token = req.headers.authorization?.replace('Bearer ', '');

    if (!token) {
      return res.status(401).json({
        success: false,
        message: '未提供认证令牌'
      });
    }

    const decoded = verifyToken(token);

    if (!decoded) {
      return res.status(401).json({
        success: false,
        message: '认证令牌无效或已过期'
      });
    }

    req.user = decoded;
    next();
  } catch (error) {
    logger.error('认证中间件错误:', error);
    return res.status(401).json({
      success: false,
      message: '认证失败'
    });
  }
};

module.exports = authMiddleware;
