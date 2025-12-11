#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Test script for memory functionality
"""

import sys
import os

# Add the current directory to the path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from agents.conversation_agent import ConversationAgent
from langchain_openai import ChatOpenAI
from pydantic import SecretStr

def test_memory_functionality():
    """Test the memory functionality"""
    print("Testing Memory Functionality...")
    
    # Initialize LLM with test API key
    llm = ChatOpenAI(
        model="glm-4",
        api_key=SecretStr("test-key"),
        base_url="https://open.bigmodel.cn/api/paas/v4",
        temperature=0.7
    )
    
    # Create agent
    agent = ConversationAgent(llm, [])
    
    # Test conversations
    test_conversations = [
        "我最近总是感到焦虑，怎么办？",
        "什么是神经可塑性？",
        "如何通过冥想缓解焦虑？",
        "大脑如何处理情绪？"
    ]
    
    print("\n=== Testing Short-Term Memory ===")
    for i, user_input in enumerate(test_conversations):
        print(f"\n对话 {i+1}:")
        print(f"用户: {user_input}")
        
        # Mock response for testing
        mock_responses = [
            "焦虑是常见的情绪反应，可以通过正念冥想来缓解。杏仁核是处理焦虑的关键脑区。",
            "神经可塑性是大脑重组自身的能力，通过形成新的神经连接来适应环境变化。",
            "冥想可以通过调节默认模式网络来减少焦虑，促进前额叶皮层的活动。",
            "情绪处理涉及多个脑区，包括杏仁核、前额叶皮层和岛叶的协同工作。"
        ]
        
        # Simulate agent response
        response = mock_responses[i]
        print(f"助手: {response}")
        
        # Add to memory
        agent.memory_manager.add_to_short_term_memory(user_input, response)
        
        # Check if should save to long-term
        if agent.memory_manager.should_save_to_long_term(user_input, response):
            agent.memory_manager.add_to_long_term_memory(user_input, response)
            print("✓ 已保存到长期记忆")
    
    # Test memory retrieval
    print("\n=== Testing Memory Retrieval ===")
    query = "焦虑"
    memories = agent.memory_manager.retrieve_relevant_memories(query)
    print(f"查询 '{query}' 的相关记忆:")
    for memory in memories:
        print(f"- {memory['content']}")
    
    # Test combined memory context
    print("\n=== Testing Combined Memory Context ===")
    current_query = "如何缓解焦虑"
    context = agent.memory_manager.get_combined_memory_context(current_query)
    print("综合记忆上下文:")
    print(context)
    
    # Test short-term memory retrieval
    print("\n=== Testing Short-Term Memory ===")
    short_term = agent.memory_manager.get_short_term_memory()
    print(f"短期记忆大小: {len(short_term)}")
    for turn in short_term:
        print(f"{turn['role']}: {turn['content']}")
    
    # Test memory clearing
    print("\n=== Testing Memory Clearing ===")
    agent.memory_manager.clear_short_term_memory()
    short_term_after_clear = agent.memory_manager.get_short_term_memory()
    print(f"清空后短期记忆大小: {len(short_term_after_clear)}")
    
    print("\n=== Memory Test Completed ===")

if __name__ == "__main__":
    test_memory_functionality()