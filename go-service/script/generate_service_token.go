package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// 🔑 必须和你 AuthMiddleware 中 cfg.JWT.Secret 一致！
	secret := "your-very-strong-secret-key-here"

	// 🧾 构造永久有效的服务 Token（无 exp 字段 = 永不过期）
	claims := jwt.MapClaims{
		"service_id": "python-ai-service", // 标识调用方
		"role":       "internal",
		"scope":      "all", // 可选：定义权限范围
		// 注意：这里故意不设置 "exp"，使其永久有效
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	fmt.Println("永久有效的服务 Token:")
	fmt.Println(tokenString)
}
