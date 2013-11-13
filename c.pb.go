// Code generated by protoc-gen-go.
// source: c.proto
// DO NOT EDIT!

package upax_go

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type UpaxClientMsg_Tag int32

const (
	UpaxClientMsg_Intro     UpaxClientMsg_Tag = 0
	UpaxClientMsg_ItsMe     UpaxClientMsg_Tag = 1
	UpaxClientMsg_KeepAlive UpaxClientMsg_Tag = 2
	UpaxClientMsg_Query     UpaxClientMsg_Tag = 3
	UpaxClientMsg_Get       UpaxClientMsg_Tag = 4
	UpaxClientMsg_Put       UpaxClientMsg_Tag = 5
	UpaxClientMsg_Bye       UpaxClientMsg_Tag = 6
	UpaxClientMsg_Ack       UpaxClientMsg_Tag = 10
	UpaxClientMsg_Data      UpaxClientMsg_Tag = 11
	UpaxClientMsg_NotFound  UpaxClientMsg_Tag = 12
	UpaxClientMsg_Error     UpaxClientMsg_Tag = 13
)

var UpaxClientMsg_Tag_name = map[int32]string{
	0:  "Intro",
	1:  "ItsMe",
	2:  "KeepAlive",
	3:  "Query",
	4:  "Get",
	5:  "Put",
	6:  "Bye",
	10: "Ack",
	11: "Data",
	12: "NotFound",
	13: "Error",
}
var UpaxClientMsg_Tag_value = map[string]int32{
	"Intro":     0,
	"ItsMe":     1,
	"KeepAlive": 2,
	"Query":     3,
	"Get":       4,
	"Put":       5,
	"Bye":       6,
	"Ack":       10,
	"Data":      11,
	"NotFound":  12,
	"Error":     13,
}

func (x UpaxClientMsg_Tag) Enum() *UpaxClientMsg_Tag {
	p := new(UpaxClientMsg_Tag)
	*p = x
	return p
}
func (x UpaxClientMsg_Tag) String() string {
	return proto.EnumName(UpaxClientMsg_Tag_name, int32(x))
}
func (x UpaxClientMsg_Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *UpaxClientMsg_Tag) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(UpaxClientMsg_Tag_value, data, "UpaxClientMsg_Tag")
	if err != nil {
		return err
	}
	*x = UpaxClientMsg_Tag(value)
	return nil
}

type UpaxClientMsg struct {
	Op               *UpaxClientMsg_Tag   `protobuf:"varint,1,opt,enum=upax_go.UpaxClientMsg_Tag" json:"Op,omitempty"`
	MsgN             *uint64              `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	ID               []byte               `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	Salt             []byte               `protobuf:"bytes,4,opt" json:"Salt,omitempty"`
	Sig              []byte               `protobuf:"bytes,5,opt" json:"Sig,omitempty"`
	YourMsgN         *uint64              `protobuf:"varint,6,opt" json:"YourMsgN,omitempty"`
	YourID           []byte               `protobuf:"bytes,7,opt" json:"YourID,omitempty"`
	ErrCode          *uint64              `protobuf:"varint,8,opt" json:"ErrCode,omitempty"`
	ErrDesc          *string              `protobuf:"bytes,9,opt" json:"ErrDesc,omitempty"`
	Hash             []byte               `protobuf:"bytes,10,opt" json:"Hash,omitempty"`
	Payload          []byte               `protobuf:"bytes,11,opt" json:"Payload,omitempty"`
	ClientInfo       *UpaxClientMsg_Token `protobuf:"bytes,12,opt" json:"ClientInfo,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *UpaxClientMsg) Reset()         { *m = UpaxClientMsg{} }
func (m *UpaxClientMsg) String() string { return proto.CompactTextString(m) }
func (*UpaxClientMsg) ProtoMessage()    {}

func (m *UpaxClientMsg) GetOp() UpaxClientMsg_Tag {
	if m != nil && m.Op != nil {
		return *m.Op
	}
	return 0
}

func (m *UpaxClientMsg) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *UpaxClientMsg) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *UpaxClientMsg) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *UpaxClientMsg) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func (m *UpaxClientMsg) GetYourMsgN() uint64 {
	if m != nil && m.YourMsgN != nil {
		return *m.YourMsgN
	}
	return 0
}

func (m *UpaxClientMsg) GetYourID() []byte {
	if m != nil {
		return m.YourID
	}
	return nil
}

func (m *UpaxClientMsg) GetErrCode() uint64 {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return 0
}

func (m *UpaxClientMsg) GetErrDesc() string {
	if m != nil && m.ErrDesc != nil {
		return *m.ErrDesc
	}
	return ""
}

func (m *UpaxClientMsg) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *UpaxClientMsg) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *UpaxClientMsg) GetClientInfo() *UpaxClientMsg_Token {
	if m != nil {
		return m.ClientInfo
	}
	return nil
}

type UpaxClientMsg_Token struct {
	Name             *string `protobuf:"bytes,1,opt" json:"Name,omitempty"`
	ID               []byte  `protobuf:"bytes,3,opt" json:"ID,omitempty"`
	CommsKey         []byte  `protobuf:"bytes,4,opt" json:"CommsKey,omitempty"`
	SigKey           []byte  `protobuf:"bytes,5,opt" json:"SigKey,omitempty"`
	Salt             []byte  `protobuf:"bytes,6,opt" json:"Salt,omitempty"`
	DigSig           []byte  `protobuf:"bytes,7,opt" json:"DigSig,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *UpaxClientMsg_Token) Reset()         { *m = UpaxClientMsg_Token{} }
func (m *UpaxClientMsg_Token) String() string { return proto.CompactTextString(m) }
func (*UpaxClientMsg_Token) ProtoMessage()    {}

func (m *UpaxClientMsg_Token) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UpaxClientMsg_Token) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *UpaxClientMsg_Token) GetCommsKey() []byte {
	if m != nil {
		return m.CommsKey
	}
	return nil
}

func (m *UpaxClientMsg_Token) GetSigKey() []byte {
	if m != nil {
		return m.SigKey
	}
	return nil
}

func (m *UpaxClientMsg_Token) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *UpaxClientMsg_Token) GetDigSig() []byte {
	if m != nil {
		return m.DigSig
	}
	return nil
}

type UpaxClientMsg_LogEntry struct {
	Index            *int64  `protobuf:"varint,1,opt" json:"Index,omitempty"`
	Timestamp        *int64  `protobuf:"varint,2,opt" json:"Timestamp,omitempty"`
	ContentKey       []byte  `protobuf:"bytes,3,opt" json:"ContentKey,omitempty"`
	Owner            []byte  `protobuf:"bytes,4,opt" json:"Owner,omitempty"`
	Src              *string `protobuf:"bytes,5,opt" json:"Src,omitempty"`
	Path             *string `protobuf:"bytes,6,opt" json:"Path,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *UpaxClientMsg_LogEntry) Reset()         { *m = UpaxClientMsg_LogEntry{} }
func (m *UpaxClientMsg_LogEntry) String() string { return proto.CompactTextString(m) }
func (*UpaxClientMsg_LogEntry) ProtoMessage()    {}

func (m *UpaxClientMsg_LogEntry) GetIndex() int64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *UpaxClientMsg_LogEntry) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *UpaxClientMsg_LogEntry) GetContentKey() []byte {
	if m != nil {
		return m.ContentKey
	}
	return nil
}

func (m *UpaxClientMsg_LogEntry) GetOwner() []byte {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *UpaxClientMsg_LogEntry) GetSrc() string {
	if m != nil && m.Src != nil {
		return *m.Src
	}
	return ""
}

func (m *UpaxClientMsg_LogEntry) GetPath() string {
	if m != nil && m.Path != nil {
		return *m.Path
	}
	return ""
}

func init() {
	proto.RegisterEnum("upax_go.UpaxClientMsg_Tag", UpaxClientMsg_Tag_name, UpaxClientMsg_Tag_value)
}
