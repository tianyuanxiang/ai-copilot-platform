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

type AiDeleteConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeleteConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeleteConversationLogic {
	return &AiDeleteConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeleteConversationLogic) AiDeleteConversation(req *types.AiConversationPathReq) (resp *types.AiCommonResp, err error) {
	userID, err := userIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.ConversationId) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 不能为空")
	}
	if _, err := l.svcCtx.AiChatClient.DeleteConversation(l.ctx, &aichatclient.DeleteConversationReq{
		UserId:         userID,
		ConversationId: req.ConversationId,
	}); err != nil {
		return nil, err
	}
	return &types.AiCommonResp{Message: "ok"}, nil
}
