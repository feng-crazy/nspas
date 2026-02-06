import os
import requests
import structlog
from dotenv import load_dotenv
from langchain.agents import create_agent, AgentState
from langchain.tools import tool, ToolRuntime

# 加载环境变量
load_dotenv()

# 配置日志
logger = structlog.get_logger(__name__)

memory_instructions = """You have access to a tool that allows you to look up the user's conversation history."""

# 获取GO_SERVICE_URL配置
GO_SERVICE_URL = os.getenv("GO_SERVICE_URL", "http://localhost:8080")
GO_SERVICE_TOKEN = os.getenv("GO_SERVICE_TOKEN", "sk-asdfghjklqwertyuiopzxcvbnm")

class UserState(AgentState):
    user_id: str
    conversation_id: str

@tool
def get_user_history(
    runtime: ToolRuntime
) -> str:
    """Look up user's conversation history."""
    user_id = runtime.state["user_id"]
    conversation_id = runtime.state["conversation_id"]
    
    logger.info("获取用户历史记录", user_id=user_id, conversation_id=conversation_id)
    
    # 构建API URL
    go_service_url = GO_SERVICE_URL + "/api/memory/" + conversation_id
    
    try:
        # 调用go-service的/memory/:conversation_id接口获取历史记录
        headers = {
            "User-ID": user_id,
            # 微服务之间认证建议：
            # 使用JWT Token：添加Authorization: Bearer <token>请求头
            "Authorization": "Bearer " + GO_SERVICE_TOKEN
        }
        response = requests.get(go_service_url, headers=headers, timeout=10)
        
        # 检查响应状态码
        response.raise_for_status()
        
        # 解析响应数据
        history_data = response.json()
        
        logger.info("成功获取用户历史记录", user_id=user_id, conversation_id=conversation_id, history_length=len(history_data.get("messages", [])))
        
        # 将历史记录格式化为字符串
        history_str = ""
        for message in history_data.get("messages", []):
            role = message.get("role", "unknown")
            content = message.get("content", "")
            history_str += f"{role}: {content}\n"
        
        return history_str
    except requests.exceptions.RequestException as e:
        logger.error("获取用户历史记录失败", user_id=user_id, conversation_id=conversation_id, error=str(e))
        return "获取历史记录失败"
    except ValueError as e:
        logger.error("解析历史记录失败", user_id=user_id, conversation_id=conversation_id, error=str(e))
        return "解析历史记录失败"



