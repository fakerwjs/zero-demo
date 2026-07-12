// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"net/http"

	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/logic/notification"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/svc"
	"github.com/fakerwjs/zero-demo/app/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetNotificationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notification.NewGetNotificationLogic(r.Context(), svcCtx)
		resp, err := l.GetNotification(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
