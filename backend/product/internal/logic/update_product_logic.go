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

type UpdateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductLogic) UpdateProduct(in *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	if in.Id <= 0 || in.Name == "" || in.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}
	_, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
  if err != nil {
      return nil, status.Error(codes.NotFound, "商品不存在")
  }
	err = l.svcCtx.ProductModel.Update(l.ctx,&model.Product{
		Id: in.Id,
		Name: in.Name,
		Description: in.Description,
		Price: in.Price,
		Stock: int64(in.Stock),
	})
	if err != nil {
		return nil,status.Error(codes.Internal,"更新商品失败")
}
	return &product.UpdateProductResponse{}, nil
}
