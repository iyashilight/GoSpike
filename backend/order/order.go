package main

import (
	"context"
	"flag"
	"fmt"

	"order/internal/config"
	"order/internal/mqs"
	"order/internal/server"
	"order/internal/svc"
	"order/order"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order.RegisterOrderServer(grpcServer, server.NewOrderServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	go mqs.StartConsumer(ctx, svcCtx)
	s.Start()
}
