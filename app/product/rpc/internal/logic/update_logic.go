package logic

import (
	"context"
	"time"

	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *product.UpdateProductReq) (*product.UpdateProductResp, error) {
	item, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "商品不存在")
	}
	if in.Name != "" {
		item.Name = in.Name
	}
	if in.Description != "" {
		item.Description = in.Description
	}
	if in.Price > 0 {
		item.Price = in.Price
	}
	if in.Stock >= 0 {
		item.Stock = in.Stock
	}
	if in.Image != "" {
		item.Image = in.Image
	}
	item.UpdatedAt = time.Now().Unix()
	if err := l.svcCtx.ProductModel.Update(l.ctx, item); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &product.UpdateProductResp{Success: true}, nil
}
