package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListKbMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKbMemberLogic {
	return &ListKbMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询公共知识库成员列表。
func (l *ListKbMemberLogic) ListKbMember(in *pb.ListKbMemberReq) (*pb.ListKbMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListKbMemberResp{}, nil
}
