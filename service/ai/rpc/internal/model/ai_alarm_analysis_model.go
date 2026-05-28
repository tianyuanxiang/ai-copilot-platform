package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AiAlarmAnalysisModel = (*customAiAlarmAnalysisModel)(nil)

type (
	// AiAlarmAnalysisModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiAlarmAnalysisModel.
	AiAlarmAnalysisModel interface {
		aiAlarmAnalysisModel
		withSession(session sqlx.Session) AiAlarmAnalysisModel
		InsertReturningID(ctx context.Context, data *AiAlarmAnalysis) (int64, error)
	}

	customAiAlarmAnalysisModel struct {
		*defaultAiAlarmAnalysisModel
	}
)

// NewAiAlarmAnalysisModel returns a model for the database table.
func NewAiAlarmAnalysisModel(conn sqlx.SqlConn) AiAlarmAnalysisModel {
	return &customAiAlarmAnalysisModel{
		defaultAiAlarmAnalysisModel: newAiAlarmAnalysisModel(conn),
	}
}

func (m *customAiAlarmAnalysisModel) withSession(session sqlx.Session) AiAlarmAnalysisModel {
	return NewAiAlarmAnalysisModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customAiAlarmAnalysisModel) InsertReturningID(ctx context.Context, data *AiAlarmAnalysis) (int64, error) {
	query := `insert into "public"."ai_alarm_analysis" (user_id,trace_id,farm_code,tower_code,alarm_code,title,content,evidence,status) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	var id int64
	err := m.conn.QueryRowCtx(ctx, &id, query, data.UserId, data.TraceId, data.FarmCode, data.TowerCode, data.AlarmCode, data.Title, data.Content, data.Evidence, data.Status)
	return id, err
}
