#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Test script to verify LLM connectivity and basic functionality
"""

import os
import sys
from dotenv import load_dotenv

# 设置TOKENIZERS_PARALLELISM环境变量以避免huggingface/tokenizers的警告
os.environ["TOKENIZERS_PARALLELISM"] = "false"

# Load environment variables from .env file if it exists
load_dotenv()

# Add the current directory to the path to import modules
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from langchain_openai import ChatOpenAI
from pydantic import SecretStr

def test_llm_connection():
    """
    Test basic LLM connection and response
    """
    print("Testing LLM connection...")
    
    # Get configuration from environment variables
    api_key = os.getenv("OPENAI_API_KEY")
    base_url = os.getenv("OPENAI_BASE_URL")
    model_name = os.getenv("MODEL_NAME", "glm-4")
    
    # Use default values if not provided
    if not api_key:
        api_key = "you-are-key"
        print("Using default API key")
    
    if not base_url:
        base_url = "https://open.bigmodel.cn/api/paas/v4"
        print(f"Using default base URL: {base_url}")
    
    print(f"Model name: {model_name}")
    
    # Initialize the LLM
    llm = ChatOpenAI(
        model=model_name,
        api_key=SecretStr(api_key),
        base_url=base_url,
        temperature=0.7,
        max_completion_tokens=1000
    )
    
    # Test prompt
    test_prompt = "请用一句话解释什么是神经可塑性"
    
    print(f"\nSending test prompt: {test_prompt}")
    
    try:
        # Send request to LLM
        response = llm.invoke(test_prompt)
        print("\nLLM Response:")
        print("-" * 50)
        print(response.content)
        print("-" * 50)
        print("✅ LLM test successful!")
        return True
        
    except Exception as e:
        print(f"❌ LLM test failed with error: {str(e)}")
        return False

def test_simple_prompt():
    """
    Test a simple prompt to verify basic functionality
    """
    print("\n" + "="*50)
    print("Testing simple prompt...")
    
    # Get configuration
    api_key = os.getenv("OPENAI_API_KEY")
    base_url = os.getenv("OPENAI_BASE_URL")
    model_name = os.getenv("MODEL_NAME", "glm-4")
    
    if not api_key:
        api_key = "you-are-key"
    
    if not base_url:
        base_url = "https://open.bigmodel.cn/api/paas/v4"
    
    # Initialize LLM
    llm = ChatOpenAI(
        model=model_name,
        api_key=SecretStr(api_key),
        base_url=base_url,
        temperature=0.3,  # Lower temperature for more deterministic response
        max_completion_tokens=500
    )
    
    # Simple test
    prompt = "Hello, are you available?"
    
    print(f"Prompt: {prompt}")
    
    try:
        response = llm.invoke(prompt)
        print("Response:", response.content)
        print("✅ Simple prompt test successful!")
        return True
    except Exception as e:
        print(f"❌ Simple prompt test failed: {str(e)}")
        return False

if __name__ == "__main__":
    print("Neuro Guide - LLM Connection Test")
    print("="*50)
    
    # Test basic connection
    success1 = test_llm_connection()
    
    # Test simple prompt
    success2 = test_simple_prompt()
    
    print("\n" + "="*50)
    if success1 and success2:
        print("🎉 All tests passed! LLM is working correctly.")
        sys.exit(0)
    else:
        print("💥 Some tests failed. Please check your configuration.")
        sys.exit(1)