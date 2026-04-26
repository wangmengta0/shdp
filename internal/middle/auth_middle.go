package middle

import (
	"net/http"
	"shdp/internal/model"
	"shdp/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RefreshTokenMiddle(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.Next()
			return
		}
		key := utils.LoginUserKey + token
		userMap, err := rdb.HGetAll(c.Request.Context(), key).Result()
		if err != nil || len(userMap) == 0 {
			c.Next()
			return
		}
		id, _ := strconv.ParseInt(userMap["id"], 10, 64)
		userDTO := model.UserDTO{
			ID:       id,
			NickName: userMap["nickName"],
			Icon:     userMap["icon"],
		}
		rdb.Expire(c.Request.Context(), key, utils.LoginUserTTL)
		c.Set("user", &userDTO)
		c.Next()
	}
}
func RequireAuthMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get("user")
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":  false,
				"errorMsg": "未登录或登录已过期",
			})
			return
		}
		c.Next()
	}
}
