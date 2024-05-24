package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	hpb "grpc_interceptor/proto/hello"
	"log"
	"net"
	"time"
)

type server struct {
	hpb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, req *hpb.HelloRequest) (*hpb.HelloReply, error) {
	format := time.Now().Format("2017-01-17 17:30:17")
	return &hpb.HelloReply{Reply: "Hello " + req.GetName() + "---" + format}, nil
}

func main() {
	//监听端口
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("listen failed: %V", err)
		return
	}
	//创建grpc服务，添加拦截器
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	//注册服务
	hpb.RegisterHelloServer(grpcServer, &server{})
	//启动服务
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server failed: %v", err)
		return
	}
}

// unaryInterceptor 一元拦截器：记录请求日志
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//预处理
	start := time.Now()
	//从传入上下文获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("couldn't parse incoming context metadata")
	}
	//检索客户端操作系统，如果它不存在，则此值位空
	os := md.Get("client-os")

	//RPC调用真正的执行罗逻辑
	m, err := handler(ctx, req)

	//后处理
	end := time.Now()
	log.Printf("client-OS: '%v' start time: %s, end time: %s, err: %v\n", os, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, nil
}
