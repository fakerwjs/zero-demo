package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"

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

func (l *ListLogic) List(in *payment.ListPaymentReq) (*payment.ListPaymentResp, error) {
	list, err := l.svcCtx.PaymentModel.FindPageByUser(l.ctx, in.UserId, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.PaymentModel.FindPageByUserCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	var resp []*payment.PaymentInfo
	for _, item := range list {
		resp = append(resp, &payment.PaymentInfo{
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
		})
	}
	return &payment.ListPaymentResp{
		List:  resp,
		Total: total,
	}, nil
}
