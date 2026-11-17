package logic

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"flashsale/flashsale"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateActivityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateActivityLogic {
	return &CreateActivityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateActivityLogic) CreateActivity(req *types.CreateActivityRequest) (resp *types.CreateActivityResponse, err error) {
	_, err = l.svcCtx.FlashsaleRpc.CreateActivity(l.ctx, &flashsale.CreateActivityRequest{
		ProductId: req.ProductId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Stock:     req.Stock,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateActivityResponse{}, nil
}
