package product_client

import (
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Product interface {
	product.ProductClient
}

type productClient struct {
	product.ProductClient
}

func NewProduct(c zrpc.RpcClientConf) Product {
	return &productClient{
		ProductClient: product.NewProductClient(zrpc.MustNewClient(c).Conn()),
	}
}

func NewProductWithConn(conn *grpc.ClientConn) Product {
	return &productClient{
		ProductClient: product.NewProductClient(conn),
	}
}
