package redislock

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SimpleRedisLock struct {
	name     string
	rdb      *redis.Client
	val      string
	stopChan chan struct{} //用于通知看门狗停止续期的信号通道
}

func NewSimpleRedisLock(name string, rdb *redis.Client) *SimpleRedisLock {
	return &SimpleRedisLock{
		name:     name,
		rdb:      rdb,
		val:      uuid.New().String(),
		stopChan: make(chan struct{}),
	}
}

// 防止误删锁：只有当锁的值与当前持有者的值相同时才删除锁
var unlockScripts = redis.NewScript(`
if redis.call("get",KEYS[1])==ARGV[1] then
	return redis.call("del",KEYS[1])
else 
	return 0
end
`)
var renewScripts = redis.NewScript(`
if redis.call("get",KEYS[1]==ARGV[1] then
	return redis.call("expire",KEYS[1],ARGV[2])
else
	return 0
end
`)

func (l *SimpleRedisLock) TryLock(ctx context.Context, timeout time.Duration) bool {
	key := "lock:" + l.name
	ok, err := l.rdb.SetNX(ctx, key, l.val, timeout).Result()
	if err != nil {
		return false
	}
	go l.startWatchDog(timeout)
	return ok
}
func (l *SimpleRedisLock) Unlock(ctx context.Context) {
	select {
	case l.stopChan <- struct{}{}:
	default:
	}
	key := "lock:" + l.name
	err := unlockScripts.Run(ctx, l.rdb, []string{key}, l.val).Err()
	if err != nil {
		fmt.Printf("释放分布式锁异常: %v\n", err)
	}
}
func (l *SimpleRedisLock) startWatchDog(timeout time.Duration) {
	key := "lock:" + l.name
	ticker := time.NewTicker(timeout / 3)
	defer ticker.Stop()
	for {
		select {
		case <-l.stopChan:
			err := renewScripts.Run(context.Background(), l.rdb, []string{key}, int(timeout.Seconds())).Err()
			if err != nil {
				return
			}
		case <-l.stopChan:
			return
		}

	}
}
