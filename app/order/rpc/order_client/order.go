package order_client

import (
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Order interface {
	order.OrderClient
}

type orderClient struct {
	order.OrderClient
}

func NewOrder(c zrpc.RpcClientConf) Order {
	return &orderClient{
		OrderClient: order.NewOrderClient(zrpc.MustNewClient(c).Conn()),
	}
}

func NewOrderWithConn(conn *grpc.ClientConn) Order {
	return &orderClient{
		OrderClient: order.NewOrderClient(conn),
	}
}
