package routes

import (
	"neuro-guide-go-service/config"
	"neuro-guide-go-service/controllers"

	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router and routes
func InitRouter(cfg *config.Config) *gin.Engine {
	// Initialize controllers
	controllers.InitUserController(cfg)

	r := gin.Default()

	// 为了保持当前实现，我们直接使用路由组来组织API
	api := r.Group("/api")
	{
		// 用户相关路由
		user := api.Group("/user")
		{
			// 微信登录
			user.POST("/wechat-login", controllers.WeChatLogin)

			// 绑定手机号
			user.POST("/bind-phone", func(c *gin.Context) {
				// 实际项目中应该有中间件来处理认证
				// 这里只是一个示例
				controllers.BindPhoneNumber(c)
			})

			// 获取用户资料
			user.GET("/profile/:id", controllers.GetUserProfile)

			// 更新用户资料
			user.PUT("/profile/:id", controllers.UpdateUserProfile)
		}
	}

	return r
}
