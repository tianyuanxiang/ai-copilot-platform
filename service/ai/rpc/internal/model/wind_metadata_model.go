package model

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

var farmDatabaseFallback = map[string]string{
	"FY": "fuyu",
	"YS": "yushu",
}

var deviceTypeStableFallback = map[string]string{
	"STM":  "strain",
	"ACC":  "accel",
	"ACCX": "accel",
	"ACCY": "accel",
	"INS":  "inclinometer",
	"INSX": "inclinometer",
	"INSY": "inclinometer",
	"ATS":  "tension",
	"WPR":  "radar",
	"JMT":  "joint_meter",
	"HLS":  "hydrostatic",
	"ULS":  "ultrasonic_level",
	"GNSS": "gnss",
}

var stableFieldsFallback = map[string][]string{
	"strain":           {"strain"},
	"accel":            {"accel"},
	"inclinometer":     {"x", "y"},
	"tension":          {"tension"},
	"radar":            {"d", "rws", "veer", "raws", "ti", "hw_shub", "direction_hub", "v_sheer", "h_sheer", "hw_shigh", "direction_high", "hw_slow", "direction_low"},
	"joint_meter":      {"joint"},
	"hydrostatic":      {"settlement", "temperature", "pressure"},
	"ultrasonic_level": {"height"},
	"gnss":             {"longitude", "latitude", "vertical", "horizontal", "ordinate", "vertical_offset", "horizontal_offset", "ordinate_offset"},
}

var alarmDeviceTypeFallback = map[string]int64{
	"STM":  1,
	"ACC":  2,
	"ACCX": 2,
	"ACCY": 2,
	"INS":  3,
	"INSX": 3,
	"INSY": 3,
	"ATS":  4,
	"WPR":  5,
	"JMT":  6,
	"HLS":  7,
	"ULS":  8,
	"GNSS": 9,
	"IPC":  10,
}

type (
	WindMetadataModel interface {
		ListFarms(ctx context.Context, keyword string) ([]*WindFarm, error)
		ListTurbines(ctx context.Context, farmCode string, keyword string) ([]*WindTurbine, error)
		ListDevices(ctx context.Context, farmCode string, towerCode string, deviceTypeCode string, keyword string) ([]*WindDeviceView, error)
		FarmDatabase(ctx context.Context, farmCode string) string
		StableForDeviceType(ctx context.Context, deviceTypeCode string) string
		FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedField string) []string
		AlarmDeviceType(deviceTypeCode string) int64
	}

	windMetadataModel struct {
		db *gorm.DB
	}

	WindFarm struct {
		FarmId     int64  `gorm:"column:farm_id"`
		FarmCode   string `gorm:"column:farm_code"`
		FarmName   string `gorm:"column:farm_name"`
		Province   string `gorm:"column:province"`
		Location   string `gorm:"column:location"`
		TDDatabase string `gorm:"column:td_database"`
		AiEnabled  bool   `gorm:"column:ai_enabled"`
	}

	WindTurbine struct {
		TurbineId int64  `gorm:"column:turbine_id"`
		TowerId   int64  `gorm:"column:tower_id"`
		TowerCode string `gorm:"column:tower_code"`
		TowerName string `gorm:"column:tower_name"`
		FarmCode  string `gorm:"column:farm_code"`
		FarmName  string `gorm:"column:farm_name"`
		RiskLevel string `gorm:"column:risk_level"`
		AiEnabled bool   `gorm:"column:ai_enabled"`
	}

	WindDeviceView struct {
		DeviceId       int64  `gorm:"column:device_id"`
		DeviceCode     string `gorm:"column:device_code"`
		DeviceTypeCode string `gorm:"column:device_type_code"`
		DeviceTypeName string `gorm:"column:device_type_name"`
		TowerId        int64  `gorm:"column:tower_id"`
		TowerCode      string `gorm:"column:tower_code"`
		StructureCode  string `gorm:"column:structure_code"`
		StructureName  string `gorm:"column:structure_name"`
		TDStable       string `gorm:"column:td_stable"`
		Status         int64  `gorm:"column:status"`
		AiEnabled      bool   `gorm:"column:ai_enabled"`
	}
)

func NewWindMetadataModel(db *gorm.DB) WindMetadataModel {
	return &windMetadataModel{db: db}
}

func (m *windMetadataModel) ListFarms(ctx context.Context, keyword string) ([]*WindFarm, error) {
	var rows []*WindFarm
	query := m.db.WithContext(ctx).Table("wind_farm").Where("is_delete = 0")
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("farm_code like ? or farm_name like ?", like, like)
	}
	err := query.Order("farm_id asc").Find(&rows).Error
	return rows, err
}

func (m *windMetadataModel) ListTurbines(ctx context.Context, farmCode string, keyword string) ([]*WindTurbine, error) {
	var rows []*WindTurbine
	query := m.db.WithContext(ctx).Table("wind_turbine").Where("is_delete = 0")
	if farmCode = strings.TrimSpace(farmCode); farmCode != "" {
		query = query.Where("farm_code = ?", strings.ToUpper(farmCode))
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("tower_code like ? or farm_name like ?", like, like)
	}
	err := query.Select("tower_id as turbine_id, tower_id, tower_code, tower_code as tower_name, farm_code, farm_name, risk_level, ai_enabled").
		Order("farm_code asc, tower_code asc").
		Find(&rows).Error
	return rows, err
}

func (m *windMetadataModel) ListDevices(ctx context.Context, farmCode string, towerCode string, deviceTypeCode string, keyword string) ([]*WindDeviceView, error) {
	var rows []*WindDeviceView
	query := m.db.WithContext(ctx).Table("wind_device d").
		Select(`d.device_id, d.device_code, d.device_type_code, dt.device_type_name,
			d.tower_id, t.tower_code, d.structure_code, st.structure_name,
			coalesce(nullif(d.td_stable, ''), dt.td_stable) as td_stable,
			d.status, d.ai_enabled`).
		Joins("left join wind_turbine t on t.tower_id = d.tower_id").
		Joins("left join wind_device_type dt on dt.device_type_code = d.device_type_code").
		Joins("left join wind_structure_type st on st.structure_code = d.structure_code").
		Where("d.is_delete = 0")
	if farmCode = strings.TrimSpace(farmCode); farmCode != "" {
		query = query.Where("t.farm_code = ?", strings.ToUpper(farmCode))
	}
	if towerCode = strings.TrimSpace(towerCode); towerCode != "" {
		query = query.Where("t.tower_code = ?", towerCode)
	}
	if deviceTypeCode = strings.TrimSpace(deviceTypeCode); deviceTypeCode != "" {
		query = query.Where("d.device_type_code = ?", strings.ToUpper(deviceTypeCode))
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("d.device_code like ? or dt.device_type_name like ?", like, like)
	}
	err := query.Order("t.farm_code asc, t.tower_code asc, d.device_type_code asc, d.device_code asc").Find(&rows).Error
	return rows, err
}

func (m *windMetadataModel) FarmDatabase(ctx context.Context, farmCode string) string {
	code := strings.ToUpper(strings.TrimSpace(farmCode))
	if m.db != nil && code != "" {
		var row struct {
			TDDatabase string `gorm:"column:td_database"`
		}
		err := m.db.WithContext(ctx).Raw(`select td_database from wind_farm where farm_code = ? and is_delete = 0 limit 1`, code).Scan(&row).Error
		if err == nil && row.TDDatabase != "" {
			return row.TDDatabase
		}
	}
	if v := farmDatabaseFallback[code]; v != "" {
		return v
	}
	return strings.ToLower(code)
}

func (m *windMetadataModel) StableForDeviceType(ctx context.Context, deviceTypeCode string) string {
	code := strings.ToUpper(strings.TrimSpace(deviceTypeCode))
	if m.db != nil && code != "" {
		var row struct {
			TDStable string `gorm:"column:td_stable"`
		}
		err := m.db.WithContext(ctx).Raw(`select td_stable from wind_device_type where device_type_code = ? and is_delete = 0 limit 1`, code).Scan(&row).Error
		if err == nil && row.TDStable != "" {
			return row.TDStable
		}
	}
	return deviceTypeStableFallback[code]
}

func (m *windMetadataModel) FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedField string) []string {
	code := strings.ToUpper(strings.TrimSpace(deviceTypeCode))
	field := strings.TrimSpace(requestedField)
	allowed := make([]string, 0)
	if m.db != nil && code != "" {
		var rows []struct {
			ColumnName string `gorm:"column:column_name"`
		}
		err := m.db.WithContext(ctx).Raw(`select column_name from wind_device_meta where device_type_code = ? and ai_enabled = true order by ord asc`, code).Scan(&rows).Error
		if err == nil && len(rows) > 0 {
			for _, row := range rows {
				if IsSafeIdentifier(row.ColumnName) {
					allowed = append(allowed, row.ColumnName)
				}
			}
		}
	}
	if len(allowed) == 0 {
		stable := deviceTypeStableFallback[code]
		allowed = stableFieldsFallback[stable]
	}
	if field == "" {
		return allowed
	}
	for _, item := range allowed {
		if field == item {
			return []string{field}
		}
	}
	return nil
}

func (m *windMetadataModel) AlarmDeviceType(deviceTypeCode string) int64 {
	return alarmDeviceTypeFallback[strings.ToUpper(strings.TrimSpace(deviceTypeCode))]
}
