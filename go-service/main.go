package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/controllers"
	"github.com/nspas/go-service/database"
	"github.com/nspas/go-service/logger"
	"github.com/nspas/go-service/middleware"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		logger.Info(context.Background(), "Application shutdown completed")
	}()

	logger.Info(context.Background(), "Application starting",
		slog.String("server_port", fmt.Sprintf("%d", cfg.Server.Port)),
		slog.String("log_level", cfg.Log.Level),
		slog.String("log_path", cfg.Log.Path))

	// 连接数据库
	logger.Info(context.Background(), "Connecting to database")
	err := database.Connect(cfg)
	if err != nil {
		logger.Error(context.Background(), "Failed to connect to database", slog.Any("error", err))
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info(context.Background(), "Database connected successfully")
	defer database.Close()

	// 获取数据库实例
	db := database.Client.Database(cfg.Database.Database)

	// 创建Gin实例
	r := gin.Default()

	// 请求ID中间件
	r.Use(func(c *gin.Context) {
		// 生成请求ID
		requestID := uuid.New().String()
		// 添加到请求头
		c.Header("X-Request-ID", requestID)
		// 添加到上下文
		ctx := logger.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		// 记录请求开始
		logger.Info(ctx, "Request started",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()))
		// 记录请求结束
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)
		logger.Info(ctx, "Request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", duration))
	})

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化控制器
	logger.Info(context.Background(), "Initializing controllers")
	userController := controllers.NewUserController(cfg, db)
	conversationController := controllers.NewConversationController()
	toolController := controllers.NewToolController()
	aiController := controllers.NewAIController(cfg)
	wechatController := controllers.NewWeChatController(cfg, db)

	// 公开路由
	public := r.Group("/api")
	{
		// 用户认证
		public.POST("/auth/register", userController.Register)
		public.POST("/auth/login", userController.Login)

		// 微信登录
		public.GET("/auth/wechat", wechatController.GetWeChatAuthURL)
		public.GET("/auth/wechat/callback", wechatController.WeChatCallback)

		// AI聊天（公开接口，后续可以添加认证）
		public.POST("/ai/chat", aiController.Chat)
	}

	// 受保护路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// 用户相关
		protected.GET("/user", userController.GetCurrentUser)

		// 对话相关
		protected.POST("/conversations", conversationController.CreateConversation)
		protected.GET("/conversations", conversationController.GetUserConversations)
		protected.GET("/conversations/:id", conversationController.GetConversation)
		protected.PUT("/conversations/:id", conversationController.UpdateConversation)
		protected.DELETE("/conversations/:id", conversationController.DeleteConversation)

		// 工具相关
		protected.POST("/tools", toolController.SaveTool)
		protected.GET("/tools", toolController.GetUserTools)
		protected.GET("/tools/:id", toolController.GetTool)
		protected.DELETE("/tools/:id", toolController.DeleteTool)
	}

	// 启动服务器
	port := cfg.Server.Port
	addr := fmt.Sprintf(":%d", port)
	logger.Info(context.Background(), "Server starting", slog.String("address", addr))

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(context.Background(), "Failed to start server", slog.Any("error", err))
		log.Fatalf("Failed to start server: %v", err)
	}
}
