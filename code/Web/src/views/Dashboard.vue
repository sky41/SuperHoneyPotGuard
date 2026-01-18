<template>
  <div>
    <a-row :gutter="16" style="margin-bottom: 24px">
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic
            title="总用户数"
            :value="stats.userCount"
            :value-style="{ color: '#3f8600' }"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic
            title="角色数"
            :value="stats.roleCount"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <TeamOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic
            title="权限数"
            :value="stats.permissionCount"
            :value-style="{ color: '#722ed1' }"
          >
            <template #prefix>
              <SafetyOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic
            title="操作日志"
            :value="stats.logCount"
            :value-style="{ color: '#cf1322' }"
          >
            <template #prefix>
              <FileTextOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-bottom: 24px">
      <a-col :span="24">
        <a-card title="HFish 蜜罐数据" :loading="hfishLoading">
          <template #extra>
            <a-space>
              <a-button @click="fetchHFishData">
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
                  :value="hfishSysInfo.totalHoneypots"
                  :value-style="{ color: '#52c41a' }"
                >
                  <template #prefix>
                    <SecurityScanOutlined />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card size="small">
                <a-statistic
                  title="运行中蜜罐"
                  :value="hfishSysInfo.activeHoneypots"
                  :value-style="{ color: '#1890ff' }"
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
                  title="总攻击次数"
                  :value="hfishSysInfo.totalAttacks"
                  :value-style="{ color: '#f5222d' }"
                >
                  <template #prefix>
                    <WarningOutlined />
                  </template>
                </a-statistic>
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card size="small">
                <a-descriptions :column="1" bordered>
                  <a-descriptions-item label="最后攻击时间">
                    {{ formatTime(hfishSysInfo.lastAttackTime) }}
                  </a-descriptions-item>
                  <a-descriptions-item label="系统状态">
                    <a-tag :color="hfishSysInfo.systemStatus === 'running' ? 'green' : 'red'">
                      {{ hfishSysInfo.systemStatus === 'running' ? '运行中' : '已停止' }}
                    </a-tag>
                  </a-descriptions-item>
                </a-descriptions>
              </a-card>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :span="24">
        <a-card title="攻击来源 IP" :loading="hfishLoading">
          <a-table
            :columns="attackIPColumns"
            :data-source="hfishAttackIPs"
            :pagination="attackIPPagination"
            @change="handleAttackIPTableChange"
            row-key="id"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'action'">
                <a-space>
                  <a-button type="link" size="small" @click="viewAttackDetails(record.ip)">
                    <EyeOutlined />
                    查看详情
                  </a-button>
                  <a-popconfirm
                    title="确定要封禁该 IP 吗？"
                    @confirm="() => blockIP(record.ip)"
                  >
                    <a-button type="link" size="small" danger>
                      <StopOutlined />
                      封禁
                    </a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-modal
      v-model:open="detailModalVisible"
      title="攻击详情"
      :footer="null"
      width="1000px"
    >
      <a-table
        :columns="attackDetailColumns"
        :data-source="hfishAttackDetails"
        :loading="detailLoading"
        :pagination="false"
        row-key="id"
        size="small"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'attackType'">
            <a-tag :color="getAttackTypeColor(record.attackType)">
              {{ record.attackType }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'protocol'">
            <a-tag color="blue">{{ record.protocol }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-modal>

    <!-- <a-card title="欢迎使用 SuperHoneyPotGuard">
      <p>这是一个基于 HFish 蜜罐平台的主动防御系统。</p>
      <p>系统功能包括：</p>
      <ul>
        <li>用户管理功能</li>
        <li>HFish 数据汇总面板功能</li>
        <li>手动和自动封禁 IP 功能</li>
        <li>自动封禁策略制定功能</li>
      </ul>
    </a-card> -->
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { UserOutlined, TeamOutlined, SafetyOutlined, FileTextOutlined, ReloadOutlined, EyeOutlined, StopOutlined, SecurityScanOutlined, CheckCircleOutlined, WarningOutlined } from '@ant-design/icons-vue'
import { dashboardAPI, hfishAPI } from '@/api'

const loading = ref(false)
const hfishLoading = ref(false)
const detailLoading = ref(false)
const stats = ref({
  userCount: 0,
  roleCount: 0,
  permissionCount: 0,
  logCount: 0
})
const hfishSysInfo = ref({
  totalHoneypots: 0,
  activeHoneypots: 0,
  totalAttacks: 0,
  lastAttackTime: '',
  systemStatus: ''
})
const hfishAttackIPs = ref([])
const hfishAttackDetails = ref([])
const detailModalVisible = ref(false)

const attackIPColumns = [
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 150 },
  { title: '攻击次数', dataIndex: 'count', key: 'count', width: 100 },
  { title: '首次发现', dataIndex: 'firstSeen', key: 'firstSeen', width: 180 },
  { title: '最后发现', dataIndex: 'lastSeen', key: 'lastSeen', width: 180 },
  { title: '操作', key: 'action', width: 200 }
]

const attackDetailColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 150 },
  { title: '攻击类型', dataIndex: 'attackType', key: 'attackType', width: 120 },
  { title: '协议', dataIndex: 'protocol', key: 'protocol', width: 100 },
  { title: '端口', dataIndex: 'port', key: 'port', width: 80 },
  { title: '载荷', dataIndex: 'payload', key: 'payload', width: 200, ellipsis: true },
  { title: '请求时间', dataIndex: 'requestTime', key: 'requestTime', width: 180 },
  { title: '账号', dataIndex: 'account', key: 'account', width: 150 }
]

const attackIPPagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  fetchStats()
  fetchHFishData()
})

const fetchStats = async () => {
  loading.value = true
  try {
    const res = await dashboardAPI.getStats()
    stats.value = res.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchHFishData = async () => {
  hfishLoading.value = true
  try {
    await Promise.all([
      fetchHFishSysInfo(),
      fetchHFishAttackIPs()
    ])
  } catch (error) {
    console.error('获取 HFish 数据失败:', error)
  } finally {
    hfishLoading.value = false
  }
}

const fetchHFishSysInfo = async () => {
  try {
    const res = await hfishAPI.getSysInfo()
    hfishSysInfo.value = res.data
  } catch (error) {
    console.error('获取系统信息失败:', error)
  }
}

const fetchHFishAttackIPs = async () => {
  try {
    const res = await hfishAPI.getAttackIPs()
    hfishAttackIPs.value = res.data
    attackIPPagination.value.total = res.data.length
  } catch (error) {
    console.error('获取攻击 IP 失败:', error)
  }
}

const handleAttackIPTableChange = (pag) => {
  attackIPPagination.value.current = pag.current
  attackIPPagination.value.pageSize = pag.pageSize
}

const viewAttackDetails = async (ip) => {
  detailLoading.value = true
  detailModalVisible.value = true
  try {
    const res = await hfishAPI.getAttackDetails()
    hfishAttackDetails.value = res.data.filter(item => item.ip === ip)
  } catch (error) {
    console.error('获取攻击详情失败:', error)
  } finally {
    detailLoading.value = false
  }
}

const blockIP = async (ip) => {
  try {
    await hfishAPI.blockIP({ ip, reason: '手动封禁' })
    alert('封禁成功')
    fetchHFishData()
  } catch (error) {
    console.error('封禁失败:', error)
  }
}

const getAttackTypeColor = (type) => {
  const colors = {
    'SSH': 'red',
    'HTTP': 'orange',
    'FTP': 'blue',
    'TELNET': 'purple',
    'MYSQL': 'green'
  }
  return colors[type] || 'default'
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}
</script>

<style scoped>
</style>
