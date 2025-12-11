#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Practice Plan Tool for Neuro Guide
This module implements a LangChain tool for generating neuroscience-based practice plans
by connecting to an MCP server
"""

import sys
import os
# Add the parent directory to the path to import modules
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


from langchain_mcp_adapters.client import MultiServerMCPClient  


ngag_client = MultiServerMCPClient(  
    {
        "Neuro Guide APP Generator  MCP Server": {
            "transport": "stdio",  # Local subprocess communication
            "command": "python",
            # Absolute path to your math_server.py file
            "args": ["/Users/hedengfeng/workspace/neuro-guide/python-ai-service/tools/mcp/mcp_server.py"],
        },
    }
)

# Get tools from the MCP server
ngag_tools = None

async def initialize_ngag_tools():
    global ngag_tools
    ngag_tools = await ngag_client.get_tools()