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

type AiDeletePersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeletePersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeletePersonalKbLogic {
	return &AiDeletePersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeletePersonalKbLogic) AiDeletePersonalKb(req *types.AiPersonalKbPathReq) (resp *types.AiCommonResp, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}

	// 3. 调用RPC删除知识库
	if _, err := l.svcCtx.AiKnowledgeClient.DeleteKnowledgeBase(l.ctx, &aiknowledgeclient.DeleteKnowledgeBaseReq{
		KbId:       req.KbId,
		UserId:     userID,
		OperatorId: userID,
	}); err != nil {
		return nil, err
	}

	// 4. 返回删除成功响应
	return &types.AiCommonResp{Message: "ok"}, nil
}
