package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationLogic {
	return &GetNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNotificationLogic) GetNotification(in *notification.GetNotificationReq) (*notification.GetNotificationResp, error) {
	n, err := l.svcCtx.NotificationModel.FindOne(l.ctx, in.Id)
	if err == model.ErrNotFound {
		return nil, status.Error(codes.NotFound, "通知不存在")
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &notification.GetNotificationResp{
		Notification: &notification.NotificationInfo{
			Id:       n.Id,
			UserId:   n.UserId,
			Title:    n.Title,
			Content:  n.Content,
			Channel:  notification.Channel(n.Channel),
			IsRead:   n.IsRead == 1,
			CreateAt: n.CreateAt.Unix(),
		},
	}, nil
}
