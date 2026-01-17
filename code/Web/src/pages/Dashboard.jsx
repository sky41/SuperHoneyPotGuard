import React from 'react';
import { Card, Row, Col, Statistic } from 'antd';
import { UserOutlined, TeamOutlined, SafetyOutlined, FileTextOutlined } from '@ant-design/icons';

const Dashboard = () => {
  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="角色数"
              value={0}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="权限数"
              value={0}
              prefix={<SafetyOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="操作日志"
              value={0}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>
      <Card title="欢迎使用 SuperHoneyPotGuard">
        <p>这是一个基于 HFish 蜜罐平台的主动防御系统。</p>
        <p>系统功能包括：</p>
        <ul>
          <li>用户管理功能</li>
          <li>HFish 数据汇总面板功能</li>
          <li>手动和自动封禁 IP 功能</li>
          <li>自动封禁策略制定功能</li>
        </ul>
      </Card>
    </div>
  );
};

export default Dashboard;
