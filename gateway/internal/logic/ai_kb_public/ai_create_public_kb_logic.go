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

type AiCreatePublicKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCreatePublicKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCreatePublicKbLogic {
	return &AiCreatePublicKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCreatePublicKbLogic) AiCreatePublicKb(req *types.AiCreatePublicKbReq) (resp *types.AiCreatePublicKbResp, err error) {
	// 1. 提取用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if req.DomainId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "domainId 不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "name 不能为空")
	}

	// 3. 调用RPC创建公共知识库
	rpcResp, err := l.svcCtx.AiKnowledgeClient.CreateKnowledgeBase(l.ctx, &aiknowledgeclient.CreateKnowledgeBaseReq{
		KbType:      "public",
		DomainId:    req.DomainId,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Visibility:  strings.TrimSpace(req.Visibility),
		Status:      int64(req.Status),
		UserId:      userID,
		OperatorId:  userID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 返回创建结果
	return &types.AiCreatePublicKbResp{
		KbId: rpcResp.KbId,
	}, nil
}
