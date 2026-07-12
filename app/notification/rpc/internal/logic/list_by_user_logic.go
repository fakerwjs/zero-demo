package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/notification/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/notification/rpc/notification"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListByUserLogic {
	return &ListByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListByUserLogic) ListByUser(in *notification.ListByUserReq) (*notification.ListByUserResp, error) {
	page, size := in.Page, in.Size
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size

	rows, err := l.svcCtx.NotificationModel.ListByUserId(l.ctx, in.UserId, offset, size)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	total, err := l.svcCtx.NotificationModel.CountByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	list := make([]*notification.NotificationInfo, 0, len(rows))
	for _, n := range rows {
		list = append(list, &notification.NotificationInfo{
			Id:       n.Id,
			UserId:   n.UserId,
			Title:    n.Title,
			Content:  n.Content,
			Channel:  notification.Channel(n.Channel),
			IsRead:   n.IsRead == 1,
			CreateAt: n.CreateAt.Unix(),
		})
	}

	return &notification.ListByUserResp{
		List:  list,
		Total: total,
	}, nil
}
