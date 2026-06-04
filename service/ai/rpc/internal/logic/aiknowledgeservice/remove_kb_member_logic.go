package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveKbMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveKbMemberLogic {
	return &RemoveKbMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 移除公共知识库成员。
func (l *RemoveKbMemberLogic) RemoveKbMember(in *pb.RemoveKbMemberReq) (*pb.Empty, error) {
	if in.KbId <= 0 || in.UserId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 和 userId 不能为空")
	}

	// 检查知识库是否存在
	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}

	// 操作者必须是知识库 owner 或 manager。
	canManage, err := canManageKnowledgeBaseMembers(l.ctx, l.svcCtx, kb, in.OperatorId)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "只有 manager 才能移除成员")
	}

	// 查找成员记录
	member, err := l.svcCtx.AiKbMemberModel.FindOneByKbIdUserId(l.ctx, in.KbId, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "该用户不是知识库成员")
		}
		return nil, err
	}

	if err := l.svcCtx.AiKbMemberModel.Delete(l.ctx, member.Id); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
