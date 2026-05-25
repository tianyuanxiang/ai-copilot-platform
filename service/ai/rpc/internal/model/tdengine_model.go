package model

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var safeIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type (
	TdengineModel interface {
		QueryCount(ctx context.Context, database string, stable string, where string) (int64, error)
		QueryRows(ctx context.Context, database string, stable string, fields []string, where string, page int64, pageSize int64) ([]map[string]string, []string, error)
		QueryLatestRow(ctx context.Context, database string, stable string, fields []string, where string) (map[string]string, []string, error)
		QueryAlarms(ctx context.Context, database string, where string, page int64, pageSize int64) ([]map[string]string, error)
		IsConfigured() bool
	}

	tdengineModel struct {
		db *sql.DB
	}
)

func NewTdengineModel(db *sql.DB) TdengineModel {
	return &tdengineModel{db: db}
}

func IsSafeIdentifier(value string) bool {
	return safeIdentifierPattern.MatchString(strings.TrimSpace(value))
}

func SafeIdentifier(value string) (string, error) {
	value = strings.TrimSpace(value)
	if !IsSafeIdentifier(value) {
		return "", fmt.Errorf("invalid identifier: %s", value)
	}
	return value, nil
}

func TimeWhere(startTime, endTime string) []string {
	var where []string
	if strings.TrimSpace(startTime) != "" && strings.TrimSpace(endTime) != "" {
		where = append(where, fmt.Sprintf("ts >= '%s' AND ts < '%s'", escapeSQLLiteral(startTime), escapeSQLLiteral(endTime)))
	}
	return where
}

func DeviceWhere(towerCode, deviceCode string) []string {
	var where []string
	if towerCode = strings.TrimSpace(towerCode); towerCode != "" {
		where = append(where, fmt.Sprintf("tower_id=%d", parseNumber(towerCode)))
	}
	if deviceCode = strings.TrimSpace(deviceCode); deviceCode != "" {
		where = append(where, fmt.Sprintf("device_channel=%d", parseNumber(deviceCode)))
	}
	return where
}

func JoinWhere(parts []string) string {
	if len(parts) == 0 {
		return "1=1"
	}
	return strings.Join(parts, " and ")
}

func (m *tdengineModel) IsConfigured() bool {
	return m != nil && m.db != nil
}

func (m *tdengineModel) QueryCount(ctx context.Context, database string, stable string, where string) (int64, error) {
	if !m.IsConfigured() {
		return 0, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return 0, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return 0, err
	}
	var total int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE %s", database, stable, where)
	if err := m.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (m *tdengineModel) QueryRows(ctx context.Context, database string, stable string, fields []string, where string, page int64, pageSize int64) ([]map[string]string, []string, error) {
	if !m.IsConfigured() {
		return nil, fields, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, nil, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return nil, nil, err
	}
	safeFields, err := safeSelectFields(fields)
	if err != nil {
		return nil, nil, err
	}
	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s ORDER BY ts DESC %s", strings.Join(safeFields, ","), database, stable, where, normalizeTDLimit(page, pageSize))
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	return scanTDengineRows(rows)
}

func (m *tdengineModel) QueryLatestRow(ctx context.Context, database string, stable string, fields []string, where string) (map[string]string, []string, error) {
	rows, columns, err := m.QueryRows(ctx, database, stable, fields, where, 1, 1)
	if err != nil || len(rows) == 0 {
		return nil, columns, err
	}
	return rows[0], columns, nil
}

func (m *tdengineModel) QueryAlarms(ctx context.Context, database string, where string, page int64, pageSize int64) ([]map[string]string, error) {
	if !m.IsConfigured() {
		return nil, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT ts,tower_id,device_channel,device_type,alarm_location,alarm_level,alarm_code,alarm_value,status FROM %s.alarm WHERE %s ORDER BY ts DESC %s", database, where, normalizeTDLimit(page, pageSize))
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maps, _, err := scanTDengineRows(rows)
	return maps, err
}

func safeSelectFields(fields []string) ([]string, error) {
	safeFields := make([]string, 0, len(fields)+1)
	safeFields = append(safeFields, "ts")
	for _, field := range fields {
		field, err := SafeIdentifier(field)
		if err != nil {
			return nil, err
		}
		safeFields = append(safeFields, field)
	}
	return safeFields, nil
}

func escapeSQLLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func parseNumber(value string) int64 {
	n, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return n
}

func normalizeTDLimit(page int64, pageSize int64) string {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)
}

func scanTDengineRows(rows *sql.Rows) ([]map[string]string, []string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	out := make([]map[string]string, 0)
	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		item := make(map[string]string, len(cols))
		for i, col := range cols {
			item[col] = formatTDValue(values[i])
		}
		out = append(out, item)
	}
	return out, cols, rows.Err()
}

func formatTDValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05.000")
	default:
		return fmt.Sprint(v)
	}
}
