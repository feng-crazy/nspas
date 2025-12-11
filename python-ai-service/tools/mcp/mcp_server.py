#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
MCP Server for Neuro Guide Practice Plan Generation
This module implements an MCP server that uses LLM to generate neuroscience-based practice plans
"""

import os
from typing import List, Dict, Any
from mcp.server.fastmcp import FastMCP
from pydantic import SecretStr
from langchain_openai import ChatOpenAI

# Import the web app generator agent
from web_app_agent import WebAppGeneratorAgent


# Initialize the MCP server
ngag_mcp = FastMCP("Neuro Guide APP Generator  MCP Server")

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

@ngag_mcp.tool()
def generate_web_app(title: str) -> str:
    """
    Generate a web application for displaying a practice plan
    
    Args:
        title: Title of the practice plan
        
    Returns:
        HTML code for a web application
    """
    return WebAppGeneratorAgent(llm).generate_web_app(title)



if __name__ == "__main__":
    ngag_mcp.run(transport="stdio")