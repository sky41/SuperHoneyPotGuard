<template>
  <div>
    <a-row :gutter="16" style="margin-bottom: 24px">
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="总用户数"
            :value="0"
            :value-style="{ color: '#3f8600' }"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="角色数"
            :value="0"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <TeamOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="权限数"
            :value="0"
            :value-style="{ color: '#722ed1' }"
          >
            <template #prefix>
              <SafetyOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="操作日志"
            :value="0"
            :value-style="{ color: '#cf1322' }"
          >
            <template #prefix>
              <FileTextOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-card title="欢迎使用 SuperHoneyPotGuard">
      <p>这是一个基于 HFish 蜜罐平台的主动防御系统。</p>
      <p>系统功能包括：</p>
      <ul>
        <li>用户管理功能</li>
        <li>HFish 数据汇总面板功能</li>
        <li>手动和自动封禁 IP 功能</li>
        <li>自动封禁策略制定功能</li>
      </ul>
    </a-card>
  </div>
</template>

<script setup>
import { UserOutlined, TeamOutlined, SafetyOutlined, FileTextOutlined } from '@ant-design/icons-vue'
</script>

<style scoped>
</style>
