package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	pb "stream_client/proto"
	"strings"
	"time"
)

func main() {
	//创建连接，无安全
	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		return
	}
	//结束时关闭连接
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	//创建调用服务的客户端
	client := pb.NewHelloClient(conn)

	//调用不同的服务
	runSayHellostreamServer(client)
	runSayHellostreamClient(client)
	runSayHellostream(client)
}

// 客户端调用服务端流式
func runSayHellostreamServer(client pb.HelloClient) {
	//创建一个超时取消ctx，用于超时取消goroutine
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//返回时，延时调用取消函数
	defer cancel()
	//grpc客户端，生成一个服务函数的客户端stream，用于接受数据
	stram, err := client.SayHellostreamServer(ctx, &pb.HelloRequest{Name: "xuji"})
	if err != nil {
		log.Fatalf("client.SayHellostreamServer failed: %v", err)
	}
	for {
		//接受到服务端的流式数据，当收到EOF或错误时退出
		resp, err := stram.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.SayHellostreamServer recv failed: %v", err)
			return
		}
		//打印输出
		fmt.Println(resp.GetReply())
	}
}

// 客户端调用客户端流式
func runSayHellostreamClient(client pb.HelloClient) {
	//创建一个超时取消ctx，用于超时取消goroutine
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//返回时，延时调用取消函数
	defer cancel()
	//grpc客户端，生成一个服务函数的客户端stream，用于接受数据
	stream, err := client.SayHellostreamClient(ctx)
	if err != nil {
		log.Fatalf("client.SayHellostreamServer failed: %v", err)
	}

	names := []string{"吉吉国王 ", "xuji ", "老板"}
	for _, name := range names {
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("client.SayHellostreamServer send failed: %v", err)
			return
		}
	}
	//接受数据
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.SayHellostreamServer send failed: %v", err)
	}
	//打印输出
	fmt.Println(resp.GetReply())
}

// 双向流式
func runSayHellostream(client pb.HelloClient) {
	//创建一个超时取消ctx，用于超时取消goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	//返回时，延时调用取消函数
	defer cancel()
	//grpc客户端，生成一个服务函数的客户端stream，用于接受数据
	stream, err := client.SayHellostream(ctx)
	if err != nil {
		log.Fatalf("client.SayHellostreamServer failed: %v", err)
	}

	waitc := make(chan struct{})
	//开一个协程接受服务端的响应
	go func() {
		for {
			// 接收服务端返回的响应
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("c.BidiHello stream.Recv() failed, err: %v", err)
			}
			fmt.Printf("AI：%s\n", in.GetReply())
		}
	}()
	//从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin) //获取输入生成读对象
	for {
		cmd, _ := reader.ReadString('\n') //读到换行
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		//将获取到的数据发送至服务端
		err := stream.Send(&pb.HelloRequest{Name: cmd})
		if err != nil {
			log.Fatalf("c.BidiHello stream.Send(%v) failed: %v", cmd, err)
		}
	}
	//关闭发送端
	_ = stream.CloseSend()
	<-waitc
}
