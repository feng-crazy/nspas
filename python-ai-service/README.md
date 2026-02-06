# NSPAS Python AI Service

基于神经科学原理的AI对话服务，提供智能分析、修行映射和辅助功能。

## 📋 项目简介

NSPAS (NeuroScience Practice Assistant System) Python AI Service 是一个基于深度学习和神经科学理论的智能对话系统。该服务提供了三种核心AI能力：

- **神经科学分析** (`analysis`)：将用户的心理症状转化为神经科学解释
- **修行映射** (`mapping`)：生成个性化的修行实践方案
- **修行助手** (`assistant`)：提供交互式的修行指导和工具

## 🏗️ 技术架构

### 核心组件
- **FastAPI**: 高性能Web框架
- **LangChain DeepAgents**: 智能Agent框架
- **LangGraph**: 对话流程管理
- **Chroma Vector Store**: 向量数据库存储
- **HuggingFace Embeddings**: 文本嵌入模型
- **BM25**: 传统信息检索算法

### Agentic RAG 实现
采用先进的Agent驱动检索增强生成模式：
```
[用户查询] → [LLM决策] → {需要外部知识?} → [构造检索查询] → [向量检索] → [结果评估] → [最终响应]
```

## 🚀 快速开始

### 环境要求
- Python 3.11+
- uv 包管理器
- OpenAI API Key 或其他大语言模型访问权限

### 安装步骤

1. 克隆项目
```bash
git clone <repository-url>
cd nspas/python-ai-service
```

2. 安装依赖
```bash
uv sync
```

3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置你的API密钥和其他设置
```

4. 启动服务
```bash
uv run python app.py
# 或者使用uvicorn直接启动
uvicorn app:app --reload --host 0.0.0.0 --port 8000
```

服务将在 `http://localhost:8000` 启动。

## 📡 API 接口

### 主要端点

#### POST `/ai/chat`
核心聊天接口，根据会话类型调用不同Agent。

**请求体**:
```json
{
  "messages": [
    {
      "content": "我想了解焦虑的神经机制",
      "is_user": true
    }
  ],
  "conversation_type": "analysis",
  "user_id": "user123",
  "conversation_id": "conv456"
}
```

**参数说明**:
- `messages`: 消息历史列表
- `conversation_type`: 会话类型 ("analysis"|"mapping"|"assistant")
- `user_id`: 用户标识符
- `conversation_id`: 会话标识符

**响应示例**:
```json
{
  "response": "从神经科学角度来看，焦虑主要涉及杏仁核的过度激活...",
  "conversation_id": "conv456"
}
```

### Agent 类型说明

| 类型 | Agent | 功能描述 |
|------|-------|----------|
| analysis | NeuroscienceAnalysisAgent | 分析心理症状的神经科学机制 |
| mapping | PracticeMappingAgent | 生成个性化的修行实践方案 |
| assistant | PracticeAssistantAgent | 提供交互式修行指导 |

## 🔧 开发指南

### 项目结构
```
python-ai-service/
├── agent/                 # AI Agent模块
│   ├── ns_analysis.py     # 神经科学分析Agent
│   ├── practice_map.py    # 修行映射Agent
│   ├── practice_assistant.py # 修行助手Agent
│   └── *.py              # 相关工具和Prompt模板
├── tests/                # 测试文件
│   ├── test_*.py         # 单元测试
├── app.py               # 主应用入口
├── pyproject.toml       # 项目配置
└── README.md           # 本文档
```

### 添加新功能

1. **创建新的Agent**
```python
from deepagents import BaseAgent

class NewAgent(BaseAgent):
    def __init__(self):
        super().__init__()
        # 初始化Agent逻辑
    
    async def process(self, messages: List[Dict]) -> str:
        # 实现处理逻辑
        pass
```

2. **注册到主应用**
```python
# 在app.py中添加
AGENTS["new_type"] = NewAgent()
```

### 运行测试
```bash
# 运行所有测试
uv run pytest

# 运行特定测试
uv run pytest tests/test_ns_analysis.py

# 运行测试并显示覆盖率
uv run pytest --cov=agent
```

## 📚 核心功能详解

### 1. 神经科学分析 (analysis)
将用户描述的心理症状转换为科学的神经机制解释，帮助用户理解自身认知和情感模式。

### 2. 修行映射 (mapping)  
基于用户输入生成个性化的修行实践方案，结合神经科学原理提供可执行的训练计划。

### 3. 修行助手 (assistant)
提供交互式的修行指导，能够生成HTML格式的练习页面和工具。

## 🔍 Agentic RAG 特性

### 自主决策机制
- LLM自动判断是否需要外部知识检索
- 智能构造检索查询关键词
- 动态评估检索结果质量

### 检索流程
1. 用户提问 → Agent接收
2. LLM分析是否需要外部知识
3. 如需要 → 构造查询词 → 向量检索 → 结果评估
4. 如不需要 → 直接生成回答
5. 返回最终响应

## ⚙️ 配置说明

### 环境变量 (.env)
```env
# API配置
OPENAI_API_KEY=your_openai_key_here
MODEL_NAME=gpt-4-turbo-preview
API_BASE_URL=https://api.openai.com/v1

# 数据库配置
CHROMA_HOST=localhost
CHROMA_PORT=8000

# 日志配置
LOG_LEVEL=INFO
```

### 知识库管理
- 向量存储使用Chroma数据库
- 文档切片使用langchain-text-splitters
- 嵌入模型使用HuggingFace BGE embeddings
- 检索算法结合BM25和向量相似度

## 🤝 集成说明

### 与Go服务集成
本服务设计为与Go后端服务配合使用：
- 通过HTTP API通信
- 共享用户认证和会话管理
- 统一的数据存储格式

### API兼容性
- 保持与go-service的请求/响应格式一致
- 支持流式和非流式响应
- 统一的错误处理机制

## 📈 性能优化

### 当前性能指标
- 响应时间：< 2秒（平均）
- 并发支持：100+ 同时连接
- 内存占用：< 500MB

### 优化建议
1. 使用连接池管理数据库连接
2. 实施缓存策略减少重复计算
3. 异步处理长时间运行的任务
4. 监控和优化向量检索性能

## 🛠️ 故障排除

### 常见问题

**Q: 启动时报错 "ModuleNotFoundError"**
A: 确保使用uv安装依赖：`uv sync`

**Q: API调用返回401错误**
A: 检查.env文件中的API密钥配置是否正确

**Q: 检索功能不工作**
A: 确认Chroma数据库服务正在运行，且文档已正确索引

### 日志查看
```bash
# 查看应用日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log
```

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [LangChain](https://github.com/langchain-ai/langchain) - AI应用框架
- [FastAPI](https://fastapi.tiangolo.com/) - Web框架
- [Chroma](https://github.com/chroma-core/chroma) - 向量数据库
- [HuggingFace](https://huggingface.co/) - 嵌入模型

## 📞 联系方式

如有问题或建议，请提交Issue或联系项目维护者。

---
*NSPAS Python AI Service - 让神经科学智慧触手可及*