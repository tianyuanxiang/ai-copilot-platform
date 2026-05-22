// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_report

import (
	"net/http"

	"ai-copilot-platform/gateway/internal/logic/ai_report"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiGetDailyReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AiDailyReportPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ai_report.NewAiGetDailyReportLogic(r.Context(), svcCtx)
		resp, err := l.AiGetDailyReport(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
