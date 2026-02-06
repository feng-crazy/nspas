"""AI服务主应用"""
import json
from typing import List, Dict, Any

from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import StreamingResponse
import uvicorn

from dotenv import load_dotenv
load_dotenv(dotenv_path='.env')

app = FastAPI(title="AI Service API")

from agent.ns_analysis import ns_analysis_agent
from agent.practice_map import practice_map_agent
from agent.practice_assistant import practice_assistant_agent




# Agent注册表
AGENTS = {
    "analysis": ns_analysis_agent,
    "mapping": practice_map_agent,
    "assistant": practice_assistant_agent
}

# 消息结构
class Message:
    """消息结构"""
    def __init__(self, content: str, is_user: bool):
        self.content = content
        self.is_user = is_user

# 转换消息格式
def convert_messages(messages: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """转换消息格式为LangChain格式"""
    converted = []
    for msg in messages:
        role = "user" if msg.get("is_user", False) else "assistant"
        converted.append({
            "role": role,
            "content": msg.get("content", "")
        })
    return converted

# 主聊天接口
@app.post("/ai/chat")
async def chat_endpoint(request: Request):
    """AI聊天接口，与go-service兼容"""
    try:
        # 解析请求
        req_data = await request.json()
        user_id = req_data.get("user_id")
        conversation_id = req_data.get("conversation_id")
        messages = req_data.get("messages", [])
        conversation_type = req_data.get("conversation_type", "analysis")
        
        # 验证参数
        if not user_id or not conversation_id:
            raise HTTPException(status_code=400, detail="Missing required parameters")
        
        # 获取对应的agent
        agent = AGENTS.get(conversation_type)
        if not agent:
            raise HTTPException(status_code=400, detail=f"Invalid conversation type: {conversation_type}")
        
        # 转换消息格式
        converted_messages = convert_messages(messages)
        
        # 构建输入数据
        input_data = {
            "messages": converted_messages,
            "user_id": user_id,
            "conversation_id": conversation_id
        }
        
        # 检查是否需要流式响应
        if request.headers.get("Accept") == "text/event-stream":
            # 流式响应
            async def event_generator():
                try:
                    # 调用agent的流式接口
                    async for chunk in agent.astream(input_data):
                        # 提取内容
                        if "model" in chunk:
                            ai_message = chunk["model"]["messages"][-1]
                            if hasattr(ai_message, 'content') and ai_message.content:
                                # 构建SSE事件
                                event_data = {
                                    "content": ai_message.content,
                                    "full_content": ai_message.content,
                                    "conversation_id": conversation_id,
                                    "messages": messages
                                }
                                yield f"data: {json.dumps(event_data)}\n\n"
                except Exception as e:
                    error_data = {
                        "error": str(e)
                    }
                    yield f"data: {json.dumps(error_data)}\n\n"
            
            # 返回SSE响应
            return StreamingResponse(
                event_generator(),
                media_type="text/event-stream",
                headers={
                    "Cache-Control": "no-cache",
                    "Connection": "keep-alive"
                }
            )
        else:
            # 非流式响应
            result = await agent.ainvoke(input_data)
            response = {
                "content": result["messages"][-1].content,
                "conversation_id": conversation_id,
                "messages": messages
            }
            return response
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# 健康检查接口
@app.get("/health")
async def health_check():
    """健康检查接口"""
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
