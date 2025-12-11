# 神经科学修行助手小程序

这是一个基于微信小程序平台开发的应用，旨在帮助用户通过神经科学原理进行心理修行和自我提升。

## 架构说明

### 技术栈
- **框架**: 微信小程序原生框架
- **语言**: JavaScript
- **样式**: WXSS (WeiXin Style Sheets)
- **模板**: WXML (WeiXin Markup Language)

### 项目结构
```
mini-program/
├── pages/              # 页面目录
│   ├── login/          # 登录页面（新增）
│   ├── chat/           # 聊天页面
│   ├── index/          # 首页
│   └── practice/       # 修行计划页面
├── utils/              # 工具类
├── images/             # 图片资源
├── app.js              # 小程序入口文件
├── app.json            # 全局配置
├── app.wxss            # 全局样式
└── sitemap.json        # 小程序搜索引擎收录配置
```

### 功能模块
1. **登录页面** ([pages/login/](pages/login/))
   - 微信授权登录
   - 开发者登录（测试用）
   - 用户身份验证和令牌存储

2. **首页** ([pages/index/](pages/index/))
   - 导航到聊天和修行计划页面
   - 展示应用基本信息

3. **聊天页面** ([pages/chat/](pages/chat/))
   - 与AI助手进行对话
   - 查看聊天历史记录
   - 清除聊天历史

4. **修行计划页面** ([pages/practice/](pages/practice/))
   - 根据用户状态生成个性化修行计划
   - 查看已有修行计划
   - 记录修行打卡

### 数据流
1. 用户界面操作触发事件处理函数
2. 通过 [utils/api.js](utils/api.js) 调用后端 RESTful API
3. 后端处理请求并返回结果
4. 小程序接收响应并更新视图

## 开发手册

### 环境准备
1. 下载并安装微信开发者工具
2. 使用微信开发者工具打开本项目目录

### 代码规范
- 使用驼峰命名法命名变量和函数
- 页面文件按功能模块组织在 [pages/](pages/) 目录下
- 公共组件和工具函数放在 [utils/](utils/) 目录下
- 页面级样式写在各自页面的 `.wxss` 文件中
- 全局样式写在 [app.wxss](app.wxss) 中

### 添加新页面
1. 在 [pages/](pages/) 目录下创建新页面文件夹
2. 创建页面必需的四个文件：`.js`、`.json`、`.wxml`、`.wxss`
3. 在 [app.json](app.json) 的 `pages` 数组中注册新页面路径

### API 调用
所有后端 API 调用封装在 [utils/api.js](utils/api.js) 中：
- 使用 `request` 函数进行统一封装的网络请求
- 每个模块有专门的 API 对象（如 [userAPI](file:///Users/hedengfeng/workspace/nspas/python-ai-service/agents/conversation_agent.py#L21-L21), [chatAPI](utils/api.js#L57-L65) 等）
- 自动处理身份验证头部信息

### 微信登录实现
1. 新增了登录页面，作为用户首次进入小程序时的默认页面
2. 实现了微信授权登录功能，调用微信 `wx.login` 接口获取临时凭证
3. 通过后端接口验证登录凭证并获取用户信息及访问令牌
4. 开发者登录功能用于开发和测试环境，可直接输入用户ID登录
5. 用户信息和令牌存储在本地缓存中，通过 `wx.setStorageSync` 和 `wx.getStorageSync` 管理

### 注意事项
1. 修改后端 API 地址时需更新 [utils/api.js](utils/api.js) 中的 `API_BASE_URL`
2. 所有异步操作使用 `async/await` 语法处理
3. 错误处理应通过 `try/catch` 或 Promise 的 `.catch()` 方法
4. 登录状态检查已集成到各主要页面，确保用户在未登录状态下无法访问核心功能

## 运行手册

### 开发环境运行
1. 使用微信开发者工具导入项目
2. 在工具中点击"编译"按钮运行项目
3. 可使用模拟器或真机调试功能

### 配置项
1. **后端服务地址**:
   - 修改 [utils/api.js](utils/api.js) 中的 `API_BASE_URL` 为实际的后端地址
   - 确保后端服务已经运行并且可访问

2. **用户认证**:
   - 默认使用本地存储中的 `auth_token` 或 `user_id` 作为认证令牌
   - 可通过 `wx.setStorageSync('user_id', 'your_user_id')` 设置用户ID（开发者登录）
   - 正式环境下通过微信登录获取 `auth_token`

### 页面说明
1. **登录页面** (`pages/login/login`)
   - 用户首次进入小程序时显示的页面
   - 提供微信登录和开发者登录两种方式
   - 登录成功后自动跳转到首页

2. **首页** (`pages/index/index`)
   - 提供导航至聊天和修行计划页面的入口
   - 显示应用标题和简介

3. **聊天页面** (`pages/chat/chat`)
   - 实时显示对话内容
   - 支持发送消息和查看历史记录
   - 提供清除历史记录功能

4. **修行计划页面** (`pages/practice/practice`)
   - 可根据用户输入的状态描述生成修行计划
   - 显示用户的所有修行计划
   - 支持选择特定计划查看详情

### 调试技巧
1. 使用微信开发者工具的控制台查看日志输出
2. 利用网络面板监控 API 请求和响应
3. 在真机上预览以获得更真实的用户体验

### 发布流程
1. 在微信公众平台完成小程序基本信息配置
2. 使用微信开发者工具上传代码
3. 在微信公众平台提交审核
4. 审核通过后发布上线