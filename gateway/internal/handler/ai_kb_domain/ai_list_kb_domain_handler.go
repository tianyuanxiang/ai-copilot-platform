// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_domain

import (
	"net/http"

	"ai-copilot-platform/gateway/internal/logic/ai_kb_domain"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiListKbDomainHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AiListKbDomainReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ai_kb_domain.NewAiListKbDomainLogic(r.Context(), svcCtx)
		resp, err := l.AiListKbDomain(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
