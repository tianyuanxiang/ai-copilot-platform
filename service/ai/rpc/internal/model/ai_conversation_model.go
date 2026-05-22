package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiConversationModel = (*customAiConversationModel)(nil)

type (
	// AiConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiConversationModel.
	AiConversationModel interface {
		aiConversationModel
		withSession(session sqlx.Session) AiConversationModel
	}

	customAiConversationModel struct {
		*defaultAiConversationModel
	}
)

// NewAiConversationModel returns a model for the database table.
func NewAiConversationModel(conn sqlx.SqlConn) AiConversationModel {
	return &customAiConversationModel{
		defaultAiConversationModel: newAiConversationModel(conn),
	}
}

func (m *customAiConversationModel) withSession(session sqlx.Session) AiConversationModel {
	return NewAiConversationModel(sqlx.NewSqlConnFromSession(session))
}
