import os
import logging

from pydantic import SecretStr   

from langchain_community.chat_models import ChatTongyi

# 配置日志
logger = logging.getLogger(__name__)

# === 1. 初始化 Qwen 模型 ===
dashscope_api_key = os.getenv("DASHSCOPE_API_KEY")
if dashscope_api_key is None:
    logger.error("环境变量 DASHSCOPE_API_KEY 未设置")
    raise ValueError("环境变量 DASHSCOPE_API_KEY 未设置")

logger.info("DashScope API key loaded: " + dashscope_api_key)

qwen_model = ChatTongyi(
    model="qwen-turbo",  # 也可用 qwen-max, qwen-turbo 等
    api_key=dashscope_api_key,
    streaming=True
)


logger.info("ChatTongyi model initialized")