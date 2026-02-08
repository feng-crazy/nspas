import os
import asyncio
import logging
from langchain_mcp_adapters.client import MultiServerMCPClient
# from langchain_mcp_adapters.sessions import StdioConnection, SSEConnection

# 配置日志
logger = logging.getLogger(__name__)

# 全局变量存储 MCP 工具
mcp_tools = None

async def get_mcp_tools():
    """获取 MCP 工具列表"""
    try:
        logger.info("Creating MultiServerMCPClient for memory server")
        client = MultiServerMCPClient(
            {
                "memory": {
                    "command": "npx",
                    "args": ["-y", "@modelcontextprotocol/server-memory"],
                    "transport": "stdio",
                }
            }
        )
        logger.info("Getting MCP tools from client")
        mcp_tools = await client.get_tools()
        logger.info(f"Retrieved {len(mcp_tools)} MCP tools")
        # mcp_prompt =  await client.get_prompt()
        return mcp_tools
    except Exception as e:
        logger.error(f"Error getting MCP tools: {str(e)}")
        raise

async def initialize_mcp_tools():
    """初始化 MCP 工具，在应用启动时调用"""
    global mcp_tools
    try:
        logger.info("Initializing MCP tools...")
        mcp_tools = await get_mcp_tools()
        logger.info(f"MCP tools initialized successfully, loaded {len(mcp_tools)} tools")
        return mcp_tools
    except Exception as e:
        logger.error(f"Error initializing MCP tools: {str(e)}")
        raise
