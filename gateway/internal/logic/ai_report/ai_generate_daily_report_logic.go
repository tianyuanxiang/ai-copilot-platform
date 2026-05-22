// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_report

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGenerateDailyReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGenerateDailyReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGenerateDailyReportLogic {
	return &AiGenerateDailyReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGenerateDailyReportLogic) AiGenerateDailyReport(req *types.AiDailyReportReq) (resp *types.AiDailyReportResp, err error) {
	// todo: add your logic here and delete this line

	return
}
