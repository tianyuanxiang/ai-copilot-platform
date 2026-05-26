package aiknowledgeservicelogic

import (
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListKbMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKbMemberLogic {
	return &ListKbMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询公共知识库成员列表。
func (l *ListKbMemberLogic) ListKbMember(in *pb.ListKbMemberReq) (*pb.ListKbMemberResp, error) {
	if in.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 检查知识库是否存在及用户是否有访问权限
	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}
	if !canAccessKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有访问该知识库的权限")
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	members, total, err := l.svcCtx.AiKbMemberModel.ListByKbID(l.ctx, in.KbId, page, pageSize, in.Keyword, in.Role, in.HasRole)
	if err != nil {
		return nil, err
	}

	// 批量获取用户信息
	items := make([]*pb.KbMemberItem, 0, len(members))
	for _, m := range members {
		item := &pb.KbMemberItem{
			KbId:      m.KbId,
			UserId:    m.UserId,
			Role:      m.Role,
			CreatedAt: formatDocumentTime(m.CreatedAt),
		}

		// 通过 SysRpc 获取用户名和昵称
		userResp, err := l.svcCtx.PermRpc.GetUserById(l.ctx, &permclient.GetUserByIdReq{UserId: m.UserId})
		if err == nil && userResp != nil {
			item.Username = userResp.Username
			item.Nickname = userResp.Nickname
		}

		items = append(items, item)
	}

	// 过滤关键字匹配（如果 keyword 非空且数据库层无法完全过滤）
	if kw := strings.TrimSpace(in.Keyword); kw != "" {
		filtered := make([]*pb.KbMemberItem, 0)
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Username), strings.ToLower(kw)) ||
				strings.Contains(strings.ToLower(item.Nickname), strings.ToLower(kw)) {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	return &pb.ListKbMemberResp{
		Total: total,
		List:  items,
	}, nil
}
