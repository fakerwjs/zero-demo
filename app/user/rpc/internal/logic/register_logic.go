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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 用户名唯一性校验
	if _, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username); err == nil {
		return nil, status.Error(codes.AlreadyExists, "用户名已存在")
	} else if err != model.ErrNotFound {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 密码 bcrypt 加盐哈希
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: in.Username,
		Password: string(hashed),
		Mobile:   in.Mobile,
		Nickname: in.Username,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, _ := res.LastInsertId()

	return &user.RegisterResp{
		UserInfo: &user.UserInfo{
			Id:       id,
			Username: in.Username,
			Mobile:   in.Mobile,
			Nickname: in.Username,
		},
	}, nil
}
