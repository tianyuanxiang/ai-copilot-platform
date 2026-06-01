// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_agent

import (
	"context"
	"net/http"
	"strings"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindResumeAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindResumeAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindResumeAgentLogic {
	return &AiWindResumeAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindResumeAgentLogic) AiWindResumeAgent(req *types.AiWindAgentResumeReq, w http.ResponseWriter) error {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	agentSessionID := strings.TrimSpace(req.AgentSessionId)
	if agentSessionID == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "agentSessionId 不能为空")
	}

	rpcStream, err := l.svcCtx.AiWindAgentClient.ResumeAgentStream(l.ctx, &pb.WindAgentResumeReq{
		UserId:         userID,
		AgentSessionId: agentSessionID,
		Action:         strings.TrimSpace(req.Action),
		Content:        strings.TrimSpace(req.Content),
	})
	if err != nil {
		return err
	}

	return forwardWindAgentStream(w, rpcStream, agentSessionID)
}
