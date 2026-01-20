package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	PythonAI PythonAIConfig
	WeChat   WeChatConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret  string
	Expires int // 过期时间（小时）
}

type PythonAIConfig struct {
	BaseURL string
}

type WeChatConfig struct {
	AppID       string
	AppSecret   string
	RedirectURI string
	Scope       string
}

func LoadConfig() *Config {
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		log.Printf("Invalid SERVER_PORT, using default: %v", err)
		port = 8080
	}

	expires, err := strconv.Atoi(getEnv("JWT_EXPIRES", "24"))
	if err != nil {
		log.Printf("Invalid JWT_EXPIRES, using default: %v", err)
		expires = 24
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			URI:      getEnv("DATABASE_URI", "mongodb://localhost:27017"),
			Database: getEnv("DATABASE_NAME", "nspas"),
		},
		JWT: JWTConfig{
			Secret:  getEnv("JWT_SECRET", "your-secret-key"),
			Expires: expires,
		},
		PythonAI: PythonAIConfig{
			BaseURL: getEnv("PYTHON_AI_BASE_URL", "http://localhost:5000"),
		},
		WeChat: WeChatConfig{
			AppID:       getEnv("WECHAT_APP_ID", "your-wechat-app-id"),
			AppSecret:   getEnv("WECHAT_APP_SECRET", "your-wechat-app-secret"),
			RedirectURI: getEnv("WECHAT_REDIRECT_URI", "http://localhost:8080/api/auth/wechat/callback"),
			Scope:       getEnv("WECHAT_SCOPE", "snsapi_userinfo"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
