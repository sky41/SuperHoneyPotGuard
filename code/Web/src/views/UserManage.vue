<template>
  <a-card title="用户管理">
    <template #extra>
      <a-button type="primary" @click="handleCreate">
        <template #icon>
          <PlusOutlined />
        </template>
        新增用户
      </a-button>
    </template>

    <a-table
      :columns="columns"
      :data-source="users"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'roles'">
          <a-tag v-for="role in record.roles" :key="role.id" color="blue">
            {{ role.roleName }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleEdit(record)">
              <EditOutlined />
              编辑
            </a-button>
            <a-button type="link" size="small" @click="handleResetPassword(record.id)">
              <LockOutlined />
              重置密码
            </a-button>
            <a-popconfirm
              v-if="record.status === 1"
              title="确定要禁用该用户吗？"
              @confirm="() => handleStatusChange(record.id, 0)"
            >
              <a-button type="link" size="small" danger>
                <StopOutlined />
                禁用
              </a-button>
            </a-popconfirm>
            <a-popconfirm
              v-else
              title="确定要启用该用户吗？"
              @confirm="() => handleStatusChange(record.id, 1)"
            >
              <a-button type="link" size="small" style="color: #52c41a">
                <CheckCircleOutlined />
                启用
              </a-button>
            </a-popconfirm>
            <a-popconfirm
              title="确定要删除该用户吗？"
              @confirm="() => handleDelete(record.id)"
            >
              <a-button type="link" size="small" danger>
                <DeleteOutlined />
                删除
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="modalType === 'create' ? '新增用户' : '编辑用户'"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      width="600px"
    >
      <a-form :model="formState" layout="vertical">
        <a-form-item
          label="用户名"
          name="username"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <a-input v-model:value="formState.username" placeholder="请输入用户名" :disabled="modalType === 'edit'" />
        </a-form-item>
        <a-form-item
          v-if="modalType === 'create'"
          label="密码"
          name="password"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password v-model:value="formState.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item
          label="邮箱"
          name="email"
          :rules="[{ type: 'email', message: '请输入有效的邮箱地址' }]"
        >
          <a-input v-model:value="formState.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model:value="formState.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="真实姓名">
          <a-input v-model:value="formState.realName" placeholder="请输入真实姓名" />
        </a-form-item>
        <a-form-item
          label="角色"
          name="roleIds"
          :rules="[{ required: true, message: '请选择角色' }]"
        >
          <a-select
            v-model:value="formState.roleIds"
            mode="multiple"
            placeholder="请选择角色"
            :options="roleOptions"
          />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="formState.status" placeholder="请选择状态">
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">禁用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  LockOutlined,
  StopOutlined,
  CheckCircleOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import { userAPI, roleAPI } from '@/api'

const loading = ref(false)
const users = ref([])
const roles = ref([])
const modalVisible = ref(false)
const modalType = ref('create')
const formState = ref({
  username: '',
  password: '',
  email: '',
  phone: '',
  realName: '',
  roleIds: [],
  status: 1
})

const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username', key: 'username' },
  { title: '邮箱', dataIndex: 'email', key: 'email' },
  { title: '真实姓名', dataIndex: 'realName', key: 'realName' },
  { title: '角色', dataIndex: 'roles', key: 'roles' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '最后登录时间', dataIndex: 'lastLoginTime', key: 'lastLoginTime' },
  { title: '操作', key: 'action', width: 300 }
]

const roleOptions = computed(() => {
  return roles.value.map(role => ({
    label: role.roleName,
    value: role.id
  }))
})

onMounted(() => {
  fetchUsers()
  fetchRoles()
})

const fetchUsers = async (params = {}) => {
  loading.value = true
  try {
    const res = await userAPI.getList({
      page: pagination.value.current,
      pageSize: pagination.value.pageSize,
      ...params
    })
    users.value = res.data.list
    pagination.value.total = res.data.total
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const res = await roleAPI.getAll()
    roles.value = res.data
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

const handleTableChange = (pag) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchUsers()
}

const handleCreate = () => {
  modalType.value = 'create'
  formState.value = {
    username: '',
    password: '',
    email: '',
    phone: '',
    realName: '',
    roleIds: [],
    status: 1
  }
  modalVisible.value = true
}

const handleEdit = (record) => {
  modalType.value = 'edit'
  formState.value = {
    username: record.username,
    email: record.email,
    phone: record.phone,
    realName: record.realName,
    roleIds: record.roles?.map(r => r.id) || [],
    status: record.status
  }
  modalVisible.value = true
}

const handleModalOk = async () => {
  try {
    if (modalType.value === 'create') {
      await userAPI.create(formState.value)
      message.success('创建成功')
    } else {
      await userAPI.update(currentUser.value.id, formState.value)
      message.success('更新成功')
    }
    modalVisible.value = false
    fetchUsers()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

const handleModalCancel = () => {
  modalVisible.value = false
}

const handleDelete = async (id) => {
  try {
    await userAPI.delete(id)
    message.success('删除成功')
    fetchUsers()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleStatusChange = async (id, status) => {
  try {
    await userAPI.updateStatus(id, status)
    message.success('状态更新成功')
    fetchUsers()
  } catch (error) {
    console.error('状态更新失败:', error)
  }
}

const handleResetPassword = (id) => {
  const newPassword = ref('')
  Modal.confirm({
    title: '重置密码',
    content: () => h('div', [
      h('a-input-password', {
        placeholder: '请输入新密码',
        style: { marginBottom: '16px' },
        onChange: (e) => {
          newPassword.value = e.target.value
        }
      })
    ]),
    onOk: async () => {
      try {
        await userAPI.resetPassword(id, newPassword.value)
        message.success('密码重置成功')
      } catch (error) {
        console.error('密码重置失败:', error)
      }
    }
  })
}

const currentUser = ref(null)
</script>

<style scoped>
</style>
