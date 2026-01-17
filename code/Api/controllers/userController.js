const pool = require('../config/database');
const { hashPassword } = require('../utils/password');
const logger = require('../utils/logger');

class UserController {
  async getUserList(req, res) {
    try {
      const { page = 1, pageSize = 10, username, status } = req.query;
      const offset = (page - 1) * pageSize;

      let whereClause = 'WHERE 1=1';
      const params = [];

      if (username) {
        whereClause += ' AND username LIKE ?';
        params.push(`%${username}%`);
      }

      if (status !== undefined) {
        whereClause += ' AND status = ?';
        params.push(status);
      }

      const [users] = await pool.query(`
        SELECT id, username, email, phone, real_name, status, last_login_time, last_login_ip, created_at, updated_at
        FROM users
        ${whereClause}
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
      `, [...params, parseInt(pageSize), offset]);

      const [countResult] = await pool.query(`
        SELECT COUNT(*) as total
        FROM users
        ${whereClause}
      `, params);

      const total = countResult[0].total;

      for (const user of users) {
        const [roles] = await pool.query(`
          SELECT r.id, r.role_name, r.role_code
          FROM roles r
          INNER JOIN user_roles ur ON r.id = ur.role_id
          WHERE ur.user_id = ?
        `, [user.id]);
        user.roles = roles;
      }

      res.json({
        success: true,
        data: {
          list: users,
          total,
          page: parseInt(page),
          pageSize: parseInt(pageSize)
        }
      });
    } catch (error) {
      logger.error('获取用户列表失败:', error);
      res.status(500).json({
        success: false,
        message: '获取用户列表失败'
      });
    }
  }

  async getUserById(req, res) {
    try {
      const { id } = req.params;

      const [users] = await pool.query(
        'SELECT id, username, email, phone, real_name, status, last_login_time, last_login_ip, created_at, updated_at FROM users WHERE id = ?',
        [id]
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
        WHERE ur.user_id = ?
      `, [id]);

      users[0].roles = roles;

      res.json({
        success: true,
        data: users[0]
      });
    } catch (error) {
      logger.error('获取用户详情失败:', error);
      res.status(500).json({
        success: false,
        message: '获取用户详情失败'
      });
    }
  }

  async createUser(req, res) {
    try {
      const { username, password, email, phone, realName, roleIds, status = 1 } = req.body;

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
        'INSERT INTO users (username, password, email, phone, real_name, status, created_by) VALUES (?, ?, ?, ?, ?, ?, ?)',
        [username, hashedPassword, email, phone, realName, status, req.user.id]
      );

      if (roleIds && roleIds.length > 0) {
        const roleValues = roleIds.map(roleId => [result.insertId, roleId]);
        await pool.query(
          'INSERT INTO user_roles (user_id, role_id) VALUES ?',
          [roleValues]
        );
      }

      logger.info(`创建用户成功: ${username}, 操作人: ${req.user.username}`);

      res.status(201).json({
        success: true,
        message: '创建用户成功',
        data: {
          id: result.insertId,
          username
        }
      });
    } catch (error) {
      logger.error('创建用户失败:', error);
      res.status(500).json({
        success: false,
        message: '创建用户失败'
      });
    }
  }

  async updateUser(req, res) {
    try {
      const { id } = req.params;
      const { email, phone, realName, roleIds, status } = req.body;

      const [users] = await pool.query('SELECT id FROM users WHERE id = ?', [id]);

      if (users.length === 0) {
        return res.status(404).json({
          success: false,
          message: '用户不存在'
        });
      }

      await pool.query(
        'UPDATE users SET email = ?, phone = ?, real_name = ?, status = ?, updated_by = ? WHERE id = ?',
        [email, phone, realName, status, req.user.id, id]
      );

      if (roleIds !== undefined) {
        await pool.query('DELETE FROM user_roles WHERE user_id = ?', [id]);

        if (roleIds.length > 0) {
          const roleValues = roleIds.map(roleId => [id, roleId]);
          await pool.query(
            'INSERT INTO user_roles (user_id, role_id) VALUES ?',
            [roleValues]
          );
        }
      }

      logger.info(`更新用户成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '更新用户成功'
      });
    } catch (error) {
      logger.error('更新用户失败:', error);
      res.status(500).json({
        success: false,
        message: '更新用户失败'
      });
    }
  }

  async deleteUser(req, res) {
    try {
      const { id } = req.params;

      if (id === req.user.id) {
        return res.status(400).json({
          success: false,
          message: '不能删除当前登录用户'
        });
      }

      const [users] = await pool.query('SELECT id FROM users WHERE id = ?', [id]);

      if (users.length === 0) {
        return res.status(404).json({
          success: false,
          message: '用户不存在'
        });
      }

      await pool.query('DELETE FROM users WHERE id = ?', [id]);

      logger.info(`删除用户成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '删除用户成功'
      });
    } catch (error) {
      logger.error('删除用户失败:', error);
      res.status(500).json({
        success: false,
        message: '删除用户失败'
      });
    }
  }

  async updateUserStatus(req, res) {
    try {
      const { id } = req.params;
      const { status } = req.body;

      if (id === req.user.id) {
        return res.status(400).json({
          success: false,
          message: '不能修改当前登录用户状态'
        });
      }

      const [users] = await pool.query('SELECT id FROM users WHERE id = ?', [id]);

      if (users.length === 0) {
        return res.status(404).json({
          success: false,
          message: '用户不存在'
        });
      }

      await pool.query(
        'UPDATE users SET status = ?, updated_by = ? WHERE id = ?',
        [status, req.user.id, id]
      );

      logger.info(`更新用户状态成功: ID=${id}, status=${status}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '更新用户状态成功'
      });
    } catch (error) {
      logger.error('更新用户状态失败:', error);
      res.status(500).json({
        success: false,
        message: '更新用户状态失败'
      });
    }
  }

  async resetPassword(req, res) {
    try {
      const { id } = req.params;
      const { newPassword } = req.body;

      const [users] = await pool.query('SELECT id FROM users WHERE id = ?', [id]);

      if (users.length === 0) {
        return res.status(404).json({
          success: false,
          message: '用户不存在'
        });
      }

      const hashedPassword = await hashPassword(newPassword);

      await pool.query(
        'UPDATE users SET password = ?, updated_by = ? WHERE id = ?',
        [hashedPassword, req.user.id, id]
      );

      logger.info(`重置用户密码成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '重置密码成功'
      });
    } catch (error) {
      logger.error('重置用户密码失败:', error);
      res.status(500).json({
        success: false,
        message: '重置密码失败'
      });
    }
  }
}

module.exports = new UserController();
