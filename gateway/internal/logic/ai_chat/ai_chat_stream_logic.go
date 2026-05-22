// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_chat

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiChatStreamLogic {
	return &AiChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiChatStreamLogic) AiChatStream(req *types.AiChatReq) error {
	// todo: add your logic here and delete this line

	return nil
}
