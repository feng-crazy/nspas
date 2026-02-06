import os
from typing import Literal
from tavily import TavilyClient
import asyncio


from deepagents import create_deep_agent
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_community.chat_models import ChatZhipuAI
from langchain.messages import AIMessage, HumanMessage, SystemMessage
from .search_tools import internet_search, research_instructions
from .memory_tools import get_user_history, UserState
from .mcp_tools import mcp_tools
from .neuroscience_analysis_prompt import NeuroscienceAnalysisPrompt

tavily_client = TavilyClient(api_key=os.getenv("TAVILY_API_KEY"))

glm_streaming_chat = ChatZhipuAI(
    model="glm-4",
    api_key=os.getenv("ZHIPU_API_KEY"),
    temperature=0.5,
    streaming=True,
)

tools = [internet_search, get_user_history]
tools.extend(mcp_tools)

ns_analysis_agent = create_deep_agent(
    model=glm_streaming_chat,
    tools=tools,
    context_schema= UserState,
    system_prompt=NeuroscienceAnalysisPrompt.get_system_prompt()
)


if __name__ == "__main__":
    print("Running NS Analysis Agent")
    result = ns_analysis_agent.invoke({"messages": [{"role": "user", "content": "什么是神经科学？"}],
                                        "user_id": "user_123", "conversation_id": "conv_456"})
    print(result)
    # Print the agent's response
    print(result["messages"][-1].content)