package logic

import (
	"context"

	"api/errorx"
	"api/internal/svc"
	"api/internal/types"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderRequest) (resp *types.GetOrderResponse, err error) {
	rpcResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
		OrderId: req.OrderId,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return nil, errorx.NewCodeError(errorx.ErrOrderNotFound)
		}
		return nil, err
	}
	var items []types.OrderItemParam
	for _, item := range rpcResp.Items {
		items = append(items, types.OrderItemParam{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}
	return &types.GetOrderResponse{
		OrderId:    rpcResp.OrderId,
		UserId:     rpcResp.UserId,
		TotalPrice: rpcResp.TotalPrice,
		Status:     rpcResp.Status,
		Items:      items,
		CreatedAt:  rpcResp.CreatedAt,
	}, nil
}
