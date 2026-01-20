package services

import "errors"

var (
	// 通用错误
	ErrInternalServer = errors.New("internal server error")
	ErrNotFound       = errors.New("resource not found")

	// 用户相关错误
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// 对话相关错误
	ErrConversationNotFound = errors.New("conversation not found")

	// 工具相关错误
	ErrToolNotFound = errors.New("tool not found")
)
