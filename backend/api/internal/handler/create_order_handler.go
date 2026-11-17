package handler

import (
	"net/http"

	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"
	"api/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, r, err)
			return
		}

		userId, err := GetUserIdFromCtx(r.Context())
		if err != nil {
			response.Error(w, r, err)
			return
		}
		req.UserId = userId

		l := logic.NewCreateOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateOrder(&req)
		if err != nil {
			response.Error(w, r, err)
			return
		}
		response.Success(w, r, resp)
	}
}
