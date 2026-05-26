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

type AiListKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListKbMemberLogic {
	return &AiListKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AiListKbMember 查询公共知识库成员列表
//
// 提取当前用户ID，校验 kbId 参数，
// 调用 RPC ListKbMember 获取成员分页数据，
// 通过 DTO 映射返回网关层响应
func (l *AiListKbMemberLogic) AiListKbMember(req *types.AiListKbMemberReq) (resp *types.AiListKbMemberResp, err error) {
	// 1. 提取当前用户ID，校验登录状态
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 2. 校验必填参数
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}

	// 3. 构建角色过滤条件
	role := strings.TrimSpace(req.Role)

	// 4. 调用 RPC 查询成员列表
	list, err := l.svcCtx.AiKnowledgeClient.ListKbMember(l.ctx, &aiknowledgeclient.ListKbMemberReq{
		Page:    int64(req.Page),
		PageSize: int64(req.PageSize),
		KbId:    req.KbId,
		UserId:  userID,
		Keyword: strings.TrimSpace(req.Keyword),
		Role:    role,
		HasRole: role != "",
	})
	if err != nil {
		return nil, err
	}

	// 5. 映射 RPC 响应为网关类型
	items := make([]types.AiKbMemberItem, 0, len(list.List))
	for _, item := range list.List {
		items = append(items, kbMemberItemFromRPC(item))
	}

	return &types.AiListKbMemberResp{
		Total: list.Total,
		List:  items,
	}, nil
}
