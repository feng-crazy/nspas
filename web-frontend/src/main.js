import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 添加路由守卫检查登录状态
router.beforeEach((to, from, next) => {
  // 检查是否存在有效的认证信息
  const isAuthenticated = localStorage.getItem('auth_token') || localStorage.getItem('user_id')
  
  // 如果访问的是登录页，且已经登录，则重定向到首页
  if (to.name === 'login' && isAuthenticated) {
    next({ name: 'chat' })
  } 
  // 如果访问受保护的页面且未登录，则重定向到登录页
  else if (to.matched.some(record => record.meta.requiresAuth) && !isAuthenticated) {
    next({ name: 'login' })
  } 
  // 其他情况允许访问
  else {
    next()
  }
})

const app = createApp(App)

app.use(router)

app.mount('#app')