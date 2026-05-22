package aireportservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDailyReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDailyReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDailyReportLogic {
	return &GetDailyReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询单份 AI 安全日报详情。
func (l *GetDailyReportLogic) GetDailyReport(in *pb.GetDailyReportReq) (*pb.DailyReportResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DailyReportResp{}, nil
}
