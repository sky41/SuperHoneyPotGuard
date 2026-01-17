<template>
  <a-layout style="min-height: 100vh">
    <a-layout-sider
      v-model:collapsed="collapsed"
      collapsible
      :trigger="null"
    >
      <div class="logo">
        {{ collapsed ? 'SHG' : 'SuperHoneyPotGuard' }}
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        theme="dark"
        @click="handleMenuClick"
      >
        <a-menu-item key="/">
          <template #icon>
            <DashboardOutlined />
          </template>
          首页
        </a-menu-item>
        <a-menu-item key="/user">
          <template #icon>
            <UserOutlined />
          </template>
          用户管理
        </a-menu-item>
        <a-menu-item key="/role">
          <template #icon>
            <TeamOutlined />
          </template>
          角色管理
        </a-menu-item>
        <a-menu-item key="/permission">
          <template #icon>
            <SafetyOutlined />
          </template>
          权限管理
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-title">{{ pageTitle }}</div>
        <a-dropdown>
          <div class="user-info">
            <a-avatar>
              <UserOutlined />
            </a-avatar>
            <span class="username">{{ user?.username || '用户' }}</span>
          </div>
          <template #overlay>
            <a-menu @click="handleUserMenuClick">
              <a-menu-item key="logout">
                <LogoutOutlined />
                退出登录
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </a-layout-header>

      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  DashboardOutlined,
  UserOutlined,
  TeamOutlined,
  SafetyOutlined,
  LogoutOutlined
} from '@ant-design/icons-vue'
import { authAPI } from '@/api'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)
const selectedKeys = ref([route.path])
const user = ref(null)

const pageTitle = computed(() => {
  const titles = {
    '/': '首页',
    '/user': '用户管理',
    '/role': '角色管理',
    '/permission': '权限管理'
  }
  return titles[route.path] || '首页'
})

onMounted(async () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
})

const handleMenuClick = ({ key }) => {
  selectedKeys.value = [key]
  router.push(key)
}

const handleUserMenuClick = async ({ key }) => {
  if (key === 'logout') {
    try {
      await authAPI.logout()
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      message.success('注销成功')
      router.push('/login')
    } catch (error) {
      console.error('注销失败:', error)
    }
  }
}
</script>

<style scoped>
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: bold;
  background: #001529;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-title {
  font-size: 18px;
  font-weight: 500;
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-size: 14px;
}

.content {
  margin: 24px;
  padding: 24px;
  background: #fff;
  border-radius: 8px;
}
</style>
