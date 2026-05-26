package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WindCameraRecordModel = (*customWindCameraRecordModel)(nil)

type (
	// WindCameraRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindCameraRecordModel.
	WindCameraRecordModel interface {
		windCameraRecordModel
		withSession(session sqlx.Session) WindCameraRecordModel
	}

	customWindCameraRecordModel struct {
		*defaultWindCameraRecordModel
	}
)

// NewWindCameraRecordModel returns a model for the database table.
func NewWindCameraRecordModel(conn sqlx.SqlConn) WindCameraRecordModel {
	return &customWindCameraRecordModel{
		defaultWindCameraRecordModel: newWindCameraRecordModel(conn),
	}
}

func (m *customWindCameraRecordModel) withSession(session sqlx.Session) WindCameraRecordModel {
	return NewWindCameraRecordModel(sqlx.NewSqlConnFromSession(session))
}
