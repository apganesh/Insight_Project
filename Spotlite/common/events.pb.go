// Code generated by protoc-gen-go. DO NOT EDIT.
// source: events.proto

/*
Package insight is a generated protocol buffer package.

It is generated from these files:
	events.proto

It has these top-level messages:
	Driver
	Rider
*/
package insight

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

type Driver struct {
	Id        int64   `protobuf:"varint,1,opt,name=Id" json:"Id,omitempty"`
	Lat       float64 `protobuf:"fixed64,2,opt,name=Lat" json:"Lat,omitempty"`
	Lng       float64 `protobuf:"fixed64,3,opt,name=Lng" json:"Lng,omitempty"`
	Radius    float64 `protobuf:"fixed64,4,opt,name=Radius" json:"Radius,omitempty"`
	Timestamp int64   `protobuf:"varint,6,opt,name=Timestamp" json:"Timestamp,omitempty"`
	Status    int64   `protobuf:"varint,5,opt,name=Status" json:"Status,omitempty"`
}

func (m *Driver) Reset()                    { *m = Driver{} }
func (m *Driver) String() string            { return proto.CompactTextString(m) }
func (*Driver) ProtoMessage()               {}
func (*Driver) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Driver) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Driver) GetLat() float64 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Driver) GetLng() float64 {
	if m != nil {
		return m.Lng
	}
	return 0
}

func (m *Driver) GetRadius() float64 {
	if m != nil {
		return m.Radius
	}
	return 0
}

func (m *Driver) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Driver) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

type Rider struct {
	Id        int64   `protobuf:"varint,1,opt,name=Id" json:"Id,omitempty"`
	SLat      float64 `protobuf:"fixed64,2,opt,name=sLat" json:"sLat,omitempty"`
	SLng      float64 `protobuf:"fixed64,3,opt,name=sLng" json:"sLng,omitempty"`
	ELat      float64 `protobuf:"fixed64,4,opt,name=eLat" json:"eLat,omitempty"`
	ELng      float64 `protobuf:"fixed64,5,opt,name=eLng" json:"eLng,omitempty"`
	Timestamp int64   `protobuf:"varint,6,opt,name=Timestamp" json:"Timestamp,omitempty"`
}

func (m *Rider) Reset()                    { *m = Rider{} }
func (m *Rider) String() string            { return proto.CompactTextString(m) }
func (*Rider) ProtoMessage()               {}
func (*Rider) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Rider) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Rider) GetSLat() float64 {
	if m != nil {
		return m.SLat
	}
	return 0
}

func (m *Rider) GetSLng() float64 {
	if m != nil {
		return m.SLng
	}
	return 0
}

func (m *Rider) GetELat() float64 {
	if m != nil {
		return m.ELat
	}
	return 0
}

func (m *Rider) GetELng() float64 {
	if m != nil {
		return m.ELng
	}
	return 0
}

func (m *Rider) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*Driver)(nil), "insight.Driver")
	proto.RegisterType((*Rider)(nil), "insight.Rider")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Matcher service

type MatcherClient interface {
	GetDriver(ctx context.Context, in *Rider, opts ...grpc.CallOption) (*Driver, error)
	AddDriver(ctx context.Context, in *Driver, opts ...grpc.CallOption) (*Driver, error)
}

type matcherClient struct {
	cc *grpc.ClientConn
}

func NewMatcherClient(cc *grpc.ClientConn) MatcherClient {
	return &matcherClient{cc}
}

func (c *matcherClient) GetDriver(ctx context.Context, in *Rider, opts ...grpc.CallOption) (*Driver, error) {
	out := new(Driver)
	err := grpc.Invoke(ctx, "/insight.Matcher/GetDriver", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matcherClient) AddDriver(ctx context.Context, in *Driver, opts ...grpc.CallOption) (*Driver, error) {
	out := new(Driver)
	err := grpc.Invoke(ctx, "/insight.Matcher/AddDriver", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Matcher service

type MatcherServer interface {
	GetDriver(context.Context, *Rider) (*Driver, error)
	AddDriver(context.Context, *Driver) (*Driver, error)
}

func RegisterMatcherServer(s *grpc.Server, srv MatcherServer) {
	s.RegisterService(&_Matcher_serviceDesc, srv)
}

func _Matcher_GetDriver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Rider)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatcherServer).GetDriver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/insight.Matcher/GetDriver",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatcherServer).GetDriver(ctx, req.(*Rider))
	}
	return interceptor(ctx, in, info, handler)
}

func _Matcher_AddDriver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Driver)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatcherServer).AddDriver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/insight.Matcher/AddDriver",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatcherServer).AddDriver(ctx, req.(*Driver))
	}
	return interceptor(ctx, in, info, handler)
}

var _Matcher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "insight.Matcher",
	HandlerType: (*MatcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDriver",
			Handler:    _Matcher_GetDriver_Handler,
		},
		{
			MethodName: "AddDriver",
			Handler:    _Matcher_AddDriver_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "events.proto",
}

func init() { proto.RegisterFile("events.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xcd, 0x6e, 0xbb, 0x65, 0x07, 0x51, 0x99, 0x83, 0x04, 0xf1, 0x50, 0xf6, 0xd4, 0xd3,
	0x0a, 0xfa, 0x04, 0x82, 0x20, 0x05, 0xbd, 0xac, 0xbe, 0x40, 0x34, 0x21, 0x8d, 0xd0, 0xb4, 0x24,
	0xd3, 0xbe, 0x81, 0x17, 0x9f, 0x5a, 0x32, 0x89, 0x2d, 0xc8, 0xd2, 0xdb, 0x3f, 0xdf, 0xcc, 0xec,
	0x7e, 0x4c, 0xe0, 0xdc, 0xec, 0x8d, 0xa7, 0xd8, 0x6f, 0xc3, 0x86, 0x36, 0x38, 0x73, 0x3e, 0x3a,
	0xbb, 0xa2, 0xee, 0x47, 0x40, 0xf3, 0x14, 0xdc, 0xde, 0x04, 0xbc, 0x80, 0x6a, 0xa9, 0xa5, 0x98,
	0x8b, 0x45, 0x3d, 0x54, 0x4b, 0x8d, 0x57, 0x50, 0xbf, 0x28, 0x92, 0xd5, 0x5c, 0x2c, 0xc4, 0x90,
	0x22, 0x13, 0x6f, 0x65, 0x5d, 0x88, 0xb7, 0x78, 0x0d, 0xcd, 0xa0, 0xb4, 0xdb, 0x45, 0x39, 0x61,
	0x58, 0x2a, 0xbc, 0x85, 0xf6, 0xdd, 0xad, 0x4d, 0x24, 0xb5, 0xde, 0xca, 0x86, 0x3f, 0x79, 0x04,
	0x69, 0xeb, 0x8d, 0x14, 0xed, 0xa2, 0x9c, 0x72, 0xab, 0x54, 0xdd, 0xb7, 0x80, 0xe9, 0xe0, 0xf4,
	0x88, 0x0b, 0xc2, 0x24, 0x1e, 0x65, 0x38, 0x67, 0x76, 0xd0, 0xe1, 0x9c, 0x98, 0x49, 0x73, 0xd9,
	0x86, 0x73, 0x66, 0xde, 0xf2, 0xbf, 0x98, 0x79, 0x7b, 0xda, 0xef, 0xfe, 0x0b, 0x66, 0xaf, 0x8a,
	0x3e, 0x57, 0x26, 0x60, 0x0f, 0xed, 0xb3, 0xa1, 0xbf, 0x0b, 0xf5, 0xe5, 0x6c, 0x3d, 0x5b, 0xde,
	0x5c, 0x1e, 0xea, 0x3c, 0xd0, 0x9d, 0xe1, 0x1d, 0xb4, 0x8f, 0x5a, 0x97, 0xf9, 0xff, 0xfd, 0x91,
	0x85, 0x8f, 0x86, 0x1f, 0xe4, 0xe1, 0x37, 0x00, 0x00, 0xff, 0xff, 0x95, 0x3b, 0x7a, 0x11, 0xa0,
	0x01, 0x00, 0x00,
}
