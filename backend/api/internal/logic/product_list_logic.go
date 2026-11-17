package logic

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductListLogic {
	return &ProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductListLogic) ProductList(req *types.ProductListRequest) (resp *types.ProductListResponse, err error) {
	rpcResp, err := l.svcCtx.ProductRpc.ListProducts(l.ctx, &product.ListProductsRequest{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	var items []types.GetProductResponse
	for _, item := range rpcResp.Items {
		items = append(items, types.GetProductResponse{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       item.Stock,
			CreatedAt:   item.CreatedAt,
		})
	}
	return &types.ProductListResponse{
		Items: items,
		Total: rpcResp.Total,
	}, nil
}
