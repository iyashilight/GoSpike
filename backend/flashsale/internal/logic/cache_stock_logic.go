package logic

import (
	"context"
	"fmt"
	"strconv"

	"flashsale/flashsale"
	"flashsale/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CacheStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCacheStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheStockLogic {
	return &CacheStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CacheStockLogic) CacheStock(in *flashsale.CacheStockRequest) (*flashsale.CacheStockResponse, error) {
	activityKey := fmt.Sprintf("activity:%d", in.ProductId)
	exists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, activityKey)
	if err != nil || !exists {
		return nil, status.Error(codes.NotFound, "秒杀活动不存在，请先创建活动")
	}

	key := fmt.Sprintf("stock:%d", in.ProductId)
	err = l.svcCtx.Redis.SetCtx(l.ctx, key, strconv.Itoa(int(in.Stock)))
	if err != nil {
		return nil, status.Error(codes.Internal, "缓存库存失败")
	}
	return &flashsale.CacheStockResponse{}, nil
}
