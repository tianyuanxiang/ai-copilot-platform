// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_agent

import (
	"go-zero-rpc/common/response"
	"net/http"

	"ai-copilot-platform/gateway/internal/logic/ai_wind_agent"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiWindCopilotAgentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AiWindAgentRunReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := ai_wind_agent.NewAiWindCopilotAgentLogic(r.Context(), svcCtx)
		err := l.AiWindCopilotAgent(&req, w)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response.OK(w, r)
		}
	}
}
