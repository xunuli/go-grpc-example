package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	pb "stream_server/proto"
	"strings"
)

type server struct {
	pb.UnimplementedHelloServer
}

// 服务端流式
func (s *server) SayHellostreamServer(in *pb.HelloRequest, stream pb.Hello_SayHellostreamServerServer) error {
	//定义返回的响应
	words := []string{
		"你好 ",
		"hello ",
		"hhhhhh ",
		"hihihi ",
	}
	for _, word := range words {
		data := &pb.HelloReply{
			Reply: word + in.GetName(),
		}
		//将多个数据采用send方法发送出去
		err := stream.Send(data)
		if err != nil {
			log.Fatalf("send failed: %v", err)
			return err
		}
	}
	return nil
}

// 客户端流式
func (s *server) SayHellostreamClient(stream pb.Hello_SayHellostreamClientServer) error {
	//返回响应
	reply := "你好： "
	//接受客户端发来的流式数据
	for {
		resq, err := stream.Recv()
		//当读取到文件结束符
		if err == io.EOF {
			//最终统一回复
			return stream.SendAndClose(&pb.HelloReply{
				Reply: reply,
			})
		}
		if err != nil {
			log.Fatalf("")
		}
		//未读取完，拼接返回值
		reply = reply + resq.GetName()
	}
}

// 双向流式
func (s *server) SayHellostream(stream pb.Hello_SayHellostreamServer) error {
	//接受客户端的流式数据
	for {
		//接受流式数据
		resq, err := stream.Recv()
		if err == io.EOF {
			//读取到文件结束符，直接返回
			return nil
		}
		if err != nil {
			log.Fatalf("recv failed:%v", err)
			return err
		}
		//对接受到的数据，生成响应
		reply := magic(resq.GetName())
		//返回流式响应
		err = stream.Send(&pb.HelloReply{Reply: reply})
		if err != nil {
			return err
		}
	}
}

// 对输入requesu进行处理
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

// 主函数，创建服务端grpc服务
func main() {
	//监听端口号
	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
		return
	}
	fmt.Println("I am listing······")
	//创建一个grpc服务，无安全连接
	grpcServer := grpc.NewServer()
	//注册服务到grpc服务
	pb.RegisterHelloServer(grpcServer, &server{})
	//启动服务
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("server failed: %v", err)
		return
	}
	fmt.Println("I am Called!")
}
