from setuptools import setup, find_packages

setup(
    name="python-ai-service",
    version="0.1.0",
    description="AI服务API，支持聊天、流式聊天、RAG等功能",
    packages=find_packages(include=["config", "logger", "rag", "memory", "prompt", "agent", "app"]),
    install_requires=[
        "deepagents>=0.3.6",
        "langchain>=1.2.6",
        "langchain-community>=0.4.1",
        "langgraph>=1.0.6",
        "fastapi>=0.100.0",
        "uvicorn>=0.23.0",
        "pydantic-settings>=2.0.0",
        "python-dotenv>=1.0.0",
        "structlog>=23.2.0",
        "openai>=1.0.0",
        "langchain-chroma>=0.1.0",
        "langchain-text-splitters>=0.0.1",
        "rank_bm25>=0.2.2",
        "sentence-transformers>=2.2.0",
        "huggingface-hub>=0.16.0",
        "chromadb>=0.4.0",
        "requests>=2.31.0"
    ],
    extras_require={
        "test": [
            "pytest>=7.0.0",
            "pytest-asyncio>=0.21.0",
            "httpx>=0.24.0",
            "fastapi.testclient>=0.100.0"
        ]
    },
    entry_points={
        "console_scripts": [
            "python-ai-service=app:main",
        ],
    },
)
