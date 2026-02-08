import os
import logging
from typing import Literal
from pathlib import Path
from tavily import TavilyClient

# 配置日志
logger = logging.getLogger(__name__)

logger.info("Initializing search tools")

# 初始化 Tavily 客户端
tavily_api_key = os.getenv("TAVILY_API_KEY")
if tavily_api_key is None:
    logger.warning("TAVILY_API_KEY not set, internet search may not work")
else:
    logger.info("Tavily API key loaded")

tavily_client = TavilyClient(api_key=tavily_api_key)
logger.info("Tavily client initialized")

def internet_search(
    query: str,
    max_results: int = 5,
    topic: Literal["general", "news", "finance"] = "general",
    include_raw_content: bool = False,
):
    """Run a web search"""
    try:
        logger.info(f"Performing internet search: query='{query}', max_results={max_results}, topic={topic}")
        results = tavily_client.search(
            query,
            max_results=max_results,
            include_raw_content=include_raw_content,
            topic=topic,
        )
        logger.info(f"Search completed successfully, found {len(results.get('results', []))} results")
        return results
    except Exception as e:
        logger.error(f"Error performing internet search: {str(e)}")
        raise


# System prompt to steer the agent to be an expert researcher
research_instructions = """You are an expert researcher. Your job is to conduct thorough research and then write a polished report.

You have access to an internet search tool as your primary means of gathering information.

## `internet_search`

Use this to run an internet search for a given query. You can specify the max number of results to return, the topic, and whether raw content should be included.
"""