package aisecurityservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSecurityAlertStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSecurityAlertStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSecurityAlertStatusLogic {
	return &UpdateSecurityAlertStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新安全告警处理状态。
func (l *UpdateSecurityAlertStatusLogic) UpdateSecurityAlertStatus(in *pb.UpdateSecurityAlertStatusReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
