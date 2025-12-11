# 神经科学修行助手 - Web前端

## 项目概述
这是一个基于Vue 3的神经科学修行助手Web前端应用，提供了聊天和练习两个主要功能界面。

## 技术栈
- **框架**: Vue 3 (Composition API)
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **构建工具**: Vite 5
- **主要依赖**:
  - vue@^3.4.0
  - vue-router@^4.3.0
  - axios@^1.6.0

## 项目结构
```
web-frontend/
├── src/
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   ├── router/           # 路由配置
│   │   └── index.js      # 路由定义
│   ├── services/         # API服务
│   │   └── api.js        # API封装
│   ├── views/            # 页面组件
│   │   ├── ChatView.vue  # 聊天界面
│   │   ├── LoginView.vue # 登录界面
│   │   └── PracticeView.vue # 练习界面
│   ├── App.vue           # 根组件
│   └── main.js           # 应用入口
├── index.html            # 入口HTML
└── vite.config.js        # Vite配置
```

## 路由配置
应用包含三个主要路由：
- `/login` (login): 登录界面
- `/` (chat): 聊天界面
- `/practice`: 练习界面

## 运行指南

### 安装依赖
```bash
npm install
```

### 开发环境运行
```bash
npm run dev
```
应用将在 `http://localhost:3000` 启动

### 生产环境构建
```bash
npm run build
```

### 预览构建结果
```bash
npm run preview
```

### 环境变量配置
创建 `.env` 文件配置开发环境变量：
```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WECHAT_APP_ID=your_wechat_app_id_for_dev
```

## 微信登录集成

### 微信SDK集成
项目支持微信登录功能，具体集成方式请参考 [微信登录集成指南](../doc/微信登录集成指南.md)

### 认证机制
- 使用localStorage存储认证信息
- 支持微信登录和开发者登录两种方式
- 路由级别的认证保护

## 代理配置
开发环境下，API请求会被代理到 `http://localhost:8080`

## 注意事项
1. 确保后端服务已启动
2. 开发环境需要Node.js 16+
3. 生产构建前请检查环境变量配置
4. 微信登录功能需要在微信公众平台进行相关配置