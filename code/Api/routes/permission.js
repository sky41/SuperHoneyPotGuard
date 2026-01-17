const express = require('express');
const { body, validationResult } = require('express-validator');
const permissionController = require('../controllers/permissionController');
const authMiddleware = require('../middleware/auth');
const checkPermission = require('../middleware/permission');

const router = express.Router();

router.use(authMiddleware);

router.get('/tree', checkPermission('permission:manage'), permissionController.getPermissionTree);

router.get('/all', checkPermission('permission:manage'), permissionController.getAllPermissions);

router.get('/:id', checkPermission('permission:manage'), permissionController.getPermissionById);

router.post('/', checkPermission('permission:manage'), [
  body('permissionName').notEmpty().withMessage('权限名称不能为空'),
  body('permissionCode').notEmpty().withMessage('权限编码不能为空'),
  body('permissionType').notEmpty().withMessage('权限类型不能为空')
], (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({
      success: false,
      message: '参数验证失败',
      errors: errors.array()
    });
  }
  next();
}, permissionController.createPermission);

router.put('/:id', checkPermission('permission:manage'), permissionController.updatePermission);

router.delete('/:id', checkPermission('permission:manage'), permissionController.deletePermission);

module.exports = router;
