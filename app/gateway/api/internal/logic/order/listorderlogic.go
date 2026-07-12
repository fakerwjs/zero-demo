package order

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogic {
	return &ListOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrderLogic) ListOrder(req *types.ListOrderReq) (resp *types.ListOrderResp, err error) {
	userId := jwtx.UserIDFromCtx(l.ctx)
	r, err := l.svcCtx.OrderRpc.List(l.ctx, &order.ListOrderReq{
		UserId: userId,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.OrderInfo, 0, len(r.List))
	for _, item := range r.List {
		items := make([]types.OrderItem, 0, len(item.Items))
		for _, i := range item.Items {
			items = append(items, types.OrderItem{
				ProductId: i.ProductId,
				Name:      i.Name,
				Price:     i.Price,
				Quantity:  i.Quantity,
			})
		}
		list = append(list, types.OrderInfo{
			Id:         item.Id,
			UserId:     item.UserId,
			OrderNo:    item.OrderNo,
			TotalPrice: item.TotalPrice,
			Status:     int32(item.Status),
			Items:      items,
			CreateAt:   item.CreateAt,
			UpdateAt:   item.UpdateAt,
		})
	}

	return &types.ListOrderResp{
		Total: r.Total,
		List:  list,
	}, nil
}
