// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_report

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetDailyReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetDailyReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetDailyReportLogic {
	return &AiGetDailyReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetDailyReportLogic) AiGetDailyReport(req *types.AiDailyReportPathReq) (resp *types.AiDailyReportResp, err error) {
	// todo: add your logic here and delete this line

	return
}
