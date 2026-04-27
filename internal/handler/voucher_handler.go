package handler

import (
	"net/http"
	"shdp/internal/model"
	"shdp/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	voucherService *service.VoucherService
}

func NewVoucherHandler(voucherService *service.VoucherService) *VoucherHandler {
	return &VoucherHandler{voucherService: voucherService}
}
func (h *VoucherHandler) Seckill(c *gin.Context) {
	voucherIDStr := c.Param("id")
	voucherID, err := strconv.ParseInt(voucherIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"seccess":  false,
			"errorMsg": "无效的优惠券ID",
		})
		return
	}
	userRaw, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "errorMsg": "未登录或登录已过期"})
		return
	}
	user, ok := userRaw.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errorMsg": "系统内部状态异常"})
		return
	}
	orderId, err := h.voucherService.SeckillVoucher(c.Request.Context(), voucherID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errorMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "orderId": orderId})
}
