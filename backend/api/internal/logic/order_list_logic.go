package logic

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderListLogic) OrderList(req *types.OrderListRequest) (resp *types.OrderListResponse, err error) {
	rpcResp, err := l.svcCtx.OrderRpc.ListOrders(l.ctx, &order.ListOrdersRequest{
		UserId: req.UserId,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	var items []types.GetOrderResponse
	for _, item := range rpcResp.Items {
		var orderItems []types.OrderItemParam
		for _, oi := range item.Items {
			orderItems = append(orderItems, types.OrderItemParam{
				ProductId:   oi.ProductId,
				ProductName: oi.ProductName,
				Price:       oi.Price,
				Quantity:    oi.Quantity,
			})
		}
		items = append(items, types.GetOrderResponse{
			OrderId:    item.OrderId,
			UserId:     item.UserId,
			TotalPrice: item.TotalPrice,
			Status:     item.Status,
			Items:      orderItems,
			CreatedAt:  item.CreatedAt,
		})
	}
	return &types.OrderListResponse{
		Items: items,
		Total: rpcResp.Total,
	}, nil
}
