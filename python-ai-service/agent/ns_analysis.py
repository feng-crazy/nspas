import os
import logging
from typing import Literal
from tavily import TavilyClient


# 配置日志
logger = logging.getLogger(__name__)

from deepagents import create_deep_agent
from .search_tools import internet_search
from .memory_tools import get_user_history, UserState
from .mcp_tools import mcp_tools
from .neuroscience_analysis_prompt import NeuroscienceAnalysisPrompt

logger.info("Initializing NS Analysis Agent")

tavily_client = TavilyClient(api_key=os.getenv("TAVILY_API_KEY"))
logger.info("Tavily client initialized")

# glm_streaming_chat = ChatZhipuAI(
#     model="glm-4",
#     api_key=os.getenv("ZHIPU_API_KEY"),
#     temperature=0.5,
#     streaming=True,
# )
# logger.info("ChatZhipuAI model initialized")

from .llm import qwen_model

tools = [internet_search, get_user_history]
tools.extend(mcp_tools)
logger.info(f"Tools loaded:")

ns_analysis_agent = create_deep_agent(
    model=qwen_model,
    tools=tools,
    context_schema= UserState,
    system_prompt=NeuroscienceAnalysisPrompt.get_system_prompt()
)
logger.info("NS Analysis Agent created successfully")


if __name__ == "__main__":
    logger.info("Running NS Analysis Agent")
    result = ns_analysis_agent.invoke({"messages": [{"role": "user", "content": "什么是神经科学？"}],
                                        "user_id": "user_123", "conversation_id": "conv_456"})
    logger.debug(f"Agent result: {result}")
    # Print the agent's response
    logger.info(f"Agent response: {result['messages'][-1].content}")