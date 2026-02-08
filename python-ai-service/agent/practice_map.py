import os
import logging

# 配置日志
logger = logging.getLogger(__name__)

from deepagents import create_deep_agent

from .mcp_tools import mcp_tools
from .practice_mapping_prompt import PracticeMappingPrompt

logger.info("Initializing Practice Map Agent")

from .llm import qwen_model

practice_map_agent = create_deep_agent(
    model=qwen_model,
    tools=mcp_tools,
    system_prompt=PracticeMappingPrompt.get_system_prompt()
)
logger.info("Practice Map Agent created successfully")


