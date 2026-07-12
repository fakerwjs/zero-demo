// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	notificationpb "github.com/fakerwjs/zero-demo/app/notification/rpc/notification"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification_client"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNotificationLogic {
	return &SendNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendNotificationLogic) SendNotification(req *types.SendNotificationReq) (resp *types.SendNotificationResp, err error) {
	userId := jwtx.UserIDFromCtx(l.ctx)
	r, err := l.svcCtx.NotificationRpc.Send(l.ctx, &notification_client.SendReq{
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
		Channel: notificationpb.Channel(req.Channel),
	})
	if err != nil {
		return nil, err
	}
	return &types.SendNotificationResp{
		Id:      r.Id,
		Success: r.Success,
	}, nil
}
