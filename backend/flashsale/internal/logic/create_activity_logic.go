package logic

import (
	"context"
	"fmt"

	"flashsale/flashsale"
	"flashsale/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateActivityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateActivityLogic {
	return &CreateActivityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateActivityLogic) CreateActivity(in *flashsale.CreateActivityRequest) (*flashsale.CreateActivityResponse, error) {
	if in.StartTime >= in.EndTime {
		return nil, status.Error(codes.InvalidArgument, "开始时间必须早于结束时间")
	}
	if in.Stock <= 0 {
		return nil, status.Error(codes.InvalidArgument, "库存必须大于 0")
	}

	key := fmt.Sprintf("activity:%d", in.ProductId)
	err := l.svcCtx.Redis.HmsetCtx(l.ctx, key, map[string]string{
		"start_time": fmt.Sprintf("%d", in.StartTime),
		"end_time":   fmt.Sprintf("%d", in.EndTime),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "创建秒杀活动失败")
	}

	// 预热库存
	stockKey := fmt.Sprintf("stock:%d", in.ProductId)
	if err := l.svcCtx.Redis.SetCtx(l.ctx, stockKey, fmt.Sprintf("%d", in.Stock)); err != nil {
		return nil, status.Error(codes.Internal, "预热库存失败")
	}

	return &flashsale.CreateActivityResponse{}, nil
}
