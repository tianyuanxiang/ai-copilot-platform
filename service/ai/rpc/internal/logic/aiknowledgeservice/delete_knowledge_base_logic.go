package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeBaseLogic {
	return &DeleteKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除知识库，并由实现层清理关联文档和索引。
func (l *DeleteKnowledgeBaseLogic) DeleteKnowledgeBase(in *pb.DeleteKnowledgeBaseReq) (*pb.Empty, error) {
	if in.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}

	// 检查维护权限
	canMaintain, err := canMaintainKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId)
	if err != nil {
		return nil, err
	}
	if !canMaintain {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有删除该知识库的权限")
	}

	// 删除知识库（数据库级联删除会处理关联的 documents、chunks、members）
	if err := l.svcCtx.AiKnowledgeBaseModel.Delete(l.ctx, in.KbId); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
