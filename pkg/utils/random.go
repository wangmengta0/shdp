package utils

import (
	"math/rand"
	"time"
)

func RandomDuration(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + time.Duration(r.Int63n(int64(max-min)))
}
