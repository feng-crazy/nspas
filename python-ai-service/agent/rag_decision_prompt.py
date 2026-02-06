"""RAG决策提示模板模块"""

from typing import Optional


class RagDecisionPrompt:
    """RAG决策提示模板类"""

    @staticmethod
    def generate_prompt(query: str, retrieved_docs: str) -> str:
        """
        生成RAG决策相关的提示

        Args:
            query: 用户查询
            retrieved_docs: 检索到的文档

        Returns:
            格式化的提示字符串
        """
        return f"""
你是RAG决策专家，请基于以下检索到的文档回答用户查询：
检索文档：{retrieved_docs}
用户查询：{query}
要求：综合文档内容给出准确回答，若文档中无相关信息则说明无法回答。
        """.strip()


    @classmethod
    def get_system_prompt(cls) -> str:
        """
        获取系统级提示

        Returns:
            系统提示字符串
        """
        return "你是一个RAG决策专家，擅长根据检索到的信息为用户提供精准答案。"