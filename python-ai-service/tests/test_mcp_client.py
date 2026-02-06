import os
from typing import Literal
import asyncio
from deepagents import create_deep_agent
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_mcp_adapters.sessions import StdioConnection, SSEConnection
from langchain_community.chat_models import ChatZhipuAI

from tavily import TavilyClient


from dotenv import load_dotenv

load_dotenv(dotenv_path='../.env')

tavily_client = TavilyClient(api_key=os.getenv("TAVILY_API_KEY"))

glm_streaming_chat = ChatZhipuAI(
    model="glm-4",
    api_key=os.getenv("ZHIPU_API_KEY"),
    temperature=0.5,
    streaming=True,
)

async def main():
    # 创建正确的Connection对象
    math_connection: StdioConnection = {
        "transport": "stdio",
        "command": "python",
        # Absolute path to your math_server.py file
        "args": ["/path/to/math_server.py"],
    }
    
    weather_connection: SSEConnection = {
        "transport": "sse",
        "url": "http://localhost:8000/mcp",
    }
    
    client = MultiServerMCPClient(
        connections={
            "math": math_connection,
            "weather": weather_connection,
        }
    )

    tools = await client.get_tools()  
    agent = create_deep_agent(
        model=glm_streaming_chat,
        tools=tools
    )
    math_response = await agent.ainvoke(
        {"messages": [{"role": "user", "content": "what's (3 + 5) x 12?"}]}
    )
    weather_response = await agent.ainvoke(
        {"messages": [{"role": "user", "content": "what is the weather in nyc?"}]}
    )
    print(math_response)
    print(weather_response)

if __name__ == "__main__":
    asyncio.run(main())