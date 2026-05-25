// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_alarm

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindQueryAlarmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindQueryAlarmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindQueryAlarmsLogic {
	return &AiWindQueryAlarmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindQueryAlarmsLogic) AiWindQueryAlarms(req *types.AiWindAlarmQueryReq) (resp *types.AiWindAlarmQueryResp, err error) {
	result, err := l.svcCtx.AiWindAlarmClient.QueryAlarms(l.ctx, &pb.WindAlarmQueryReq{
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceTypeCode: req.DeviceTypeCode,
		DeviceCode:     req.DeviceCode,
		AlarmCode:      req.AlarmCode,
		AlarmLevel:     int64(req.AlarmLevel),
		Status:         int64(req.Status),
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Page:           int64(req.Page),
		PageSize:       int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.AiWindAlarmItem, 0, len(result.List))
	for _, item := range result.List {
		if item == nil {
			continue
		}
		list = append(list, types.AiWindAlarmItem{
			Ts:            item.Ts,
			FarmCode:      item.FarmCode,
			TowerCode:     item.TowerCode,
			DeviceChannel: item.DeviceChannel,
			DeviceType:    item.DeviceType,
			AlarmLocation: item.AlarmLocation,
			AlarmLevel:    item.AlarmLevel,
			AlarmCode:     item.AlarmCode,
			AlarmValue:    item.AlarmValue,
			Status:        item.Status,
		})
	}
	return &types.AiWindAlarmQueryResp{Total: result.Total, List: list, EvidenceJson: result.EvidenceJson, Message: result.Message}, nil
}
