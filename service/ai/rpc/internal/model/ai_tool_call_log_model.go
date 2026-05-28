package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiToolCallLogModel = (*customAiToolCallLogModel)(nil)

type (
	// AiToolCallLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiToolCallLogModel.
	AiToolCallLogModel interface {
		aiToolCallLogModel
		withSession(session sqlx.Session) AiToolCallLogModel
		ListByFilter(ctx context.Context, userID int64, traceID string, toolName string, status string, page int64, pageSize int64) (int64, []AiToolCallLog, error)
	}

	customAiToolCallLogModel struct {
		*defaultAiToolCallLogModel
		db *gorm.DB
	}
)

// NewAiToolCallLogModel returns a model for the database table.
func NewAiToolCallLogModel(conn sqlx.SqlConn, db *gorm.DB) AiToolCallLogModel {
	return &customAiToolCallLogModel{
		defaultAiToolCallLogModel: newAiToolCallLogModel(conn),
		db:                        db,
	}
}

func (m *customAiToolCallLogModel) withSession(session sqlx.Session) AiToolCallLogModel {
	return NewAiToolCallLogModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiToolCallLogModel) ListByFilter(ctx context.Context, userID int64, traceID string, toolName string, status string, page int64, pageSize int64) (int64, []AiToolCallLog, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	query := m.db.WithContext(ctx).Table("ai_tool_call_log")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if traceID != "" {
		query = query.Where("trace_id = ?", traceID)
	}
	if toolName != "" {
		query = query.Where("tool_name = ?", toolName)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []AiToolCallLog
	err := query.Order("created_at desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Scan(&rows).Error
	return total, rows, err
}
