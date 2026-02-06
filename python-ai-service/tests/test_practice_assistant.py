import os
import sys
from pathlib import Path
from unittest.mock import patch, MagicMock


from dotenv import load_dotenv

# 获取当前文件所在目录的父目录路径，构建.env文件的路径
current_dir = Path(__file__).parent
dotenv_path = current_dir.parent / ".env"
load_dotenv(dotenv_path=str(dotenv_path))

from agent.practice_assistant import practice_assistant_agent

if __name__ == "__main__":
    print("test practice_assistant_agent")
    result = practice_assistant_agent.invoke({"messages": [{"role": "user", "content": "什么是知行合一"}]})
    print(result)
    # Print the agent's response
    print(result["messages"][-1].content)