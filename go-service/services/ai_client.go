package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
)

// AIClient 定义AI服务客户端接口
type AIClient interface {
	// Chat 调用AI服务进行对话
	Chat(ctx context.Context, messages []Message, convType string) (string, error)
	// StreamChat 流式调用AI服务进行对话
	StreamChat(ctx context.Context, messages []Message, convType string) (<-chan string, error)
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
	logger.Info(ctx, "AI chat request started",
		slog.String("conversation_type", convType),
		slog.Int("message_count", len(messages)))

	// 创建请求体
	reqBody := AIChatRequest{
		Messages:         messages,
		ConversationType: convType,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error(ctx, "Failed to marshal request body", slog.Any("error", err))
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	logger.Debug(ctx, "AI chat request body", 
		slog.String("body", string(jsonData)))

	// 创建HTTP请求
	reqURL := fmt.Sprintf("%s/chat", c.cfg.PythonAI.BaseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		reqURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		logger.Error(ctx, "Failed to create HTTP request", 
			slog.String("url", reqURL),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	logger.Info(ctx, "Sending request to AI service", 
		slog.String("url", reqURL))
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		logger.Error(ctx, "Failed to send request to AI service", 
			slog.String("url", reqURL),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		logger.Error(ctx, "Failed to read AI service response", 
			slog.String("url", reqURL),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查响应状态码
	if httpResp.StatusCode != http.StatusOK {
		logger.Error(ctx, "AI service returned error status", 
			slog.String("url", reqURL),
			slog.Int("status_code", httpResp.StatusCode),
			slog.String("response", string(respBody)))
		return "", fmt.Errorf("AI service returned status code: %d, body: %s", httpResp.StatusCode, string(respBody))
	}

	// 解析响应体
	var resp AIChatResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		logger.Error(ctx, "Failed to unmarshal AI service response", 
			slog.String("url", reqURL),
			slog.String("response", string(respBody)),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	logger.Info(ctx, "AI chat request completed successfully", 
		slog.String("conversation_type", convType),
		slog.String("response", resp.Content[:50]+"..."))

	return resp.Content, nil
}

// StreamChat 流式调用python-ai-service进行对话
func (c *HTTPClient) StreamChat(ctx context.Context, messages []Message, convType string) (<-chan string, error) {
	logger.Info(ctx, "AI stream chat request started",
		slog.String("conversation_type", convType),
		slog.Int("message_count", len(messages)))

	// 创建请求体
	reqBody := AIChatRequest{
		Messages:         messages,
		ConversationType: convType,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error(ctx, "Failed to marshal request body", slog.Any("error", err))
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	logger.Debug(ctx, "AI stream chat request body", 
		slog.String("body", string(jsonData)))

	// 创建HTTP请求
	reqURL := fmt.Sprintf("%s/stream-chat", c.cfg.PythonAI.BaseURL)
	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		reqURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		logger.Error(ctx, "Failed to create HTTP request", 
			slog.String("url", reqURL),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	logger.Info(ctx, "Sending stream request to AI service", 
		slog.String("url", reqURL))
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		logger.Error(ctx, "Failed to send stream request to AI service", 
			slog.String("url", reqURL),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查响应状态码
	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		logger.Error(ctx, "AI service returned error status", 
			slog.String("url", reqURL),
			slog.Int("status_code", httpResp.StatusCode),
			slog.String("response", string(body)))
		return nil, fmt.Errorf("AI service returned status code: %d, body: %s", httpResp.StatusCode, string(body))
	}

	// 创建结果channel
	resultChan := make(chan string)

	// 启动goroutine处理流式响应
	go func() {
		defer func() {
			httpResp.Body.Close()
			close(resultChan)
		}()

		// 创建bufio.Reader用于读取响应体
		reader := bufio.NewReader(httpResp.Body)

		for {
			// 读取一行数据
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					logger.Error(ctx, "Failed to read stream response", slog.Any("error", err))
				}
				break
			}

			// 去除换行符
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 解析响应数据
			var resp AIChatResponse
			if err := json.Unmarshal([]byte(line), &resp); err != nil {
				logger.Error(ctx, "Failed to unmarshal stream response", 
					slog.String("line", line),
					slog.Any("error", err))
				continue
			}

			// 发送响应内容到channel
			select {
			case <-ctx.Done():
				logger.Info(ctx, "Stream chat context cancelled")
				return
			case resultChan <- resp.Content:
				// 发送成功
			}
		}

		logger.Info(ctx, "AI stream chat request completed successfully")
	}()

	return resultChan, nil
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
	logger.Info(ctx, "Mock AI chat request started",
		slog.String("conversation_type", convType),
		slog.Int("message_count", len(messages)))

	if c.MockError != nil {
		logger.Error(ctx, "Mock AI chat returned error", 
			slog.String("conversation_type", convType),
			slog.Any("error", c.MockError))
		return "", c.MockError
	}

	// 如果没有设置mock响应，返回默认响应
	if c.MockResponse == "" {
		defaultResp := getDefaultResponse(convType)
		logger.Info(ctx, "Mock AI chat returned default response", 
			slog.String("conversation_type", convType),
			slog.String("response", defaultResp[:50]+"..."))
		return defaultResp, nil
	}

	logger.Info(ctx, "Mock AI chat returned custom response", 
		slog.String("conversation_type", convType),
		slog.String("response", c.MockResponse[:50]+"..."))

	return c.MockResponse, nil
}

// StreamChat 模拟流式输出AI响应
func (c *MockAIClient) StreamChat(ctx context.Context, messages []Message, convType string) (<-chan string, error) {
	logger.Info(ctx, "Mock AI stream chat request started",
		slog.String("conversation_type", convType),
		slog.Int("message_count", len(messages)))

	if c.MockError != nil {
		logger.Error(ctx, "Mock AI stream chat returned error", 
			slog.String("conversation_type", convType),
			slog.Any("error", c.MockError))
		return nil, c.MockError
	}

	// 如果没有设置mock响应，返回默认响应
	response := c.MockResponse
	if response == "" {
		response = getDefaultResponse(convType)
	}

	// 创建结果channel
	resultChan := make(chan string)

	// 启动goroutine模拟流式输出
	go func() {
		defer close(resultChan)

		// 逐字符发送响应，模拟打字机效果
		for _, char := range response {
			select {
			case <-ctx.Done():
				logger.Info(ctx, "Mock AI stream chat context cancelled")
				return
			case resultChan <- string(char):
				// 模拟真实的打字速度，每100毫秒发送一个字符
				time.Sleep(100 * time.Millisecond)
			}
		}

		logger.Info(ctx, "Mock AI stream chat completed successfully", 
			slog.String("conversation_type", convType))
	}()

	return resultChan, nil
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
