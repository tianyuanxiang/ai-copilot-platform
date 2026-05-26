package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ WindDeviceMetaModel = (*customWindDeviceMetaModel)(nil)

type (
	// WindDeviceMetaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindDeviceMetaModel.
	WindDeviceMetaModel interface {
		windDeviceMetaModel
		withSession(session sqlx.Session) WindDeviceMetaModel
		FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedField string) []string
	}

	customWindDeviceMetaModel struct {
		*defaultWindDeviceMetaModel
		db *gorm.DB
	}
)

// NewWindDeviceMetaModel returns a model for the database table.
func NewWindDeviceMetaModel(conn sqlx.SqlConn, db *gorm.DB) WindDeviceMetaModel {
	return &customWindDeviceMetaModel{
		defaultWindDeviceMetaModel: newWindDeviceMetaModel(conn),
		db:                         db,
	}
}

func (m *customWindDeviceMetaModel) withSession(session sqlx.Session) WindDeviceMetaModel {
	return NewWindDeviceMetaModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customWindDeviceMetaModel) FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedField string) []string {
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
