package logic

import (
	"context"

	"product/internal/svc"
	"product/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DecrStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrStockLogic {
	return &DecrStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecrStockLogic) DecrStock(in *product.DecrStockRequest) (*product.DecrStockResponse, error) {
	if in.Id <= 0 || in.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}
	query := "update product set stock = stock - ? where id = ? and stock >= ?"
	result, err := l.svcCtx.Conn.ExecCtx(l.ctx, query, in.Quantity, in.Id, in.Quantity)
	if err != nil {
		return nil,status.Error(codes.Internal,"扣库存失败")
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, status.Error(codes.ResourceExhausted,"库存不足")
	}
	return &product.DecrStockResponse{}, nil
}
