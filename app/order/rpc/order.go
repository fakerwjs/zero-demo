package main

import (
	"flag"
	"fmt"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/config"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/consumer"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/server"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	if ctx.RabbitMQ != nil {
		consumer := consumer.NewOrderTimeoutConsumer(ctx)
		if err := ctx.RabbitMQ.Consume(consumer.Consume); err != nil {
			logx.Error("Failed to start order timeout consumer: %v", err)
		} else {
			logx.Info("Order timeout consumer started")
		}
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order.RegisterOrderServer(grpcServer, server.NewOrderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	if ctx.RabbitMQ != nil {
		defer ctx.RabbitMQ.Close()
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
