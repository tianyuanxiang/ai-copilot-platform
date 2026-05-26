// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_conversation

import (
	"context"
	"strings"

	aichatclient "ai-copilot-platform/ai-rpc/client/aichatservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiUpdateConversationTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdateConversationTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdateConversationTitleLogic {
	return &AiUpdateConversationTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiUpdateConversationTitleLogic) AiUpdateConversationTitle(req *types.AiUpdateConversationTitleReq) (resp *types.AiCommonResp, err error) {
	userID, err := userIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.ConversationId) == "" || strings.TrimSpace(req.Title) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 和 title 不能为空")
	}
	if _, err := l.svcCtx.AiChatClient.UpdateConversationTitle(l.ctx, &aichatclient.UpdateConversationTitleReq{
		UserId:         userID,
		ConversationId: req.ConversationId,
		Title:          req.Title,
	}); err != nil {
		return nil, err
	}
	return &types.AiCommonResp{Message: "ok"}, nil
}
