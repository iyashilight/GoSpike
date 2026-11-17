package logic

import (
	"context"
	"time"

	"product/internal/model"
	"product/internal/svc"
	"product/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductLogic) GetProduct(in *product.GetProductRequest) (*product.GetProductResponse, error) {
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
  if err != nil {
      if err == model.ErrNotFound {
          return nil, status.Error(codes.NotFound, "商品不存在")
      }
      return nil, status.Error(codes.Internal, "查询商品失败")
  }
	return &product.GetProductResponse{
		Id: p.Id,
		Name: p.Name,
		Description: p.Description,
		Price: p.Price,
		Stock: int32(p.Stock),
		CreatedAt: p.CreatedAt.Format(time.DateTime),
	}, nil
}
