// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_conversation

import (
	"context"

	aichatclient "ai-copilot-platform/ai-rpc/client/aichatservice"
	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListConversationLogic {
	return &AiListConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListConversationLogic) AiListConversation(req *types.AiListConversationReq) (resp *types.AiListConversationResp, err error) {
	userID, err := userIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.AiChatClient.ListConversation(l.ctx, &aichatclient.ListConversationReq{
		UserId:   userID,
		Page:     req.Page,
		PageSize: req.PageSize,
		KbId:     req.KbId,
		HasKbId:  req.KbId > 0,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AiListConversationResp{
		Total: rpcResp.Total,
		List:  make([]types.AiConversationItem, 0, len(rpcResp.List)),
	}
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		resp.List = append(resp.List, types.AiConversationItem{
			ConversationId: item.ConversationId,
			KbId:           item.KbId,
			Title:          item.Title,
			LatestMessage:  item.LatestMessage,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return resp, nil
}

func userIDFromCtx(ctx context.Context) (int64, error) {
	userID := middleware.GetUserIdFromCtx(ctx)
	if userID <= 0 {
		return 0, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	return userID, nil
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
