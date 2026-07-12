package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *product.ListProductReq) (*product.ListProductResp, error) {
	list, err := l.svcCtx.ProductModel.FindPageList(l.ctx, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.ProductModel.FindPageListCount(l.ctx)
	if err != nil {
		return nil, err
	}
	var resp []*product.ProductInfo
	for _, item := range list {
		resp = append(resp, &product.ProductInfo{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       item.Stock,
			Image:       item.Image,
			CreateAt:    item.CreatedAt,
			UpdateAt:    item.UpdatedAt,
		})
	}
	return &product.ListProductResp{
		List:  resp,
		Total: total,
	}, nil
}
