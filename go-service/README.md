# 神经科学修行助手 - 后端服务

## 项目概述
这是一个基于Go语言的RESTful API服务，为神经科学修行助手提供后端支持。服务包含用户认证、聊天交互、修行计划生成和练习记录管理等核心功能。

## 技术栈
- **主框架**: Gin (v1.11.0)
- **数据库**: MongoDB (go.mongodb.org/mongo-driver v1.15.0)
- **测试框架**: testify (v1.11.1)
- **主要依赖**:
  - Gin Web框架
  - MongoDB驱动
  - 测试工具链

## 项目目录结构
```
go-service/
├── config/          # 配置管理
│   └── config.go    # 配置加载和解析
├── controllers/      # API控制器
│   ├── chat_controller.go
│   ├── health_controller.go
│   ├── practice_plan_controller.go
│   ├── practice_record_controller.go
│   └── user_controller.go
├── database/         # 数据库连接和操作
│   └── database.go
├── middleware/       # 中间件
│   └── auth.go
├── models/           # 数据模型
│   ├── chat_message.go
│   ├── practice_plan.go
│   ├── practice_record.go
│   └── user.go
├── services/         # 业务逻辑服务
│   ├── chat_service.go
│   ├── practice_plan_service.go
│   ├── practice_record_service.go
│   └── user_service.go
├── Dockerfile        # Docker构建配置
├── README.md         # 当前文档
├── go.mod           # Go模块定义
├── go.sum           # 依赖校验和
└── main.go          # 应用入口
```

## 核心组件说明

### 配置管理
配置通过环境变量加载，支持以下参数：
- `MONGODB_URI`: MongoDB连接字符串 (默认: "mongodb://localhost:27017")
- `MONGODB_NAME`: 数据库名称 (默认: "neuro_guide")
- `PYTHON_AI_SERVICE_URL`: Python AI服务地址 (默认: "http://localhost:8000")
- `PORT`: 服务端口 (默认: "8080")
- `WECHAT_APP_ID`: 微信公众平台AppID (默认: "")
- `WECHAT_APP_SECRET`: 微信公众平台AppSecret (默认: "")

### 数据库连接
使用MongoDB作为主数据库，通过`database.go`管理连接：
- `Connect()`: 建立数据库连接
- `Disconnect()`: 关闭数据库连接

### API路由
路由分为以下几组：
1. **基础路由**:
   - `/`: 欢迎信息
   - `/api/health`: 健康检查

2. **用户相关路由**:
   - `/api/user/login`: 用户登录
   - `/api/user/wechat-login`: 微信登录
   - `/api/user/profile/:id`: 获取/更新用户资料

3. **聊天相关路由**:
   - `/api/chat/message`: 发送/接收消息
   - `/api/chat/history`: 获取聊天记录

4. **修行计划相关路由**:
   - `/api/plan/generate`: 生成修行计划
   - `/api/plan/list`: 获取计划列表

5. **练习记录相关路由**:
   - `/api/record/checkin`: 记录练习
   - `/api/record/list`: 获取练习记录

### 认证中间件
提供两种认证方式：
1. `AuthMiddleware()`: 必选认证中间件
2. `OptionalAuthMiddleware()`: 可选认证中间件

## 运行指南

### 环境准备
1. 安装Go 1.24+
2. 安装MongoDB
3. 设置必要的环境变量

### 运行步骤
1. 获取依赖:
```bash
go mod download
```

2. 运行开发服务器:
```bash
go run main.go
```

3. 访问API:
```
http://localhost:8080/api/health
```

### 生产环境构建
```bash
go build -o neuro-guide-service
./neuro-guide-service
```

### 使用Docker运行
1. 构建镜像:
```bash
docker build -t neuro-guide-service .
```

2. 运行容器:
```bash
docker run -d -p 8080:8080 --name neuro-guide-service neuro-guide-service
```

## 测试
运行单元测试:
```bash
go test -v ./services/...
```

## 注意事项
1. 确保Python AI服务已启动并监听指定端口
2. 开发环境与生产环境使用相同的认证机制
3. MongoDB连接字符串包含认证信息时要注意安全性