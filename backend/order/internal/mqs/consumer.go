package mqs

import (
	"context"
	"encoding/json"
	"order/internal/svc"
	"product/product"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type purchaseMsg struct {
	UserId    int64 `json:"user_id"`
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func StartConsumer(ctx context.Context, svcCtx *svc.ServiceContext) {
	ch, err := svcCtx.GetChannel()
	if err != nil {
		logx.Errorf("创建 RabbitMQ channel 失败: %v", err)
		return
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"order_flash_purchase",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logx.Errorf("声明队列失败: %v", err)
		return
	}
	msgs, err := ch.Consume(
		"order_flash_purchase",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logx.Errorf("注册消费者失败:%v", err)
		return
	}
	logx.Info("Order 消费者已启动,等待秒杀消息...")

	for {
		select {
		case <-ctx.Done():
			logx.Info("消费者退出")
			return
		case msg := <-msgs:
			var p purchaseMsg
			if err := json.Unmarshal(msg.Body, &p); err != nil {
				logx.Errorf("消息解析失败: %v", err)
				msg.Nack(false,true)
				continue
			}
			processPurchase(ctx, svcCtx, &p, msg)
		}
	}
}

func processPurchase(ctx context.Context, svcCtx *svc.ServiceContext, p *purchaseMsg, delivery amqp.Delivery) {
	// 1. 先查商品信息（只读，不修改数据）
	productResp, err := svcCtx.ProductRpc.GetProduct(ctx, &product.GetProductRequest{
		Id: p.ProductId,
	})
	if err != nil {
		logx.Errorf("获取商品信息失败: %v", err)
		return
	}

	// 2. 先创建订单（status = pending），再扣 MySQL 库存。
	//    避免库存扣了但订单没创建成功的不一致问题。
	var orderId int64
	err = svcCtx.SqlConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		result, err := s.ExecCtx(ctx,
			"insert into `orders` (`user_id`, `total_price`, `status`) values (?, ?, ?)",
			p.UserId, productResp.Price*float64(p.Quantity), "pending",
		)
		if err != nil {
			return err
		}
		orderId, _ = result.LastInsertId()
		_, err = s.ExecCtx(ctx,
			"insert into `order_item` (`order_id`, `product_id`, `product_name`, `price`, `quantity`) values (?, ?, ?, ?, ?)",
			orderId, p.ProductId, productResp.Name, productResp.Price, p.Quantity,
		)
		return err
	})
	if err != nil {
		logx.Errorf("创建订单失败: %v", err)
		delivery.Nack(false, true)
		return
	}

	// 3. 扣 MySQL 库存（WHERE stock >= ? 防止超卖）
	_, err = svcCtx.ProductRpc.DecrStock(ctx, &product.DecrStockRequest{
		Id:       p.ProductId,
		Quantity: p.Quantity,
	})
	if err != nil {
		logx.Errorf("扣 MySQL 库存失败，订单 %d 标记为 failed: %v", orderId, err)
		if updateErr := svcCtx.SqlConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
			_, execErr := s.ExecCtx(ctx, "update `orders` set `status` = ? where `id` = ?", "failed", orderId)
			return execErr
		}); updateErr != nil {
			logx.Errorf("更新订单 %d 状态为 failed 失败: %v", orderId, updateErr)
		}
		delivery.Nack(false, true)
		return
	}

	// 4. 更新订单状态为 paid
	if updateErr := svcCtx.SqlConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		_, execErr := s.ExecCtx(ctx, "update `orders` set `status` = ? where `id` = ?", "paid", orderId)
		return execErr
	}); updateErr != nil {
		logx.Errorf("更新订单 %d 状态为 paid 失败: %v", orderId, updateErr)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
	logx.Infof("订单创建成功: order_id=%d, user_id=%d, product_id=%d", orderId, p.UserId, p.ProductId)
}
