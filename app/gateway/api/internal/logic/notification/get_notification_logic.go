// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationLogic {
	return &GetNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationLogic) GetNotification(req *types.GetNotificationReq) (resp *types.GetNotificationResp, err error) {
	r, err := l.svcCtx.NotificationRpc.GetNotification(l.ctx, &notification_client.GetNotificationReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	n := r.Notification
	return &types.GetNotificationResp{
		Notification: types.NotificationItem{
			Id:       n.Id,
			UserId:   n.UserId,
			Title:    n.Title,
			Content:  n.Content,
			Channel:  int32(n.Channel),
			IsRead:   n.IsRead,
			CreateAt: n.CreateAt,
		},
	}, nil
}
