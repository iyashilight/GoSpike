package logic

import (
	"context"
	"fmt"

	"order/internal/model"
	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	if in.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}

	o, err := l.svcCtx.OrdersModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "订单不存在")
		}
		return nil, status.Error(codes.Internal, "查询订单失败")
	}

	if in.UserId > 0 && o.UserId != in.UserId {
		return nil, status.Error(codes.PermissionDenied, "无权操作此订单")
	}

	if o.Status != "pending" {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("当前订单状态为 %s，无法取消", o.Status))
	}

	err = l.svcCtx.SqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, execErr := session.ExecCtx(ctx,
			"update `orders` set `status` = ? where `id` = ? and `status` = ?",
			"cancelled", in.OrderId, "pending",
		)
		return execErr
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "取消订单失败")
	}

	return &order.CancelOrderResponse{}, nil
}
