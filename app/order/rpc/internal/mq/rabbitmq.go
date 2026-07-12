package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type RabbitMQ struct {
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
	queue      string
	routingKey string
}

func NewRabbitMQ(addr, exchange, queue, routingKey string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQ{
		conn:       conn,
		channel:    channel,
		exchange:   exchange,
		queue:      queue,
		routingKey: routingKey,
	}, nil
}

func (r *RabbitMQ) PublishDelayMessage(message string, delay time.Duration) error {
	args := amqp091.Table{
		"x-dead-letter-exchange":    r.exchange,
		"x-dead-letter-routing-key": r.routingKey,
	}

	delayQueue := r.queue + "_delay"
	_, err := r.channel.QueueDeclare(
		delayQueue,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare delay queue: %w", err)
	}

	err = r.channel.PublishWithContext(
		context.Background(),
		"",
		delayQueue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Expiration:  fmt.Sprintf("%d", delay.Milliseconds()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Consume(handler func(msg string) error) error {
	msgs, err := r.channel.Consume(
		r.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			logx.Info("Received message: %s", msg.Body)
			if err := handler(string(msg.Body)); err != nil {
				logx.Error("Failed to handle message: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
