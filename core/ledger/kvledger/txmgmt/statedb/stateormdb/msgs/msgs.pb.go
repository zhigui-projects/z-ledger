// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msgs.proto

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

type VersionFieldProto struct {
	VersionBytes         []byte   `protobuf:"bytes,1,opt,name=version_bytes,json=versionBytes,proto3" json:"version_bytes,omitempty"`
	Metadata             []byte   `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionFieldProto) Reset()         { *m = VersionFieldProto{} }
func (m *VersionFieldProto) String() string { return proto.CompactTextString(m) }
func (*VersionFieldProto) ProtoMessage()    {}
func (*VersionFieldProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_952909143bb80d72, []int{0}
}

func (m *VersionFieldProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionFieldProto.Unmarshal(m, b)
}
func (m *VersionFieldProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionFieldProto.Marshal(b, m, deterministic)
}
func (m *VersionFieldProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionFieldProto.Merge(m, src)
}
func (m *VersionFieldProto) XXX_Size() int {
	return xxx_messageInfo_VersionFieldProto.Size(m)
}
func (m *VersionFieldProto) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionFieldProto.DiscardUnknown(m)
}

var xxx_messageInfo_VersionFieldProto proto.InternalMessageInfo

func (m *VersionFieldProto) GetVersionBytes() []byte {
	if m != nil {
		return m.VersionBytes
	}
	return nil
}

func (m *VersionFieldProto) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterType((*VersionFieldProto)(nil), "msgs.VersionFieldProto")
}

func init() { proto.RegisterFile("msgs.proto", fileDescriptor_952909143bb80d72) }

var fileDescriptor_952909143bb80d72 = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0xce, 0x3f, 0x0b, 0xc2, 0x30,
	0x10, 0x05, 0x70, 0x2a, 0x22, 0x12, 0xea, 0x60, 0xa7, 0xe2, 0x24, 0xba, 0x38, 0x35, 0x83, 0xdf,
	0xa0, 0x83, 0xb3, 0xff, 0x70, 0x70, 0x91, 0xa4, 0x39, 0xd3, 0x60, 0xe3, 0x95, 0xcb, 0x59, 0xec,
	0xb7, 0x97, 0xfe, 0xc1, 0xe9, 0xde, 0xfd, 0xde, 0xf2, 0x84, 0xf0, 0xc1, 0x86, 0xac, 0x26, 0x64,
	0x4c, 0xa6, 0x5d, 0xde, 0x5c, 0xc5, 0xf2, 0x06, 0x14, 0x1c, 0xbe, 0x0f, 0x0e, 0x2a, 0x73, 0xec,
	0xab, 0xad, 0x58, 0x34, 0x03, 0x3e, 0x74, 0xcb, 0x10, 0xd2, 0x68, 0x1d, 0xed, 0xe2, 0x73, 0x3c,
	0x62, 0xde, 0x59, 0xb2, 0x12, 0x73, 0x0f, 0xac, 0x8c, 0x62, 0x95, 0x4e, 0xfa, 0xfe, 0xff, 0xe7,
	0x97, 0xfb, 0xc9, 0x3a, 0x2e, 0x3f, 0x3a, 0x2b, 0xd0, 0xcb, 0xb2, 0xad, 0x81, 0x2a, 0x30, 0x16,
	0x48, 0x3e, 0x95, 0x26, 0x57, 0xc8, 0x02, 0x09, 0xe4, 0x48, 0xaf, 0x66, 0x0c, 0xfc, 0xf5, 0xd6,
	0xb3, 0x0c, 0xac, 0x18, 0x8c, 0x1e, 0x2e, 0x92, 0x37, 0x5a, 0x76, 0x53, 0xf5, 0xac, 0xdf, 0xbd,
	0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0xd2, 0xea, 0xb4, 0x94, 0xc5, 0x00, 0x00, 0x00,
}
