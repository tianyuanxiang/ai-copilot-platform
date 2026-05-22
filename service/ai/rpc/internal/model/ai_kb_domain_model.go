package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiKbDomainModel = (*customAiKbDomainModel)(nil)

type (
	// AiKbDomainModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiKbDomainModel.
	AiKbDomainModel interface {
		aiKbDomainModel
		withSession(session sqlx.Session) AiKbDomainModel
	}

	customAiKbDomainModel struct {
		*defaultAiKbDomainModel
	}
)

// NewAiKbDomainModel returns a model for the database table.
func NewAiKbDomainModel(conn sqlx.SqlConn) AiKbDomainModel {
	return &customAiKbDomainModel{
		defaultAiKbDomainModel: newAiKbDomainModel(conn),
	}
}

func (m *customAiKbDomainModel) withSession(session sqlx.Session) AiKbDomainModel {
	return NewAiKbDomainModel(sqlx.NewSqlConnFromSession(session))
}
