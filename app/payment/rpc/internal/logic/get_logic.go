package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"

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

func (l *GetLogic) Get(in *payment.GetPaymentReq) (*payment.GetPaymentResp, error) {
	item, err := l.svcCtx.PaymentModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "支付记录不存在")
	}
	return &payment.GetPaymentResp{
		Payment: &payment.PaymentInfo{
			Id:            item.Id,
			UserId:        item.UserId,
			OrderId:       item.OrderId,
			OrderNo:       item.OrderNo,
			Amount:        item.Amount,
			Method:        payment.PaymentMethod(item.Method),
			Status:        payment.PaymentStatus(item.Status),
			TransactionId: item.TransactionId,
			CreateAt:      item.CreatedAt,
			UpdateAt:      item.UpdatedAt,
		},
	}, nil
}
