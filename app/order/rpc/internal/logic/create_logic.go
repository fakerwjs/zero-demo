package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

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

func (l *CreateLogic) Create(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	if len(in.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "订单项不能为空")
	}

	var totalPrice int64
	for _, item := range in.Items {
		if item.Quantity <= 0 {
			return nil, status.Error(codes.InvalidArgument, "商品数量必须大于0")
		}
		productInfo, err := l.svcCtx.ProductRpc.Get(l.ctx, &product.GetProductReq{Id: item.ProductId})
		if err != nil {
			return nil, status.Error(codes.NotFound, "商品不存在")
		}
		item.Name = productInfo.Product.Name
		item.Price = productInfo.Product.Price
		totalPrice += item.Price * item.Quantity
	}

	for _, item := range in.Items {
		if _, err := l.svcCtx.ProductRpc.DeductStock(l.ctx, &product.DeductStockReq{Id: item.ProductId, Amount: item.Quantity}); err != nil {
			return nil, err
		}
	}

	itemsJSON, _ := json.Marshal(in.Items)
	orderNo := generateOrderNo()
	now := time.Now().Unix()

	res, err := l.svcCtx.OrderModel.Insert(l.ctx, &model.Order{
		UserId:     in.UserId,
		OrderNo:    orderNo,
		TotalPrice: totalPrice,
		Status:     int64(order.OrderStatus_ORDER_STATUS_PENDING),
		Items:      string(itemsJSON),
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, _ := res.LastInsertId()

	if l.svcCtx.RabbitMQ != nil {
		go func() {
			delay := 15 * time.Minute
			if err := l.svcCtx.RabbitMQ.PublishDelayMessage(strconv.FormatInt(id, 10), delay); err != nil {
				logx.Error("Failed to publish order timeout message: %v", err)
			}
		}()
	}

	return &order.CreateOrderResp{
		Id:      id,
		OrderNo: orderNo,
		Success: true,
	}, nil
}

var orderNoCounter int64

func generateOrderNo() string {
	timestamp := time.Now().UnixMilli()
	seq := atomic.AddInt64(&orderNoCounter, 1) % 10000
	return "ORD" + strconv.FormatInt(timestamp, 10) + fmt.Sprintf("%04d", seq)
}
