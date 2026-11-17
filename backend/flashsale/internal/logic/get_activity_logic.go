package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"flashsale/flashsale"
	"flashsale/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetActivityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActivityLogic {
	return &GetActivityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetActivityLogic) GetActivity(in *flashsale.GetActivityRequest) (*flashsale.GetActivityResponse, error) {
	key := fmt.Sprintf("activity:%d", in.ProductId)
	exists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, key)
	if err != nil || !exists {
		return nil, status.Error(codes.NotFound, "秒杀活动不存在")
	}

	vals, err := l.svcCtx.Redis.HgetallCtx(l.ctx, key)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询秒杀活动失败")
	}

	startTime, _ := strconv.ParseInt(vals["start_time"], 10, 64)
	endTime, _ := strconv.ParseInt(vals["end_time"], 10, 64)
	now := time.Now().Unix()

	stockKey := fmt.Sprintf("stock:%d", in.ProductId)
	stockStr, err := l.svcCtx.Redis.GetCtx(l.ctx, stockKey)
	stock := int32(0)
	if err == nil {
		s, _ := strconv.Atoi(stockStr)
		stock = int32(s)
	}

	return &flashsale.GetActivityResponse{
		ProductId: in.ProductId,
		StartTime: startTime,
		EndTime:   endTime,
		Stock:     stock,
		Active:    now >= startTime && now <= endTime,
	}, nil
}
