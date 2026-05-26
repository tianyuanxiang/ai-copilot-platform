// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_member

import (
	"context"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiRemoveKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiRemoveKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiRemoveKbMemberLogic {
	return &AiRemoveKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AiRemoveKbMember 移除公共知识库成员
//
// 提取当前用户ID作为操作者，校验 kbId、userId 参数，
// 调用 RPC RemoveKbMember 完成成员移除
func (l *AiRemoveKbMemberLogic) AiRemoveKbMember(req *types.AiKbMemberPathReq) (resp *types.AiCommonResp, err error) {
	// 1. 提取当前用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}
	if req.UserId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "userId 不能为空")
	}

	// 3. 调用 RPC 移除成员
	if _, err := l.svcCtx.AiKnowledgeClient.RemoveKbMember(l.ctx, &aiknowledgeclient.RemoveKbMemberReq{
		KbId:       req.KbId,
		UserId:     req.UserId,
		OperatorId: userID,
	}); err != nil {
		return nil, err
	}

	return &types.AiCommonResp{Message: "ok"}, nil
}
