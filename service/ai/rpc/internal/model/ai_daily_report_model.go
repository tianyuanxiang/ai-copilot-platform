package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiDailyReportModel = (*customAiDailyReportModel)(nil)

type (
	// AiDailyReportModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiDailyReportModel.
	AiDailyReportModel interface {
		aiDailyReportModel
		withSession(session sqlx.Session) AiDailyReportModel
	}

	customAiDailyReportModel struct {
		*defaultAiDailyReportModel
	}
)

// NewAiDailyReportModel returns a model for the database table.
func NewAiDailyReportModel(conn sqlx.SqlConn) AiDailyReportModel {
	return &customAiDailyReportModel{
		defaultAiDailyReportModel: newAiDailyReportModel(conn),
	}
}

func (m *customAiDailyReportModel) withSession(session sqlx.Session) AiDailyReportModel {
	return NewAiDailyReportModel(sqlx.NewSqlConnFromSession(session))
}
