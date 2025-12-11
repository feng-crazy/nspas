#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Neuro Guide Python AI Service
This service implements the AI layer using LangChain and LlamaIndex
"""

import os
# 设置TOKENIZERS_PARALLELISM环境变量以避免huggingface/tokenizers的警告
os.environ["TOKENIZERS_PARALLELISM"] = "false"

from fastapi import FastAPI
from pydantic import BaseModel, SecretStr
from typing import List, Dict, Any, Optional

# Initialize FastAPI app
app = FastAPI(
    title="Neuro Guide Python AI Service",
    description="AI layer for Neuro Guide application",
    version="0.1.0"
)

class HealthResponse(BaseModel):
    status: str
    message: str

class ChatRequest(BaseModel):
    user_id: str
    message: str
    context: Optional[List[Dict[str, Any]]] = None

class ChatResponse(BaseModel):
    response: str

@app.get("/")
async def root():
    return {"message": "Welcome to Neuro Guide Python AI Service"}

@app.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(
        status="ok",
        message="Python AI service is running"
    )

# Import our conversation agent
from agents.conversation_agent import ConversationAgent

# Import OpenAI for standard interface
from langchain_openai import ChatOpenAI

# Initialize the LLM with OpenAI standard interface
# This supports multiple providers like ZhipuAI, Qwen, etc.
# You need to set the appropriate environment variables:
# For ZhipuAI: OPENAI_API_KEY and OPENAI_BASE_URL
# For Qwen: OPENAI_API_KEY (and optionally OPENAI_BASE_URL if using custom endpoint)
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

# Import tools
from tools.practice_plan_tool import ngag_tools
# Create tools
tools = [
    ngag_tools
]

# Create the agent
agent = ConversationAgent(llm, tools)

@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """
    Handle chat requests with memory support
    """
    try:
        response = agent.handle_conversation(request.user_id, request.message, request.context)
        return ChatResponse(response=response)
    except Exception as e:
        return ChatResponse(response=f"Error processing request: {str(e)}")

@app.post("/chat/clear_memory")
async def clear_memory(user_id: str):
    """
    Clear conversation memory for a specific user
    """
    try:
        agent.memory_manager.clear_short_term_memory(user_id)
        return {"status": "success", "message": "Memory cleared"}
    except Exception as e:
        return {"status": "error", "message": str(e)}

@app.get("/chat/memory")
async def get_memory(user_id: str):
    """
    Get current memory state for a specific user
    """
    try:
        short_term = agent.memory_manager.get_short_term_memory(user_id)
        return {
            "short_term_memory": short_term,
            "memory_size": len(short_term)
        }
    except Exception as e:
        return {"status": "error", "message": str(e)}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)