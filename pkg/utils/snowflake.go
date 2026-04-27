package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	// Epoch 设定一个起始时间戳 (例如 2024-01-01 00:00:00 UTC)
	// 【绝对红线】：这个值一旦投入生产，千万不能修改，否则会导致生成重复 ID！
	Epoch int64 = 1704067200000

	// 各部分占据的位数
	NodeBits uint8 = 10 // 机器码占据 10 位 (最多支持 1024 台服务器)
	StepBits uint8 = 12 // 序列号占据 12 位 (每台机器每毫秒最多生成 4096 个 ID)

	// 最大值 (用于与运算校验边界)
	// -1 ^ (-1 << 10) 算出 1023
	MaxNode int64 = -1 ^ (-1 << NodeBits)
	MaxStep int64 = -1 ^ (-1 << StepBits)

	// 移位偏移量 (组装 64 位数时要把各个部分推到对应的位置)
	TimeShift uint8 = NodeBits + StepBits // 时间戳需要向左推 22 位
	NodeShift uint8 = StepBits            // 机器码需要向左推 12 位
)

// Snowflake 雪花算法核心结构体
type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	node      int64
	step      int64
}

// 全局单例
var (
	node *Snowflake
	once sync.Once
)

// InitSnowflake 初始化全局雪花节点
// nodeID: 当前机器的唯一编号 (0 ~ 1023)
func InitSnowflake(nodeID int64) error {
	if nodeID < 0 || nodeID > MaxNode {
		return errors.New("机器节点ID超出允许范围 (0-1023)")
	}
	once.Do(func() {
		node = &Snowflake{
			timestamp: 0,
			node:      nodeID,
			step:      0,
		}
	})
	return nil
}

// GenerateSnowflakeID 对外暴露的生成全局唯一 ID 方法
func GenerateSnowflakeID() int64 {
	if node == nil {
		// 容错兜底：如果没有显式初始化，默认当做 1 号机
		InitSnowflake(1)
	}
	return node.generate()
}

// generate 内部实际的生成逻辑
func (s *Snowflake) generate() int64 {
	// 加锁保证并发安全
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前时间的毫秒数
	now := time.Now().UnixNano() / 1e6

	if now == s.timestamp {
		// 处于同一毫秒内，序列号自增
		s.step = (s.step + 1) & MaxStep

		if s.step == 0 {
			// 如果同一毫秒内序列号用完了 (超过了 4095)
			// 阻塞等待到下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 不在同一毫秒，序列号重置为 0
		s.step = 0
	}

	// 记录本次生成的时间戳
	s.timestamp = now

	// 【灵魂核心】：通过位运算拼接 64 位整数
	id := ((now - Epoch) << TimeShift) |
		(s.node << NodeShift) |
		(s.step)

	return id
}
