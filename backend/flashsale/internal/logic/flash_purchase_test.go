package logic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestDecrStockScript_StockNotExists(t *testing.T) {
	rds, err := redis.NewRedis(redis.RedisConf{Host: "127.0.0.1:6380"})
	if err != nil {
		t.Skip("Redis 不可用，跳过测试")
	}
	ctx := context.Background()
	key := fmt.Sprintf("test:stock:%d", time.Now().UnixNano())

	result, err := rds.ScriptRunCtx(ctx, decrStockScript, []string{key}, 1)
	if err != nil {
		t.Fatalf("脚本执行失败: %v", err)
	}
	if result != int64(-1) {
		t.Fatalf("期望 -1（不存在），得到 %d", result)
	}
}

func TestDecrStockScript_StockNotEnough(t *testing.T) {
	rds, err := redis.NewRedis(redis.RedisConf{Host: "127.0.0.1:6380"})
	if err != nil {
		t.Skip("Redis 不可用，跳过测试")
	}
	ctx := context.Background()
	key := fmt.Sprintf("test:stock:%d", time.Now().UnixNano())

	// 先设置库存为 3，尝试扣 5
	_ = rds.SetCtx(ctx, key, "3")

	result, err := rds.ScriptRunCtx(ctx, decrStockScript, []string{key}, 5)
	if err != nil {
		t.Fatalf("脚本执行失败: %v", err)
	}
	if result != int64(0) {
		t.Fatalf("期望 0（库存不足），得到 %d", result)
	}

	// 验证库存没有被扣
	val, _ := rds.GetCtx(ctx, key)
	if val != "3" {
		t.Fatalf("库存不应被扣，期望 3，得到 %s", val)
	}
}

func TestDecrStockScript_Success(t *testing.T) {
	rds, err := redis.NewRedis(redis.RedisConf{Host: "127.0.0.1:6380"})
	if err != nil {
		t.Skip("Redis 不可用，跳过测试")
	}
	ctx := context.Background()
	key := fmt.Sprintf("test:stock:%d", time.Now().UnixNano())

	_ = rds.SetCtx(ctx, key, "10")

	result, err := rds.ScriptRunCtx(ctx, decrStockScript, []string{key}, 3)
	if err != nil {
		t.Fatalf("脚本执行失败: %v", err)
	}
	if result != int64(1) {
		t.Fatalf("期望 1（成功），得到 %d", result)
	}

	val, _ := rds.GetCtx(ctx, key)
	if val != "7" {
		t.Fatalf("库存应为 10-3=7，得到 %s", val)
	}
}
