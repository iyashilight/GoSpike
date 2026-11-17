package middleware

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func RateLimit(rds *redis.Redis, next http.HandlerFunc) http.HandlerFunc {
	limiter := limit.NewPeriodLimit(1, 100, rds, "rate:flash")
	return func(w http.ResponseWriter, r *http.Request) {
		key := "global"
		if v := r.Context().Value("user_id"); v != nil {
			if id, ok := v.(float64); ok {
				key = strconv.FormatInt(int64(id), 10)
			}
		}

		code, err := limiter.Take(key)
		if err != nil {
			next(w, r)
			return
		}

		if code == limit.OverQuota {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"code":429,"msg":"请求过于频繁"}`))
			return
		}

		next(w, r)
	}
}
