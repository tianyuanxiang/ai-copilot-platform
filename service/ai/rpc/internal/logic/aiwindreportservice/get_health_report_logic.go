package aiwindreportservicelogic

import (
	"context"
	"fmt"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHealthReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHealthReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHealthReportLogic {
	return &GetHealthReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetHealthReportLogic) GetHealthReport(in *pb.WindHealthReportGetReq) (*pb.WindScaffoldResp, error) {
	evidence := map[string]any{"scaffold": true, "report_id": in.ReportId, "todo": "后续从 ai_health_report 查询报告正文和 evidence。"}
	return &pb.WindScaffoldResp{
		Id:           in.ReportId,
		TraceId:      model.WindTraceID(),
		Title:        fmt.Sprintf("健康报告 %d", in.ReportId),
		Content:      "健康报告查询脚手架已预留。",
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      "get health report scaffold ready",
	}, nil
}
