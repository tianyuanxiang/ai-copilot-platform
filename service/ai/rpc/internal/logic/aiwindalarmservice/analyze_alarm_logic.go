package aiwindalarmservicelogic

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

type AnalyzeAlarmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnalyzeAlarmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalyzeAlarmLogic {
	return &AnalyzeAlarmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnalyzeAlarmLogic) AnalyzeAlarm(in *pb.WindAlarmAnalyzeReq) (*pb.WindScaffoldResp, error) {
	traceID := model.WindTraceID()

	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(l.ctx, l.svcCtx, in.FarmCode, in.TowerCode, in.AlarmCode, in.StartTime, in.EndTime, 0, true)
	if err != nil {
		return nil, err
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(in.EvidenceJson, alarmEvidence)

	payload := engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(in.UserId),
		TraceID:      traceID,
		FarmCode:     in.FarmCode,
		TowerCode:    in.TowerCode,
		AlarmCode:    in.AlarmCode,
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}

	startedAt := time.Now()
	draft, err := l.svcCtx.EngineCallClient.WindAlarmSummary(l.ctx, payload)
	if err != nil {
		aiwinddraft.WriteToolCallLog(l.ctx, l.svcCtx, in.UserId, traceID, "wind_alarm_analysis", payload, map[string]any{}, startedAt, "failed", err.Error())
		return nil, err
	}
	aiwinddraft.WriteToolCallLog(l.ctx, l.svcCtx, in.UserId, traceID, "wind_alarm_analysis", payload, draft, startedAt, "success", "")

	status := draft.Status
	if status == "" {
		status = "draft"
	}
	content := aiwinddraft.DraftContent(draft)

	// 数据库只存轻量 evidence 摘要，不存全量原始告警记录
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)

	id, err := l.svcCtx.AiAlarmAnalysisModel.InsertReturningID(l.ctx, &model.AiAlarmAnalysis{
		UserId:    in.UserId,
		TraceId:   traceID,
		FarmCode:  in.FarmCode,
		TowerCode: in.TowerCode,
		AlarmCode: in.AlarmCode,
		Title:     draft.Title,
		Content:   content,
		Evidence:  evidenceSummary,
		Status:    status,
	})
	if err != nil {
		return nil, err
	}

	// API 响应只返回轻量 evidence 摘要，不回传原始查询数据
	return &pb.WindScaffoldResp{
		Id:           id,
		TraceId:      traceID,
		Title:        draft.Title,
		Content:      content,
		EvidenceJson: evidenceSummary,
		Message:      fmt.Sprintf("%s; status=%s", draft.Message, status),
	}, nil
}
