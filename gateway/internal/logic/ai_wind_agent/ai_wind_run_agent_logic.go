// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_agent

import (
	"context"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindRunAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindRunAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindRunAgentLogic {
	return &AiWindRunAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindRunAgentLogic) AiWindRunAgent(req *types.AiWindAgentRunReq) (resp *types.AiWindAgentRunResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	result, err := l.svcCtx.AiWindAgentClient.RunAgent(l.ctx, &pb.WindAgentRunReq{
		AgentSessionId: req.AgentSessionId,
		Input:          req.Input,
		UserId:         userID,
	})
	if err != nil {
		return nil, err
	}
	calls := make([]types.AiWindToolCallItem, 0, len(result.ToolCalls))
	for _, item := range result.ToolCalls {
		if item == nil {
			continue
		}
		calls = append(calls, types.AiWindToolCallItem{
			ToolCallId:    item.ToolCallId,
			ToolName:      item.ToolName,
			Status:        item.Status,
			ArgumentsJson: item.ArgumentsJson,
			ResultJson:    item.ResultJson,
			Message:       item.Message,
			CreatedAt:     item.CreatedAt,
		})
	}
	return &types.AiWindAgentRunResp{Answer: result.Answer, TraceId: result.TraceId, ToolCalls: calls}, nil
}
