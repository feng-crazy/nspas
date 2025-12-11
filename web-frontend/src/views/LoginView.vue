<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo">
        <h1>神经科学修行助手</h1>
        <p>探索大脑奥秘，理解思维本质</p>
      </div>
      
      <div class="login-options">
        <button 
          class="wechat-login-btn" 
          @click="handleWeChatLogin"
          :disabled="isLoggingIn"
        >
          <span class="icon-wechat">🌐</span>
          {{ isLoggingIn ? '登录中...' : '微信登录' }}
        </button>
        
        <div class="divider">
          <span>或</span>
        </div>
        
        <div class="dev-login">
          <input 
            v-model="devUserId" 
            placeholder="开发者用户ID (仅限开发环境)" 
            class="dev-input"
          />
          <button @click="handleDevLogin" class="dev-login-btn">
            开发者登录
          </button>
        </div>
      </div>
      
      <div v-if="loginError" class="error-message">
        {{ loginError }}
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { userAPI } from '../services/api'

// 全局微信回调处理函数
window.handleWeChatOAuthCallback = async (code, state) => {
  try {
    const storedState = localStorage.getItem('wechat_oauth_state')
    if (state !== storedState) {
      throw new Error('状态不匹配，可能存在CSRF攻击')
    }
    
    // 调用后端API完成登录
    const response = await userAPI.wechatLogin(code)
    
    if (response.data && response.data.token && response.data.user) {
      // 保存token和用户信息
      localStorage.setItem('auth_token', response.data.token)
      localStorage.setItem('user_id', response.data.user.id)
      if (response.data.user.is_guest) {
        localStorage.setItem('is_guest', 'true')
      } else {
        localStorage.removeItem('is_guest')
      }
      
      // 跳转到主页
      router.push('/')
    } else {
      throw new Error('无效的响应格式')
    }
  } catch (error) {
    console.error('微信登录回调处理失败:', error)
    // 在LoginView组件中存储错误信息
    if (window.LoginView && window.LoginView.loginError) {
      window.LoginView.loginError.value = '微信登录失败，请重试'
    } else {
      // 如果LoginView组件不可用，显示一个简单的错误提示
      alert('微信登录失败，请重试')
    }
  }
}

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const isLoggingIn = ref(false)
    const loginError = ref('')
    const devUserId = ref('')

    // 从环境变量中获取配置
    const WECHAT_APP_ID = import.meta.env.VITE_WECHAT_APP_ID || ''
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || window.location.origin
    
    // 检查URL中是否有微信回调参数（在组件挂载时处理）
    onMounted(() => {
      const urlParams = new URLSearchParams(window.location.search)
      const code = urlParams.get('code')
      const state = urlParams.get('state')
      
      if (code && state) {
        // 清除URL参数
        window.history.replaceState({}, document.title, window.location.pathname)
        // 调用全局回调处理函数
        if (window.handleWeChatOAuthCallback) {
          window.handleWeChatOAuthCallback(code, state)
        }
      }
    })

    // 微信登录处理函数
    const handleWeChatLogin = async () => {
      try {
        isLoggingIn.value = true
        loginError.value = ''
        
        if (!WECHAT_APP_ID) {
          throw new Error('微信AppID未配置，请检查环境变量')
        }

        // 构建回调URL
        const redirectUri = encodeURIComponent(`${API_BASE_URL}/wechat-callback`)
        const state = 'login_' + Date.now() // 防止CSRF攻击
        
        // 存储state到localStorage，用于验证回调
        localStorage.setItem('wechat_oauth_state', state)
        
        // 跳转到微信登录页面
        window.location.href = `https://open.weixin.qq.com/connect/qrconnect?appid=${WECHAT_APP_ID}&redirect_uri=${redirectUri}&response_type=code&scope=snsapi_login&state=${state}#wechat_redirect`
      } catch (error) {
        console.error('微信登录准备失败:', error)
        loginError.value = '微信登录准备失败，请稍后重试'
        isLoggingIn.value = false
      }
    }

    // 处理微信回调
    const handleWeChatCallback = async (code, state) => {
      try {
        const storedState = localStorage.getItem('wechat_oauth_state')
        if (state !== storedState) {
          throw new Error('状态不匹配，可能存在CSRF攻击')
        }
        
        isLoggingIn.value = true
        loginError.value = ''
        
        // 调用后端API完成登录
        const response = await userAPI.wechatLogin(code)
        
        // 保存token和用户信息
        localStorage.setItem('auth_token', response.data.token)
        localStorage.setItem('user_id', response.data.user.id)
        if (response.data.user.is_guest) {
          localStorage.setItem('is_guest', 'true')
        } else {
          localStorage.removeItem('is_guest')
        }
        
        // 跳转到主页
        router.push('/')
      } catch (error) {
        console.error('微信登录回调处理失败:', error)
        loginError.value = '微信登录失败，请重试'
      } finally {
        isLoggingIn.value = false
      }
    }
    
    // 开发者登录处理函数
    const handleDevLogin = async () => {
      if (!devUserId.value.trim()) {
        loginError.value = '请输入用户ID'
        return
      }
      
      try {
        isLoggingIn.value = true
        loginError.value = ''
        
        // 保存开发者用户ID
        localStorage.setItem('user_id', devUserId.value)
        localStorage.setItem('is_guest', 'true') // 开发者登录作为游客
        localStorage.removeItem('auth_token') // 移除token，使用user_id模式
        
        // 跳转到主页
        router.push('/')
      } catch (error) {
        console.error('开发者登录失败:', error)
        loginError.value = '开发者登录失败，请重试'
      } finally {
        isLoggingIn.value = false
      }
    }
    
    return {
      isLoggingIn,
      loginError,
      devUserId,
      handleWeChatLogin,
      handleDevLogin,
      // 暴露给路由守卫使用的方法
      handleWeChatCallback
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  padding: 40px;
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.logo h1 {
  margin: 0 0 10px;
  color: #333;
  font-size: 28px;
}

.logo p {
  margin: 0 0 30px;
  color: #666;
  font-size: 16px;
}

.wechat-login-btn {
  width: 100%;
  padding: 15px;
  background: #07c160;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.wechat-login-btn:hover:not(:disabled) {
  background: #06ad56;
}

.wechat-login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.icon-wechat {
  font-size: 20px;
}

.divider {
  margin: 25px 0;
  position: relative;
  text-align: center;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #eee;
}

.divider span {
  background: white;
  padding: 0 15px;
  color: #999;
  position: relative;
}

.dev-login {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.dev-input {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.dev-login-btn {
  padding: 12px;
  background: #42b983;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
}

.dev-login-btn:hover {
  background: #359c6d;
}

.error-message {
  margin-top: 20px;
  padding: 12px;
  background: #fee;
  color: #c33;
  border-radius: 6px;
  border: 1px solid #fcc;
}
</style>
