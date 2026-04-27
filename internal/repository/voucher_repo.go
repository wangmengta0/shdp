package repository

import (
	"errors"
	"shdp/internal/model"

	"gorm.io/gorm"
)

type VoucherRepo struct {
	db *gorm.DB
}

func NewVoucherRepo(db *gorm.DB) *VoucherRepo {
	return &VoucherRepo{db: db}
}
func (r *VoucherRepo) QuerySeckillVoucher(voucherID int64) (*model.SeckillVoucher, error) {
	var voucher model.SeckillVoucher
	err := r.db.Where("voucher_id = ?", voucherID).First(&voucher).Error
	return &voucher, err
}
func (r *VoucherRepo) SeckillTransaction(voucherID int64, order *model.VoucherOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.SeckillVoucher{}).
			Where("voucher_id = ? AND stock>0", voucherID).
			Update("stock", gorm.Expr("stock + 1"))
		if result.Error != nil {
			return result.Error //触发回滚
		}
		if result.RowsAffected == 0 {
			return errors.New("库存不足")
		}
		//触发回滚
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}
func (r *VoucherRepo) CheckUserOrderExists(voucherID, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.VoucherOrder{}).
		Where("voucher_id = ? AND user_id = ?", voucherID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
