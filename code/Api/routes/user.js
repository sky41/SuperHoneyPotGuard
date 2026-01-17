const express = require('express');
const { body, validationResult } = require('express-validator');
const userController = require('../controllers/userController');
const authMiddleware = require('../middleware/auth');
const checkPermission = require('../middleware/permission');

const router = express.Router();

router.use(authMiddleware);

router.get('/list', checkPermission('user:list'), userController.getUserList);

router.get('/:id', checkPermission('user:list'), userController.getUserById);

router.post('/', checkPermission('user:add'), [
  body('username').notEmpty().withMessage('用户名不能为空').isLength({ min: 3, max: 50 }).withMessage('用户名长度为3-50个字符'),
  body('password').notEmpty().withMessage('密码不能为空').isLength({ min: 6 }).withMessage('密码长度至少6个字符'),
  body('email').optional().isEmail().withMessage('邮箱格式不正确'),
  body('phone').optional().isMobilePhone('zh-CN').withMessage('手机号格式不正确')
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
}, userController.createUser);

router.put('/:id', checkPermission('user:edit'), userController.updateUser);

router.delete('/:id', checkPermission('user:delete'), userController.deleteUser);

router.patch('/:id/status', checkPermission('user:edit'), [
  body('status').isIn([0, 1]).withMessage('状态值不正确')
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
}, userController.updateUserStatus);

router.post('/:id/reset-password', checkPermission('user:edit'), [
  body('newPassword').notEmpty().withMessage('新密码不能为空').isLength({ min: 6 }).withMessage('密码长度至少6个字符')
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
}, userController.resetPassword);

module.exports = router;
