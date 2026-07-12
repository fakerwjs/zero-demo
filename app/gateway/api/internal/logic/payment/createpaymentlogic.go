package payment

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaymentLogic) CreatePayment(req *types.CreatePaymentReq) (resp *types.CreatePaymentResp, err error) {
	userId := jwtx.UserIDFromCtx(l.ctx)

	orderInfo, err := l.svcCtx.OrderRpc.Get(l.ctx, &order.GetOrderReq{Id: req.OrderId})
	if err != nil {
		return nil, err
	}

	r, err := l.svcCtx.PaymentRpc.Create(l.ctx, &payment.CreatePaymentReq{
		UserId:  userId,
		OrderId: req.OrderId,
		OrderNo: orderInfo.Order.OrderNo,
		Amount:  req.Amount,
		Method:  payment.PaymentMethod(req.Method),
	})
	if err != nil {
		return nil, err
	}

	return &types.CreatePaymentResp{
		Id:      r.Id,
		PayUrl:  r.PayUrl,
		Success: r.Success,
	}, nil
}
