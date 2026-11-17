package svc

import (
	"fmt"
	"log"
	"order/internal/config"
	"order/internal/model"
	"product/product_client"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	OrdersModel    model.OrdersModel
	OrderItemModel model.OrderItemModel
	ProductRpc     product_client.Product
	SqlConn        sqlx.SqlConn
	conn           *amqp.Connection
	mu             sync.RWMutex
	notify         chan *amqp.Error
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:         c,
		OrdersModel:    model.NewOrdersModel(sqlx.NewMysql(c.DataSource), c.Cache),
		OrderItemModel: model.NewOrderItemModel(sqlx.NewMysql(c.DataSource), c.Cache),
		ProductRpc:     product_client.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		SqlConn:        sqlx.NewMysql(c.DataSource),
	}
	svc.connect(c.RabbitMQ)
	return svc
}

func (s *ServiceContext) connect(cfg config.RabbitMQConf) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("RabbitMQ 连接失败，5秒后重试: %v", err)
		time.AfterFunc(5*time.Second, func() { s.connect(cfg) })
		return
	}

	s.mu.Lock()
	s.conn = conn
	s.notify = conn.NotifyClose(make(chan *amqp.Error))
	s.mu.Unlock()

	log.Println("RabbitMQ 连接成功")

	go func() {
		err := <-s.notify
		log.Printf("RabbitMQ 连接断开 (%v)，正在重连...", err)
		s.mu.Lock()
		s.conn = nil
		s.mu.Unlock()
		s.connect(cfg)
	}()
}

func (s *ServiceContext) GetChannel() (*amqp.Channel, error) {
	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()

	if conn == nil || conn.IsClosed() {
		return nil, fmt.Errorf("RabbitMQ 未连接")
	}
	return conn.Channel()
}
