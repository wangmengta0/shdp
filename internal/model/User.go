package model

import "time"

type User struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone      string    `gorm:"type:varchar(11);uniqueIndex" json:"phone"`
	Password   string    `gorm:"type:varchar(128)" json:"-"` // 隐藏密码
	NickName   string    `gorm:"type:varchar(32)" json:"nickName"`
	Icon       string    `gorm:"type:varchar(255)" json:"icon"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
