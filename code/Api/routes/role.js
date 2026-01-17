const express = require('express');
const { body, validationResult } = require('express-validator');
const roleController = require('../controllers/roleController');
const authMiddleware = require('../middleware/auth');
const checkPermission = require('../middleware/permission');

const router = express.Router();

router.use(authMiddleware);

router.get('/list', checkPermission('role:list'), roleController.getRoleList);

router.get('/all', checkPermission('role:list'), roleController.getAllRoles);

router.get('/:id', checkPermission('role:list'), roleController.getRoleById);

router.post('/', checkPermission('role:add'), [
  body('roleName').notEmpty().withMessage('角色名称不能为空'),
  body('roleCode').notEmpty().withMessage('角色编码不能为空')
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
}, roleController.createRole);

router.put('/:id', checkPermission('role:edit'), roleController.updateRole);

router.delete('/:id', checkPermission('role:delete'), roleController.deleteRole);

module.exports = router;
