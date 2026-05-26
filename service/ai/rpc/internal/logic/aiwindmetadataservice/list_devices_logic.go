package aiwindmetadataservicelogic

import (
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDevicesLogic {
	return &ListDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDevicesLogic) ListDevices(in *pb.ListWindDeviceReq) (*pb.ListWindDeviceResp, error) {
	var rows []struct {
		DeviceId       int64  `gorm:"column:device_id"`
		DeviceCode     string `gorm:"column:device_code"`
		DeviceTypeCode string `gorm:"column:device_type_code"`
		DeviceTypeName string `gorm:"column:device_type_name"`
		TowerId        int64  `gorm:"column:tower_id"`
		TowerCode      string `gorm:"column:tower_code"`
		StructureCode  string `gorm:"column:structure_code"`
		StructureName  string `gorm:"column:structure_name"`
		Status         int64  `gorm:"column:status"`
	}

	query := l.svcCtx.Orm.Table("wind_device wd").
		Select("wd.device_id, wd.device_code, wd.device_type_code, wd.device_type_name, wd.tower_id, wt.tower_code, wd.structure_code, wd.structure_name, wd.status").
		Joins("left join wind_tower wt on wt.tower_id = wd.tower_id").
		Joins("left join wind_device_type wdt on wdt.device_type_id = wd.device_type_id").
		Where("wd.is_delete = 0")

	if farmCode := strings.TrimSpace(in.FarmCode); farmCode != "" {
		query = query.Where("wt.farm_code = ?", strings.ToUpper(farmCode))
	}
	if towerCode := strings.TrimSpace(in.TowerCode); towerCode != "" {
		query = query.Where("wt.tower_code = ?", towerCode)
	}
	if deviceTypeCode := strings.TrimSpace(in.DeviceTypeCode); deviceTypeCode != "" {
		query = query.Where("wd.device_type_code = ?", strings.ToUpper(deviceTypeCode))
	}
	if keyword := strings.TrimSpace(in.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("wd.device_code like ? or wd.device_type_name like ?", like, like)
	}
	if err := query.Order("wt.farm_code asc, wt.tower_code asc, wd.device_type_code asc, wd.device_code asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]*pb.WindDeviceItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &pb.WindDeviceItem{
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
	return &pb.ListWindDeviceResp{Total: int64(len(items)), List: items, Message: "wind device metadata scaffold ready"}, nil
}
