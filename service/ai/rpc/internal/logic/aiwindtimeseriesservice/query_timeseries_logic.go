package aiwindtimeseriesservicelogic

import (
	"context"
	"fmt"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"go-zero-rpc/common/xerr"

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
	deviceTypeCode := strings.ToUpper(strings.TrimSpace(in.DeviceTypeCode))
	if deviceTypeCode == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "DeviceTypeCode不能为空")
	}

	database := l.svcCtx.WindFarmModel.FarmDatabase(l.ctx, in.FarmCode)
	stableName := model.DeviceTypeStableFallback[deviceTypeCode]
	if stableName == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("不支持的设备类型: %s", deviceTypeCode))
	}

	displayMeta, err := l.svcCtx.WindDeviceMetaModel.DisplayMetaForDeviceType(l.ctx, deviceTypeCode)
	if err != nil {
		l.Logger.Errorf("DisplayMetaForDeviceType failed %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "通过设备类型获取测点元数据失败")
	}
	fields, err := displayMeta.ResolveFields(in.Field)
	if err != nil {
		l.Logger.Errorf("ResolveFields failed %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "设备类型不存在该测点字段")
	}

	whereParts := append(model.TimeWhere(in.StartTime, in.EndTime), model.DeviceWhere(in.TowerCode, in.DeviceCode)...)
	where := model.JoinWhere(whereParts)
	isRadar := deviceTypeCode == wprDeviceTypeCode

	resolvedIndexID := in.IndexId
	var matchedDistanceM float64
	if isRadar && resolvedIndexID == 0 && in.RadarDistanceM > 0 && l.svcCtx.TdengineModel.IsConfigured() {
		resolvedIndexID, matchedDistanceM, err = l.svcCtx.TdengineModel.ResolveRadarIndexByDistance(
			l.ctx, database, stableName, where, float64(in.RadarDistanceM),
		)
		if err != nil {
			l.Logger.Errorf("Get radar indexId failed %v", err)
			return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "通过距离获取雷达 index_id 失败")
		}
	}
	indexIdIgnored := !isRadar && in.IndexId != 0

	var total int64
	var rows []map[string]string
	if isRadar {
		countWhere := where
		if resolvedIndexID > 0 {
			countWhere = fmt.Sprintf("(%s) AND index_id=%d", where, resolvedIndexID)
		}
		total, err = l.svcCtx.TdengineModel.QueryCount(l.ctx, database, stableName, countWhere)
		if err != nil {
			l.Logger.Errorf("Query radar data count failed %v", err)
			return nil, err
		}
		rows, _, err = l.svcCtx.TdengineModel.QueryRadarRows(l.ctx, database, stableName, fields, where, resolvedIndexID, in.Page, in.PageSize)
		if err != nil {
			l.Logger.Errorf("Query radar data failed %v", err)
			return nil, err
		}
	} else {
		total, err = l.svcCtx.TdengineModel.QueryCount(l.ctx, database, stableName, where)
		if err != nil {
			l.Logger.Errorf("Query data count failed %v", err)
			return nil, err
		}
		rows, _, err = l.svcCtx.TdengineModel.QueryRows(l.ctx, database, stableName, fields, where, in.Page, in.PageSize)
		if err != nil {
			l.Logger.Errorf("Query data rows failed %v", err)
			return nil, err
		}
	}

	points := make([]*pb.WindDataPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, &pb.WindDataPoint{Ts: row["ts"], Values: row})
	}

	evidence := map[string]any{
		"source":          "tdengine",
		"scaffold":        !l.svcCtx.TdengineModel.IsConfigured(),
		"farmCode":        in.FarmCode,
		"towerCode":       in.TowerCode,
		"deviceTypeCode":  deviceTypeCode,
		"deviceTypeName":  displayMeta.DeviceTypeName,
		"fieldLabels":     displayMeta.FieldLabels,
		"fieldUnits":      displayMeta.FieldUnits,
		"database":        database,
		"stable":          stableName,
		"fields":          fields,
		"where":           where,
		"returnedRecords": len(points),
		"isRadar":         isRadar,
	}
	if isRadar {
		evidence["indexId"] = resolvedIndexID
		evidence["requestedDistanceM"] = in.RadarDistanceM
		evidence["matchedDistanceM"] = matchedDistanceM
		if resolvedIndexID == 0 {
			allIndexIds := make([]int, 10)
			for i := range allIndexIds {
				allIndexIds[i] = i + 1
			}
			evidence["selectedIndexIds"] = allIndexIds
			evidence["indexCount"] = 10
		} else {
			evidence["selectedIndexIds"] = []int64{resolvedIndexID}
			evidence["indexCount"] = 1
		}
	}
	if indexIdIgnored {
		evidence["indexIdIgnored"] = true
	}

	message := "TDengine timeseries query completed"
	if !l.svcCtx.TdengineModel.IsConfigured() {
		message = "TDengine is not configured; returning scaffold evidence only"
	}

	return &pb.WindTimeseriesQueryResp{
		Total:        total,
		Database:     database,
		Stable:       stableName,
		Fields:       fields,
		Points:       points,
		EvidenceJson: model.WindEvidenceJSON(evidence),
		Message:      message,
	}, nil
}
