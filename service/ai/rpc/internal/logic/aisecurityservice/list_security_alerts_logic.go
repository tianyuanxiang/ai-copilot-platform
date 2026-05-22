package aisecurityservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSecurityAlertsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSecurityAlertsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSecurityAlertsLogic {
	return &ListSecurityAlertsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询安全告警列表。
func (l *ListSecurityAlertsLogic) ListSecurityAlerts(in *pb.ListSecurityAlertsReq) (*pb.ListSecurityAlertsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListSecurityAlertsResp{}, nil
}
