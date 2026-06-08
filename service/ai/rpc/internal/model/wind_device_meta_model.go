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
		DisplayMetaForDeviceType(ctx context.Context, deviceTypeCode string) (WindDeviceDisplayMeta, error)
	}

	customWindDeviceMetaModel struct {
		*defaultWindDeviceMetaModel
		db *gorm.DB
	}

	WindDeviceDisplayMeta struct {
		DeviceTypeCode string            `json:"deviceTypeCode"`
		DeviceTypeName string            `json:"deviceTypeName"`
		Fields         []string          `json:"fields"`
		FieldLabels    map[string]string `json:"fieldLabels"`
		FieldUnits     map[string]string `json:"fieldUnits"`
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
	displayMeta, err := m.DisplayMetaForDeviceType(ctx, deviceTypeCode)
	if err != nil {
		return nil, err
	}
	return displayMeta.ResolveFields(requestedFields)
}

func (m *customWindDeviceMetaModel) DisplayMetaForDeviceType(ctx context.Context, deviceTypeCode string) (WindDeviceDisplayMeta, error) {
	code := strings.ToUpper(strings.TrimSpace(deviceTypeCode))
	displayMeta := newWindDeviceDisplayMeta(code)
	if m.db != nil && code != "" {
		var rows []struct {
			DeviceTypeName string `gorm:"column:device_type_name"`
			ColumnName     string `gorm:"column:column_name"`
			DisplayName    string `gorm:"column:display_name"`
			Unit           string `gorm:"column:unit"`
		}
		err := m.db.WithContext(ctx).Table("wind_device_type AS dt").
			Select("dt.device_type_name, wm.column_name, wm.display_name, wm.unit").
			Joins("LEFT JOIN wind_device_meta AS wm ON wm.device_type_code = dt.device_type_code AND wm.ai_enabled = ?", true).
			Where("dt.device_type_code = ? AND dt.is_delete = ?", code, 0).
			Order("wm.ord asc").
			Scan(&rows).Error
		if err == nil && len(rows) > 0 {
			for _, row := range rows {
				if strings.TrimSpace(row.DeviceTypeName) != "" {
					displayMeta.DeviceTypeName = strings.TrimSpace(row.DeviceTypeName)
				}
				if IsSafeIdentifier(row.ColumnName) {
					displayMeta.addField(row.ColumnName, row.DisplayName, row.Unit)
				}
			}
		}
	}
	displayMeta.fillFallbackFields()
	return displayMeta, nil
}

func (m WindDeviceDisplayMeta) ResolveFields(requestedFields []string) ([]string, error) {
	fields := normalizeRequestedFields(requestedFields)
	if len(fields) == 0 {
		return append([]string(nil), m.Fields...), nil
	}
	allowedSet := make(map[string]struct{}, len(m.Fields))
	for _, item := range m.Fields {
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

func (m WindDeviceDisplayMeta) FieldLabel(field string) string {
	if label := strings.TrimSpace(m.FieldLabels[field]); label != "" {
		return label
	}
	return field
}

func (m WindDeviceDisplayMeta) FieldUnit(field string) string {
	return strings.TrimSpace(m.FieldUnits[field])
}

func newWindDeviceDisplayMeta(code string) WindDeviceDisplayMeta {
	return WindDeviceDisplayMeta{
		DeviceTypeCode: code,
		DeviceTypeName: code,
		FieldLabels:    make(map[string]string),
		FieldUnits:     make(map[string]string),
	}
}

func (m *WindDeviceDisplayMeta) addField(columnName, displayName, unit string) {
	columnName = strings.TrimSpace(columnName)
	if columnName == "" {
		return
	}
	m.Fields = append(m.Fields, columnName)
	if label := strings.TrimSpace(displayName); label != "" {
		m.FieldLabels[columnName] = label
	} else {
		m.FieldLabels[columnName] = columnName
	}
	if unit = strings.TrimSpace(unit); unit != "" {
		m.FieldUnits[columnName] = unit
	}
}

func (m *WindDeviceDisplayMeta) fillFallbackFields() {
	if len(m.Fields) > 0 {
		for _, field := range m.Fields {
			if _, ok := m.FieldLabels[field]; !ok {
				m.FieldLabels[field] = field
			}
		}
		return
	}
	stable := DeviceTypeStableFallback[m.DeviceTypeCode]
	for _, field := range stableFieldsFallback[stable] {
		m.addField(field, field, "")
	}
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
