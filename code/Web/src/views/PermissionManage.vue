<template>
  <a-card title="权限管理">
    <template #extra>
      <a-button type="primary" @click="handleCreate">
        <template #icon>
          <PlusOutlined />
        </template>
        新增权限
      </a-button>
    </template>

    <a-table
      :columns="columns"
      :data-source="permissions"
      :loading="loading"
      :pagination="false"
      row-key="id"
      :default-expand-all-rows="true"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'permissionType'">
          <a-tag :color="getPermissionTypeColor(record.permissionType)">
            {{ getPermissionTypeLabel(record.permissionType) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'icon'">
          <component :is="getIcon(record.icon)" v-if="record.icon" />
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
            <a-button type="link" size="small" @click="handleCreateChild(record)">
              <PlusOutlined />
              新增子权限
            </a-button>
            <a-popconfirm
              title="确定要删除该权限吗？"
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
      :title="modalType === 'create' ? '新增权限' : '编辑权限'"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      width="600px"
    >
      <a-form :model="formState" layout="vertical">
        <a-form-item
          label="权限名称"
          name="permissionName"
          :rules="[{ required: true, message: '请输入权限名称' }]"
        >
          <a-input v-model:value="formState.permissionName" placeholder="请输入权限名称" />
        </a-form-item>
        <a-form-item
          label="权限编码"
          name="permissionCode"
          :rules="[{ required: true, message: '请输入权限编码' }]"
        >
          <a-input v-model:value="formState.permissionCode" placeholder="请输入权限编码" :disabled="modalType === 'edit'" />
        </a-form-item>
        <a-form-item
          label="权限类型"
          name="permissionType"
          :rules="[{ required: true, message: '请选择权限类型' }]"
        >
          <a-select v-model:value="formState.permissionType" placeholder="请选择权限类型">
            <a-select-option value="menu">菜单</a-select-option>
            <a-select-option value="button">按钮</a-select-option>
            <a-select-option value="api">接口</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="父权限">
          <a-tree-select
            v-model:value="formState.parentId"
            :tree-data="permissionTreeOptions"
            placeholder="请选择父权限"
            allow-clear
            tree-default-expand-all
            :field-names="{ children: 'children', title: 'permissionName', value: 'id' }"
          />
        </a-form-item>
        <a-form-item label="路由路径">
          <a-input v-model:value="formState.path" placeholder="请输入路由路径" />
        </a-form-item>
        <a-form-item label="组件路径">
          <a-input v-model:value="formState.component" placeholder="请输入组件路径" />
        </a-form-item>
        <a-form-item label="图标">
          <a-select
            v-model:value="formState.icon"
            placeholder="请选择图标"
            show-search
            :filter-option="filterIconOption"
          >
            <a-select-option v-for="icon in iconOptions" :key="icon" :value="icon">
              <component :is="icon" />
              {{ icon }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="formState.sortOrder" :min="0" placeholder="请输入排序" style="width: 100%" />
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
  </a-card>
</template>

<script setup>
import { ref, onMounted, h } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SettingOutlined,
  UserOutlined,
  TeamOutlined,
  SafetyOutlined,
  FileTextOutlined,
  DashboardOutlined,
  HomeOutlined,
  MenuOutlined,
  AppstoreOutlined,
  FolderOutlined,
  FileOutlined,
  SearchOutlined,
  ReloadOutlined,
  ExportOutlined,
  ImportOutlined,
  DownloadOutlined,
  UploadOutlined,
  PrinterOutlined,
  CopyOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  InfoCircleOutlined,
  WarningOutlined,
  ExclamationCircleOutlined,
  QuestionCircleOutlined,
  LockOutlined,
  UnlockOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
  EditTwoTone,
  DeleteTwoTone,
  PlusCircleOutlined,
  MinusCircleOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  ArrowLeftOutlined,
  ArrowRightOutlined,
  UpCircleOutlined,
  DownCircleOutlined,
  LeftCircleOutlined,
  RightCircleOutlined
} from '@ant-design/icons-vue'
import { permissionAPI } from '@/api'

const loading = ref(false)
const permissions = ref([])
const modalVisible = ref(false)
const modalType = ref('create')
const formState = ref({
  permissionName: '',
  permissionCode: '',
  permissionType: 'menu',
  parentId: null,
  path: '',
  component: '',
  icon: '',
  sortOrder: 0,
  description: '',
  status: 1
})

const iconOptions = [
  'SettingOutlined',
  'UserOutlined',
  'TeamOutlined',
  'SafetyOutlined',
  'FileTextOutlined',
  'DashboardOutlined',
  'HomeOutlined',
  'MenuOutlined',
  'AppstoreOutlined',
  'FolderOutlined',
  'FileOutlined',
  'SearchOutlined',
  'ReloadOutlined',
  'ExportOutlined',
  'ImportOutlined',
  'DownloadOutlined',
  'UploadOutlined',
  'PrinterOutlined',
  'CopyOutlined',
  'CheckCircleOutlined',
  'CloseCircleOutlined',
  'InfoCircleOutlined',
  'WarningOutlined',
  'ExclamationCircleOutlined',
  'QuestionCircleOutlined',
  'LockOutlined',
  'UnlockOutlined',
  'EyeOutlined',
  'EyeInvisibleOutlined',
  'PlusCircleOutlined',
  'MinusCircleOutlined',
  'ArrowUpOutlined',
  'ArrowDownOutlined',
  'ArrowLeftOutlined',
  'ArrowRightOutlined',
  'UpCircleOutlined',
  'DownCircleOutlined',
  'LeftCircleOutlined',
  'RightCircleOutlined'
]

const iconMap = {
  SettingOutlined,
  UserOutlined,
  TeamOutlined,
  SafetyOutlined,
  FileTextOutlined,
  DashboardOutlined,
  HomeOutlined,
  MenuOutlined,
  AppstoreOutlined,
  FolderOutlined,
  FileOutlined,
  SearchOutlined,
  ReloadOutlined,
  ExportOutlined,
  ImportOutlined,
  DownloadOutlined,
  UploadOutlined,
  PrinterOutlined,
  CopyOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  InfoCircleOutlined,
  WarningOutlined,
  ExclamationCircleOutlined,
  QuestionCircleOutlined,
  LockOutlined,
  UnlockOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
  PlusCircleOutlined,
  MinusCircleOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  ArrowLeftOutlined,
  ArrowRightOutlined,
  UpCircleOutlined,
  DownCircleOutlined,
  LeftCircleOutlined,
  RightCircleOutlined
}

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '权限名称', dataIndex: 'permissionName', key: 'permissionName' },
  { title: '权限编码', dataIndex: 'permissionCode', key: 'permissionCode' },
  { title: '权限类型', dataIndex: 'permissionType', key: 'permissionType' },
  { title: '路由路径', dataIndex: 'path', key: 'path' },
  { title: '组件路径', dataIndex: 'component', key: 'component' },
  { title: '图标', dataIndex: 'icon', key: 'icon' },
  { title: '排序', dataIndex: 'sortOrder', key: 'sortOrder' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'action', width: 250 }
]

const permissionTreeOptions = ref([])

onMounted(() => {
  fetchPermissions()
})

const fetchPermissions = async () => {
  loading.value = true
  try {
    const res = await permissionAPI.getTree()
    permissions.value = res.data
    permissionTreeOptions.value = res.data
  } catch (error) {
    console.error('获取权限列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  modalType.value = 'create'
  formState.value = {
    permissionName: '',
    permissionCode: '',
    permissionType: 'menu',
    parentId: null,
    path: '',
    component: '',
    icon: '',
    sortOrder: 0,
    description: '',
    status: 1
  }
  modalVisible.value = true
}

const handleCreateChild = (record) => {
  modalType.value = 'create'
  formState.value = {
    permissionName: '',
    permissionCode: '',
    permissionType: 'button',
    parentId: record.id,
    path: '',
    component: '',
    icon: '',
    sortOrder: 0,
    description: '',
    status: 1
  }
  modalVisible.value = true
}

const handleEdit = (record) => {
  modalType.value = 'edit'
  formState.value = {
    permissionName: record.permissionName,
    permissionCode: record.permissionCode,
    permissionType: record.permissionType,
    parentId: record.parentId || null,
    path: record.path,
    component: record.component,
    icon: record.icon,
    sortOrder: record.sortOrder,
    description: record.description,
    status: record.status
  }
  modalVisible.value = true
}

const handleModalOk = async () => {
  try {
    if (modalType.value === 'create') {
      await permissionAPI.create(formState.value)
      message.success('创建成功')
    } else {
      await permissionAPI.update(currentPermission.value.id, formState.value)
      message.success('更新成功')
    }
    modalVisible.value = false
    fetchPermissions()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

const handleModalCancel = () => {
  modalVisible.value = false
}

const handleDelete = async (id) => {
  try {
    await permissionAPI.delete(id)
    message.success('删除成功')
    fetchPermissions()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const getPermissionTypeLabel = (type) => {
  const map = {
    menu: '菜单',
    button: '按钮',
    api: '接口'
  }
  return map[type] || type
}

const getPermissionTypeColor = (type) => {
  const map = {
    menu: 'blue',
    button: 'green',
    api: 'orange'
  }
  return map[type] || 'default'
}

const getIcon = (iconName) => {
  return iconMap[iconName] || null
}

const filterIconOption = (input, option) => {
  return option.value.toLowerCase().includes(input.toLowerCase())
}

const currentPermission = ref(null)
</script>

<style scoped>
</style>
