// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ab.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("ab.proto", fileDescriptor_23048cda704d433a) }

var fileDescriptor_23048cda704d433a = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x48, 0x4c, 0xd2, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x92, 0xe2, 0xcb, 0xc8, 0x2f, 0x29, 0x2e, 0x29,
	0x4d, 0x4b, 0x83, 0x88, 0x19, 0xd9, 0x70, 0xf1, 0x3b, 0x96, 0xe4, 0xe7, 0x66, 0x26, 0x3b, 0x15,
	0xe5, 0x27, 0xa6, 0x24, 0x27, 0x16, 0x97, 0x08, 0x69, 0x72, 0x71, 0x22, 0x38, 0xdc, 0x7a, 0x05,
	0x49, 0x7a, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x52, 0xc8, 0x1c, 0x0d, 0x46, 0x03, 0xc6,
	0x24, 0x36, 0xb0, 0x21, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x26, 0xae, 0x12, 0x4b, 0x64,
	0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AtomicBroadcastClient is the client API for AtomicBroadcast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AtomicBroadcastClient interface {
	Broadcast(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcast_BroadcastClient, error)
}

type atomicBroadcastClient struct {
	cc *grpc.ClientConn
}

func NewAtomicBroadcastClient(cc *grpc.ClientConn) AtomicBroadcastClient {
	return &atomicBroadcastClient{cc}
}

func (c *atomicBroadcastClient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcast_BroadcastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AtomicBroadcast_serviceDesc.Streams[0], "/pb.AtomicBroadcast/Broadcast", opts...)
	if err != nil {
		return nil, err
	}
	x := &atomicBroadcastBroadcastClient{stream}
	return x, nil
}

type AtomicBroadcast_BroadcastClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type atomicBroadcastBroadcastClient struct {
	grpc.ClientStream
}

func (x *atomicBroadcastBroadcastClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *atomicBroadcastBroadcastClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AtomicBroadcastServer is the server API for AtomicBroadcast service.
type AtomicBroadcastServer interface {
	Broadcast(AtomicBroadcast_BroadcastServer) error
}

// UnimplementedAtomicBroadcastServer can be embedded to have forward compatible implementations.
type UnimplementedAtomicBroadcastServer struct {
}

func (*UnimplementedAtomicBroadcastServer) Broadcast(srv AtomicBroadcast_BroadcastServer) error {
	return status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}

func RegisterAtomicBroadcastServer(s *grpc.Server, srv AtomicBroadcastServer) {
	s.RegisterService(&_AtomicBroadcast_serviceDesc, srv)
}

func _AtomicBroadcast_Broadcast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AtomicBroadcastServer).Broadcast(&atomicBroadcastBroadcastServer{stream})
}

type AtomicBroadcast_BroadcastServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type atomicBroadcastBroadcastServer struct {
	grpc.ServerStream
}

func (x *atomicBroadcastBroadcastServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *atomicBroadcastBroadcastServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AtomicBroadcast_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AtomicBroadcast",
	HandlerType: (*AtomicBroadcastServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Broadcast",
			Handler:       _AtomicBroadcast_Broadcast_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "ab.proto",
}