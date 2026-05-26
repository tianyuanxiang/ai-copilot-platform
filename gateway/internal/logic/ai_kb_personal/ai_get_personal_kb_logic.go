// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

import (
	"context"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetPersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetPersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetPersonalKbLogic {
	return &AiGetPersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetPersonalKbLogic) AiGetPersonalKb(req *types.AiPersonalKbPathReq) (resp *types.AiPersonalKbItem, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}

	// 3. 调用RPC查询知识库详情
	item, err := l.svcCtx.AiKnowledgeClient.GetKnowledgeBase(l.ctx, &aiknowledgeclient.GetKnowledgeBaseReq{
		KbId:   req.KbId,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 映射RPC响应到网关类型并返回
	result := personalKbItemFromRPC(item)
	return &result, nil
}
