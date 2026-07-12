package svc

import (
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/config"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/mq"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product_client"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.OrderModel
	ProductRpc product_client.Product
	RabbitMQ   *mq.RabbitMQ
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	var rmq *mq.RabbitMQ
	if c.RabbitMQ.Addr != "" {
		var err error
		rmq, err = mq.NewRabbitMQ(
			c.RabbitMQ.Addr,
			c.RabbitMQ.Exchange,
			c.RabbitMQ.Queue,
			c.RabbitMQ.RoutingKey,
		)
		if err != nil {
			logx.Error("Failed to connect to RabbitMQ, order timeout feature disabled: %v", err)
			rmq = nil
		} else {
			logx.Info("Connected to RabbitMQ successfully")
		}
	}

	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(conn, c.Cache),
		ProductRpc: product_client.NewProduct(c.ProductRpc),
		RabbitMQ:   rmq,
	}
}
