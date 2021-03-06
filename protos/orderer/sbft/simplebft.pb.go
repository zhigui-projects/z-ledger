// Code generated by protoc-gen-go. DO NOT EDIT.
// source: orderer/sbft/simplebft.proto

package sbft // import "github.com/hyperledger/fabric/protos/orderer/sbft"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MultiChainMsg struct {
	ChainID              string   `protobuf:"bytes,1,opt,name=chainID,proto3" json:"chainID,omitempty"`
	Msg                  *Msg     `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MultiChainMsg) Reset()         { *m = MultiChainMsg{} }
func (m *MultiChainMsg) String() string { return proto.CompactTextString(m) }
func (*MultiChainMsg) ProtoMessage()    {}
func (*MultiChainMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{0}
}
func (m *MultiChainMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MultiChainMsg.Unmarshal(m, b)
}
func (m *MultiChainMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MultiChainMsg.Marshal(b, m, deterministic)
}
func (dst *MultiChainMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiChainMsg.Merge(dst, src)
}
func (m *MultiChainMsg) XXX_Size() int {
	return xxx_messageInfo_MultiChainMsg.Size(m)
}
func (m *MultiChainMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiChainMsg.DiscardUnknown(m)
}

var xxx_messageInfo_MultiChainMsg proto.InternalMessageInfo

func (m *MultiChainMsg) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

func (m *MultiChainMsg) GetMsg() *Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

type Msg struct {
	// Types that are valid to be assigned to Type:
	//	*Msg_Request
	//	*Msg_Preprepare
	//	*Msg_Prepare
	//	*Msg_Commit
	//	*Msg_ViewChange
	//	*Msg_NewView
	//	*Msg_Checkpoint
	//	*Msg_Hello
	Type                 isMsg_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Msg) Reset()         { *m = Msg{} }
func (m *Msg) String() string { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()    {}
func (*Msg) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{1}
}
func (m *Msg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Msg.Unmarshal(m, b)
}
func (m *Msg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Msg.Marshal(b, m, deterministic)
}
func (dst *Msg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Msg.Merge(dst, src)
}
func (m *Msg) XXX_Size() int {
	return xxx_messageInfo_Msg.Size(m)
}
func (m *Msg) XXX_DiscardUnknown() {
	xxx_messageInfo_Msg.DiscardUnknown(m)
}

var xxx_messageInfo_Msg proto.InternalMessageInfo

type isMsg_Type interface {
	isMsg_Type()
}

type Msg_Request struct {
	Request *Request `protobuf:"bytes,1,opt,name=request,proto3,oneof"`
}

type Msg_Preprepare struct {
	Preprepare *Preprepare `protobuf:"bytes,2,opt,name=preprepare,proto3,oneof"`
}

type Msg_Prepare struct {
	Prepare *Subject `protobuf:"bytes,3,opt,name=prepare,proto3,oneof"`
}

type Msg_Commit struct {
	Commit *Subject `protobuf:"bytes,4,opt,name=commit,proto3,oneof"`
}

type Msg_ViewChange struct {
	ViewChange *Signed `protobuf:"bytes,5,opt,name=view_change,json=viewChange,proto3,oneof"`
}

type Msg_NewView struct {
	NewView *NewView `protobuf:"bytes,6,opt,name=new_view,json=newView,proto3,oneof"`
}

type Msg_Checkpoint struct {
	Checkpoint *Checkpoint `protobuf:"bytes,7,opt,name=checkpoint,proto3,oneof"`
}

type Msg_Hello struct {
	Hello *Hello `protobuf:"bytes,8,opt,name=hello,proto3,oneof"`
}

func (*Msg_Request) isMsg_Type() {}

func (*Msg_Preprepare) isMsg_Type() {}

func (*Msg_Prepare) isMsg_Type() {}

func (*Msg_Commit) isMsg_Type() {}

func (*Msg_ViewChange) isMsg_Type() {}

func (*Msg_NewView) isMsg_Type() {}

func (*Msg_Checkpoint) isMsg_Type() {}

func (*Msg_Hello) isMsg_Type() {}

func (m *Msg) GetType() isMsg_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Msg) GetRequest() *Request {
	if x, ok := m.GetType().(*Msg_Request); ok {
		return x.Request
	}
	return nil
}

func (m *Msg) GetPreprepare() *Preprepare {
	if x, ok := m.GetType().(*Msg_Preprepare); ok {
		return x.Preprepare
	}
	return nil
}

func (m *Msg) GetPrepare() *Subject {
	if x, ok := m.GetType().(*Msg_Prepare); ok {
		return x.Prepare
	}
	return nil
}

func (m *Msg) GetCommit() *Subject {
	if x, ok := m.GetType().(*Msg_Commit); ok {
		return x.Commit
	}
	return nil
}

func (m *Msg) GetViewChange() *Signed {
	if x, ok := m.GetType().(*Msg_ViewChange); ok {
		return x.ViewChange
	}
	return nil
}

func (m *Msg) GetNewView() *NewView {
	if x, ok := m.GetType().(*Msg_NewView); ok {
		return x.NewView
	}
	return nil
}

func (m *Msg) GetCheckpoint() *Checkpoint {
	if x, ok := m.GetType().(*Msg_Checkpoint); ok {
		return x.Checkpoint
	}
	return nil
}

func (m *Msg) GetHello() *Hello {
	if x, ok := m.GetType().(*Msg_Hello); ok {
		return x.Hello
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Msg) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Msg_OneofMarshaler, _Msg_OneofUnmarshaler, _Msg_OneofSizer, []interface{}{
		(*Msg_Request)(nil),
		(*Msg_Preprepare)(nil),
		(*Msg_Prepare)(nil),
		(*Msg_Commit)(nil),
		(*Msg_ViewChange)(nil),
		(*Msg_NewView)(nil),
		(*Msg_Checkpoint)(nil),
		(*Msg_Hello)(nil),
	}
}

func _Msg_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Msg)
	// type
	switch x := m.Type.(type) {
	case *Msg_Request:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Request); err != nil {
			return err
		}
	case *Msg_Preprepare:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Preprepare); err != nil {
			return err
		}
	case *Msg_Prepare:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Prepare); err != nil {
			return err
		}
	case *Msg_Commit:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Commit); err != nil {
			return err
		}
	case *Msg_ViewChange:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ViewChange); err != nil {
			return err
		}
	case *Msg_NewView:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.NewView); err != nil {
			return err
		}
	case *Msg_Checkpoint:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Checkpoint); err != nil {
			return err
		}
	case *Msg_Hello:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Hello); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Msg.Type has unexpected type %T", x)
	}
	return nil
}

func _Msg_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Msg)
	switch tag {
	case 1: // type.request
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Request)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Request{msg}
		return true, err
	case 2: // type.preprepare
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Preprepare)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Preprepare{msg}
		return true, err
	case 3: // type.prepare
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Subject)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Prepare{msg}
		return true, err
	case 4: // type.commit
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Subject)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Commit{msg}
		return true, err
	case 5: // type.view_change
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Signed)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_ViewChange{msg}
		return true, err
	case 6: // type.new_view
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(NewView)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_NewView{msg}
		return true, err
	case 7: // type.checkpoint
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Checkpoint)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Checkpoint{msg}
		return true, err
	case 8: // type.hello
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Hello)
		err := b.DecodeMessage(msg)
		m.Type = &Msg_Hello{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Msg_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Msg)
	// type
	switch x := m.Type.(type) {
	case *Msg_Request:
		s := proto.Size(x.Request)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Preprepare:
		s := proto.Size(x.Preprepare)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Prepare:
		s := proto.Size(x.Prepare)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Commit:
		s := proto.Size(x.Commit)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_ViewChange:
		s := proto.Size(x.ViewChange)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_NewView:
		s := proto.Size(x.NewView)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Checkpoint:
		s := proto.Size(x.Checkpoint)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Msg_Hello:
		s := proto.Size(x.Hello)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Request struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{2}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (dst *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(dst, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type SeqView struct {
	View                 uint64   `protobuf:"varint,1,opt,name=view,proto3" json:"view,omitempty"`
	Seq                  uint64   `protobuf:"varint,2,opt,name=seq,proto3" json:"seq,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SeqView) Reset()         { *m = SeqView{} }
func (m *SeqView) String() string { return proto.CompactTextString(m) }
func (*SeqView) ProtoMessage()    {}
func (*SeqView) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{3}
}
func (m *SeqView) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeqView.Unmarshal(m, b)
}
func (m *SeqView) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeqView.Marshal(b, m, deterministic)
}
func (dst *SeqView) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeqView.Merge(dst, src)
}
func (m *SeqView) XXX_Size() int {
	return xxx_messageInfo_SeqView.Size(m)
}
func (m *SeqView) XXX_DiscardUnknown() {
	xxx_messageInfo_SeqView.DiscardUnknown(m)
}

var xxx_messageInfo_SeqView proto.InternalMessageInfo

func (m *SeqView) GetView() uint64 {
	if m != nil {
		return m.View
	}
	return 0
}

func (m *SeqView) GetSeq() uint64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

type BatchHeader struct {
	Seq                  uint64   `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
	PrevHash             []byte   `protobuf:"bytes,2,opt,name=prev_hash,json=prevHash,proto3" json:"prev_hash,omitempty"`
	DataHash             []byte   `protobuf:"bytes,3,opt,name=data_hash,json=dataHash,proto3" json:"data_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchHeader) Reset()         { *m = BatchHeader{} }
func (m *BatchHeader) String() string { return proto.CompactTextString(m) }
func (*BatchHeader) ProtoMessage()    {}
func (*BatchHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{4}
}
func (m *BatchHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatchHeader.Unmarshal(m, b)
}
func (m *BatchHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatchHeader.Marshal(b, m, deterministic)
}
func (dst *BatchHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchHeader.Merge(dst, src)
}
func (m *BatchHeader) XXX_Size() int {
	return xxx_messageInfo_BatchHeader.Size(m)
}
func (m *BatchHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BatchHeader proto.InternalMessageInfo

func (m *BatchHeader) GetSeq() uint64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *BatchHeader) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *BatchHeader) GetDataHash() []byte {
	if m != nil {
		return m.DataHash
	}
	return nil
}

type Batch struct {
	Header               []byte            `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Payloads             [][]byte          `protobuf:"bytes,2,rep,name=payloads,proto3" json:"payloads,omitempty"`
	Signatures           map[uint64][]byte `protobuf:"bytes,3,rep,name=signatures,proto3" json:"signatures,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Batch) Reset()         { *m = Batch{} }
func (m *Batch) String() string { return proto.CompactTextString(m) }
func (*Batch) ProtoMessage()    {}
func (*Batch) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{5}
}
func (m *Batch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Batch.Unmarshal(m, b)
}
func (m *Batch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Batch.Marshal(b, m, deterministic)
}
func (dst *Batch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Batch.Merge(dst, src)
}
func (m *Batch) XXX_Size() int {
	return xxx_messageInfo_Batch.Size(m)
}
func (m *Batch) XXX_DiscardUnknown() {
	xxx_messageInfo_Batch.DiscardUnknown(m)
}

var xxx_messageInfo_Batch proto.InternalMessageInfo

func (m *Batch) GetHeader() []byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Batch) GetPayloads() [][]byte {
	if m != nil {
		return m.Payloads
	}
	return nil
}

func (m *Batch) GetSignatures() map[uint64][]byte {
	if m != nil {
		return m.Signatures
	}
	return nil
}

type Preprepare struct {
	Seq                  *SeqView `protobuf:"bytes,1,opt,name=seq,proto3" json:"seq,omitempty"`
	Batch                *Batch   `protobuf:"bytes,2,opt,name=batch,proto3" json:"batch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Preprepare) Reset()         { *m = Preprepare{} }
func (m *Preprepare) String() string { return proto.CompactTextString(m) }
func (*Preprepare) ProtoMessage()    {}
func (*Preprepare) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{6}
}
func (m *Preprepare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Preprepare.Unmarshal(m, b)
}
func (m *Preprepare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Preprepare.Marshal(b, m, deterministic)
}
func (dst *Preprepare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Preprepare.Merge(dst, src)
}
func (m *Preprepare) XXX_Size() int {
	return xxx_messageInfo_Preprepare.Size(m)
}
func (m *Preprepare) XXX_DiscardUnknown() {
	xxx_messageInfo_Preprepare.DiscardUnknown(m)
}

var xxx_messageInfo_Preprepare proto.InternalMessageInfo

func (m *Preprepare) GetSeq() *SeqView {
	if m != nil {
		return m.Seq
	}
	return nil
}

func (m *Preprepare) GetBatch() *Batch {
	if m != nil {
		return m.Batch
	}
	return nil
}

type Subject struct {
	Seq                  *SeqView `protobuf:"bytes,1,opt,name=seq,proto3" json:"seq,omitempty"`
	Digest               []byte   `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Subject) Reset()         { *m = Subject{} }
func (m *Subject) String() string { return proto.CompactTextString(m) }
func (*Subject) ProtoMessage()    {}
func (*Subject) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{7}
}
func (m *Subject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subject.Unmarshal(m, b)
}
func (m *Subject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subject.Marshal(b, m, deterministic)
}
func (dst *Subject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subject.Merge(dst, src)
}
func (m *Subject) XXX_Size() int {
	return xxx_messageInfo_Subject.Size(m)
}
func (m *Subject) XXX_DiscardUnknown() {
	xxx_messageInfo_Subject.DiscardUnknown(m)
}

var xxx_messageInfo_Subject proto.InternalMessageInfo

func (m *Subject) GetSeq() *SeqView {
	if m != nil {
		return m.Seq
	}
	return nil
}

func (m *Subject) GetDigest() []byte {
	if m != nil {
		return m.Digest
	}
	return nil
}

type ViewChange struct {
	View                 uint64     `protobuf:"varint,1,opt,name=view,proto3" json:"view,omitempty"`
	Pset                 []*Subject `protobuf:"bytes,2,rep,name=pset,proto3" json:"pset,omitempty"`
	Qset                 []*Subject `protobuf:"bytes,3,rep,name=qset,proto3" json:"qset,omitempty"`
	Checkpoint           *Batch     `protobuf:"bytes,4,opt,name=checkpoint,proto3" json:"checkpoint,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ViewChange) Reset()         { *m = ViewChange{} }
func (m *ViewChange) String() string { return proto.CompactTextString(m) }
func (*ViewChange) ProtoMessage()    {}
func (*ViewChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{8}
}
func (m *ViewChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ViewChange.Unmarshal(m, b)
}
func (m *ViewChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ViewChange.Marshal(b, m, deterministic)
}
func (dst *ViewChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ViewChange.Merge(dst, src)
}
func (m *ViewChange) XXX_Size() int {
	return xxx_messageInfo_ViewChange.Size(m)
}
func (m *ViewChange) XXX_DiscardUnknown() {
	xxx_messageInfo_ViewChange.DiscardUnknown(m)
}

var xxx_messageInfo_ViewChange proto.InternalMessageInfo

func (m *ViewChange) GetView() uint64 {
	if m != nil {
		return m.View
	}
	return 0
}

func (m *ViewChange) GetPset() []*Subject {
	if m != nil {
		return m.Pset
	}
	return nil
}

func (m *ViewChange) GetQset() []*Subject {
	if m != nil {
		return m.Qset
	}
	return nil
}

func (m *ViewChange) GetCheckpoint() *Batch {
	if m != nil {
		return m.Checkpoint
	}
	return nil
}

type Signed struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Signed) Reset()         { *m = Signed{} }
func (m *Signed) String() string { return proto.CompactTextString(m) }
func (*Signed) ProtoMessage()    {}
func (*Signed) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{9}
}
func (m *Signed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Signed.Unmarshal(m, b)
}
func (m *Signed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Signed.Marshal(b, m, deterministic)
}
func (dst *Signed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signed.Merge(dst, src)
}
func (m *Signed) XXX_Size() int {
	return xxx_messageInfo_Signed.Size(m)
}
func (m *Signed) XXX_DiscardUnknown() {
	xxx_messageInfo_Signed.DiscardUnknown(m)
}

var xxx_messageInfo_Signed proto.InternalMessageInfo

func (m *Signed) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Signed) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type NewView struct {
	View                 uint64             `protobuf:"varint,1,opt,name=view,proto3" json:"view,omitempty"`
	Vset                 map[uint64]*Signed `protobuf:"bytes,2,rep,name=vset,proto3" json:"vset,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Xset                 *Subject           `protobuf:"bytes,3,opt,name=xset,proto3" json:"xset,omitempty"`
	Batch                *Batch             `protobuf:"bytes,4,opt,name=batch,proto3" json:"batch,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *NewView) Reset()         { *m = NewView{} }
func (m *NewView) String() string { return proto.CompactTextString(m) }
func (*NewView) ProtoMessage()    {}
func (*NewView) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{10}
}
func (m *NewView) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewView.Unmarshal(m, b)
}
func (m *NewView) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewView.Marshal(b, m, deterministic)
}
func (dst *NewView) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewView.Merge(dst, src)
}
func (m *NewView) XXX_Size() int {
	return xxx_messageInfo_NewView.Size(m)
}
func (m *NewView) XXX_DiscardUnknown() {
	xxx_messageInfo_NewView.DiscardUnknown(m)
}

var xxx_messageInfo_NewView proto.InternalMessageInfo

func (m *NewView) GetView() uint64 {
	if m != nil {
		return m.View
	}
	return 0
}

func (m *NewView) GetVset() map[uint64]*Signed {
	if m != nil {
		return m.Vset
	}
	return nil
}

func (m *NewView) GetXset() *Subject {
	if m != nil {
		return m.Xset
	}
	return nil
}

func (m *NewView) GetBatch() *Batch {
	if m != nil {
		return m.Batch
	}
	return nil
}

type Checkpoint struct {
	Seq                  uint64   `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
	Digest               []byte   `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	Signature            []byte   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Checkpoint) Reset()         { *m = Checkpoint{} }
func (m *Checkpoint) String() string { return proto.CompactTextString(m) }
func (*Checkpoint) ProtoMessage()    {}
func (*Checkpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{11}
}
func (m *Checkpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Checkpoint.Unmarshal(m, b)
}
func (m *Checkpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Checkpoint.Marshal(b, m, deterministic)
}
func (dst *Checkpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Checkpoint.Merge(dst, src)
}
func (m *Checkpoint) XXX_Size() int {
	return xxx_messageInfo_Checkpoint.Size(m)
}
func (m *Checkpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Checkpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Checkpoint proto.InternalMessageInfo

func (m *Checkpoint) GetSeq() uint64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *Checkpoint) GetDigest() []byte {
	if m != nil {
		return m.Digest
	}
	return nil
}

func (m *Checkpoint) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type Hello struct {
	Batch                *Batch   `protobuf:"bytes,1,opt,name=batch,proto3" json:"batch,omitempty"`
	NewView              *NewView `protobuf:"bytes,2,opt,name=new_view,json=newView,proto3" json:"new_view,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hello) Reset()         { *m = Hello{} }
func (m *Hello) String() string { return proto.CompactTextString(m) }
func (*Hello) ProtoMessage()    {}
func (*Hello) Descriptor() ([]byte, []int) {
	return fileDescriptor_simplebft_41cd103f09603772, []int{12}
}
func (m *Hello) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hello.Unmarshal(m, b)
}
func (m *Hello) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hello.Marshal(b, m, deterministic)
}
func (dst *Hello) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hello.Merge(dst, src)
}
func (m *Hello) XXX_Size() int {
	return xxx_messageInfo_Hello.Size(m)
}
func (m *Hello) XXX_DiscardUnknown() {
	xxx_messageInfo_Hello.DiscardUnknown(m)
}

var xxx_messageInfo_Hello proto.InternalMessageInfo

func (m *Hello) GetBatch() *Batch {
	if m != nil {
		return m.Batch
	}
	return nil
}

func (m *Hello) GetNewView() *NewView {
	if m != nil {
		return m.NewView
	}
	return nil
}

func init() {
	proto.RegisterType((*MultiChainMsg)(nil), "simplebft.MultiChainMsg")
	proto.RegisterType((*Msg)(nil), "simplebft.Msg")
	proto.RegisterType((*Request)(nil), "simplebft.Request")
	proto.RegisterType((*SeqView)(nil), "simplebft.SeqView")
	proto.RegisterType((*BatchHeader)(nil), "simplebft.BatchHeader")
	proto.RegisterType((*Batch)(nil), "simplebft.Batch")
	proto.RegisterMapType((map[uint64][]byte)(nil), "simplebft.Batch.SignaturesEntry")
	proto.RegisterType((*Preprepare)(nil), "simplebft.Preprepare")
	proto.RegisterType((*Subject)(nil), "simplebft.Subject")
	proto.RegisterType((*ViewChange)(nil), "simplebft.ViewChange")
	proto.RegisterType((*Signed)(nil), "simplebft.Signed")
	proto.RegisterType((*NewView)(nil), "simplebft.NewView")
	proto.RegisterMapType((map[uint64]*Signed)(nil), "simplebft.NewView.VsetEntry")
	proto.RegisterType((*Checkpoint)(nil), "simplebft.Checkpoint")
	proto.RegisterType((*Hello)(nil), "simplebft.Hello")
}

func init() {
	proto.RegisterFile("orderer/sbft/simplebft.proto", fileDescriptor_simplebft_41cd103f09603772)
}

var fileDescriptor_simplebft_41cd103f09603772 = []byte{
	// 747 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x55, 0xdb, 0x6e, 0xd3, 0x4c,
	0x10, 0xae, 0x63, 0xe7, 0x34, 0xe9, 0xff, 0xff, 0xfd, 0x57, 0x50, 0x59, 0xa5, 0x17, 0x91, 0x41,
	0x25, 0x17, 0x10, 0x97, 0x16, 0x09, 0x54, 0x09, 0x09, 0xb5, 0x20, 0x02, 0xa8, 0x08, 0xb9, 0xa8,
	0x12, 0xbd, 0xa0, 0x72, 0x9c, 0xa9, 0x6d, 0x9a, 0xd8, 0xce, 0xee, 0x26, 0x69, 0x5e, 0x86, 0x0b,
	0x9e, 0x83, 0x37, 0xe2, 0x25, 0xd0, 0x1e, 0x12, 0xbb, 0x39, 0x54, 0x48, 0xbe, 0xd8, 0x9d, 0xf9,
	0xe6, 0xf0, 0xcd, 0x61, 0x0d, 0xbb, 0x29, 0xed, 0x21, 0x45, 0xea, 0xb2, 0xee, 0x15, 0x77, 0x59,
	0x3c, 0xc8, 0xfa, 0xd8, 0xbd, 0xe2, 0xed, 0x8c, 0xa6, 0x3c, 0x25, 0xf5, 0xb9, 0xc0, 0xf9, 0x08,
	0xff, 0x9c, 0x8e, 0xfa, 0x3c, 0x3e, 0x89, 0xfc, 0x38, 0x39, 0x65, 0x21, 0xb1, 0xa1, 0x1a, 0x88,
	0xf3, 0xfb, 0x37, 0xb6, 0xd1, 0x34, 0x5a, 0x75, 0x6f, 0x76, 0x25, 0x4d, 0x30, 0x07, 0x2c, 0xb4,
	0x4b, 0x4d, 0xa3, 0xd5, 0x38, 0xf8, 0xb7, 0x9d, 0x3b, 0x3d, 0x65, 0xa1, 0x27, 0x54, 0xce, 0x0f,
	0x13, 0x4c, 0xe1, 0xa3, 0x0d, 0x55, 0x8a, 0xc3, 0x11, 0x32, 0x2e, 0x7d, 0x34, 0x0e, 0x48, 0x01,
	0xed, 0x29, 0x4d, 0x67, 0xc3, 0x9b, 0x81, 0xc8, 0x0b, 0x80, 0x8c, 0xa2, 0xf8, 0x7c, 0x8a, 0x3a,
	0xc0, 0xfd, 0x82, 0xc9, 0xe7, 0xb9, 0xb2, 0xb3, 0xe1, 0x15, 0xa0, 0x22, 0xd0, 0xcc, 0xca, 0x5c,
	0x0a, 0x74, 0x36, 0xea, 0x7e, 0xc7, 0x40, 0x06, 0x9a, 0xe1, 0x9f, 0x40, 0x25, 0x48, 0x07, 0x83,
	0x98, 0xdb, 0xd6, 0x1d, 0x70, 0x8d, 0x21, 0xcf, 0xa1, 0x31, 0x8e, 0x71, 0x72, 0x19, 0x44, 0x7e,
	0x12, 0xa2, 0x5d, 0x96, 0x26, 0xff, 0x17, 0x4d, 0xe2, 0x30, 0xc1, 0x9e, 0xc8, 0x49, 0xe0, 0x4e,
	0x24, 0x8c, 0xb8, 0x50, 0x4b, 0x70, 0x72, 0x29, 0x24, 0x76, 0x65, 0x29, 0xca, 0x27, 0x9c, 0x9c,
	0xc7, 0x38, 0x11, 0x49, 0x25, 0xea, 0x28, 0xd8, 0x07, 0x11, 0x06, 0xd7, 0x59, 0x1a, 0x27, 0xdc,
	0xae, 0x2e, 0xb1, 0x3f, 0x99, 0x2b, 0x45, 0xa4, 0x1c, 0x4a, 0x5a, 0x50, 0x8e, 0xb0, 0xdf, 0x4f,
	0xed, 0x9a, 0xb4, 0xd9, 0x2a, 0xd8, 0x74, 0x84, 0xbc, 0xb3, 0xe1, 0x29, 0xc0, 0x71, 0x05, 0x2c,
	0x3e, 0xcd, 0xd0, 0x79, 0x08, 0x55, 0x5d, 0x7e, 0xd1, 0xe7, 0xcc, 0x9f, 0xf6, 0x53, 0xbf, 0x27,
	0x7b, 0xb4, 0xe9, 0xcd, 0xae, 0x8e, 0x0b, 0xd5, 0x33, 0x1c, 0xca, 0xd4, 0x08, 0x58, 0x92, 0x87,
	0x40, 0x58, 0x9e, 0x3c, 0x93, 0x2d, 0x30, 0x19, 0x0e, 0x65, 0x97, 0x2c, 0x4f, 0x1c, 0x9d, 0xaf,
	0xd0, 0x38, 0xf6, 0x79, 0x10, 0x75, 0xd0, 0xef, 0x21, 0x9d, 0x01, 0x8c, 0x39, 0x80, 0x3c, 0x80,
	0x7a, 0x46, 0x71, 0x7c, 0x19, 0xf9, 0x2c, 0x92, 0x86, 0x9b, 0x5e, 0x4d, 0x08, 0x3a, 0x3e, 0x8b,
	0x84, 0xb2, 0xe7, 0x73, 0x5f, 0x29, 0x4d, 0xa5, 0x14, 0x02, 0xa1, 0x74, 0x7e, 0x19, 0x50, 0x96,
	0xbe, 0xc9, 0x36, 0x54, 0x22, 0xe9, 0x5f, 0xa7, 0xab, 0x6f, 0x64, 0x07, 0x6a, 0x3a, 0x71, 0x66,
	0x97, 0x9a, 0xa6, 0x74, 0xad, 0xef, 0xe4, 0x35, 0x00, 0x8b, 0xc3, 0xc4, 0xe7, 0x23, 0x8a, 0xcc,
	0x36, 0x9b, 0x66, 0xab, 0x71, 0xd0, 0x2c, 0x54, 0x49, 0x7a, 0x96, 0x5d, 0x54, 0x90, 0xb7, 0x09,
	0xa7, 0x53, 0xaf, 0x60, 0xb3, 0xf3, 0x0a, 0xfe, 0x5b, 0x50, 0x0b, 0x7a, 0xd7, 0x38, 0x9d, 0xd1,
	0xbb, 0xc6, 0x29, 0xb9, 0x07, 0xe5, 0xb1, 0xdf, 0x1f, 0xa1, 0xa6, 0xa6, 0x2e, 0x47, 0xa5, 0x97,
	0x86, 0x73, 0x01, 0x90, 0xcf, 0x2e, 0x79, 0x94, 0x17, 0x66, 0x61, 0xf4, 0x54, 0xb9, 0x55, 0xb1,
	0xf6, 0xa0, 0xdc, 0x15, 0x79, 0xe9, 0x3d, 0xd8, 0x5a, 0xcc, 0xd7, 0x53, 0x6a, 0xe7, 0x1d, 0x54,
	0xf5, 0xc8, 0xfe, 0xa5, 0xe3, 0x6d, 0xa8, 0xf4, 0xe2, 0x50, 0x2c, 0xa5, 0xca, 0x53, 0xdf, 0x9c,
	0x9f, 0x06, 0xc0, 0x79, 0x3e, 0xbf, 0xab, 0x7a, 0xbe, 0x07, 0x56, 0xc6, 0x90, 0xcb, 0x02, 0xaf,
	0xdc, 0x1a, 0x4f, 0xea, 0x05, 0x6e, 0x28, 0x70, 0xe6, 0x7a, 0x9c, 0xd0, 0x93, 0xfd, 0x5b, 0x23,
	0x6f, 0xad, 0x21, 0x5a, 0xc0, 0x38, 0x47, 0x50, 0x51, 0xdb, 0x26, 0xf2, 0x13, 0xe3, 0xa1, 0xc7,
	0x40, 0x9e, 0xc9, 0x2e, 0xd4, 0xe7, 0x4d, 0xd3, 0xec, 0x72, 0x81, 0xf3, 0xdb, 0x80, 0xaa, 0xde,
	0xbb, 0x95, 0xec, 0xf6, 0xc1, 0x1a, 0xe7, 0xec, 0x76, 0x97, 0xb7, 0xb5, 0x7d, 0xce, 0x90, 0xab,
	0xe1, 0x90, 0x48, 0xc1, 0xf3, 0x46, 0xf1, 0x34, 0xd6, 0xf1, 0xbc, 0x51, 0x38, 0xdd, 0x4b, 0xeb,
	0xce, 0x5e, 0xee, 0x7c, 0x80, 0xfa, 0x3c, 0xc4, 0x8a, 0x01, 0x7b, 0x5c, 0x1c, 0xb0, 0x55, 0x4f,
	0x50, 0x71, 0xe6, 0xbe, 0x00, 0xe4, 0x2f, 0xc6, 0x8a, 0x65, 0x5c, 0x33, 0x06, 0xb7, 0x6b, 0x68,
	0x2e, 0xd6, 0xf0, 0x1b, 0x94, 0xe5, 0x9b, 0x92, 0x53, 0x32, 0xee, 0xa4, 0x44, 0x9e, 0x16, 0x9e,
	0xc1, 0xd2, 0xba, 0x67, 0x70, 0xfe, 0x08, 0x1e, 0x1f, 0x5e, 0x3c, 0x0b, 0x63, 0x1e, 0x8d, 0xba,
	0xed, 0x20, 0x1d, 0xb8, 0xd1, 0x34, 0x43, 0xda, 0xc7, 0x5e, 0x88, 0xd4, 0xbd, 0xf2, 0xbb, 0x34,
	0x0e, 0x5c, 0xf9, 0xe7, 0x62, 0x6e, 0xf1, 0xbf, 0xd6, 0xad, 0x48, 0xe1, 0xe1, 0x9f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xd9, 0xef, 0x03, 0x06, 0xee, 0x06, 0x00, 0x00,
}
