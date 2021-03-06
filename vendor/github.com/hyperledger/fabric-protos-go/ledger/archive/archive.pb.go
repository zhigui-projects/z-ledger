// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ledger/archive/archive.proto

package archive

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

type ArchiveMetaInfo struct {
	LastSentFileSuffix    int32            `protobuf:"varint,1,opt,name=last_sent_file_suffix,json=lastSentFileSuffix,proto3" json:"last_sent_file_suffix,omitempty"`
	LastArchiveFileSuffix int32            `protobuf:"varint,2,opt,name=last_archive_file_suffix,json=lastArchiveFileSuffix,proto3" json:"last_archive_file_suffix,omitempty"`
	FileProofs            map[int32]string `protobuf:"bytes,3,rep,name=file_proofs,json=fileProofs,proto3" json:"file_proofs,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ChannelId             string           `protobuf:"bytes,4,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}         `json:"-"`
	XXX_unrecognized      []byte           `json:"-"`
	XXX_sizecache         int32            `json:"-"`
}

func (m *ArchiveMetaInfo) Reset()         { *m = ArchiveMetaInfo{} }
func (m *ArchiveMetaInfo) String() string { return proto.CompactTextString(m) }
func (*ArchiveMetaInfo) ProtoMessage()    {}
func (*ArchiveMetaInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_8515fd65c6f6c1d3, []int{0}
}

func (m *ArchiveMetaInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ArchiveMetaInfo.Unmarshal(m, b)
}
func (m *ArchiveMetaInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ArchiveMetaInfo.Marshal(b, m, deterministic)
}
func (m *ArchiveMetaInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArchiveMetaInfo.Merge(m, src)
}
func (m *ArchiveMetaInfo) XXX_Size() int {
	return xxx_messageInfo_ArchiveMetaInfo.Size(m)
}
func (m *ArchiveMetaInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ArchiveMetaInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ArchiveMetaInfo proto.InternalMessageInfo

func (m *ArchiveMetaInfo) GetLastSentFileSuffix() int32 {
	if m != nil {
		return m.LastSentFileSuffix
	}
	return 0
}

func (m *ArchiveMetaInfo) GetLastArchiveFileSuffix() int32 {
	if m != nil {
		return m.LastArchiveFileSuffix
	}
	return 0
}

func (m *ArchiveMetaInfo) GetFileProofs() map[int32]string {
	if m != nil {
		return m.FileProofs
	}
	return nil
}

func (m *ArchiveMetaInfo) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func init() {
	proto.RegisterType((*ArchiveMetaInfo)(nil), "archive.ArchiveMetaInfo")
	proto.RegisterMapType((map[int32]string)(nil), "archive.ArchiveMetaInfo.FileProofsEntry")
}

func init() { proto.RegisterFile("ledger/archive/archive.proto", fileDescriptor_8515fd65c6f6c1d3) }

var fileDescriptor_8515fd65c6f6c1d3 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0xe9, 0xe6, 0x94, 0xbd, 0x1d, 0x26, 0x41, 0xa1, 0x88, 0xc2, 0xf0, 0xd4, 0x83, 0x4b,
	0x51, 0x41, 0x45, 0xf0, 0xa0, 0xa0, 0xb0, 0x83, 0x20, 0xdd, 0xcd, 0x4b, 0xc9, 0xba, 0x97, 0x36,
	0x18, 0x93, 0x91, 0xa6, 0xc3, 0x7e, 0x21, 0x3f, 0xa7, 0x34, 0xc9, 0x64, 0xdb, 0x29, 0xef, 0xbd,
	0xff, 0xfb, 0xbd, 0xfc, 0xf3, 0x02, 0xe7, 0x12, 0x97, 0x25, 0x9a, 0x94, 0x99, 0xa2, 0x12, 0x6b,
	0xdc, 0x9c, 0x74, 0x65, 0xb4, 0xd5, 0xe4, 0x28, 0xa4, 0x97, 0xbf, 0x3d, 0x18, 0x3f, 0xfb, 0xf8,
	0x1d, 0x2d, 0x9b, 0x29, 0xae, 0xc9, 0x35, 0x9c, 0x4a, 0x56, 0xdb, 0xbc, 0x46, 0x65, 0x73, 0x2e,
	0x24, 0xe6, 0x75, 0xc3, 0xb9, 0xf8, 0x89, 0xa3, 0x49, 0x94, 0x0c, 0x32, 0xd2, 0x89, 0x73, 0x54,
	0xf6, 0x4d, 0x48, 0x9c, 0x3b, 0x85, 0xdc, 0x43, 0xec, 0x90, 0x30, 0x76, 0x87, 0xea, 0x39, 0xca,
	0x8d, 0x0c, 0x37, 0x6d, 0x81, 0x33, 0x18, 0xb9, 0xde, 0x95, 0xd1, 0x9a, 0xd7, 0x71, 0x7f, 0xd2,
	0x4f, 0x46, 0x37, 0x09, 0xdd, 0xb8, 0xdd, 0xb3, 0x46, 0x3b, 0xf2, 0xc3, 0xb5, 0xbe, 0x2a, 0x6b,
	0xda, 0x0c, 0xf8, 0x7f, 0x81, 0x5c, 0x00, 0x14, 0x15, 0x53, 0x0a, 0x65, 0x2e, 0x96, 0xf1, 0xc1,
	0x24, 0x4a, 0x86, 0xd9, 0x30, 0x54, 0x66, 0xcb, 0xb3, 0x27, 0x18, 0xef, 0xd1, 0xe4, 0x18, 0xfa,
	0x5f, 0xd8, 0x86, 0x67, 0x75, 0x21, 0x39, 0x81, 0xc1, 0x9a, 0xc9, 0x06, 0x9d, 0xe9, 0x61, 0xe6,
	0x93, 0xc7, 0xde, 0x43, 0xf4, 0xc2, 0xe1, 0x4a, 0x9b, 0x92, 0x56, 0xed, 0x0a, 0x8d, 0x5f, 0x2d,
	0xe5, 0x6c, 0x61, 0x44, 0xe1, 0x37, 0x5a, 0xd3, 0x50, 0x0c, 0xce, 0x3f, 0xef, 0x4a, 0x61, 0xab,
	0x66, 0x41, 0x0b, 0xfd, 0x9d, 0x6e, 0x41, 0xa9, 0x87, 0xa6, 0x1e, 0x9a, 0x96, 0x3a, 0xdd, 0xfd,
	0xa7, 0xc5, 0xa1, 0x53, 0x6e, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xff, 0x5f, 0x76, 0xd9, 0xc0,
	0x01, 0x00, 0x00,
}
