#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Memory Manager for Neuro Guide
Implements long-term and short-term memory using LangChain and ChromaDB
"""

import os
from typing import List, Dict, Any, Optional
from datetime import datetime, timedelta
import json

from langchain_chroma import Chroma
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_core.documents import Document
from langchain_text_splitters import RecursiveCharacterTextSplitter

class MemoryManager:
    """Manages both short-term and long-term memory for the conversation agent"""
    
    def __init__(self, persist_directory: str = "./storage/memory"):
        """
        Initialize the memory manager
        
        Args:
            persist_directory: Directory to persist vector store
        """
        self.persist_directory = persist_directory
        os.makedirs(persist_directory, exist_ok=True)
        
        # Initialize embeddings
        self.embeddings = HuggingFaceEmbeddings(
            model_name="sentence-transformers/all-MiniLM-L6-v2"
        )
        
        # Initialize short-term memory (simple list-based buffer)
        self.short_term_memory: Dict[str, List[Dict[str, Any]]] = {}  # Changed to dict with user_id as key
        self.max_short_term_memory = 6  # Keep last 6 conversation parts (3 turns)
        
        # Initialize long-term memory (vector store)
        self.long_term_memory = self._initialize_long_term_memory()
        
        # Text splitter for processing documents
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200
        )
    
    def _initialize_long_term_memory(self) -> Chroma:
        """Initialize the long-term memory vector store"""
        return Chroma(
            persist_directory=self.persist_directory,
            embedding_function=self.embeddings,
            collection_name="conversation_memory"
        )
    
    def add_to_short_term_memory(self, user_id: str, user_input: str, agent_response: str) -> None:
        """
        Add a conversation turn to short-term memory for a specific user
        
        Args:
            user_id: User identifier
            user_input: User's message
            agent_response: Agent's response
        """
        if user_id not in self.short_term_memory:
            self.short_term_memory[user_id] = []
            
        conversation_turn = {
            "role": "user",
            "content": user_input,
            "timestamp": datetime.now().isoformat()
        }
        self.short_term_memory[user_id].append(conversation_turn)
        
        conversation_turn = {
            "role": "assistant",
            "content": agent_response,
            "timestamp": datetime.now().isoformat()
        }
        self.short_term_memory[user_id].append(conversation_turn)
        
        # Keep only the most recent conversations
        if len(self.short_term_memory[user_id]) > self.max_short_term_memory:
            self.short_term_memory[user_id] = self.short_term_memory[user_id][-self.max_short_term_memory:]

    def add_to_long_term_memory(self, user_id: str, user_input: str, agent_response: str, 
                               metadata: Optional[Dict[str, Any]] = None) -> None:
        """
        Add important conversation to long-term memory for a specific user
        
        Args:
            user_id: User identifier
            user_input: User's message
            agent_response: Agent's response
            metadata: Additional metadata for the memory
        """
        if metadata is None:
            metadata = {}
        
        # Create a document with the conversation
        conversation_text = f"User: {user_input}\nAgent: {agent_response}"
        
        # Add timestamp and other metadata
        metadata.update({
            "user_id": user_id,
            "timestamp": datetime.now().isoformat(),
            "type": "conversation_memory"
        })
        
        document = Document(
            page_content=conversation_text,
            metadata=metadata
        )
        
        # Add to vector store
        self.long_term_memory.add_documents([document])
    
    def retrieve_relevant_memories(self, user_id: str, query: str, k: int = 5) -> List[Dict[str, Any]]:
        """
        Retrieve relevant memories from long-term memory for a specific user
        
        Args:
            user_id: User identifier
            query: Query to search for relevant memories
            k: Number of memories to retrieve
            
        Returns:
            List of relevant memories with metadata
        """
        try:
            # Search for relevant documents with user filter
            results = self.long_term_memory.similarity_search(
                query, 
                k=k,
                filter={"user_id": user_id}
            )
            
            memories = []
            for doc in results:
                memory = {
                    "content": doc.page_content,
                    "metadata": doc.metadata,
                    "relevance_score": 1.0  # Placeholder for similarity score
                }
                memories.append(memory)
            
            return memories
        except Exception as e:
            print(f"Error retrieving memories: {e}")
            return []
    
    def get_short_term_memory(self, user_id: str) -> List[Dict[str, Any]]:
        """
        Get current short-term memory context for a specific user
        
        Args:
            user_id: User identifier
            
        Returns:
            List of recent conversation turns
        """
        return self.short_term_memory.get(user_id, []).copy()
    
    def should_retrieve_memory(self, query: str) -> bool:
        """
        Determine if memory retrieval is needed based on the query
        
        Args:
            query: User's query
            
        Returns:
            Whether memory retrieval is needed
        """
        # Keywords that indicate the need for context or memory
        memory_keywords = [
            "我们之前", "上次", "以前", "历史", "记录", "讨论过", 
            "还记得", "我想起来了", "回顾", "总结一下", "之前的",
            "context", "previous", "before", "history", "last time"
        ]
        
        # Always retrieve memory for certain types of questions
        question_starts = [
            "我们讨论到哪里了", "我们刚才聊了什么", "能回顾一下吗", 
            "总结一下刚才的内容", "之前的对话"
        ]
        
        query_lower = query.lower().strip()
        
        # Check if query contains memory-related keywords
        has_memory_keyword = any(keyword in query_lower for keyword in memory_keywords)
        
        # Check if query starts with question patterns that need context
        is_context_question = any(query_lower.startswith(pattern) for pattern in question_starts)
        
        # For meaningful queries (not just greetings), we might want to check memory
        is_meaningful_query = len(query_lower) > 5
        
        # Skip memory retrieval for simple greetings or short phrases
        greetings = ["你好", "您好", "hi", "hello", "hey"]
        is_greeting = any(query_lower.startswith(g) for g in greetings)
        
        # Don't retrieve memory for very short queries that are likely greetings
        if len(query_lower) <= 5 and is_greeting:
            return False
            
        return has_memory_keyword or is_context_question or (is_meaningful_query and not is_greeting)

    def get_combined_memory_context(self, user_id: str, current_query: str) -> str:
        """
        Get combined memory context for prompt for a specific user
        
        Args:
            user_id: User identifier
            current_query: Current user query
            
        Returns:
            Formatted memory context string
        """
        # Check if we should retrieve memory based on the query
        if not self.should_retrieve_memory(current_query):
            # Return only short-term memory for recent interactions
            short_term = self.get_short_term_memory(user_id)
            if short_term:
                context_parts = ["最近对话:"]
                for turn in short_term:
                    context_parts.append(f"{turn['role']}: {turn['content']}")
                return "\n".join(context_parts)
            return "无相关记忆"
        
        # Get short-term memory
        short_term = self.get_short_term_memory(user_id)
        
        # Get relevant long-term memories
        long_term_memories = self.retrieve_relevant_memories(user_id, current_query)
        
        # Format the context
        context_parts = []
        
        # Add short-term memory
        if short_term:
            context_parts.append("最近对话:")
            for turn in short_term:
                context_parts.append(f"{turn['role']}: {turn['content']}")
        
        # Add relevant long-term memories
        if long_term_memories:
            context_parts.append("\n相关记忆:")
            for memory in long_term_memories:
                context_parts.append(f"- {memory['content']}")
                if 'timestamp' in memory['metadata']:
                    context_parts.append(f"  时间: {memory['metadata']['timestamp']}")
        
        return "\n".join(context_parts) if context_parts else "无相关记忆"
    
    def clear_short_term_memory(self, user_id: str) -> None:
        """Clear short-term memory for a specific user"""
        if user_id in self.short_term_memory:
            self.short_term_memory[user_id] = []
    
    def persist_memory(self) -> None:
        """Persist long-term memory to disk"""
        # In newer versions of langchain-chroma, persistence is handled automatically
        # The persist() method is no longer needed
        pass
    
    def should_save_to_long_term(self, user_input: str, agent_response: str) -> bool:
        """
        Determine if a conversation should be saved to long-term memory
        
        Args:
            user_input: User's message
            agent_response: Agent's response
            
        Returns:
            Whether to save to long-term memory
        """
        # Save conversations that contain important keywords
        important_keywords = [
            "焦虑", "抑郁", "压力", "失眠", "正念", "冥想", 
            "修行", "无我", "知行合一", "神经科学", "脑科学",
            "计划", "训练", "练习", "方案"
        ]
        
        # Check if conversation contains important topics
        conversation_text = user_input + " " + agent_response
        has_important_topic = any(keyword in conversation_text for keyword in important_keywords)
        
        # Check if it's a meaningful conversation (not just greetings)
        is_meaningful = len(user_input.strip()) > 10 and len(agent_response.strip()) > 20
        
        return has_important_topic and is_meaningful

# Example usage
if __name__ == "__main__":
    # Test the memory manager
    memory_manager = MemoryManager()
    
    # Add some test conversations
    test_conversations = [
        ("我最近总是感到焦虑，怎么办？", 
         "焦虑是常见的情绪反应，可以通过正念冥想来缓解。杏仁核是处理焦虑的关键脑区。"),
        ("什么是神经可塑性？",
         "神经可塑性是大脑重组自身的能力，通过形成新的神经连接来适应环境变化。"),
        ("你好", "你好！我是神经科学助手，有什么可以帮您的？")
    ]
    
    for user_input, agent_response in test_conversations:
        # Add to short-term memory
        memory_manager.add_to_short_term_memory(user_input, agent_response)
        
        # Check if should save to long-term memory
        if memory_manager.should_save_to_long_term(user_input, agent_response):
            memory_manager.add_to_long_term_memory(user_input, agent_response)
    
    # Test memory retrieval
    query = "焦虑和神经科学"
    memories = memory_manager.retrieve_relevant_memories(query)
    print("相关记忆:")
    for memory in memories:
        print(f"- {memory['content']}")
    
    # Get combined context
    context = memory_manager.get_combined_memory_context("如何缓解焦虑")
    print("\n综合记忆上下文:")
    print(context)