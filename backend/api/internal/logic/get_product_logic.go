package logic

import (
	"context"

	"api/errorx"
	"api/internal/svc"
	"api/internal/types"
	"product/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductRequest) (resp *types.GetProductResponse, err error) {
	rpcResp, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductRequest{
		Id: req.Id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return nil, errorx.NewCodeError(errorx.ErrProductNotFound)
		}
		return nil, err
	}
	return &types.GetProductResponse{
		Id:          rpcResp.Id,
		Name:        rpcResp.Name,
		Description: rpcResp.Description,
		Price:       rpcResp.Price,
		Stock:       rpcResp.Stock,
		CreatedAt:   rpcResp.CreatedAt,
	}, nil
}
