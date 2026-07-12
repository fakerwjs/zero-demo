package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLogic) Get(in *product.GetProductReq) (*product.GetProductResp, error) {
	item, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "商品不存在")
	}
	return &product.GetProductResp{
		Product: &product.ProductInfo{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       item.Stock,
			Image:       item.Image,
			CreateAt:    item.CreatedAt,
			UpdateAt:    item.UpdatedAt,
		},
	}, nil
}
