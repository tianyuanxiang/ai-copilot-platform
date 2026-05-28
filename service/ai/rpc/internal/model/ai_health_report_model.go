package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	AiHealthReportModel interface {
		InsertReturningID(ctx context.Context, data *AiHealthReport) (int64, error)
		FindOne(ctx context.Context, id int64) (*AiHealthReport, error)
		withSession(session sqlx.Session) AiHealthReportModel
	}

	customAiHealthReportModel struct {
		conn  sqlx.SqlConn
		table string
	}

	AiHealthReport struct {
		Id         int64        `db:"id"`
		UserId     int64        `db:"user_id"`
		TraceId    string       `db:"trace_id"`
		ReportType string       `db:"report_type"`
		FarmCode   string       `db:"farm_code"`
		TowerCode  string       `db:"tower_code"`
		StartTime  sql.NullTime `db:"start_time"`
		EndTime    sql.NullTime `db:"end_time"`
		Title      string       `db:"title"`
		Content    string       `db:"content"`
		Evidence   string       `db:"evidence"`
		Status     string       `db:"status"`
		CreatedAt  time.Time    `db:"created_at"`
	}
)

func NewAiHealthReportModel(conn sqlx.SqlConn) AiHealthReportModel {
	return &customAiHealthReportModel{
		conn:  conn,
		table: `"public"."ai_health_report"`,
	}
}

func (m *customAiHealthReportModel) withSession(session sqlx.Session) AiHealthReportModel {
	return NewAiHealthReportModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customAiHealthReportModel) InsertReturningID(ctx context.Context, data *AiHealthReport) (int64, error) {
	query := fmt.Sprintf("insert into %s (user_id,trace_id,report_type,farm_code,tower_code,start_time,end_time,title,content,evidence,status) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id", m.table)
	var id int64
	err := m.conn.QueryRowCtx(ctx, &id, query, data.UserId, data.TraceId, data.ReportType, data.FarmCode, data.TowerCode, data.StartTime, data.EndTime, data.Title, data.Content, data.Evidence, data.Status)
	return id, err
}

func (m *customAiHealthReportModel) FindOne(ctx context.Context, id int64) (*AiHealthReport, error) {
	query := fmt.Sprintf("select id,user_id,trace_id,report_type,farm_code,tower_code,start_time,end_time,title,content,evidence,status,created_at from %s where id = $1 limit 1", m.table)
	var resp AiHealthReport
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
