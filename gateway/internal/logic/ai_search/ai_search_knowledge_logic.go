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
	if req.KbId <= 0 || strings.TrimSpace(req.Query) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId、query 不能为空")
	}
	if req.TopK > 20 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "TopK 不能大于20")
	}
	if req.TopK <= 0 {
		req.TopK = 5
	}
	l.svcCtx.AiKnowledgeClient.SearchKnowledge(l.ctx, &aiknowledgeclient.SearchKnowledgeReq{})
	return
}
