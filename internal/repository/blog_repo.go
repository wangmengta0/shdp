package repository

import (
	"shdp/internal/model"

	"gorm.io/gorm"
)

type BlogRepo struct {
	db *gorm.DB
}

func NewBlogRepo(db *gorm.DB) *BlogRepo {
	return &BlogRepo{db: db}
}
func (r *BlogRepo) UpdateBlogLiked(blogId int64, step int) error {
	return r.db.Model(&model.Blog{}).
		Where("id = ?", blogId).
		Update("liked", gorm.Expr("liked+?", step)).Error
}
