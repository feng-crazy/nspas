import os
import logging


# 配置日志
logger = logging.getLogger(__name__)

from langchain.chat_models import init_chat_model

from deepagents import create_deep_agent

from .mcp_tools import mcp_tools

from .practice_assistant_prompt import PracticeAssistantPrompt


logger.info("Initializing Practice Assistant Agent")

from .llm import qwen_model

practice_assistant_agent = create_deep_agent(
    model=qwen_model,
    tools=mcp_tools,
    system_prompt=PracticeAssistantPrompt.get_system_prompt()
)
logger.info("Practice Assistant Agent created successfully")


def __main__():
    logger.info("Running Practice Assistant Agent")
    result = practice_assistant_agent.invoke({"messages": [{"role": "user", "content": "What is langgraph?"}]})

    # Print the agent's response
    logger.info(f"Agent response: {result['messages'][-1].content}")