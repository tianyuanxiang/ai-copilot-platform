// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_chat

import (
	aichatclient "ai-copilot-platform/ai-rpc/client/aichatservice"
	"ai-copilot-platform/ai-rpc/pb"
	"context"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"
	"strings"

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
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if strings.TrimSpace(req.Question) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "question 不能为空")
	}

	rpcResp, err := l.svcCtx.AiChatClient.RagChat(l.ctx, &aichatclient.RagChatReq{
		UserId:         userID,
		HasKbId:        req.KbId > 0,
		KbId:           req.KbId,
		ConversationId: req.ConversationId,
		Question:       req.Question,
		AnswerMode:     req.AnswerMode,
		SearchScope:    req.SearchScope,
		HasDomainId:    req.DomainId > 0,
		DomainId:       req.DomainId,
		DocumentIds:    req.DocumentIds,
	})
	if err != nil {
		return nil, err
	}

	return &types.AiChatResp{
		Answer:         rpcResp.Answer,
		TraceId:        rpcResp.TraceId,
		Mode:           rpcResp.Mode,
		ConversationId: rpcResp.ConversationId,
		Citations:      citationsFromRPC(rpcResp.Citations),
	}, nil
}

func citationsFromRPC(citations []*pb.Citation) []types.AiCitation {
	items := make([]types.AiCitation, 0, len(citations))
	for _, item := range citations {
		if item == nil {
			continue
		}
		items = append(items, types.AiCitation{
			DocumentId: item.DocumentId,
			ChunkId:    item.ChunkId,
			Title:      item.Title,
			Snippet:    item.Snippet,
			Score:      item.Score,
		})
	}
	return items
}
