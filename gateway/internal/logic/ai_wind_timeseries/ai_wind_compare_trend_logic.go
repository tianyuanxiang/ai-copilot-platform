// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_timeseries

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindCompareTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindCompareTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindCompareTrendLogic {
	return &AiWindCompareTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindCompareTrendLogic) AiWindCompareTrend(req *types.AiWindTrendCompareReq) (resp *types.AiWindTrendCompareResp, err error) {
	result, err := l.svcCtx.AiWindTimeseriesClient.CompareTrend(l.ctx, &pb.WindTrendCompareReq{
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceTypeCode: req.DeviceTypeCode,
		Field:          req.Field,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &types.AiWindTrendCompareResp{Summary: result.Summary, EvidenceJson: result.EvidenceJson, Message: result.Message}, nil
}
