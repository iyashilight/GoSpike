package svc

import (
	"flashsale/internal/config"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *redis.Redis
	Rabbit    *amqp.Connection
	mu        sync.RWMutex
	channel   *amqp.Channel
	channelMu sync.Mutex
	notify    chan *amqp.Error
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.BizRedis),
	}
	svc.connect(c.RabbitMQ)
	return svc
}

func (s *ServiceContext) connect(cfg config.RabbitMQConf) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("RabbitMQ 连接失败,5秒后重试: %v", err)
		time.AfterFunc(5*time.Second, func() { s.connect(cfg) })
		return
	}

	s.mu.Lock()
	s.Rabbit = conn
	s.notify = conn.NotifyClose(make(chan *amqp.Error))
	s.mu.Unlock()

	// 连接成功后初始化 channel
	s.initChannel()

	log.Println("RabbitMQ 连接成功")

	go func() {
		err := <-s.notify
		log.Printf("RabbitMQ 连接断开 (%v)，正在重连...", err)
		s.mu.Lock()
		s.Rabbit = nil
		s.channel = nil
		s.mu.Unlock()
		s.connect(cfg)
	}()
}

func (s *ServiceContext) initChannel() {
	ch, err := s.Rabbit.Channel()
	if err != nil {
		log.Printf("创建 RabbitMQ channel 失败: %v", err)
		return
	}
	_, err = ch.QueueDeclare(
		"order_flash_purchase",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("声明队列失败: %v", err)
		ch.Close()
		return
	}
	s.mu.Lock()
	s.channel = ch
	s.mu.Unlock()
}

func (s *ServiceContext) PublishMessage(body []byte) error {
	s.channelMu.Lock()
	defer s.channelMu.Unlock()

	s.mu.RLock()
	ch := s.channel
	s.mu.RUnlock()

	if ch == nil || ch.IsClosed() {
		return fmt.Errorf("RabbitMQ channel 未就绪")
	}

	return ch.Publish("", "order_flash_purchase", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}