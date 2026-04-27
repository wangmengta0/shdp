package service

import (
	"context"
	"errors"
	"fmt"
	"shdp/internal/middle/RabbitMQ"
	"shdp/internal/model"
	"shdp/internal/repository"
	"shdp/pkg/utils"

	"github.com/redis/go-redis/v9"
)

var seckillScript = redis.NewScript(`
-- 1. 获取参数
local stockKey = KEYS[1]
local orderKey = KEYS[2]
local userId = ARGV[1]

-- 2. 判断库存是否充足
-- tonumber 将字符串转为数字，如果 key 不存在默认当做 0
if (tonumber(redis.call('get', stockKey) or "0") <= 0) then
    return 1 -- 错误码 1：库存不足
end

-- 3. 判断一人一单 (该用户ID是否已经在 Set 集合中)
if (redis.call('sismember', orderKey, userId) == 1) then
    return 2 -- 错误码 2：重复下单
end

-- 4. 具备抢购资格，开始内存预扣减
redis.call('incrby', stockKey, -1)
redis.call('sadd', orderKey, userId)

return 0 -- 返回 0 代表抢购成功
`)

type VoucherService struct {
	repo *repository.VoucherRepo
	rdb  *redis.Client
}

func NewVoucherService(repo *repository.VoucherRepo, rdb *redis.Client) *VoucherService {
	return &VoucherService{
		repo: repo,
		rdb:  rdb,
	}
}
func (s *VoucherService) SeckillVoucher(ctx context.Context, voucherId int64, userId int64) (int64, error) {
	stockKey := fmt.Sprintf("seckill:stock:%d", voucherId)
	orderKey := fmt.Sprintf("seckill:order:%d", voucherId)
	result, err := seckillScript.Run(ctx, s.rdb, []string{stockKey, orderKey}, userId).Result()
	if err != nil {
		return 0, errors.New("系统异常，请稍后尝试")
	}
	if result == 1 {
		return 0, errors.New("已抢空")
	}
	if result == 2 {
		return 0, errors.New("您已经抢到过该优惠券了，把机会留给别人吧")
	}
	orderID := utils.GenerateSnowflakeID()
	orderMsg := &model.VoucherOrder{
		ID:        orderID,
		UserID:    userId,
		VoucherID: voucherId,
		PayType:   1,
		Status:    1,
	}
	err = RabbitMQ.PublishSeckillOrder(orderMsg)
	if err != nil {
		return 0, errors.New("系统异常，请稍后尝试")
	}
	return orderID, nil
	//voucher, err := s.repo.QuerySeckillVoucher(voucherId)
	//if err != nil {
	//	return 0, errors.New("优惠券不存在")
	//}
	//now := time.Now()
	//if now.Before(voucher.BeginTime) {
	//	return 0, errors.New("秒杀还未开始")
	//}
	//if now.After(voucher.EndTime) {
	//	return 0, errors.New("秒杀已经结束")
	//}
	//if voucher.Stock < 1 {
	//	return 0, errors.New("已抢空")
	//}
	////TODO一人一单校验
	//lockName := "order:" + strconv.FormatInt(userId, 10)
	//lock := redislock.NewSimpleRedisLock(lockName, s.rdb)
	//isLocked := lock.TryLock(ctx, 5*time.Second)
	//if !isLocked {
	//	return 0, errors.New("操作频繁，请勿重复点击")
	//}
	//defer lock.Unlock(ctx)
	//exists, err := s.repo.CheckUserOrderExists(voucherId, userId)
	//if err != nil {
	//	return 0, errors.New("系统查询订单异常，请稍后尝试")
	//}
	//if exists {
	//	return 0, errors.New("您已经抢到过该优惠券了，把机会留给别人吧")
	//}
	//orderID := utils.GenerateSnowflakeID()
	//order := &model.VoucherOrder{
	//	ID:        orderID,
	//	UserID:    userId,
	//	VoucherID: voucherId,
	//	PayType:   1,
	//	Status:    1,
	//}
	//err = s.repo.SeckillTransaction(voucherId, order)
	//if err != nil {
	//	return 0, err
	//}
	//return orderID, nil
}
