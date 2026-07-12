// Package mq 提供基于 RabbitMQ 的轻量发布/消费封装。
// 采用惰性连接：连接失败不影响服务启动，发布/消费时自动重连。
package mq

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher 发布者：惰性连接，线程安全。支持立即发送和延迟发送（死信交换机）。
type Publisher struct {
	url           string
	queue         string
	delayQueue    string
	delayExchange string
	mu            sync.Mutex
	conn          *amqp.Connection
	ch            *amqp.Channel
}

// NewPublisher 创建发布者（此时不建立连接）。
func NewPublisher(url, queue string) *Publisher {
	return &Publisher{url: url, queue: queue}
}

// NewPublisherWithDelay 创建支持延迟发送的发布者。
func NewPublisherWithDelay(url, queue, delayQueue, delayExchange string) *Publisher {
	return &Publisher{
		url:           url,
		queue:         queue,
		delayQueue:    delayQueue,
		delayExchange: delayExchange,
	}
}

func (p *Publisher) ensure() error {
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil {
		return nil
	}
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if _, err = ch.QueueDeclare(p.queue, true, false, false, false, nil); err != nil {
		_ = conn.Close()
		return err
	}
	if p.delayQueue != "" && p.delayExchange != "" {
		if err := ch.ExchangeDeclare(p.delayExchange, "direct", true, false, false, false, nil); err != nil {
			_ = conn.Close()
			return err
		}
		if _, err := ch.QueueDeclare(p.delayQueue, true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": p.queue,
		}); err != nil {
			_ = conn.Close()
			return err
		}
		if err := ch.QueueBind(p.delayQueue, p.delayQueue, p.delayExchange, false, nil); err != nil {
			_ = conn.Close()
			return err
		}
	}
	p.conn, p.ch = conn, ch
	return nil
}

// Publish 发布一条持久化消息到队列。
func (p *Publisher) Publish(ctx context.Context, body []byte) error {
	return p.PublishWithDelay(ctx, body, 0)
}

// PublishWithDelay 发布消息，支持延迟发送。delaySec > 0 时使用延迟队列（死信交换机）。
func (p *Publisher) PublishWithDelay(ctx context.Context, body []byte, delaySec int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ensure(); err != nil {
		return err
	}
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	}
	if delaySec > 0 && p.delayQueue != "" && p.delayExchange != "" {
		publishing.Expiration = strconv.FormatInt(delaySec*1000, 10)
		return p.ch.PublishWithContext(ctx, p.delayExchange, p.delayQueue, false, false, publishing)
	}
	return p.ch.PublishWithContext(ctx, "", p.queue, false, false, publishing)
}

// Close 关闭连接。
func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}

// Consumer 消费者：断线自动重连重试。
type Consumer struct {
	url     string
	queue   string
	handler func([]byte) error
}

// NewConsumer 创建消费者。
func NewConsumer(url, queue string, handler func([]byte) error) *Consumer {
	return &Consumer{url: url, queue: queue, handler: handler}
}

// Start 在后台协程持续消费；连接失败会每 5s 重试。
func (c *Consumer) Start() {
	go func() {
		for {
			if err := c.consume(); err != nil {
				logx.Errorf("[mq] consume error: %v, retry in 5s", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func (c *Consumer) consume() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err = ch.QueueDeclare(c.queue, true, false, false, false, nil); err != nil {
		return err
	}

	msgs, err := ch.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	logx.Infof("[mq] consuming queue=%s", c.queue)

	for d := range msgs {
		if err := c.handler(d.Body); err != nil {
			logx.Errorf("[mq] handle message error: %v", err)
			_ = d.Nack(false, true) // 重回队列
			continue
		}
		_ = d.Ack(false)
	}
	return amqp.ErrClosed
}
