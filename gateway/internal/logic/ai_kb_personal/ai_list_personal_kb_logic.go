// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

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

type AiListPersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListPersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListPersonalKbLogic {
	return &AiListPersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListPersonalKbLogic) AiListPersonalKb(req *types.AiListPersonalKbReq) (resp *types.AiListPersonalKbResp, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 调用RPC查询个人知识库列表
	rpcResp, err := l.svcCtx.AiKnowledgeClient.ListKnowledgeBase(l.ctx, &aiknowledgeclient.ListKnowledgeBaseReq{
		Page:       int64(req.Page),
		PageSize:   int64(req.PageSize),
		UserId:     userID,
		KbType:     "personal",
		HasKbType:  true,
		Keyword:    strings.TrimSpace(req.Keyword),
		Status:     int64(req.Status),
		HasStatus:  req.Status != 0,
	})
	if err != nil {
		return nil, err
	}

	// 3. 映射RPC响应到网关类型
	items := make([]types.AiPersonalKbItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		items = append(items, personalKbItemFromRPC(item))
	}

	// 4. 返回列表结果
	return &types.AiListPersonalKbResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
