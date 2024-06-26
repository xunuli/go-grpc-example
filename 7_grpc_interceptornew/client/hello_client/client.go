package main

import (
	"context"
	"fmt"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	hpb "grpc_interceptor/proto/hello"
	"log"
	"runtime"
	"time"
)

// unaryInterceptor（客户端一元拦截器） 一个简单的 unary interceptor 示例。
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// pre-processing（预处理）
	fmt.Println("我是第一个拦截器")
	start := time.Now()
	cos := runtime.GOOS
	ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)
	//调用RPC方法，一般在拦截器中调用 invoker 能达到调用 RPC 方法的效果，当然底层也是 gRPC 在处理。
	err := invoker(ctx, method, req, reply, cc, opts...) // invoking RPC method
	// post-processing（后处理）
	end := time.Now()
	log.Printf("client-OS: '%v' start time: %s, end time: %s, err: %v", cos, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return err
}

// unaryInterceptorTwo（客户端一元拦截器） 第二个
func unaryInterceptorTow() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		//预处理
		fmt.Println("我是第二个拦截器")
		//调用当前方法
		err := invoker(ctx, method, req, reply, cc, opts...)

		//后处理
		fmt.Println("哈哈哈啊哈哈")
		return err
	}
}

func main() {
	//建立连接  链式拦截器

	conn, err := grpc.NewClient("localhost:9092",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithUnaryInterceptor(unaryInterceptor), //一元拦截器
		grpc.WithUnaryInterceptor(middleware.ChainUnaryClient(unaryInterceptor, unaryInterceptorTow())),
	)
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("conn close failed: %v", err)
			return
		}
	}(conn)

	//注册客户端
	client := hpb.NewHelloClient(conn)

	//调用函数
	runSayHello(client)
}

func runSayHello(client hpb.HelloClient) {
	//设置超时取消
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	for i := 0; i < 1; i++ {
		resp, err := client.SayHello(ctx, &hpb.HelloRequest{Name: "xuji"})
		if err != nil {
			log.Fatalf("call server failed: %v", err)
			return
		}
		fmt.Printf("Resp: %v\n", resp.GetReply())
	}
}
