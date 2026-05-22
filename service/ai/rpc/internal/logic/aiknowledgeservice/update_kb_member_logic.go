package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateKbMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateKbMemberLogic {
	return &UpdateKbMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新公共知识库成员角色。
func (l *UpdateKbMemberLogic) UpdateKbMember(in *pb.UpdateKbMemberReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
