package aiwindagentservicelogic

import (
	"context"
	"fmt"
	"go-zero-rpc/common/xerr"
	"time"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/logic/aiwinddraft"
	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTicketDraftLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTicketDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTicketDraftLogic {
	return &CreateTicketDraftLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTicketDraftLogic) CreateTicketDraft(in *pb.WindTicketDraftReq) (*pb.WindScaffoldResp, error) {
	traceID := model.WindTraceID()

	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(l.ctx, l.svcCtx, in.FarmCode, in.TowerCode, in.AlarmCode, "", "", 0, true)
	if err != nil {
		l.Logger.Errorf("Draft build alarm evidence failed: %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "建立告警事件失败")
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(in.EvidenceJson, alarmEvidence)

	payload := engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(in.UserId),
		TraceID:      traceID,
		FarmCode:     in.FarmCode,
		TowerCode:    in.TowerCode,
		AlarmCode:    in.AlarmCode,
		Priority:     "normal",
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}
	startedAt := time.Now()
	draft, err := l.svcCtx.EngineCallClient.WindTicketDraft(l.ctx, payload)
	if err != nil {
		l.Logger.Errorf("告警数据AI分析失败: %v", err)
		aiwinddraft.WriteToolCallLog(l.ctx, l.svcCtx, in.UserId, traceID, "wind_ticket_draft", payload, map[string]any{}, startedAt, "failed", err.Error())
		return nil, err
	}
	aiwinddraft.WriteToolCallLog(l.ctx, l.svcCtx, in.UserId, traceID, "wind_ticket_draft", payload, draft, startedAt, "success", "")

	status := draft.Status
	if status == "" {
		status = "draft"
	}
	content := aiwinddraft.DraftContent(draft)

	// 数据库只存轻量 evidence 摘要
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)

	id, err := l.svcCtx.AiMaintenanceTicketDraftModel.InsertReturningID(l.ctx, &model.AiMaintenanceTicketDraft{
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
		l.Logger.Errorf("Draft insert failed: %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "插入AI维修工单草稿表失败")
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
