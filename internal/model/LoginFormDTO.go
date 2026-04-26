package model

// UserDTO 存入 Redis 和返回给前端的脱敏用户信息
type UserDTO struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}
