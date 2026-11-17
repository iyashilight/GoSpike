package logic

import (
	"context"

	"api/errorx"
	"api/internal/svc"
	"api/internal/types"
	"flashsale/flashsale"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FlashPurchaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlashPurchaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlashPurchaseLogic {
	return &FlashPurchaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlashPurchaseLogic) FlashPurchase(req *types.FlashPurchaseRequest) (resp *types.FlashPurchaseResponse, err error) {
	rpcResp, err := l.svcCtx.FlashsaleRpc.FlashPurchase(l.ctx, &flashsale.FlashPurchaseRequest{
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.ResourceExhausted:
			return nil, errorx.NewCodeError(errorx.ErrStockNotEnough)
		case codes.NotFound:
			return nil, errorx.NewCodeError(errorx.ErrProductNotFound)
		default:
			return nil, err
		}
	}
	return &types.FlashPurchaseResponse{
		OrderId: rpcResp.OrderId,
		Status:  rpcResp.Status,
	}, nil
}
