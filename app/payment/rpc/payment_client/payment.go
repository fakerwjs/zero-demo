package payment_client

import (
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Payment interface {
	payment.PaymentClient
}

type paymentClient struct {
	payment.PaymentClient
}

func NewPayment(c zrpc.RpcClientConf) Payment {
	return &paymentClient{
		PaymentClient: payment.NewPaymentClient(zrpc.MustNewClient(c).Conn()),
	}
}

func NewPaymentWithConn(conn *grpc.ClientConn) Payment {
	return &paymentClient{
		PaymentClient: payment.NewPaymentClient(conn),
	}
}
