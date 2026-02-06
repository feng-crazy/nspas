import os
from typing import Literal
from tavily import TavilyClient
from pydantic import SecretStr

from langchain.chat_models import init_chat_model

from deepagents import create_deep_agent

from .mcp_tools import mcp_tools

from .practice_assistant_prompt import PracticeAssistantPrompt



from langchain_community.chat_models import ChatTongyi

# === 1. 初始化 Qwen 模型 ===
qwen_model = ChatTongyi(
    model="qwen-plus-latest",  # 也可用 qwen-max, qwen-turbo 等
    api_key=SecretStr(os.getenv("DASHSCOPE_API_KEY")),
    streaming=True,
    # 注意：如果使用“思考模式”，需开启流式（见下文说明）
)

practice_assistant_agent = create_deep_agent(
    model=qwen_model,
    tools=mcp_tools,
    system_prompt=PracticeAssistantPrompt.get_system_prompt()
)


def __main__():
    result = practice_assistant_agent.invoke({"messages": [{"role": "user", "content": "What is langgraph?"}]})

    # Print the agent's response
    print(result["messages"][-1].content)