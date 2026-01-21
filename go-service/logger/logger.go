package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/nspas/go-service/config"
)

// Logger 全局日志实例
var (
	globalLogger *slog.Logger
	once         sync.Once
)

// RequestIDKey 上下文请求ID键
const RequestIDKey = "request_id"

// InitLogger 初始化日志
func InitLogger(cfg *config.LogConfig) error {
	var err error
	once.Do(func() {
		// 创建日志目录
		if err := os.MkdirAll(cfg.Path, 0755); err != nil {
			err = fmt.Errorf("failed to create log directory: %w", err)
			return
		}

		// 解析日志等级
		level := parseLogLevel(cfg.Level)

		// 创建日志处理器
		handlers := make([]slog.Handler, 0, 2)

		// 控制台处理器（文本格式）
		consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{
						Key:   slog.TimeKey,
						Value: slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000")),
					}
				}
				return a
			},
		})
		handlers = append(handlers, consoleHandler)

		// 文件处理器（JSON格式）
		fileHandler, fileErr := newFileHandler(cfg)
		if fileErr != nil {
			err = fmt.Errorf("failed to create file handler: %w", fileErr)
			return
		}
		handlers = append(handlers, fileHandler)

		// 创建多处理器
		var multiHandler slog.Handler
		if len(handlers) == 1 {
			multiHandler = handlers[0]
		} else {
			multiHandler = newMultiHandler(handlers...)
		}

		// 创建日志记录器
		globalLogger = slog.New(multiHandler)

		// 替换默认日志记录器
		slog.SetDefault(globalLogger)
	})

	return err
}

// parseLogLevel 解析日志等级
func parseLogLevel(levelStr string) slog.Level {
	levelStr = strings.ToLower(levelStr)
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// newFileHandler 创建文件日志处理器
func newFileHandler(cfg *config.LogConfig) (slog.Handler, error) {
	// 获取当前时间作为日志文件名的一部分
	now := time.Now()
	logFileName := fmt.Sprintf("app-%s.log", now.Format("2006-01-02"))
	logFilePath := filepath.Join(cfg.Path, logFileName)

	// 打开或创建日志文件
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// 创建JSON格式处理器
	return slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     parseLogLevel(cfg.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 自定义时间格式
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z07:00")),
				}
			}
			return a
		},
	}), nil
}

// newMultiHandler 创建多处理器
func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

// multiHandler 多日志处理器

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multiple handler errors: %v", errs)
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Debug 调试日志
func Debug(ctx context.Context, msg string, args ...any) {
	logWithSource(ctx, slog.LevelDebug, msg, args...)
}

// Info 信息日志
func Info(ctx context.Context, msg string, args ...any) {
	logWithSource(ctx, slog.LevelInfo, msg, args...)
}

// Warn 警告日志
func Warn(ctx context.Context, msg string, args ...any) {
	logWithSource(ctx, slog.LevelWarn, msg, args...)
}

// Error 错误日志
func Error(ctx context.Context, msg string, args ...any) {
	logWithSource(ctx, slog.LevelError, msg, args...)
}

// logWithSource 带调用位置的日志记录
func logWithSource(ctx context.Context, level slog.Level, msg string, args ...any) {
	// 检查日志是否启用
	if globalLogger == nil || !globalLogger.Enabled(ctx, level) {
		return
	}

	// 获取调用位置
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // 跳过当前函数和log函数
	frame, _ := runtime.CallersFrames(pcs[:]).Next()

	// 提取文件名（仅保留最后两级目录）
	parts := strings.Split(frame.File, string(filepath.Separator))
	if len(parts) > 2 {
		frame.File = strings.Join(parts[len(parts)-2:], string(filepath.Separator))
	}

	// 添加请求ID（如果有）
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		args = append(args, slog.String(RequestIDKey, requestID))
	}

	// 添加调用位置
	args = append(args, slog.String("source", fmt.Sprintf("%s:%d", frame.File, frame.Line)))

	// 创建日志记录
	record := slog.NewRecord(time.Now(), level, msg, frame.PC)
	record.Add(args...)

	// 记录日志
	if err := globalLogger.Handler().Handle(ctx, record); err != nil {
		fmt.Fprintf(os.Stderr, "failed to log: %v\n", err)
	}
}

// GetLogger 获取日志实例
func GetLogger() *slog.Logger {
	return globalLogger
}

// WithRequestID 创建带请求ID的上下文
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID 从上下文中获取请求ID
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
