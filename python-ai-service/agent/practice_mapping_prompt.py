"""修行映射提示模板模块"""

from typing import Optional


class PracticeMappingPrompt:
    """修行映射提示模板类"""

    @staticmethod
    def generate_prompt(user_input: str, practice_data: str) -> str:
        """
        生成修行映射相关的提示

        Args:
            user_input: 用户输入内容
            practice_data: 修行相关数据

        Returns:
            格式化的提示字符串
        """
        return f"""
你是修行指导专家，请根据以下修行数据为用户提供指导：
修行数据：{practice_data}
用户输入：{user_input}
要求：结合修行数据给出针对性建议，语言温和且专业。
        """.strip()


    @classmethod
    def get_system_prompt(cls) -> str:
        """
        获取系统级提示

        Returns:
            系统提示字符串
        """
        return "你是一名学识渊博，熟悉各种哲学和修行语录以及心理学概念的神经科学的专家，擅长用神经科学角度解释各种修行道理和方法，比如从神经科学角度解释知行合一的道理。"