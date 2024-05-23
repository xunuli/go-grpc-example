package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	pb "grpc_nameresolver/name_server/proto"
	"log"
	"time"
)

// 自定义name resolver

// 固定shceme（协议方案），和服务名称（网址）
const (
	myScheme      = "xuji"
	myServiceName = "resolver.xuji.com"
)

// 服务ip地址
var addrs = []string{"127.0.0.1:8972"}

// xujiResolver 自定义name resolver.需实现resolver接口
type xujiResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *xujiResolver) ResolveNow(o resolver.ResolveNowOptions) {
	// 直接从map中取出对于的addrList
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*xujiResolver) Close() {}

// xujiResolverBuilder 需实现 Builder 接口
type xujiResolverBuilder struct{}

func (*xujiResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &xujiResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myServiceName: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (*xujiResolverBuilder) Scheme() string { return myScheme }

// 注册reslver
func init() {
	// 注册 xujiResolverBuilder
	resolver.Register(&xujiResolverBuilder{})
}

// 调用服务
func runSayHello(ctx context.Context, client pb.HelloClient) {
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "xuji"})
	if err != nil {
		log.Fatalf("call err")
		return
	}
	fmt.Printf("resp: %v\n", resp.GetReply())
}

func main() {
	//建立连接
	conn, err := grpc.NewClient(
		"xuji:///resolver.xuji.com", //
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connet: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("did not close connet: %v", err)
			return
		}
	}(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	client := pb.NewHelloClient(conn)
	//调用服务
	for i := 0; i < 10; i++ {
		runSayHello(ctx, client)
	}
}
