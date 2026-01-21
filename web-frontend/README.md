# Web Frontend

## 项目介绍

Web Frontend 是一个基于 React + TypeScript + Vite 构建的现代化前端应用，为 NSPAS 系统提供用户界面。该应用采用模块化设计，包含多个功能页面，支持用户认证、数据分析、地图展示、智能助手等核心功能。

## 概述

该前端应用使用了以下技术栈：
- React 19.2.0
- TypeScript 5.9.3
- Vite 7.2.5
- React Router 7.12.0
- ky (HTTP 客户端)
- React Markdown (Markdown 渲染)

应用采用了组件化开发方式，通过 React Context 管理用户认证状态，使用 React Router 实现页面路由，并提供了响应式的用户界面。

## 项目结构

```
web-frontend/
├── public/             # 静态资源
├── src/
│   ├── assets/         # 图片等静态资源
│   ├── components/     # 可复用组件
│   │   ├── ChatInterface.tsx       # 聊天界面组件
│   │   ├── ConversationHistory.tsx # 对话历史组件
│   │   ├── ConversationLayout.tsx  # 对话布局组件
│   │   └── Navbar.tsx              # 导航栏组件
│   ├── context/        # React Context
│   │   └── AuthContext.tsx         # 认证上下文
│   ├── pages/          # 页面组件
│   │   ├── Home.tsx                # 首页
│   │   ├── Analysis.tsx            # 分析页面
│   │   ├── Assistant.tsx           # 助手页面
│   │   ├── History.tsx             # 历史记录页面
│   │   ├── Login.tsx               # 登录页面
│   │   ├── Mapping.tsx             # 地图页面
│   │   └── Tools.tsx               # 工具页面
│   ├── services/       # 服务
│   │   └── api.ts                  # API 服务
│   ├── types/          # TypeScript 类型定义
│   ├── App.tsx         # 应用主组件
│   ├── main.tsx        # 应用入口
│   └── index.css       # 全局样式
├── .gitignore
├── eslint.config.js    # ESLint 配置
├── index.html
├── package.json        # 项目配置和依赖
├── tsconfig.json       # TypeScript 配置
└── vite.config.ts      # Vite 配置
```

## 路由配置

应用使用 React Router 进行路由管理，主要路由如下：

| 路径 | 组件 | 说明 | 权限 |
|------|------|------|------|
| `/` | Home | 首页 | 需要认证 |
| `/analysis` | Analysis | 分析页面 | 需要认证 |
| `/mapping` | Mapping | 地图页面 | 需要认证 |
| `/assistant` | Assistant | 智能助手页面 | 需要认证 |
| `/tools` | Tools | 工具页面 | 需要认证 |
| `/history` | History | 历史记录页面 | 需要认证 |
| `/login` | Login | 登录页面 | 无需认证 |

应用实现了路由保护机制，未登录用户将被重定向到登录页面。

## 运行指南

### 前置条件

- Node.js 18.x 或更高版本
- npm 或 yarn

### 安装依赖

```bash
# 使用 npm
npm install

# 或使用 yarn
yarn install
```

### 开发模式运行

```bash
# 使用 npm
npm run dev

# 或使用 yarn
yarn dev
```

应用将在 http://localhost:5173 启动。

### 构建生产版本

```bash
# 使用 npm
npm run build

# 或使用 yarn
yarn build
```

构建产物将输出到 `dist` 目录。

### 预览生产版本

```bash
# 使用 npm
npm run preview

# 或使用 yarn
yarn preview
```

### 代码检查

```bash
# 使用 npm
npm run lint

# 或使用 yarn
yarn lint
```

## 功能特性

- **用户认证**：基于 React Context 实现的认证状态管理
- **数据分析**：提供数据可视化和分析功能
- **地图展示**：集成地图服务，展示地理数据
- **智能助手**：提供基于 AI 的智能助手功能
- **历史记录**：记录和管理用户操作历史
- **工具集**：提供各种实用工具

## 技术亮点

- **TypeScript**：提供类型安全，减少运行时错误
- **组件化**：模块化设计，提高代码复用性
- **响应式**：适配不同屏幕尺寸
- **现代化**：使用最新的 React 19 特性
- **性能优化**：基于 Vite 构建，提供快速的开发和构建体验
