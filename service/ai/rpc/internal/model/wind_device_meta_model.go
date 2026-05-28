package model

import (
	"context"
	"errors"
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
		FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedFields []string) ([]string, error)
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

func (m *customWindDeviceMetaModel) FieldsForDeviceType(ctx context.Context, deviceTypeCode string, requestedFields []string) ([]string, error) {
	code := strings.ToUpper(strings.TrimSpace(deviceTypeCode))

	allowed := make([]string, 0)
	if m.db != nil && code != "" {
		var rows []struct {
			ColumnName  string `gorm:"column:column_name"`
			DisplayName string `gorm:"column:display_name"`
		}
		err := m.db.WithContext(ctx).Table("wind_device_meta").
			Where("device_type_code = ?", code).
			Where("ai_enabled = ?", true).
			Order("ord asc").
			Find(&rows).Error
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

	fields := normalizeRequestedFields(requestedFields)
	if len(fields) == 0 {
		return allowed, nil
	}
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, item := range allowed {
		allowedSet[item] = struct{}{}
	}
	selected := make([]string, 0, len(fields))
	for _, field := range fields {
		if _, ok := allowedSet[field]; !ok {
			return nil, errors.New("field not found")
		}
		selected = append(selected, field)
	}
	return selected, nil
}

func normalizeRequestedFields(fields []string) []string {
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		for _, item := range strings.Split(field, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				result = append(result, item)
			}
		}
	}
	return result
}
