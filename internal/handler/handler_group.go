package handler

// Group 将所有模块的 Handler 聚合在一起
type Group struct {
	User    *UserHandler
	Voucher *VoucherHandler
	Blog    *BlogHandler
	Follow  *FollowHandler // 新增
}

func NewGroup(user *UserHandler, voucher *VoucherHandler, blog *BlogHandler, follow *FollowHandler) *Group {
	return &Group{
		User:    user,
		Voucher: voucher,
		Blog:    blog,
		Follow:  follow,
	}
}
