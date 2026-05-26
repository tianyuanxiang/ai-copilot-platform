// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_conversation

import (
	"net/http"

	"ai-copilot-platform/gateway/internal/logic/ai_conversation"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiDeleteConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AiConversationPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ai_conversation.NewAiDeleteConversationLogic(r.Context(), svcCtx)
		resp, err := l.AiDeleteConversation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
