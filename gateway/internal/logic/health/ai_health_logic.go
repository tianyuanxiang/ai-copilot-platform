// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package health

import (
	"context"

	"ai-copilot-platform/ai-rpc/client/aistatusservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiHealthLogic {
	return &AiHealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiHealthLogic) AiHealth() (resp *types.AiHealthResp, err error) {
	rpcResp, err := l.svcCtx.AiStatusClient.Health(l.ctx, &aistatusservice.HealthReq{})
	if err != nil {
		l.Logger.Errorf("health rpc call failed: %v", err)
		return nil, err
	}

	return &types.AiHealthResp{
		Status:      rpcResp.Status,
		EngineReady: rpcResp.EngineReady,
		Message:     rpcResp.Message,
	}, nil
}
