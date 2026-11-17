package logic

import (
	"context"
	"time"

	"product/internal/svc"
	"product/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListProductsLogic) ListProducts(in *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	size := in.Size
	if size < 1 || size > 100 {
		size = 20
	}

	items, err := l.svcCtx.ProductModel.FindPage(l.ctx, page, size)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询商品列表失败")
	}
	total, err := l.svcCtx.ProductModel.Count(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询商品总数失败")
	}

	var pbItems []*product.GetProductResponse
	for _, item := range items {
		pbItems = append(pbItems, &product.GetProductResponse{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       int32(item.Stock),
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
		})
	}

	return &product.ListProductsResponse{
		Items: pbItems,
		Total: total,
	}, nil
}
