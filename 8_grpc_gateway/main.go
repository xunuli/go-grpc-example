package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc_gateway/proto/helloworld"
	"log"
	"net"
	"net/http"
	"strings"
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Reply: "Hello world! " + req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	//创建一个grpc server对象
	grpcServer := grpc.NewServer()
	//注册服务
	pb.RegisterHelloServer(grpcServer, &server{})

	//创建一个grpc路由，将http请求路由到grpc方法
	gwmux := runtime.NewServeMux()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//注册路由
	err = pb.RegisterHelloHandlerFromEndpoint(context.Background(), gwmux, "localhost:8091", opts)
	if err != nil {
		log.Fatalln("Failed to register gwmux:", err)
	}
	//创建一个http路由
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	//定义HTTP server配置
	gwServer := &http.Server{
		Addr:    "localhost:8091",
		Handler: grpcHandlerFunc(grpcServer, mux), //请求的统一入口
	}
	log.Println("Serving on http://127.0.0.1:8091")
	log.Fatalln(gwServer.Serve(lis)) // 启动HTTP服务
}

// grpcHandlerFunc 将gRPC请求和HTTP请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
