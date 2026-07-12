package product

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductReq) (resp *types.GetProductResp, err error) {
	r, err := l.svcCtx.ProductRpc.Get(l.ctx, &product.GetProductReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.GetProductResp{
		Product: types.ProductInfo{
			Id:          r.Product.Id,
			Name:        r.Product.Name,
			Description: r.Product.Description,
			Price:       r.Product.Price,
			Stock:       r.Product.Stock,
			Image:       r.Product.Image,
			CreateAt:    r.Product.CreateAt,
			UpdateAt:    r.Product.UpdateAt,
		},
	}, nil
}
