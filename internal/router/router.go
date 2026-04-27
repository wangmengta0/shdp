package router

import (
	"net/http"
	"shdp/internal/handler"
	"shdp/internal/middle"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetUpRouter 形参简化：传入 handlers *handler.Group
func SetUpRouter(r *gin.Engine, rdb *redis.Client, handlers *handler.Group) {
	// 全局 Token 刷新中间件
	r.Use(middle.RefreshTokenMiddle(rdb))

	// 1. 公开接口组 (无需登录)
	publicUserGroup := r.Group("/api/user")
	{
		publicUserGroup.POST("/code", handlers.User.SendCode)
		publicUserGroup.POST("/login", handlers.User.Login)
	}

	// 2. 私有接口组 (必须登录) - 基础路径设为 /api
	privateGroup := r.Group("/api")
	privateGroup.Use(middle.RequireAuthMiddle())
	{
		// /api/user 相关私有接口
		userPrivate := privateGroup.Group("/user")
		{
			userPrivate.GET("/me", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data":    user,
				})
			})
		}

		// /api/voucher 相关私有接口 (逻辑拆分更清晰)
		voucherPrivate := privateGroup.Group("/voucher")
		{
			voucherPrivate.POST("/seckill/:id", handlers.Voucher.Seckill)
		}
	}
}
