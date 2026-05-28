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

type AiWindQueryTimeseriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindQueryTimeseriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindQueryTimeseriesLogic {
	return &AiWindQueryTimeseriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindQueryTimeseriesLogic) AiWindQueryTimeseries(req *types.AiWindTimeseriesQueryReq) (resp *types.AiWindTimeseriesQueryResp, err error) {
	result, err := l.svcCtx.AiWindTimeseriesClient.QueryTimeseries(l.ctx, &pb.WindTimeseriesQueryReq{
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceCode:     req.DeviceCode,
		DeviceTypeCode: req.DeviceTypeCode,
		Field:          req.Field,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Page:           int64(req.Page),
		PageSize:       int64(req.PageSize),
		IndexId:        req.IndexId,
	})
	if err != nil {
		return nil, err
	}
	points := make([]types.AiWindDataPoint, 0, len(result.Points))
	for _, point := range result.Points {
		if point == nil {
			continue
		}
		points = append(points, types.AiWindDataPoint{Ts: point.Ts, Values: point.Values})
	}
	return &types.AiWindTimeseriesQueryResp{
		Total:        result.Total,
		Database:     result.Database,
		Stable:       result.Stable,
		Fields:       result.Fields,
		Points:       points,
		EvidenceJson: result.EvidenceJson,
		Message:      result.Message,
	}, nil
}
