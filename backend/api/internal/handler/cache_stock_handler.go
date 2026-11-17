package handler

import (
	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"
	"api/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)
func CacheStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CacheStockRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, r, err)
			return
		}

		l := logic.NewCacheStockLogic(r.Context(),svcCtx)
		resp, err := l.CacheStock(&req)
		if err != nil {
			response.Error(w, r, err)
			return
		}
		response.Success(w, r, resp)
	}
}