Page({
  data: {

  },
  onLoad() {
    // 检查用户是否已登录
    this.checkLoginStatus()
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
  navigateToChat() {
    // 再次检查登录状态
    this.checkLoginStatus()
    
    wx.navigateTo({
      url: '/pages/chat/chat',
    })
  },
  navigateToPractice() {
    // 再次检查登录状态
    this.checkLoginStatus()
    
    wx.navigateTo({
      url: '/pages/practice/practice',
    })
  }
})