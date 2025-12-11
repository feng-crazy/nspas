from typing import Dict, Any

class ToolsPrompt:
    """Class containing prompt templates for tool-related queries"""
    
    @staticmethod
    def get_tool_usage_prompt(tools_description: str, user_input: str) -> str:
        """
        Get the prompt template for requesting tool usage
        
        Returns:
            Prompt template string
        """
        return f"""
            你是一位神经科学助手，拥有以下工具可用：

            {tools_description}

            请分析用户的输入，确定：
            1. 是否应该使用工具
            2. 使用哪个工具（如果有）  
            3. 调用该工具的参数

            用户输入: {user_input}

            请以严格的JSON格式响应，包含以下字段：
            - should_use_tool: boolean (是否使用工具)
            - tool_name: string (工具名称，如果不使用工具则为null)
            - tool_parameters: object (工具参数，如果不使用工具则为null)
            - reasoning: string (你的推理过程)

            只返回JSON格式的响应，不要有其他内容。

            示例响应格式：
            {{
                "should_use_tool": true,
                "tool_name": "generate_web_app",
                "tool_parameters": {{"title": "7天正念减压修行计划"}},
                "reasoning": "用户明确要求生成网页应用，generate_web_app工具专门用于此目的"
            }}
            """
    
    @staticmethod
    def get_tool_response_prompt(user_input: str, tool_name: str, result: str) -> str:
        """
        Get the prompt template for incorporating tool response into final answer
        
        Returns:
            Prompt template string
        """
        return f"""
            你是一位神经科学顾问。刚才你使用了{tool_name}工具来处理用户的请求。

            用户输入: {user_input}
            工具执行结果: {result}

            请给用户提供一个完整的响应，解释：
            1. 你使用了什么工具来处理这个请求
            2. 工具生成了什么内容
            3. 这个内容如何帮助用户
            4. 接下来的建议步骤

            请提供友好、专业的回应，让用户感觉被理解和帮助。

            只返回最终的用户响应，不要有其他内容。
                        """