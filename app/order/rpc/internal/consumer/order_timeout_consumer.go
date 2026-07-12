package consumer

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/fakerwjs/zero-demo/app/order/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/order/rpc/order"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderTimeoutConsumer struct {
	svcCtx *svc.ServiceContext
}

func NewOrderTimeoutConsumer(svcCtx *svc.ServiceContext) *OrderTimeoutConsumer {
	return &OrderTimeoutConsumer{
		svcCtx: svcCtx,
	}
}

func (c *OrderTimeoutConsumer) Consume(msg string) error {
	orderId, err := strconv.ParseInt(msg, 10, 64)
	if err != nil {
		logx.Error("Failed to parse order ID: %v", err)
		return nil
	}

	ctx := context.Background()
	orderInfo, err := c.svcCtx.OrderModel.FindOne(ctx, orderId)
	if err != nil {
		logx.Error("Failed to find order: %v", err)
		return nil
	}

	if orderInfo.Status != int64(order.OrderStatus_ORDER_STATUS_PENDING) {
		logx.Info("Order %d is not pending, skip cancel", orderId)
		return nil
	}

	var items []order.OrderItem
	if orderInfo.Items != "" {
		if err := json.Unmarshal([]byte(orderInfo.Items), &items); err != nil {
			logx.Error("Failed to parse order items: %v", err)
			return nil
		}
	}

	for _, item := range items {
		if _, err := c.svcCtx.ProductRpc.DeductStock(ctx, &product.DeductStockReq{Id: item.ProductId, Amount: -item.Quantity}); err != nil {
			logx.Error("Failed to restock product %d: %v", item.ProductId, err)
		}
	}

	if err := c.svcCtx.OrderModel.UpdateStatus(ctx, orderId, int64(order.OrderStatus_ORDER_STATUS_CANCELLED), time.Now().Unix()); err != nil {
		logx.Error("Failed to cancel order %d: %v", orderId, err)
		return err
	}

	logx.Info("Order %d cancelled due to timeout", orderId)
	return nil
}
