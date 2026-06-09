package model

import (
	"context"
	"database/sql"
	"fmt"
	"math"
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
		// QueryRadarRows 查询 WPR 雷达时序数据，支持按 index_id 过滤距离层。
		// indexId=0 表示不过滤，返回所有距离层；indexId>0 只返回指定距离层。
		// 结果按 ts DESC, index_id ASC 排序。
		QueryRadarRows(ctx context.Context, database string, stable string, fields []string, where string, indexId int64, page int64, pageSize int64) ([]map[string]string, []string, error)
		// QueryAggStats 针对普通传感器按字段返回聚合统计（min/max/avg/count/first/last）。
		// 返回 map[fieldName]AggStat。
		QueryAggStats(ctx context.Context, database string, stable string, fields []string, where string) (map[string]AggStat, error)
		// QueryRadarAggStats 针对 WPR 雷达按 index_id 分层返回聚合统计。
		// 返回 map[indexId]map[fieldName]AggStat。
		QueryRadarAggStats(ctx context.Context, database string, stable string, fields []string, where string) (map[int]map[string]AggStat, error)
		IsConfigured() bool

		// QueryBasicStats 趋势专用：查询各字段的 COUNT/MIN/MAX/AVG。
		QueryBasicStats(ctx context.Context, database string, stable string, fields []string, where string) (map[string]BasicStat, error)
		// QueryBoundaryRows 趋势专用：查询各字段的首行和尾行数据。
		QueryBoundaryRows(ctx context.Context, database string, stable string, fields []string, where string) (first map[string]string, last map[string]string, err error)
		// QueryMinuteBuckets 趋势专用：分钟级降采样查询，返回各字段在各时间窗口的聚合。
		QueryMinuteBuckets(ctx context.Context, database string, stable string, fields []string, where string) ([]BucketRow, error)

		// QueryAlarmAggregates 一次性查询当前条件下的告警总数、首末时间戳。
		QueryAlarmAggregates(ctx context.Context, database string, where string) (map[string]string, error)
		// QueryAlarmGroupBy 按指定列分组统计告警数量。
		QueryAlarmGroupBy(ctx context.Context, database string, groupByColumn string, where string, limit int64) ([]map[string]string, error)
		// QueryAlarmTimeBuckets 按时间桶统计告警数量。
		QueryAlarmTimeBuckets(ctx context.Context, database string, where string, interval string) ([]map[string]string, error)
		// 通过测量距离查询雷达index_id
		ResolveRadarIndexByDistance(ctx context.Context, database string, stable string, where string, distanceM float64) (int64, float64, error)
	}

	tdengineModel struct {
		db *sql.DB
	}

	// BasicStat 趋势专用：单个字段的基础聚合统计结果。
	BasicStat struct {
		Count sql.NullInt64
		Min   sql.NullFloat64
		Max   sql.NullFloat64
		Avg   sql.NullFloat64
	}

	// BucketRow 趋势专用：分钟级降采样查询的一行结果。
	BucketRow struct {
		Field  string
		Wstart string
		Avg    sql.NullFloat64
		Min    sql.NullFloat64
		Max    sql.NullFloat64
		Count  sql.NullInt64
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
func JoinWhere(parts []string) string {
	if len(parts) == 0 {
		return "1=1"
	}
	return strings.Join(parts, " and ")
}

func TimeWhere(startTime, endTime string) []string {
	var where []string
	startTime = strings.TrimSpace(startTime)
	endTime = strings.TrimSpace(endTime)
	if startTime != "" && endTime != "" {
		where = append(where, fmt.Sprintf("ts >= '%s' AND ts < '%s'", escapeSQL(startTime), escapeSQL(endTime)))
	}
	return where
}

func DeviceWhere(towerCode, deviceCode string) []string {
	var where []string
	towerCode = strings.TrimSpace(towerCode)
	if towerCode != "" {
		where = append(where, fmt.Sprintf("tower_id=%d", parseTdNumber(towerCode)))
	}
	deviceCode = strings.TrimSpace(deviceCode)
	if deviceCode != "" {
		where = append(where, fmt.Sprintf("device_channel=%d", parseTdNumber(deviceCode)))
	}
	return where
}

func escapeSQL(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func parseTdNumber(value string) int64 {
	n, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return n
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
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)
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
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)
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
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maps, _, err := scanTDengineRows(rows)
	return maps, err
}

// QueryRadarRows 查询 WPR 雷达时序数据，支持按 index_id 过滤距离层。
// indexId=0 表示不过滤，返回所有距离层；indexId>0 只返回指定距离层。
func (m *tdengineModel) QueryRadarRows(ctx context.Context, database string, stable string, fields []string, where string, indexId int64, page int64, pageSize int64) ([]map[string]string, []string, error) {
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
	// 雷达字段列表包含 index_id
	radarFields := make([]string, 0, len(fields)+2)
	radarFields = append(radarFields, "ts", "index_id")
	for _, f := range fields {
		if f == "ts" || f == "index_id" {
			continue
		}
		sf, err := SafeIdentifier(f)
		if err != nil {
			return nil, nil, err
		}
		radarFields = append(radarFields, sf)
	}
	// 追加 index_id 过滤条件
	fullWhere := where
	if indexId > 0 {
		fullWhere = fmt.Sprintf("(%s) AND index_id=%d", where, indexId)
	}
	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s ORDER BY ts DESC, index_id ASC %s",
		strings.Join(radarFields, ","), database, stable, fullWhere, normalizeTDLimit(page, pageSize))
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	return scanTDengineRows(rows)
}

// QueryAggStats 针对普通传感器按字段返回聚合统计（所有字段一次查询）。
func (m *tdengineModel) QueryAggStats(ctx context.Context, database string, stable string, fields []string, where string) (map[string]AggStat, error) {
	if !m.IsConfigured() {
		return nil, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, nil
	}

	// 1. 动态构建 SELECT 列：COUNT(f),MIN(f),MAX(f),AVG(f),STDDEV(f),FIRST(f),LAST(f) x len(fields)
	allCols := make([]string, 0, len(fields)*7+2)
	for _, f := range fields {
		sf, err := SafeIdentifier(f)
		if err != nil {
			continue
		}
		allCols = append(allCols,
			fmt.Sprintf("COUNT(%s)", sf),
			fmt.Sprintf("MIN(%s)", sf),
			fmt.Sprintf("MAX(%s)", sf),
			fmt.Sprintf("AVG(%s)", sf),
			fmt.Sprintf("STDDEV(%s)", sf),
			fmt.Sprintf("FIRST(%s)", sf),
			fmt.Sprintf("LAST(%s)", sf))
	}
	allCols = append(allCols, "FIRST(ts)", "LAST(ts)")

	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s",
		strings.Join(allCols, ","), database, stable, where)
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)

	// 2. 动态构建 scan 目标（每个字段 7 个聚合 + 2 个 ts）
	targets := make([]any, 0, len(fields)*7+2)
	aggPtrs := make([]*AggStat, len(fields))
	for i := range fields {
		a := &AggStat{}
		aggPtrs[i] = a
		targets = append(targets, &a.Count, &a.Min, &a.Max, &a.Avg, &a.Stddev, &a.First, &a.Last)
	}
	var firstTs, lastTs sql.NullString
	targets = append(targets, &firstTs, &lastTs)

	// 3. 执行查询
	err = m.db.QueryRowContext(ctx, query).Scan(targets...)
	if err != nil {
		return nil, err
	}

	// 4. 组装结果
	result := make(map[string]AggStat, len(fields))
	for i, field := range fields {
		a := aggPtrs[i]
		if a == nil {
			continue
		}
		if firstTs.Valid {
			a.FirstTs = firstTs.String
		}
		if lastTs.Valid {
			a.LastTs = lastTs.String
		}
		result[field] = *a
	}
	return result, nil
}

// QueryRadarAggStats 针对 WPR 雷达按 index_id 分层返回聚合统计（所有字段一次查询）。
func (m *tdengineModel) QueryRadarAggStats(ctx context.Context, database string, stable string, fields []string, where string) (map[int]map[string]AggStat, error) {
	if !m.IsConfigured() {
		return nil, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, nil
	}

	// 1. 动态构建 SELECT 列（不含 index_id，它作为第一列单独加）
	allCols := make([]string, 0, len(fields)*7+2)
	for _, f := range fields {
		sf, err := SafeIdentifier(f)
		if err != nil {
			continue
		}
		allCols = append(allCols,
			fmt.Sprintf("COUNT(%s)", sf),
			fmt.Sprintf("MIN(%s)", sf),
			fmt.Sprintf("MAX(%s)", sf),
			fmt.Sprintf("AVG(%s)", sf),
			fmt.Sprintf("STDDEV(%s)", sf),
			fmt.Sprintf("FIRST(%s)", sf),
			fmt.Sprintf("LAST(%s)", sf))
	}
	allCols = append(allCols, "FIRST(ts)", "LAST(ts)")

	query := fmt.Sprintf("SELECT index_id,%s FROM %s.%s WHERE %s PARTITION BY index_id",
		strings.Join(allCols, ","), database, stable, where)
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)

	// 2. 执行查询（多行，每行一个 index_id）
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]map[string]AggStat)
	for rows.Next() {
		var indexId int
		targets := make([]any, 0, len(fields)*7+3) // indexId + field aggs + 2 ts
		targets = append(targets, &indexId)

		perIndexAggs := make([]*AggStat, len(fields))
		for i := range fields {
			a := &AggStat{}
			perIndexAggs[i] = a
			targets = append(targets, &a.Count, &a.Min, &a.Max, &a.Avg, &a.Stddev, &a.First, &a.Last)
		}
		var firstTs, lastTs sql.NullString
		targets = append(targets, &firstTs, &lastTs)

		if scanErr := rows.Scan(targets...); scanErr != nil {
			continue
		}

		inner := make(map[string]AggStat, len(fields))
		for i, field := range fields {
			a := perIndexAggs[i]
			if a == nil {
				continue
			}
			if firstTs.Valid {
				a.FirstTs = firstTs.String
			}
			if lastTs.Valid {
				a.LastTs = lastTs.String
			}
			inner[field] = *a
		}
		result[indexId] = inner
	}
	return result, rows.Err()
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

func normalizeTDLimit(page int64, pageSize int64) string {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
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

// QueryBasicStats 趋势专用：查询各字段的 COUNT/MIN/MAX/AVG。
// 每个字段单独查询，以确保 SafeIdentifier 校验覆盖到每个字段。
func (m *tdengineModel) QueryBasicStats(ctx context.Context, database, stable string, fields []string, where string) (map[string]BasicStat, error) {
	if !m.IsConfigured() {
		return nil, fmt.Errorf("tdengine is not configured")
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return nil, nil
	}

	result := make(map[string]BasicStat, len(fields))
	for _, f := range fields {
		sf, err := SafeIdentifier(f)
		if err != nil {
			continue
		}
		query := fmt.Sprintf("SELECT COUNT(%s), MIN(%s), MAX(%s), AVG(%s) FROM %s.%s WHERE %s",
			sf, sf, sf, sf, database, stable, where)
		fmt.Printf("[TDengine] [%s]: %s\n", database, query)

		var stat BasicStat
		err = m.db.QueryRowContext(ctx, query).Scan(&stat.Count, &stat.Min, &stat.Max, &stat.Avg)
		if err != nil {
			return result, err
		}

		result[f] = stat
	}
	return result, nil
}

// QueryBoundaryRows 趋势专用：查询各字段的首行和尾行数据。
// 返回 first 和 last 两个 map[string]string，分别对应首行和尾行的所有字段值。
func (m *tdengineModel) QueryBoundaryRows(ctx context.Context, database, stable string, fields []string, where string) (map[string]string, map[string]string, error) {
	if !m.IsConfigured() {
		return nil, nil, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, nil, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return nil, nil, err
	}

	safeFields, err := safeBoundarySelect(fields)
	if err != nil {
		return nil, nil, err
	}
	cols := strings.Join(safeFields, ",")

	// 首行（按时间升序第一条）
	firstQuery := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s ORDER BY ts ASC LIMIT 1",
		cols, database, stable, where)
	fmt.Printf("[TDengine] [%s]: %s\n", database, firstQuery)
	first, err := m.queryOneRow(ctx, firstQuery, safeFields)
	if err != nil {
		return nil, nil, err
	}

	// 尾行（按时间降序第一条）
	lastQuery := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s ORDER BY ts DESC LIMIT 1",
		cols, database, stable, where)
	fmt.Printf("[TDengine] [%s]: %s\n", database, lastQuery)
	last, err := m.queryOneRow(ctx, lastQuery, safeFields)
	if err != nil {
		return nil, nil, err
	}

	return first, last, nil
}

// queryOneRow 执行单行查询，返回 map[columnName]value。
func (m *tdengineModel) queryOneRow(ctx context.Context, query string, columns []string) (map[string]string, error) {
	row := m.db.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return nil, row.Err()
	}

	targets := make([]any, len(columns))
	ptrs := make([]any, len(columns))
	for i := range targets {
		ptrs[i] = &targets[i]
	}
	if err := row.Scan(ptrs...); err != nil {
		return nil, err
	}

	result := make(map[string]string, len(columns))
	for i, col := range columns {
		result[col] = formatTDValue(targets[i])
	}
	return result, nil
}

// safeBoundarySelect 构建边界查询的 SELECT 列（包含 ts 和所有字段），返回安全的列名列表。
func safeBoundarySelect(fields []string) ([]string, error) {
	safeFields := make([]string, 0, len(fields)+1)
	safeFields = append(safeFields, "ts")
	for _, field := range fields {
		sf, err := SafeIdentifier(field)
		if err != nil {
			return nil, err
		}
		safeFields = append(safeFields, sf)
	}
	return safeFields, nil
}

// QueryMinuteBuckets 趋势专用：分钟级降采样查询。
// 每个字段单独执行 INTERVAL(1m) 查询，结果合并到统一的 BucketRow 切片中。
func (m *tdengineModel) QueryMinuteBuckets(ctx context.Context, database, stable string, fields []string, where string) ([]BucketRow, error) {
	if !m.IsConfigured() {
		return nil, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, nil
	}

	var allRows []BucketRow
	for _, f := range fields {
		sf, err := SafeIdentifier(f)
		if err != nil {
			continue
		}
		query := fmt.Sprintf("SELECT _wstart, AVG(%s), MIN(%s), MAX(%s), COUNT(%s) FROM %s.%s WHERE %s INTERVAL(1h)",
			sf, sf, sf, sf, database, stable, where)
		fmt.Printf("[TDengine] [%s]: %s\n", database, query)

		rows, err := m.db.QueryContext(ctx, query)
		if err != nil {
			return allRows, err
		}
		defer rows.Close()

		for rows.Next() {
			var br BucketRow
			br.Field = f
			if scanErr := rows.Scan(&br.Wstart, &br.Avg, &br.Min, &br.Max, &br.Count); scanErr != nil {
				continue
			}
			allRows = append(allRows, br)
		}
		if closeErr := rows.Close(); closeErr != nil {
			return allRows, closeErr
		}
	}
	return allRows, nil
}

var alarmGroupByAllowlist = map[string]bool{
	"alarm_level": true,
	"status":      true,
	"alarm_code":  true,
	"tower_id":    true,
}

var alarmIntervalPattern = regexp.MustCompile(`^\d+[smhdw]$`)

// QueryAlarmAggregates 一次性查询当前条件下的告警总数、首末时间戳。
func (m *tdengineModel) QueryAlarmAggregates(ctx context.Context, database string, where string) (map[string]string, error) {
	if !m.IsConfigured() {
		return nil, nil
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT COUNT(*),FIRST(ts),LAST(ts) FROM %s.alarm WHERE %s",
		database, where)
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)

	var total sql.NullString
	var firstTs, lastTs sql.NullString
	if err := m.db.QueryRowContext(ctx, query).Scan(&total, &firstTs, &lastTs); err != nil {
		return nil, err
	}
	return map[string]string{
		"total":    nullStr(total),
		"first_ts": nullStr(firstTs),
		"last_ts":  nullStr(lastTs),
	}, nil
}

// QueryAlarmGroupBy 按指定列分组统计告警数量。
func (m *tdengineModel) QueryAlarmGroupBy(ctx context.Context, database string, groupByColumn string, where string, limit int64) ([]map[string]string, error) {
	if !m.IsConfigured() {
		return nil, nil
	}
	if !alarmGroupByAllowlist[groupByColumn] {
		return nil, fmt.Errorf("alarm group by column not allowed: %s", groupByColumn)
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT %s,COUNT(*) AS cnt FROM %s.alarm WHERE %s GROUP BY %s ORDER BY cnt DESC",
		groupByColumn, database, where, groupByColumn)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maps, _, err := scanTDengineRows(rows)
	return maps, err
}

// QueryAlarmTimeBuckets 按时间桶统计告警数量。
func (m *tdengineModel) QueryAlarmTimeBuckets(ctx context.Context, database string, where string, interval string) ([]map[string]string, error) {
	if !m.IsConfigured() {
		return nil, fmt.Errorf("tdengine not configured")
	}
	if !alarmIntervalPattern.MatchString(interval) {
		return nil, fmt.Errorf("invalid alarm interval: %s", interval)
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT _wstart,COUNT(*) AS cnt FROM %s.alarm WHERE %s INTERVAL(%s)",
		database, where, interval)
	fmt.Printf("[TDengine] [%s]: %s\n", database, query)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maps, _, err := scanTDengineRows(rows)
	return maps, err
}

func (m *tdengineModel) ResolveRadarIndexByDistance(ctx context.Context, database string, stable string, where string, distanceM float64) (int64, float64, error) {
	if !m.IsConfigured() {
		return 0, 0, fmt.Errorf("tdengine not configured")
	}
	database, err := SafeIdentifier(database)
	if err != nil {
		return 0, 0, err
	}
	stable, err = SafeIdentifier(stable)
	if err != nil {
		return 0, 0, err
	}

	query := fmt.Sprintf(
		"SELECT index_id, LAST(d) FROM %s.%s WHERE %s PARTITION BY index_id",
		database, stable, where,
	)

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var bestIndex int64
	var bestDistance float64
	bestDiff := math.MaxFloat64

	for rows.Next() {
		var indexID int64
		var d sql.NullFloat64
		if err := rows.Scan(&indexID, &d); err != nil {
			return 0, 0, err
		}
		if !d.Valid {
			continue
		}
		diff := math.Abs(d.Float64 - distanceM)
		if diff < bestDiff {
			bestDiff = diff
			bestIndex = indexID
			bestDistance = d.Float64
		}
	}

	return bestIndex, bestDistance, rows.Err()
}

func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
