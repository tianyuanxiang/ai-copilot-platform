package aiwindagentservicelogic

import (
	"context"

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
	evidence := map[string]any{
		"scaffold":       true,
		"farm_code":      in.FarmCode,
		"tower_code":     in.TowerCode,
		"alarm_code":     in.AlarmCode,
		"input_evidence": in.EvidenceJson,
		"todo":           "后续补充检查项、备件、人员、工单系统提交前确认。",
	}
	return &pb.WindScaffoldResp{
		TraceId:      traceID,
		Title:        "维修工单草稿",
		Content:      "维修工单草稿脚手架已预留。一期不自动提交工单。",
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      "maintenance ticket draft scaffold ready",
	}, nil
}
