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
func (s *ShopService) SaveShop2Redis(ctx context.Context, id int, expireSeconds int) error {
	shop, err := s.repo.QueryById(id)
	if err != nil || shop == nil {
		return errors.New("商户不存在")
	}
	redisDate := model.RedisData{
		Data:       shop,
		ExpireTime: time.Now().Add(time.Duration(expireSeconds) * time.Second),
	}
	dataBytes, _ := json.Marshal(redisDate)
	return s.rdb.Set(ctx, fmt.Sprintf("%s%d", utils.CacheShopKey, id), string(dataBytes), 0).Err()
}
func (s *ShopService) QueryWithLogicalExpire(ctx context.Context, id int) (*model.Shop, error) {
	shopKey := fmt.Sprintf("%s%d", utils.CacheShopKey, id)
	jsonStr, err := s.rdb.Get(ctx, shopKey).Result()
	if err == redis.Nil || jsonStr == "" {
		return nil, nil
	} else if err != nil {
		return nil, errors.New("Redis 服务异常")
	}
	var redisData model.RedisData
	json.Unmarshal([]byte(jsonStr), &redisData)
	shopBytes, _ := json.Marshal(redisData.Data)
	var shop model.Shop
	json.Unmarshal(shopBytes, &shop)
	if time.Now().Before(redisData.ExpireTime) {
		return &shop, nil
	}
	lockKey := fmt.Sprintf("%s%d", utils.CacheShopKey, id)
	isLocked := s.tryLock(ctx, lockKey)
	if isLocked {
		ckeckStr, _ := s.rdb.Get(ctx, shopKey).Result()
		var checkData model.RedisData
		json.Unmarshal([]byte(ckeckStr), &checkData)
		if time.Now().Before(checkData.ExpireTime) {
			s.unlock(ctx, lockKey)
			json.Unmarshal([]byte(jsonStr), &redisData)
			shopBytes, _ := json.Marshal(redisData.Data)
			json.Unmarshal(shopBytes, &shop)
			return &shop, nil
		}
		go func() {
			defer s.unlock(context.Background(), lockKey)
			s.SaveShop2Redis(context.Background(), id, 30*60)
		}()
	}
	return &shop, nil
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
	jitter := utils.RandomDuration(1*time.Minute, 5*time.Minute)
	finalTTL := utils.CacheNullTTL + jitter
	shop, err := s.repo.QueryById(id)
	if err != nil || shop == nil {
		s.rdb.Set(ctx, shopKey, "", finalTTL)
		return nil, errors.New("商户不存在")
	}
	shopBytes, _ := json.Marshal(shop)
	finalTTL = utils.CacheShopTTL + jitter
	s.rdb.Set(ctx, shopKey, string(shopBytes), finalTTL)
	return shop, nil
}
func (s *ShopService) tryLock(ctx context.Context, lockKey string) bool {
	ok, _ := s.rdb.SetNX(ctx, lockKey, "1", utils.LockShopTTL).Result()
	return ok
}

func (s *ShopService) unlock(ctx context.Context, key string) {
	s.rdb.Del(ctx, key)
}
