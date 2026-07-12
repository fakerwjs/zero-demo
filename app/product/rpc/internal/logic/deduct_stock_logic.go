package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductStockLogic) DeductStock(in *product.DeductStockReq) (*product.DeductStockResp, error) {
	if in.Amount == 0 {
		return nil, status.Error(codes.InvalidArgument, "数量不能为0")
	}
	if in.Amount > 0 {
		rows, err := l.svcCtx.ProductModel.DeductStock(l.ctx, in.Id, in.Amount)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if rows == 0 {
			return nil, status.Error(codes.InvalidArgument, "库存不足")
		}
	} else {
		productInfo, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
		if err != nil {
			return nil, status.Error(codes.NotFound, "商品不存在")
		}
		productInfo.Stock -= in.Amount
		if err := l.svcCtx.ProductModel.Update(l.ctx, productInfo); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &product.DeductStockResp{Success: true}, nil
}
