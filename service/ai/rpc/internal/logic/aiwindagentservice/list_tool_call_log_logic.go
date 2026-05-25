package aiwindagentservicelogic

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

func (l *ListToolCallLogLogic) ListToolCallLog(in *pb.WindListToolCallLogReq) (*pb.WindListToolCallLogResp, error) {
	// 一期先返回空列表，表结构和查询入口已预留；后续在工具执行处写入 ai_tool_call_log 后补分页查询。
	return &pb.WindListToolCallLogResp{Total: 0, List: []*pb.ToolCall{}}, nil
}
