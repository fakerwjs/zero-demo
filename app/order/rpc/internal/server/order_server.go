package server

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/logic"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderServer struct {
	order.UnimplementedOrderServer
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

func (s *OrderServer) Create(ctx context.Context, in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	return logic.NewCreateLogic(ctx, s.svcCtx).Create(in)
}

func (s *OrderServer) Get(ctx context.Context, in *order.GetOrderReq) (*order.GetOrderResp, error) {
	return logic.NewGetLogic(ctx, s.svcCtx).Get(in)
}

func (s *OrderServer) List(ctx context.Context, in *order.ListOrderReq) (*order.ListOrderResp, error) {
	return logic.NewListLogic(ctx, s.svcCtx).List(in)
}

func (s *OrderServer) UpdateStatus(ctx context.Context, in *order.UpdateOrderStatusReq) (*order.UpdateOrderStatusResp, error) {
	return logic.NewUpdateStatusLogic(ctx, s.svcCtx).UpdateStatus(in)
}
