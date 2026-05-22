package aichatservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RagChatStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRagChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RagChatStreamLogic {
	return &RagChatStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 执行流式 RAG 问答，逐步返回 token、引用和完成事件。
func (l *RagChatStreamLogic) RagChatStream(in *pb.RagChatReq, stream pb.AiChatService_RagChatStreamServer) error {
	// todo: add your logic here and delete this line

	return nil
}
