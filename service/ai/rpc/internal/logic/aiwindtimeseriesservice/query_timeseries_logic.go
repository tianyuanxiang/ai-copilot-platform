package aiwindtimeseriesservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"

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

// QueryTimeseries 查询风机时序数据。
// WPR 雷达支持按 index_id 过滤距离层，indexId=0 返回全部距离层。
// 其他设备类型忽略 indexId 字段。
func (l *QueryTimeseriesLogic) QueryTimeseries(in *pb.WindTimeseriesQueryReq) (*pb.WindTimeseriesQueryResp, error) {
	database := l.svcCtx.WindFarmModel.FarmDatabase(l.ctx, in.FarmCode)

	if in.DeviceTypeCode == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "DeviceTypeCode不能为空")
	}
	stableName := model.DeviceTypeStableFallback[in.DeviceTypeCode]

	fields, err := l.svcCtx.WindDeviceMetaModel.FieldsForDeviceType(l.ctx, in.DeviceTypeCode, in.Field)
	if err != nil {
		l.Logger.Errorf("FieldsForDeviceType failed %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "通过设备类型获取测点名称失败")
	}

	whereParts := append(model.TimeWhere(in.StartTime, in.EndTime), model.DeviceWhere(in.TowerCode, in.DeviceCode)...)
	where := model.JoinWhere(whereParts)

	isRadar := in.DeviceTypeCode == "WPR"

	// indexId 只对 WPR 雷达生效；非雷达设备传了 indexId 时忽略并在 evidence 标记
	indexIdIgnored := !isRadar && in.IndexId != 0

	var total int64
	var rows []map[string]string

	if isRadar {
		// 雷达查询：含 index_id 字段，支持按距离层过滤
		total, err = l.svcCtx.TdengineModel.QueryCount(l.ctx, database, stableName, where)
		if err != nil {
			return nil, err
		}
		rows, _, err = l.svcCtx.TdengineModel.QueryRadarRows(l.ctx, database, stableName, fields, where, in.IndexId, in.Page, in.PageSize)
		if err != nil {
			return nil, err
		}
	} else {
		// 普通传感器查询
		total, err = l.svcCtx.TdengineModel.QueryCount(l.ctx, database, stableName, where)
		if err != nil {
			return nil, err
		}
		rows, _, err = l.svcCtx.TdengineModel.QueryRows(l.ctx, database, stableName, fields, where, in.Page, in.PageSize)
		if err != nil {
			return nil, err
		}
	}

	points := make([]*pb.WindDataPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, &pb.WindDataPoint{Ts: row["ts"], Values: row})
	}

	// 构造 evidence
	evidence := map[string]any{
		"source":          "tdengine",
		"scaffold":        !l.svcCtx.TdengineModel.IsConfigured(),
		"farmCode":        in.FarmCode,
		"towerCode":       in.TowerCode,
		"deviceTypeCode":  in.DeviceTypeCode,
		"database":        database,
		"stable":          stableName,
		"fields":          fields,
		"where":           where,
		"returnedRecords": len(points),
		"isRadar":         isRadar,
	}
	if isRadar {
		// 记录雷达 index 信息
		evidence["indexId"] = in.IndexId
		if in.IndexId == 0 {
			// 返回全部 10 个距离层
			allIndexIds := make([]int, 10)
			for i := range allIndexIds {
				allIndexIds[i] = i + 1
			}
			evidence["selectedIndexIds"] = allIndexIds
			evidence["indexCount"] = 10
		} else {
			evidence["selectedIndexIds"] = []int64{in.IndexId}
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
