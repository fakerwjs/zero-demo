package logic

import (
	"context"
	"time"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStatusLogic) UpdateStatus(in *order.UpdateOrderStatusReq) (*order.UpdateOrderStatusResp, error) {
	item, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "订单不存在")
	}
	if !canTransition(item.Status, int64(in.Status)) {
		return nil, status.Error(codes.InvalidArgument, "状态不允许变更")
	}
	if err := l.svcCtx.OrderModel.UpdateStatus(l.ctx, in.Id, int64(in.Status), time.Now().Unix()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.UpdateOrderStatusResp{Success: true}, nil
}

func canTransition(from, to int64) bool {
	switch from {
	case int64(order.OrderStatus_ORDER_STATUS_PENDING):
		return to == int64(order.OrderStatus_ORDER_STATUS_PAID) || to == int64(order.OrderStatus_ORDER_STATUS_CANCELLED)
	case int64(order.OrderStatus_ORDER_STATUS_PAID):
		return to == int64(order.OrderStatus_ORDER_STATUS_SHIPPED) || to == int64(order.OrderStatus_ORDER_STATUS_CANCELLED)
	case int64(order.OrderStatus_ORDER_STATUS_SHIPPED):
		return to == int64(order.OrderStatus_ORDER_STATUS_COMPLETED)
	case int64(order.OrderStatus_ORDER_STATUS_COMPLETED):
		return false
	case int64(order.OrderStatus_ORDER_STATUS_CANCELLED):
		return false
	default:
		return false
	}
}
