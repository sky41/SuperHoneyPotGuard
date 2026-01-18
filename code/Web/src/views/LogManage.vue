<template>
  <a-card title="操作日志">
    <template #extra>
      <a-space>
        <a-button @click="handleRefresh">
          <template #icon>
            <ReloadOutlined />
          </template>
          刷新
        </a-button>
        <a-popconfirm
          title="确定要清空所有日志吗？"
          @confirm="handleClear"
        >
          <a-button danger>
            <template #icon>
              <DeleteOutlined />
            </template>
            清空日志
          </a-button>
        </a-popconfirm>
      </a-space>
    </template>

    <a-form layout="inline" style="margin-bottom: 16px">
      <a-form-item label="用户名">
        <a-input
          v-model:value="filters.username"
          placeholder="请输入用户名"
          style="width: 200px"
          @pressEnter="handleSearch"
        />
      </a-form-item>
      <a-form-item label="操作类型">
        <a-input
          v-model:value="filters.operation"
          placeholder="请输入操作类型"
          style="width: 200px"
          @pressEnter="handleSearch"
        />
      </a-form-item>
      <a-form-item label="状态">
        <a-select
          v-model:value="filters.status"
          placeholder="请选择状态"
          style="width: 150px"
          allow-clear
        >
          <a-select-option value="">全部</a-select-option>
          <a-select-option value="1">成功</a-select-option>
          <a-select-option value="0">失败</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" @click="handleSearch">
          <template #icon>
            <SearchOutlined />
          </template>
          搜索
        </a-button>
        <a-button @click="handleReset">
          重置
        </a-button>
      </a-form-item>
    </a-form>

    <a-table
      :columns="columns"
      :data-source="logs"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '成功' : '失败' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'method'">
          <a-tag :color="getMethodColor(record.method)">
            {{ record.method }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleViewDetail(record)">
              <EyeOutlined />
              查看详情
            </a-button>
            <a-popconfirm
              title="确定要删除该日志吗？"
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
      v-model:open="detailModalVisible"
      title="日志详情"
      :footer="null"
      width="800px"
    >
      <a-descriptions bordered :column="2">
        <a-descriptions-item label="日志ID">
          {{ currentLog.id }}
        </a-descriptions-item>
        <a-descriptions-item label="用户名">
          {{ currentLog.username || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="操作类型">
          {{ currentLog.operation || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="请求方法">
          <a-tag :color="getMethodColor(currentLog.method)">
            {{ currentLog.method || '-' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="请求URL">
          {{ currentLog.url || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="IP地址">
          {{ currentLog.ip || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="地理位置">
          {{ currentLog.location || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="请求参数">
          <pre style="max-height: 200px; overflow: auto; background: #f5f5f5; padding: 8px; border-radius: 4px;">{{ currentLog.params || '-' }}</pre>
        </a-descriptions-item>
        <a-descriptions-item label="返回结果">
          <pre style="max-height: 200px; overflow: auto; background: #f5f5f5; padding: 8px; border-radius: 4px;">{{ currentLog.result || '-' }}</pre>
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="currentLog.status === 1 ? 'green' : 'red'">
            {{ currentLog.status === 1 ? '成功' : '失败' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="错误信息" v-if="currentLog.errorMsg">
          <a-tag color="red">
            {{ currentLog.errorMsg }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="执行时间">
          {{ currentLog.executeTime }} ms
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">
          {{ formatTime(currentLog.createdAt) }}
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </a-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  ReloadOutlined,
  DeleteOutlined,
  SearchOutlined,
  EyeOutlined
} from '@ant-design/icons-vue'
import { logAPI } from '@/api'

const loading = ref(false)
const logs = ref([])
const detailModalVisible = ref(false)
const currentLog = ref({})
const filters = ref({
  username: '',
  operation: '',
  status: ''
})

const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username', key: 'username', width: 120 },
  { title: '操作类型', dataIndex: 'operation', key: 'operation', width: 150 },
  { title: '请求方法', dataIndex: 'method', key: 'method', width: 100 },
  { title: '请求URL', dataIndex: 'url', key: 'url', width: 300, ellipsis: true },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 150 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 80 },
  { title: '执行时间(ms)', dataIndex: 'executeTime', key: 'executeTime', width: 120 },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' }
]

onMounted(() => {
  fetchLogs()
})

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await logAPI.getList({
      page: pagination.value.current,
      pageSize: pagination.value.pageSize,
      ...filters.value
    })
    logs.value = res.data.list
    pagination.value.total = res.data.total
  } catch (error) {
    console.error('获取日志列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchLogs()
}

const handleSearch = () => {
  pagination.value.current = 1
  fetchLogs()
}

const handleReset = () => {
  filters.value = {
    username: '',
    operation: '',
    status: ''
  }
  pagination.value.current = 1
  fetchLogs()
}

const handleRefresh = () => {
  fetchLogs()
}

const handleViewDetail = (record) => {
  currentLog.value = record
  detailModalVisible.value = true
}

const handleDelete = async (id) => {
  try {
    await logAPI.delete(id)
    message.success('删除成功')
    fetchLogs()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleClear = async () => {
  try {
    await logAPI.clear()
    message.success('清空成功')
    pagination.value.total = 0
    logs.value = []
  } catch (error) {
    console.error('清空失败:', error)
  }
}

const getMethodColor = (method) => {
  const colors = {
    'GET': 'blue',
    'POST': 'green',
    'PUT': 'orange',
    'DELETE': 'red',
    'PATCH': 'purple'
  }
  return colors[method] || 'default'
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
