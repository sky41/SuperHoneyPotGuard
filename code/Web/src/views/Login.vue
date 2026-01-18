<template>
  <div class="login-container">
    <a-card :title="isLogin ? 'SuperHoneyPotGuard' : '注册新用户'" class="login-card">
      <a-form
        :model="formState"
        @finish="handleFinish"
        layout="vertical"
      >
        <a-form-item
          label="用户名"
          name="username"
          :rules="[{ required: true, message: '请输入用户名' }, { min: 3, max: 50, message: '用户名长度为3-50个字符' }]"
        >
          <a-input
            v-model:value="formState.username"
            placeholder="请输入用户名"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          v-if="!isLogin"
          label="邮箱"
          name="email"
          :rules="[{ type: 'email', message: '请输入有效的邮箱地址' }]"
        >
          <a-input
            v-model:value="formState.email"
            placeholder="请输入邮箱"
          >
            <template #suffix>
              <a-button
                type="primary"
                :disabled="countdown > 0"
                @click="sendVerificationCode"
                :loading="sendingCode"
                size="small"
              >
                {{ countdown > 0 ? `${countdown}秒后重发` : '发送验证码' }}
              </a-button>
            </template>
            <template #prefix>
              <MailOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          v-if="!isLogin"
          label="验证码"
          name="code"
          :rules="[{ required: true, message: '请输入验证码' }]"
        >
          <a-input
            v-model:value="formState.code"
            placeholder="请输入6位验证码"
            maxlength="6"
          >
            <template #prefix>
              <SafetyOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          label="密码"
          name="password"
          :rules="[
            { required: true, message: '请输入密码' },
            { min: 6, message: '密码至少6个字符' }
          ]"
        >
          <a-input-password
            v-model:value="formState.password"
            placeholder="请输入密码"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item
          v-if="!isLogin"
          label="确认密码"
          name="confirmPassword"
          :rules="[
            { required: true, message: '请确认密码' },
            { validator: validatePassword }
          ]"
        >
          <a-input-password
            v-model:value="formState.confirmPassword"
            placeholder="请再次输入密码"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item
          v-if="!isLogin"
          label="真实姓名"
          name="realName"
        >
          <a-input
            v-model:value="formState.realName"
            placeholder="请输入真实姓名（可选）"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            block
            :loading="loading"
          >
            {{ isLogin ? '登录' : '注册' }}
          </a-button>
        </a-form-item>

        <a-form-item>
          <a-button
            type="link"
            block
            @click="toggleMode"
          >
            {{ isLogin ? '还没有账号？立即注册' : '已有账号？立即登录' }}
          </a-button>
        </a-form-item>

        <a-form-item v-if="isLogin">
          <a-button
            type="link"
            block
            @click="goToResetPassword"
          >
            忘记密码？
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined, MailOutlined, SafetyOutlined } from '@ant-design/icons-vue'
import { authAPI } from '@/api'

const router = useRouter()
const loading = ref(false)
const sendingCode = ref(false)
const isLogin = ref(true)
const countdown = ref(0)
let countdownTimer = null

const formState = ref({
  username: '',
  password: '',
  email: '',
  code: '',
  confirmPassword: '',
  realName: ''
})

const toggleMode = () => {
  isLogin.value = !isLogin.value
  formState.value = {
    username: '',
    password: '',
    email: '',
    code: '',
    confirmPassword: '',
    realName: ''
  }
  resetCountdown()
}

const resetCountdown = () => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
  countdown.value = 0
}

const sendVerificationCode = async () => {
  if (!formState.value.email) {
    message.error('请先输入邮箱地址')
    return
  }

  sendingCode.value = true
  try {
    await authAPI.sendVerificationCode({ email: formState.value.email })
    message.success('验证码已发送到您的邮箱，有效期5分钟')
    startCountdown()
  } catch (error) {
    console.error('发送验证码失败:', error)
    message.error('发送验证码失败，请重试')
  } finally {
    sendingCode.value = false
  }
}

const startCountdown = () => {
  countdown.value = 60
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      resetCountdown()
    }
  }, 1000)
}

const validatePassword = async (rule, value) => {
  if (value !== formState.value.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const handleFinish = async () => {
  loading.value = true
  try {
    if (isLogin.value) {
      const res = await authAPI.login({
        username: formState.value.username,
        password: formState.value.password
      })
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      message.success('登录成功')
      router.push('/')
    } else {
      await authAPI.register({
        username: formState.value.username,
        password: formState.value.password,
        email: formState.value.email,
        code: formState.value.code,
        realName: formState.value.realName || undefined
      })
      message.success('注册成功，请登录')
      isLogin.value = true
      formState.value = {
        username: '',
        password: '',
        email: '',
        code: '',
        confirmPassword: '',
        realName: ''
      }
    }
  } catch (error) {
    console.error(isLogin.value ? '登录失败:' : '注册失败:', error)
    message.error(isLogin.value ? '登录失败' : '注册失败')
  } finally {
    loading.value = false
  }
}

const goToResetPassword = () => {
  router.push('/reset-password')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
