package aiagentservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListToolCallLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListToolCallLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListToolCallLogLogic {
	return &ListToolCallLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询 Agent 工具调用日志。
func (l *ListToolCallLogLogic) ListToolCallLog(in *pb.ListToolCallLogReq) (*pb.ListToolCallLogResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListToolCallLogResp{}, nil
}
