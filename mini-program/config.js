// 小程序配置
const config = {
  // API基础地址，根据实际部署环境修改
  API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8080/api',
  
  // 微信登录配置
  WECHAT_APP_ID: 'your_wechat_app_id', // 替换为实际的微信AppID
  
  // 路由白名单（不需要登录即可访问的页面）
  WHITE_LIST: [
    '/pages/login/login',
    '/pages/chat/chat' // 微信登录后直接进入聊天页
  ]
}

export default config