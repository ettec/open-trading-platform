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
	Size                 *Decimal64 `protobuf:"bytes,1,opt,name=size,proto3" json:"size,omitempty"`
	Price                *Decimal64 `protobuf:"bytes,2,opt,name=price,proto3" json:"price,omitempty"`
	EntryId              string     `protobuf:"bytes,3,opt,name=entryId,proto3" json:"entryId,omitempty"`
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

func (m *ClobLine) GetSize() *Decimal64 {
	if m != nil {
		return m.Size
	}
	return nil
}

func (m *ClobLine) GetPrice() *Decimal64 {
	if m != nil {
		return m.Price
	}
	return nil
}

func (m *ClobLine) GetEntryId() string {
	if m != nil {
		return m.EntryId
	}
	return ""
}

type ClobQuote struct {
	ListingId            int32       `protobuf:"varint,1,opt,name=listingId,proto3" json:"listingId,omitempty"`
	Bids                 []*ClobLine `protobuf:"bytes,2,rep,name=bids,proto3" json:"bids,omitempty"`
	Offers               []*ClobLine `protobuf:"bytes,3,rep,name=offers,proto3" json:"offers,omitempty"`
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

func (m *ClobQuote) GetBids() []*ClobLine {
	if m != nil {
		return m.Bids
	}
	return nil
}

func (m *ClobQuote) GetOffers() []*ClobLine {
	if m != nil {
		return m.Offers
	}
	return nil
}

func init() {
	proto.RegisterType((*ClobLine)(nil), "model.ClobLine")
	proto.RegisterType((*ClobQuote)(nil), "model.ClobQuote")
}

func init() { proto.RegisterFile("clobquote.proto", fileDescriptor_eff833333d312bfe) }

var fileDescriptor_eff833333d312bfe = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0xce, 0xc9, 0x4f,
	0x2a, 0x2c, 0xcd, 0x2f, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0xcd, 0x4f,
	0x49, 0xcd, 0x91, 0xe2, 0x49, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0x83, 0x08, 0x2a, 0x15, 0x71, 0x71,
	0x38, 0xe7, 0xe4, 0x27, 0xf9, 0x64, 0xe6, 0xa5, 0x0a, 0xa9, 0x70, 0xb1, 0x14, 0x67, 0x56, 0xa5,
	0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x09, 0xe8, 0x81, 0xd5, 0xeb, 0xb9, 0xa4, 0x26, 0x67,
	0xe6, 0x26, 0xe6, 0x98, 0x99, 0x04, 0x81, 0x65, 0x85, 0xd4, 0xb8, 0x58, 0x0b, 0x8a, 0x32, 0x93,
	0x53, 0x25, 0x98, 0x70, 0x28, 0x83, 0x48, 0x0b, 0x49, 0x70, 0xb1, 0xa7, 0xe6, 0x95, 0x14, 0x55,
	0x7a, 0xa6, 0x48, 0x30, 0x2b, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x4a, 0xe5, 0x5c, 0x9c, 0x20,
	0x3b, 0x03, 0x41, 0x6e, 0x13, 0x92, 0xe1, 0xe2, 0xcc, 0xc9, 0x2c, 0x2e, 0xc9, 0xcc, 0x4b, 0xf7,
	0x4c, 0x01, 0xdb, 0xcc, 0x1a, 0x84, 0x10, 0x10, 0x52, 0xe6, 0x62, 0x49, 0xca, 0x4c, 0x29, 0x96,
	0x60, 0x52, 0x60, 0xd6, 0xe0, 0x36, 0xe2, 0x87, 0xda, 0x05, 0x73, 0x71, 0x10, 0x58, 0x52, 0x48,
	0x9d, 0x8b, 0x2d, 0x3f, 0x2d, 0x2d, 0xb5, 0xa8, 0x58, 0x82, 0x19, 0xbb, 0x32, 0xa8, 0xb4, 0x13,
	0x7b, 0x14, 0x24, 0x0c, 0x92, 0xd8, 0xc0, 0x9e, 0x37, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x65,
	0x9a, 0x56, 0xfb, 0x24, 0x01, 0x00, 0x00,
}