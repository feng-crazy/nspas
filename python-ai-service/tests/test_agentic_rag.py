import os
import sys
from pathlib import Path
from dotenv import load_dotenv

# 获取当前文件所在目录的父目录路径，构建.env文件的路径
current_dir = Path(__file__).parent
dotenv_path = current_dir.parent / ".env"
load_dotenv(dotenv_path=str(dotenv_path))

# 将项目根目录添加到Python路径中，以便能够导入agent模块
sys.path.insert(0, str(current_dir.parent))

from agent.rag_tools import agentic_rag_search, rag_tool_instance

def test_agentic_rag_search():
    """测试Agentic RAG搜索功能"""
    print("\n=== 测试Agentic RAG搜索功能 ===\n")
    
    # 测试用例1：需要检索的查询
    print("1. 测试需要检索的查询：")
    query1 = "什么是冥想？"
    print(f"查询：{query1}")
    try:
        # 模拟ToolRuntime
        class MockToolRuntime:
            def __init__(self):
                self.state = {"user_id": "test_user", "conversation_id": "test_conv"}
        
        result1 = agentic_rag_search(query1, MockToolRuntime())
        print(f"结果：{result1}")
        print(f"结果长度：{len(result1)}字符")
        print("✓ 测试通过")
    except Exception as e:
        print(f"✗ 测试失败：{str(e)}")
    
    # 测试用例2：不需要检索的查询
    print("\n2. 测试不需要检索的查询：")
    query2 = "你好，今天天气怎么样？"
    print(f"查询：{query2}")
    try:
        result2 = agentic_rag_search(query2, MockToolRuntime())
        print(f"结果：{result2}")
        print(f"结果长度：{len(result2)}字符")
        print("✓ 测试通过")
    except Exception as e:
        print(f"✗ 测试失败：{str(e)}")
    
    # 测试用例3：复杂查询
    print("\n3. 测试复杂查询：")
    query3 = "神经科学如何解释冥想对大脑的影响？"
    print(f"查询：{query3}")
    try:
        result3 = agentic_rag_search(query3, MockToolRuntime())
        print(f"结果：{result3}")
        print(f"结果长度：{len(result3)}字符")
        print("✓ 测试通过")
    except Exception as e:
        print(f"✗ 测试失败：{str(e)}")
    
    # 测试用例4：检查向量库状态
    print("\n4. 检查向量库状态：")
    try:
        # 检查向量库中文档数量
        collection = rag_tool_instance.vector_store._collection
        count = collection.count()
        print(f"向量库中文档数量：{count}")
        print("✓ 测试通过")
    except Exception as e:
        print(f"✗ 测试失败：{str(e)}")

if __name__ == "__main__":
    print("测试Agentic RAG工具")
    print("=" * 50)
    
    try:
        test_agentic_rag_search()
        print("\n" + "=" * 50)
        print("🎉 所有测试完成！")
    except KeyboardInterrupt:
        print("\n测试被中断")
    finally:
        # 停止文档监控器
        rag_tool_instance.stop()
