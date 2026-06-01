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

type ResumeAgentStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResumeAgentStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeAgentStreamLogic {
	return &ResumeAgentStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResumeAgentStreamLogic) ResumeAgentStream(in *pb.WindAgentResumeReq, stream pb.AiWindAgentService_ResumeAgentStreamServer) error {
	if in == nil {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "请求不能为空")
	}
	if in.UserId <= 0 {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "userId 必须大于 0")
	}
	agentSessionId := strings.TrimSpace(in.AgentSessionId)
	if agentSessionId == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 不能为空")
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	switch action {
	case "approve", "reject":
	case "clarify":
		if strings.TrimSpace(in.Content) == "" {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "clarify 必须携带 content")
		}
	default:
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "action 仅支持 approve、reject 或 clarify")
	}

	// 恢复请求仍由 Python checkpoint 接管，Go 不解释审批或澄清业务。
	return l.svcCtx.EngineCallClient.EngineWindAgentResumeStream(l.ctx, engine.WindAgentResumeRequest{
		UserID:         in.UserId,
		AgentSessionId: agentSessionId,
		Action:         action,
		Content:        strings.TrimSpace(in.Content),
	}, func(event engine.WindAgentStreamEvent) error {
		return sendWindAgentStreamEvent(stream, event)
	})
}
