import os
import sys
from pathlib import Path
from unittest.mock import patch, MagicMock


from dotenv import load_dotenv

# 获取当前文件所在目录的父目录路径，构建.env文件的路径
current_dir = Path(__file__).parent
dotenv_path = current_dir.parent / ".env"
load_dotenv(dotenv_path=str(dotenv_path))

# 将项目根目录添加到Python路径中，以便能够导入agent模块
sys.path.insert(0, str(current_dir.parent))

from agent.ns_analysis import ns_analysis_agent

if __name__ == "__main__":
    print("test NS Analysis Agent")
    # 流式调用（关键：使用 .stream()）
    messages = [{"role": "user", "content": "简要回答什么是神经科学？20字以内"}]

    print("🤖 Agent is responding (streaming):\n")

    # 使用 .stream() 获取每一步的增量更新
    for chunk in ns_analysis_agent.stream({"messages": messages, 
                                           "user_id": "user_123", 
                                           "conversation_id": "conv_456"}):
        # chunk 是一个 dict，可能包含 'agent', 'tools', 'subagent' 等节点的输出
        # 我们只关心最终 AI 的回复（通常在 'agent' 节点中）
        if "model" in chunk:
            ai_message = chunk["model"]["messages"][-1]
            if hasattr(ai_message, 'content') and ai_message.content:
                # Token 级别的流式内容会逐步追加
                print(ai_message.content, end="", flush=True)
        # 你也可以打印工具调用等中间步骤用于调试
        #print(chunk)
        #print("\n-----------\n")

    print("\n\n✅ Done.")