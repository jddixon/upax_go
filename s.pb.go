// Code generated by protoc-gen-go.
// source: s.proto
// DO NOT EDIT!

package upax_go

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type UpaxClusterMsg_Tag int32

const (
	UpaxClusterMsg_ItsMe     UpaxClusterMsg_Tag = 0
	UpaxClusterMsg_KeepAlive UpaxClusterMsg_Tag = 1
	UpaxClusterMsg_Get       UpaxClusterMsg_Tag = 2
	UpaxClusterMsg_IHave     UpaxClusterMsg_Tag = 3
	UpaxClusterMsg_Put       UpaxClusterMsg_Tag = 4
	UpaxClusterMsg_Bye       UpaxClusterMsg_Tag = 5
	UpaxClusterMsg_Ack       UpaxClusterMsg_Tag = 10
	UpaxClusterMsg_Data      UpaxClusterMsg_Tag = 11
	UpaxClusterMsg_Error     UpaxClusterMsg_Tag = 12
)

var UpaxClusterMsg_Tag_name = map[int32]string{
	0:  "ItsMe",
	1:  "KeepAlive",
	2:  "Get",
	3:  "IHave",
	4:  "Put",
	5:  "Bye",
	10: "Ack",
	11: "Data",
	12: "Error",
}
var UpaxClusterMsg_Tag_value = map[string]int32{
	"ItsMe":     0,
	"KeepAlive": 1,
	"Get":       2,
	"IHave":     3,
	"Put":       4,
	"Bye":       5,
	"Ack":       10,
	"Data":      11,
	"Error":     12,
}

func (x UpaxClusterMsg_Tag) Enum() *UpaxClusterMsg_Tag {
	p := new(UpaxClusterMsg_Tag)
	*p = x
	return p
}
func (x UpaxClusterMsg_Tag) String() string {
	return proto.EnumName(UpaxClusterMsg_Tag_name, int32(x))
}
func (x UpaxClusterMsg_Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *UpaxClusterMsg_Tag) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(UpaxClusterMsg_Tag_value, data, "UpaxClusterMsg_Tag")
	if err != nil {
		return err
	}
	*x = UpaxClusterMsg_Tag(value)
	return nil
}

type UpaxClusterMsg struct {
	Op               *UpaxClusterMsg_Tag `protobuf:"varint,1,opt,enum=upax_go.UpaxClusterMsg_Tag" json:"Op,omitempty"`
	MsgN             *uint64             `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte              `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	Salt             []byte              `protobuf:"bytes,4,opt" json:"Salt,omitempty"`
	Sig              []byte              `protobuf:"bytes,5,opt" json:"Sig,omitempty"`
	YourMsgN         *uint64             `protobuf:"varint,6,opt" json:"YourMsgN,omitempty"`
	YourID           []byte              `protobuf:"bytes,7,opt" json:"YourID,omitempty"`
	ErrCode          *uint64             `protobuf:"varint,8,opt" json:"ErrCode,omitempty"`
	ErrDesc          *string             `protobuf:"bytes,9,opt" json:"ErrDesc,omitempty"`
	Hash             []byte              `protobuf:"bytes,10,opt" json:"Hash,omitempty"`
	Payload          []byte              `protobuf:"bytes,11,opt" json:"Payload,omitempty"`
	Index            *int64              `protobuf:"varint,32,opt" json:"Index,omitempty"`
	Timestamp        *int64              `protobuf:"varint,33,opt" json:"Timestamp,omitempty"`
	ContentKey       []byte              `protobuf:"bytes,34,opt" json:"ContentKey,omitempty"`
	Owner            []byte              `protobuf:"bytes,35,opt" json:"Owner,omitempty"`
	Src              *string             `protobuf:"bytes,36,opt" json:"Src,omitempty"`
	Path             *string             `protobuf:"bytes,37,opt" json:"Path,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *UpaxClusterMsg) Reset()         { *m = UpaxClusterMsg{} }
func (m *UpaxClusterMsg) String() string { return proto.CompactTextString(m) }
func (*UpaxClusterMsg) ProtoMessage()    {}

func (m *UpaxClusterMsg) GetOp() UpaxClusterMsg_Tag {
	if m != nil && m.Op != nil {
		return *m.Op
	}
	return 0
}

func (m *UpaxClusterMsg) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *UpaxClusterMsg) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *UpaxClusterMsg) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *UpaxClusterMsg) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func (m *UpaxClusterMsg) GetYourMsgN() uint64 {
	if m != nil && m.YourMsgN != nil {
		return *m.YourMsgN
	}
	return 0
}

func (m *UpaxClusterMsg) GetYourID() []byte {
	if m != nil {
		return m.YourID
	}
	return nil
}

func (m *UpaxClusterMsg) GetErrCode() uint64 {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return 0
}

func (m *UpaxClusterMsg) GetErrDesc() string {
	if m != nil && m.ErrDesc != nil {
		return *m.ErrDesc
	}
	return ""
}

func (m *UpaxClusterMsg) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *UpaxClusterMsg) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *UpaxClusterMsg) GetIndex() int64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *UpaxClusterMsg) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *UpaxClusterMsg) GetContentKey() []byte {
	if m != nil {
		return m.ContentKey
	}
	return nil
}

func (m *UpaxClusterMsg) GetOwner() []byte {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *UpaxClusterMsg) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

func (m *UpaxClusterMsg) GetPath() string {
	if m != nil && m.Path != nil {
		return *m.Path
	}
	return ""
}

func init() {
	proto.RegisterEnum("upax_go.UpaxClusterMsg_Tag", UpaxClusterMsg_Tag_name, UpaxClusterMsg_Tag_value)
}
