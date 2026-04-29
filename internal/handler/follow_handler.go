package handler

import (
	"net/http"
	"shdp/internal/model"
	"shdp/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService *service.FollowService
}

func NewFollowHandler(followService *service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// Follow 关注/取关
func (h *FollowHandler) Follow(c *gin.Context) {
	followIdStr := c.Param("id")
	isFollowStr := c.Param("isFollow") // 通常前端发 PUT /follow/:id/:isFollow

	followUserId, _ := strconv.ParseInt(followIdStr, 10, 64)
	isFollow, _ := strconv.ParseBool(isFollowStr)

	userRaw, _ := c.Get("user")
	user := userRaw.(*model.UserDTO)

	err := h.followService.Follow(c.Request.Context(), user.ID, followUserId, isFollow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// IsFollow 检查是否关注
func (h *FollowHandler) IsFollow(c *gin.Context) {
	followIdStr := c.Param("id")
	followUserId, _ := strconv.ParseInt(followIdStr, 10, 64)

	userRaw, _ := c.Get("user")
	user := userRaw.(*model.UserDTO)

	isFollow, err := h.followService.IsFollow(user.ID, followUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": isFollow})
}

// Common 获取共同关注
func (h *FollowHandler) Common(c *gin.Context) {
	targetIdStr := c.Param("id")
	targetUserId, _ := strconv.ParseInt(targetIdStr, 10, 64)

	userRaw, _ := c.Get("user")
	user := userRaw.(*model.UserDTO)

	users, err := h.followService.CommonFollows(c.Request.Context(), user.ID, targetUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
