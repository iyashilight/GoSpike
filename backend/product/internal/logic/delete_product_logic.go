package logic

import (
	"context"

	"product/internal/model"
	"product/internal/svc"
	"product/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProductLogic) DeleteProduct(in *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}
	err := l.svcCtx.ProductModel.Delete(l.ctx,in.Id)
	if err != nil {
		if err == model.ErrNotFound {
		return nil,status.Error(codes.NotFound,"商品不存在")
	}
	return nil,status.Error(codes.Internal,"删除商品失败")
}
	return &product.DeleteProductResponse{}, nil
}
