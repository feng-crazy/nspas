package config

type Config struct {
	Port               string // 服务运行端口
	AppName            string // 应用名称
	Environment        string // 运行环境 (development/production)
	MongoDBURI         string // MongoDB连接URI
	MongoDBName        string // MongoDB数据库名称
	WeChatAppID        string // 微信AppID
	WeChatSecret       string // 微信AppSecret
	PythonAIServiceURL string // Python AI服务URL
}
