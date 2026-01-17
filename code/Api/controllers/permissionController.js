const pool = require('../config/database');
const logger = require('../utils/logger');

class PermissionController {
  async getPermissionTree(req, res) {
    try {
      const [permissions] = await pool.query(`
        SELECT id, permission_name, permission_code, permission_type, parent_id, path, component, icon, sort_order, description, status
        FROM permissions
        ORDER BY sort_order ASC, id ASC
      `);

      const buildTree = (parentId = 0) => {
        return permissions
          .filter(p => p.parent_id === parentId)
          .map(p => ({
            ...p,
            children: buildTree(p.id)
          }));
      };

      const tree = buildTree(0);

      res.json({
        success: true,
        data: tree
      });
    } catch (error) {
      logger.error('获取权限树失败:', error);
      res.status(500).json({
        success: false,
        message: '获取权限树失败'
      });
    }
  }

  async getAllPermissions(req, res) {
    try {
      const [permissions] = await pool.query(`
        SELECT id, permission_name, permission_code, permission_type, parent_id, path, component, icon, sort_order, description, status
        FROM permissions
        WHERE status = 1
        ORDER BY sort_order ASC, id ASC
      `);

      res.json({
        success: true,
        data: permissions
      });
    } catch (error) {
      logger.error('获取所有权限失败:', error);
      res.status(500).json({
        success: false,
        message: '获取所有权限失败'
      });
    }
  }

  async getPermissionById(req, res) {
    try {
      const { id } = req.params;

      const [permissions] = await pool.query(
        'SELECT * FROM permissions WHERE id = ?',
        [id]
      );

      if (permissions.length === 0) {
        return res.status(404).json({
          success: false,
          message: '权限不存在'
        });
      }

      res.json({
        success: true,
        data: permissions[0]
      });
    } catch (error) {
      logger.error('获取权限详情失败:', error);
      res.status(500).json({
        success: false,
        message: '获取权限详情失败'
      });
    }
  }

  async createPermission(req, res) {
    try {
      const { permissionName, permissionCode, permissionType, parentId, path, component, icon, sortOrder, description, status = 1 } = req.body;

      const [existingPermissions] = await pool.query(
        'SELECT id FROM permissions WHERE permission_code = ?',
        [permissionCode]
      );

      if (existingPermissions.length > 0) {
        return res.status(400).json({
          success: false,
          message: '权限编码已存在'
        });
      }

      const [result] = await pool.query(
        'INSERT INTO permissions (permission_name, permission_code, permission_type, parent_id, path, component, icon, sort_order, description, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',
        [permissionName, permissionCode, permissionType, parentId || 0, path, component, icon, sortOrder || 0, description, status]
      );

      logger.info(`创建权限成功: ${permissionName}, 操作人: ${req.user.username}`);

      res.status(201).json({
        success: true,
        message: '创建权限成功',
        data: {
          id: result.insertId,
          permissionName
        }
      });
    } catch (error) {
      logger.error('创建权限失败:', error);
      res.status(500).json({
        success: false,
        message: '创建权限失败'
      });
    }
  }

  async updatePermission(req, res) {
    try {
      const { id } = req.params;
      const { permissionName, permissionType, parentId, path, component, icon, sortOrder, description, status } = req.body;

      const [permissions] = await pool.query('SELECT id FROM permissions WHERE id = ?', [id]);

      if (permissions.length === 0) {
        return res.status(404).json({
          success: false,
          message: '权限不存在'
        });
      }

      await pool.query(
        'UPDATE permissions SET permission_name = ?, permission_type = ?, parent_id = ?, path = ?, component = ?, icon = ?, sort_order = ?, description = ?, status = ? WHERE id = ?',
        [permissionName, permissionType, parentId || 0, path, component, icon, sortOrder || 0, description, status, id]
      );

      logger.info(`更新权限成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '更新权限成功'
      });
    } catch (error) {
      logger.error('更新权限失败:', error);
      res.status(500).json({
        success: false,
        message: '更新权限失败'
      });
    }
  }

  async deletePermission(req, res) {
    try {
      const { id } = req.params;

      const [children] = await pool.query('SELECT COUNT(*) as count FROM permissions WHERE parent_id = ?', [id]);

      if (children[0].count > 0) {
        return res.status(400).json({
          success: false,
          message: '该权限下还有子权限，无法删除'
        });
      }

      const [permissions] = await pool.query('SELECT id FROM permissions WHERE id = ?', [id]);

      if (permissions.length === 0) {
        return res.status(404).json({
          success: false,
          message: '权限不存在'
        });
      }

      await pool.query('DELETE FROM permissions WHERE id = ?', [id]);

      logger.info(`删除权限成功: ID=${id}, 操作人: ${req.user.username}`);

      res.json({
        success: true,
        message: '删除权限成功'
      });
    } catch (error) {
      logger.error('删除权限失败:', error);
      res.status(500).json({
        success: false,
        message: '删除权限失败'
      });
    }
  }
}

module.exports = new PermissionController();
