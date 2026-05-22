package aiobservabilityservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLlmTraceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLlmTraceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLlmTraceLogic {
	return &GetLlmTraceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询指定 traceId 下的 LLM 调用链路。
func (l *GetLlmTraceLogic) GetLlmTrace(in *pb.GetLlmTraceReq) (*pb.LlmTraceResp, error) {
	// todo: add your logic here and delete this line

	return &pb.LlmTraceResp{}, nil
}
