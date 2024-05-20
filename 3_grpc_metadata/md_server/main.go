package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	pb "md_server/proto"
	"net"
	"strconv"
	"time"
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	//通过defer设置trailer，处理完数据后发送元数据trailer
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		grpc.SetTrailer(ctx, trailer)
	}()

	//从客户端读取请求上下文中的metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	//如果metadata存在
	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata:\n")
		//如果token中的信息不匹配，t是[]string类型，len(t)<1则说明没有信息，t[0]!="···"则说明不匹配
		if len(t) < 1 || t[0] != "app-test-xuji" {
			return nil, status.Errorf(codes.Unauthenticated, "认证失败！")
		}
	}

	//创建和发送header
	header := metadata.New(map[string]string{"location": "BeiJing"})
	grpc.SendHeader(ctx, header)
	fmt.Printf("request received: %v, say hello...\n", in)

	return &pb.HelloReply{Reply: "Hello " + in.GetName()}, nil

}

func (s *server) SayHelloStream(stream pb.Hello_SayHelloStreamServer) error {
	//通过defer设置trailer，处理完数据后发送元数据trailer，记录函数的返回时间
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		stream.SetTrailer(trailer)
	}()

	// 从client读取metadata.
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "BidirectionalStreamingSayHello: failed to get metadata")
	}
	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// 创建和发送header.
	header := metadata.New(map[string]string{"location": "X2Q"})
	stream.SendHeader(header)

	// 读取请求数据发送响应数据.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending reply\n", in)
		if err := stream.Send(&pb.HelloReply{Reply: in.Name}); err != nil {
			return err
		}
	}
}

func main() {
	listerner, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	//创建grpc服务
	grpcServer := grpc.NewServer()
	//注册
	pb.RegisterHelloServer(grpcServer, &server{})
	//启动grpc服务
	err = grpcServer.Serve(listerner)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
