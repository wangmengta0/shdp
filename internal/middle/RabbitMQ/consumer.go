package RabbitMQ

import (
	"context"
	"log"
	"shdp/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartCacheDeleteConsumer(rdb *redis.Client) {
	msg, err := channel.Consume(
		config.Conf.RabbitMQ.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("消费队列失败: %v", err)
	}
	go func() {
		for d := range msg {
			key := string(d.Body)
			log.Printf("收到补偿任务：尝试删除 Key [%s]", key)
			err := rdb.Del(context.Background(), key).Err()
			if err != nil {
				log.Printf("补偿删除失败 [%s]: %v，等待重试...", key, err)
				time.Sleep(5 * time.Second)
				d.Nack(false, true)
			} else {
				log.Printf("补偿删除成功 [%s]", key)
				d.Ack(false)
			}
		}
	}()
}
