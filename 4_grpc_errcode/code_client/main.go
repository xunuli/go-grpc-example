package main

import (
	pb "code_client/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

//grpc客户端，调用server端的SayHello方法

var name = flag.String("name", "xuji", "通过-name告诉server你是谁")

func main() {
	//解析命令参数
	flag.Parse()

	//创建连接
	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	//创建grpc客户端
	client := pb.NewHelloClient(conn)
	//调用服务端的方法
	runSayHello(client)

}

func runSayHello(client pb.HelloClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		//将err转为status
		s := status.Convert(err)
		for _, d := range s.Details() { //获取details
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				fmt.Printf("Quota failure: %s\n", info)
			default:
				fmt.Printf("Unexpected type: %s\n", info)
			}
		}
		fmt.Printf("server failed: %v", err)
		return
	}
	//打印响应
	log.Printf("resp: %v\n", resp.GetReply())
}
