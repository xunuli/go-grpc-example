package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "hello_server/proto"
	"log"
	"net"
)

// 自定义服务类
type server struct {
	pb.UnimplementedHelloServer
}

// 服务提供的方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("I am called!")
	return &pb.HelloReply{Reply: "Hello " + in.GetName()}, nil
}

// 启用服务
func main() {
	//监听端口
	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
		return
	}

	//创建一个grpc服务
	grpcServer := grpc.NewServer()
	//将自定义服务注册进去
	pb.RegisterHelloServer(grpcServer, &server{})
	//启动服务
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("server failed: %v", err)
		return
	}
	fmt.Println("I am ready!")
}
