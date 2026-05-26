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

type AiUpdatePersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdatePersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdatePersonalKbLogic {
	return &AiUpdatePersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiUpdatePersonalKbLogic) AiUpdatePersonalKb(req *types.AiUpdatePersonalKbReq) (resp *types.AiCommonResp, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}

	// 3. 构建更新请求，根据可选字段是否提供设置has_xxx标志
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if _, err := l.svcCtx.AiKnowledgeClient.UpdateKnowledgeBase(l.ctx, &aiknowledgeclient.UpdateKnowledgeBaseReq{
		KbId:           req.KbId,
		Name:           name,
		HasName:        name != "",
		Description:    description,
		HasDescription: description != "",
		Status:         int64(req.Status),
		HasStatus:      req.Status != 0,
		UserId:         userID,
		OperatorId:     userID,
	}); err != nil {
		return nil, err
	}

	// 4. 返回更新成功响应
	return &types.AiCommonResp{Message: "ok"}, nil
}
