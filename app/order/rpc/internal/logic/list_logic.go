package logic

import (
	"context"
	"encoding/json"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *order.ListOrderReq) (*order.ListOrderResp, error) {
	list, err := l.svcCtx.OrderModel.FindPageByUser(l.ctx, in.UserId, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.OrderModel.FindPageByUserCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	var resp []*order.OrderInfo
	for _, item := range list {
		var items []*order.OrderItem
		_ = json.Unmarshal([]byte(item.Items), &items)
		resp = append(resp, &order.OrderInfo{
			Id:         item.Id,
			UserId:     item.UserId,
			OrderNo:    item.OrderNo,
			TotalPrice: item.TotalPrice,
			Status:     order.OrderStatus(item.Status),
			Items:      items,
			CreateAt:   item.CreatedAt,
			UpdateAt:   item.UpdatedAt,
		})
	}
	return &order.ListOrderResp{
		List:  resp,
		Total: total,
	}, nil
}
