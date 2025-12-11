#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Web App Generator Agent for Neuro Guide
This module implements an agent that generates web applications for practice plans
"""

from typing import List, Dict, Any

class WebAppGeneratorAgent:
    """Agent for generating web applications to display practice plans"""
    
    def __init__(self, llm):
        """
        Initialize the web app generator agent
        
        Args:
            llm: Language model to use for generating web applications
        """
        self.llm = llm
    
    def generate_web_app(self, title: str) -> str:
        """
        Generate a web application for displaying the practice plan using LLM
        
        Args:
            title: Title of the practice plan
            tasks: List of tasks in the practice plan
            
        Returns:
            HTML code for a web application
        """
        # Create a prompt for the LLM to generate HTML code
        prompt = f"""
        你是一个专业的前端开发工程师，擅长制作简洁美观的HTML页面。
        同时你也是一位基于神经科学的冥想和正念老师。
        请根据以下信息生成一个展示修行计划的HTML页面：
        
        标题: {title}
        包含：
        - 针对特定神经回路的日常修行
        - 每项修行的科学依据
        - 简单可行的操作指导
        
        要求:
        1. 页面应具有良好的视觉效果和用户体验
        2. 使用绿色系(#42b983)作为主色调，体现健康和成长主题
        3. 响应式设计，适配移动设备
        4. 每个任务应该有清晰的展示区域
        5. 科学依据部分应该以引用形式展示
        6. 返回纯净的HTML代码，不需要任何解释或其他内容
        7. HTML应包含必要的CSS样式，不要使用外部样式表
        8. 确保代码可以直接保存为.html文件并在浏览器中运行
        9. 以结构化的计划格式呈现，包含明确的每日任务。
        """
        
        # Generate HTML using LLM
        response = self.llm.invoke(prompt)
        return response.content.strip()
    
    @staticmethod
    def get_practice_plan_prompt() -> str:
        """
        Get the prompt template for generating practice plans
        
        Returns:
            Prompt template string
        """
        return """
        你是一位基于神经科学的冥想和正念老师。
        请根据以下信息生成个性化的修行计划：
        1. 用户的情绪状态: {emotional_state}
        2. 识别的神经机制: {neural_mechanisms}
        3. 用户的目标: {user_goals}
        4. 科学证据: {scientific_evidence}
        
        创建一个{days}天的计划，包含：
        - 针对特定神经回路的日常修行
        - 每项修行的科学依据
        - 简单可行的操作指导
        
        以结构化的计划格式呈现，包含明确的每日任务。
        """