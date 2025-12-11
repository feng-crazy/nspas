import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    const userId = localStorage.getItem('user_id')
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    } else if (userId) {
      // For development, use user_id as token
      config.headers.Authorization = `Bearer ${userId}`
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_id')
      // Redirect to login page
      window.location.href = '/#/login'
    }
    return Promise.reject(error)
  }
)

// User API
export const userAPI = {
  login: (wechatID, nickname, avatar) => 
    api.post('/user/login', { wechat_id: wechatID, nickname, avatar }),
  
  // 更新微信登录方法，处理回调参数
  wechatLogin: (code, state) => {
    // 在实际应用中，这个方法应该由路由或全局回调处理
    // 这里保留以兼容现有代码
    return api.post('/user/wechat-login', { code })
  },
  
  // 添加新的微信登录回调处理方法（实际在全局处理）
  completeWeChatLogin: (code, state) => {
    return api.post('/user/wechat-login', { code })
  },
  
  getProfile: (userId) => 
    api.get(`/user/profile/${userId}`),
  
  updateProfile: (userId, data) => 
    api.put(`/user/profile/${userId}`, data)
}

// Chat API
export const chatAPI = {
  sendMessage: (message, context = []) => 
    api.post('/chat/message', { message, context }),
  
  getHistory: (limit = 50) => 
    api.get('/chat/history', { params: { limit } }),
  
  clearHistory: () => 
    api.delete('/chat/history')
}

// Practice Plan API
export const planAPI = {
  createPlan: (planData) => 
    api.post('/plan/generate', planData),
  
  getPlans: () => 
    api.get('/plan/list'),
  
  getPlan: (planId) => 
    api.get(`/plan/${planId}`),
  
  deletePlan: (planId) => 
    api.delete(`/plan/${planId}`)
}

// Practice Record API
export const recordAPI = {
  createRecord: (recordData) => 
    api.post('/record/checkin', recordData),
  
  getRecords: (planId = null) => 
    api.get('/record/list', { params: planId ? { plan_id: planId } : {} }),
  
  getRecord: (recordId) => 
    api.get(`/record/${recordId}`),
  
  updateRecord: (recordId, data) => 
    api.put(`/record/${recordId}`, data)
}

export default api