package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiMessageModel = (*customAiMessageModel)(nil)

type (
	// AiMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiMessageModel.
	AiMessageModel interface {
		aiMessageModel
		withSession(session sqlx.Session) AiMessageModel
	}

	customAiMessageModel struct {
		*defaultAiMessageModel
	}
)

// NewAiMessageModel returns a model for the database table.
func NewAiMessageModel(conn sqlx.SqlConn) AiMessageModel {
	return &customAiMessageModel{
		defaultAiMessageModel: newAiMessageModel(conn),
	}
}

func (m *customAiMessageModel) withSession(session sqlx.Session) AiMessageModel {
	return NewAiMessageModel(sqlx.NewSqlConnFromSession(session))
}
