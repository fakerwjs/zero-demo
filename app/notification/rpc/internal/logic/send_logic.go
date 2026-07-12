package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification"
	"github.com/fakerwjs/zero-demo/pkg/mongox"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendLogic) Send(in *notification.SendReq) (*notification.SendResp, error) {
	// 1. 落库（默认未读）
	res, err := l.svcCtx.NotificationModel.Insert(l.ctx, &model.Notification{
		UserId:  in.UserId,
		Title:   in.Title,
		Content: in.Content,
		Channel: int64(in.Channel),
		IsRead:  0,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, _ := res.LastInsertId()

	// 2. 投递到 RabbitMQ，由消费者异步按渠道真正发送（短信/邮件/推送）
	msg, _ := json.Marshal(map[string]any{
		"id":      id,
		"userId":  in.UserId,
		"title":   in.Title,
		"content": in.Content,
		"channel": in.Channel,
	})
	if err := l.svcCtx.Publisher.PublishWithDelay(l.ctx, msg, 0); err != nil {
		// MQ 投递失败不阻断主流程，记录后续可补偿
		l.Errorf("publish notification to mq failed: %v", err)
	}

	// 3. 写 MongoDB 审计日志（不阻断主流程）
	if l.svcCtx.Mongo != nil {
		if err := l.svcCtx.Mongo.InsertNotificationLog(l.ctx, &mongox.NotificationLog{
			NotificationID: id,
			UserID:         in.UserId,
			Title:          in.Title,
			Content:        in.Content,
			Channel:        int32(in.Channel),
			CreatedAt:      time.Now(),
		}); err != nil {
			l.Errorf("write mongo audit log failed: %v", err)
		}
	}

	return &notification.SendResp{
		Id:      id,
		Success: true,
	}, nil
}
