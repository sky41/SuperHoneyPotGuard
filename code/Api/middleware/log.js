const pool = require('../config/database');
const logger = require('../utils/logger');

const logMiddleware = async (req, res, next) => {
  const startTime = Date.now();

  const originalSend = res.send;
  res.send = function (data) {
    res.locals.responseData = data;
    originalSend.call(this, data);
  };

  res.on('finish', async () => {
    try {
      const executeTime = Date.now() - startTime;
      const responseData = res.locals.responseData;
      let status = 1;

      if (typeof responseData === 'string') {
        try {
          const parsed = JSON.parse(responseData);
          status = parsed.success ? 1 : 0;
        } catch (e) {
          status = 0;
        }
      }

      const userId = req.user?.id || null;
      const username = req.user?.username || 'anonymous';

      await pool.query(`
        INSERT INTO operation_logs (user_id, username, operation, method, url, ip, params, result, status, execute_time)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
      `, [
        userId,
        username,
        req.route?.path || req.path,
        req.method,
        req.originalUrl,
        req.ip || req.connection.remoteAddress,
        JSON.stringify(req.body || req.query),
        typeof responseData === 'string' ? responseData.substring(0, 500) : JSON.stringify(responseData).substring(0, 500),
        status,
        executeTime
      ]);
    } catch (error) {
      logger.error('日志记录错误:', error);
    }
  });

  next();
};

module.exports = logMiddleware;
