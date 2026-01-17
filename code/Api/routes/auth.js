const express = require('express');
const { body, validationResult } = require('express-validator');
const authController = require('../controllers/authController');
const authMiddleware = require('../middleware/auth');
const { authLimiter } = require('../middleware/rateLimit');

const router = express.Router();

router.post('/register', [
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
}, authController.register);

router.post('/login', authLimiter, [
  body('username').notEmpty().withMessage('用户名不能为空'),
  body('password').notEmpty().withMessage('密码不能为空')
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
}, authController.login);

router.post('/logout', authMiddleware, authController.logout);

router.get('/current', authMiddleware, authController.getCurrentUser);

module.exports = router;
