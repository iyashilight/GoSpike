package middleware

import (
	"net/http"

	"api/errorx"
	"api/response"

	"github.com/zeromicro/go-zero/core/logx"
)

func Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logx.Errorf("panic recovered: %v", rec)
				response.Error(w, r, errorx.NewCodeError(errorx.ErrUnknown))
			}
		}()
		next(w, r)
	}
}