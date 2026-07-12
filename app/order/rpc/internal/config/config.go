package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf
	ProductRpc zrpc.RpcClientConf
	RabbitMQ   RabbitMQConf
}

type RabbitMQConf struct {
	Addr       string
	Exchange   string
	Queue      string
	RoutingKey string
}
