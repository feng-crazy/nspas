// API配置
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080/api' // 开发环境地址，根据实际情况修改

// 获取用户token
function getToken() {
  return wx.getStorageSync('auth_token') || wx.getStorageSync('user_id')
}

// 获取用户ID
function getUserId() {
  return wx.getStorageSync('user_id')
}

// 判断是否是游客用户
function isGuest() {
  return wx.getStorageSync('is_guest') === 'true'
}

// 通用请求方法
function request(url, method = 'GET', data = {}) {
  return new Promise((resolve, reject) => {
    const token = getToken()
    const userId = getUserId()
    
    const header = {
      'Content-Type': 'application/json'
    }
    
    if (token) {
      header['Authorization'] = `Bearer ${token}`
    } else if (userId) {
      header['Authorization'] = `Bearer ${userId}`
    }
    
    wx.request({
      url: `${API_BASE_URL}${url}`,
      method: method,
      data: data,
      header: header,
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(res.data)
        } else if (res.statusCode === 401) {
          // 未授权，清除本地存储
          wx.removeStorageSync('auth_token')
          wx.removeStorageSync('user_id')
          reject(new Error('未授权，请重新登录'))
        } else {
          reject(new Error(res.data.error || '请求失败'))
        }
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}

// User API
export const userAPI = {
  // 微信登录（包含用户信息）
  wechatLogin: (code, userInfo = null) => {
    const requestData = { code }
    if (userInfo) {
      requestData.userInfo = userInfo
    }
    return request({
      url: '/user/wechat-login',
      method: 'POST',
      data: requestData
    })
  },
  
  // 绑定手机号
  bindPhoneNumber: (encryptedData, iv) => request({
    url: '/user/bind-phone',
    method: 'POST',
    data: { encryptedData, iv }
  }),
  
  getProfile: (userId) => 
    request(`/user/profile/${userId}`, 'GET'),
  
  updateProfile: (userId, data) => 
    request(`/user/profile/${userId}`, 'PUT', data)
}

// Chat API
export const chatAPI = {
  sendMessage: (message, context = []) => 
    request('/chat/message', 'POST', { message, context }),
  
  getHistory: (limit = 50) => 
    request('/chat/history', 'GET', { limit }),
  
  clearHistory: () => 
    request('/chat/history', 'DELETE')
}

// Practice Plan API
export const planAPI = {
  createPlan: (planData) => 
    request('/plan/generate', 'POST', planData),
  
  getPlans: () => 
    request('/plan/list', 'GET'),
  
  getPlan: (planId) => 
    request(`/plan/${planId}`, 'GET'),
  
  deletePlan: (planId) => 
    request(`/plan/${planId}`, 'DELETE')
}

// Practice Record API
export const recordAPI = {
  createRecord: (recordData) => 
    request('/record/checkin', 'POST', recordData),
  
  getRecords: (planId = null) => 
    request('/record/list', 'GET', planId ? { plan_id: planId } : {}),
  
  getRecord: (recordId) => 
    request(`/record/${recordId}`, 'GET'),
  
  updateRecord: (recordId, data) => 
    request(`/record/${recordId}`, 'PUT', data)
}

export default {
  request,
  API_BASE_URL,
  userAPI,
  chatAPI,
  planAPI,
  recordAPI,
  isGuest
}