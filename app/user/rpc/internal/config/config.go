package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// MySQL 连接串
	DataSource string
	// Redis 缓存（go-zero model 缓存）
	Cache cache.CacheConf
}
