<template>
  <div id="app">
    <header v-if="$route.name !== 'login'">
      <h1>神经科学修行助手</h1>
      <nav>
        <router-link to="/">聊天</router-link>
        <router-link to="/practice">修行</router-link>
        <button @click="logout" class="logout-btn">退出</button>
      </nav>
    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: 'App',
  methods: {
    logout() {
      // 清除本地存储的认证信息
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_id')
      
      // 重定向到登录页
      this.$router.push('/login')
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

header {
  background-color: #42b983;
  color: white;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

header h1 {
  margin: 0;
  font-size: 24px;
}

nav {
  display: flex;
  gap: 20px;
}

nav a {
  color: white;
  text-decoration: none;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background 0.3s;
}

nav a.router-link-active,
nav a:hover {
  background: rgba(255, 255, 255, 0.2);
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.3s;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>