package logic

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"flashsale/flashsale"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"api/errorx"
)

type GetActivityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActivityLogic {
	return &GetActivityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetActivityLogic) GetActivity(req *types.GetActivityRequest) (resp *types.GetActivityResponse, err error) {
	rpcResp, err := l.svcCtx.FlashsaleRpc.GetActivity(l.ctx, &flashsale.GetActivityRequest{
		ProductId: req.ProductId,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return nil, errorx.NewCodeError(errorx.ErrProductNotFound)
		}
		return nil, err
	}
	return &types.GetActivityResponse{
		ProductId: rpcResp.ProductId,
		StartTime: rpcResp.StartTime,
		EndTime:   rpcResp.EndTime,
		Stock:     rpcResp.Stock,
		Active:    rpcResp.Active,
	}, nil
}
