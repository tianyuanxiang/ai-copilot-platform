package aiwindreportservicelogic

import (
	"context"
	"fmt"
	"time"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/logic/aiwinddraft"
	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateHealthReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateHealthReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateHealthReportLogic {
	return &GenerateHealthReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateHealthReportLogic) GenerateHealthReport(in *pb.WindHealthReportReq) (*pb.WindScaffoldResp, error) {
	traceID := model.WindTraceID()

	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(l.ctx, l.svcCtx, in.FarmCode, in.TowerCode, "", in.StartTime, in.EndTime, 0, false)
	if err != nil {
		l.Logger.Errorf("Generate health build alarm evidence failed: %v", err)
		return nil, err
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(in.EvidenceJson, alarmEvidence)

	payload := engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(in.UserId),
		TraceID:      traceID,
		FarmCode:     in.FarmCode,
		TowerCode:    in.TowerCode,
		ReportType:   in.ReportType,
		StartTime:    in.StartTime,
		EndTime:      in.EndTime,
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}
	startedAt := time.Now()
	draft, err := l.svcCtx.EngineCallClient.WindHealthReportDraft(l.ctx, payload)
	if err != nil {
		l.Logger.Errorf("健康报告AI分析失败: %v", err)
		aiwinddraft.WriteToolCallLog(l.ctx, l.svcCtx, in.UserId, traceID, "wind_health_report_draft", payload, map[string]any{}, startedAt, "failed", err.Error())
		return nil, err
	}
	aiwinddraft.WriteToolCallLog(l.ctx, l.svcCtx, in.UserId, traceID, "wind_health_report_draft", payload, draft, startedAt, "success", "")

	status := draft.Status
	if status == "" {
		status = "draft"
	}
	reportType := in.ReportType
	if reportType == "" {
		reportType = "health"
	}
	content := aiwinddraft.DraftContent(draft)

	// 数据库只存轻量 evidence 摘要
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)

	id, err := l.svcCtx.AiHealthReportModel.InsertReturningID(l.ctx, &model.AiHealthReport{
		UserId:     in.UserId,
		TraceId:    traceID,
		ReportType: reportType,
		FarmCode:   in.FarmCode,
		TowerCode:  in.TowerCode,
		StartTime:  aiwinddraft.ParseNullableTime(in.StartTime),
		EndTime:    aiwinddraft.ParseNullableTime(in.EndTime),
		Title:      draft.Title,
		Content:    content,
		Evidence:   evidenceSummary,
		Status:     status,
	})
	if err != nil {
		l.Logger.Errorf("AiHealthReportModel insertReturningID failed: %v", err)
		return nil, err
	}

	// API 响应只返回轻量 evidence 摘要
	return &pb.WindScaffoldResp{
		Id:           id,
		TraceId:      traceID,
		Title:        draft.Title,
		Content:      content,
		EvidenceJson: evidenceSummary,
		Message:      fmt.Sprintf("%s; status=%s", draft.Message, status),
	}, nil
}
