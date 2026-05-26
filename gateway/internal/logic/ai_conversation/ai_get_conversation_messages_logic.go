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

type AiGetConversationMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetConversationMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetConversationMessagesLogic {
	return &AiGetConversationMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetConversationMessagesLogic) AiGetConversationMessages(req *types.AiConversationMessagesReq) (resp *types.AiConversationMessagesResp, err error) {
	userID, err := userIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.ConversationId) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 不能为空")
	}
	rpcResp, err := l.svcCtx.AiChatClient.GetConversationMessages(l.ctx, &aichatclient.GetConversationMessagesReq{
		UserId:         userID,
		ConversationId: req.ConversationId,
		Page:           req.Page,
		PageSize:       req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AiConversationMessagesResp{
		Total: rpcResp.Total,
		List:  make([]types.AiMessageItem, 0, len(rpcResp.List)),
	}
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		resp.List = append(resp.List, types.AiMessageItem{
			MessageId:      item.MessageId,
			ConversationId: item.ConversationId,
			Role:           item.Role,
			Content:        item.Content,
			Citations:      citationsFromRPC(item.Citations),
			TraceId:        item.TraceId,
			CreatedAt:      item.CreatedAt,
		})
	}
	return resp, nil
}
