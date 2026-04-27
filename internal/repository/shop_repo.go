package repository

import (
	"shdp/internal/model"

	"gorm.io/gorm"
)

type ShopRepo struct {
	db *gorm.DB
}

func NewShopRepo(db *gorm.DB) *ShopRepo {
	return &ShopRepo{db: db}
}
func (repo *ShopRepo) QueryById(id int) (*model.Shop, error) {
	var shop model.Shop
	err := repo.db.First(&shop, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &shop, nil
}
func (repo *ShopRepo) Update(shop *model.Shop) error {
	return repo.db.Model(&model.Shop{}).Where("id = ?", shop.ID).Updates(shop).Error
}
