// Code generated by protoc-gen-go. DO NOT EDIT.
// source: static-data-service.proto

package api

import (
	"github.com/ettec/open-trading-platform/go/model"
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

type ListingId struct {
	ListingId            int32    `protobuf:"varint,1,opt,name=listingId,proto3" json:"listingId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListingId) Reset()         { *m = ListingId{} }
func (m *ListingId) String() string { return proto.CompactTextString(m) }
func (*ListingId) ProtoMessage()    {}
func (*ListingId) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f596fb8d7650c9d, []int{0}
}

func (m *ListingId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListingId.Unmarshal(m, b)
}
func (m *ListingId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListingId.Marshal(b, m, deterministic)
}
func (m *ListingId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListingId.Merge(m, src)
}
func (m *ListingId) XXX_Size() int {
	return xxx_messageInfo_ListingId.Size(m)
}
func (m *ListingId) XXX_DiscardUnknown() {
	xxx_messageInfo_ListingId.DiscardUnknown(m)
}

var xxx_messageInfo_ListingId proto.InternalMessageInfo

func (m *ListingId) GetListingId() int32 {
	if m != nil {
		return m.ListingId
	}
	return 0
}

type ListingIds struct {
	ListingIds           []int32  `protobuf:"varint,1,rep,packed,name=listingIds,proto3" json:"listingIds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListingIds) Reset()         { *m = ListingIds{} }
func (m *ListingIds) String() string { return proto.CompactTextString(m) }
func (*ListingIds) ProtoMessage()    {}
func (*ListingIds) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f596fb8d7650c9d, []int{1}
}

func (m *ListingIds) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListingIds.Unmarshal(m, b)
}
func (m *ListingIds) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListingIds.Marshal(b, m, deterministic)
}
func (m *ListingIds) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListingIds.Merge(m, src)
}
func (m *ListingIds) XXX_Size() int {
	return xxx_messageInfo_ListingIds.Size(m)
}
func (m *ListingIds) XXX_DiscardUnknown() {
	xxx_messageInfo_ListingIds.DiscardUnknown(m)
}

var xxx_messageInfo_ListingIds proto.InternalMessageInfo

func (m *ListingIds) GetListingIds() []int32 {
	if m != nil {
		return m.ListingIds
	}
	return nil
}

type Listings struct {
	Listings             []*model.Listing `protobuf:"bytes,1,rep,name=listings,proto3" json:"listings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Listings) Reset()         { *m = Listings{} }
func (m *Listings) String() string { return proto.CompactTextString(m) }
func (*Listings) ProtoMessage()    {}
func (*Listings) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f596fb8d7650c9d, []int{2}
}

func (m *Listings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Listings.Unmarshal(m, b)
}
func (m *Listings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Listings.Marshal(b, m, deterministic)
}
func (m *Listings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Listings.Merge(m, src)
}
func (m *Listings) XXX_Size() int {
	return xxx_messageInfo_Listings.Size(m)
}
func (m *Listings) XXX_DiscardUnknown() {
	xxx_messageInfo_Listings.DiscardUnknown(m)
}

var xxx_messageInfo_Listings proto.InternalMessageInfo

func (m *Listings) GetListings() []*model.Listing {
	if m != nil {
		return m.Listings
	}
	return nil
}

type MatchParameters struct {
	SymbolMatch          string   `protobuf:"bytes,1,opt,name=symbolMatch,proto3" json:"symbolMatch,omitempty"`
	NameMatch            string   `protobuf:"bytes,2,opt,name=nameMatch,proto3" json:"nameMatch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MatchParameters) Reset()         { *m = MatchParameters{} }
func (m *MatchParameters) String() string { return proto.CompactTextString(m) }
func (*MatchParameters) ProtoMessage()    {}
func (*MatchParameters) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f596fb8d7650c9d, []int{3}
}

func (m *MatchParameters) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MatchParameters.Unmarshal(m, b)
}
func (m *MatchParameters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MatchParameters.Marshal(b, m, deterministic)
}
func (m *MatchParameters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MatchParameters.Merge(m, src)
}
func (m *MatchParameters) XXX_Size() int {
	return xxx_messageInfo_MatchParameters.Size(m)
}
func (m *MatchParameters) XXX_DiscardUnknown() {
	xxx_messageInfo_MatchParameters.DiscardUnknown(m)
}

var xxx_messageInfo_MatchParameters proto.InternalMessageInfo

func (m *MatchParameters) GetSymbolMatch() string {
	if m != nil {
		return m.SymbolMatch
	}
	return ""
}

func (m *MatchParameters) GetNameMatch() string {
	if m != nil {
		return m.NameMatch
	}
	return ""
}

func init() {
	proto.RegisterType((*ListingId)(nil), "staticdataservice.ListingId")
	proto.RegisterType((*ListingIds)(nil), "staticdataservice.ListingIds")
	proto.RegisterType((*Listings)(nil), "staticdataservice.Listings")
	proto.RegisterType((*MatchParameters)(nil), "staticdataservice.MatchParameters")
}

func init() { proto.RegisterFile("static-data-service.proto", fileDescriptor_1f596fb8d7650c9d) }

var fileDescriptor_1f596fb8d7650c9d = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0x9b, 0x96, 0x48, 0x73, 0x83, 0x4a, 0xaf, 0x9b, 0x5a, 0xab, 0x84, 0x59, 0x55, 0xb1,
	0x59, 0x54, 0x70, 0xe5, 0x4a, 0x04, 0xa9, 0x3f, 0xa0, 0xe9, 0x46, 0xdc, 0xdd, 0x26, 0x43, 0x1d,
	0xc8, 0x4f, 0xc9, 0x0c, 0x82, 0xef, 0xe7, 0x83, 0x49, 0x67, 0xa6, 0x49, 0xa8, 0x55, 0x71, 0x99,
	0xf3, 0x9d, 0x1c, 0x3e, 0x2e, 0x03, 0x87, 0x52, 0x91, 0x12, 0xf1, 0x38, 0x21, 0x45, 0x63, 0xc9,
	0xcb, 0x77, 0x11, 0xf3, 0x70, 0x59, 0x16, 0xaa, 0xc0, 0x9e, 0x41, 0x2b, 0x62, 0xc1, 0x60, 0x37,
	0x15, 0x52, 0x89, 0x7c, 0x61, 0x1a, 0xec, 0x14, 0xbc, 0x07, 0x13, 0x4c, 0x13, 0x1c, 0x82, 0x97,
	0xae, 0x3f, 0xfa, 0x4e, 0xe0, 0x8c, 0xdc, 0xa8, 0x0e, 0xd8, 0x39, 0x40, 0x55, 0x95, 0x78, 0x02,
	0x50, 0x21, 0xd9, 0x77, 0x82, 0xce, 0xc8, 0x8d, 0x1a, 0x09, 0xbb, 0x84, 0xae, 0x6d, 0x4b, 0x3c,
	0x83, 0xae, 0x25, 0xa6, 0xe9, 0x4f, 0xf6, 0xc2, 0xac, 0x48, 0x78, 0x1a, 0xda, 0x4a, 0x54, 0x71,
	0xf6, 0x0c, 0xfb, 0x8f, 0xa4, 0xe2, 0xb7, 0x27, 0x2a, 0x29, 0xe3, 0x8a, 0x97, 0x12, 0x03, 0xf0,
	0xe5, 0x47, 0x36, 0x2f, 0x52, 0x0d, 0xb4, 0x98, 0x17, 0x35, 0xa3, 0x95, 0x78, 0x4e, 0x19, 0x37,
	0xbc, 0xad, 0x79, 0x1d, 0x4c, 0x3e, 0xdb, 0xd0, 0x9b, 0xe9, 0x43, 0xdc, 0x90, 0xa2, 0x99, 0x39,
	0x04, 0xde, 0x01, 0xde, 0x72, 0x65, 0x05, 0x74, 0x51, 0xe4, 0x0b, 0x64, 0xe1, 0xb7, 0x93, 0x85,
	0x1b, 0x3e, 0x83, 0x0d, 0x79, 0xd6, 0xc2, 0x17, 0x38, 0xa8, 0xb7, 0xe4, 0xbf, 0xc6, 0x8e, 0xb6,
	0x74, 0xd6, 0x43, 0xac, 0x85, 0x57, 0x00, 0xf5, 0x32, 0x0e, 0x7f, 0x2e, 0x4f, 0x93, 0x2d, 0x5e,
	0xf7, 0xe0, 0x37, 0xbc, 0xf0, 0xf8, 0xb7, 0xdf, 0xff, 0x52, 0xb9, 0x76, 0x5f, 0x3b, 0xb4, 0x14,
	0xf3, 0x1d, 0xfd, 0x70, 0x2e, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x6d, 0xb3, 0x94, 0xe2, 0x77,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StaticDataServiceClient is the client API for StaticDataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StaticDataServiceClient interface {
	GetListingMatching(ctx context.Context, in *MatchParameters, opts ...grpc.CallOption) (*model.Listing, error)
	GetListingsMatching(ctx context.Context, in *MatchParameters, opts ...grpc.CallOption) (*Listings, error)
	GetListing(ctx context.Context, in *ListingId, opts ...grpc.CallOption) (*model.Listing, error)
	GetListings(ctx context.Context, in *ListingIds, opts ...grpc.CallOption) (*Listings, error)
}

type staticDataServiceClient struct {
	cc *grpc.ClientConn
}

func NewStaticDataServiceClient(cc *grpc.ClientConn) StaticDataServiceClient {
	return &staticDataServiceClient{cc}
}

func (c *staticDataServiceClient) GetListingMatching(ctx context.Context, in *MatchParameters, opts ...grpc.CallOption) (*model.Listing, error) {
	out := new(Listing)
	err := c.cc.Invoke(ctx, "/staticdataservice.StaticDataService/GetListingMatching", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staticDataServiceClient) GetListingsMatching(ctx context.Context, in *MatchParameters, opts ...grpc.CallOption) (*Listings, error) {
	out := new(Listings)
	err := c.cc.Invoke(ctx, "/staticdataservice.StaticDataService/GetListingsMatching", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staticDataServiceClient) GetListing(ctx context.Context, in *ListingId, opts ...grpc.CallOption) (*model.Listing, error) {
	out := new(Listing)
	err := c.cc.Invoke(ctx, "/staticdataservice.StaticDataService/GetListing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staticDataServiceClient) GetListings(ctx context.Context, in *ListingIds, opts ...grpc.CallOption) (*Listings, error) {
	out := new(Listings)
	err := c.cc.Invoke(ctx, "/staticdataservice.StaticDataService/GetListings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StaticDataServiceServer is the server API for StaticDataService service.
type StaticDataServiceServer interface {
	GetListingMatching(context.Context, *MatchParameters) (*model.Listing, error)
	GetListingsMatching(context.Context, *MatchParameters) (*Listings, error)
	GetListing(context.Context, *ListingId) (*model.Listing, error)
	GetListings(context.Context, *ListingIds) (*Listings, error)
}

// UnimplementedStaticDataServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStaticDataServiceServer struct {
}

func (*UnimplementedStaticDataServiceServer) GetListingMatching(ctx context.Context, req *MatchParameters) (*model.Listing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListingMatching not implemented")
}
func (*UnimplementedStaticDataServiceServer) GetListingsMatching(ctx context.Context, req *MatchParameters) (*Listings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListingsMatching not implemented")
}
func (*UnimplementedStaticDataServiceServer) GetListing(ctx context.Context, req *ListingId) (*model.Listing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListing not implemented")
}
func (*UnimplementedStaticDataServiceServer) GetListings(ctx context.Context, req *ListingIds) (*Listings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListings not implemented")
}

func RegisterStaticDataServiceServer(s *grpc.Server, srv StaticDataServiceServer) {
	s.RegisterService(&_StaticDataService_serviceDesc, srv)
}

func _StaticDataService_GetListingMatching_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticDataServiceServer).GetListingMatching(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/staticdataservice.StaticDataService/GetListingMatching",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticDataServiceServer).GetListingMatching(ctx, req.(*MatchParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaticDataService_GetListingsMatching_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticDataServiceServer).GetListingsMatching(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/staticdataservice.StaticDataService/GetListingsMatching",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticDataServiceServer).GetListingsMatching(ctx, req.(*MatchParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaticDataService_GetListing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListingId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticDataServiceServer).GetListing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/staticdataservice.StaticDataService/GetListing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticDataServiceServer).GetListing(ctx, req.(*ListingId))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaticDataService_GetListings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListingIds)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticDataServiceServer).GetListings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/staticdataservice.StaticDataService/GetListings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticDataServiceServer).GetListings(ctx, req.(*ListingIds))
	}
	return interceptor(ctx, in, info, handler)
}

var _StaticDataService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "staticdataservice.StaticDataService",
	HandlerType: (*StaticDataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetListingMatching",
			Handler:    _StaticDataService_GetListingMatching_Handler,
		},
		{
			MethodName: "GetListingsMatching",
			Handler:    _StaticDataService_GetListingsMatching_Handler,
		},
		{
			MethodName: "GetListing",
			Handler:    _StaticDataService_GetListing_Handler,
		},
		{
			MethodName: "GetListings",
			Handler:    _StaticDataService_GetListings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "static-data-service.proto",
}
