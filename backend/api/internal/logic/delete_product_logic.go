// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"product/product"

	"api/errorx"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLogic) DeleteProduct(req *types.DeleteProductRequest) (resp *types.DeleteProductResponse, err error) {
	_, err = l.svcCtx.ProductRpc.DeleteProduct(l.ctx, &product.DeleteProductRequest{
		Id: req.Id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, errorx.NewCodeError(errorx.ErrProductNotFound)
		case codes.Internal:
			return nil, errorx.NewCodeError(errorx.ErrDeleteProductFailed)
		default:
			return nil, err
		}
	}
	return &types.DeleteProductResponse{}, nil
}
