package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nspas/go-service/config"
)

// AIClient 定义AI服务客户端接口
type AIClient interface {
	// Chat 调用AI服务进行对话
	Chat(ctx context.Context, messages []Message, convType string) (string, error)
}

// Message 定义消息结构
type Message struct {
	Content string `json:"content"`
	IsUser  bool   `json:"is_user"`
}

// AIChatRequest 定义AI聊天请求结构
type AIChatRequest struct {
	Messages         []Message `json:"messages"`
	ConversationType string    `json:"conversation_type"`
}

// AIChatResponse 定义AI聊天响应结构
type AIChatResponse struct {
	Content string `json:"content"`
}

// HTTPClient 实现AIClient接口，用于实际调用python-ai-service
type HTTPClient struct {
	cfg     *config.Config
	client  *http.Client
}

// NewHTTPClient 创建一个新的HTTPClient实例
func NewHTTPClient(cfg *config.Config) *HTTPClient {
	return &HTTPClient{
		cfg:     cfg,
		client:  &http.Client{},
	}
}

// Chat 调用python-ai-service进行对话
func (c *HTTPClient) Chat(ctx context.Context, messages []Message, convType string) (string, error) {
	// 创建请求体
	reqBody := AIChatRequest{
		Messages:         messages,
		ConversationType: convType,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/chat", c.cfg.PythonAI.BaseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查响应状态码
	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI service returned status code: %d, body: %s", httpResp.StatusCode, string(respBody))
	}

	// 解析响应体
	var resp AIChatResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return resp.Content, nil
}

// MockAIClient 实现AIClient接口，用于mock AI服务调用
type MockAIClient struct {
	// MockResponse 用于设置mock响应
	MockResponse string
	// MockError 用于设置mock错误
	MockError error
}

// NewMockAIClient 创建一个新的MockAIClient实例
func NewMockAIClient() *MockAIClient {
	return &MockAIClient{}
}

// Chat 返回mock响应或错误
func (c *MockAIClient) Chat(ctx context.Context, messages []Message, convType string) (string, error) {
	if c.MockError != nil {
		return "", c.MockError
	}

	// 如果没有设置mock响应，返回默认响应
	if c.MockResponse == "" {
		return getDefaultResponse(convType), nil
	}

	return c.MockResponse, nil
}

// SetMockResponse 设置mock响应
func (c *MockAIClient) SetMockResponse(resp string) {
	c.MockResponse = resp
}

// SetMockError 设置mock错误
func (c *MockAIClient) SetMockError(err error) {
	c.MockError = err
}

// getDefaultResponse 返回默认响应，根据对话类型
func getDefaultResponse(convType string) string {
	switch convType {
	case "analysis":
		return "## 神经科学分析\n\n欢迎使用神经科学分析功能。请描述您的思维过程或行为，我将从神经科学角度为您分析。"
	case "mapping":
		return "## 修行映射\n\n欢迎使用修行映射功能。请输入修行语录或概念，我将为您映射到脑科学机制与神经通路。"
	case "assistant":
		return "## 修行小助手\n\n欢迎使用修行小助手。请描述您的需求，我将为您生成个性化的脑科学修行工具。"
	default:
		return "欢迎使用神经科学AI修行助手。请输入您的问题或需求。"
	}
}
