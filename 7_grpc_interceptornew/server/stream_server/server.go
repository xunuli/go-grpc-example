package main

import (
	"fmt"
	"google.golang.org/grpc"
	spb "grpc_interceptor/proto/stream"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type streamserver struct {
	spb.UnimplementedStreamServer
}

// 服务端流式，Server是stream，client是普通rpc
// 客户端发送一次请求，服务端通过流式响应多次发送数据
func (s *streamserver) List(req *spb.StreamRequest, stream spb.Stream_ListServer) error {
	//具体返回多少个response根据业务逻辑调整
	for i := 0; i <= 6; i++ {
		//通过send不断发送数据
		err := stream.Send(&spb.StreamResponse{
			Pt: &spb.StreamPoint{
				Name:  req.Pt.Name,
				Value: req.Pt.Value + int32(i),
			},
		})

		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

// 客户端流式RPC
// 客户端通过流式多次发送RPC请求给服务端，服务端返回一次响应
func (s *streamserver) Record(stream spb.Stream_RecordServer) error {
	//for循环接收客户端发送的消息
	for {
		//通过rec不断获取客户端send的请求
		resq, err := stream.Recv()
		//err == io.EOF表示获取到全部数据
		if err == io.EOF {
			//在客户端发送完毕后服务端即可返回响应
			return stream.SendAndClose(&spb.StreamResponse{
				Pt: &spb.StreamPoint{
					Name:  "gRPC stream Server: Record",
					Value: 1,
				},
			})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %s, pt.value: %d\n", resq.Pt.Name, resq.Pt.Value)
		time.Sleep(time.Second)
	}
	return nil
}

// 双向流式,由客户端发起流式RPC方法，服务端以同样的流式RPC方法返回响应
func (s *streamserver) Route(stream spb.Stream_RouteServer) error {
	var (
		wg    sync.WaitGroup
		msgCh = make(chan *spb.StreamPoint)
	)
	//发送响应
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for v := range msgCh {
			//发送数据
			err := stream.Send(&spb.StreamResponse{
				Pt: &spb.StreamPoint{
					Name:  v.GetName(),
					Value: int32(i),
				},
			})
			if err != nil {
				fmt.Println("send err:", err)
				continue
			}
			i++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		//接收数据
		for {
			resq, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("recv err: %v", err)
			}
			log.Printf("sream.Recv pt.name: %s, pt.value: %d", resq.Pt.Name, resq.Pt.Value)
			msgCh <- &spb.StreamPoint{
				Name: "grpc stream server: route",
			}
		}
		close(msgCh)
	}()
	wg.Wait() //等待任务结束
	return nil
}

func main() {
	//监听端口
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
		return
	}
	//创建grpc服务
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(streamServerInterceptor))
	//注册服务
	spb.RegisterStreamServer(grpcServer, &streamserver{})
	//启动服务
	err = grpcServer.Serve(lis)
}

// streamServerInterceptor 流式拦截器
func streamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 包装 grpc.ServerStream
	//重载RecvMsg和SendMsg方法
	wrapper := newStreamServer(ss)
	//调用服务方法
	err := handler(srv, wrapper)
	if err != nil {
		log.Fatalf("RPC failed with error %v", err)
	}
	return err
}

// 包装grpc.ServerStream
type streamServer struct {
	grpc.ServerStream
}

// 新建streamServer
func newStreamServer(ss grpc.ServerStream) grpc.ServerStream {
	return &streamServer{ss}
}

// 重写接收方法Recvmsg
func (s *streamServer) RecvMsg(m interface{}) error {
	//在这里，我们可以对接收到的消息执行额外的逻辑
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

// 重写发送方法Sendmsg
func (s *streamServer) SendMsg(m interface{}) error {
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := s.ServerStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}
