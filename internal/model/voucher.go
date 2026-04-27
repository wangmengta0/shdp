package model

import "time"

// Voucher 优惠券基本信息表
type Voucher struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID      int64     `gorm:"index" json:"shopId"`                  // 归属的商户 ID
	Title       string    `gorm:"type:varchar(255)" json:"title"`       // 代金券标题
	SubTitle    string    `gorm:"type:varchar(255)" json:"subTitle"`    // 副标题
	Rules       string    `gorm:"type:varchar(1024)" json:"rules"`      // 使用规则
	PayValue    int64     `gorm:"type:int unsigned" json:"payValue"`    // 支付金额，单位：分 (避免浮点数精度丢失)
	ActualValue int64     `gorm:"type:int unsigned" json:"actualValue"` // 实际抵扣金额，单位：分
	Type        int       `gorm:"type:tinyint unsigned" json:"type"`    // 券类型: 0:普通券, 1:秒杀券
	Status      int       `gorm:"type:tinyint unsigned" json:"status"`  // 状态: 1:上架, 2:下架, 3:过期
	CreateTime  time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime  time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// SeckillVoucher 秒杀优惠券库存与时间表 (与 Voucher 一对一)
type SeckillVoucher struct {
	VoucherID  int64     `gorm:"primaryKey" json:"voucherId"` // 关联 Voucher.ID，不再使用自增ID
	Stock      int       `gorm:"type:int" json:"stock"`       // 剩余库存 (秒杀扣减核心字段)
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	BeginTime  time.Time `json:"beginTime"` // 秒杀生效时间
	EndTime    time.Time `json:"endTime"`   // 秒杀失效时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

// VoucherOrder 优惠券秒杀订单表
type VoucherOrder struct {
	ID         int64      `gorm:"primaryKey" json:"id"`                 // 订单ID (使用雪花算法生成，保证全局唯一)
	UserID     int64      `gorm:"index" json:"userId"`                  // 下单用户ID
	VoucherID  int64      `gorm:"index" json:"voucherId"`               // 购买的代金券ID
	PayType    int        `gorm:"type:tinyint unsigned" json:"payType"` // 支付方式: 1:微信, 2:支付宝
	Status     int        `gorm:"type:tinyint unsigned" json:"status"`  // 订单状态: 1:未支付, 2:已支付, 3:已核销, 4:已取消, 5:退款中, 6:已退款
	CreateTime time.Time  `gorm:"autoCreateTime" json:"createTime"`     // 下单时间
	PayTime    *time.Time `json:"payTime"`                              // 支付时间 (使用指针允许为 nil)
	UseTime    *time.Time `json:"useTime"`                              // 核销时间 (使用指针允许为 nil)
	RefundTime *time.Time `json:"refundTime"`                           // 退款时间 (使用指针允许为 nil)
	UpdateTime time.Time  `gorm:"autoUpdateTime" json:"updateTime"`
}
