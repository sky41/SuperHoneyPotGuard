const pool = require('../config/database');
const logger = require('../utils/logger');

class RoleController {
  async getRoleList(req, res) {
    try {
      const { page = 1, pageSize = 10, roleName, status } = req.query;
      const offset = (page - 1) * pageSize;

      let whereClause = 'WHERE 1=1';
      const params = [];

      if (roleName) {
        whereClause += ' AND role_name LIKE ?';
        params.push(`%${roleName}%`);
      }

      if (status !== undefined) {
        whereClause += ' AND status = ?';
        params.push(status);
      }

      const [roles] = await pool.query(`
        SELECT id, role_name, role_code, description, status, created_at, updated_at
        FROM roles
        ${whereClause}
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
      `, [...params, parseInt(pageSize), offset]);

      const [countResult] = await pool.query(`
        SELECT COUNT(*) as total
        FROM roles
        ${whereClause}
      `, params);

      const total = countResult[0].total;

      for (const role of roles) {
        const [permissions] = await pool.query(`
          SELECT p.id, p.permission_name, p.permission_code, p.permission_type
          FROM permissions p
          INNER JOIN role_permissions rp ON p.id = rp.permission_id
          WHERE rp.role_id = ?
        `, [role.id]);
        role.permissions = permissions;
      }

      res.json({
        success: true,
        data: {
          list: roles,
          total,
          page: parseInt(page),
          pageSize: parseInt(pageSize)
        }
      });
    } catch (error) {
      logger.error('获取角色列表失败:', error);
      res.status(500).json({
        success: false,
        message: '获取角色列表失败'
      });
    }
  }

  async getAllRoles(req, res) {
    try {
      const [roles] = await pool.query(`
        SELECT id, role_name, role_code, description, status
        FROM roles
        WHERE status = 1
        ORDER BY id ASC
      `);

      res.json({
        success: true,
        data: roles
      });
    } catch (error) {
      logger.error('获取所有角色失败:', error);
      res.status(500).json({
        success: false,
        message: '获取所有角色失败'
      });
    }
  }

  async getRoleById(req, res) {
    try {
      const { id } = req.params;

      const [roles] = await pool.query(
        'SELECT id, role_name, role_code, description, status, created_at, updated_at FROM roles WHERE id = ?',
        [id]
      );

      if (roles.length === 0) {
        return res.status(404).json({
          success: false,
          message: '角色不存在'
        });
      }

      const [permissions] = await pool.query(`
        SELECT p.id, p.permission_name, p.permission_code, p.permission_type
        FROM permissions p
        INNER JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = ?
      `, [id]);

      roles[0].permissions = permissions;

      res.json({
        success: true,
        data: roles[0]
      });
    } catch (error) {
      logger.error('获取角色详情失败:', error);
      res.status(500).json({
        success: false,
        message: '获取角色详情失败'
      });
    }
  }

  async createRole(req, res) {
    try {
      const { roleName, roleCode, description, permissionIds, status = 1 } = req.body;

      const [existingRoles] = await pool.query(
        'SELECT id FROM roles WHERE role_name = ? OR role_code = ?',
        [roleName, roleCode]
      );

      if (existingRoles.length > 0) {
        return res.status(400).json({
          success: false,
          message: '角色名称或角色编码已存在'
        });
      }

      const [result] = await pool.query(
        'INSERT INTO roles (role_name, role_code, description, status, created_by) VALUES (?, ?, ?, ?, ?)',
        [roleName, roleCode, description, status, req.user.id]
      );

      if (permissionIds && permissionIds.length > 0) {
        const permissionValues = permissionIds.map(permissionId => [result.insertId, permissionId]);
        await pool.query(
          'INSERT INTO role_permissions (role_id, permission_id) VALUES ?',
          [permissionValues]
        );
      }

      logger.info(`创建角色成功: ${roleName}, 操作人: ${req.user.username}`);

      res.status(201).json({
        success: true,
        message: '创建角色成功',
        data: {
          id: result.insertId,
          roleName
        }
      });
    } catch (error) {
      logger.error('创建角色失败:', error);
      res.status(500).json({
        success: false,
        message: '创建角色失败'
      });
    }
  }

  async updateRole(req, res) {
    try {
      const { id } = req.params;
      const { roleName, description, permissionIds, status } = req.body;

      const [roles] = await pool.query('SELECT id FROM roles WHERE id = ?', [id]);

      if (roles.length === 0) {
        return res.status(404).json({
          success: false,
          message: '角色不存在'
        });
      }

      await pool.query(
        'UPDATE roles SET role_name = ?, description = ?, status = ?, updated_by = ? WHERE id = ?',
        [roleName, description, status, req.user.id, id]
      );

      if (permissionIds !== undefined) {
        await pool.query('DELETE FROM role_permissions WHERE role_id = ?', [id]);

        if (permissionIds.length > 0) {
          const permissionValues = permissionIds.map(permissionId => [id, permissionId]);
          await pool.query(
            'INSERT INTO role_permissions (role_id, permission_id) VALUES ?',
            [permissionValues]
          );
        }
      }

      logger.info(`更新角色成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '更新角色成功'
      });
    } catch (error) {
      logger.error('更新角色失败:', error);
      res.status(500).json({
        success: false,
        message: '更新角色失败'
      });
    }
  }

  async deleteRole(req, res) {
    try {
      const { id } = req.params;

      const [users] = await pool.query('SELECT COUNT(*) as count FROM user_roles WHERE role_id = ?', [id]);

      if (users[0].count > 0) {
        return res.status(400).json({
          success: false,
          message: '该角色下还有用户，无法删除'
        });
      }

      const [roles] = await pool.query('SELECT id FROM roles WHERE id = ?', [id]);

      if (roles.length === 0) {
        return res.status(404).json({
          success: false,
          message: '角色不存在'
        });
      }

      await pool.query('DELETE FROM roles WHERE id = ?', [id]);

      logger.info(`删除角色成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '删除角色成功'
      });
    } catch (error) {
      logger.error('删除角色失败:', error);
      res.status(500).json({
        success: false,
        message: '删除角色失败'
      });
    }
  }
}

module.exports = new RoleController();
