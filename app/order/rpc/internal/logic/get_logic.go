package logic

import (
	"context"
	"encoding/json"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLogic) Get(in *order.GetOrderReq) (*order.GetOrderResp, error) {
	item, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "订单不存在")
	}
	var items []*order.OrderItem
	_ = json.Unmarshal([]byte(item.Items), &items)
	return &order.GetOrderResp{
		Order: &order.OrderInfo{
			Id:         item.Id,
			UserId:     item.UserId,
			OrderNo:    item.OrderNo,
			TotalPrice: item.TotalPrice,
			Status:     order.OrderStatus(item.Status),
			Items:      items,
			CreateAt:   item.CreatedAt,
			UpdateAt:   item.UpdatedAt,
		},
	}, nil
}
