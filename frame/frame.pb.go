// Code generated by protoc-gen-go. DO NOT EDIT.
// source: frame.proto

/*
Package frame is a generated protocol buffer package.

It is generated from these files:
	frame.proto

It has these top-level messages:
	Request
	Person
*/
package frame

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

type Request struct {
	Id               *int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Message          *string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Request) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

type Person struct {
	Id               *int32  `protobuf:"varint,1,req,name=id,def=0" json:"id,omitempty"`
	Name             *string `protobuf:"bytes,2,req,name=name,def=doesn't exist" json:"name,omitempty"`
	Age              *int32  `protobuf:"varint,3,req,name=age,def=0" json:"age,omitempty"`
	City             *string `protobuf:"bytes,4,req,name=city,def=doesn't exist" json:"city,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Person) Reset()                    { *m = Person{} }
func (m *Person) String() string            { return proto.CompactTextString(m) }
func (*Person) ProtoMessage()               {}
func (*Person) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

const Default_Person_Id int32 = 0
const Default_Person_Name string = "doesn't exist"
const Default_Person_Age int32 = 0
const Default_Person_City string = "doesn't exist"

func (m *Person) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return Default_Person_Id
}

func (m *Person) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return Default_Person_Name
}

func (m *Person) GetAge() int32 {
	if m != nil && m.Age != nil {
		return *m.Age
	}
	return Default_Person_Age
}

func (m *Person) GetCity() string {
	if m != nil && m.City != nil {
		return *m.City
	}
	return Default_Person_City
}

func init() {
	proto.RegisterType((*Request)(nil), "frame.Request")
	proto.RegisterType((*Person)(nil), "frame.Person")
}

func init() { proto.RegisterFile("frame.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 160 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0xcc, 0xb1, 0x0a, 0xc2, 0x30,
	0x14, 0x85, 0x61, 0x92, 0xb6, 0x96, 0x5e, 0x51, 0x30, 0x2e, 0x19, 0x6b, 0x17, 0x33, 0x89, 0xe0,
	0xd6, 0xa7, 0x90, 0xbc, 0x41, 0x31, 0x57, 0xc9, 0x90, 0x44, 0x73, 0xa3, 0xe8, 0xdb, 0x4b, 0x82,
	0xdd, 0x1c, 0x7f, 0x0e, 0xdf, 0x81, 0xe5, 0x35, 0x4e, 0x0e, 0x0f, 0xf7, 0x18, 0x52, 0x10, 0x4d,
	0x89, 0xe1, 0x04, 0xad, 0xc6, 0xc7, 0x13, 0x29, 0x89, 0x35, 0x70, 0x6b, 0x24, 0xeb, 0x99, 0x6a,
	0x34, 0xb7, 0x46, 0x48, 0x68, 0x1d, 0x12, 0x4d, 0x37, 0x94, 0xbc, 0x67, 0xaa, 0xd3, 0x73, 0x0e,
	0x2f, 0x58, 0x9c, 0x31, 0x52, 0xf0, 0x62, 0xf3, 0x33, 0x5c, 0x35, 0x23, 0x3b, 0x16, 0xb6, 0x83,
	0xda, 0x4f, 0x2e, 0x1b, 0xae, 0xba, 0x71, 0x65, 0x02, 0x92, 0xdf, 0xa7, 0x1e, 0xdf, 0x96, 0x92,
	0x2e, 0x93, 0xd8, 0x42, 0x95, 0x5f, 0xab, 0x99, 0xe5, 0xca, 0xee, 0x62, 0xd3, 0x47, 0xd6, 0x7f,
	0x5d, 0x9e, 0xbe, 0x01, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x2b, 0x64, 0xed, 0xc1, 0x00, 0x00, 0x00,
}