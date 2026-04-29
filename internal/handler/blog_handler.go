package handler

import (
	"net/http"
	"shdp/internal/model"
	"shdp/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService *service.BlogService
}

func NewBlogHandler(blogService *service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}
func (h *BlogHandler) LikeBlog(c *gin.Context) {
	idStr := c.Param("id")
	blogId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的博客ID",
		})
	}
	userRaw, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录或登录已过期",
		})
	}
	user, ok := userRaw.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "系统内部状态异常",
		})
	}
	err = h.blogService.LikeBlog(c.Request.Context(), blogId, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
func (h *BlogHandler) QueryBlogLikes(c *gin.Context) {
	idStr := c.Param("id")
	blogId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的博客ID"})
		return
	}

	users, err := h.blogService.QueryBlogLikes(c.Request.Context(), blogId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
