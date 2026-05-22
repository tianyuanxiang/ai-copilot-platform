// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package health

import (
	"go-zero-rpc/common/response"
	"net/http"

	"ai-copilot-platform/gateway/internal/logic/health"
	"ai-copilot-platform/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewAiHealthLogic(r.Context(), svcCtx)
		resp, err := l.AiHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
