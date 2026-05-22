package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiToolCallLogModel = (*customAiToolCallLogModel)(nil)

type (
	// AiToolCallLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiToolCallLogModel.
	AiToolCallLogModel interface {
		aiToolCallLogModel
		withSession(session sqlx.Session) AiToolCallLogModel
	}

	customAiToolCallLogModel struct {
		*defaultAiToolCallLogModel
	}
)

// NewAiToolCallLogModel returns a model for the database table.
func NewAiToolCallLogModel(conn sqlx.SqlConn) AiToolCallLogModel {
	return &customAiToolCallLogModel{
		defaultAiToolCallLogModel: newAiToolCallLogModel(conn),
	}
}

func (m *customAiToolCallLogModel) withSession(session sqlx.Session) AiToolCallLogModel {
	return NewAiToolCallLogModel(sqlx.NewSqlConnFromSession(session))
}
