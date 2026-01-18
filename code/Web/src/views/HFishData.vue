<template>
  <a-card title="HFish 蜜罐数据" :loading="loading">
    <template #extra>
      <a-space>
        <a-button @click="fetchData">
          <template #icon>
            <ReloadOutlined />
          </template>
          刷新
        </a-button>
      </a-space>
    </template>

    <a-row :gutter="16">
      <a-col :span="6">
        <a-card size="small">
          <a-statistic
            title="蜜罐总数"
            :value="sysInfo.total_honeypots"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <DatabaseOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small">
          <a-statistic
            title="在线蜜罐"
            :value="sysInfo.total_online_honeypots"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small">
          <a-statistic
            title="离线蜜罐"
            :value="sysInfo.total_offline_honeypots"
            :value-style="{ color: '#f5222d' }"
          >
            <template #prefix>
              <CloseCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small">
          <a-descriptions :column="1" bordered size="small">
            <a-descriptions-item label="系统状态">
              <a-tag color="green">运行中</a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="蜜罐详情" size="small">
          <a-row :gutter="8">
            <a-col :span="12" v-for="(count, type) in honeypotCounts" :key="type">
              <a-card size="small" :title="type">
                <a-statistic
                  :value="count"
                  :value-style="{ fontSize: '16px' }"
                >
                  <template #prefix>
                    <AppstoreOutlined />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="客户端节点" size="small">
          <a-table
            :columns="clientColumns"
            :data-source="sysInfo.clients"
            :pagination="false"
            row-key="name"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'ip'">
                <a-tag color="blue">{{ record.ip }}</a-tag>
              </template>
              <template v-else-if="column.key === 'createTime'">
                {{ formatTime(record.create_time) }}
              </template>
              <template v-else-if="column.key === 'honeypots'">
                <a-tag v-if="record.honeypots && record.honeypots.length > 0" color="green">
                  {{ record.honeypots.length }} 个蜜罐
                </a-tag>
                <a-tag v-else color="default">无蜜罐</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="内置节点蜜罐状态" size="small">
          <a-table
            :columns="honeypotColumns"
            :data-source="builtinHoneypots"
            :pagination="false"
            row-key="name"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'state'">
                <a-tag :color="getHoneypotStateColor(record.state)">
                  {{ getHoneypotStateText(record.state) }}
                </a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </a-card>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  ReloadOutlined,
  DatabaseOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  AppstoreOutlined
} from '@ant-design/icons-vue'
import { hfishAPI } from '@/api'

const loading = ref(false)
const sysInfo = ref({
  total_honeypots: 0,
  total_online_honeypots: 0,
  total_offline_honeypots: 0,
  honeypot_self_cnt: {},
  clients: []
})

const clientColumns = [
  { title: '节点名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 300 },
  { title: '创建时间', dataIndex: 'createTime', key: 'createTime', width: 180 },
  { title: '蜜罐数量', dataIndex: 'honeypots', key: 'honeypots', width: 120 }
]

const honeypotColumns = [
  { title: '蜜罐类型', dataIndex: 'type', key: 'type', width: 150 },
  { title: '蜜罐名称', dataIndex: 'name', key: 'name', width: 200 },
  { title: '状态', dataIndex: 'state', key: 'state', width: 120 }
]

const honeypotCounts = computed(() => {
  if (!sysInfo.value.honeypot_self_cnt) return {}
  return Object.entries(sysInfo.value.honeypot_self_cnt).map(([type, count]) => {
    return {
      type: type.split('|')[1] || type,
      count: count
    }
  })
})

const builtinHoneypots = computed(() => {
  const builtinNode = sysInfo.value.clients.find(client => client.name === '内置节点')
  if (!builtinNode || !builtinNode.honeypots) return []
  return builtinNode.honeypots
})

onMounted(() => {
  fetchData()
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await hfishAPI.getSysInfo()
    sysInfo.value = res.data
  } catch (error) {
    console.error('获取 HFish 数据失败:', error)
    message.error('获取 HFish 数据失败')
  } finally {
    loading.value = false
  }
}

const getHoneypotStateColor = (state) => {
  const colors = {
    1: 'green',
    2: 'orange',
    3: 'red'
  }
  return colors[state] || 'default'
}

const getHoneypotStateText = (state) => {
  const texts = {
    1: '在线',
    2: '离线',
    3: '异常'
  }
  return texts[state] || '未知'
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}
</script>

<style scoped>
</style>
