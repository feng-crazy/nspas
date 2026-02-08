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
	"regexp"
	"strconv"
	"strings"

	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
)

// AIClient 定义AI服务客户端接口
type AIClient interface {
	// StreamChat 流式调用AI服务进行对话
	StreamChat(ctx context.Context, userID, conversationID string, messages []Message, convType string) (<-chan string, error)
}

// Message 定义消息结构
type Message struct {
	Content string `json:"content"`
	IsUser  bool   `json:"isUser"`
}

// AIChatRequest 定义AI聊天请求结构
type AIChatRequest struct {
	UserID           string    `json:"user_id"`
	ConversationID   string    `json:"conversation_id"`
	Messages         []Message `json:"messages"`
	ConversationType string    `json:"conversation_type"`
}

// HTTPClient 实现AIClient接口，用于实际调用python-ai-service
type HTTPClient struct {
	cfg    *config.Config
	client *http.Client
}

// NewHTTPClient 创建一个新的HTTPClient实例
func NewHTTPClient(cfg *config.Config) *HTTPClient {
	return &HTTPClient{
		cfg:    cfg,
		client: &http.Client{},
	}
}

// StreamChat 流式调用python-ai-service进行对话
// 直接透传python-ai-service返回的流式数据，保留data:前缀
func (c *HTTPClient) StreamChat(ctx context.Context, userID, conversationID string, messages []Message, convType string) (<-chan string, error) {
	logger.Info(ctx, "AI stream chat request started",
		slog.String("user_id", userID),
		slog.String("conversation_id", conversationID),
		slog.String("conversation_type", convType),
		slog.Int("message_count", len(messages)))

	// 创建请求体
	reqBody := AIChatRequest{
		UserID:           userID,
		ConversationID:   conversationID,
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
	httpReq.Header.Set("Accept", "text/event-stream")

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
			close(resultChan)
			// 确保在goroutine结束时关闭响应体
			if httpResp.Body != nil {
				httpResp.Body.Close()
			}
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

			// 直接透传数据，保留原始格式（包括data:前缀）
			select {
			case <-ctx.Done():
				logger.Info(ctx, "Stream chat context cancelled")
				return
			case resultChan <- line:
				// 发送成功
				logData := parseAndFormatLogData(line)
				logger.Debug(ctx, "Stream response sent", slog.String("data", logData))
			}
		}

		logger.Info(ctx, "AI stream chat request completed successfully")
	}()

	return resultChan, nil
}

// parseAndFormatLogData 解码Unicode转义序列并返回可读的日志格式
func parseAndFormatLogData(sseData string) string {
	// 使用 strconv.Unquote 来解码Unicode转义序列
	// 首先尝试添加引号以便Unquote可以正确处理
	processedData := sseData

	// 如果数据已经是SSE格式（以 "data: " 开头），我们需要特别处理
	if strings.HasPrefix(sseData, "data: ") {
		contentPart := strings.TrimPrefix(sseData, "data: ")

		// 尝试解码可能存在的Unicode转义序列
		decodedContent, err := strconv.Unquote(strings.ReplaceAll(`"`+contentPart+`"`, `\"`, `\\"`))
		if err != nil {
			// 如果上面的方法失败，尝试另一种方式：直接替换常见的转义序列
			decodedContent = decodeUnicodeEscapes(contentPart)
		}

		processedData = "data: " + decodedContent
	} else {
		// 对于非SSE格式的数据，尝试解码Unicode转义序列
		decodedContent, err := strconv.Unquote(strings.ReplaceAll(`"`+sseData+`"`, `\"`, `\\"`))
		if err != nil {
			decodedContent = decodeUnicodeEscapes(sseData)
		}
		processedData = decodedContent
	}

	return processedData
}

// decodeUnicodeEscapes 手动解码Unicode转义序列
func decodeUnicodeEscapes(input string) string {
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		codePointHex := strings.TrimPrefix(match, "\\u")
		codePoint, err := strconv.ParseInt(codePointHex, 16, 32)
		if err != nil {
			return match // 如果解析失败，返回原始匹配项
		}
		return string(rune(codePoint))
	})
}
