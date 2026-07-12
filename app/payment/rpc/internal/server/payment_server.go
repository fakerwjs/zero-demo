package server

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/logic"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServer
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaymentServer(svcCtx *svc.ServiceContext) *PaymentServer {
	return &PaymentServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

func (s *PaymentServer) Create(ctx context.Context, in *payment.CreatePaymentReq) (*payment.CreatePaymentResp, error) {
	return logic.NewCreateLogic(ctx, s.svcCtx).Create(in)
}

func (s *PaymentServer) Get(ctx context.Context, in *payment.GetPaymentReq) (*payment.GetPaymentResp, error) {
	return logic.NewGetLogic(ctx, s.svcCtx).Get(in)
}

func (s *PaymentServer) List(ctx context.Context, in *payment.ListPaymentReq) (*payment.ListPaymentResp, error) {
	return logic.NewListLogic(ctx, s.svcCtx).List(in)
}

func (s *PaymentServer) Notify(ctx context.Context, in *payment.NotifyPaymentReq) (*payment.NotifyPaymentResp, error) {
	return logic.NewNotifyLogic(ctx, s.svcCtx).Notify(in)
}
