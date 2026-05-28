package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ WindDeviceTypeModel = (*customWindDeviceTypeModel)(nil)

type (
	// WindDeviceTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindDeviceTypeModel.
	WindDeviceTypeModel interface {
		windDeviceTypeModel
		withSession(session sqlx.Session) WindDeviceTypeModel
		FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedField string) []string
		AlarmDeviceType(deviceTypeCode string) int64
	}

	customWindDeviceTypeModel struct {
		*defaultWindDeviceTypeModel
		db *gorm.DB
	}
)

// NewWindDeviceTypeModel returns a model for the database table.
func NewWindDeviceTypeModel(conn sqlx.SqlConn, db *gorm.DB) WindDeviceTypeModel {
	return &customWindDeviceTypeModel{
		defaultWindDeviceTypeModel: newWindDeviceTypeModel(conn),
		db:                         db,
	}
}

func (m *customWindDeviceTypeModel) withSession(session sqlx.Session) WindDeviceTypeModel {
	return NewWindDeviceTypeModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customWindDeviceTypeModel) FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedField string) []string {
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
		stable := DeviceTypeStableFallback[code]
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

func (m *customWindDeviceTypeModel) AlarmDeviceType(deviceTypeCode string) int64 {
	return alarmDeviceTypeFallback[strings.ToUpper(strings.TrimSpace(deviceTypeCode))]
}
