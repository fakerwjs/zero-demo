// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/config"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/middleware"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification_client"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order_client"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment_client"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product_client"
	"github.com/fakerwjs/zero-demo/app/user/rpc/user_client"
	"github.com/fakerwjs/zero-demo/pkg/miniox"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	UserRpc         user_client.User
	NotificationRpc notification_client.Notification
	ProductRpc      product_client.Product
	OrderRpc        order_client.Order
	PaymentRpc      payment_client.Payment
	Minio           *miniox.Client
	RateLimit       rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	// MinIO 客户端（失败不阻断启动，仅头像上传不可用）
	var mc *miniox.Client
	m, err := miniox.New(miniox.Config{
		Endpoint:        c.Minio.Endpoint,
		AccessKeyID:     c.Minio.AccessKeyID,
		SecretAccessKey: c.Minio.SecretAccessKey,
		UseSSL:          c.Minio.UseSSL,
		Bucket:          c.Minio.Bucket,
		PublicBaseURL:   c.Minio.PublicBaseURL,
	})
	if err != nil {
		logx.Errorf("init minio failed: %v", err)
	} else {
		mc = m
	}

	return &ServiceContext{
		Config:          c,
		UserRpc:         user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
		NotificationRpc: notification_client.NewNotification(zrpc.MustNewClient(c.NotificationRpc)),
		ProductRpc:      product_client.NewProduct(c.ProductRpc),
		OrderRpc:        order_client.NewOrder(c.OrderRpc),
		PaymentRpc:      payment_client.NewPayment(c.PaymentRpc),
		Minio:           mc,
		RateLimit:       middleware.NewRateLimitMiddleware(c.RateLimit.Rate, c.RateLimit.Burst).Handle,
	}
}
