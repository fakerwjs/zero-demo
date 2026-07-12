// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/user/rpc/user_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	r, err := l.svcCtx.UserRpc.Register(l.ctx, &user_client.RegisterReq{
		Username: req.Username,
		Password: req.Password,
		Mobile:   req.Mobile,
	})
	if err != nil {
		return nil, err
	}
	return &types.RegisterResp{
		Id:       r.UserInfo.Id,
		Username: r.UserInfo.Username,
		Mobile:   r.UserInfo.Mobile,
	}, nil
}
