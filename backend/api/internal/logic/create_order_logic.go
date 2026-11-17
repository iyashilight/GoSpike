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

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (resp *types.CreateOrderResponse, err error) {
	var pbItems []*order.OrderItem
	for _, item := range req.Items {
		pbItems = append(pbItems, &order.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}

	rpcResp, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderRequest{
		UserId: req.UserId,
		Items:  pbItems,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Internal:
			return nil, errorx.NewCodeError(errorx.ErrCreateOrderFailed)
		case codes.ResourceExhausted:
			return nil, errorx.NewCodeError(errorx.ErrStockNotEnough)
		default:
			return nil, err
		}
	}
	return &types.CreateOrderResponse{
		OrderId: rpcResp.OrderId,
	}, nil
}
