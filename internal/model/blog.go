package model

import "time"

type Blog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID     int64     `json:"shopId"`
	UserID     int64     `json:"userId"`
	Title      string    `json:"title"`
	Images     string    `json:"images"`
	Content    string    `json:"content"`
	Liked      int       `json:"liked"` // 点赞数量
	Comments   int       `json:"comments"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`

	// 非数据库字段，用于返回给前端
	IsLike bool `gorm:"-" json:"isLike"`
}
