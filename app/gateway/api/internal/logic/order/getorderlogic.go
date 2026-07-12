package order

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (resp *types.GetOrderResp, err error) {
	r, err := l.svcCtx.OrderRpc.Get(l.ctx, &order.GetOrderReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	items := make([]types.OrderItem, 0, len(r.Order.Items))
	for _, item := range r.Order.Items {
		items = append(items, types.OrderItem{
			ProductId: item.ProductId,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
	}

	return &types.GetOrderResp{
		Order: types.OrderInfo{
			Id:         r.Order.Id,
			UserId:     r.Order.UserId,
			OrderNo:    r.Order.OrderNo,
			TotalPrice: r.Order.TotalPrice,
			Status:     int32(r.Order.Status),
			Items:      items,
			CreateAt:   r.Order.CreateAt,
			UpdateAt:   r.Order.UpdateAt,
		},
	}, nil
}
