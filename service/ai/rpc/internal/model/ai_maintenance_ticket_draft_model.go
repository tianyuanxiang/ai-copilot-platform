package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AiMaintenanceTicketDraftModel = (*customAiMaintenanceTicketDraftModel)(nil)

type (
	// AiMaintenanceTicketDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiMaintenanceTicketDraftModel.
	AiMaintenanceTicketDraftModel interface {
		aiMaintenanceTicketDraftModel
		withSession(session sqlx.Session) AiMaintenanceTicketDraftModel
		InsertReturningID(ctx context.Context, data *AiMaintenanceTicketDraft) (int64, error)
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

func (m *customAiMaintenanceTicketDraftModel) InsertReturningID(ctx context.Context, data *AiMaintenanceTicketDraft) (int64, error) {
	query := `insert into "public"."ai_maintenance_ticket_draft" (user_id,trace_id,farm_code,tower_code,alarm_code,title,content,evidence,status) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	var id int64
	err := m.conn.QueryRowCtx(ctx, &id, query, data.UserId, data.TraceId, data.FarmCode, data.TowerCode, data.AlarmCode, data.Title, data.Content, data.Evidence, data.Status)
	return id, err
}
