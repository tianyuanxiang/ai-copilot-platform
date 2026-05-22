// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_chat

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiChatLogic {
	return &AiChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiChatLogic) AiChat(req *types.AiChatReq) (resp *types.AiChatResp, err error) {
	// todo: add your logic here and delete this line

	return
}
