# 神经科学修行助手 (Neuro Science Practice Assistant)

基于神经科学原理的AI对话系统，通过分析用户情感状态、认知模式及行为倾向，提供实时神经机制解读、修行语录的神经科学解释、权威知识库问答和个性化脑科学修行方案。

## 项目结构

```
nspas/
├── go-service/          # Go后端服务
│   ├── controllers/     # 控制器
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   ├── middleware/      # 中间件
│   ├── config/          # 配置
│   ├── database/        # 数据库连接
│   └── main.go          # 入口文件
├── python-ai-service/   # Python AI服务
│   ├── agents/          # AI Agent
│   ├── memory/          # 记忆管理
│   ├── retrievers/      # 知识检索
│   ├── prompts/         # Prompt模板
│   ├── tools/           # 工具函数
│   └── main.py          # 入口文件
├── web-frontend/        # Web前端
│   ├── src/
│   │   ├── views/       # 页面视图
│   │   ├── services/    # API服务
│   │   └── router/      # 路由配置
│   └── package.json
├── mini-program/        # 微信小程序
│   ├── pages/           # 页面
│   ├── utils/           # 工具函数
│   └── app.js
├── doc/                 # 文档
│   ├── API接口文档.md
│   ├── 部署文档.md
│   └── 用户手册.md
├── scripts/             # 脚本
│   ├── deploy.sh        # 部署脚本
│   └── backup.sh        # 备份脚本
├── docker-compose.yml   # Docker编排
└── README.md
```

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 1.29+
- Node.js 16+ (前端开发)
- Go 1.19+ (Go服务开发)
- Python 3.9+ (Python AI服务开发)

### 使用Docker Compose部署

1. **克隆项目**
```bash
git clone <repository-url>
cd nspas
```

2. **配置环境变量**
```bash
cp .env.example .env
# 编辑.env文件，设置必要的环境变量
```

3. **启动服务**
```bash
# 使用部署脚本
./scripts/deploy.sh

# 或手动启动
docker-compose up -d
```

4. **访问服务**
- Web前端: http://localhost
- Go API: http://localhost:8080/api
- Python AI服务: http://localhost:8000

### 本地开发

#### Go服务开发

```bash
cd go-service
go mod download
go run main.go
```

#### Python AI服务开发

```bash
cd python-ai-service
pip install -r requirements.txt
export OPENAI_API_KEY=your_api_key
python main.py
```

#### Web前端开发

```bash
cd web-frontend
npm install
npm run dev
```

## 功能特性

### 核心功能

1. **智能对话引擎**
   - 用户多轮情感分析
   - 神经科学映射
   - 动态知识调用

2. **思维症状-机制解释器**
   - 将用户描述的症状转换为神经科学解释
   - 提供基于脑科学的建议

3. **修行语录的神经科学解释器**
   - 将哲学和修行语录转换为神经科学解释
   - 帮助理解"知行合一"等概念

4. **交互式知识问答**
   - 支持多轮对话
   - 上下文关联理解
   - 知识库检索与大模型结合

5. **修行方案生成器**
   - 根据用户状态生成个性化修行方案
   - 包含科学依据的训练计划

6. **修行记录功能**
   - 记录修行进展
   - 打卡记录
   - 成长感悟记录

## API文档

详细的API文档请参考 [doc/API接口文档.md](doc/API接口文档.md)

## 部署文档

详细的部署文档请参考 [doc/部署文档.md](doc/部署文档.md)

## 技术栈

### 后端
- **Go**: Gin框架，高并发请求处理
- **Python**: LangChain + LlamaIndex，AI Agent实现
- **MongoDB**: 数据存储

### 前端
- **Web**: Vue 3 + Vite
- **小程序**: 微信小程序原生开发

### 基础设施
- **Docker**: 容器化部署
- **Nginx**: 反向代理和负载均衡

## 开发计划

- [x] 核心对话引擎和知识库MVP
- [x] Go服务层开发
- [x] Python AI层开发
- [x] Web前端开发
- [x] 微信小程序开发
- [x] API接口开发
- [x] 部署配置
- [ ] 单元测试和集成测试
- [ ] 性能优化
- [ ] 监控和日志系统

## 贡献指南

1. Fork项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

如有问题或建议，请提交Issue或联系项目维护者。

## 免责声明

本应用仅供学习和参考使用，不提供医疗诊断或治疗建议。如有严重的心理健康问题，请咨询专业医生或心理治疗师。
