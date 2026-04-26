package model

type LoginFormDTO struct {
	Phone string `json:"phone" binding:"required"` // 可配合 validator 校验手机号格式
	Code  string `json:"code"`
	// Password string `json:"password"` // 如果支持密码登录可扩展
}
