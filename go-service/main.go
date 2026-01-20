package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/controllers"
	"github.com/nspas/go-service/database"
	"github.com/nspas/go-service/middleware"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 连接数据库
	err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 获取数据库实例
	db := database.Client.Database(cfg.Database.Database)

	// 创建Gin实例
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化控制器
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
	log.Printf("Server running on port %d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
