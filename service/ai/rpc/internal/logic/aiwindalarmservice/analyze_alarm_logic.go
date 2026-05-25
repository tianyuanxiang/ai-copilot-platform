package aiwindalarmservicelogic

import (
	"context"

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
	evidence := map[string]any{
		"scaffold":       true,
		"farm_code":      in.FarmCode,
		"tower_code":     in.TowerCode,
		"alarm_code":     in.AlarmCode,
		"input_evidence": in.EvidenceJson,
		"todo":           "后续补充告警前后时序窗口、SOP 检索、历史案例召回和归因排序。",
	}
	return &pb.WindScaffoldResp{
		TraceId:      traceID,
		Title:        "告警归因草稿",
		Content:      "告警归因脚手架已预留。一期只记录输入证据和后续补全步骤。",
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      "alarm analysis scaffold ready",
	}, nil
}
