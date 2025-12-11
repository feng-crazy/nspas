# 神经科学修行助手 API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api` (开发环境)
- **Content-Type**: `application/json`
- **认证方式**: Bearer Token (在Header中传递: `Authorization: Bearer <token>`)

## 通用响应格式

### 成功响应
```json
{
  "data": {...}
}
```

### 错误响应
```json
{
  "error": "错误信息"
}
```

## 接口列表

### 1. 健康检查

**GET** `/health`

检查服务是否正常运行

**响应示例:**
```json
{
  "status": "ok",
  "message": "Go service is running"
}
```

---

### 2. 用户相关接口

#### 2.1 用户登录/注册

**POST** `/user/login`

**请求体:**
```json
{
  "wechat_id": "string",
  "nickname": "string",
  "avatar": "string"
}
```

**响应示例:**
```json
{
  "id": "user_id",
  "wechat_id": "wx_xxx",
  "nickname": "用户昵称",
  "avatar": "头像URL",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 2.2 微信登录

**POST** `/user/wechat-login`

**请求体:**
```json
{
  "code": "微信登录凭证code"
}
```

**响应示例:**
```json
{
  "user": {
    "id": "user_id",
    "wechat_id": "wx_xxx",
    "nickname": "用户昵称",
    "avatar": "头像URL",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "用户访问令牌"
}
```

#### 2.3 获取用户信息

**GET** `/user/profile/:id`

**路径参数:**
- `id`: 用户ID

**响应示例:**
```json
{
  "id": "user_id",
  "wechat_id": "wx_xxx",
  "nickname": "用户昵称",
  "avatar": "头像URL",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 2.4 更新用户信息

**PUT** `/user/profile/:id`

**路径参数:**
- `id`: 用户ID

**请求体:**
```json
{
  "nickname": "新昵称",
  "avatar": "新头像URL"
}
```

---

### 3. 对话相关接口

#### 3.1 发送消息

**POST** `/chat/message`

**请求体:**
```json
{
  "message": "用户消息内容",
  "context": [
    {
      "role": "user",
      "message": "之前的消息"
    }
  ]
}
```

**响应示例:**
```json
{
  "response": "AI回复内容"
}
```

#### 3.2 获取对话历史

**GET** `/chat/history?limit=50`

**查询参数:**
- `limit`: 返回消息数量限制（默认50）

**响应示例:**
```json
{
  "messages": [
    {
      "id": "msg_id",
      "user_id": "user_id",
      "message": "消息内容",
      "role": "user",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 3.3 清除对话历史

**DELETE** `/chat/history`

**响应示例:**
```json
{
  "message": "Chat history cleared"
}
```

---

### 4. 修行方案相关接口

#### 4.1 创建修行方案

**POST** `/plan/generate`

**请求体:**
```json
{
  "title": "方案标题",
  "days": 7,
  "tasks": [
    {
      "day": 1,
      "title": "任务标题",
      "description": "任务描述",
      "scientific_basis": "科学依据"
    }
  ]
}
```

**响应示例:**
```json
{
  "id": "plan_id",
  "user_id": "user_id",
  "title": "方案标题",
  "days": 7,
  "tasks": [...],
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### 4.2 获取方案列表

**GET** `/plan/list`

**响应示例:**
```json
{
  "plans": [
    {
      "id": "plan_id",
      "user_id": "user_id",
      "title": "方案标题",
      "days": 7,
      "tasks": [...],
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 4.3 获取方案详情

**GET** `/plan/:id`

**路径参数:**
- `id`: 方案ID

#### 4.4 删除方案

**DELETE** `/plan/:id`

**路径参数:**
- `id`: 方案ID

---

### 5. 修行记录相关接口

#### 5.1 创建修行记录

**POST** `/record/checkin`

**请求体:**
```json
{
  "plan_id": "plan_id",
  "date": "2024-01-01T00:00:00Z",
  "completed_tasks": ["task1", "task2"],
  "reflection": "今日感悟"
}
```

**响应示例:**
```json
{
  "id": "record_id",
  "user_id": "user_id",
  "plan_id": "plan_id",
  "date": "2024-01-01T00:00:00Z",
  "completed_tasks": ["task1", "task2"],
  "reflection": "今日感悟",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 5.2 获取记录列表

**GET** `/record/list?plan_id=xxx`

**查询参数:**
- `plan_id`: 可选，方案ID（如果提供则只返回该方案的记录）

**响应示例:**
```json
{
  "records": [
    {
      "id": "record_id",
      "user_id": "user_id",
      "plan_id": "plan_id",
      "date": "2024-01-01T00:00:00Z",
      "completed_tasks": ["task1"],
      "reflection": "感悟",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 5.3 获取记录详情

**GET** `/record/:id`

**路径参数:**
- `id`: 记录ID

#### 5.4 更新记录

**PUT** `/record/:id`

**路径参数:**
- `id`: 记录ID

**请求体:**
```json
{
  "completed_tasks": ["task1", "task2"],
  "reflection": "更新后的感悟"
}
```

---

## 错误码

| HTTP状态码 | 说明 |
|-----------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 认证说明

大部分接口需要认证，认证方式为在请求头中添加：

```
Authorization: Bearer <token>
```

对于开发环境，也可以使用查询参数：
```
?user_id=<user_id>
```