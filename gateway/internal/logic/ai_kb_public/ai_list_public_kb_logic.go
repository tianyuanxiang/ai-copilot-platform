// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_public

import (
	"context"
	"strings"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListPublicKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListPublicKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListPublicKbLogic {
	return &AiListPublicKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListPublicKbLogic) AiListPublicKb(req *types.AiListPublicKbReq) (resp *types.AiListPublicKbResp, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 调用RPC查询公共知识库列表
	visibility := strings.TrimSpace(req.Visibility)
	rpcResp, err := l.svcCtx.AiKnowledgeClient.ListKnowledgeBase(l.ctx, &aiknowledgeclient.ListKnowledgeBaseReq{
		Page:         int64(req.Page),
		PageSize:     int64(req.PageSize),
		UserId:       userID,
		KbType:       "public",
		HasKbType:    true,
		DomainId:     req.DomainId,
		HasDomainId:  req.DomainId > 0,
		Keyword:      strings.TrimSpace(req.Keyword),
		Visibility:   visibility,
		HasVisibility: visibility != "",
		Status:       int64(req.Status),
		HasStatus:    req.Status != 0,
	})
	if err != nil {
		return nil, err
	}

	// 3. 映射RPC响应到网关类型
	items := make([]types.AiPublicKbItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		items = append(items, publicKbItemFromRPC(item))
	}

	// 4. 返回列表结果
	return &types.AiListPublicKbResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
