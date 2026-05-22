package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiSecurityEventModel = (*customAiSecurityEventModel)(nil)

type (
	// AiSecurityEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiSecurityEventModel.
	AiSecurityEventModel interface {
		aiSecurityEventModel
		withSession(session sqlx.Session) AiSecurityEventModel
	}

	customAiSecurityEventModel struct {
		*defaultAiSecurityEventModel
	}
)

// NewAiSecurityEventModel returns a model for the database table.
func NewAiSecurityEventModel(conn sqlx.SqlConn) AiSecurityEventModel {
	return &customAiSecurityEventModel{
		defaultAiSecurityEventModel: newAiSecurityEventModel(conn),
	}
}

func (m *customAiSecurityEventModel) withSession(session sqlx.Session) AiSecurityEventModel {
	return NewAiSecurityEventModel(sqlx.NewSqlConnFromSession(session))
}
