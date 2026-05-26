package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ WindDeviceModel = (*customWindDeviceModel)(nil)

type (
	// WindDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindDeviceModel.
	WindDeviceModel interface {
		windDeviceModel
		withSession(session sqlx.Session) WindDeviceModel
		ListDevices(ctx context.Context, farmCode string, towerCode string, deviceTypeCode string, keyword string) ([]*WindDeviceView, error)
	}

	customWindDeviceModel struct {
		*defaultWindDeviceModel
		db *gorm.DB
	}
)

// NewWindDeviceModel returns a model for the database table.
func NewWindDeviceModel(conn sqlx.SqlConn, db *gorm.DB) WindDeviceModel {
	return &customWindDeviceModel{
		defaultWindDeviceModel: newWindDeviceModel(conn),
		db:                     db,
	}
}

func (m *customWindDeviceModel) withSession(session sqlx.Session) WindDeviceModel {
	return NewWindDeviceModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customWindDeviceModel) ListDevices(ctx context.Context, farmCode string, towerCode string, deviceTypeCode string, keyword string) ([]*WindDeviceView, error) {
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
