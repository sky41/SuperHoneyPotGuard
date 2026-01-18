<template>
  <a-card title="角色管理">
    <template #extra>
      <a-button type="primary" @click="handleCreate">
        <template #icon>
          <PlusOutlined />
        </template>
        新增角色
      </a-button>
    </template>

    <a-table
      :columns="columns"
      :data-source="roles"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
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
            <a-button type="link" size="small" @click="handleAssignPermissions(record)">
              <SafetyOutlined />
              分配权限
            </a-button>
            <a-popconfirm
              v-if="record.status === 1"
              title="确定要禁用该角色吗？"
              @confirm="() => handleStatusChange(record.id, 0)"
            >
              <a-button type="link" size="small" danger>
                <StopOutlined />
                禁用
              </a-button>
            </a-popconfirm>
            <a-popconfirm
              v-else
              title="确定要启用该角色吗？"
              @confirm="() => handleStatusChange(record.id, 1)"
            >
              <a-button type="link" size="small" style="color: #52c41a">
                <CheckCircleOutlined />
                启用
              </a-button>
            </a-popconfirm>
            <a-popconfirm
              title="确定要删除该角色吗？"
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
      :title="modalType === 'create' ? '新增角色' : '编辑角色'"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      width="600px"
    >
      <a-form :model="formState" layout="vertical">
        <a-form-item
          label="角色名称"
          name="roleName"
          :rules="[{ required: true, message: '请输入角色名称' }]"
        >
          <a-input v-model:value="formState.roleName" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item
          label="角色编码"
          name="roleCode"
          :rules="[{ required: true, message: '请输入角色编码' }]"
        >
          <a-input v-model:value="formState.roleCode" placeholder="请输入角色编码" :disabled="modalType === 'edit'" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="formState.description" placeholder="请输入描述" :rows="4" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="formState.status" placeholder="请选择状态">
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">禁用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:open="permissionModalVisible"
      title="分配权限"
      @ok="handlePermissionModalOk"
      @cancel="handlePermissionModalCancel"
      width="600px"
    >
      <a-tree
        v-model:checkedKeys="checkedPermissions"
        checkable
        :tree-data="permissionTree"
        :field-names="{ children: 'children', title: 'permissionName', key: 'id' }"
      />
    </a-modal>
  </a-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  SafetyOutlined,
  StopOutlined,
  CheckCircleOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import { roleAPI, permissionAPI } from '@/api'

const loading = ref(false)
const roles = ref([])
const modalVisible = ref(false)
const modalType = ref('create')
const formState = ref({
  roleName: '',
  roleCode: '',
  description: '',
  status: 1
})

const permissionModalVisible = ref(false)
const permissionTree = ref([])
const checkedPermissions = ref([])
const currentRoleId = ref(null)

const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '角色名称', dataIndex: 'roleName', key: 'roleName' },
  { title: '角色编码', dataIndex: 'roleCode', key: 'roleCode' },
  { title: '描述', dataIndex: 'description', key: 'description' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  { title: '操作', key: 'action', width: 350 }
]

onMounted(() => {
  fetchRoles()
})

const fetchRoles = async (params = {}) => {
  loading.value = true
  try {
    const res = await roleAPI.getList({
      page: pagination.value.current,
      pageSize: pagination.value.pageSize,
      ...params
    })
    roles.value = res.data.list
    pagination.value.total = res.data.total
  } catch (error) {
    console.error('获取角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchPermissions = async () => {
  try {
    const res = await permissionAPI.getTree()
    permissionTree.value = res.data
  } catch (error) {
    console.error('获取权限树失败:', error)
  }
}

const handleTableChange = (pag) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchRoles()
}

const handleCreate = () => {
  modalType.value = 'create'
  formState.value = {
    roleName: '',
    roleCode: '',
    description: '',
    status: 1
  }
  modalVisible.value = true
}

const handleEdit = (record) => {
  modalType.value = 'edit'
  formState.value = {
    roleName: record.roleName,
    roleCode: record.roleCode,
    description: record.description,
    status: record.status
  }
  modalVisible.value = true
}

const handleModalOk = async () => {
  try {
    if (modalType.value === 'create') {
      await roleAPI.create(formState.value)
      message.success('创建成功')
    } else {
      await roleAPI.update(currentRole.value.id, formState.value)
      message.success('更新成功')
    }
    modalVisible.value = false
    fetchRoles()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

const handleModalCancel = () => {
  modalVisible.value = false
}

const handleDelete = async (id) => {
  try {
    await roleAPI.delete(id)
    message.success('删除成功')
    fetchRoles()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleStatusChange = async (id, status) => {
  try {
    await roleAPI.update(id, { status })
    message.success('状态更新成功')
    fetchRoles()
  } catch (error) {
    console.error('状态更新失败:', error)
  }
}

const handleAssignPermissions = async (record) => {
  currentRoleId.value = record.id
  await fetchPermissions()
  checkedPermissions.value = record.permissions?.map(p => p.id) || []
  permissionModalVisible.value = true
}

const handlePermissionModalOk = async () => {
  try {
    await roleAPI.update(currentRoleId.value, { permissionIds: checkedPermissions.value })
    message.success('权限分配成功')
    permissionModalVisible.value = false
    fetchRoles()
  } catch (error) {
    console.error('权限分配失败:', error)
  }
}

const handlePermissionModalCancel = () => {
  permissionModalVisible.value = false
}

const currentRole = ref(null)
</script>

<style scoped>
</style>
