package utils

import "time"

const (
	LoginCodeKey = "login:code:"
	LoginCodeTTL = 2 * time.Minute
	LoginUserKey = "login:token:"
	LoginUserTTL = 30 * time.Minute
	CacheShopKey = "cache:shop:"
	CacheShopTTL = 30 * time.Minute
	CacheNullTTL = 2 * time.Minute // 缓存穿透：空对象的过期时间(设短一点)
	LockShopKey  = "lock:shop:"
	LockShopTTL  = 10 * time.Second // 分布式锁的兜底过期时间
	BlogLikedKey = "blog:liked:"
	FollowsKey   = "follows:"
)
