// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_report

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindGetHealthReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindGetHealthReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindGetHealthReportLogic {
	return &AiWindGetHealthReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindGetHealthReportLogic) AiWindGetHealthReport(req *types.AiWindHealthReportPathReq) (resp *types.AiWindScaffoldResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	result, err := l.svcCtx.AiWindReportClient.GetHealthReport(l.ctx, &pb.WindHealthReportGetReq{ReportId: req.ReportId, UserId: userID})
	if err != nil {
		return nil, err
	}
	return &types.AiWindScaffoldResp{
		Id:           result.Id,
		TraceId:      result.TraceId,
		Title:        result.Title,
		Content:      result.Content,
		EvidenceJson: result.EvidenceJson,
		Message:      result.Message,
	}, nil
}
