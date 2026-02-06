# 编写python-ai-service的代码

## 1. 编写agent模块

### 1.1 创建三个独立的Agent类，使用langchain的deepagent

* `NeuroscienceAnalysisAgent`：神经科学分析agent

* `PracticeMappingAgent`：修行映射agent

* `PracticeAssistantAgent`：修行小助手（修行工具）agent

参考 https://docs.langchain.com/oss/python/deepagents/quickstart

## 2. 实现Agentic RAG逻辑

### 2.1 编写范围
* 只有NeuroscienceAnalysisAgent和PracticeMappingAgent两个agent需要agentic rag的逻辑

### 2.2 决策机制

* 让LLM自主判断是否需要外部知识，
* 让LLM构造检索查询关键字
* 让LLM评估检索内容

### 2.2 检索流程
[User Query] -> [Agentic Rag (LLM + Reasoning act Loop) -> LLM决策：是否需要外部知识？ -> 若需要 → 自主构造检索查询关键字 → 调用检索器 → 评估结果 → 决定是否重试 -> 结果ok就返回检索内容string。若不需要 → 返回结果空string]

### 2.3 检索库

向量存储使用langchain-chroma， embedding使用HuggingFaceBgeEmbeddings， 文档切片使用langchain-text-splitters，rank使用rank_bm25，不用实现上传文档的接口，就直接扫描本地的目录，有文档更新时触发

## 3. 编写prompt模板

### 3.1 编写agent的prompt模板

* 神经科学分析模板

* 修行映射模板

* 修行小助手模板（专注于生成HTML页面）

### 3.2 添加rag决策prompt

* 用于让LLM判断是否需要检索的prompt

* 用于让LLM构造检索查询关键字prompt

* 用于让LLM评估检索内容的prompt

### 3.3 添加历史会话决策prompt

* 用于让LLM判断是否需要历史对话内容
* 如果需要历史对话就使用tools，从工具中去读去对话记忆



## 4. 编写app.py

### 4.1 接口实现
* 支持非流式响应, 实现参考go-service的service/ai_client.go文件
* 根据conversation\_type调用不同的agent


### 4.2 保持与go-service的兼容性

* 确保API格式不变

* 确保响应格式不变


## 5. 编写工具模块

### 5.1 实现对话记忆工具
该工具仅有NeuroscienceAnalysisAgent使用
* 决策是否需要历史对话内容

* 如果需要历史对话，从go-service的/memory/:conversation_id接口获取

### 5.2 编写检索工具逻辑

* 将agentic RAG逻辑封装为一个deepagents的一个tools


## 6. 测试和验证

### 6.1 单元测试

* 测试每个agent的基本功能

* 测试Agentic RAG逻辑

* 测试使用对话记忆的逻辑

### 6.2 集成测试

* 测试与go-service的交互

* 测试不同conversation\_type的处理

## 关键技术点
1. 一定要去使用doc-langchain的mcp服务阅读langchain的文档，因为我们使用了最新langchain版本去搭建，不要从零开始构建， doc-langchain上有一些例子。
2. 项目是使用uv来管理的依赖包的
3. **兼容性**：确保与现有go-service的兼容性

## 预期结果

1. 三个独立的agent，每个有不同的功能和prompt模板
2. 实现Agentic RAG，LLM自主决策是否需要检索
3. 根据conversation\_type调用不同的agent
4. 实现神经科学分析agent，可以然后LLM自主决策是否使用对话记忆

## 参考信息
1. [Langchain](https://docs.langchain.com/)
2. LLM的模型是云上的，具体调用的url和key，model_name在python-ai-service下的.env文件中
