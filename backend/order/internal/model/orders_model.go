package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrdersModel = (*customOrdersModel)(nil)

type (
	OrdersModel interface {
		ordersModel
		FindPageByUserId(ctx context.Context, userId int64, page, size int32) ([]*Orders, error)
		FindPage(ctx context.Context, page, size int32) ([]*Orders, error)
		CountByUserId(ctx context.Context, userId int64) (int64, error)
		Count(ctx context.Context) (int64, error)
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

func NewOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrdersModel {
	return &customOrdersModel{
		defaultOrdersModel: newOrdersModel(conn, c, opts...),
	}
}

func (m *customOrdersModel) FindPageByUserId(ctx context.Context, userId int64, page, size int32) ([]*Orders, error) {
	offset := (page - 1) * size
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by id desc limit ?, ?", ordersRows, m.table)
	var resp []*Orders
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, offset, size)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customOrdersModel) FindPage(ctx context.Context, page, size int32) ([]*Orders, error) {
	offset := (page - 1) * size
	query := fmt.Sprintf("select %s from %s order by id desc limit ?, ?", ordersRows, m.table)
	var resp []*Orders
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, offset, size)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customOrdersModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId)
	if err != nil {
		return 0, err
	}
	return resp, nil
}

func (m *customOrdersModel) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s", m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return 0, err
	}
	return resp, nil
}
