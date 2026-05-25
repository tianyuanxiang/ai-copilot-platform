package aiwindtimeseriesservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryTimeseriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryTimeseriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTimeseriesLogic {
	return &QueryTimeseriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryTimeseriesLogic) QueryTimeseries(in *pb.WindTimeseriesQueryReq) (*pb.WindTimeseriesQueryResp, error) {
	database := l.svcCtx.WindMetadataModel.FarmDatabase(l.ctx, in.FarmCode)
	stable := l.svcCtx.WindMetadataModel.StableForDeviceType(l.ctx, in.DeviceTypeCode)
	fields := l.svcCtx.WindMetadataModel.FieldsForDeviceType(l.ctx, in.DeviceTypeCode, in.Field)
	whereParts := append(model.TimeWhere(in.StartTime, in.EndTime), model.DeviceWhere(in.TowerCode, in.DeviceCode)...)
	where := model.JoinWhere(whereParts)

	total, err := l.svcCtx.TdengineModel.QueryCount(l.ctx, database, stable, where)
	if err != nil {
		return nil, err
	}
	rows, _, err := l.svcCtx.TdengineModel.QueryRows(l.ctx, database, stable, fields, where, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}
	points := make([]*pb.WindDataPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, &pb.WindDataPoint{Ts: row["ts"], Values: row})
	}
	evidence := map[string]any{
		"source":           "tdengine",
		"scaffold":         !l.svcCtx.TdengineModel.IsConfigured(),
		"farm_code":        in.FarmCode,
		"database":         database,
		"stable":           stable,
		"fields":           fields,
		"where":            where,
		"returned_records": len(points),
	}
	message := "TDengine timeseries query scaffold ready"
	if !l.svcCtx.TdengineModel.IsConfigured() {
		message = "TDengine is not configured; returning scaffold evidence only"
	}
	return &pb.WindTimeseriesQueryResp{
		Total:        total,
		Database:     database,
		Stable:       stable,
		Fields:       fields,
		Points:       points,
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      message,
	}, nil
}
