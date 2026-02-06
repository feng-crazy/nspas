import os
import structlog
from pathlib import Path
from typing import List, Optional, Tuple, Dict
from langchain_chroma import Chroma
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_community.document_loaders import DirectoryLoader, TextLoader, PyPDFLoader, Docx2txtLoader
from langchain_core.documents import Document
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
from rank_bm25 import BM25Okapi
from sklearn.feature_extraction.text import TfidfVectorizer
import numpy as np
from langchain.agents import AgentState
from langchain.tools import tool, ToolRuntime
from dotenv import load_dotenv

import time

# 加载环境变量
load_dotenv()

# 配置日志
logger = structlog.get_logger(__name__)

# 全局配置
DATA_DIR = "/Users/hedengfeng/workspace/nspas/data/"
CHROMA_DB_PATH = "./chroma_db"
EMBEDDING_MODEL = "BAAI/bge-small-zh-v1.5"

# 文档加载器映射
LOADER_MAPPING = {
    ".txt": TextLoader,
    ".md": TextLoader,
    ".pdf": PyPDFLoader,
    ".docx": Docx2txtLoader,
    ".doc": Docx2txtLoader,
}

class DocumentMonitor(FileSystemEventHandler):
    """文档目录监控器，用于监控文件变更并更新向量库"""
    
    def __init__(self, rag_tool):
        self.rag_tool = rag_tool
    
    def on_any_event(self, event):
        """处理文件系统事件"""
        if event.is_directory:
            return
        
        # 只处理支持的文件类型
        file_ext = Path(event.src_path).suffix.lower()
        if file_ext not in LOADER_MAPPING:
            return
        
        logger.info(f"检测到文件变更: {event.event_type} - {event.src_path}")
        
        # 延迟更新，避免频繁触发
        time.sleep(1)
        
        # 更新向量库
        try:
            self.rag_tool.update_vector_store()
            logger.info("向量库更新成功")
        except Exception as e:
            logger.error(f"向量库更新失败: {str(e)}")

class AgenticRAGTool:
    """Agentic RAG工具类，实现智能检索决策循环"""
    
    def __init__(self):
        """初始化Agentic RAG工具"""
        logger.info("初始化Agentic RAG工具")
        
        # 初始化嵌入模型
        self.embeddings = HuggingFaceEmbeddings(
            model_name=EMBEDDING_MODEL,
            model_kwargs={'device': 'cpu'},
            encode_kwargs={'normalize_embeddings': True}
        )
        
        # 初始化文本分割器
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len
        )
        
        # 初始化向量存储
        self.vector_store = self._init_vector_store()
        
        # 初始化文档监控器
        self._init_document_monitor()
        
        logger.info("Agentic RAG工具初始化完成")
    
    def _init_vector_store(self) -> Chroma:
        """初始化向量存储"""
        logger.info(f"初始化向量存储，路径: {CHROMA_DB_PATH}")
        
        # 创建向量存储
        vector_store = Chroma(
            persist_directory=CHROMA_DB_PATH,
            embedding_function=self.embeddings
        )
        
        # 检查向量存储是否为空
        if vector_store._collection.count() == 0:
            logger.info("向量存储为空，初始加载文档")
            # 保存向量存储到实例变量
            self.vector_store = vector_store
            # 初始加载文档
            self.update_vector_store()
        else:
            logger.info("使用现有向量存储")
        
        return vector_store
    
    def _init_document_monitor(self):
        """初始化文档监控器"""
        logger.info(f"初始化文档监控器，监控目录: {DATA_DIR}")
        
        # 创建事件处理器
        event_handler = DocumentMonitor(self)
        
        # 创建观察者
        self.observer = Observer()
        self.observer.schedule(event_handler, DATA_DIR, recursive=True)
        
        # 启动监控线程
        self.observer.start()
        logger.info("文档监控器启动成功")
    
    def load_documents(self) -> List[Document]:
        """加载目录中的所有文档"""
        logger.info(f"加载文档，目录: {DATA_DIR}")
        
        # 创建目录加载器
        loader = DirectoryLoader(
            DATA_DIR,
            glob="*.*",
            loader_cls=TextLoader,
            loader_kwargs={"encoding": "utf-8"}
        )
        
        # 加载文档
        documents = loader.load()
        logger.info(f"成功加载 {len(documents)} 个文档")
        
        return documents
    
    def split_documents(self, documents: List[Document]) -> List[Document]:
        """分割文档为小块"""
        logger.info(f"分割文档，原始文档数: {len(documents)}")
        
        # 分割文档
        split_docs = self.text_splitter.split_documents(documents)
        logger.info(f"分割后文档块数: {len(split_docs)}")
        
        return split_docs
    
    def update_vector_store(self):
        """更新向量存储"""
        logger.info("更新向量存储")
        
        # 加载文档
        documents = self.load_documents()
        if not documents:
            logger.warning("没有找到可加载的文档")
            return
        
        # 分割文档
        split_docs = self.split_documents(documents)
        
        # 重新创建向量存储（清空现有数据）
        self.vector_store = Chroma(
            persist_directory=CHROMA_DB_PATH,
            embedding_function=self.embeddings
        )
        
        # 添加新文档到向量存储
        self.vector_store.add_documents(split_docs)
        
        logger.info(f"向量存储更新完成，添加了 {len(split_docs)} 个文档块")
    
    def retrieve_docs(self, query: str, k: int = 10) -> List[Document]:
        """检索相关文档"""
        logger.info(f"检索文档，查询: {query}")
        
        # 向量检索
        docs = self.vector_store.similarity_search(query, k=k)
        logger.info(f"向量检索到 {len(docs)} 个文档")
        
        # 如果检索到的文档数较少，直接返回
        if len(docs) <= 1:
            return docs
        
        # 使用BM25重排序
        reranked_docs = self._rerank_docs(query, docs)
        logger.info("文档重排序完成")
        
        return reranked_docs
    
    def _rerank_docs(self, query: str, docs: List[Document], top_k: int = 5) -> List[Document]:
        """使用BM25对检索到的文档进行重排序"""
        # 提取文档内容
        doc_texts = [doc.page_content for doc in docs]
        
        # 构建BM25模型
        tokenized_docs = [text.split() for text in doc_texts]
        bm25 = BM25Okapi(tokenized_docs)
        
        # 对查询进行分词
        tokenized_query = query.split()
        
        # 计算BM25分数
        scores = bm25.get_scores(tokenized_query)
        
        # 按分数排序文档
        doc_score_pairs = list(zip(docs, scores))
        doc_score_pairs.sort(key=lambda x: x[1], reverse=True)
        
        # 返回前top_k个文档
        return [doc for doc, score in doc_score_pairs[:top_k]]
    
    def generate_retrieval_query(self, original_query: str, llm: Optional[any] = None) -> str:
        """生成优化的检索查询"""
        # 简单实现：如果没有LLM，直接返回原始查询
        # 后续可以集成LLM来生成更优化的检索查询
        return original_query
    
    def decide_if_needs_retrieval(self, query: str, llm: Optional[any] = None) -> bool:
        """决定是否需要检索外部知识"""
        # 简单实现：如果查询包含特定关键词，需要检索
        # 后续可以集成LLM来进行更智能的决策
        retrieval_keywords = ["什么是", "解释", "定义", "原理", "如何", "步骤", "方法", "技巧", "神经科学", "冥想", "修行"]
        return any(keyword in query for keyword in retrieval_keywords)
    
    def run(self, query: str, runtime: Optional[ToolRuntime] = None) -> str:
        """执行Agentic RAG流程
        
        Args:
            query: 用户查询
            runtime: Tool运行时信息
            
        Returns:
            检索到的相关内容或空字符串
        """
        logger.info(f"执行Agentic RAG流程，查询: {query}")
        
        # 1. 决定是否需要检索
        needs_retrieval = self.decide_if_needs_retrieval(query)
        logger.info(f"检索决策: {'需要' if needs_retrieval else '不需要'}")
        
        if not needs_retrieval:
            return ""
        
        # 2. 生成优化的检索查询
        retrieval_query = self.generate_retrieval_query(query)
        logger.info(f"优化后的检索查询: {retrieval_query}")
        
        # 3. 执行检索
        retrieved_docs = self.retrieve_docs(retrieval_query)
        
        # 4. 处理检索结果
        if not retrieved_docs:
            logger.info("没有检索到相关文档")
            return ""
        
        # 5. 格式化检索结果
        formatted_result = self._format_retrieved_docs(retrieved_docs)
        logger.info(f"检索结果格式化完成，长度: {len(formatted_result)}")
        
        return formatted_result
    
    def _format_retrieved_docs(self, docs: List[Document]) -> str:
        """格式化检索到的文档"""
        result = []
        for i, doc in enumerate(docs, 1):
            result.append(f"=== 文档 {i} ===")
            result.append(doc.page_content)
            result.append("")
        
        return "\n".join(result)
    
    def stop(self):
        """停止文档监控器"""
        if hasattr(self, 'observer'):
            self.observer.stop()
            self.observer.join()
            logger.info("文档监控器已停止")

# 初始化全局Agentic RAG工具实例
rag_tool_instance = AgenticRAGTool()

@tool
def agentic_rag_search(
    query: str,
    runtime: ToolRuntime
) -> str:
    """执行智能检索，根据LLM判断动态调用检索器并返回结果。
    
    Args:
        query: 用户查询
        runtime: Tool运行时信息
        
    Returns:
        检索到的相关内容或空字符串
    """
    logger.info(f"调用agentic_rag_search工具，查询: {query}")
    
    try:
        # 调用Agentic RAG工具的run方法
        result = rag_tool_instance.run(query, runtime)
        logger.info(f"agentic_rag_search工具执行完成，结果长度: {len(result)}")
        return result
    except Exception as e:
        logger.error(f"agentic_rag_search工具执行失败: {str(e)}")
        return f"检索失败: {str(e)}"

# Agentic RAG工具说明
agentic_rag_instructions = """你可以使用智能检索工具来获取外部知识。

## `agentic_rag_search`

使用这个工具来执行智能检索，它会根据查询内容自动决定是否需要检索外部知识，并返回相关结果。

参数：
- query: 要检索的查询字符串

返回：
- 检索到的相关内容或空字符串（如果不需要检索）
"""
