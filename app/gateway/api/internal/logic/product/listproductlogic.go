package product

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductReq) (resp *types.ListProductResp, err error) {
	r, err := l.svcCtx.ProductRpc.List(l.ctx, &product.ListProductReq{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.ProductInfo, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.ProductInfo{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       item.Stock,
			Image:       item.Image,
			CreateAt:    item.CreateAt,
			UpdateAt:    item.UpdateAt,
		})
	}

	return &types.ListProductResp{
		Total: r.Total,
		List:  list,
	}, nil
}
