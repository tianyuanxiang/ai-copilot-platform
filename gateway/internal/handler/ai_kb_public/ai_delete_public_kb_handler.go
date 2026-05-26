// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_public

import (
	"net/http"

	"ai-copilot-platform/gateway/internal/logic/ai_kb_public"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiDeletePublicKbHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AiPublicKbPathReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := ai_kb_public.NewAiDeletePublicKbLogic(r.Context(), svcCtx)
		resp, err := l.AiDeletePublicKb(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
