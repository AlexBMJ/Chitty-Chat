// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package chatservice_grpc

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

// ChittychatClient is the client API for Chittychat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittychatClient interface {
	Join(ctx context.Context, in *Message, opts ...grpc.CallOption) (Chittychat_JoinClient, error)
	Leave(ctx context.Context, in *IdMessage, opts ...grpc.CallOption) (*Empty, error)
	Publish(ctx context.Context, in *TextMessage, opts ...grpc.CallOption) (*Empty, error)
	Broadcast(ctx context.Context, in *TextMessage, opts ...grpc.CallOption) (*Empty, error)
}

type chittychatClient struct {
	cc grpc.ClientConnInterface
}

func NewChittychatClient(cc grpc.ClientConnInterface) ChittychatClient {
	return &chittychatClient{cc}
}

func (c *chittychatClient) Join(ctx context.Context, in *Message, opts ...grpc.CallOption) (Chittychat_JoinClient, error) {
	stream, err := c.cc.NewStream(ctx, &Chittychat_ServiceDesc.Streams[0], "/chatservice.Chittychat/Join", opts...)
	if err != nil {
		return nil, err
	}
	x := &chittychatJoinClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Chittychat_JoinClient interface {
	Recv() (*TextMessage, error)
	grpc.ClientStream
}

type chittychatJoinClient struct {
	grpc.ClientStream
}

func (x *chittychatJoinClient) Recv() (*TextMessage, error) {
	m := new(TextMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chittychatClient) Leave(ctx context.Context, in *IdMessage, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/chatservice.Chittychat/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittychatClient) Publish(ctx context.Context, in *TextMessage, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/chatservice.Chittychat/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittychatClient) Broadcast(ctx context.Context, in *TextMessage, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/chatservice.Chittychat/Broadcast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittychatServer is the server API for Chittychat service.
// All implementations must embed UnimplementedChittychatServer
// for forward compatibility
type ChittychatServer interface {
	Join(*Message, Chittychat_JoinServer) error
	Leave(context.Context, *IdMessage) (*Empty, error)
	Publish(context.Context, *TextMessage) (*Empty, error)
	Broadcast(context.Context, *TextMessage) (*Empty, error)
	mustEmbedUnimplementedChittychatServer()
}

// UnimplementedChittychatServer must be embedded to have forward compatible implementations.
type UnimplementedChittychatServer struct {
}

func (UnimplementedChittychatServer) Join(*Message, Chittychat_JoinServer) error {
	return status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedChittychatServer) Leave(context.Context, *IdMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedChittychatServer) Publish(context.Context, *TextMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedChittychatServer) Broadcast(context.Context, *TextMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}
func (UnimplementedChittychatServer) mustEmbedUnimplementedChittychatServer() {}

// UnsafeChittychatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittychatServer will
// result in compilation errors.
type UnsafeChittychatServer interface {
	mustEmbedUnimplementedChittychatServer()
}

func RegisterChittychatServer(s grpc.ServiceRegistrar, srv ChittychatServer) {
	s.RegisterService(&Chittychat_ServiceDesc, srv)
}

func _Chittychat_Join_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Message)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChittychatServer).Join(m, &chittychatJoinServer{stream})
}

type Chittychat_JoinServer interface {
	Send(*TextMessage) error
	grpc.ServerStream
}

type chittychatJoinServer struct {
	grpc.ServerStream
}

func (x *chittychatJoinServer) Send(m *TextMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _Chittychat_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittychatServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.Chittychat/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittychatServer).Leave(ctx, req.(*IdMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chittychat_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittychatServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.Chittychat/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittychatServer).Publish(ctx, req.(*TextMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chittychat_Broadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittychatServer).Broadcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chatservice.Chittychat/Broadcast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittychatServer).Broadcast(ctx, req.(*TextMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Chittychat_ServiceDesc is the grpc.ServiceDesc for Chittychat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chittychat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatservice.Chittychat",
	HandlerType: (*ChittychatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Leave",
			Handler:    _Chittychat_Leave_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _Chittychat_Publish_Handler,
		},
		{
			MethodName: "Broadcast",
			Handler:    _Chittychat_Broadcast_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Join",
			Handler:       _Chittychat_Join_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chatservice/chatservice.proto",
}