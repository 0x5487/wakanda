// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pkg/gateway/proto/job.proto

package proto

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type SendJobRequest struct {
	Jobs                 []*Job   `protobuf:"bytes,1,rep,name=Jobs,proto3" json:"Jobs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendJobRequest) Reset()         { *m = SendJobRequest{} }
func (m *SendJobRequest) String() string { return proto.CompactTextString(m) }
func (*SendJobRequest) ProtoMessage()    {}
func (*SendJobRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_66325226e622ec2a, []int{0}
}
func (m *SendJobRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendJobRequest.Unmarshal(m, b)
}
func (m *SendJobRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendJobRequest.Marshal(b, m, deterministic)
}
func (m *SendJobRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendJobRequest.Merge(m, src)
}
func (m *SendJobRequest) XXX_Size() int {
	return xxx_messageInfo_SendJobRequest.Size(m)
}
func (m *SendJobRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendJobRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendJobRequest proto.InternalMessageInfo

func (m *SendJobRequest) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

type Job struct {
	Type                 string   `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	TargetID             string   `protobuf:"bytes,2,opt,name=TargetID,proto3" json:"TargetID,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_66325226e622ec2a, []int{1}
}
func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (m *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(m, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Job) GetTargetID() string {
	if m != nil {
		return m.TargetID
	}
	return ""
}

func (m *Job) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type EmptyRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyRequest) Reset()         { *m = EmptyRequest{} }
func (m *EmptyRequest) String() string { return proto.CompactTextString(m) }
func (*EmptyRequest) ProtoMessage()    {}
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_66325226e622ec2a, []int{2}
}
func (m *EmptyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyRequest.Unmarshal(m, b)
}
func (m *EmptyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyRequest.Marshal(b, m, deterministic)
}
func (m *EmptyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyRequest.Merge(m, src)
}
func (m *EmptyRequest) XXX_Size() int {
	return xxx_messageInfo_EmptyRequest.Size(m)
}
func (m *EmptyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyRequest proto.InternalMessageInfo

type EmptyReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyReply) Reset()         { *m = EmptyReply{} }
func (m *EmptyReply) String() string { return proto.CompactTextString(m) }
func (*EmptyReply) ProtoMessage()    {}
func (*EmptyReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_66325226e622ec2a, []int{3}
}
func (m *EmptyReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyReply.Unmarshal(m, b)
}
func (m *EmptyReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyReply.Marshal(b, m, deterministic)
}
func (m *EmptyReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyReply.Merge(m, src)
}
func (m *EmptyReply) XXX_Size() int {
	return xxx_messageInfo_EmptyReply.Size(m)
}
func (m *EmptyReply) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyReply.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyReply proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SendJobRequest)(nil), "proto.SendJobRequest")
	proto.RegisterType((*Job)(nil), "proto.Job")
	proto.RegisterType((*EmptyRequest)(nil), "proto.EmptyRequest")
	proto.RegisterType((*EmptyReply)(nil), "proto.EmptyReply")
}

func init() { proto.RegisterFile("pkg/gateway/proto/job.proto", fileDescriptor_66325226e622ec2a) }

var fileDescriptor_66325226e622ec2a = []byte{
	// 209 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0xcd, 0x6a, 0x84, 0x30,
	0x14, 0x85, 0x9b, 0x6a, 0x8b, 0xbd, 0x15, 0xa1, 0x81, 0x42, 0xb0, 0x50, 0x24, 0x2b, 0x57, 0x5a,
	0x2c, 0xf4, 0x09, 0xec, 0xc2, 0x2c, 0xd5, 0x17, 0x48, 0xda, 0x8b, 0xf4, 0x37, 0x19, 0xcd, 0xcc,
	0x90, 0xb7, 0x1f, 0xcc, 0x84, 0x19, 0x66, 0x95, 0x73, 0x3e, 0x0e, 0xb9, 0x1f, 0x3c, 0x99, 0x9f,
	0xa9, 0x9e, 0xa4, 0xc5, 0xbd, 0x74, 0xb5, 0x99, 0xb5, 0xd5, 0xf5, 0xb7, 0x56, 0x95, 0x4f, 0xf4,
	0xc6, 0x3f, 0xfc, 0x05, 0xb2, 0x01, 0xff, 0x3f, 0x85, 0x56, 0x3d, 0x6e, 0xb6, 0xb8, 0x58, 0xfa,
	0x0c, 0xb1, 0xd0, 0x6a, 0x61, 0xa4, 0x88, 0xca, 0xfb, 0x06, 0x8e, 0xf3, 0x6a, 0x1d, 0x78, 0xce,
	0x3b, 0x88, 0x84, 0x56, 0x94, 0x42, 0x3c, 0x3a, 0x83, 0x8c, 0x14, 0xa4, 0xbc, 0xeb, 0x7d, 0xa6,
	0x39, 0x24, 0xa3, 0x9c, 0x27, 0xb4, 0x5d, 0xcb, 0xae, 0x3d, 0x3f, 0xf5, 0x75, 0xdf, 0x4a, 0x2b,
	0x59, 0x54, 0x90, 0x32, 0xed, 0x7d, 0xe6, 0x19, 0xa4, 0xef, 0x7f, 0xc6, 0xba, 0x70, 0x9a, 0xa7,
	0x00, 0xa1, 0x9b, 0x5f, 0xd7, 0xb4, 0x00, 0x42, 0xab, 0x01, 0xe7, 0xdd, 0xd7, 0x07, 0xd2, 0x37,
	0x48, 0x82, 0xe8, 0x42, 0x1f, 0x83, 0xd4, 0xa5, 0x79, 0xfe, 0x10, 0xf0, 0xf9, 0x0f, 0x7e, 0xa5,
	0x6e, 0x3d, 0x7b, 0x3d, 0x04, 0x00, 0x00, 0xff, 0xff, 0x99, 0x50, 0x8d, 0xb1, 0x0d, 0x01, 0x00,
	0x00,
}
