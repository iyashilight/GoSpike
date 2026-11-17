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

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrdersLogic) ListOrders(in *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	size := in.Size
	if size < 1 || size > 100 {
		size = 20
	}

	var items []*order.GetOrderResponse
	var total int64
	var err error

	if in.UserId > 0 {
		orders, findErr := l.svcCtx.OrdersModel.FindPageByUserId(l.ctx, in.UserId, page, size)
		if findErr != nil {
			return nil, status.Error(codes.Internal, "查询订单列表失败")
		}
		total, err = l.svcCtx.OrdersModel.CountByUserId(l.ctx, in.UserId)
		if err != nil {
			return nil, status.Error(codes.Internal, "查询订单总数失败")
		}
		for _, o := range orders {
			pbItems, itemErr := l.getOrderItems(o.Id)
			if itemErr != nil {
				return nil, itemErr
			}
			items = append(items, &order.GetOrderResponse{
				OrderId:    o.Id,
				UserId:     o.UserId,
				TotalPrice: o.TotalPrice,
				Status:     o.Status,
				Items:      pbItems,
				CreatedAt:  o.CreatedAt.Format(time.DateTime),
			})
		}
	} else {
		orders, findErr := l.svcCtx.OrdersModel.FindPage(l.ctx, page, size)
		if findErr != nil {
			return nil, status.Error(codes.Internal, "查询订单列表失败")
		}
		total, err = l.svcCtx.OrdersModel.Count(l.ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "查询订单总数失败")
		}
		for _, o := range orders {
			pbItems, itemErr := l.getOrderItems(o.Id)
			if itemErr != nil {
				return nil, itemErr
			}
			items = append(items, &order.GetOrderResponse{
				OrderId:    o.Id,
				UserId:     o.UserId,
				TotalPrice: o.TotalPrice,
				Status:     o.Status,
				Items:      pbItems,
				CreatedAt:  o.CreatedAt.Format(time.DateTime),
			})
		}
	}

	return &order.ListOrdersResponse{
		Items: items,
		Total: total,
	}, nil
}

func (l *ListOrdersLogic) getOrderItems(orderId int64) ([]*order.OrderItem, error) {
	orderItems, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, orderId)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询订单项失败")
	}
	var pbItems []*order.OrderItem
	for _, item := range orderItems {
		pbItems = append(pbItems, &order.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
		})
	}
	return pbItems, nil
}
