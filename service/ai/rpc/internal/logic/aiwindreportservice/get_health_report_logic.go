package aiwindreportservicelogic

import (
	"context"
	"fmt"
	"go-zero-rpc/common/xerr"

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
	report, err := l.svcCtx.AiHealthReportModel.FindOne(l.ctx, in.ReportId)
	if err != nil {
		l.Logger.Errorf("GetHealthReport err: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if in.UserId > 0 && report.UserId != in.UserId {
		return nil, xerr.NewCodeError(xerr.ErrForbidden)
	}
	return &pb.WindScaffoldResp{
		Id:           in.ReportId,
		TraceId:      report.TraceId,
		Title:        report.Title,
		Content:      report.Content,
		EvidenceJson: report.Evidence,
		Message:      fmt.Sprintf("health report loaded; status=%s", report.Status),
	}, nil
}
