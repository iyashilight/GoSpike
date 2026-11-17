package logic

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"api/errorx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderRequest) (resp *types.CancelOrderResponse, err error) {
	_, err = l.svcCtx.OrderRpc.CancelOrder(l.ctx, &order.CancelOrderRequest{
		OrderId: req.OrderId,
		UserId:  req.UserId,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return nil, errorx.NewCodeError(errorx.ErrOrderNotFound)
		}
		return nil, err
	}
	return &types.CancelOrderResponse{}, nil
}
