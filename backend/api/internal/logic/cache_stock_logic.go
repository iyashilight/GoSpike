package logic

import (
	"api/errorx"
	"api/internal/svc"
	"api/internal/types"
	"context"
	"flashsale/flashsale"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type CacheStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCacheStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheStockLogic {
	return &CacheStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CacheStockLogic) CacheStock(req *types.CacheStockRequest) (resp *types.CacheStockResponse, err error) {
	_, err = l.svcCtx.FlashsaleRpc.CacheStock(l.ctx,&flashsale.CacheStockRequest{
		ProductId: req.ProductId,
		Stock: req.Stock,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Internal:
			return nil, errorx.NewCodeError(errorx.ErrCacheStockFailed)
		default:
			return nil,err
		}
	}
	return &types.CacheStockResponse{},nil
}
