#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Neuroscience Prompt Templates for Neuro Guide
This module contains improved prompt templates and routing for neuroscience-related queries
"""

from typing import Dict, Any, Optional
from enum import Enum
import re
import os

from langchain_core.prompts import PromptTemplate


class IntentType(Enum):
    """Enumeration of possible intents"""
    CONCEPT_EXPLANATION = "concept_explanation"
    SYMPTOM_ANALYSIS = "symptom_analysis"
    DEFAULT = "default"


class NeurosciencePrompt:
    """Class containing improved prompt templates and routing logic for neuroscience-related queries"""
    
    # Keywords for intent classification
    SYMPTOM_KEYWORDS = {
        "焦虑", "抑郁", "失眠", "压力", "紧张", "惊恐", "强迫", "创伤", 
        "恐惧", "疑病", "躯体化", "情绪波动", "注意力不集中", "记忆减退",
        "心悸", "头晕", "头痛", "疲劳", "麻木", "疼痛", "恐慌", "躁狂",
        "抽搐", "震颤", "癫痫", "痴呆", "幻觉", "妄想", "多动", "自闭",
        "焦虑症", "抑郁症", "失眠症", "恐惧症", "强迫症", "创伤后应激",
        "PTSD", "ADHD", "OCD", "bipolar", "autism", "schizophrenia",
        "panic", "stress", "fatigue", "insomnia", "depression", "anxiety"
    }
    
    CONCEPT_KEYWORDS = {
        "正念", "禅", "冥想", "修行", "无我", "知行合一", "意识", "自由意志",
        "神经可塑性", "默认模式网络", "前额叶", "杏仁核", "海马体", "皮质醇",
        "血清素", "多巴胺", "GABA", "内啡肽", "催产素", "神经递质",
        "神经回路", "脑区", "脑波", "alpha波", "theta波", "gamma波",
        "冥想科学", "正念减压", "MBSR", "MBCT", "慈悲冥想", "专注冥想",
        "觉知", "内观", "禅修", "呼吸法", "身体扫描", "观息法",
        "mindfulness", "meditation", "neuroplasticity", "default mode network",
        "prefrontal cortex", "amygdala", "serotonin", "dopamine", "awareness"
    }

    def __init__(self, use_hf_api: bool = False, hf_api_token: Optional[str] = None):
        """
        Initialize the NeurosciencePrompt class
        
        Args:
            use_hf_api: Whether to use HuggingFace Inference API for classification
            hf_api_token: HuggingFace API token (required if use_hf_api=True)
        """
        self.use_hf_api = use_hf_api
        self.hf_api_token = hf_api_token or os.getenv("HF_API_TOKEN")
        
        if self.use_hf_api and not self.hf_api_token:
            raise ValueError("HuggingFace API token is required when use_hf_api=True. "
                           "Either pass hf_api_token or set HF_API_TOKEN environment variable.")
        
        if self.use_hf_api:
            try:
                import requests
                self.requests = requests
            except ImportError:
                raise ImportError("requests package is required for HuggingFace API. "
                                "Install with: pip install requests")

    @staticmethod
    def get_concept_explanation_prompt() -> str:
        """
        Get the prompt template for explaining neuroscience concepts
        
        Returns:
            Prompt template string
        """
        return """
        你是一位融合神经科学与东方修行传统的学者。请用以下结构回答：

        【神经机制】
        - 涉及的核心脑区（如默认模式网络、前额叶皮层等）
        - 关键神经递质或激素（如 GABA、血清素、皮质醇）
        - 相关神经可塑性变化

        【实践联结】
        - 如何通过正念/冥想/日常练习影响上述机制？
        - 提供一个可操作的小练习（<50字）

        【注意】避免诊断或治疗建议。若用户描述明显症状，请建议寻求专业帮助。

        用户问题：{user_input}
        上下文：{context}
        知识库：{knowledge_base}
        """.strip()

    @staticmethod
    def get_symptom_analysis_prompt() -> str:
        """
        Get the prompt template for analyzing symptoms from a neuroscience perspective
        
        Returns:
            Prompt template string
        """
        return """
        你是一位临床神经科学家。请按以下框架分析：

        1. 【症状映射】将用户描述的症状映射到已知神经精神综合征（如广泛性焦虑、PTSD 等）
        2. 【神经回路】说明失调的脑网络（如杏仁核-前额叶通路、HPA 轴）
        3. 【干预证据】列出 2-3 项有实证支持的非药物干预（如 HRV 生物反馈、正念认知疗法）
        4. 【警示】明确说明何时应转介临床医生

        用户输入：{user_input}
        上下文：{context}
        知识库：{knowledge_base}
            """.strip()

    def classify_intent_hf_api(self, user_input: str) -> IntentType:
        """
        Classify intent using HuggingFace Inference API
        
        Args:
            user_input: The user's input string
            
        Returns:
            IntentType enum value
        """
        api_url = "https://api-inference.huggingface.co/models/facebook/bart-large-mnli"
        
        payload = {
            "inputs": user_input,
            "parameters": {
                "candidate_labels": ["symptom analysis", "concept explanation"]
            }
        }
        
        headers = {
            "Authorization": f"Bearer {self.hf_api_token}",
            "Content-Type": "application/json"
        }
        
        try:
            response = self.requests.post(api_url, json=payload, headers=headers)
            response.raise_for_status()
            result = response.json()
            
            # Get the label with highest score
            top_label = result["labels"][0]
            
            if "symptom" in top_label.lower():
                return IntentType.SYMPTOM_ANALYSIS
            else:
                return IntentType.CONCEPT_EXPLANATION
                
        except Exception as e:
            print(f"HuggingFace API error: {e}. Falling back to keyword classification.")
            return self.classify_intent_keyword(user_input)

    def classify_intent_keyword(self, user_input: str) -> IntentType:
        """
        Classify the intent of the user input using keyword matching.
        
        Args:
            user_input: The user's input string
            
        Returns:
            IntentType enum value
        """
        input_lower = user_input.lower()
        
        # Check for symptom keywords
        has_symptom = any(keyword in input_lower for keyword in NeurosciencePrompt.SYMPTOM_KEYWORDS)
        
        # Check for concept keywords
        has_concept = any(keyword in input_lower for keyword in NeurosciencePrompt.CONCEPT_KEYWORDS)
        
        # Determine intent based on keyword matches
        if has_symptom and not has_concept:
            return IntentType.SYMPTOM_ANALYSIS
        elif has_concept:
            return IntentType.CONCEPT_EXPLANATION
        else:
            # Default to concept explanation for safety
            return IntentType.CONCEPT_EXPLANATION

    def classify_intent(self, user_input: str) -> IntentType:
        """
        Classify the intent of the user input using either HuggingFace API or keyword matching
        
        Args:
            user_input: The user's input string
            
        Returns:
            IntentType enum value
        """
        if self.use_hf_api:
            return self.classify_intent_hf_api(user_input)
        else:
            return self.classify_intent_keyword(user_input)

    @staticmethod
    def get_prompt_template(intent: IntentType) -> PromptTemplate:
        """
        Get the appropriate prompt template based on intent
        
        Args:
            intent: The classified intent
            
        Returns:
            PromptTemplate object
        """
        if intent == IntentType.SYMPTOM_ANALYSIS:
            return PromptTemplate.from_template(
                NeurosciencePrompt.get_symptom_analysis_prompt()
            )
        else:
            return PromptTemplate.from_template(
                NeurosciencePrompt.get_concept_explanation_prompt()
            )

    def route_prompt(self, input_dict: Dict[str, Any]) -> Dict[str, Any]:
        """
        Route to the appropriate prompt based on input classification
        
        Args:
            input_dict: Dictionary containing 'input', 'context', and 'knowledge_base'
            
        Returns:
            Dictionary with prompt template and filled values
        """
        user_input = input_dict.get("input", "")
        context = input_dict.get("context", "")
        knowledge_base = input_dict.get("knowledge_base", "")
        
        intent = self.classify_intent(user_input)
        prompt_template = NeurosciencePrompt.get_prompt_template(intent)
        
        return {
            "prompt": prompt_template,
            "intent": intent.value,
            "user_input": user_input,
            "context": context,
            "knowledge_base": knowledge_base,
        }

    @staticmethod
    def get_all_prompts() -> Dict[str, str]:
        """
        Get all available prompt templates as a dictionary
        
        Returns:
            Dictionary mapping prompt names to template strings
        """
        return {
            "concept_explanation": NeurosciencePrompt.get_concept_explanation_prompt(),
            "symptom_analysis": NeurosciencePrompt.get_symptom_analysis_prompt()
        }


# Example usage
if __name__ == "__main__":
    # Example 1: Using keyword-based classification (default)
    print("=== Using Keyword Classification ===")
    prompt_templates_keyword = NeurosciencePrompt(use_hf_api=False)
    
    test_inputs = [
        {"input": "什么是正念冥想的神经机制？", "context": "", "knowledge_base": ""},
        {"input": "我总是焦虑失眠怎么办？", "context": "", "knowledge_base": ""},
        {"input": "神经可塑性是什么意思？", "context": "", "knowledge_base": ""},
        {"input": "我经常心悸和紧张", "context": "", "knowledge_base": ""},
    ]
    
    for i, test_input in enumerate(test_inputs):
        print(f"\n--- Test Case {i+1} ---")
        print(f"Input: {test_input['input']}")
        
        result = prompt_templates_keyword.route_prompt(test_input)
        print(f"Intent: {result['intent']}")
        print(f"Prompt template preview:\n{result['prompt'].template[:200]}...")
    
    # Example 2: Using HuggingFace API (requires API token)
    print("\n\n=== Using HuggingFace API Classification ===")
    try:
        # This will only work if you have a valid HF_API_TOKEN set
        prompt_templates_hf = NeurosciencePrompt(
            use_hf_api=True,
            hf_api_token=os.getenv("HF_API_TOKEN")  # Will use environment variable
        )
        
        # Test with the same inputs
        for i, test_input in enumerate(test_inputs):
            print(f"\n--- HF API Test Case {i+1} ---")
            print(f"Input: {test_input['input']}")
            
            result = prompt_templates_hf.route_prompt(test_input)
            print(f"Intent: {result['intent']}")
            
    except ValueError as e:
        print(f"Note: {e}")
        print("To use HuggingFace API, set your HF_API_TOKEN environment variable or pass it as parameter.")
        print("Get your token from: https://huggingface.co/settings/tokens")