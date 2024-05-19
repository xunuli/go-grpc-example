package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "hello_client/proto"
	"log"
	"time"
)

func main() {
	//连接到服务端,此处表示禁用安全传输
	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	//基于此链接，创建一个客户端，用于调用服务端
	client := pb.NewHelloClient(conn)
	//初始化一个超时取消context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("I am ready!")
	repo, err := client.SayHello(ctx, &pb.HelloRequest{Name: "xuji"})
	if err != nil {
		log.Fatalf("called failed: %v", err)
		return
	}
	fmt.Println(repo.GetReply())
}
