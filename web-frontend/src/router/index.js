import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import PracticeView from '../views/PracticeView.vue'
import LoginView from '../views/LoginView.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/',
    name: 'chat',
    component: ChatView,
    meta: { requiresAuth: true }
  },
  {
    path: '/practice',
    name: 'practice',
    component: PracticeView,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 添加导航守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('auth_token') || localStorage.getItem('user_id')
  
  // 如果路由需要认证但没有token，则跳转到登录页
  if (to.matched.some(record => record.meta.requiresAuth) && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router