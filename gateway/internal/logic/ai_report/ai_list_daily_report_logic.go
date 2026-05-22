// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_report

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListDailyReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListDailyReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListDailyReportLogic {
	return &AiListDailyReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListDailyReportLogic) AiListDailyReport(req *types.AiDailyReportQueryReq) (resp *types.AiListDailyReportResp, err error) {
	// todo: add your logic here and delete this line

	return
}
