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
	// RabbitMQ（发送通知异步投递）
	RabbitMQ struct {
		URL           string
		Queue         string
		DelayQueue    string
		DelayExchange string
	}
	// MongoDB（通知审计日志）
	Mongo struct {
		URI        string
		Database   string
		Collection string
	}
	// 真实投递网关（邮件/短信）。留空则降级为仅日志。
	Sender struct {
		TestEmail  string
		TestMobile string
		SMTP       struct {
			Host     string
			Port     int
			Username string
			Password string
			From     string
		}
		SMS struct {
			Endpoint  string
			AccessKey string
			SignName  string
			Template  string
		}
	}
}
