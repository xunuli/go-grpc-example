package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime" // 注意v2版本
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc_gateway/proto/helloworld"
	"log"
	"net"
	"net/http"
	"testing"
)

type server1 struct {
	pb.UnimplementedHelloServer
}

func (s *server1) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Reply: "Hello world! " + req.GetName()}, nil
}

func Test(t *testing.T) {
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
		return
	}
	//创建一个grpc服务
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServer(grpcServer, &server{})
	// 启动gRPC Server
	log.Println("Serving gRPC on localhost:9092")
	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	//创建一个连接到我们刚刚启动的gRPC服务器的客户端连接
	//grpc-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	//创建用来处理http请求的HTTP多路复用器，可以理解为一个理由，负责将传入的HTTP请求转发到相应的gRPC方法
	gwmux := runtime.NewServeMux()
	//注册路由
	err = pb.RegisterHelloHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	gwServer := &http.Server{
		Addr:    ":9090",
		Handler: gwmux,
	}

	// 9090端口提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:9090")
	log.Fatalln(gwServer.ListenAndServe())
}
