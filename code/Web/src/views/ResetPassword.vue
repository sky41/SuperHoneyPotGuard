<template>
  <div class="reset-password-container">
    <a-card title="重置密码" class="reset-card">
      <a-form
        :model="formState"
        @finish="handleFinish"
        layout="vertical"
      >
        <a-form-item
          label="邮箱"
          name="email"
          :rules="[{ type: 'email', message: '请输入有效的邮箱地址' }]"
        >
          <a-input
            v-model:value="formState.email"
            placeholder="请输入邮箱"
            :disabled="step > 1"
          >
            <template #prefix>
              <MailOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          v-if="step === 1"
          label="验证码"
          name="code"
          :rules="[{ required: true, message: '请输入验证码' }]"
        >
          <a-input
            v-model:value="formState.code"
            placeholder="请输入6位验证码"
            maxlength="6"
          >
            <template #suffix>
              <a-button
                type="primary"
                :disabled="countdown > 0"
                :loading="sendingCode"
                @click="sendResetCode"
                size="small"
              >
                {{ countdown > 0 ? `${countdown}秒后重发` : '发送验证码' }}
              </a-button>
            </template>
            <template #prefix>
              <SafetyOutlined />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item
          v-if="step === 2"
          label="新密码"
          name="newPassword"
          :rules="[
            { required: true, message: '请输入新密码' },
            { min: 6, message: '密码至少6个字符' }
          ]"
        >
          <a-input-password
            v-model:value="formState.newPassword"
            placeholder="请输入新密码（至少6个字符）"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item
          v-if="step === 2"
          label="确认密码"
          name="confirmPassword"
          :rules="[
            { required: true, message: '请确认新密码' },
            { validator: validatePassword }
          ]"
        >
          <a-input-password
            v-model:value="formState.confirmPassword"
            placeholder="请再次输入新密码"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            v-if="step === 1"
            type="primary"
            html-type="submit"
            block
            :loading="loading"
            @click="nextStep"
          >
            下一步
          </a-button>
          <a-button
            v-if="step === 2"
            type="primary"
            html-type="submit"
            block
            :loading="loading"
          >
            重置密码
          </a-button>
          <a-button
            type="link"
            block
            @click="goBack"
          >
            返回登录
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
import { MailOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons-vue'
import { authAPI } from '@/api'

const router = useRouter()
const loading = ref(false)
const sendingCode = ref(false)
const step = ref(1)
const countdown = ref(0)
let countdownTimer = null

const formState = ref({
  email: '',
  code: '',
  newPassword: '',
  confirmPassword: ''
})

const resetCountdown = () => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
  countdown.value = 0
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

const sendResetCode = async () => {
  if (!formState.value.email) {
    message.error('请先输入邮箱地址')
    return
  }

  sendingCode.value = true
  try {
    await authAPI.sendResetPasswordCode({ email: formState.value.email })
    message.success('验证码已发送到您的邮箱，有效期5分钟')
    startCountdown()
  } catch (error) {
    console.error('发送验证码失败:', error)
    message.error('发送验证码失败，请重试')
  } finally {
    sendingCode.value = false
  }
}

const validatePassword = async (rule, value) => {
  if (value !== formState.value.newPassword) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const nextStep = () => {
  if (!formState.value.code) {
    message.error('请输入验证码')
    return
  }
  step.value = 2
}

const handleFinish = async () => {
  loading.value = true
  try {
    await authAPI.resetPassword({
      email: formState.value.email,
      code: formState.value.code,
      newPassword: formState.value.newPassword
    })
    message.success('密码重置成功，请使用新密码登录')
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error) {
    console.error('密码重置失败:', error)
    message.error('密码重置失败')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/login')
}
</script>

<style scoped>
.reset-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.reset-card {
  width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
