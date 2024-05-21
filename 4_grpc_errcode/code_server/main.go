package main

import (
	pb "code_server/proto"
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
)

// 自定义服务hello
// 为服务设置访问限制，每隔name只能调用一次SayHello方法，超过此限制就返回一个请求错误
type server struct {
	pb.UnimplementedHelloServer
	mu    sync.Mutex     //锁
	count map[string]int //记录每一个name的请求次数
}

// 实现服务
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	//加锁
	s.mu.Lock()
	defer s.mu.Unlock()
	//记录用户的请求次数
	s.count[in.GetName()]++
	//超过一次就返回错误
	if s.count[in.GetName()] > 1 {
		//创建一个错误，添加描述信息
		st := status.New(codes.ResourceExhausted, "Request limite exceeded")
		//为错误添加详细的描述信息
		ds, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name: %s", in.GetName()),
					Description: "限制每一个用户调用一次",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}

	//正常返回响应
	reply := "Hello " + in.GetName()
	return &pb.HelloReply{Reply: reply}, nil
}

func main() {
	//监听端口
	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen, err: %v", err)
		return
	}
	//创建grpc服务
	grpcServer := grpc.NewServer()
	//注册服务，并将其初始化
	pb.RegisterHelloServer(grpcServer, &server{count: make(map[string]int)})
	//启动服务
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to server: %v", err)
		return
	}
}
