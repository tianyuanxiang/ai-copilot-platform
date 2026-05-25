// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_report

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindGenerateHealthReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindGenerateHealthReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindGenerateHealthReportLogic {
	return &AiWindGenerateHealthReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindGenerateHealthReportLogic) AiWindGenerateHealthReport(req *types.AiWindHealthReportReq) (resp *types.AiWindScaffoldResp, err error) {
	result, err := l.svcCtx.AiWindReportClient.GenerateHealthReport(l.ctx, &pb.WindHealthReportReq{
		ReportType:   req.ReportType,
		FarmCode:     req.FarmCode,
		TowerCode:    req.TowerCode,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		EvidenceJson: req.EvidenceJson,
		UserId:       0,
	})
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
