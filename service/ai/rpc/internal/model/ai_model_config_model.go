package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiModelConfigModel = (*customAiModelConfigModel)(nil)

type (
	// AiModelConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiModelConfigModel.
	AiModelConfigModel interface {
		aiModelConfigModel
		withSession(session sqlx.Session) AiModelConfigModel
	}

	customAiModelConfigModel struct {
		*defaultAiModelConfigModel
	}
)

// NewAiModelConfigModel returns a model for the database table.
func NewAiModelConfigModel(conn sqlx.SqlConn) AiModelConfigModel {
	return &customAiModelConfigModel{
		defaultAiModelConfigModel: newAiModelConfigModel(conn),
	}
}

func (m *customAiModelConfigModel) withSession(session sqlx.Session) AiModelConfigModel {
	return NewAiModelConfigModel(sqlx.NewSqlConnFromSession(session))
}
