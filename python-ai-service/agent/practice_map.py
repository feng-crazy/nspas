import os

from deepagents import create_deep_agent

from langchain_community.chat_models import ChatTongyi

from .mcp_tools import mcp_tools
from .practice_mapping_prompt import PracticeMappingPrompt

# === 1. 初始化 Qwen 模型 ===
dashscope_api_key = os.getenv("DASHSCOPE_API_KEY")
if dashscope_api_key is None:
    raise ValueError("环境变量 DASHSCOPE_API_KEY 未设置")

qwen_model = ChatTongyi(
    model="qwen-plus-latest",  # 也可用 qwen-max, qwen-turbo 等
    api_key=dashscope_api_key,  # type: ignore
    streaming=True,
)

practice_map_agent = create_deep_agent(
    model=qwen_model,
    tools=mcp_tools,
    system_prompt=PracticeMappingPrompt.get_system_prompt()
)


