package handler

import (
	"context"
	"fmt"
)

func GetUserIdFromCtx(ctx context.Context) (int64, error) {
	v := ctx.Value("user_id")
	if v == nil {
		return 0, fmt.Errorf("user_id not found in context")
	}
	// JWT claims 中的数值类型解码为 float64
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("user_id has unexpected type %T", v)
	}
	return int64(f), nil
}
