#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Conversation Agent for Neuro Guide
This agent handles user conversations and provides neuroscientific explanations
"""

import sys
import os
import json
from typing import List, Dict, Any, Optional

from langchain_core.runnables.base import Runnable
from langchain_core.runnables import RunnableLambda, RunnablePassthrough, RunnableBranch
from langchain_core.output_parsers import StrOutputParser
from langchain_core.messages import HumanMessage, AIMessage, SystemMessage


# Import memory manager
from memory.memory_manager import MemoryManager
from prompts.tool_prompt import ToolsPrompt

# Add the parent directory to the path to import tools
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# 添加日志功能
import logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def log_prompt(x):
    # x 通常是 PromptValue 或字符串，取决于前一步输出
    if hasattr(x, 'to_string'):
        prompt_str = x.to_string()
    else:
        prompt_str = str(x)
    logger.info("【发送给大模型的 Prompt】:\n" + prompt_str)
    return x  # 原样返回，不改变数据流

LogPromptRunnable = RunnableLambda(log_prompt)

class ConversationAgent:
    """Agent responsible for handling user conversations"""
    
    def __init__(self, llm: Any, tools: Optional[List[Any]] = None):
        """
        Initialize the conversation agent
        
        Args:
            llm: The language model to use for generating responses
            tools: List of tools available to the agent
        """
        self.llm = llm
        self.tools = tools or []
        # Initialize memory manager
        self.memory_manager = MemoryManager()
        # Create a chain using pipe operator with retriever
        self.chain: Runnable[Any, Any] = self._create_chain()
    
    def _create_chain(self) -> Runnable:
        """
        Create a chain using pipe operator for processing user input with retrieval and prompt routing
        
        Returns:
            A LangChain Runnable
        """
        # Import the neuroscience retriever
        from retrievers.neuroscience_retriever import NeuroscienceRetriever
        from prompts.neuroscience_prompt import NeurosciencePrompt
        
        # Create retriever instance
        retriever = NeuroscienceRetriever()
        
        prompt = NeurosciencePrompt()

        def format_docs(docs):
            return "\n".join([d.page_content for d in docs])
        
        # First retrieve context, then route to appropriate prompt
        chain = (
            {"knowledge_base": retriever | format_docs, "input": RunnablePassthrough()}
            | RunnableLambda(prompt.route_prompt)
            | RunnableLambda(lambda x: x["prompt"].format(**x))
            | LogPromptRunnable
            | self.llm
            | StrOutputParser()
        )
        
        return chain
    
    def handle_conversation(self, user_id: str, user_input: str, context: Optional[List[Dict[str, Any]]] = None) -> str:
        """
        Handle a conversation turn with the user
        
        Args:
            user_id: The user's ID
            user_input: The user's input message
            context: Previous conversation context
            
        Returns:
            The agent's response
        """
        try:
            # Get memory context before processing
            memory_context = self.memory_manager.get_combined_memory_context(user_id, user_input)
            
            # Combine with external context if provided
            full_context = memory_context
            if context:
                formatted_context = self._format_context(context)
                full_context = f"{memory_context}\n\n{formatted_context}" if memory_context else formatted_context
            
            # Check if we need to use any tools
            tool_response = self._check_and_use_tools(user_input)
            if tool_response:
                response = tool_response
            else:
                # Use the chain to process the input with memory context
                if full_context:
                    # Modify the input to include memory context
                    enhanced_input = f"{user_input}\n\nContext: {full_context}"
                    logger.info(f"用户输入: {user_input}")
                    logger.info(f"完整上下文: {full_context}")
                    response = self.chain.invoke(enhanced_input)
                else:
                    logger.info(f"用户输入: {user_input}")
                    response = self.chain.invoke(user_input)
            
            # Add to short-term memory
            self.memory_manager.add_to_short_term_memory(user_id, user_input, response)
            
            # Check if should save to long-term memory
            if self.memory_manager.should_save_to_long_term(user_input, response):
                self.memory_manager.add_to_long_term_memory(user_id, user_input, response)
                self.memory_manager.persist_memory()
            
            return response
        except Exception as e:
            # Fallback response in case of error
            return f"抱歉，处理您的请求时出现错误: {str(e)}"
    
    def _format_context(self, context: List[Dict[str, Any]]) -> str:
        """
        Format conversation context for the prompt
        
        Args:
            context: List of previous conversation turns
            
        Returns:
            Formatted context string
        """
        formatted = "Previous conversation:\n"
        for turn in context:
            role = turn.get("role", "unknown")
            message = turn.get("message", "")
            formatted += f"{role.capitalize()}: {message}\n"
        return formatted
    
    def _check_and_use_tools(self, user_input: str) -> Optional[str]:
        """
        Check if any tools should be used based on user input and use them if needed
        
        Args:
            user_input: The user's input message
            
        Returns:
            Tool response if a tool was used, None otherwise
        """
        # If no tools are available, return None
        if not self.tools:
            return None
        
        try:
            # Format tools information for the prompt
            tools_description = "\n".join([
                f"- {getattr(tool, 'name', type(tool).__name__)}: {getattr(tool, 'description', 'No description available')}" 
                for tool in self.tools
            ])
            
            # Create an enhanced prompt for tool detection with JSON response requirement
            tool_detection_prompt = ToolsPrompt.get_tool_usage_prompt(tools_description, user_input)
            
            # Get LLM decision
            logger.info(f"Checking tool usage for input: {user_input}")
            tool_decision = self.llm.invoke(tool_detection_prompt)
            
            # Parse the JSON response
            decision_data = json.loads(tool_decision.content.strip())
            
            logger.info(f"Tool decision: {decision_data}")
            
            if not decision_data.get("should_use_tool", False):
                logger.info("LLM decided not to use any tools")
                return None
            
            tool_name = decision_data.get("tool_name")
            tool_parameters = decision_data.get("tool_parameters", {})
            
            if not tool_name:
                logger.warning("LLM indicated tool usage but no tool name provided")
                return None
            
            # Find the appropriate tool
            selected_tool = None
            for tool in self.tools:
                tool_attr_name = getattr(tool, 'name', type(tool).__name__)
                if tool_attr_name == tool_name:
                    selected_tool = tool
                    break
            
            if not selected_tool:
                logger.warning(f"Requested tool '{tool_name}' not found in available tools")
                return None
            
            # Execute the tool with the provided parameters
            logger.info(f"Executing tool '{tool_name}' with parameters: {tool_parameters}")
            
            if hasattr(selected_tool, 'invoke'):
                result = selected_tool.invoke(tool_parameters)
            elif hasattr(selected_tool, '__call__'):
                result = selected_tool(tool_parameters)
            else:
                logger.error(f"Tool '{tool_name}' does not have a callable interface")
                return None
            
            logger.info(f"Tool execution successful")
            return f"Tool '{tool_name}' executed successfully: {result}"            
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing JSON response from LLM: {e}")
            # logger.error(f"Raw response: {tool_decision.content if 'tool_decision' in locals() else 'N/A'}")
            return None
        except Exception as e:
            logger.error(f"Error in tool checking or usage: {e}")
            # Fallback to normal conversation flow
            return None

# Example usage (for testing purposes)
if __name__ == "__main__":
    # Import and use actual tools
    from tools.practice_plan_tool import ngag_tools
    # For testing, we'll create a simple mock LLM
    class MockLLM:
        def invoke(self, prompt: str) -> str:
            return f"Mock response to: {prompt}"
    
    llm = MockLLM()
    tools = [
        ngag_tools
    ]
    agent = ConversationAgent(llm, tools)
    
    # Test conversation
    response = agent.handle_conversation("我最近总是感到焦虑，能帮我分析一下吗？")
    print(response)