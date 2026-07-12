package logic

import (
	"context"
	"time"

	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/product/rpc/internal/svc"
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

func (l *CreateLogic) Create(in *product.CreateProductReq) (*product.CreateProductResp, error) {
	now := time.Now().Unix()
	res, err := l.svcCtx.ProductModel.Insert(l.ctx, &model.Product{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		Image:       in.Image,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, _ := res.LastInsertId()
	return &product.CreateProductResp{
		Id:      id,
		Success: true,
	}, nil
}
