package handler

import (
	"net/http"

	"api/internal/svc"
	"api/response"
)

func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, r, map[string]string{"status": "ok"})
	}
}
