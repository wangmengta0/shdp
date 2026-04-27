package handler

// Group 将所有模块的 Handler 聚合在一起
type Group struct {
	User    *UserHandler
	Voucher *VoucherHandler
	// 未来增加 ShopHandler、OrderHandler 直接往这里加即可
}

// NewGroup 统一的构造入口
func NewGroup(user *UserHandler, voucher *VoucherHandler) *Group {
	return &Group{
		User:    user,
		Voucher: voucher,
	}
}
