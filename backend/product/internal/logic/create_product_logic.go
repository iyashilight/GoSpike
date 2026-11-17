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

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductLogic) CreateProduct(in *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	if in.Name == "" || in.Price <= 0 || in.Stock <= 0 {
		return nil, status.Error(codes.InvalidArgument, "参数无效")
	}
	result, err := l.svcCtx.ProductModel.Insert(l.ctx,&model.Product{
		Name: in.Name,
		Description: in.Description,
		Price: in.Price,
		Stock: int64(in.Stock),
	})
	if err != nil {
		return nil,status.Error(codes.Internal,"创建商品失败")
	}
	id, _ := result.LastInsertId()
	return &product.CreateProductResponse{Id: id}, nil
}
