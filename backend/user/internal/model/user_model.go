package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindOneByEmail(ctx context.Context,email string) (*User, error)
		withSession(session sqlx.Session) UserModel
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) FindOneByEmail(ctx context.Context,email string) (*User,error) {
	var resp User
	query := "select " + userRows +" from " +m.table + " where `email` = ? limit 1"
	err := m.QueryRowNoCacheCtx(ctx,&resp,query,email)
	switch err {
	case nil:
		return &resp,nil
	case sqlc.ErrNotFound:
		return nil,ErrNotFound
	default:
		return nil,err
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return &customUserModel{
		defaultUserModel: &defaultUserModel{
			CachedConn: sqlc.NewConn(sqlx.NewSqlConnFromSession(session),nil),
			table: m.table,
		},
	}
}