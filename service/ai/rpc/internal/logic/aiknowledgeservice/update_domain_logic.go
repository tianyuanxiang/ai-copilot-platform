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

type UpdateDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDomainLogic {
	return &UpdateDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新公共知识库领域基础信息。
func (l *UpdateDomainLogic) UpdateDomain(in *pb.UpdateDomainReq) (*pb.Empty, error) {
	if in.DomainId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "domainId 不能为空")
	}

	existing, err := l.svcCtx.AiKbDomainModel.FindOne(l.ctx, in.DomainId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "领域不存在")
		}
		return nil, err
	}

	// 应用部分更新
	if in.HasName {
		existing.Name = strings.TrimSpace(in.Name)
	}
	if in.HasCode {
		code := strings.TrimSpace(in.Code)
		if code != existing.Code {
			_, err := l.svcCtx.AiKbDomainModel.FindOneByCode(l.ctx, code)
			if err == nil {
				return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "领域编码已存在")
			}
			if err != model.ErrNotFound {
				return nil, err
			}
		}
		existing.Code = code
	}
	if in.HasDescription {
		existing.Description = strings.TrimSpace(in.Description)
	}
	if in.HasSort {
		existing.Sort = in.Sort
	}
	if in.HasStatus {
		existing.Status = in.Status
	}

	if err := l.svcCtx.AiKbDomainModel.Update(l.ctx, existing); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
