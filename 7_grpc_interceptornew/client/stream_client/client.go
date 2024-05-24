package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	spb "grpc_interceptor/proto/stream"
	"io"
	"log"
	"time"
)

// 流式拦截器
// streamInterceptor（流拦截器） 一个简单的 stream interceptor 示例。
func streamclientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newStreamClient(s), nil
}

type streamClient struct {
	grpc.ClientStream
}

func newStreamClient(c grpc.ClientStream) grpc.ClientStream {
	return &streamClient{c}
}

// RecvMsg从流中接收消息
func (e *streamClient) RecvMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := e.ClientStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

// RecvMsg从流中接收消息
func (e *streamClient) SendMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := e.ClientStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}

func main() {
	//建立连接
	conn, err := grpc.NewClient("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(streamclientInterceptor),
	)
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("colse conn failed: %v", err)
			return
		}
	}(conn)
	//注册客户端
	client := spb.NewStreamClient(conn)
	//调用
	err = printList(client, &spb.StreamRequest{Pt: &spb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 1111}})
	err = printRecord(client, &spb.StreamRequest{Pt: &spb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2222}})
	err = printRoute(client, &spb.StreamRequest{Pt: &spb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 3333}})
}

// 服务端流式，接收服务端流式信息
func printList(client spb.StreamClient, resq *spb.StreamRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.List(ctx, resq)
	if err != nil {
		return err
	}
	//for循环获取服务端的信息
	for {
		//通过Recv()不断获取服务端推送的信息
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pt.name:%s, pt.value:%d\n", resp.Pt.Name, resp.Pt.Value)
	}
	return nil
}

// 客户端流式，发送流式信息
func printRecord(client spb.StreamClient, resq *spb.StreamRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//获取stream
	stream, err := client.Record(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		//通过send不断推送消息
		err := stream.Send(resq)
		if err != nil {
			return err
		}
	}

	//发送完成后通过stream.CloseAndRecv()关闭stream并接收服务端返回结果
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	return nil
}

// 双向流式，发送流式信息，接收流式信息
func printRoute(client spb.StreamClient, resq *spb.StreamRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	done := make(chan bool)
	stream, err := client.Route(ctx)
	if err != nil {
		return err
	}
	//开两个goroutine用于接收和发送
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("server closed")
				done <- true
				break
			}
			if err != nil {
				continue
			}
			log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			err := stream.Send(resq)
			if err != nil {
				log.Printf("send error:%v\n", err)
			}
			time.Sleep(time.Second)
		}
		//发送完毕关闭
		err = stream.CloseSend()
		if err != nil {
			log.Printf("Send error:%v\n", err)
			return
		}
	}()
	//阻塞至读取完成
	<-done
	return nil
}
