package RabbitMQ

import (
	"context"
	"encoding/json"
	"log"
	"shdp/config"
	"shdp/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

var channel *amqp.Channel

const SeckillQueueName = "seckill.order.queue"

func InitRabbitMQ() {
	conn, err := amqp.Dial(config.Conf.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("无法连接RabbitMQ:%v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法开启Channel:%v", err)
	}
	_, err = ch.QueueDeclare(
		config.Conf.RabbitMQ.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("无法声明队列:%v", err)
	}
	_, err = ch.QueueDeclare(
		SeckillQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("无法声明秒杀队列: %v", err)
	}
	channel = ch
}
func PublishDeleteKey(key string) error {
	return channel.PublishWithContext(
		context.Background(),
		"",
		config.Conf.RabbitMQ.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(key),
			DeliveryMode: amqp.Persistent,
		},
	)
}
func PublishSeckillOrder(order *model.VoucherOrder) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return channel.PublishWithContext(
		context.Background(),
		"",
		SeckillQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
