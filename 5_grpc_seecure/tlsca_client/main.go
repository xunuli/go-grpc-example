package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	pb "tlsca_client/proto"
)

// Token认证
// 基于grpc的PerRPCCredentials接口 定义了需要为每个 RPC（如 oauth2）附加安全信息的凭证的通用接口。
// type PerRPCCredentials interface {
// GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
// RequireTransportSecurity() bool
// }
// 可以采用这种方式，也可以采用之前的metadata的方式
type ClientTokenAuth struct{}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{ //元数据信息
		"appid":  "xuji",
		"appkey": "123123",
	}, nil
}

func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return true //开启tls
}

func main() {
	//TLS加密连接的凭证
	creds, _ := credentials.NewClientTLSFromFile("C:\\Users\\xuji\\Desktop\\goProject\\go-grpc-example\\5_grpc_seecure\\key\\test.pem", "*.localhost")
	//安全连接到server端
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithPerRPCCredentials(&ClientTokenAuth{}))
	conn, err := grpc.NewClient("localhost:9092", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("did not close: %v", err)
			return
		}
	}(conn)

	//创建客户端，建立连接
	client := pb.NewHelloClient(conn)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "xuji"})
	if err != nil {
		log.Fatalf("no resp: %v", err)
		return
	}
	fmt.Printf("resp: %v", resp.GetResponseMsg())
}
