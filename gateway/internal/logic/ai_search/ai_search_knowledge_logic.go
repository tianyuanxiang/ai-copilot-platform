// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_search

import (
	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"context"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"
	"strings"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiSearchKnowledgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiSearchKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiSearchKnowledgeLogic {
	return &AiSearchKnowledgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiSearchKnowledgeLogic) AiSearchKnowledge(req *types.AiSearchKnowledgeReq) (resp *types.AiSearchKnowledgeResp, err error) {

	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if strings.TrimSpace(req.Query) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "query 不能为空")
	}
	if req.TopK > 20 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "TopK 不能大于20")
	}
	if req.TopK <= 0 {
		req.TopK = 5
	}

	searchResult, err := l.svcCtx.AiKnowledgeClient.SearchKnowledge(l.ctx, &aiknowledgeclient.SearchKnowledgeReq{
		UserId:      userId,
		HasKbId:     req.KbId > 0,
		KbId:        req.KbId,
		Query:       req.Query,
		TopK:        int64(req.TopK),
		AnswerMode:  req.AnswerMode,
		SearchScope: req.SearchScope,
		HasDomainId: req.DomainId > 0,
		DomainId:    req.DomainId,
		DocumentIds: req.DocumentIds,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.AiChunkItem, 0, len(searchResult.Chunks))
	if searchResult.Chunks == nil {
		return &types.AiSearchKnowledgeResp{
			Chunks:  nil,
			Mode:    searchResult.Mode,
			Message: searchResult.Message,
		}, err
	}
	for _, doc := range searchResult.Chunks {
		items = append(items, types.AiChunkItem{
			DocumentId: doc.DocumentId,
			ChunkId:    doc.ChunkId,
			Title:      doc.Title,
			Snippet:    doc.Snippet,
			Content:    doc.Content,
			Score:      doc.Score,
			Source:     doc.Source,
		})
	}

	return &types.AiSearchKnowledgeResp{
		Chunks:  items,
		Mode:    searchResult.Mode,
		Message: searchResult.Message,
	}, err
}
