package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc_nameresolver/name_server/proto"
	"net"
)

// 定义端口
var port = flag.Int("port", 8972, "服务端口")

// 定义服务
type server struct {
	pb.UnimplementedHelloServer
	Addr string
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	//拼接字符串
	reply := fmt.Sprintf("Hello %s. [from %s]", req.GetName(), s.Addr)
	return &pb.HelloReply{Reply: reply}, nil
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	//启动服务
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to listen, err:%v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServer(grpcServer, &server{Addr: addr})
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve,err:%v\n", err)
		return
	}
}
