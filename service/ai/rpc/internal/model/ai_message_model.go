package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiMessageModel = (*customAiMessageModel)(nil)

type (
	// AiMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiMessageModel.
	AiMessageModel interface {
		aiMessageModel
		withSession(session sqlx.Session) AiMessageModel
		ListByConversation(ctx context.Context, conversationID int64, userID int64, page int64, pageSize int64) ([]MessageListItem, int64, error)
		ListRecentByConversation(ctx context.Context, conversationID int64, userID int64, limit int64) ([]AiMessage, error)
		ListActiveByConversation(ctx context.Context, conversationID int64, userID int64) ([]AiMessage, error)
		SoftDeleteByIDUserID(ctx context.Context, conversationID int64, messageID int64, userID int64, deletedBy int64) error
	}

	customAiMessageModel struct {
		*defaultAiMessageModel
		db *gorm.DB
	}

	MessageListItem struct {
		MessageId      int64     `db:"message_id"`
		ConversationId int64     `db:"conversation_id"`
		Role           string    `db:"role"`
		Content        string    `db:"content"`
		Citations      string    `db:"citations"`
		TraceId        string    `db:"trace_id"`
		CreatedAt      time.Time `db:"created_at"`
	}
)

// NewAiMessageModel returns a model for the database table.
func NewAiMessageModel(conn sqlx.SqlConn, db *gorm.DB) AiMessageModel {
	return &customAiMessageModel{
		defaultAiMessageModel: newAiMessageModel(conn),
		db:                    db,
	}
}

func (m *customAiMessageModel) withSession(session sqlx.Session) AiMessageModel {
	return NewAiMessageModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiMessageModel) ListByConversation(ctx context.Context, conversationID int64, userID int64, page int64, pageSize int64) ([]MessageListItem, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	db := m.db.WithContext(ctx).Table("ai_message as m").
		Joins("join ai_conversation c on c.id = m.conversation_id").
		Where("m.conversation_id = ? and c.user_id = ? and c.deleted_at is null and m.deleted_at is null", conversationID, userID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []MessageListItem
	err := db.Select("m.id as message_id,m.conversation_id,m.role,m.content,m.citations::text as citations,m.trace_id,m.created_at").
		Order("m.created_at asc, m.id asc").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *customAiMessageModel) ListRecentByConversation(ctx context.Context, conversationID int64, userID int64, limit int64) ([]AiMessage, error) {
	if limit <= 0 {
		limit = 8
	}
	var desc []AiMessage
	err := m.db.WithContext(ctx).Table("ai_message as m").
		Select(aiMessageSelectColumns("m")).
		Joins("join ai_conversation c on c.id = m.conversation_id").
		Where("m.conversation_id = ? and c.user_id = ? and c.deleted_at is null and m.deleted_at is null", conversationID, userID).
		Order("m.created_at desc, m.id desc").
		Limit(int(limit)).
		Scan(&desc).Error
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(desc)-1; i < j; i, j = i+1, j-1 {
		desc[i], desc[j] = desc[j], desc[i]
	}
	return desc, nil
}

func (m *customAiMessageModel) ListActiveByConversation(ctx context.Context, conversationID int64, userID int64) ([]AiMessage, error) {
	var list []AiMessage
	err := m.db.WithContext(ctx).Table("ai_message as m").
		Select(aiMessageSelectColumns("m")).
		Joins("join ai_conversation c on c.id = m.conversation_id").
		Where("m.conversation_id = ? and c.user_id = ? and c.deleted_at is null and m.deleted_at is null", conversationID, userID).
		Order("m.created_at asc, m.id asc").
		Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customAiMessageModel) SoftDeleteByIDUserID(ctx context.Context, conversationID int64, messageID int64, userID int64, deletedBy int64) error {
	result := m.db.WithContext(ctx).Table("ai_message as m").
		Where("m.id = ? and m.conversation_id = ? and m.deleted_at is null", messageID, conversationID).
		Where("exists (select 1 from ai_conversation c where c.id = m.conversation_id and c.user_id = ? and c.deleted_at is null)", userID).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		})
	return gormNotFoundIfNoRows(result)
}

func aiMessageSelectColumns(alias string) string {
	return alias + ".id," +
		alias + ".conversation_id," +
		alias + ".role," +
		alias + ".content," +
		alias + ".citations::text as citations," +
		alias + ".trace_id," +
		alias + ".created_at"
}
