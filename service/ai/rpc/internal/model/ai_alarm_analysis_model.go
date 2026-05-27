package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiAlarmAnalysisModel = (*customAiAlarmAnalysisModel)(nil)

type (
	// AiAlarmAnalysisModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiAlarmAnalysisModel.
	AiAlarmAnalysisModel interface {
		aiAlarmAnalysisModel
		withSession(session sqlx.Session) AiAlarmAnalysisModel
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
