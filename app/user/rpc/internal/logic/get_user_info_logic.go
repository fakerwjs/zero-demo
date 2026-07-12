package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/user/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/user/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err == model.ErrNotFound {
		return nil, status.Error(codes.NotFound, "用户不存在")
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.GetUserInfoResp{
		UserInfo: &user.UserInfo{
			Id:       u.Id,
			Username: u.Username,
			Mobile:   u.Mobile,
			Nickname: u.Nickname,
			CreateAt: u.CreateAt.Unix(),
		},
	}, nil
}
