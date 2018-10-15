// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/dispatcher/proto/dispatcher.proto

package proto

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

// The request message containing the user's name.
type CommandRequest struct {
	ReqID                string   `protobuf:"bytes,1,opt,name=ReqID,proto3" json:"ReqID,omitempty"`
	OP                   string   `protobuf:"bytes,2,opt,name=OP,proto3" json:"OP,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandRequest) Reset()         { *m = CommandRequest{} }
func (m *CommandRequest) String() string { return proto.CompactTextString(m) }
func (*CommandRequest) ProtoMessage()    {}
func (*CommandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dispatcher_6bba419722e7e866, []int{0}
}
func (m *CommandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandRequest.Unmarshal(m, b)
}
func (m *CommandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandRequest.Marshal(b, m, deterministic)
}
func (dst *CommandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandRequest.Merge(dst, src)
}
func (m *CommandRequest) XXX_Size() int {
	return xxx_messageInfo_CommandRequest.Size(m)
}
func (m *CommandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommandRequest proto.InternalMessageInfo

func (m *CommandRequest) GetReqID() string {
	if m != nil {
		return m.ReqID
	}
	return ""
}

func (m *CommandRequest) GetOP() string {
	if m != nil {
		return m.OP
	}
	return ""
}

func (m *CommandRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// The response message containing the greetings
type CommandReply struct {
	ReqID                string   `protobuf:"bytes,1,opt,name=ReqID,proto3" json:"ReqID,omitempty"`
	OP                   string   `protobuf:"bytes,2,opt,name=OP,proto3" json:"OP,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandReply) Reset()         { *m = CommandReply{} }
func (m *CommandReply) String() string { return proto.CompactTextString(m) }
func (*CommandReply) ProtoMessage()    {}
func (*CommandReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_dispatcher_6bba419722e7e866, []int{1}
}
func (m *CommandReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandReply.Unmarshal(m, b)
}
func (m *CommandReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandReply.Marshal(b, m, deterministic)
}
func (dst *CommandReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandReply.Merge(dst, src)
}
func (m *CommandReply) XXX_Size() int {
	return xxx_messageInfo_CommandReply.Size(m)
}
func (m *CommandReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandReply.DiscardUnknown(m)
}

var xxx_messageInfo_CommandReply proto.InternalMessageInfo

func (m *CommandReply) GetReqID() string {
	if m != nil {
		return m.ReqID
	}
	return ""
}

func (m *CommandReply) GetOP() string {
	if m != nil {
		return m.OP
	}
	return ""
}

func (m *CommandReply) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*CommandRequest)(nil), "proto.CommandRequest")
	proto.RegisterType((*CommandReply)(nil), "proto.CommandReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DispatcherClient is the client API for Dispatcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DispatcherClient interface {
	// handle a command
	HandleCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandReply, error)
}

type dispatcherClient struct {
	cc *grpc.ClientConn
}

func NewDispatcherClient(cc *grpc.ClientConn) DispatcherClient {
	return &dispatcherClient{cc}
}

func (c *dispatcherClient) HandleCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandReply, error) {
	out := new(CommandReply)
	err := c.cc.Invoke(ctx, "/proto.Dispatcher/HandleCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DispatcherServer is the server API for Dispatcher service.
type DispatcherServer interface {
	// handle a command
	HandleCommand(context.Context, *CommandRequest) (*CommandReply, error)
}

func RegisterDispatcherServer(s *grpc.Server, srv DispatcherServer) {
	s.RegisterService(&_Dispatcher_serviceDesc, srv)
}

func _Dispatcher_HandleCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatcherServer).HandleCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Dispatcher/HandleCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatcherServer).HandleCommand(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Dispatcher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Dispatcher",
	HandlerType: (*DispatcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleCommand",
			Handler:    _Dispatcher_HandleCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/dispatcher/proto/dispatcher.proto",
}

func init() {
	proto.RegisterFile("pkg/dispatcher/proto/dispatcher.proto", fileDescriptor_dispatcher_6bba419722e7e866)
}

var fileDescriptor_dispatcher_6bba419722e7e866 = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2d, 0xc8, 0x4e, 0xd7,
	0x4f, 0xc9, 0x2c, 0x2e, 0x48, 0x2c, 0x49, 0xce, 0x48, 0x2d, 0xd2, 0x2f, 0x28, 0xca, 0x2f, 0xc9,
	0x47, 0x12, 0xd0, 0x03, 0x0b, 0x08, 0xb1, 0x82, 0x29, 0x25, 0x2f, 0x2e, 0x3e, 0xe7, 0xfc, 0xdc,
	0xdc, 0xc4, 0xbc, 0x94, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x11, 0x2e, 0xd6, 0xa0,
	0xd4, 0x42, 0x4f, 0x17, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x08, 0x47, 0x88, 0x8f, 0x8b,
	0xc9, 0x3f, 0x40, 0x82, 0x09, 0x2c, 0xc4, 0xe4, 0x1f, 0x20, 0x24, 0xc4, 0xc5, 0xe2, 0x92, 0x58,
	0x92, 0x28, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x13, 0x04, 0x66, 0x2b, 0x79, 0x70, 0xf1, 0xc0, 0xcd,
	0x2a, 0xc8, 0xa9, 0x24, 0xdf, 0x24, 0x23, 0x6f, 0x2e, 0x2e, 0x17, 0xb8, 0x83, 0x85, 0x6c, 0xb9,
	0x78, 0x3d, 0x12, 0xf3, 0x52, 0x72, 0x52, 0xa1, 0xa6, 0x0b, 0x89, 0x42, 0xfc, 0xa0, 0x87, 0xea,
	0x72, 0x29, 0x61, 0x74, 0xe1, 0x82, 0x9c, 0x4a, 0x25, 0x86, 0x24, 0x36, 0xb0, 0xa8, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x7c, 0xa8, 0xe1, 0x8c, 0x19, 0x01, 0x00, 0x00,
}
