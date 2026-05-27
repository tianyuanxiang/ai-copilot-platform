package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiMaintenanceTicketDraftModel = (*customAiMaintenanceTicketDraftModel)(nil)

type (
	// AiMaintenanceTicketDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiMaintenanceTicketDraftModel.
	AiMaintenanceTicketDraftModel interface {
		aiMaintenanceTicketDraftModel
		withSession(session sqlx.Session) AiMaintenanceTicketDraftModel
	}

	customAiMaintenanceTicketDraftModel struct {
		*defaultAiMaintenanceTicketDraftModel
	}
)

// NewAiMaintenanceTicketDraftModel returns a model for the database table.
func NewAiMaintenanceTicketDraftModel(conn sqlx.SqlConn) AiMaintenanceTicketDraftModel {
	return &customAiMaintenanceTicketDraftModel{
		defaultAiMaintenanceTicketDraftModel: newAiMaintenanceTicketDraftModel(conn),
	}
}

func (m *customAiMaintenanceTicketDraftModel) withSession(session sqlx.Session) AiMaintenanceTicketDraftModel {
	return NewAiMaintenanceTicketDraftModel(sqlx.NewSqlConnFromSession(session))
}
