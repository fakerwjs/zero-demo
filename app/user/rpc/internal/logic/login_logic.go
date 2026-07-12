package logic

import (
	"context"

	"github.com/fakerwjs/zero-demo/app/user/rpc/internal/model"
	"github.com/fakerwjs/zero-demo/app/user/rpc/internal/svc"
	"github.com/fakerwjs/zero-demo/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err == model.ErrNotFound {
		return nil, status.Error(codes.NotFound, "用户不存在")
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 校验密码哈希
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}

	// JWT 令牌由 user-api 层用其 AccessSecret 统一签发，RPC 只返回用户身份。
	return &user.LoginResp{
		Id:           u.Id,
		AccessToken:  "",
		AccessExpire: 0,
	}, nil
}
