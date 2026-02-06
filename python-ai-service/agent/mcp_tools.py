
import os
import asyncio
from langchain_mcp_adapters.client import MultiServerMCPClient
# from langchain_mcp_adapters.sessions import StdioConnection, SSEConnection

async def get_mcp_tools():
    client = MultiServerMCPClient(
        {
            "memory": {
                "command": "npx",
                "args": ["-y", "@modelcontextprotocol/server-memory"],
                "transport": "stdio",
            }
        }
    )
    mcp_tools = await client.get_tools()
    # mcp_prompt =  await client.get_prompt()
    return mcp_tools

mcp_tools = asyncio.run(get_mcp_tools())