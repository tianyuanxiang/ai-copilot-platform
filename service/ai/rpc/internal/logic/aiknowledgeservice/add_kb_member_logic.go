package aiknowledgeservicelogic

import (
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddKbMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddKbMemberLogic {
	return &AddKbMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加公共知识库成员并授予角色。
func (l *AddKbMemberLogic) AddKbMember(in *pb.AddKbMemberReq) (*pb.Empty, error) {
	if in.KbId <= 0 || in.UserId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 和 userId 不能为空")
	}
	role := strings.TrimSpace(in.Role)
	if role == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "role 不能为空")
	}
	if role != "viewer" && role != "editor" && role != "manager" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "role 只支持 viewer、editor、manager")
	}

	// 检查知识库是否存在
	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}

	// 操作者必须有 manager 权限
	canMaintain, err := canMaintainKnowledgeBase(l.ctx, l.svcCtx, kb, in.OperatorId)
	if err != nil {
		return nil, err
	}
	if !canMaintain {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "只有 manager 才能添加成员")
	}

	// 检查用户是否已经是成员
	_, err = l.svcCtx.AiKbMemberModel.FindByKbIDUserID(l.ctx, in.KbId, in.UserId)
	if err == nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "该用户已是知识库成员")
	}
	if err != model.ErrNotFound {
		return nil, err
	}

	_, err = l.svcCtx.AiKbMemberModel.Insert(l.ctx, &model.AiKbMember{
		KbId:      in.KbId,
		UserId:    in.UserId,
		Role:      role,
		CreatedBy: in.OperatorId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
