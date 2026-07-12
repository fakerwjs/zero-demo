package server

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/logic"
	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductServer struct {
	product.UnimplementedProductServer
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductServer(svcCtx *svc.ServiceContext) *ProductServer {
	return &ProductServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

func (s *ProductServer) Create(ctx context.Context, in *product.CreateProductReq) (*product.CreateProductResp, error) {
	return logic.NewCreateLogic(ctx, s.svcCtx).Create(in)
}

func (s *ProductServer) Get(ctx context.Context, in *product.GetProductReq) (*product.GetProductResp, error) {
	return logic.NewGetLogic(ctx, s.svcCtx).Get(in)
}

func (s *ProductServer) List(ctx context.Context, in *product.ListProductReq) (*product.ListProductResp, error) {
	return logic.NewListLogic(ctx, s.svcCtx).List(in)
}

func (s *ProductServer) Update(ctx context.Context, in *product.UpdateProductReq) (*product.UpdateProductResp, error) {
	return logic.NewUpdateLogic(ctx, s.svcCtx).Update(in)
}

func (s *ProductServer) Delete(ctx context.Context, in *product.DeleteProductReq) (*product.DeleteProductResp, error) {
	return logic.NewDeleteLogic(ctx, s.svcCtx).Delete(in)
}

func (s *ProductServer) DeductStock(ctx context.Context, in *product.DeductStockReq) (*product.DeductStockResp, error) {
	return logic.NewDeductStockLogic(ctx, s.svcCtx).DeductStock(in)
}
