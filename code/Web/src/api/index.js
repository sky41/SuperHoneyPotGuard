import request from '@/utils/request'

export const authAPI = {
  register: (data) => request.post('/auth/register', data),
  login: (data) => request.post('/auth/login', data),
  logout: () => request.post('/auth/logout'),
  getCurrentUser: () => request.get('/auth/current'),
  sendVerificationCode: (data) => request.post('/auth/send-verification-code', data),
  sendResetPasswordCode: (data) => request.post('/auth/send-reset-code', data),
  resetPassword: (data) => request.post('/password/reset', data)
}

export const userAPI = {
  getList: (params) => request.get('/user/list', { params }),
  getById: (id) => request.get(`/user/${id}`),
  create: (data) => request.post('/user', data),
  update: (id, data) => request.put(`/user/${id}`, data),
  delete: (id) => request.delete(`/user/${id}`),
  updateStatus: (id, status) => request.patch(`/user/${id}/status`, { status }),
  resetPassword: (id, newPassword) => request.post(`/user/${id}/reset-password`, { newPassword })
}

export const roleAPI = {
  getList: (params) => request.get('/role/list', { params }),
  getAll: () => request.get('/role/all'),
  getById: (id) => request.get(`/role/${id}`),
  create: (data) => request.post('/role', data),
  update: (id, data) => request.put(`/role/${id}`, data),
  delete: (id) => request.delete(`/role/${id}`)
}

export const permissionAPI = {
  getTree: () => request.get('/permission/tree'),
  getAll: () => request.get('/permission/all'),
  getById: (id) => request.get(`/permission/${id}`),
  create: (data) => request.post('/permission', data),
  update: (id, data) => request.put(`/permission/${id}`, data),
  delete: (id) => request.delete(`/permission/${id}`)
}

export const dashboardAPI = {
  getStats: () => request.get('/dashboard/stats')
}

export const logAPI = {
  getList: (params) => request.get('/log/list', { params }),
  getById: (id) => request.get(`/log/${id}`),
  delete: (id) => request.delete(`/log/${id}`),
  clear: () => request.delete('/log/clear')
}

export const hfishAPI = {
  getAttackIPs: () => request.post('/hfish/attack/ips'),
  getAttackDetails: () => request.post('/hfish/attack/details'),
  getAccountInfo: () => request.post('/hfish/account/info'),
  getSysInfo: () => request.get('/hfish/sys/info'),
  blockIP: (data) => request.post('/hfish/block/ip', data)
}
