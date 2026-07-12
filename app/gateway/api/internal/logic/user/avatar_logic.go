// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/fakerwjs/zero-demo/pkg/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxAvatarSize = 5 << 20 // 5MB

type AvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AvatarLogic {
	return &AvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Avatar 接收 multipart 表单字段 file，上传到 MinIO，返回可访问 URL。
func (l *AvatarLogic) Avatar(r *http.Request) (resp *types.AvatarResp, err error) {
	if l.svcCtx.Minio == nil {
		return nil, errors.New("对象存储未就绪")
	}

	if err = r.ParseMultipartForm(maxAvatarSize); err != nil {
		return nil, fmt.Errorf("解析上传表单失败: %w", err)
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("缺少文件字段 file: %w", err)
	}
	defer file.Close()

	userId := jwtx.UserIDFromCtx(l.ctx)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}
	objectName := fmt.Sprintf("user_%d/avatar%s", userId, ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := l.svcCtx.Minio.Put(l.ctx, objectName, contentType, file, header.Size)
	if err != nil {
		return nil, fmt.Errorf("上传对象存储失败: %w", err)
	}

	return &types.AvatarResp{Url: url}, nil
}
