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

type CreateDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDomainLogic {
	return &CreateDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建公共知识库领域。
func (l *CreateDomainLogic) CreateDomain(in *pb.CreateDomainReq) (*pb.CreateDomainResp, error) {
	name := strings.TrimSpace(in.Name)
	code := strings.TrimSpace(in.Code)
	if name == "" || code == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "name 和 code 不能为空")
	}

	// 检查 code 唯一性
	_, err := l.svcCtx.AiKbDomainModel.FindOneByCode(l.ctx, code)
	if err == nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "领域编码已存在")
	}
	if err != model.ErrNotFound {
		return nil, err
	}

	sort := in.Sort
	if sort <= 0 {
		sort = 0
	}
	status := in.Status
	if status <= 0 {
		status = 1
	}

	result, err := l.svcCtx.AiKbDomainModel.Insert(l.ctx, &model.AiKbDomain{
		Name:        name,
		Code:        code,
		Description: strings.TrimSpace(in.Description),
		Sort:        sort,
		Status:      status,
		CreatedBy:   in.OperatorId,
	})
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &pb.CreateDomainResp{DomainId: id}, nil
}
