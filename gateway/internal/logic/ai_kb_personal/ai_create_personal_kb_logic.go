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

type AiCreatePersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCreatePersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCreatePersonalKbLogic {
	return &AiCreatePersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCreatePersonalKbLogic) AiCreatePersonalKb(req *types.AiCreatePersonalKbReq) (resp *types.AiCreatePersonalKbResp, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if strings.TrimSpace(req.Name) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "name 不能为空")
	}

	// 3. 调用RPC创建个人知识库
	rpcResp, err := l.svcCtx.AiKnowledgeClient.CreateKnowledgeBase(l.ctx, &aiknowledgeclient.CreateKnowledgeBaseReq{
		KbType:      "personal",
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		UserId:      userID,
		OperatorId:  userID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 返回创建结果
	return &types.AiCreatePersonalKbResp{
		KbId: rpcResp.KbId,
	}, nil
}
