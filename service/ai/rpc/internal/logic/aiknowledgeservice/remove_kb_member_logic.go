package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
