// Code generated by protoc-gen-go. DO NOT EDIT.
// source: clobquote.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type ClobLine struct {
	BidSize              *Decimal64 `protobuf:"bytes,1,opt,name=bidSize,proto3" json:"bidSize,omitempty"`
	BidPrice             *Decimal64 `protobuf:"bytes,2,opt,name=bidPrice,proto3" json:"bidPrice,omitempty"`
	AskSize              *Decimal64 `protobuf:"bytes,3,opt,name=askSize,proto3" json:"askSize,omitempty"`
	AskPrice             *Decimal64 `protobuf:"bytes,4,opt,name=askPrice,proto3" json:"askPrice,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ClobLine) Reset()         { *m = ClobLine{} }
func (m *ClobLine) String() string { return proto.CompactTextString(m) }
func (*ClobLine) ProtoMessage()    {}
func (*ClobLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_eff833333d312bfe, []int{0}
}

func (m *ClobLine) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClobLine.Unmarshal(m, b)
}
func (m *ClobLine) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClobLine.Marshal(b, m, deterministic)
}
func (m *ClobLine) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClobLine.Merge(m, src)
}
func (m *ClobLine) XXX_Size() int {
	return xxx_messageInfo_ClobLine.Size(m)
}
func (m *ClobLine) XXX_DiscardUnknown() {
	xxx_messageInfo_ClobLine.DiscardUnknown(m)
}

var xxx_messageInfo_ClobLine proto.InternalMessageInfo

func (m *ClobLine) GetBidSize() *Decimal64 {
	if m != nil {
		return m.BidSize
	}
	return nil
}

func (m *ClobLine) GetBidPrice() *Decimal64 {
	if m != nil {
		return m.BidPrice
	}
	return nil
}

func (m *ClobLine) GetAskSize() *Decimal64 {
	if m != nil {
		return m.AskSize
	}
	return nil
}

func (m *ClobLine) GetAskPrice() *Decimal64 {
	if m != nil {
		return m.AskPrice
	}
	return nil
}

type ClobQuote struct {
	ListingId            int32       `protobuf:"varint,1,opt,name=listingId,proto3" json:"listingId,omitempty"`
	Depth                []*ClobLine `protobuf:"bytes,2,rep,name=depth,proto3" json:"depth,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ClobQuote) Reset()         { *m = ClobQuote{} }
func (m *ClobQuote) String() string { return proto.CompactTextString(m) }
func (*ClobQuote) ProtoMessage()    {}
func (*ClobQuote) Descriptor() ([]byte, []int) {
	return fileDescriptor_eff833333d312bfe, []int{1}
}

func (m *ClobQuote) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClobQuote.Unmarshal(m, b)
}
func (m *ClobQuote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClobQuote.Marshal(b, m, deterministic)
}
func (m *ClobQuote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClobQuote.Merge(m, src)
}
func (m *ClobQuote) XXX_Size() int {
	return xxx_messageInfo_ClobQuote.Size(m)
}
func (m *ClobQuote) XXX_DiscardUnknown() {
	xxx_messageInfo_ClobQuote.DiscardUnknown(m)
}

var xxx_messageInfo_ClobQuote proto.InternalMessageInfo

func (m *ClobQuote) GetListingId() int32 {
	if m != nil {
		return m.ListingId
	}
	return 0
}

func (m *ClobQuote) GetDepth() []*ClobLine {
	if m != nil {
		return m.Depth
	}
	return nil
}

func init() {
	proto.RegisterType((*ClobLine)(nil), "model.ClobLine")
	proto.RegisterType((*ClobQuote)(nil), "model.ClobQuote")
}

func init() { proto.RegisterFile("clobquote.proto", fileDescriptor_eff833333d312bfe) }

var fileDescriptor_eff833333d312bfe = []byte{
	// 213 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0xce, 0xc9, 0x4f,
	0x2a, 0x2c, 0xcd, 0x2f, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0xcd, 0x4f,
	0x49, 0xcd, 0x91, 0xe2, 0x49, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0x83, 0x08, 0x2a, 0xed, 0x63, 0xe4,
	0xe2, 0x70, 0xce, 0xc9, 0x4f, 0xf2, 0xc9, 0xcc, 0x4b, 0x15, 0xd2, 0xe2, 0x62, 0x4f, 0xca, 0x4c,
	0x09, 0xce, 0xac, 0x4a, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x36, 0x12, 0xd0, 0x03, 0xeb, 0xd1,
	0x73, 0x49, 0x4d, 0xce, 0xcc, 0x4d, 0xcc, 0x31, 0x33, 0x09, 0x82, 0x29, 0x10, 0xd2, 0xe1, 0xe2,
	0x48, 0xca, 0x4c, 0x09, 0x28, 0xca, 0x4c, 0x4e, 0x95, 0x60, 0xc2, 0xa1, 0x18, 0xae, 0x02, 0x64,
	0x72, 0x62, 0x71, 0x36, 0xd8, 0x64, 0x66, 0x5c, 0x26, 0x43, 0x15, 0x80, 0x4c, 0x4e, 0x2c, 0xce,
	0x86, 0x98, 0xcc, 0x82, 0xcb, 0x64, 0x98, 0x0a, 0xa5, 0x00, 0x2e, 0x4e, 0x90, 0xfb, 0x03, 0x41,
	0x1e, 0x15, 0x92, 0xe1, 0xe2, 0xcc, 0xc9, 0x2c, 0x2e, 0xc9, 0xcc, 0x4b, 0xf7, 0x4c, 0x01, 0x7b,
	0x81, 0x35, 0x08, 0x21, 0x20, 0xa4, 0xca, 0xc5, 0x9a, 0x92, 0x5a, 0x50, 0x92, 0x21, 0xc1, 0xa4,
	0xc0, 0xac, 0xc1, 0x6d, 0xc4, 0x0f, 0x35, 0x15, 0xe6, 0xfd, 0x20, 0x88, 0xac, 0x13, 0x7b, 0x14,
	0x24, 0xa4, 0x92, 0xd8, 0xc0, 0x41, 0x64, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xcd, 0xf0, 0xfd,
	0x63, 0x4a, 0x01, 0x00, 0x00,
}