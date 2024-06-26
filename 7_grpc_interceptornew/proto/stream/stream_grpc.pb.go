// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: proto/stream/stream.proto

package stream

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StreamClient is the client API for Stream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamClient interface {
	// 服务端流式RPC
	List(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (Stream_ListClient, error)
	// 客户端流式RPC
	Record(ctx context.Context, opts ...grpc.CallOption) (Stream_RecordClient, error)
	// 双向流式RPC
	Route(ctx context.Context, opts ...grpc.CallOption) (Stream_RouteClient, error)
}

type streamClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamClient(cc grpc.ClientConnInterface) StreamClient {
	return &streamClient{cc}
}

func (c *streamClient) List(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (Stream_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[0], "/proto.stream.Stream/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Stream_ListClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type streamListClient struct {
	grpc.ClientStream
}

func (x *streamListClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamClient) Record(ctx context.Context, opts ...grpc.CallOption) (Stream_RecordClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[1], "/proto.stream.Stream/Record", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamRecordClient{stream}
	return x, nil
}

type Stream_RecordClient interface {
	Send(*StreamRequest) error
	CloseAndRecv() (*StreamResponse, error)
	grpc.ClientStream
}

type streamRecordClient struct {
	grpc.ClientStream
}

func (x *streamRecordClient) Send(m *StreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamRecordClient) CloseAndRecv() (*StreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamClient) Route(ctx context.Context, opts ...grpc.CallOption) (Stream_RouteClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[2], "/proto.stream.Stream/Route", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamRouteClient{stream}
	return x, nil
}

type Stream_RouteClient interface {
	Send(*StreamRequest) error
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type streamRouteClient struct {
	grpc.ClientStream
}

func (x *streamRouteClient) Send(m *StreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamRouteClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamServer is the server API for Stream service.
// All implementations must embed UnimplementedStreamServer
// for forward compatibility
type StreamServer interface {
	// 服务端流式RPC
	List(*StreamRequest, Stream_ListServer) error
	// 客户端流式RPC
	Record(Stream_RecordServer) error
	// 双向流式RPC
	Route(Stream_RouteServer) error
	mustEmbedUnimplementedStreamServer()
}

// UnimplementedStreamServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServer struct {
}

func (UnimplementedStreamServer) List(*StreamRequest, Stream_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedStreamServer) Record(Stream_RecordServer) error {
	return status.Errorf(codes.Unimplemented, "method Record not implemented")
}
func (UnimplementedStreamServer) Route(Stream_RouteServer) error {
	return status.Errorf(codes.Unimplemented, "method Route not implemented")
}
func (UnimplementedStreamServer) mustEmbedUnimplementedStreamServer() {}

// UnsafeStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServer will
// result in compilation errors.
type UnsafeStreamServer interface {
	mustEmbedUnimplementedStreamServer()
}

func RegisterStreamServer(s grpc.ServiceRegistrar, srv StreamServer) {
	s.RegisterService(&Stream_ServiceDesc, srv)
}

func _Stream_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServer).List(m, &streamListServer{stream})
}

type Stream_ListServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type streamListServer struct {
	grpc.ServerStream
}

func (x *streamListServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Stream_Record_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServer).Record(&streamRecordServer{stream})
}

type Stream_RecordServer interface {
	SendAndClose(*StreamResponse) error
	Recv() (*StreamRequest, error)
	grpc.ServerStream
}

type streamRecordServer struct {
	grpc.ServerStream
}

func (x *streamRecordServer) SendAndClose(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamRecordServer) Recv() (*StreamRequest, error) {
	m := new(StreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Stream_Route_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServer).Route(&streamRouteServer{stream})
}

type Stream_RouteServer interface {
	Send(*StreamResponse) error
	Recv() (*StreamRequest, error)
	grpc.ServerStream
}

type streamRouteServer struct {
	grpc.ServerStream
}

func (x *streamRouteServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamRouteServer) Recv() (*StreamRequest, error) {
	m := new(StreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Stream_ServiceDesc is the grpc.ServiceDesc for Stream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.stream.Stream",
	HandlerType: (*StreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _Stream_List_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Record",
			Handler:       _Stream_Record_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Route",
			Handler:       _Stream_Route_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/stream/stream.proto",
}
