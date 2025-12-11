// pages/login/login.js
const { userAPI } = require('../../utils/api.js')

Page({
  data: {
    isLoggingIn: false,
    hasAuth: false, // 是否已经获取用户授权
    isGuest: false, // 是否是游客用户
    userInfo: null // 用户基本信息
  },

  onLoad() {
    // 页面加载时检查是否已经登录
    const token = wx.getStorageSync('auth_token')
    const userId = wx.getStorageSync('user_id')
    const isGuest = wx.getStorageSync('is_guest') === 'true'
    
    if (token || userId) {
      // 如果已登录，直接跳转到首页
      wx.switchTab({
        url: '/pages/index/index'
      })
    } else {
      // 设置游客状态
      this.setData({ isGuest: isGuest })
    }
  },

  // 微信授权登录
  onWechatLogin() {
    if (this.data.isLoggingIn) return
    
    this.setData({ isLoggingIn: true })
    
    // 调用微信登录接口
    wx.login({
      success: (res) => {
        if (res.code) {
          // 获取用户授权信息
          this.getUserProfile()
        } else {
          console.log('登录失败！' + res.errMsg)
          wx.showToast({
            title: '登录失败，请重试',
            icon: 'none'
          })
          this.setData({ isLoggingIn: false })
        }
      },
      fail: (err) => {
        console.error('微信登录失败', err)
        wx.showToast({
          title: '微信登录失败',
          icon: 'none'
        })
        this.setData({ isLoggingIn: false })
      }
    })
  },

  // 获取用户授权信息
  getUserProfile() {
    wx.getUserProfile({
      desc: '用于完善用户资料',
      success: (res) => {
        // 保存用户基本信息
        this.setData({
          userInfo: res.userInfo,
          hasAuth: true
        })
        
        // 获取微信登录凭证
        wx.login({
          success: (loginRes) => {
            if (loginRes.code) {
              // 将微信登录凭证和用户信息发送到后端进行验证
              this.handleBackendLogin(loginRes.code, res.userInfo)
            } else {
              console.log('登录失败！' + loginRes.errMsg)
              wx.showToast({
                title: '登录失败，请重试',
                icon: 'none'
              })
              this.setData({ isLoggingIn: false })
            }
          }
        })
      },
      fail: (err) => {
        console.error('获取用户信息失败', err)
        // 即使没有获取用户信息，仍然尝试登录（作为游客）
        wx.login({
          success: (loginRes) => {
            if (loginRes.code) {
              this.handleBackendLogin(loginRes.code)
            } else {
              console.log('登录失败！' + loginRes.errMsg)
              wx.showToast({
                title: '登录失败，请重试',
                icon: 'none'
              })
              this.setData({ isLoggingIn: false })
            }
          }
        })
      }
    })
  },

  // 绑定手机号
  onBindPhoneNumber(e) {
    if (e.detail.errMsg === 'getPhoneNumber:fail user deny') {
      wx.showToast({
        title: '您拒绝了手机号授权',
        icon: 'none'
      })
      return
    }
    
    // 这里需要调用后端接口解密手机号
    // 实际项目中，你需要将encryptedData和iv发送到后端解密
    this.handleBindPhoneNumber(e.detail)
  },

  // 处理手机号绑定
  handleBindPhoneNumber(detail) {
    // 获取当前用户信息
    const token = wx.getStorageSync('auth_token')
    const userId = wx.getStorageSync('user_id')
    
    if (!userId) {
      wx.showToast({
        title: '请先完成微信登录',
        icon: 'none'
      })
      return
    }
    
    // 模拟调用后端接口解密手机号
    // 实际项目中，你需要将encryptedData和iv发送到后端解密
    wx.showLoading({
      title: '绑定中...',
    })
    
    // 模拟延迟
    setTimeout(() => {
      wx.hideLoading()
      wx.showToast({
        title: '手机号绑定成功',
        icon: 'success'
      })
      
      // 更新用户信息（实际项目中应该从后端获取最新信息）
      this.setData({
        isGuest: false // 绑定手机号后不再是游客
      })
      wx.setStorageSync('is_guest', 'false')
      
      // 跳转到首页
      setTimeout(() => {
        wx.switchTab({
          url: '/pages/index/index'
        })
      }, 1000)
    }, 1000)
  },

  // 处理后端登录逻辑
  handleBackendLogin(code, userInfo = null) {
    // 准备发送到后端的数据
    const requestData = {
      code: code
    }
    
    // 如果有用户信息，也一并发送
    if (userInfo) {
      requestData.userInfo = userInfo
    }
    
    wx.request({
      url: 'https://your-api-domain.com/api/user/wechat-login', // 替换为你的后端微信登录接口
      method: 'POST',
      data: requestData,
      success: (res) => {
        if (res.statusCode === 200 && res.data.token) {
          // 登录成功，保存用户信息和token
          wx.setStorageSync('auth_token', res.data.token)
          wx.setStorageSync('user_id', res.data.user.id)
          
          // 判断是否是游客
          const isGuest = res.data.user.is_guest || false
          this.setData({ isGuest: isGuest })
          wx.setStorageSync('is_guest', isGuest ? 'true' : 'false')
          
          // 显示成功消息并跳转到首页
          wx.showToast({
            title: '登录成功',
            icon: 'success'
          })
          
          // 延迟跳转以显示提示
          setTimeout(() => {
            wx.switchTab({
              url: '/pages/index/index'
            })
          }, 1000)
        } else {
          wx.showToast({
            title: res.data.message || '登录失败',
            icon: 'none'
          })
        }
      },
      fail: (err) => {
        console.error('后端登录失败', err)
        wx.showToast({
          title: '网络错误，请重试',
          icon: 'none'
        })
      },
      complete: () => {
        this.setData({ isLoggingIn: false })
      }
    })
  },

  // 开发者登录（用于开发测试）
  onDevLogin() {
    // 在开发者工具中使用固定用户ID登录
    wx.showModal({
      title: '开发者登录',
      content: '请输入用户ID（仅开发环境使用）',
      editable: true,
      placeholderText: '例如：dev_user_123',
      success: (res) => {
        if (res.confirm && res.content) {
          const userId = res.content.trim()
          if (userId) {
            wx.setStorageSync('user_id', userId)
            wx.setStorageSync('is_guest', 'true')
            wx.removeStorageSync('auth_token')
            
            wx.showToast({
              title: '登录成功',
              icon: 'success'
            })
            
            setTimeout(() => {
              wx.switchTab({
                url: '/pages/index/index'
              })
            }, 1000)
          } else {
            wx.showToast({
              title: '请输入有效用户ID',
              icon: 'none'
            })
          }
        }
      }
    })
  }
})