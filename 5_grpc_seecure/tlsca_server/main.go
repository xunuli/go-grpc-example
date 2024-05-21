package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	pb "tlsca_server/proto"
)

// 自定义服务
type server struct {
	pb.UnimplementedHelloServer
}

// 重写SayHello方法
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	//获取元数据信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "未传输元数据")
	}
	//对比元数据信息
	var appid string
	var appkey string
	if v, ok := md["appid"]; ok {
		appid = v[0]
	}
	if v, ok := md["appkey"]; ok {
		appkey = v[0]
	}
	if appid != "xuji" || appkey != "123123" {
		return nil, errors.New("Token 不正确")
	}

	return &pb.HelloReply{ResponseMsg: "hello" + req.RequestName}, nil

}

func main() {
	//安全凭证
	creds, _ := credentials.NewServerTLSFromFile("C:\\Users\\xuji\\Desktop\\goProject\\go-grpc-example\\5_grpc_seecure\\key\\test.pem",
		"C:\\Users\\xuji\\Desktop\\goProject\\go-grpc-example\\5_grpc_seecure\\key\\test.key")

	//开启端口
	listen, err := net.Listen("tcp", ":9092")
	//创建一个grpc服务，通过凭证实现tls安全连接
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	//注册服务
	pb.RegisterHelloServer(grpcServer, &server{})
	//启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}
