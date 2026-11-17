package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"flashsale/flashsale"
	"flashsale/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KEYS[1] = stock key (stock:{product_id})
// KEYS[2] = activity key (activity:{product_id})
// ARGV[1] = quantity
//
// Returns:
//   -2 活动未开始
//   -3 活动已结束
//   -4 商品不存在或未预热
//    0 库存不足
//    1 扣库存成功
var decrStockScript = redis.NewScript(`
		local stock_key = KEYS[1]
		local activity_key = KEYS[2]
		local quantity = tonumber(ARGV[1])

		-- 检查活动是否存在
		local start_time = redis.call("HGET", activity_key, "start_time")
		if not start_time then
			return -4 -- 活动不存在
		end
		local end_time = redis.call("HGET", activity_key, "end_time")

		-- 使用 Redis 服务器时间校验收
		local redis_time = redis.call("TIME")
		local now = tonumber(redis_time[1])

		if now < tonumber(start_time) then
			return -2 -- 活动未开始
		end
		if now > tonumber(end_time) then
			return -3 -- 活动已结束
		end

		local stock = redis.call("GET", stock_key)
		if not stock then
			return -4 -- 商品不存在/未预热
		end
		if tonumber(stock) < quantity then
			return 0 -- 库存不足
		end
		redis.call("DECRBY", stock_key, quantity)
		return 1 -- 扣库存成功
	`)

type purchaseMsg struct {
	UserId    int64 `json:"user_id"`
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type FlashPurchaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlashPurchaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlashPurchaseLogic {
	return &FlashPurchaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FlashPurchaseLogic) FlashPurchase(in *flashsale.FlashPurchaseRequest) (*flashsale.FlashPurchaseResponse, error) {
	stockKey := fmt.Sprintf("stock:%d", in.ProductId)
	activityKey := fmt.Sprintf("activity:%d", in.ProductId)

	result, err := l.svcCtx.Redis.ScriptRunCtx(l.ctx, decrStockScript, []string{stockKey, activityKey}, in.Quantity)
	if err != nil {
		return nil, status.Error(codes.Internal, "秒杀失败")
	}
	switch result {
	case int64(-4):
		return nil, status.Error(codes.NotFound, "商品不存在或者未预热")
	case int64(-3):
		return nil, status.Error(codes.FailedPrecondition, "秒杀活动已结束")
	case int64(-2):
		return nil, status.Error(codes.FailedPrecondition, "秒杀活动尚未开始")
	case int64(0):
		return nil, status.Error(codes.ResourceExhausted, "库存不足")
	default:
		body, _ := json.Marshal(purchaseMsg{
			UserId:    in.UserId,
			ProductId: in.ProductId,
			Quantity:  in.Quantity,
		})
		if err := l.svcCtx.PublishMessage(body); err != nil {
			logx.Errorf("秒杀消息发送失败: %v", err)
			return nil, status.Error(codes.Internal, "消息发送失败")
		}
		return &flashsale.FlashPurchaseResponse{
			OrderId: 0,
			Status:  "pending",
		}, nil
	}
}
