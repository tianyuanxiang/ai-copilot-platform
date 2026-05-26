package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDomainLogic {
	return &DeleteDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除公共知识库领域。
func (l *DeleteDomainLogic) DeleteDomain(in *pb.DeleteDomainReq) (*pb.Empty, error) {
	if in.DomainId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "domainId 不能为空")
	}

	_, err := l.svcCtx.AiKbDomainModel.FindOne(l.ctx, in.DomainId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "领域不存在")
		}
		return nil, err
	}

	if err := l.svcCtx.AiKbDomainModel.Delete(l.ctx, in.DomainId); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
