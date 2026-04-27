package RabbitMQ

import (
	"encoding/json"
	"log"
	"shdp/internal/model"
	"shdp/internal/repository"
)

func StartSeckillOrderConsumer(repo *repository.VoucherRepo) {
	msgs, err := channel.Consume(
		SeckillQueueName,
		"seckill_worker",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("开启秒杀队列消费失败: %v", err)
	}
	log.Println("🚀 MQ 秒杀异步落盘消费者已启动，等待处理订单...")
	go func() {
		for d := range msgs {
			var order model.VoucherOrder
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("解析秒杀订单消息失败: %v", err)
				d.Nack(false, false) // 格式错误，直接丢弃
				continue
			}
			err = repo.SeckillTransaction(order.VoucherID, &order)
			if err != nil {
				log.Printf("订单异步落盘失败 [OrderID: %d], 原因: %v。正在重试...", order.ID, err)
				// 真正的生产环境这里需要控制重试次数，防止毒消息堵塞队列
				d.Nack(false, true)
			} else {
				log.Printf("订单异步落盘成功 [OrderID: %d, UserID: %d]", order.ID, order.UserID)
				d.Ack(false) // 确认消费
			}
		}
	}()
}
