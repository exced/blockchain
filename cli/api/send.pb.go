// Code generated by protoc-gen-go. DO NOT EDIT.
// source: send.proto

/*
Package send is a generated protocol buffer package.

It is generated from these files:
	send.proto

It has these top-level messages:
	Transaction
	TransactionMessage
	Response
*/
package send

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Transaction struct {
	To     string `protobuf:"bytes,1,opt,name=to" json:"to,omitempty"`
	Amount int64  `protobuf:"varint,2,opt,name=amount" json:"amount,omitempty"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Transaction) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Transaction) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type TransactionMessage struct {
	Signature    []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Hash         []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Rsapublickey []byte `protobuf:"bytes,3,opt,name=rsapublickey,proto3" json:"rsapublickey,omitempty"`
}

func (m *TransactionMessage) Reset()                    { *m = TransactionMessage{} }
func (m *TransactionMessage) String() string            { return proto.CompactTextString(m) }
func (*TransactionMessage) ProtoMessage()               {}
func (*TransactionMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TransactionMessage) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *TransactionMessage) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *TransactionMessage) GetRsapublickey() []byte {
	if m != nil {
		return m.Rsapublickey
	}
	return nil
}

type Response struct {
	Success bool   `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Msg     string `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Response) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Transaction)(nil), "send.Transaction")
	proto.RegisterType((*TransactionMessage)(nil), "send.TransactionMessage")
	proto.RegisterType((*Response)(nil), "send.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Peer service

type PeerClient interface {
	Send(ctx context.Context, in *TransactionMessage, opts ...grpc.CallOption) (*Response, error)
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

func (c *peerClient) Send(ctx context.Context, in *TransactionMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/send.Peer/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Peer service

type PeerServer interface {
	Send(context.Context, *TransactionMessage) (*Response, error)
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

func _Peer_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransactionMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/send.Peer/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).Send(ctx, req.(*TransactionMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "send.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Peer_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "send.proto",
}

func init() { proto.RegisterFile("send.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xcd, 0x4a, 0x03, 0x31,
	0x10, 0xc7, 0xdd, 0xdd, 0x50, 0xbb, 0xe3, 0x52, 0x64, 0x0e, 0x12, 0xc4, 0x43, 0xc9, 0xa9, 0xa7,
	0x1e, 0x2a, 0x7a, 0xf0, 0x1d, 0x04, 0x89, 0xbe, 0x40, 0x9a, 0x0e, 0xdb, 0x55, 0x9b, 0x2c, 0x99,
	0xe4, 0xe0, 0xdb, 0xcb, 0x8e, 0x2e, 0x2a, 0xde, 0xfe, 0x1f, 0xe4, 0x97, 0x99, 0x01, 0x60, 0x0a,
	0x87, 0xed, 0x98, 0x62, 0x8e, 0xa8, 0x26, 0x6d, 0xee, 0xe0, 0xe2, 0x25, 0xb9, 0xc0, 0xce, 0xe7,
	0x21, 0x06, 0x5c, 0x41, 0x9d, 0xa3, 0xae, 0xd6, 0xd5, 0xa6, 0xb5, 0x75, 0x8e, 0x78, 0x05, 0x0b,
	0x77, 0x8a, 0x25, 0x64, 0x5d, 0xaf, 0xab, 0x4d, 0x63, 0xbf, 0x9d, 0x79, 0x05, 0xfc, 0xf5, 0xec,
	0x91, 0x98, 0x5d, 0x4f, 0x78, 0x03, 0x2d, 0x0f, 0x7d, 0x70, 0xb9, 0x24, 0x12, 0x48, 0x67, 0x7f,
	0x02, 0x44, 0x50, 0x47, 0xc7, 0x47, 0x21, 0x75, 0x56, 0x34, 0x1a, 0xe8, 0x12, 0xbb, 0xb1, 0xec,
	0xdf, 0x07, 0xff, 0x46, 0x1f, 0xba, 0x91, 0xee, 0x4f, 0x66, 0xee, 0x61, 0x69, 0x89, 0xc7, 0x18,
	0x98, 0x50, 0xc3, 0x39, 0x17, 0xef, 0x89, 0x59, 0xf8, 0x4b, 0x3b, 0x5b, 0xbc, 0x84, 0xe6, 0xc4,
	0xbd, 0xc0, 0x5b, 0x3b, 0xc9, 0xdd, 0x03, 0xa8, 0x27, 0xa2, 0x84, 0x3b, 0x50, 0xcf, 0x14, 0x0e,
	0xa8, 0xb7, 0xb2, 0xfd, 0xff, 0xb9, 0xaf, 0x57, 0x5f, 0xcd, 0xfc, 0x8b, 0x39, 0xdb, 0x2f, 0xe4,
	0x46, 0xb7, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x27, 0x00, 0xe2, 0x61, 0x31, 0x01, 0x00, 0x00,
}