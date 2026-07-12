package payment

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentLogic {
	return &GetPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentLogic) GetPayment(req *types.GetPaymentReq) (resp *types.GetPaymentResp, err error) {
	r, err := l.svcCtx.PaymentRpc.Get(l.ctx, &payment.GetPaymentReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.GetPaymentResp{
		Payment: types.PaymentInfo{
			Id:            r.Payment.Id,
			UserId:        r.Payment.UserId,
			OrderId:       r.Payment.OrderId,
			OrderNo:       r.Payment.OrderNo,
			Amount:        r.Payment.Amount,
			Method:        int32(r.Payment.Method),
			Status:        int32(r.Payment.Status),
			TransactionId: r.Payment.TransactionId,
			CreateAt:      r.Payment.CreateAt,
			UpdateAt:      r.Payment.UpdateAt,
		},
	}, nil
}
