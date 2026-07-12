package payment

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPaymentLogic {
	return &ListPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPaymentLogic) ListPayment(req *types.ListPaymentReq) (resp *types.ListPaymentResp, err error) {
	userId := jwtx.UserIDFromCtx(l.ctx)
	r, err := l.svcCtx.PaymentRpc.List(l.ctx, &payment.ListPaymentReq{
		UserId: userId,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.PaymentInfo, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.PaymentInfo{
			Id:            item.Id,
			UserId:        item.UserId,
			OrderId:       item.OrderId,
			OrderNo:       item.OrderNo,
			Amount:        item.Amount,
			Method:        int32(item.Method),
			Status:        int32(item.Status),
			TransactionId: item.TransactionId,
			CreateAt:      item.CreateAt,
			UpdateAt:      item.UpdateAt,
		})
	}

	return &types.ListPaymentResp{
		Total: r.Total,
		List:  list,
	}, nil
}
