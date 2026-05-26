package aiwindalarmservicelogic

import (
	"context"
	"fmt"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAlarmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryAlarmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAlarmsLogic {
	return &QueryAlarmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryAlarmsLogic) QueryAlarms(in *pb.WindAlarmQueryReq) (*pb.WindAlarmQueryResp, error) {
	database := l.svcCtx.WindFarmModel.FarmDatabase(l.ctx, in.FarmCode)
	whereParts := model.TimeWhere(in.StartTime, in.EndTime)
	whereParts = append(whereParts, model.DeviceWhere(in.TowerCode, in.DeviceCode)...)
	if strings.TrimSpace(in.AlarmCode) != "" {
		whereParts = append(whereParts, fmt.Sprintf("alarm_code=%d", parseInt(in.AlarmCode)))
	}
	if in.AlarmLevel > 0 {
		whereParts = append(whereParts, fmt.Sprintf("alarm_level=%d", in.AlarmLevel))
	}
	if n := l.svcCtx.WindDeviceTypeModel.AlarmDeviceType(in.DeviceTypeCode); n > 0 {
		whereParts = append(whereParts, fmt.Sprintf("device_type=%d", n))
	}
	whereParts = append(whereParts, fmt.Sprintf("status=%d", in.Status))
	where := model.JoinWhere(whereParts)

	total, err := l.svcCtx.TdengineModel.QueryCount(l.ctx, database, "alarm", where)
	if err != nil {
		return nil, err
	}
	rows, err := l.svcCtx.TdengineModel.QueryAlarms(l.ctx, database, where, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]*pb.WindAlarmItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &pb.WindAlarmItem{
			Ts:            row["ts"],
			FarmCode:      in.FarmCode,
			TowerCode:     row["tower_id"],
			DeviceChannel: parseInt(row["device_channel"]),
			DeviceType:    parseInt(row["device_type"]),
			AlarmLocation: row["alarm_location"],
			AlarmLevel:    parseInt(row["alarm_level"]),
			AlarmCode:     parseInt(row["alarm_code"]),
			AlarmValue:    row["alarm_value"],
			Status:        parseInt(row["status"]),
		})
	}
	evidence := map[string]any{
		"source":           "tdengine.alarm",
		"scaffold":         !l.svcCtx.TdengineModel.IsConfigured(),
		"farm_code":        in.FarmCode,
		"database":         database,
		"where":            where,
		"returned_records": len(items),
	}
	message := "TDengine alarm query scaffold ready"
	if !l.svcCtx.TdengineModel.IsConfigured() {
		message = "TDengine is not configured; returning alarm scaffold evidence only"
	}
	return &pb.WindAlarmQueryResp{Total: total, List: items, EvidenceJson: model.WindEvidenceJSON(evidence), Message: message}, nil
}

func parseInt(value string) int64 {
	var n int64
	fmt.Sscan(value, &n)
	return n
}
