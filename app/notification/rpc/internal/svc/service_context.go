package svc

import (
	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/config"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/pkg/mongox"
	"github.com/fakerwjs/zero-demo/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	NotificationModel model.NotificationModel
	Publisher         *mq.Publisher
	Mongo             *mongox.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	// MongoDB 审计日志（失败不阻断启动）
	var mc *mongox.Client
	if m, err := mongox.New(c.Mongo.URI, c.Mongo.Database, c.Mongo.Collection); err != nil {
		logx.Errorf("init mongo failed: %v", err)
	} else {
		mc = m
	}

	return &ServiceContext{
		Config:            c,
		NotificationModel: model.NewNotificationModel(conn, c.Cache),
		Publisher:         mq.NewPublisherWithDelay(c.RabbitMQ.URL, c.RabbitMQ.Queue, c.RabbitMQ.DelayQueue, c.RabbitMQ.DelayExchange),
		Mongo:             mc,
	}
}
