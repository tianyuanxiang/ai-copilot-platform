package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiSecurityAlertModel = (*customAiSecurityAlertModel)(nil)

type (
	// AiSecurityAlertModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiSecurityAlertModel.
	AiSecurityAlertModel interface {
		aiSecurityAlertModel
		withSession(session sqlx.Session) AiSecurityAlertModel
	}

	customAiSecurityAlertModel struct {
		*defaultAiSecurityAlertModel
	}
)

// NewAiSecurityAlertModel returns a model for the database table.
func NewAiSecurityAlertModel(conn sqlx.SqlConn) AiSecurityAlertModel {
	return &customAiSecurityAlertModel{
		defaultAiSecurityAlertModel: newAiSecurityAlertModel(conn),
	}
}

func (m *customAiSecurityAlertModel) withSession(session sqlx.Session) AiSecurityAlertModel {
	return NewAiSecurityAlertModel(sqlx.NewSqlConnFromSession(session))
}
