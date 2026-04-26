package router

import (
	"net/http"
	"shdp/internal/handler"
	"shdp/internal/middle"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetUpRouter(r *gin.Engine, rdb *redis.Client, userHandler *handler.UserHandler) {
	r.Use(middle.RefreshTokenMiddle(rdb))
	publicGroup := r.Group("/api/user")
	{
		publicGroup.POST("/code", userHandler.SendCode)
		publicGroup.POST("/login", userHandler.Login)
	}
	privateGroup := r.Group("/api/user")
	privateGroup.Use(middle.RequireAuthMiddle())
	{
		privateGroup.GET("/me", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    user,
			})
		})
	}
}
