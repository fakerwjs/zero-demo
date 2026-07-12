package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *payment.CreatePaymentReq) (*payment.CreatePaymentResp, error) {
	if in.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "支付金额必须大于0")
	}

	orderInfo, err := l.svcCtx.OrderRpc.Get(l.ctx, &order.GetOrderReq{Id: in.OrderId})
	if err != nil {
		return nil, status.Error(codes.NotFound, "订单不存在")
	}
	if orderInfo.Order.Status != order.OrderStatus_ORDER_STATUS_PENDING {
		return nil, status.Error(codes.InvalidArgument, "订单状态不允许支付")
	}
	if orderInfo.Order.TotalPrice != in.Amount {
		return nil, status.Error(codes.InvalidArgument, "支付金额与订单金额不一致")
	}
	if orderInfo.Order.UserId != in.UserId {
		return nil, status.Error(codes.PermissionDenied, "无权支付该订单")
	}

	now := time.Now().Unix()
	res, err := l.svcCtx.PaymentModel.Insert(l.ctx, &model.Payment{
		UserId:        in.UserId,
		OrderId:       in.OrderId,
		OrderNo:       orderInfo.Order.OrderNo,
		Amount:        in.Amount,
		Method:        int64(in.Method),
		Status:        int64(payment.PaymentStatus_PAYMENT_STATUS_PENDING),
		TransactionId: "",
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, _ := res.LastInsertId()

	payUrl := generatePayUrl(id)
	return &payment.CreatePaymentResp{
		Id:      id,
		PayUrl:  payUrl,
		Success: true,
	}, nil
}

func generatePayUrl(id int64) string {
	return "/api/v1/payment/pay/" + strconv.FormatInt(id, 10)
}
