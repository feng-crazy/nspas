"""修行小助手提示模板模块"""

from typing import Optional


class PracticeAssistantPrompt:
    """修行小助手提示模板类"""

    @staticmethod
    def generate_prompt(user_question: str, practice_context: str) -> str:
        """
        生成修行小助手相关的提示

        Args:
            user_question: 用户问题
            practice_context: 修行上下文

        Returns:
            格式化的提示字符串
        """
        return f"""
你是修行小助手，请根据以下修行上下文回答用户问题：
修行上下文：{practice_context}
用户问题：{user_question}
要求：以温和、鼓励的方式回答，体现修行智慧。
        """.strip()


    @classmethod
    def get_system_prompt(cls) -> str:
        """
        获取系统级提示

        Returns:
            系统提示字符串
        """
        return "你是一名神经科学的修行小助手，擅长基于神经科学，神经可塑性训练，大脑驯化机制来帮助用户创建自己的训练计划和训练工具。返回html格式的内容"