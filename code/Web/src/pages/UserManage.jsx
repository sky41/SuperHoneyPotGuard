import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, Select, message, Popconfirm, Tag, Card } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, LockOutlined, StopOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { userAPI, roleAPI } from '../api';

const UserManage = () => {
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [modalType, setModalType] = useState('create');
  const [currentUser, setCurrentUser] = useState(null);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [form] = Form.useForm();

  useEffect(() => {
    fetchUsers();
    fetchRoles();
  }, []);

  const fetchUsers = async (params = {}) => {
    setLoading(true);
    try {
      const res = await userAPI.getList({
        page: pagination.current,
        pageSize: pagination.pageSize,
        ...params
      });
      setUsers(res.data.list);
      setPagination(prev => ({
        ...prev,
        total: res.data.total
      }));
    } catch (error) {
      console.error('获取用户列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchRoles = async () => {
    try {
      const res = await roleAPI.getAll();
      setRoles(res.data);
    } catch (error) {
      console.error('获取角色列表失败:', error);
    }
  };

  const handleTableChange = (newPagination) => {
    setPagination(newPagination);
    fetchUsers({ page: newPagination.current, pageSize: newPagination.pageSize });
  };

  const handleCreate = () => {
    setModalType('create');
    setCurrentUser(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record) => {
    setModalType('edit');
    setCurrentUser(record);
    form.setFieldsValue({
      ...record,
      roleIds: record.roles?.map(r => r.id)
    });
    setModalVisible(true);
  };

  const handleDelete = async (id) => {
    try {
      await userAPI.delete(id);
      message.success('删除成功');
      fetchUsers();
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  const handleStatusChange = async (id, status) => {
    try {
      await userAPI.updateStatus(id, status);
      message.success('状态更新成功');
      fetchUsers();
    } catch (error) {
      console.error('状态更新失败:', error);
    }
  };

  const handleResetPassword = async (id) => {
    Modal.confirm({
      title: '重置密码',
      content: (
        <Form form={form} layout="vertical">
          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[{ required: true, message: '请输入新密码' }, { min: 6, message: '密码长度至少6个字符' }]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
        </Form>
      ),
      onOk: async () => {
        try {
          const values = await form.validateFields();
          await userAPI.resetPassword(id, values.newPassword);
          message.success('密码重置成功');
          form.resetFields();
        } catch (error) {
          console.error('密码重置失败:', error);
        }
      }
    });
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      if (modalType === 'create') {
        await userAPI.create(values);
        message.success('创建成功');
      } else {
        await userAPI.update(currentUser.id, values);
        message.success('更新成功');
      }
      setModalVisible(false);
      form.resetFields();
      fetchUsers();
    } catch (error) {
      console.error('操作失败:', error);
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username'
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email'
    },
    {
      title: '真实姓名',
      dataIndex: 'real_name',
      key: 'real_name'
    },
    {
      title: '角色',
      dataIndex: 'roles',
      key: 'roles',
      render: (roles) => roles?.map(role => (
        <Tag key={role.id} color="blue">{role.role_name}</Tag>
      ))
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      )
    },
    {
      title: '最后登录时间',
      dataIndex: 'last_login_time',
      key: 'last_login_time',
      render: (time) => time || '-'
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<LockOutlined />}
            onClick={() => handleResetPassword(record.id)}
          >
            重置密码
          </Button>
          {record.status === 1 ? (
            <Popconfirm
              title="确定要禁用该用户吗？"
              onConfirm={() => handleStatusChange(record.id, 0)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" icon={<StopOutlined />} danger>
                禁用
              </Button>
            </Popconfirm>
          ) : (
            <Popconfirm
              title="确定要启用该用户吗？"
              onConfirm={() => handleStatusChange(record.id, 1)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" icon={<CheckCircleOutlined />} style={{ color: '#52c41a' }}>
                启用
              </Button>
            </Popconfirm>
          )}
          <Popconfirm
            title="确定要删除该用户吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" icon={<DeleteOutlined />} danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <Card title="用户管理">
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增用户
        </Button>
      </Space>
      <Table
        columns={columns}
        dataSource={users}
        rowKey="id"
        loading={loading}
        pagination={pagination}
        onChange={handleTableChange}
      />
      <Modal
        title={modalType === 'create' ? '新增用户' : '编辑用户'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
        }}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" disabled={modalType === 'edit'} />
          </Form.Item>
          {modalType === 'create' && (
            <Form.Item
              name="password"
              label="密码"
              rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '密码长度至少6个字符' }]}
            >
              <Input.Password placeholder="请输入密码" />
            </Form.Item>
          )}
          <Form.Item
            name="email"
            label="邮箱"
            rules={[{ type: 'email', message: '请输入有效的邮箱地址' }]}
          >
            <Input placeholder="请输入邮箱" />
          </Form.Item>
          <Form.Item
            name="phone"
            label="手机号"
          >
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item
            name="realName"
            label="真实姓名"
          >
            <Input placeholder="请输入真实姓名" />
          </Form.Item>
          <Form.Item
            name="roleIds"
            label="角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select mode="multiple" placeholder="请选择角色" options={roles.map(r => ({ label: r.role_name, value: r.id }))} />
          </Form.Item>
          <Form.Item
            name="status"
            label="状态"
            initialValue={1}
          >
            <Select>
              <Select.Option value={1}>启用</Select.Option>
              <Select.Option value={0}>禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default UserManage;
