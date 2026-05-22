package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiLlmCallLogModel = (*customAiLlmCallLogModel)(nil)

type (
	// AiLlmCallLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiLlmCallLogModel.
	AiLlmCallLogModel interface {
		aiLlmCallLogModel
		withSession(session sqlx.Session) AiLlmCallLogModel
	}

	customAiLlmCallLogModel struct {
		*defaultAiLlmCallLogModel
	}
)

// NewAiLlmCallLogModel returns a model for the database table.
func NewAiLlmCallLogModel(conn sqlx.SqlConn) AiLlmCallLogModel {
	return &customAiLlmCallLogModel{
		defaultAiLlmCallLogModel: newAiLlmCallLogModel(conn),
	}
}

func (m *customAiLlmCallLogModel) withSession(session sqlx.Session) AiLlmCallLogModel {
	return NewAiLlmCallLogModel(sqlx.NewSqlConnFromSession(session))
}
