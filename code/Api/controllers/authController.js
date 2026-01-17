const pool = require('../config/database');
const { hashPassword, comparePassword } = require('../utils/password');
const { generateToken } = require('../utils/jwt');
const logger = require('../utils/logger');

class AuthController {
  async register(req, res) {
    try {
      const { username, password, email, phone, realName } = req.body;

      const [existingUsers] = await pool.query(
        'SELECT id FROM users WHERE username = ? OR email = ?',
        [username, email]
      );

      if (existingUsers.length > 0) {
        return res.status(400).json({
          success: false,
          message: '用户名或邮箱已存在'
        });
      }

      const hashedPassword = await hashPassword(password);

      const [result] = await pool.query(
        'INSERT INTO users (username, password, email, phone, real_name, status) VALUES (?, ?, ?, ?, ?, 1)',
        [username, hashedPassword, email, phone, realName]
      );

      const [defaultRole] = await pool.query(
        'SELECT id FROM roles WHERE role_code = ?',
        ['USER']
      );

      if (defaultRole.length > 0) {
        await pool.query(
          'INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)',
          [result.insertId, defaultRole[0].id]
        );
      }

      logger.info(`用户注册成功: ${username}`);

      res.status(201).json({
        success: true,
        message: '注册成功',
        data: {
          id: result.insertId,
          username,
          email
        }
      });
    } catch (error) {
      logger.error('用户注册失败:', error);
      res.status(500).json({
        success: false,
        message: '注册失败，请稍后重试'
      });
    }
  }

  async login(req, res) {
    try {
      const { username, password } = req.body;

      const [users] = await pool.query(
        'SELECT * FROM users WHERE username = ?',
        [username]
      );

      if (users.length === 0) {
        return res.status(401).json({
          success: false,
          message: '用户名或密码错误'
        });
      }

      const user = users[0];

      if (user.status !== 1) {
        return res.status(403).json({
          success: false,
          message: '账号已被禁用'
        });
      }

      const isPasswordValid = await comparePassword(password, user.password);

      if (!isPasswordValid) {
        return res.status(401).json({
          success: false,
          message: '用户名或密码错误'
        });
      }

      await pool.query(
        'UPDATE users SET last_login_time = NOW(), last_login_ip = ? WHERE id = ?',
        [req.ip, user.id]
      );

      const [roles] = await pool.query(`
        SELECT r.id, r.role_name, r.role_code
        FROM roles r
        INNER JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = ? AND r.status = 1
      `, [user.id]);

      const [permissions] = await pool.query(`
        SELECT DISTINCT p.permission_code
        FROM permissions p
        INNER JOIN role_permissions rp ON p.id = rp.permission_id
        INNER JOIN user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = ? AND p.status = 1
      `, [user.id]);

      const token = generateToken({
        id: user.id,
        username: user.username,
        roles: roles.map(r => r.role_code),
        permissions: permissions.map(p => p.permission_code)
      });

      logger.info(`用户登录成功: ${username}`);

      res.json({
        success: true,
        message: '登录成功',
        data: {
          token,
          user: {
            id: user.id,
            username: user.username,
            email: user.email,
            realName: user.real_name,
            roles,
            permissions
          }
        }
      });
    } catch (error) {
      logger.error('用户登录失败:', error);
      res.status(500).json({
        success: false,
        message: '登录失败，请稍后重试'
      });
    }
  }

  async logout(req, res) {
    try {
      logger.info(`用户注销成功: ${req.user.username}`);

      res.json({
        success: true,
        message: '注销成功'
      });
    } catch (error) {
      logger.error('用户注销失败:', error);
      res.status(500).json({
        success: false,
        message: '注销失败'
      });
    }
  }

  async getCurrentUser(req, res) {
    try {
      const userId = req.user.id;

      const [users] = await pool.query(
        'SELECT id, username, email, phone, real_name, status, last_login_time, created_at FROM users WHERE id = ?',
        [userId]
      );

      if (users.length === 0) {
        return res.status(404).json({
          success: false,
          message: '用户不存在'
        });
      }

      const [roles] = await pool.query(`
        SELECT r.id, r.role_name, r.role_code
        FROM roles r
        INNER JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = ? AND r.status = 1
      `, [userId]);

      const [permissions] = await pool.query(`
        SELECT DISTINCT p.id, p.permission_name, p.permission_code, p.permission_type
        FROM permissions p
        INNER JOIN role_permissions rp ON p.id = rp.permission_id
        INNER JOIN user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = ? AND p.status = 1
      `, [userId]);

      res.json({
        success: true,
        data: {
          user: users[0],
          roles,
          permissions
        }
      });
    } catch (error) {
      logger.error('获取当前用户信息失败:', error);
      res.status(500).json({
        success: false,
        message: '获取用户信息失败'
      });
    }
  }
}

module.exports = new AuthController();
