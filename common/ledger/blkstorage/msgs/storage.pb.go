// Code generated by protoc-gen-go. DO NOT EDIT.
// source: storage.proto

package msgs

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

type TxIDIndexValProto struct {
	BlkLocation          []byte   `protobuf:"bytes,1,opt,name=blk_location,json=blkLocation,proto3" json:"blk_location,omitempty"`
	TxLocation           []byte   `protobuf:"bytes,2,opt,name=tx_location,json=txLocation,proto3" json:"tx_location,omitempty"`
	TxValidationCode     int32    `protobuf:"varint,3,opt,name=tx_validation_code,json=txValidationCode,proto3" json:"tx_validation_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxIDIndexValProto) Reset()         { *m = TxIDIndexValProto{} }
func (m *TxIDIndexValProto) String() string { return proto.CompactTextString(m) }
func (*TxIDIndexValProto) ProtoMessage()    {}
func (*TxIDIndexValProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_0d2c4ccf1453ffdb, []int{0}
}

func (m *TxIDIndexValProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxIDIndexValProto.Unmarshal(m, b)
}
func (m *TxIDIndexValProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxIDIndexValProto.Marshal(b, m, deterministic)
}
func (m *TxIDIndexValProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxIDIndexValProto.Merge(m, src)
}
func (m *TxIDIndexValProto) XXX_Size() int {
	return xxx_messageInfo_TxIDIndexValProto.Size(m)
}
func (m *TxIDIndexValProto) XXX_DiscardUnknown() {
	xxx_messageInfo_TxIDIndexValProto.DiscardUnknown(m)
}

var xxx_messageInfo_TxIDIndexValProto proto.InternalMessageInfo

func (m *TxIDIndexValProto) GetBlkLocation() []byte {
	if m != nil {
		return m.BlkLocation
	}
	return nil
}

func (m *TxIDIndexValProto) GetTxLocation() []byte {
	if m != nil {
		return m.TxLocation
	}
	return nil
}

func (m *TxIDIndexValProto) GetTxValidationCode() int32 {
	if m != nil {
		return m.TxValidationCode
	}
	return 0
}

type TxDateIndexValProto struct {
	BlkNumber            int64    `protobuf:"varint,1,opt,name=blk_number,json=blkNumber,proto3" json:"blk_number,omitempty"`
	BlkLocation          []byte   `protobuf:"bytes,2,opt,name=blk_location,json=blkLocation,proto3" json:"blk_location,omitempty"`
	TxLocation           []byte   `protobuf:"bytes,3,opt,name=tx_location,json=txLocation,proto3" json:"tx_location,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxDateIndexValProto) Reset()         { *m = TxDateIndexValProto{} }
func (m *TxDateIndexValProto) String() string { return proto.CompactTextString(m) }
func (*TxDateIndexValProto) ProtoMessage()    {}
func (*TxDateIndexValProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_0d2c4ccf1453ffdb, []int{1}
}

func (m *TxDateIndexValProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxDateIndexValProto.Unmarshal(m, b)
}
func (m *TxDateIndexValProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxDateIndexValProto.Marshal(b, m, deterministic)
}
func (m *TxDateIndexValProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxDateIndexValProto.Merge(m, src)
}
func (m *TxDateIndexValProto) XXX_Size() int {
	return xxx_messageInfo_TxDateIndexValProto.Size(m)
}
func (m *TxDateIndexValProto) XXX_DiscardUnknown() {
	xxx_messageInfo_TxDateIndexValProto.DiscardUnknown(m)
}

var xxx_messageInfo_TxDateIndexValProto proto.InternalMessageInfo

func (m *TxDateIndexValProto) GetBlkNumber() int64 {
	if m != nil {
		return m.BlkNumber
	}
	return 0
}

func (m *TxDateIndexValProto) GetBlkLocation() []byte {
	if m != nil {
		return m.BlkLocation
	}
	return nil
}

func (m *TxDateIndexValProto) GetTxLocation() []byte {
	if m != nil {
		return m.TxLocation
	}
	return nil
}

func init() {
	proto.RegisterType((*TxIDIndexValProto)(nil), "msgs.txIDIndexValProto")
	proto.RegisterType((*TxDateIndexValProto)(nil), "msgs.txDateIndexValProto")
}

func init() { proto.RegisterFile("storage.proto", fileDescriptor_0d2c4ccf1453ffdb) }

var fileDescriptor_0d2c4ccf1453ffdb = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0x87, 0x49, 0xa3, 0x82, 0xdb, 0x0a, 0x1a, 0x2f, 0xbd, 0x88, 0xb5, 0xa7, 0x1e, 0xa4, 0x39,
	0xf8, 0x06, 0xda, 0x4b, 0xa1, 0x14, 0xc9, 0xc1, 0x83, 0x97, 0xb0, 0xb3, 0x3b, 0x26, 0x4b, 0x76,
	0x33, 0x65, 0x33, 0x95, 0xed, 0x03, 0xf8, 0xde, 0x92, 0x25, 0xf8, 0xf7, 0xd0, 0xeb, 0x37, 0x3f,
	0x98, 0x8f, 0x4f, 0x5c, 0x74, 0x4c, 0x5e, 0x56, 0xb8, 0xdc, 0x79, 0x62, 0xca, 0x4e, 0x5c, 0x57,
	0x75, 0xf3, 0x8f, 0x44, 0x5c, 0x71, 0x58, 0xaf, 0xd6, 0xad, 0xc6, 0xf0, 0x22, 0xed, 0x73, 0xbc,
	0xdd, 0x89, 0x09, 0xd8, 0xa6, 0xb4, 0xa4, 0x24, 0x1b, 0x6a, 0xa7, 0xc9, 0x2c, 0x59, 0x4c, 0x8a,
	0x31, 0xd8, 0x66, 0x33, 0xa0, 0xec, 0x56, 0x8c, 0x39, 0x7c, 0x2f, 0x46, 0x71, 0x21, 0x38, 0x7c,
	0x0d, 0xee, 0x45, 0xc6, 0xa1, 0x7c, 0x97, 0xd6, 0xe8, 0x08, 0x4a, 0x45, 0x1a, 0xa7, 0xe9, 0x2c,
	0x59, 0x9c, 0x16, 0x97, 0xdc, 0xff, 0x1a, 0x0e, 0x4f, 0xa4, 0x71, 0x1e, 0xc4, 0x35, 0x87, 0x95,
	0x64, 0xfc, 0x2d, 0x72, 0x23, 0x44, 0x2f, 0xd2, 0xee, 0x1d, 0xa0, 0x8f, 0x1a, 0x69, 0x71, 0x0e,
	0xb6, 0xd9, 0x46, 0xf0, 0xcf, 0x73, 0x74, 0xd4, 0x33, 0xfd, 0xeb, 0xf9, 0xb8, 0x7d, 0xdd, 0x54,
	0x86, 0xeb, 0x3d, 0x2c, 0x15, 0xb9, 0xbc, 0x3e, 0xec, 0xd0, 0x5b, 0xd4, 0x15, 0xfa, 0xfc, 0x4d,
	0x82, 0x37, 0x2a, 0x57, 0xe4, 0x1c, 0xb5, 0xf9, 0x00, 0xc1, 0x36, 0x43, 0xc7, 0xbc, 0x3e, 0x80,
	0x37, 0xfa, 0x07, 0xe8, 0x8b, 0xc2, 0x59, 0xcc, 0xfb, 0xf0, 0x19, 0x00, 0x00, 0xff, 0xff, 0x7c,
	0x44, 0x8a, 0xd2, 0x6f, 0x01, 0x00, 0x00,
}
