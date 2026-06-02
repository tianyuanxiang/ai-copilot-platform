// Package aiwindtool 的 metadata 文件实现风场、风机和设备元数据查询工具。
package aiwindtool

import (
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/pb"
)

// executeGetTurbineMetadata 直接复用 WindFarmModel、WindTowerModel 和 WindDeviceModel。
// 这里不调用现有 ListTurbinesLogic，因为该 logic 会为 towerCode 拼接展示文本；
// Agent 需要原始编码，才能继续调用时序和告警工具。
func (e *Executor) executeGetTurbineMetadata(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args GetTurbineMetadataArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		e.Logger.Errorf("decode args err:%v", err)
		return nil, err
	}

	farmCode := normalizeFarmCode(args.FarmCode)
	towerCode := strings.TrimSpace(args.TowerCode)
	deviceTypeCode := strings.ToUpper(strings.TrimSpace(args.DeviceTypeCode))

	farms := make([]*pb.WindFarmItem, 0)
	turbines := make([]*pb.WindTurbineItem, 0)
	devices := make([]*pb.WindDeviceItem, 0)

	if farmCode == "" {
		rows, err := e.svcCtx.WindFarmModel.ListFarms(ctx, "")
		if err != nil {
			e.Logger.Errorf("list farms err:%v", err)
			return nil, err
		}
		for _, row := range rows {
			farms = append(farms, &pb.WindFarmItem{
				FarmId:     row.FarmId,
				FarmCode:   row.FarmCode,
				FarmName:   row.FarmName,
				Province:   row.Province,
				Location:   row.Location,
				TdDatabase: row.TdDatabase,
				AiEnabled:  row.AiEnabled,
			})
		}
	} else {
		rows, err := e.svcCtx.WindTowerModel.ListTurbines(ctx, farmCode, "")
		if err != nil {
			e.Logger.Errorf("list tower err:%v", err)
			return nil, err
		}
		for _, row := range rows {
			turbines = append(turbines, &pb.WindTurbineItem{
				TowerId:   row.TowerId,
				TowerCode: row.TowerCode,
				FarmId:    row.FarmId,
				FarmCode:  row.FarmCode,
				FarmName:  row.FarmName,
				RiskLevel: row.RiskLevel,
				AiEnabled: row.AiEnabled,
				Remark:    row.Remark,
			})
		}
	}

	if farmCode != "" && towerCode != "" {
		rows, err := e.svcCtx.WindDeviceModel.ListDevices(ctx, farmCode, towerCode, deviceTypeCode, "")
		if err != nil {
			e.Logger.Errorf("list devices err:%v", err)
			return nil, err
		}
		for _, row := range rows {
			devices = append(devices, &pb.WindDeviceItem{
				DeviceId:       row.DeviceId,
				DeviceCode:     row.DeviceCode,
				DeviceTypeCode: row.DeviceTypeCode,
				DeviceTypeName: row.DeviceTypeName,
				TowerId:        row.TowerId,
				TowerCode:      row.TowerCode,
				StructureCode:  row.StructureCode,
				StructureName:  row.StructureName,
				Status:         row.Status,
			})
		}
	}

	resultJSON := marshalJSON(map[string]any{
		"farms":    farms,
		"turbines": turbines,
		"devices":  devices,
	})
	return &toolResult{
		ResultJSON: resultJSON,
		Message:    "wind turbine metadata query completed",
	}, nil
}
