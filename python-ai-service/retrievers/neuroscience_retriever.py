#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Neuroscience Retriever for Neuro Guide
This module implements retrieval functionality for neuroscience knowledge using LlamaIndex
"""

from typing import List, Dict, Any, Optional
from langchain_core.retrievers import BaseRetriever
from langchain_core.callbacks import CallbackManagerForRetrieverRun, AsyncCallbackManagerForRetrieverRun
from langchain_core.documents import Document
import os
import asyncio

# Import LlamaIndex components
from llama_index.core import VectorStoreIndex, SimpleDirectoryReader, StorageContext
from llama_index.core.node_parser import SimpleNodeParser
# from llama_index.embeddings.openai import OpenAIEmbedding
from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.vector_stores.chroma import ChromaVectorStore
from llama_index.core import Settings
import chromadb

class NeuroscienceRetriever(BaseRetriever):
    """Retriever for neuroscience-related knowledge using LlamaIndex"""
    
    # Use class variables for configuration since BaseRetriever is a Pydantic model
    _data_dir = os.path.join(os.path.dirname(__file__), "../data/neuroscience")
    _persist_dir = os.path.join(os.path.dirname(__file__), "../storage/chroma_db")
    _similarity_threshold = 0.6  # Minimum similarity threshold for filtering results
    
    def __init__(self, data_dir: Optional[str] = None, persist_dir: Optional[str] = None, similarity_threshold: float = 0.6):
        """
        Initialize the neuroscience retriever with LlamaIndex
        
        Args:
            data_dir: Directory containing neuroscience documents
            persist_dir: Directory to persist the vector store
            similarity_threshold: Minimum similarity score threshold for filtering results
        """
        super().__init__()
        # Update class variables if provided
        if data_dir:
            NeuroscienceRetriever._data_dir = data_dir
        if persist_dir:
            NeuroscienceRetriever._persist_dir = persist_dir
        if similarity_threshold is not None:
            NeuroscienceRetriever._similarity_threshold = similarity_threshold
        
        # Initialize in a way that doesn't violate Pydantic constraints
        self._initialize_retriever()
    
    def _initialize_retriever(self):
        """Initialize the retriever after object creation"""
        # Use a separate method to avoid setting instance attributes in __init__
        self._index = None
        self._retriever = None
        self._initialize_index()
    
    def _initialize_index(self):
        """Initialize or load the LlamaIndex vector store"""
        try:
            # Initialize embedding model
            # 配置Hugging Face嵌入模型
            Settings.embed_model = HuggingFaceEmbedding(
                model_name="BAAI/bge-small-en-v1.5",
                embed_batch_size=32,  # 批量处理大小
                max_length=512,       # 最大文本长度
            )

            # Initialize ChromaDB client
            chroma_client = chromadb.PersistentClient(path=self._persist_dir)
            chroma_collection = chroma_client.get_or_create_collection("neuroscience_knowledge")
            
            # Create vector store
            vector_store = ChromaVectorStore(chroma_collection=chroma_collection)
            storage_context = StorageContext.from_defaults(vector_store=vector_store)
            
            # Check if index already exists
            if os.path.exists(self._persist_dir) and chroma_collection.count() > 0:
                print(f"Loading existing index from {self._persist_dir}")
                self._index = VectorStoreIndex.from_vector_store(
                    vector_store, storage_context=storage_context
                )
            else:
                print(f"Creating new index from data in {self._data_dir}")
                self._build_index(storage_context)
            
            # Create retriever only if index was successfully created
            if self._index is not None:
                self._retriever = self._index.as_retriever(
                    similarity_top_k=5,
                    vector_store_query_mode="default"
                )
            
        except Exception as e:
            print(f"Error initializing LlamaIndex: {e}")
            print("Falling back to mock implementation")
            self._retriever = None
    
    def _build_index(self, storage_context: Any):
        """Build the vector index from neuroscience documents"""
        try:
            # Load documents from data directory
            if os.path.exists(self._data_dir):
                reader = SimpleDirectoryReader(self._data_dir)
                documents = reader.load_data()
                
                # Create nodes from documents
                parser = SimpleNodeParser.from_defaults()
                nodes = parser.get_nodes_from_documents(documents)
                
                # Create index
                self._index = VectorStoreIndex(
                    nodes, 
                    storage_context=storage_context,
                    embed_model=Settings.embed_model
                )
                
                # Persist the index
                storage_context.persist(persist_dir=self._persist_dir)
                print(f"Index built and persisted to {self._persist_dir}")
            else:
                print(f"Data directory {self._data_dir} not found. Using empty index.")
                # Create empty index
                self._index = VectorStoreIndex.from_documents(
                    [], 
                    storage_context=storage_context,
                    embed_model=Settings.embed_model
                )
                
        except Exception as e:
            print(f"Error building index: {e}")
            raise
    
    def _filter_documents_by_similarity(self, documents: List[Document]) -> List[Document]:
        """
        Filter documents by similarity score threshold
        
        Args:
            documents: List of Document objects with similarity scores
            
        Returns:
            List of Document objects that meet the similarity threshold
        """
        filtered_docs = []
        for doc in documents:
            # Extract score from metadata (could be 'score' or 'relevance_score')
            score = doc.metadata.get('score', doc.metadata.get('relevance_score', 0))
            print(f"Similarity score: {score}")
            if score >= self._similarity_threshold:
                filtered_docs.append(doc)
        return filtered_docs
    
    def _get_relevant_documents(
        self, query: str, *, run_manager: CallbackManagerForRetrieverRun
    ) -> List[Document]:
        """
        使用LlamaIndex检索与查询相关的神经科学知识文档

        当LlamaIndex不可用或检索失败时会回退到模拟结果

        Args:
            query: 查询字符串
            run_manager: 检索运行的回调管理器

        Returns:
            List[Document]: 包含相关文档的列表，每个文档包含:
                - page_content: 文档内容文本
                - metadata: 元数据字典，包含:
                    * score: 相关性分数
                    * id: 节点ID
                    * 其他来自LlamaIndex节点的元数据

        Raises:
            不会抛出异常，错误时会回退到模拟结果
        """
        if not hasattr(self, '_retriever') or self._retriever is None:
            # Fallback to mock results if LlamaIndex is not available
            mock_docs = self._get_mock_documents(query)
            return self._filter_documents_by_similarity(mock_docs)
        
        try:
            # Use LlamaIndex retriever to get relevant nodes
            nodes = self._retriever.retrieve(query)
            
            # Convert LlamaIndex nodes to LangChain Documents
            documents = []
            for node in nodes:
                doc = Document(
                    page_content=node.text,
                    metadata={
                        "score": node.score,
                        "id": node.node_id,
                        **node.metadata
                    }
                )
                documents.append(doc)
            
            # Filter documents by similarity threshold
            filtered_documents = self._filter_documents_by_similarity(documents)
            return filtered_documents
            
        except Exception as e:
            print(f"Error retrieving documents: {e}")
            # Fallback to mock results
            mock_docs = self._get_mock_documents(query)
            return self._filter_documents_by_similarity(mock_docs)
    
    def _get_mock_documents(self, query: str) -> List[Document]:
        """Fallback mock implementation"""
        mock_results = [
            {
                "text": "The default mode network (DMN) is a network of interacting brain regions known to have higher activity when a person is not focused on the outside world.",
                "relevance_score": 0.95
            },
            {
                "text": "Neuroplasticity refers to the brain's ability to reorganize itself by forming new neural connections throughout life.",
                "relevance_score": 0.87
            },
            {
                "text": "The amygdala is an almond-shaped cluster of nuclei in the brain that plays a key role in processing emotions, particularly fear and anxiety.",
                "relevance_score": 0.82
            },
            {
                "text": "The prefrontal cortex is the front part of the brain's cerebral cortex, responsible for executive functions like decision-making, working memory, and regulating emotions.",
                "relevance_score": 0.78
            },
            {
                "text": "Mindfulness meditation has been shown to increase gray matter density in the hippocampus and prefrontal cortex, regions associated with learning, memory, and emotional regulation.",
                "relevance_score": 0.75
            }
        ]
        
        # Convert mock results to Document objects
        documents = []
        for result in mock_results:
            doc = Document(
                page_content=result["text"],
                metadata={"relevance_score": result["relevance_score"]}
            )
            documents.append(doc)
        
        return documents
    
    async def _aget_relevant_documents(
        self, query: str, *, run_manager: AsyncCallbackManagerForRetrieverRun
    ) -> List[Document]:
        """
        Async version of _get_relevant_documents
        
        Args:
            query: The query string
            run_manager: Callback manager for the retriever run
            
        Returns:
            List of relevant Document objects
        """
        if not hasattr(self, '_retriever') or self._retriever is None:
            # Fallback to mock results if LlamaIndex is not available
            # Run mock documents in an executor to avoid blocking
            loop = asyncio.get_event_loop()
            mock_docs = await loop.run_in_executor(None, self._get_mock_documents, query)
            return self._filter_documents_by_similarity(mock_docs)
        
        try:
            # Use LlamaIndex retriever to get relevant nodes
            # Run the synchronous retrieve method in an executor to avoid blocking
            loop = asyncio.get_event_loop()
            nodes = await loop.run_in_executor(None, self._retriever.retrieve, query)
            
            # Convert LlamaIndex nodes to LangChain Documents
            documents = []
            for node in nodes:
                doc = Document(
                    page_content=node.text,
                    metadata={
                        "score": node.score,
                        "id": node.node_id,
                        **node.metadata
                    }
                )
                documents.append(doc)
            
            # Filter documents by similarity threshold
            filtered_documents = self._filter_documents_by_similarity(documents)
            return filtered_documents
            
        except Exception as e:
            print(f"Error retrieving documents: {e}")
            # Fallback to mock results
            loop = asyncio.get_event_loop()
            mock_docs = await loop.run_in_executor(None, self._get_mock_documents, query)
            return self._filter_documents_by_similarity(mock_docs)

# Example usage (for testing purposes)
if __name__ == "__main__":
    # Create retriever with a higher threshold to demonstrate filtering
    retriever = NeuroscienceRetriever(similarity_threshold=0.8)
    
    # Test retrieval
    results = retriever.invoke("Explain the default mode network")
    print("Retrieved knowledge:")
    if results:
        for i, result in enumerate(results, 1):
            score = result.metadata.get('score', result.metadata.get('relevance_score', 0))
            print(f"{i}. {result.page_content} (Score: {score:.3f})")
    else:
        print("No relevant documents found above the similarity threshold.")