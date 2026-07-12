// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification_client"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationLogic {
	return &ListNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNotificationLogic) ListNotification(req *types.ListNotificationReq) (resp *types.ListNotificationResp, err error) {
	userId := jwtx.UserIDFromCtx(l.ctx)
	r, err := l.svcCtx.NotificationRpc.ListByUser(l.ctx, &notification_client.ListByUserReq{
		UserId: userId,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.NotificationItem, 0, len(r.List))
	for _, n := range r.List {
		list = append(list, types.NotificationItem{
			Id:       n.Id,
			UserId:   n.UserId,
			Title:    n.Title,
			Content:  n.Content,
			Channel:  int32(n.Channel),
			IsRead:   n.IsRead,
			CreateAt: n.CreateAt,
		})
	}
	return &types.ListNotificationResp{
		Total: r.Total,
		List:  list,
	}, nil
}
