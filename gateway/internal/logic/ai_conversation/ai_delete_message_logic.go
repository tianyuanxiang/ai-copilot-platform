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

type AiDeleteMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeleteMessageLogic {
	return &AiDeleteMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeleteMessageLogic) AiDeleteMessage(req *types.AiDeleteMessageReq) (resp *types.AiCommonResp, err error) {
	userID, err := userIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.ConversationId) == "" || req.MessageId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 和 messageId 不能为空")
	}
	if _, err := l.svcCtx.AiChatClient.DeleteMessage(l.ctx, &aichatclient.DeleteMessageReq{
		UserId:         userID,
		ConversationId: req.ConversationId,
		MessageId:      req.MessageId,
	}); err != nil {
		return nil, err
	}
	return &types.AiCommonResp{Message: "ok"}, nil
}
