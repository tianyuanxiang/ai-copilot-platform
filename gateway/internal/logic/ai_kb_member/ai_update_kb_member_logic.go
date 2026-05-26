// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_member

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

type AiUpdateKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdateKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdateKbMemberLogic {
	return &AiUpdateKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AiUpdateKbMember 更新公共知识库成员角色
//
// 提取当前用户ID作为操作者，校验 kbId、userId、role 参数，
// 调用 RPC UpdateKbMember 完成角色更新
func (l *AiUpdateKbMemberLogic) AiUpdateKbMember(req *types.AiUpdateKbMemberReq) (resp *types.AiCommonResp, err error) {
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
	if strings.TrimSpace(req.Role) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "role 不能为空")
	}

	// 3. 调用 RPC 更新成员角色
	if _, err := l.svcCtx.AiKnowledgeClient.UpdateKbMember(l.ctx, &aiknowledgeclient.UpdateKbMemberReq{
		KbId:       req.KbId,
		UserId:     req.UserId,
		Role:       strings.TrimSpace(req.Role),
		OperatorId: userID,
	}); err != nil {
		return nil, err
	}

	return &types.AiCommonResp{Message: "ok"}, nil
}
