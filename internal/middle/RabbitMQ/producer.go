package RabbitMQ

import (
	"context"
	"log"
	"shdp/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var channel *amqp.Channel

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
