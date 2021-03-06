// Code generated by protoc-gen-go. DO NOT EDIT.
// source: swag-engine.proto

package protobuf

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

// The health request message (it's empty)
type HealthRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthRequest) Reset()         { *m = HealthRequest{} }
func (m *HealthRequest) String() string { return proto.CompactTextString(m) }
func (*HealthRequest) ProtoMessage()    {}
func (*HealthRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_swag_engine_ea8beac46599a85c, []int{0}
}
func (m *HealthRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthRequest.Unmarshal(m, b)
}
func (m *HealthRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthRequest.Marshal(b, m, deterministic)
}
func (dst *HealthRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthRequest.Merge(dst, src)
}
func (m *HealthRequest) XXX_Size() int {
	return xxx_messageInfo_HealthRequest.Size(m)
}
func (m *HealthRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HealthRequest proto.InternalMessageInfo

// The health response message containing the health status
type HealthReply struct {
	Healthstatus         string   `protobuf:"bytes,1,opt,name=healthstatus" json:"healthstatus,omitempty"`
	Processedevents      int32    `protobuf:"varint,2,opt,name=processedevents" json:"processedevents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthReply) Reset()         { *m = HealthReply{} }
func (m *HealthReply) String() string { return proto.CompactTextString(m) }
func (*HealthReply) ProtoMessage()    {}
func (*HealthReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_swag_engine_ea8beac46599a85c, []int{1}
}
func (m *HealthReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthReply.Unmarshal(m, b)
}
func (m *HealthReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthReply.Marshal(b, m, deterministic)
}
func (dst *HealthReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthReply.Merge(dst, src)
}
func (m *HealthReply) XXX_Size() int {
	return xxx_messageInfo_HealthReply.Size(m)
}
func (m *HealthReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthReply.DiscardUnknown(m)
}

var xxx_messageInfo_HealthReply proto.InternalMessageInfo

func (m *HealthReply) GetHealthstatus() string {
	if m != nil {
		return m.Healthstatus
	}
	return ""
}

func (m *HealthReply) GetProcessedevents() int32 {
	if m != nil {
		return m.Processedevents
	}
	return 0
}

// The mom post request message
type MomPostRequest struct {
	Momtype              string   `protobuf:"bytes,1,opt,name=momtype" json:"momtype,omitempty"`
	Momdefinition        string   `protobuf:"bytes,2,opt,name=momdefinition" json:"momdefinition,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MomPostRequest) Reset()         { *m = MomPostRequest{} }
func (m *MomPostRequest) String() string { return proto.CompactTextString(m) }
func (*MomPostRequest) ProtoMessage()    {}
func (*MomPostRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_swag_engine_ea8beac46599a85c, []int{2}
}
func (m *MomPostRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MomPostRequest.Unmarshal(m, b)
}
func (m *MomPostRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MomPostRequest.Marshal(b, m, deterministic)
}
func (dst *MomPostRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MomPostRequest.Merge(dst, src)
}
func (m *MomPostRequest) XXX_Size() int {
	return xxx_messageInfo_MomPostRequest.Size(m)
}
func (m *MomPostRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MomPostRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MomPostRequest proto.InternalMessageInfo

func (m *MomPostRequest) GetMomtype() string {
	if m != nil {
		return m.Momtype
	}
	return ""
}

func (m *MomPostRequest) GetMomdefinition() string {
	if m != nil {
		return m.Momdefinition
	}
	return ""
}

// The mom post response message
type MomPostReply struct {
	Deploymenterror      string   `protobuf:"bytes,1,opt,name=deploymenterror" json:"deploymenterror,omitempty"`
	Momid                string   `protobuf:"bytes,2,opt,name=momid" json:"momid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MomPostReply) Reset()         { *m = MomPostReply{} }
func (m *MomPostReply) String() string { return proto.CompactTextString(m) }
func (*MomPostReply) ProtoMessage()    {}
func (*MomPostReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_swag_engine_ea8beac46599a85c, []int{3}
}
func (m *MomPostReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MomPostReply.Unmarshal(m, b)
}
func (m *MomPostReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MomPostReply.Marshal(b, m, deterministic)
}
func (dst *MomPostReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MomPostReply.Merge(dst, src)
}
func (m *MomPostReply) XXX_Size() int {
	return xxx_messageInfo_MomPostReply.Size(m)
}
func (m *MomPostReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MomPostReply.DiscardUnknown(m)
}

var xxx_messageInfo_MomPostReply proto.InternalMessageInfo

func (m *MomPostReply) GetDeploymenterror() string {
	if m != nil {
		return m.Deploymenterror
	}
	return ""
}

func (m *MomPostReply) GetMomid() string {
	if m != nil {
		return m.Momid
	}
	return ""
}

// The mom post request message
type MomDeleteRequest struct {
	Momid                string   `protobuf:"bytes,1,opt,name=momid" json:"momid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MomDeleteRequest) Reset()         { *m = MomDeleteRequest{} }
func (m *MomDeleteRequest) String() string { return proto.CompactTextString(m) }
func (*MomDeleteRequest) ProtoMessage()    {}
func (*MomDeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_swag_engine_ea8beac46599a85c, []int{4}
}
func (m *MomDeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MomDeleteRequest.Unmarshal(m, b)
}
func (m *MomDeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MomDeleteRequest.Marshal(b, m, deterministic)
}
func (dst *MomDeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MomDeleteRequest.Merge(dst, src)
}
func (m *MomDeleteRequest) XXX_Size() int {
	return xxx_messageInfo_MomDeleteRequest.Size(m)
}
func (m *MomDeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MomDeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MomDeleteRequest proto.InternalMessageInfo

func (m *MomDeleteRequest) GetMomid() string {
	if m != nil {
		return m.Momid
	}
	return ""
}

// The mom post response message
type MomDeleteReply struct {
	Deletionerror        string   `protobuf:"bytes,1,opt,name=deletionerror" json:"deletionerror,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MomDeleteReply) Reset()         { *m = MomDeleteReply{} }
func (m *MomDeleteReply) String() string { return proto.CompactTextString(m) }
func (*MomDeleteReply) ProtoMessage()    {}
func (*MomDeleteReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_swag_engine_ea8beac46599a85c, []int{5}
}
func (m *MomDeleteReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MomDeleteReply.Unmarshal(m, b)
}
func (m *MomDeleteReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MomDeleteReply.Marshal(b, m, deterministic)
}
func (dst *MomDeleteReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MomDeleteReply.Merge(dst, src)
}
func (m *MomDeleteReply) XXX_Size() int {
	return xxx_messageInfo_MomDeleteReply.Size(m)
}
func (m *MomDeleteReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MomDeleteReply.DiscardUnknown(m)
}

var xxx_messageInfo_MomDeleteReply proto.InternalMessageInfo

func (m *MomDeleteReply) GetDeletionerror() string {
	if m != nil {
		return m.Deletionerror
	}
	return ""
}

func init() {
	proto.RegisterType((*HealthRequest)(nil), "protobuf.HealthRequest")
	proto.RegisterType((*HealthReply)(nil), "protobuf.HealthReply")
	proto.RegisterType((*MomPostRequest)(nil), "protobuf.MomPostRequest")
	proto.RegisterType((*MomPostReply)(nil), "protobuf.MomPostReply")
	proto.RegisterType((*MomDeleteRequest)(nil), "protobuf.MomDeleteRequest")
	proto.RegisterType((*MomDeleteReply)(nil), "protobuf.MomDeleteReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Engine service

type EngineClient interface {
	// Requires health
	GetHealth(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthReply, error)
	PostMoM(ctx context.Context, in *MomPostRequest, opts ...grpc.CallOption) (*MomPostReply, error)
	PostStamper(ctx context.Context, in *MomPostRequest, opts ...grpc.CallOption) (*MomPostReply, error)
	DeleteMoM(ctx context.Context, in *MomDeleteRequest, opts ...grpc.CallOption) (*MomDeleteReply, error)
	DeleteStamper(ctx context.Context, in *MomDeleteRequest, opts ...grpc.CallOption) (*MomDeleteReply, error)
}

type engineClient struct {
	cc *grpc.ClientConn
}

func NewEngineClient(cc *grpc.ClientConn) EngineClient {
	return &engineClient{cc}
}

func (c *engineClient) GetHealth(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthReply, error) {
	out := new(HealthReply)
	err := grpc.Invoke(ctx, "/protobuf.Engine/GetHealth", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *engineClient) PostMoM(ctx context.Context, in *MomPostRequest, opts ...grpc.CallOption) (*MomPostReply, error) {
	out := new(MomPostReply)
	err := grpc.Invoke(ctx, "/protobuf.Engine/PostMoM", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *engineClient) PostStamper(ctx context.Context, in *MomPostRequest, opts ...grpc.CallOption) (*MomPostReply, error) {
	out := new(MomPostReply)
	err := grpc.Invoke(ctx, "/protobuf.Engine/PostStamper", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *engineClient) DeleteMoM(ctx context.Context, in *MomDeleteRequest, opts ...grpc.CallOption) (*MomDeleteReply, error) {
	out := new(MomDeleteReply)
	err := grpc.Invoke(ctx, "/protobuf.Engine/DeleteMoM", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *engineClient) DeleteStamper(ctx context.Context, in *MomDeleteRequest, opts ...grpc.CallOption) (*MomDeleteReply, error) {
	out := new(MomDeleteReply)
	err := grpc.Invoke(ctx, "/protobuf.Engine/DeleteStamper", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Engine service

type EngineServer interface {
	// Requires health
	GetHealth(context.Context, *HealthRequest) (*HealthReply, error)
	PostMoM(context.Context, *MomPostRequest) (*MomPostReply, error)
	PostStamper(context.Context, *MomPostRequest) (*MomPostReply, error)
	DeleteMoM(context.Context, *MomDeleteRequest) (*MomDeleteReply, error)
	DeleteStamper(context.Context, *MomDeleteRequest) (*MomDeleteReply, error)
}

func RegisterEngineServer(s *grpc.Server, srv EngineServer) {
	s.RegisterService(&_Engine_serviceDesc, srv)
}

func _Engine_GetHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServer).GetHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Engine/GetHealth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServer).GetHealth(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Engine_PostMoM_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MomPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServer).PostMoM(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Engine/PostMoM",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServer).PostMoM(ctx, req.(*MomPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Engine_PostStamper_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MomPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServer).PostStamper(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Engine/PostStamper",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServer).PostStamper(ctx, req.(*MomPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Engine_DeleteMoM_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MomDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServer).DeleteMoM(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Engine/DeleteMoM",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServer).DeleteMoM(ctx, req.(*MomDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Engine_DeleteStamper_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MomDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServer).DeleteStamper(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Engine/DeleteStamper",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServer).DeleteStamper(ctx, req.(*MomDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Engine_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Engine",
	HandlerType: (*EngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHealth",
			Handler:    _Engine_GetHealth_Handler,
		},
		{
			MethodName: "PostMoM",
			Handler:    _Engine_PostMoM_Handler,
		},
		{
			MethodName: "PostStamper",
			Handler:    _Engine_PostStamper_Handler,
		},
		{
			MethodName: "DeleteMoM",
			Handler:    _Engine_DeleteMoM_Handler,
		},
		{
			MethodName: "DeleteStamper",
			Handler:    _Engine_DeleteStamper_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swag-engine.proto",
}

func init() { proto.RegisterFile("swag-engine.proto", fileDescriptor_swag_engine_ea8beac46599a85c) }

var fileDescriptor_swag_engine_ea8beac46599a85c = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0x3f, 0x4f, 0xe3, 0x40,
	0x10, 0xc5, 0xcf, 0x91, 0x92, 0x9c, 0x27, 0xf1, 0xe5, 0x6e, 0x75, 0x80, 0xe5, 0x2a, 0xb2, 0x52,
	0xb8, 0x21, 0x05, 0x48, 0x74, 0x29, 0x10, 0xa0, 0xd0, 0x18, 0x45, 0xa6, 0xa4, 0x72, 0xf0, 0x24,
	0xb1, 0xe4, 0xdd, 0x59, 0xbc, 0x1b, 0x90, 0x3f, 0x29, 0x5f, 0x07, 0xf9, 0x5f, 0xc8, 0x86, 0x54,
	0xa9, 0xac, 0xf7, 0xe6, 0xf9, 0x37, 0xb3, 0x33, 0xf0, 0x4f, 0x7d, 0xc4, 0xeb, 0x4b, 0x14, 0xeb,
	0x54, 0xe0, 0x54, 0xe6, 0xa4, 0x89, 0xfd, 0xae, 0x3e, 0xcb, 0xed, 0xca, 0x1f, 0x81, 0xf3, 0x88,
	0x71, 0xa6, 0x37, 0x11, 0xbe, 0x6d, 0x51, 0x69, 0xff, 0x05, 0x06, 0xad, 0x21, 0xb3, 0x82, 0xf9,
	0x30, 0xdc, 0x54, 0x52, 0xe9, 0x58, 0x6f, 0x95, 0x6b, 0x8d, 0xad, 0xc0, 0x8e, 0x0c, 0x8f, 0x05,
	0x30, 0x92, 0x39, 0xbd, 0xa2, 0x52, 0x98, 0xe0, 0x3b, 0x0a, 0xad, 0xdc, 0xce, 0xd8, 0x0a, 0xba,
	0xd1, 0xa1, 0xed, 0x2f, 0xe0, 0x4f, 0x48, 0x7c, 0x41, 0x4a, 0x37, 0xed, 0x98, 0x0b, 0x7d, 0x4e,
	0x5c, 0x17, 0x12, 0x1b, 0x74, 0x2b, 0xd9, 0x04, 0x1c, 0x4e, 0x3c, 0xc1, 0x55, 0x2a, 0x52, 0x9d,
	0x92, 0xa8, 0x98, 0x76, 0x64, 0x9a, 0xfe, 0x13, 0x0c, 0x77, 0xc4, 0x72, 0xde, 0x00, 0x46, 0x09,
	0xca, 0x8c, 0x0a, 0x8e, 0x42, 0x63, 0x9e, 0x53, 0xde, 0x70, 0x0f, 0x6d, 0xf6, 0x1f, 0xba, 0x9c,
	0x78, 0x9a, 0x34, 0xdc, 0x5a, 0xf8, 0x01, 0xfc, 0x0d, 0x89, 0xdf, 0x63, 0x86, 0x1a, 0xdb, 0x19,
	0x77, 0x49, 0x6b, 0x3f, 0x79, 0x53, 0xbd, 0xa5, 0x4d, 0x96, 0xbd, 0x27, 0xe0, 0x24, 0xa5, 0x4c,
	0x49, 0xec, 0x77, 0x36, 0xcd, 0xab, 0xcf, 0x0e, 0xf4, 0x1e, 0xaa, 0x63, 0xb0, 0x19, 0xd8, 0x73,
	0xd4, 0xf5, 0xba, 0xd9, 0xc5, 0xb4, 0x3d, 0xca, 0xd4, 0xb8, 0x88, 0x77, 0xf6, 0xb3, 0x20, 0xb3,
	0xc2, 0xff, 0xc5, 0x66, 0xd0, 0x2f, 0x1f, 0x1e, 0x52, 0xc8, 0xdc, 0xef, 0x8c, 0xb9, 0x60, 0xef,
	0xfc, 0x48, 0xa5, 0xfe, 0xfd, 0x16, 0x06, 0xa5, 0x7c, 0xd6, 0x31, 0x97, 0x98, 0x9f, 0x84, 0xb8,
	0x03, 0xbb, 0x5e, 0x40, 0x39, 0x83, 0x67, 0xc4, 0x8c, 0x15, 0x7a, 0xee, 0xd1, 0x5a, 0x0d, 0x99,
	0x83, 0x53, 0x1b, 0xed, 0x24, 0x27, 0x82, 0x96, 0xbd, 0xaa, 0x74, 0xfd, 0x15, 0x00, 0x00, 0xff,
	0xff, 0xc3, 0x23, 0xa5, 0x54, 0xf1, 0x02, 0x00, 0x00,
}
