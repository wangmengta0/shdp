package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shdp/internal/model"
	"shdp/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) SendCode(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "手机号不能为空"})
		return
	}
	_, err := h.userService.SendCode(c.Request.Context(), phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "验证码发送成功"})
}
func (h *UserHandler) Login(c *gin.Context) {
	var loginDTO model.LoginFormDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数格式有误"})
		return
	}
	token, err := h.userService.Login(c.Request.Context(), loginDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": token})
}
