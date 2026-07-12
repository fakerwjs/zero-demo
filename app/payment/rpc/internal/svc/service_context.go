package svc

import (
	"github.com/fakerwjs/zero-demo/app/order/rpc/order_client"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/config"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	PaymentModel model.PaymentModel
	OrderRpc     order_client.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:       c,
		PaymentModel: model.NewPaymentModel(conn, c.Cache),
		OrderRpc:     order_client.NewOrder(c.OrderRpc),
	}
}
