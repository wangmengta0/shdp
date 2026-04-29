package repository

import (
	"fmt"
	"shdp/internal/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) QueryByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *UserRepo) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
func (r *UserRepo) QueryUsersByIdsOrdered(ids []int64) ([]*model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var users []*model.User
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.FormatInt(id, 10)
	}
	orderStr := fmt.Sprintf("FIELD(id, %s)", strings.Join(idStrs, ","))

	err := r.db.Where("id IN ?", ids).Order(orderStr).Find(&users).Error
	return users, err
}
