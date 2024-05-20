package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	pb "md_client/proto"
)

func main() {
	//连接
	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close err: %v", err)
		}
	}(conn)
	//创建客户端
	client := pb.NewHelloClient(conn)
	//调用不同的服务
	runSayHello(client, "xuji")
	runSayHelloStream(client, "xuji")
}

// 普通
func runSayHello(client pb.HelloClient, name string) {
	fmt.Println("--- UnarySayHello client---")
	//创建metadata
	md := metadata.Pairs(
		"token", "app-test-xuji",
		"request_id", "1234567",
	)
	//基于metadata创建context
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	//RPC调用
	var header, trailer metadata.MD
	repo, err := client.SayHello(
		ctx,
		&pb.HelloRequest{Name: name},
		grpc.Header(&header),   //接受服务端发来的header
		grpc.Trailer(&trailer), //接受服务端发来的trailer
	)
	if err != nil {
		log.Fatalf("failed to call SayHello: %v", err)
		return
	}
	//从header中取出location
	if t, ok := header["location"]; ok {
		fmt.Printf("location from header:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Printf("location expected but doesn't exist in header")
		return
	}
	//获取响应数据
	fmt.Printf("got response: %s\n", repo.GetReply())
	//从trailer中获取timestamp
	if t, ok := trailer["timestamp"]; ok {
		fmt.Printf("timestamp from trailer:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Printf("timestamp expected but doesn't exist in trailer")
	}
}

// 流式
func runSayHelloStream(client pb.HelloClient, name string) {
	fmt.Println()
	fmt.Println("--- StreamSayHello client---")
	// 创建metadata和context.
	md := metadata.Pairs("token", "app-test-xuji")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// 使用带有metadata的context执行RPC调用.
	stream, err := client.SayHelloStream(ctx)
	if err != nil {
		log.Fatalf("failed to call BidiHello: %v\n", err)
	}

	go func() {
		// 当header到达时读取header.
		header, err := stream.Header()
		if err != nil {
			log.Fatalf("failed to get header from stream: %v", err)
		}
		// 从返回响应的header中读取数据.
		if l, ok := header["location"]; ok {
			fmt.Printf("location from header:\n")
			for i, e := range l {
				fmt.Printf(" %d. %s\n", i, e)
			}
		} else {
			log.Println("location expected but doesn't exist in header")
			return
		}
		// 发送所有的请求数据到server.
		for i := 0; i < 5; i++ {
			if err := stream.Send(&pb.HelloRequest{Name: name}); err != nil {
				log.Fatalf("failed to send streaming: %v\n", err)
			}
		}
		stream.CloseSend()
	}()

	// 读取所有的响应.
	var rpcStatus error
	fmt.Printf("got response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf(" - %s\n", r.Reply)
	}
	if rpcStatus != io.EOF {
		log.Printf("failed to finish server streaming: %v", rpcStatus)
		return
	}

	// 当RPC结束时读取trailer
	trailer := stream.Trailer()
	// 从返回响应的trailer中读取metadata.
	if t, ok := trailer["timestamp"]; ok {
		fmt.Printf("timestamp from trailer:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Printf("timestamp expected but doesn't exist in trailer")
	}

}
