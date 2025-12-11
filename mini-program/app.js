App({
  onLaunch() {
    console.log('神经科学修行助手小程序启动')
    // 启动时检查登录状态
    this.checkLoginStatus()
  },
  
  checkLoginStatus() {
    const token = wx.getStorageSync('auth_token')
    const userId = wx.getStorageSync('user_id')
    
    // 如果没有登录信息，跳转到登录页
    if (!token && !userId) {
      wx.redirectTo({
        url: '/pages/login/login'
      })
    }
  },
  
  globalData: {
    userInfo: null
  }
})