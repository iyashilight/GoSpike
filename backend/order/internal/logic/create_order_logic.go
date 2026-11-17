package logic

import (
	"context"
	"fmt"
	"product/product"

	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if in.UserId <= 0 || len(in.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}
	for _, item := range in.Items {
		if item.ProductId <= 0 || item.Quantity <= 0 || item.Price <= 0 {
			return nil, status.Error(codes.InvalidArgument, "订单项参数无效")
		}
	}
	var totalPrice float64
	for _, item := range in.Items {
		totalPrice += item.Price * float64(item.Quantity)
	}
	for _, item := range in.Items {
			_, err := l.svcCtx.ProductRpc.DecrStock(l.ctx,&product.DecrStockRequest{
				Id: item.ProductId,
				Quantity: item.Quantity,
			})
			if err != nil {
				return nil,status.Error(codes.ResourceExhausted,fmt.Sprintf("商品 %d 库存不足: %v",item.ProductId,err))
			}
		}
	var orderId int64
	err := l.svcCtx.SqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		result, err := session.ExecCtx(ctx,
			"insert into `orders` (`user_id`, `total_price`, `status`) values (?, ?, ?)",
			in.UserId, totalPrice, "pending",
		)
		if err != nil {
			return err
		}
		orderId, _ = result.LastInsertId()
		for _, item := range in.Items {
			_, err = session.ExecCtx(ctx,
				"insert into `order_item` (`order_id`, `product_id`, `product_name`, `price`, `quantity`) values (?, ?, ?, ?, ?)",
				orderId, item.ProductId, item.ProductName, item.Price, item.Quantity,
			)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("创建订单失败: %v",
			err))
	}
	return &order.CreateOrderResponse{OrderId: orderId}, nil
}
