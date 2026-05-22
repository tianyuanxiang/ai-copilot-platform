package aireportservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDailyReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDailyReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDailyReportLogic {
	return &ListDailyReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询 AI 安全日报列表。
func (l *ListDailyReportLogic) ListDailyReport(in *pb.ListDailyReportReq) (*pb.ListDailyReportResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListDailyReportResp{}, nil
}
