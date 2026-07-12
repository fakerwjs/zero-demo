// Package sender 提供通知的真实投递能力：邮件(SMTP) 与 短信(HTTP 网关)。
// 未配置凭证时降级为仅日志，便于本地运行。
package sender

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

// 渠道常量，与 proto Channel 对齐。
const (
	ChannelInbox = 1
	ChannelSMS   = 2
	ChannelEmail = 3
	ChannelPush  = 4
)

// Sender 单一渠道投递接口。
type Sender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// Dispatcher 按渠道路由到对应 Sender。
type Dispatcher struct {
	Email Sender
	SMS   Sender
	// 演示用收件人（生产应按 userId 查询用户真实邮箱/手机号）
	TestEmail  string
	TestMobile string
}

// Dispatch 根据渠道分发。inbox/push 目前仅落库+日志。
func (d *Dispatcher) Dispatch(ctx context.Context, channel int32, subject, body string) error {
	switch int(channel) {
	case ChannelEmail:
		if d.Email == nil {
			logx.Infof("[sender] EMAIL(log-only) to=%s subject=%s", d.TestEmail, subject)
			return nil
		}
		return d.Email.Send(ctx, d.TestEmail, subject, body)
	case ChannelSMS:
		if d.SMS == nil {
			logx.Infof("[sender] SMS(log-only) to=%s body=%s", d.TestMobile, body)
			return nil
		}
		return d.SMS.Send(ctx, d.TestMobile, subject, body)
	default:
		logx.Infof("[sender] channel=%d 仅站内信/推送，跳过外发", channel)
		return nil
	}
}
