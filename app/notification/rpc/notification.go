package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/config"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/server"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification"
	"github.com/fakerwjs/zero-demo/pkg/mq"
	"github.com/fakerwjs/zero-demo/pkg/sender"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// buildDispatcher 根据配置构建邮件/短信投递器；未配置凭证的渠道降级为仅日志。
func buildDispatcher(c config.Config) *sender.Dispatcher {
	d := &sender.Dispatcher{
		TestEmail:  c.Sender.TestEmail,
		TestMobile: c.Sender.TestMobile,
	}
	if c.Sender.SMTP.Host != "" {
		d.Email = sender.NewEmailSender(sender.SMTPConfig{
			Host:     c.Sender.SMTP.Host,
			Port:     c.Sender.SMTP.Port,
			Username: c.Sender.SMTP.Username,
			Password: c.Sender.SMTP.Password,
			From:     c.Sender.SMTP.From,
		})
	}
	if c.Sender.SMS.Endpoint != "" {
		d.SMS = sender.NewSMSSender(sender.SMSConfig{
			Endpoint:  c.Sender.SMS.Endpoint,
			AccessKey: c.Sender.SMS.AccessKey,
			SignName:  c.Sender.SMS.SignName,
			Template:  c.Sender.SMS.Template,
		})
	}
	return d
}

var configFile = flag.String("f", "etc/notification.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// 启动 RabbitMQ 消费者：按渠道分发到真实邮件/短信网关（未配置凭证则仅日志）
	dispatcher := buildDispatcher(c)
	mq.NewConsumer(c.RabbitMQ.URL, c.RabbitMQ.Queue, func(body []byte) error {
		var m struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Channel int32  `json:"channel"`
		}
		if err := json.Unmarshal(body, &m); err != nil {
			logx.Errorf("[notification] bad message: %s", string(body))
			return nil // 丢弃坏消息，不重回队列
		}
		return dispatcher.Dispatch(context.Background(), m.Channel, m.Title, m.Content)
	}).Start()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		notification.RegisterNotificationServer(grpcServer, server.NewNotificationServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
