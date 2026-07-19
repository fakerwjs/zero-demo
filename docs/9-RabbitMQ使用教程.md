# Zero-Demo RabbitMQ 使用教程

## 1. 概述

本教程详细介绍 RabbitMQ 在 Zero-Demo 项目中的配置、使用方法和最佳实践。

**RabbitMQ 版本**: 3.13.0

**管理界面地址**: http://localhost:15672（用户名: admin，密码: admin123456）

## 2. RabbitMQ 架构

### 2.1 架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        RabbitMQ Server                          │
│                                                                 │
│  ┌──────────┐    ┌──────────────┐    ┌──────────────────────┐  │
│  │  Exchanges│───▶│   Queues     │───▶│    Consumers         │  │
│  │  (交换机) │    │   (队列)     │    │   (消费者)           │  │
│  └────┬─────┘    └──────┬───────┘    └──────────┬───────────┘  │
│       │                 │                        │              │
│       │ 绑定(Binding)   │                        │              │
│       ▼                 ▼                        ▼              │
│  ┌──────────┐    ┌──────────────┐    ┌──────────────────────┐  │
│  │  Producers│    │   Messages   │    │    Acknowledgment   │  │
│  │  (生产者) │───▶│   (消息)     │◀───│    (确认机制)        │  │
│  └──────────┘    └──────────────┘    └──────────────────────┘  │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │                    队列类型                                │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐ │ │
│  │  │ Direct Queue│ │ Topic Queue │ │ Delayed Queue(延迟) │ │ │
│  │  └─────────────┘ └─────────────┘ └─────────────────────┘ │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 项目中的队列设计

| 队列名称 | 类型 | 用途 | 绑定键 |
|----------|------|------|--------|
| notification.send | Direct | 发送通知 | notification.send |
| order.cancel | Delayed | 订单超时取消 | order.cancel |
| payment.result | Direct | 支付结果处理 | payment.result |
| stock.update | Direct | 库存更新 | stock.update |

## 3. 环境准备

### 3.1 启动 RabbitMQ

```bash
cd deploy/docker
docker compose up -d rabbitmq
```

### 3.2 验证服务状态

```bash
# 查看 RabbitMQ 容器状态
docker compose ps rabbitmq

# 检查 RabbitMQ 服务状态
curl http://localhost:15672/api/health/checks/node -u admin:admin123456

# 查看队列列表
curl http://localhost:15672/api/queues -u admin:admin123456 | jq '.[].name'
```

### 3.3 创建延迟队列插件

RabbitMQ 默认不支持延迟队列，需要安装插件：

```bash
# 进入容器安装延迟插件
docker exec -it zero-rabbitmq rabbitmq-plugins enable rabbitmq_delayed_message_exchange

# 重启容器生效
docker restart zero-rabbitmq
```

## 4. 代码集成

### 4.1 引入依赖

```go
// go.mod
require (
    github.com/rabbitmq/amqp091-go v1.10.0
)
```

### 4.2 配置 RabbitMQ

在服务配置文件中添加 RabbitMQ 配置：

```yaml
# app/order/rpc/etc/order-rpc.yaml
RabbitMQ:
  Host: ${RABBITMQ_HOST:127.0.0.1}
  Port: ${RABBITMQ_PORT:5672}
  User: ${RABBITMQ_USER:admin}
  Password: ${RABBITMQ_PASSWORD:admin123456}
  VHost: /
```

### 4.3 创建连接工厂

```go
// pkg/rabbitmq/rabbitmq.go
package rabbitmq

import (
    "fmt"
    "log"
    "time"

    "github.com/rabbitmq/amqp091-go"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    VHost    string
}

func NewRabbitMQ(cfg Config) (*amqp091.Connection, error) {
    connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost)

    var conn *amqp091.Connection
    var err error

    for i := 0; i < 5; i++ {
        conn, err = amqp091.Dial(connStr)
        if err == nil {
            return conn, nil
        }
        log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds... (%d/5)", i+1)
        time.Sleep(5 * time.Second)
    }

    return nil, err
}
```

## 5. 生产者实现

### 5.1 普通消息发送

```go
// app/notification/rpc/internal/logic/send_notification_logic.go
func (l *SendNotificationLogic) SendNotification(in *notification.SendNotificationReq) (*notification.SendNotificationResp, error) {
    conn, err := rabbitmq.NewRabbitMQ(l.svcCtx.Config.RabbitMQ)
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    defer ch.Close()

    err = ch.PublishWithContext(l.ctx,
        "notification_exchange",
        "notification.send",
        false,
        false,
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        []byte(in.Message),
        })
    if err != nil {
        return nil, err
    }

    return &notification.SendNotificationResp{Success: true}, nil
}
```

### 5.2 延迟消息发送

```go
// app/order/rpc/internal/logic/create_order_logic.go
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
    orderNo := util.GenerateOrderNo()
    
    conn, err := rabbitmq.NewRabbitMQ(l.svcCtx.Config.RabbitMQ)
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    defer ch.Close()

    err = ch.PublishWithContext(l.ctx,
        "delayed_exchange",
        "order.cancel",
        false,
        false,
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        []byte(orderNo),
            Headers: amqp091.Table{
                "x-delay": 300000,
            },
        })
    if err != nil {
        return nil, err
    }

    return &order.CreateOrderResp{OrderNo: orderNo}, nil
}
```

## 6. 消费者实现

### 6.1 普通消费者

```go
// app/notification/rpc/internal/consumer/notification_consumer.go
package consumer

import (
    "encoding/json"
    "log"

    "github.com/rabbitmq/amqp091-go"
)

type NotificationConsumer struct {
    ch       *amqp091.Channel
    queue    string
    handler  func(message []byte) error
}

func NewNotificationConsumer(ch *amqp091.Channel, queue string, handler func(message []byte) error) *NotificationConsumer {
    return &NotificationConsumer{
        ch:      ch,
        queue:   queue,
        handler: handler,
    }
}

func (c *NotificationConsumer) Start() error {
    msgs, err := c.ch.Consume(
        c.queue,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    go func() {
        for d := range msgs {
            err := c.handler(d.Body)
            if err != nil {
                log.Printf("Failed to handle message: %v", err)
                d.Nack(false, true)
            } else {
                d.Ack(false)
            }
        }
    }()

    return nil
}
```

### 6.2 启动消费者

```go
// app/notification/rpc/internal/server/notificationrpcserver.go
func (s *NotificationRpcServer) StartConsumers() error {
    conn, err := rabbitmq.NewRabbitMQ(s.config.RabbitMQ)
    if err != nil {
        return err
    }

    ch, err := conn.Channel()
    if err != nil {
        return err
    }

    err = ch.ExchangeDeclare(
        "notification_exchange",
        "direct",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    q, err := ch.QueueDeclare(
        "notification.send",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    err = ch.QueueBind(
        q.Name,
        "notification.send",
        "notification_exchange",
        false,
        nil,
    )
    if err != nil {
        return err
    }

    consumer := NewNotificationConsumer(ch, q.Name, s.handleNotification)
    return consumer.Start()
}

func (s *NotificationRpcServer) handleNotification(message []byte) error {
    var req notification.SendNotificationReq
    err := json.Unmarshal(message, &req)
    if err != nil {
        return err
    }

    _, err = s.SendNotification(nil, &req)
    return err
}
```

## 7. 延迟队列配置

### 7.1 创建延迟交换机

```go
// pkg/rabbitmq/delayed.go
package rabbitmq

import "github.com/rabbitmq/amqp091-go"

func DeclareDelayedExchange(ch *amqp091.Channel, exchangeName string) error {
    return ch.ExchangeDeclare(
        exchangeName,
        "x-delayed-message",
        true,
        false,
        false,
        false,
        amqp091.Table{
            "x-delayed-type": "direct",
        },
    )
}
```

### 7.2 订单超时取消消费者

```go
// app/order/rpc/internal/consumer/order_cancel_consumer.go
package consumer

import (
    "log"

    "github.com/rabbitmq/amqp091-go"
)

type OrderCancelConsumer struct {
    ch      *amqp091.Channel
    queue   string
    handler func(orderNo string) error
}

func NewOrderCancelConsumer(ch *amqp091.Channel, queue string, handler func(orderNo string) error) *OrderCancelConsumer {
    return &OrderCancelConsumer{
        ch:      ch,
        queue:   queue,
        handler: handler,
    }
}

func (c *OrderCancelConsumer) Start() error {
    msgs, err := c.ch.Consume(
        c.queue,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    go func() {
        for d := range msgs {
            orderNo := string(d.Body)
            err := c.handler(orderNo)
            if err != nil {
                log.Printf("Failed to cancel order %s: %v", orderNo, err)
                d.Nack(false, true)
            } else {
                d.Ack(false)
            }
        }
    }()

    return nil
}
```

## 8. 管理界面使用

### 8.1 登录管理界面

访问 http://localhost:15672，使用以下凭据登录：
- 用户名: admin
- 密码: admin123456

### 8.2 查看队列状态

1. 点击左侧菜单 **Queues**
2. 可以查看：
   - 队列名称
   - 消息数量
   - 消费者数量
   - 入队/出队速率

### 8.3 查看消息

1. 在队列列表中找到目标队列
2. 点击 **Get Messages**
3. 可以查看消息内容

### 8.4 管理交换机

1. 点击左侧菜单 **Exchanges**
2. 可以：
   - 创建新交换机
   - 绑定队列到交换机
   - 发布测试消息

### 8.5 监控指标

1. 点击左侧菜单 **Overview**
2. 可以查看：
   - 消息速率
   - 队列数量
   - 连接数
   - 通道数

## 9. 最佳实践

### 9.1 消息可靠性

```go
// 开启消息持久化
err = ch.PublishWithContext(l.ctx,
    "notification_exchange",
    "notification.send",
    false,
    false,
    amqp091.Publishing{
        ContentType:  "application/json",
        Body:         []byte(message),
        DeliveryMode: amqp091.Persistent,
    })
```

### 9.2 消费者确认机制

```go
// 手动确认消息
msgs, err := ch.Consume(
    queueName,
    "",
    false,
    false,
    false,
    false,
    nil,
)

for d := range msgs {
    err := handleMessage(d.Body)
    if err != nil {
        d.Nack(false, true)
    } else {
        d.Ack(false)
    }
}
```

### 9.3 连接重试

```go
// 自动重连机制
func NewRabbitMQWithRetry(cfg Config, maxRetries int) (*amqp091.Connection, error) {
    var conn *amqp091.Connection
    var err error

    for i := 0; i < maxRetries; i++ {
        conn, err = amqp091.Dial(buildConnStr(cfg))
        if err == nil {
            return conn, nil
        }
        time.Sleep(time.Duration(i+1) * 5 * time.Second)
    }

    return nil, err
}
```

### 9.4 消息序列化

```go
// 使用 JSON 序列化消息
type NotificationMessage struct {
    UserID  int64  `json:"user_id"`
    Type    string `json:"type"`
    Content string `json:"content"`
}

func (m *NotificationMessage) ToBytes() ([]byte, error) {
    return json.Marshal(m)
}

func ParseNotificationMessage(data []byte) (*NotificationMessage, error) {
    var msg NotificationMessage
    err := json.Unmarshal(data, &msg)
    return &msg, err
}
```

## 10. 故障排查

### 10.1 常见问题

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 连接失败 | RabbitMQ 服务未启动 | 检查容器状态，重启服务 |
| 消息丢失 | 未开启持久化 | 设置 DeliveryMode 为 Persistent |
| 消息重复 | 消费者未正确确认 | 确保处理成功后调用 Ack |
| 延迟队列不工作 | 未安装延迟插件 | 安装 rabbitmq_delayed_message_exchange 插件 |

### 10.2 日志查看

```bash
# 查看 RabbitMQ 日志
docker logs zero-rabbitmq

# 查看队列状态
curl http://localhost:15672/api/queues/notification.send -u admin:admin123456

# 查看消息统计
curl http://localhost:15672/api/queues/notification.send/get?count=10 -u admin:admin123456
```

### 10.3 性能优化建议

1. **预创建连接和通道**：避免频繁创建和销毁连接
2. **批量发送消息**：减少网络开销
3. **合理设置队列参数**：根据业务需求设置持久化、过期时间等
4. **监控队列长度**：及时发现消息堆积问题

## 11. 测试用例

### 11.1 发送测试消息

```bash
# 通过管理 API 发送测试消息
curl -X POST http://localhost:15672/api/exchanges/%2f/notification_exchange/publish \
  -u admin:admin123456 \
  -H "Content-Type: application/json" \
  -d '{
    "routing_key": "notification.send",
    "payload": "{\"user_id\":1,\"type\":\"sms\",\"content\":\"Test message\"}",
    "payload_encoding": "string",
    "properties": {
      "delivery_mode": 2
    }
  }'
```

### 11.2 验证消息消费

```bash
# 查看队列消息数量
curl http://localhost:15672/api/queues/notification.send -u admin:admin123456 | jq '.messages'
```

---

**说明**: 本教程覆盖了 RabbitMQ 在项目中的基本使用场景，包括普通消息、延迟消息的发送和消费。在实际使用中，请根据业务需求进行适当调整。