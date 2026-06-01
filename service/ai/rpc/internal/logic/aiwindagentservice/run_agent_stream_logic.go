package aiwindagentservicelogic

import (
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunAgentStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunAgentStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunAgentStreamLogic {
	return &RunAgentStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RunAgentStreamLogic) RunAgentStream(in *pb.WindAgentRunReq, stream pb.AiWindAgentService_RunAgentStreamServer) error {
	if in == nil {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "请求不能为空")
	}
	if in.UserId <= 0 {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "userId 必须大于 0")
	}
	input := strings.TrimSpace(in.Input)
	if input == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "input 不能为空")
	}

	req := engine.WindAgentRunRequest{
		UserID:         in.UserId,
		AgentSessionId: strings.TrimSpace(in.AgentSessionId),
		Input:          input,
	}

	// 定义一个函数：Python 每返回一条 event，就执行这个函数
	handleEvent := func(event engine.WindAgentStreamEvent) error {
		return sendWindAgentStreamEvent(stream, event)
	}

	// Python Agent 负责推理和状态编排，RPC 层只校验请求并逐条转发 SSE 事件。
	return l.svcCtx.EngineCallClient.EngineWindAgentStream(l.ctx, req, handleEvent)
}
