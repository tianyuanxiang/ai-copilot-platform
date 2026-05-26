package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
func NewAiMessageModel(conn sqlx.SqlConn) AiMessageModel {
	return &customAiMessageModel{
		defaultAiMessageModel: newAiMessageModel(conn),
	}
}

func (m *customAiMessageModel) withSession(session sqlx.Session) AiMessageModel {
	return NewAiMessageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customAiMessageModel) ListByConversation(ctx context.Context, conversationID int64, userID int64, page int64, pageSize int64) ([]MessageListItem, int64, error) {
	page, pageSize = normalizePage(page, pageSize)
	countQuery := fmt.Sprintf(`
select count(1)
from %s m
join "public"."ai_conversation" c on c.id = m.conversation_id
where m.conversation_id = $1 and c.user_id = $2 and c.deleted_at is null and m.deleted_at is null`, m.table)
	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, conversationID, userID); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
select
  m.id as message_id,
  m.conversation_id,
  m.role,
  m.content,
  m.citations::text as citations,
  m.trace_id,
  m.created_at
from %s m
join "public"."ai_conversation" c on c.id = m.conversation_id
where m.conversation_id = $1 and c.user_id = $2 and c.deleted_at is null and m.deleted_at is null
order by m.created_at asc, m.id asc
limit $3 offset $4`, m.table)

	var list []MessageListItem
	if err := m.conn.QueryRowsCtx(ctx, &list, listQuery, conversationID, userID, pageSize, (page-1)*pageSize); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *customAiMessageModel) ListRecentByConversation(ctx context.Context, conversationID int64, userID int64, limit int64) ([]AiMessage, error) {
	if limit <= 0 {
		limit = 8
	}
	query := fmt.Sprintf(`
select %s
from %s m
join "public"."ai_conversation" c on c.id = m.conversation_id
where m.conversation_id = $1 and c.user_id = $2 and c.deleted_at is null and m.deleted_at is null
order by m.created_at desc, m.id desc
limit $3`, prefixedAiMessageRows("m"), m.table)

	var desc []AiMessage
	if err := m.conn.QueryRowsCtx(ctx, &desc, query, conversationID, userID, limit); err != nil {
		return nil, err
	}
	for i, j := 0, len(desc)-1; i < j; i, j = i+1, j-1 {
		desc[i], desc[j] = desc[j], desc[i]
	}
	return desc, nil
}

func (m *customAiMessageModel) ListActiveByConversation(ctx context.Context, conversationID int64, userID int64) ([]AiMessage, error) {
	query := fmt.Sprintf(`
select %s
from %s m
join "public"."ai_conversation" c on c.id = m.conversation_id
where m.conversation_id = $1 and c.user_id = $2 and c.deleted_at is null and m.deleted_at is null
order by m.created_at asc, m.id asc`, prefixedAiMessageRows("m"), m.table)

	var list []AiMessage
	if err := m.conn.QueryRowsCtx(ctx, &list, query, conversationID, userID); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customAiMessageModel) SoftDeleteByIDUserID(ctx context.Context, conversationID int64, messageID int64, userID int64, deletedBy int64) error {
	query := fmt.Sprintf(`
update %s m
set deleted_at = now(), deleted_by = $1
where m.id = $2
  and m.conversation_id = $3
  and m.deleted_at is null
  and exists (
    select 1
    from "public"."ai_conversation" c
    where c.id = m.conversation_id and c.user_id = $4 and c.deleted_at is null
  )`, m.table)
	ret, err := m.conn.ExecCtx(ctx, query, deletedBy, messageID, conversationID, userID)
	if err != nil {
		return err
	}
	return notFoundIfNoRows(ret)
}

func prefixedAiMessageRows(alias string) string {
	return alias + ".id," +
		alias + ".conversation_id," +
		alias + ".role," +
		alias + ".content," +
		alias + ".citations::text as citations," +
		alias + ".trace_id," +
		alias + ".created_at"
}
