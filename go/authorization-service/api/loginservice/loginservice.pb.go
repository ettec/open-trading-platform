// Code generated by protoc-gen-go. DO NOT EDIT.
// source: loginservice.proto

package loginservice

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

type LoginParams struct {
	User                 string   `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginParams) Reset()         { *m = LoginParams{} }
func (m *LoginParams) String() string { return proto.CompactTextString(m) }
func (*LoginParams) ProtoMessage()    {}
func (*LoginParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_dc528863bee5e9c1, []int{0}
}

func (m *LoginParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginParams.Unmarshal(m, b)
}
func (m *LoginParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginParams.Marshal(b, m, deterministic)
}
func (m *LoginParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginParams.Merge(m, src)
}
func (m *LoginParams) XXX_Size() int {
	return xxx_messageInfo_LoginParams.Size(m)
}
func (m *LoginParams) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginParams.DiscardUnknown(m)
}

var xxx_messageInfo_LoginParams proto.InternalMessageInfo

func (m *LoginParams) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *LoginParams) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Token struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Desk                 string   `protobuf:"bytes,2,opt,name=desk,proto3" json:"desk,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_dc528863bee5e9c1, []int{1}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Token) GetDesk() string {
	if m != nil {
		return m.Desk
	}
	return ""
}

func init() {
	proto.RegisterType((*LoginParams)(nil), "loginservice.LoginParams")
	proto.RegisterType((*Token)(nil), "loginservice.Token")
}

func init() { proto.RegisterFile("loginservice.proto", fileDescriptor_dc528863bee5e9c1) }

var fileDescriptor_dc528863bee5e9c1 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xc9, 0x4f, 0xcf,
	0xcc, 0x2b, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x41, 0x16, 0x53, 0xb2, 0xe5, 0xe2, 0xf6, 0x01, 0xf1, 0x03, 0x12, 0x8b, 0x12, 0x73, 0x8b, 0x85,
	0x84, 0xb8, 0x58, 0x4a, 0x8b, 0x53, 0x8b, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c,
	0x21, 0x29, 0x2e, 0x8e, 0x82, 0xc4, 0xe2, 0xe2, 0xf2, 0xfc, 0xa2, 0x14, 0x09, 0x26, 0xb0, 0x38,
	0x9c, 0xaf, 0x64, 0xc8, 0xc5, 0x1a, 0x92, 0x9f, 0x9d, 0x9a, 0x27, 0x24, 0xc2, 0xc5, 0x5a, 0x02,
	0x62, 0x40, 0x75, 0x42, 0x38, 0x20, 0xe3, 0x52, 0x52, 0x8b, 0xb3, 0xa1, 0xda, 0xc0, 0x6c, 0x23,
	0x4f, 0x2e, 0x1e, 0xb0, 0x8d, 0xc1, 0x10, 0x17, 0x08, 0x59, 0x72, 0xb1, 0x82, 0xf9, 0x42, 0x92,
	0x7a, 0x28, 0xae, 0x45, 0x72, 0x96, 0x94, 0x30, 0xaa, 0x14, 0xd8, 0x4a, 0x25, 0x86, 0x24, 0x36,
	0xb0, 0x8f, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3d, 0xb4, 0x9d, 0x66, 0xe7, 0x00, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LoginServiceClient is the client API for LoginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LoginServiceClient interface {
	Login(ctx context.Context, in *LoginParams, opts ...grpc.CallOption) (*Token, error)
}

type loginServiceClient struct {
	cc *grpc.ClientConn
}

func NewLoginServiceClient(cc *grpc.ClientConn) LoginServiceClient {
	return &loginServiceClient{cc}
}

func (c *loginServiceClient) Login(ctx context.Context, in *LoginParams, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/loginservice.LoginService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginServiceServer is the server API for LoginService service.
type LoginServiceServer interface {
	Login(context.Context, *LoginParams) (*Token, error)
}

// UnimplementedLoginServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLoginServiceServer struct {
}

func (*UnimplementedLoginServiceServer) Login(ctx context.Context, req *LoginParams) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}

func RegisterLoginServiceServer(s *grpc.Server, srv LoginServiceServer) {
	s.RegisterService(&_LoginService_serviceDesc, srv)
}

func _LoginService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loginservice.LoginService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).Login(ctx, req.(*LoginParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _LoginService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "loginservice.LoginService",
	HandlerType: (*LoginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _LoginService_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "loginservice.proto",
}