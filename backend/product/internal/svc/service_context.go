package svc

import (
	"product/internal/config"
	"product/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	ProductModel model.ProductModel
	Conn sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config: c,
		ProductModel: model.NewProductModel(conn,c.Cache),
		Conn: conn,
	}
}
