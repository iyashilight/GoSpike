// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"api/internal/config"
	"flashsale/flashsale_client"
	"order/order_client"
	"product/product_client"
	"user/user_client"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	Redis	*redis.Redis
	UserRpc    user_client.User
	ProductRpc product_client.Product
	OrderRpc   order_client.Order
	FlashsaleRpc flashsale_client.Flashsale
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Redis: redis.MustNewRedis(c.BizRedis),
		UserRpc:    user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc: product_client.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		OrderRpc:   order_client.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		FlashsaleRpc: flashsale_client.NewFlashsale(zrpc.MustNewClient(c.FlashsaleRpc)),
	}
}
