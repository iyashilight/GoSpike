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

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductRequest) (resp *types.CreateProductResponse, err error) {
	rpcResp, err := l.svcCtx.ProductRpc.CreateProduct(l.ctx, &product.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Internal {
			return nil, errorx.NewCodeError(errorx.ErrCreateProductFailed)
		}
		return nil, err
	}
	return &types.CreateProductResponse{
		Id: rpcResp.Id,
	}, nil
}
