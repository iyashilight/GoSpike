package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	ProductModel interface {
		productModel
		FindPage(ctx context.Context, page, size int32) ([]*Product, error)
		Count(ctx context.Context) (int64, error)
	}

	customProductModel struct {
		*defaultProductModel
	}
)

func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductModel {
	return &customProductModel{
		defaultProductModel: newProductModel(conn, c, opts...),
	}
}

func (m *customProductModel) FindPage(ctx context.Context, page, size int32) ([]*Product, error) {
	offset := (page - 1) * size
	query := fmt.Sprintf("select %s from %s order by id desc limit ?, ?", productRows, m.table)
	var resp []*Product
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, offset, size)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customProductModel) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s", m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return 0, err
	}
	return resp, nil
}
