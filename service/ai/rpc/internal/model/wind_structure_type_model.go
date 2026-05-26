package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WindStructureTypeModel = (*customWindStructureTypeModel)(nil)

type (
	// WindStructureTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindStructureTypeModel.
	WindStructureTypeModel interface {
		windStructureTypeModel
		withSession(session sqlx.Session) WindStructureTypeModel
	}

	customWindStructureTypeModel struct {
		*defaultWindStructureTypeModel
	}
)

// NewWindStructureTypeModel returns a model for the database table.
func NewWindStructureTypeModel(conn sqlx.SqlConn) WindStructureTypeModel {
	return &customWindStructureTypeModel{
		defaultWindStructureTypeModel: newWindStructureTypeModel(conn),
	}
}

func (m *customWindStructureTypeModel) withSession(session sqlx.Session) WindStructureTypeModel {
	return NewWindStructureTypeModel(sqlx.NewSqlConnFromSession(session))
}
