def test_app_health():
    """测试FastAPI应用健康检查端点"""
    from main import app
    from fastapi.testclient import TestClient
    
    client = TestClient(app)
    response = client.get("/health")
    
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}