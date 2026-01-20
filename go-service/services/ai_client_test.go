package services

import (
	"context"
	"testing"

	"github.com/nspas/go-service/config"
	"github.com/stretchr/testify/assert"
)

func TestMockAIClient_Chat(t *testing.T) {
	// 创建mock客户端
	mockClient := NewMockAIClient()
	
	// 设置mock响应
	mockClient.SetMockResponse("Hello, world!")
	
	// 测试聊天功能
	messages := []Message{
		{Content: "Hello", IsUser: true},
	}
	
	response, err := mockClient.Chat(context.Background(), messages, "analysis")
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", response)
}

func TestMockAIClient_Chat_WithError(t *testing.T) {
	// 创建mock客户端
	mockClient := NewMockAIClient()
	
	// 设置mock错误
	mockClient.SetMockError(assert.AnError)
	
	// 测试聊天功能
	messages := []Message{
		{Content: "Hello", IsUser: true},
	}
	
	response, err := mockClient.Chat(context.Background(), messages, "analysis")
	
	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, response)
}

func TestHTTPClient_Chat(t *testing.T) {
	// 创建HTTP客户端
	cfg := &config.Config{
		PythonAI: config.PythonAIConfig{
			BaseURL: "http://localhost:5000",
		},
	}
	client := NewHTTPClient(cfg)
	
	// 测试聊天功能（这里会失败，因为没有运行Python AI服务，主要是为了测试代码结构）
	messages := []Message{
		{Content: "Hello", IsUser: true},
	}
	
	response, err := client.Chat(context.Background(), messages, "analysis")
	
	// 验证结果
	// 由于没有实际的AI服务，这里应该返回错误
	assert.Error(t, err)
	assert.Empty(t, response)
}
