package model

import "time"

type Follow struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"uniqueIndex:idx_user_follow" json:"userId"`       // 关注者
	FollowUserID int64     `gorm:"uniqueIndex:idx_user_follow" json:"followUserId"` // 被关注者
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
}
