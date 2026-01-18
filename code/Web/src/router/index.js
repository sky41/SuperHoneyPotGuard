import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Dashboard from '@/views/Dashboard.vue'
import UserManage from '@/views/UserManage.vue'
import RoleManage from '@/views/RoleManage.vue'
import PermissionManage from '@/views/PermissionManage.vue'
import LogManage from '@/views/LogManage.vue'
import HFishData from '@/views/HFishData.vue'
import ResetPassword from '@/views/ResetPassword.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: 'user',
        name: 'UserManage',
        component: UserManage
      },
      {
        path: 'role',
        name: 'RoleManage',
        component: RoleManage
      },
      {
        path: 'permission',
        name: 'PermissionManage',
        component: PermissionManage
      },
      {
        path: 'log',
        name: 'LogManage',
        component: LogManage
      },
      {
        path: 'hfish',
        name: 'HFishData',
        component: HFishData
      }
    ]
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: ResetPassword
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && to.path !== '/reset-password' && !token) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/reset-password') && token) {
    next('/')
  } else {
    next()
  }
})

export default router
