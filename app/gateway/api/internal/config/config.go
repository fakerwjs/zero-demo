// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// 后端 RPC（etcd 服务发现）
	UserRpc         zrpc.RpcClientConf
	NotificationRpc zrpc.RpcClientConf
	ProductRpc      zrpc.RpcClientConf
	OrderRpc        zrpc.RpcClientConf
	PaymentRpc      zrpc.RpcClientConf

	// 限流：每 IP 每秒允许的请求数与突发值
	RateLimit struct {
		Rate  int
		Burst int
	}

	// MinIO 对象存储（头像上传）
	Minio struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		Bucket          string
		PublicBaseURL   string
	}
}
