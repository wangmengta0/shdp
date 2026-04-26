package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"shdp/internal/model"
	"shdp/internal/repository"
	"shdp/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type ShopService struct {
	repo *repository.ShopRepo
	rdb  *redis.Client
}

func NewShopService(repo *repository.ShopRepo, rdb *redis.Client) *ShopService {
	return &ShopService{repo: repo, rdb: rdb}
}
func (s *ShopService) QueryShopById(ctx context.Context, id int) (*model.Shop, error) {
	return s.queryWithMutex(ctx, id)
}
func (s *ShopService) queryWithMutex(ctx context.Context, id int) (*model.Shop, error) {
	shopKey := fmt.Sprintf("%s%d", utils.CacheShopKey, id)
	shopJSON, err := s.rdb.Get(ctx, shopKey).Result()
	if err == nil {
		if shopJSON == "" {
			return nil, errors.New("商户不存在")
		}
		var shop model.Shop
		json.Unmarshal([]byte(shopJSON), &shop)
		return &shop, nil
	} else if err != redis.Nil {
		return nil, errors.New("Redis 服务异常")
	}
	lockKey := fmt.Sprintf("%s%d", utils.LockShopKey, id)
	isLocked := s.tryLock(ctx, lockKey)
	if !isLocked {
		time.Sleep(time.Millisecond * 100)
		return s.queryWithMutex(ctx, id)
	}
	defer s.unlock(ctx, lockKey)
	shopJSON, err = s.rdb.Get(ctx, shopKey).Result()
	if err == nil {
		if shopJSON == "" {
			return nil, errors.New("商户不存在")
		}
		var shop model.Shop
		json.Unmarshal([]byte(shopJSON), &shop)
		return &shop, nil
	}
	shop, err := s.repo.QueryById(id)
	if err != nil || shop == nil {
		s.rdb.Set(ctx, shopKey, "", utils.CacheNullTTL)
		return nil, errors.New("商户不存在")
	}
	shopBytes, _ := json.Marshal(shop)
	s.rdb.Set(ctx, shopKey, string(shopBytes), utils.CacheNullTTL)
	return shop, nil
}
func (s *ShopService) tryLock(ctx context.Context, lockKey string) bool {
	ok, _ := s.rdb.SetNX(ctx, lockKey, "1", utils.LockShopTTL).Result()
	return ok
}

func (s *ShopService) unlock(ctx context.Context, key string) {
	s.rdb.Del(ctx, key)
}
