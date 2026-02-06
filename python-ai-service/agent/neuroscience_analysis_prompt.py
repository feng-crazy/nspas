"""神经科学分析提示模板模块"""

from typing import Optional


class NeuroscienceAnalysisPrompt:
    """神经科学分析提示模板类"""

    @staticmethod
    def generate_prompt(question: str, context: str) -> str:
        """
        生成神经科学分析相关的提示

        Args:
            question: 用户提出的问题
            context: 相关的上下文信息

        Returns:
            格式化的提示字符串
        """
        return f"""
你是神经科学分析专家，请根据以下上下文回答问题：
上下文：{context}
问题：{question}
要求：回答需基于上下文，准确简洁。
        """.strip()


    @classmethod
    def get_system_prompt(cls) -> str:
        """
        获取系统级提示

        Returns:
            系统提示字符串
        """
        return "你是一个专业的脑神经科学研究助手，擅长根据神经科学和脑机制来分析和解释用户的问题。"