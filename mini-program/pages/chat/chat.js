import { chatAPI } from '../../utils/api.js'

Page({
  data: {
    messages: [
      {
        id: 1,
        role: 'assistant',
        content: '您好！我是您的神经科学修行助手。请告诉我您当前的心理状态或想要探讨的哲学问题，我会从神经科学的角度为您解析。'
      }
    ],
    inputValue: '',
    scrollTop: 0,
    isLoading: false
  },

  onLoad() {
    // 检查登录状态
    this.checkLoginStatus()
    this.loadChatHistory()
  },
  
  checkLoginStatus() {
    const token = wx.getStorageSync('auth_token')
    const userId = wx.getStorageSync('user_id')
    
    if (!token && !userId) {
      // 未登录，跳转到登录页面
      wx.redirectTo({
        url: '/pages/login/login'
      })
      return
    }
  },

  onInput(e) {
    this.setData({
      inputValue: e.detail.value
    })
  },

  async loadChatHistory() {
    try {
      const response = await chatAPI.getHistory()
      if (response.messages && response.messages.length > 0) {
        this.setData({
          messages: response.messages.map(msg => ({
            id: msg.id || Date.now() + Math.random(),
            role: msg.role,
            content: msg.message
          }))
        })
        this.scrollToBottom()
      }
    } catch (error) {
      console.error('Failed to load chat history:', error)
    }
  },

  async sendMessage() {
    if (!this.data.inputValue.trim() || this.data.isLoading) return

    const userContent = this.data.inputValue
    
    // 添加用户消息
    const userMessage = {
      id: Date.now(),
      role: 'user',
      content: userContent
    }

    const messages = [...this.data.messages, userMessage]
    
    this.setData({
      messages,
      inputValue: '',
      isLoading: true
    })

    // 滚动到底部
    this.scrollToBottom()

    try {
      // 准备上下文
      const context = this.data.messages.slice(0, -1).map(msg => ({
        role: msg.role,
        message: msg.content
      }))
      
      // 调用后端API
      const response = await chatAPI.sendMessage(userContent, context)
      
      const aiResponse = {
        id: Date.now() + 1,
        role: 'assistant',
        content: response.response
      }
      
      this.setData({
        messages: [...messages, aiResponse],
        isLoading: false
      })
      
      this.scrollToBottom()
    } catch (error) {
      console.error('Failed to send message:', error)
      wx.showToast({
        title: '发送失败，请重试',
        icon: 'none'
      })
      
      const errorResponse = {
        id: Date.now() + 1,
        role: 'assistant',
        content: '抱歉，处理您的请求时出现错误。请稍后重试。'
      }
      
      this.setData({
        messages: [...messages, errorResponse],
        isLoading: false
      })
    }
  },

  scrollToBottom() {
    this.setData({
      scrollTop: 999999
    })
  },

  async clearHistory() {
    try {
      await chatAPI.clearHistory()
      this.setData({
        messages: [{
          id: Date.now(),
          role: 'assistant',
          content: '对话历史已清除。有什么可以帮您的吗？'
        }]
      })
      wx.showToast({
        title: '历史已清除',
        icon: 'success'
      })
    } catch (error) {
      console.error('Failed to clear history:', error)
      wx.showToast({
        title: '清除失败',
        icon: 'none'
      })
    }
  }
})