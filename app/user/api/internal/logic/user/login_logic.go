// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"time"

	"github.com/fakerwjs/zero-demo/app/user/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/user/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/user/rpc/user_client"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user_client.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 由 API 层用 AccessSecret 统一签发 JWT，claim 中带上 userId
	now := time.Now().Unix()
	expire := l.svcCtx.Config.Auth.AccessExpire
	token, err := jwtx.BuildToken(l.svcCtx.Config.Auth.AccessSecret, now, expire, rpcResp.Id)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Id:           rpcResp.Id,
		AccessToken:  token,
		AccessExpire: now + expire,
	}, nil
}
