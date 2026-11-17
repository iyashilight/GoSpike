package logic

import (
	"context"
	"time"

	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	if in.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}
	o, err := l.svcCtx.OrdersModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "订单不存在")
	}
	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, in.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询订单失败")
	}
	var pbItems []*order.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &order.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
		})
	}
	return &order.GetOrderResponse{
		OrderId:    o.Id,
		UserId:     o.UserId,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		Items:      pbItems,
		CreatedAt:  o.CreatedAt.Format(time.DateTime),
	}, nil
}
