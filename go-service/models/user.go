package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	WechatID  string    `json:"wechat_id" bson:"wechat_id"`
	Nickname  string    `json:"nickname" bson:"nickname"`
	Avatar    string    `json:"avatar" bson:"avatar"`
	IsGuest   bool      `json:"is_guest" bson:"is_guest"` // 标识是否是游客用户
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
