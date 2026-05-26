package aiknowledgeservicelogic

import (
	"context"
	"database/sql"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateKnowledgeBaseLogic {
	return &UpdateKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新知识库名称、描述、领域、可见性或状态。
func (l *UpdateKnowledgeBaseLogic) UpdateKnowledgeBase(in *pb.UpdateKnowledgeBaseReq) (*pb.Empty, error) {
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
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有修改该知识库的权限")
	}

	// 应用部分更新
	if in.HasName {
		kb.Name = strings.TrimSpace(in.Name)
	}
	if in.HasDescription {
		kb.Description = strings.TrimSpace(in.Description)
	}
	if in.HasDomainId {
		// 检查领域是否存在
		if in.DomainId > 0 {
			_, err := l.svcCtx.AiKbDomainModel.FindOne(l.ctx, in.DomainId)
			if err != nil {
				if err == model.ErrNotFound {
					return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "领域不存在")
				}
				return nil, err
			}
		}
		kb.DomainId = sql.NullInt64{Int64: in.DomainId, Valid: in.DomainId > 0}
	}
	if in.HasVisibility {
		kb.Visibility = strings.TrimSpace(in.Visibility)
	}
	if in.HasStatus {
		kb.Status = in.Status
	}

	if err := l.svcCtx.AiKnowledgeBaseModel.Update(l.ctx, kb); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
