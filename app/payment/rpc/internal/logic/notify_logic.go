package logic

import (
	"context"
	"time"

	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyLogic {
	return &NotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyLogic) Notify(in *payment.NotifyPaymentReq) (*payment.NotifyPaymentResp, error) {
	item, err := l.svcCtx.PaymentModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "支付记录不存在")
	}
	if item.Status == int64(payment.PaymentStatus_PAYMENT_STATUS_SUCCESS) {
		return &payment.NotifyPaymentResp{Success: true}, nil
	}

	now := time.Now().Unix()
	if err := l.svcCtx.PaymentModel.UpdateStatus(l.ctx, in.Id, int64(in.Status), in.TransactionId, now); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if in.Status == payment.PaymentStatus_PAYMENT_STATUS_SUCCESS {
		if _, err := l.svcCtx.OrderRpc.UpdateStatus(l.ctx, &order.UpdateOrderStatusReq{
			Id:     item.OrderId,
			Status: order.OrderStatus_ORDER_STATUS_PAID,
		}); err != nil {
			l.Errorf("update order status failed: %v", err)
		}
	}

	return &payment.NotifyPaymentResp{Success: true}, nil
}
