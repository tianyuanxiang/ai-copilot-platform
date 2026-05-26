package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AiConversationModel = (*customAiConversationModel)(nil)

type (
	// AiConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiConversationModel.
	AiConversationModel interface {
		aiConversationModel
		withSession(session sqlx.Session) AiConversationModel
		InsertReturningID(ctx context.Context, data *AiConversation) (int64, error)
		FindByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversation, error)
		FindContextByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversationContext, error)
		ListByUser(ctx context.Context, query ConversationListQuery) ([]ConversationListItem, int64, error)
		UpdateTitleByIDUserID(ctx context.Context, id int64, userID int64, title string) error
		SoftDeleteByIDUserID(ctx context.Context, id int64, userID int64, deletedBy int64) error
		UpdateSummaryByIDUserID(ctx context.Context, id int64, userID int64, summary string) error
		TouchUpdatedAt(ctx context.Context, id int64) error
	}

	customAiConversationModel struct {
		*defaultAiConversationModel
	}

	AiConversationContext struct {
		Id                  int64  `db:"id"`
		UserId              int64  `db:"user_id"`
		KbId                int64  `db:"kb_id"`
		HasKbId             bool   `db:"has_kb_id"`
		Title               string `db:"title"`
		ConversationSummary string `db:"conversation_summary"`
	}

	ConversationListQuery struct {
		Page     int64
		PageSize int64
		UserID   int64
		KbID     int64
		HasKbID  bool
		Keyword  string
	}

	ConversationListItem struct {
		Id            int64     `db:"id"`
		KbId          int64     `db:"kb_id"`
		HasKbId       bool      `db:"has_kb_id"`
		Title         string    `db:"title"`
		LatestMessage string    `db:"latest_message"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
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

func (m *customAiConversationModel) InsertReturningID(ctx context.Context, data *AiConversation) (int64, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3) returning id", m.table, aiConversationRowsExpectAutoSet)
	var id int64
	err := m.conn.QueryRowCtx(ctx, &id, query, data.UserId, data.KbId, data.Title)
	return id, err
}

func (m *customAiConversationModel) FindByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversation, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 and user_id = $2 and deleted_at is null limit 1", aiConversationRows, m.table)
	var resp AiConversation
	err := m.conn.QueryRowCtx(ctx, &resp, query, id, userID)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customAiConversationModel) FindContextByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversationContext, error) {
	query := fmt.Sprintf("select id,user_id,coalesce(kb_id, 0) as kb_id,kb_id is not null as has_kb_id,title,conversation_summary from %s where id = $1 and user_id = $2 and deleted_at is null limit 1", m.table)
	var resp AiConversationContext
	err := m.conn.QueryRowCtx(ctx, &resp, query, id, userID)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customAiConversationModel) ListByUser(ctx context.Context, query ConversationListQuery) ([]ConversationListItem, int64, error) {
	where, args := conversationListWhere(query)

	countQuery := fmt.Sprintf("select count(1) from %s c where %s", m.table, strings.Join(where, " and "))
	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePage(query.Page, query.PageSize)
	args = append(args, pageSize, (page-1)*pageSize)
	listQuery := fmt.Sprintf(`
select
  c.id,
  coalesce(c.kb_id, 0) as kb_id,
  c.kb_id is not null as has_kb_id,
  c.title,
  coalesce((
    select m.content
    from "public"."ai_message" m
    where m.conversation_id = c.id and m.deleted_at is null
    order by m.created_at desc, m.id desc
    limit 1
  ), '') as latest_message,
  c.created_at,
  c.updated_at
from %s c
where %s
order by c.updated_at desc, c.id desc
limit $%d offset $%d`, m.table, strings.Join(where, " and "), len(args)-1, len(args))

	var list []ConversationListItem
	if err := m.conn.QueryRowsCtx(ctx, &list, listQuery, args...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *customAiConversationModel) UpdateTitleByIDUserID(ctx context.Context, id int64, userID int64, title string) error {
	query := fmt.Sprintf("update %s set title = $1, updated_at = now() where id = $2 and user_id = $3 and deleted_at is null", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, title, id, userID)
	if err != nil {
		return err
	}
	return notFoundIfNoRows(ret)
}

func (m *customAiConversationModel) SoftDeleteByIDUserID(ctx context.Context, id int64, userID int64, deletedBy int64) error {
	query := fmt.Sprintf("update %s set deleted_at = now(), deleted_by = $1, updated_at = now() where id = $2 and user_id = $3 and deleted_at is null", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, deletedBy, id, userID)
	if err != nil {
		return err
	}
	return notFoundIfNoRows(ret)
}

func (m *customAiConversationModel) UpdateSummaryByIDUserID(ctx context.Context, id int64, userID int64, summary string) error {
	query := fmt.Sprintf("update %s set conversation_summary = $1, updated_at = now() where id = $2 and user_id = $3 and deleted_at is null", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, summary, id, userID)
	if err != nil {
		return err
	}
	return notFoundIfNoRows(ret)
}

func (m *customAiConversationModel) TouchUpdatedAt(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set updated_at = now() where id = $1 and deleted_at is null", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func conversationListWhere(query ConversationListQuery) ([]string, []interface{}) {
	where := []string{"c.user_id = $1", "c.deleted_at is null"}
	args := []interface{}{query.UserID}
	if query.HasKbID && query.KbID > 0 {
		args = append(args, query.KbID)
		where = append(where, fmt.Sprintf("c.kb_id = $%d", len(args)))
	}
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		args = append(args, "%"+keyword+"%")
		where = append(where, fmt.Sprintf("c.title ilike $%d", len(args)))
	}
	return where, args
}

func normalizePage(page int64, pageSize int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func notFoundIfNoRows(ret sql.Result) error {
	affected, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
