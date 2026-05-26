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

type CreateKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateKnowledgeBaseLogic {
	return &CreateKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建个人或公共知识库。
func (l *CreateKnowledgeBaseLogic) CreateKnowledgeBase(in *pb.CreateKnowledgeBaseReq) (*pb.CreateKnowledgeBaseResp, error) {
	kbType := strings.ToLower(strings.TrimSpace(in.KbType))
	name := strings.TrimSpace(in.Name)
	if kbType == "" || name == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbType 和 name 不能为空")
	}
	if kbType != "personal" && kbType != "public" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbType 只支持 personal 或 public")
	}

	kb := &model.AiKnowledgeBase{
		KbType:      kbType,
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		CreatedBy:   in.OperatorId,
	}

	// 根据类型设置不同字段
	switch kbType {
	case "personal":
		kb.OwnerUserId = sql.NullInt64{Int64: in.UserId, Valid: in.UserId > 0}
		kb.Visibility = "private"
	case "public":
		if in.DomainId <= 0 {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "公共知识库 domainId 不能为空")
		}
		// 检查领域是否存在
		_, err := l.svcCtx.AiKbDomainModel.FindOne(l.ctx, in.DomainId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "领域不存在")
			}
			return nil, err
		}
		kb.DomainId = sql.NullInt64{Int64: in.DomainId, Valid: true}
		visibility := strings.TrimSpace(in.Visibility)
		if visibility == "" {
			visibility = "private"
		}
		kb.Visibility = visibility
	}

	status := in.Status
	if status <= 0 {
		status = 1
	}
	kb.Status = status

	result, err := l.svcCtx.AiKnowledgeBaseModel.Insert(l.ctx, kb)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	// 如果是公共知识库，创建者自动成为 manager
	if kbType == "public" && in.UserId > 0 {
		_, _ = l.svcCtx.AiKbMemberModel.Insert(l.ctx, &model.AiKbMember{
			KbId:      id,
			UserId:    in.UserId,
			Role:      "manager",
			CreatedBy: in.OperatorId,
		})
	}

	return &pb.CreateKnowledgeBaseResp{KbId: id}, nil
}
