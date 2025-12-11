#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Test Agent for Neuro Guide
This script provides a simple test interface for the conversation agent
"""

import os
import sys
from pydantic import SecretStr

# Add the current directory to the path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

# Import required modules
from agents.conversation_agent import ConversationAgent
from langchain_openai import ChatOpenAI
from tools.neuroscience_tool import NeuroscienceExplanationTool
from tools.symptom_analyzer_tool import SymptomAnalyzerTool
from tools.practice_plan_tool import PracticePlanTool

# 设置TOKENIZERS_PARALLELISM环境变量以避免huggingface/tokenizers的警告
os.environ["TOKENIZERS_PARALLELISM"] = "false"

def create_test_agent():
    """
    Create a test agent with the same configuration as the main service
    """
    # Initialize the LLM with OpenAI standard interface
    # This supports multiple providers like ZhipuAI, Qwen, etc.
    api_key = os.getenv("OPENAI_API_KEY")
    base_url = os.getenv("OPENAI_BASE_URL")
    model_name = os.getenv("MODEL_NAME", "glm-4")

    if not api_key:
        # Default to ZhipuAI if no API key is provided
        api_key = "you-are-key"

    llm = ChatOpenAI(
        model=model_name,
        api_key=SecretStr(api_key),
        base_url=base_url if base_url else "https://open.bigmodel.cn/api/paas/v4",
        temperature=0.7,
        max_completion_tokens=2048
    )

    # Create tools
    tools = [
        NeuroscienceExplanationTool(),
        SymptomAnalyzerTool(),
        PracticePlanTool()
    ]

    # Create the agent
    agent = ConversationAgent(llm, tools)
    return agent

def run_test_conversation():
    """
    Run a test conversation with the agent
    """
    print("=== Neuro Guide Test Agent ===")
    print("Initializing agent...")
    
    agent = create_test_agent()
    
    print("Agent initialized successfully!")
    print("\nYou can now chat with the Neuro Guide agent.")
    print("Type 'quit' to exit, 'clear' to clear memory, or 'memory' to see current memory.")
    print("-" * 50)
    
    while True:
        try:
            user_input = input("\nYou: ").strip()
            
            if user_input.lower() == 'quit':
                print("Goodbye!")
                break
            elif user_input.lower() == 'clear':
                agent.memory_manager.clear_short_term_memory()
                print("Memory cleared.")
                continue
            elif user_input.lower() == 'memory':
                short_term = agent.memory_manager.get_short_term_memory()
                print(f"Memory contains {len(short_term)} entries:")
                for entry in short_term:
                    print(f"  {entry['role']}: {entry['content']}")
                continue
            elif not user_input:
                continue
                
            # Handle conversation
            response = agent.handle_conversation(user_input)
            print(f"\nAgent: {response}")
            
        except KeyboardInterrupt:
            print("\n\nGoodbye!")
            break
        except Exception as e:
            print(f"Error: {e}")

def run_preset_tests():
    """
    Run preset test cases
    """
    print("=== Running Preset Tests ===")
    
    agent = create_test_agent()
    
    test_cases = [
        "我最近总是感到焦虑，怎么办？",
        "什么是神经可塑性？",
        "如何通过冥想缓解焦虑？",
        "大脑如何处理情绪？",
        "我想制定一个减压计划"
    ]
    
    for i, test_input in enumerate(test_cases, 1):
        print(f"\n--- Test {i}: {test_input} ---")
        try:
            response = agent.handle_conversation(test_input)
            print(f"Agent: {response}")
        except Exception as e:
            print(f"Error in test {i}: {e}")
    
    print("\n--- Memory State ---")
    short_term = agent.memory_manager.get_short_term_memory()
    print(f"Memory contains {len(short_term)} entries:")
    for entry in short_term:
        print(f"  {entry['role']}: {entry['content']}")

if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "--preset":
        run_preset_tests()
    else:
        run_test_conversation()