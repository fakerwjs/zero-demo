// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"net/http"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/logic/user"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewAvatarLogic(r.Context(), svcCtx)
		resp, err := l.Avatar(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
