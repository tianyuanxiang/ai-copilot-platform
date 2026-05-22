package aiobservabilityservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLlmCallLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLlmCallLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLlmCallLogLogic {
	return &ListLlmCallLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询 LLM 调用日志。
func (l *ListLlmCallLogLogic) ListLlmCallLog(in *pb.ListLlmCallLogReq) (*pb.ListLlmCallLogResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListLlmCallLogResp{}, nil
}
