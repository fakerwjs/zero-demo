package product

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (resp *types.CreateProductResp, err error) {
	r, err := l.svcCtx.ProductRpc.Create(l.ctx, &product.CreateProductReq{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Image:       req.Image,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateProductResp{
		Id:      r.Id,
		Success: r.Success,
	}, nil
}
