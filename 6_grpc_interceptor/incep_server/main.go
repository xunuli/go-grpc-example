package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	pb "incep_server/proto"
	"io"
	"log"
	"net"
	"time"
)

var port = flag.Int("port", 50051, "the port to serve on")

// logger 简单打印日志
func logger(format string, a ...interface{}) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("unary echoing message %q\n", in.Message)
	return &pb.EchoResponse{Message: in.Message}, nil
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("server: error receiving from stream: %v\n", err)
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in.Message)
		err = stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			fmt.Printf("server: error send to stream: %v\n", err)
		}
	}
}

// unaryInterceptor 一元拦截器：记录请求日志
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//预处理
	start := time.Now()
	//调用RPC方法，
	m, err := handler(ctx, req)
	//后处理
	end := time.Now()
	// 记录请求参数 耗时 错误信息等数据
	logger("RPC: %s,req:%v start time: %s, end time: %s, err: %v", info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, err
}

// 分装流式服务端grpc.serverstream
type wrappedStream struct {
	grpc.ServerStream
}

// 初始化一个流式服务端
func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// 重写接收方法：实现调用前处理
func (w *wrappedStream) RecvMsg(m interface{}) error {
	logger("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

// 重写发送方法：实现调用后处理
func (w *wrappedStream) SendMsg(m interface{}) error {
	logger("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

// streamInterceptor 流式拦截器
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 包装 grpc.ServerStream 以替换 RecvMsg SendMsg这两个方法。
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		logger("RPC failed with error %v", err)
	}
	return err
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//creds, err := credentials.NewServerTLSFromFile(data.Path("x509/server.crt"), data.Path("x509/server.key"))
	//if err != nil {
	//	log.Fatalf("failed to create credentials: %v", err)
	//}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))

	pb.RegisterEchoServer(s, &server{})
	log.Println("Server gRPC on 0.0.0.0" + fmt.Sprintf(":%d", *port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
