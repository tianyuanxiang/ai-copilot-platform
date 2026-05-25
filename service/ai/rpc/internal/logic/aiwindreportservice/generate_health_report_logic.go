package aiwindreportservicelogic

import (
	"context"

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
	evidence := map[string]any{
		"scaffold":       true,
		"report_type":    in.ReportType,
		"farm_code":      in.FarmCode,
		"tower_code":     in.TowerCode,
		"start_time":     in.StartTime,
		"end_time":       in.EndTime,
		"input_evidence": in.EvidenceJson,
		"todo":           "后续补充在线率、告警统计、异常测点、SOP 引用和 LLM 总结。",
	}
	return &pb.WindScaffoldResp{
		TraceId:      traceID,
		Title:        "风机健康报告草稿",
		Content:      "健康报告生成脚手架已预留。一期返回结构化草稿，后续补齐日报、周报、单机报告和故障复盘。",
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      "health report scaffold ready",
	}, nil
}
