package main

import (
	"log"
	"os"

	"neuro-guide-go-service/config"
	"neuro-guide-go-service/database"
	"neuro-guide-go-service/routes"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// 初始化配置
	cfg := &config.Config{
		Port:               os.Getenv("PORT"),
		AppName:            "neuro-guide-go-service",
		Environment:        os.Getenv("ENVIRONMENT"),
		MongoDBURI:         os.Getenv("MONGODB_URI"),
		MongoDBName:        os.Getenv("MONGODB_NAME"),
		WeChatAppID:        os.Getenv("WECHAT_APP_ID"),
		WeChatSecret:       os.Getenv("WECHAT_APP_SECRET"),
		PythonAIServiceURL: os.Getenv("PYTHON_AI_SERVICE_URL"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080" // 默认端口
	}

	if cfg.MongoDBURI == "" {
		cfg.MongoDBURI = "mongodb://localhost:27017" // 默认MongoDB地址
	}

	if cfg.MongoDBName == "" {
		cfg.MongoDBName = "neuro_guide" // 默认数据库名
	}

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化路由
	router := routes.InitRouter(cfg)

	// 启动服务器
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
