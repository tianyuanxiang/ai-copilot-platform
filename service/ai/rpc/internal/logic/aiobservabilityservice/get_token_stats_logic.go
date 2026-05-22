package aiobservabilityservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenStatsLogic {
	return &GetTokenStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询当前用户在指定范围内的 token 用量统计。
func (l *GetTokenStatsLogic) GetTokenStats(in *pb.GetTokenStatsReq) (*pb.TokenStatsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.TokenStatsResp{}, nil
}
