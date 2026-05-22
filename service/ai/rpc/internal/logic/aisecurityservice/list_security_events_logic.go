package aisecurityservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSecurityEventsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSecurityEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSecurityEventsLogic {
	return &ListSecurityEventsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询结构化安全事件列表。
func (l *ListSecurityEventsLogic) ListSecurityEvents(in *pb.ListSecurityEventsReq) (*pb.ListSecurityEventsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListSecurityEventsResp{}, nil
}
