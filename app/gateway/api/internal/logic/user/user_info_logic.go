// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/app/user/rpc/user_client"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResp, err error) {
	userId := jwtx.UserIDFromCtx(l.ctx)
	r, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user_client.GetUserInfoReq{Id: userId})
	if err != nil {
		return nil, err
	}
	return &types.UserInfoResp{
		Id:       r.UserInfo.Id,
		Username: r.UserInfo.Username,
		Mobile:   r.UserInfo.Mobile,
		Nickname: r.UserInfo.Nickname,
		CreateAt: r.UserInfo.CreateAt,
	}, nil
}
