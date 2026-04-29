package repository

import (
	"shdp/internal/model"

	"gorm.io/gorm"
)

type FollowRepo struct {
	db *gorm.DB
}

func NewFollowRepo(db *gorm.DB) *FollowRepo {
	return &FollowRepo{db: db}
}

func (r *FollowRepo) CheckFollow(userId, followUserId int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).
		Where("user_id = ? AND follow_user_id = ?", userId, followUserId).
		Count(&count).Error
	return count > 0, err
}

// CreateFollow 新增关注
func (r *FollowRepo) CreateFollow(userId, followUserId int64) error {
	follow := &model.Follow{
		UserID:       userId,
		FollowUserID: followUserId,
	}
	return r.db.Create(follow).Error
}

// DeleteFollow 取消关注
func (r *FollowRepo) DeleteFollow(userId, followUserId int64) error {
	return r.db.Where("user_id = ? AND follow_user_id = ?", userId, followUserId).
		Delete(&model.Follow{}).Error
}
