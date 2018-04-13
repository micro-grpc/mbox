// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/pb/authorize/authorizer.proto

/*
Package mbox_authorize is a generated protocol buffer package.

It is generated from these files:
	grpc/pb/authorize/authorizer.proto

It has these top-level messages:
	AuthorizeRequest
	AuthorizeBasic
	Param
	AuthorizeResponse
*/
package mbox_authorize

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

type AuthorizeRequest struct {
	Token    string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Secret   string `protobuf:"bytes,2,opt,name=secret" json:"secret,omitempty"`
	Security bool   `protobuf:"varint,3,opt,name=security" json:"security,omitempty"`
}

func (m *AuthorizeRequest) Reset()                    { *m = AuthorizeRequest{} }
func (m *AuthorizeRequest) String() string            { return proto.CompactTextString(m) }
func (*AuthorizeRequest) ProtoMessage()               {}
func (*AuthorizeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AuthorizeRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *AuthorizeRequest) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *AuthorizeRequest) GetSecurity() bool {
	if m != nil {
		return m.Security
	}
	return false
}

type AuthorizeBasic struct {
	Token    string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Secret   string `protobuf:"bytes,2,opt,name=secret" json:"secret,omitempty"`
	Security bool   `protobuf:"varint,3,opt,name=security" json:"security,omitempty"`
}

func (m *AuthorizeBasic) Reset()                    { *m = AuthorizeBasic{} }
func (m *AuthorizeBasic) String() string            { return proto.CompactTextString(m) }
func (*AuthorizeBasic) ProtoMessage()               {}
func (*AuthorizeBasic) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AuthorizeBasic) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *AuthorizeBasic) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *AuthorizeBasic) GetSecurity() bool {
	if m != nil {
		return m.Security
	}
	return false
}

type Param struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Param) Reset()                    { *m = Param{} }
func (m *Param) String() string            { return proto.CompactTextString(m) }
func (*Param) ProtoMessage()               {}
func (*Param) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Param) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Param) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type AuthorizeResponse struct {
	UserId string   `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	User   string   `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	Params []*Param `protobuf:"bytes,3,rep,name=params" json:"params,omitempty"`
	Prefix string   `protobuf:"bytes,4,opt,name=prefix" json:"prefix,omitempty"`
}

func (m *AuthorizeResponse) Reset()                    { *m = AuthorizeResponse{} }
func (m *AuthorizeResponse) String() string            { return proto.CompactTextString(m) }
func (*AuthorizeResponse) ProtoMessage()               {}
func (*AuthorizeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AuthorizeResponse) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *AuthorizeResponse) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *AuthorizeResponse) GetParams() []*Param {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *AuthorizeResponse) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func init() {
	proto.RegisterType((*AuthorizeRequest)(nil), "mbox.authorize.AuthorizeRequest")
	proto.RegisterType((*AuthorizeBasic)(nil), "mbox.authorize.AuthorizeBasic")
	proto.RegisterType((*Param)(nil), "mbox.authorize.Param")
	proto.RegisterType((*AuthorizeResponse)(nil), "mbox.authorize.AuthorizeResponse")
}

func init() { proto.RegisterFile("grpc/pb/authorize/authorizer.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x4f, 0x4e, 0xf3, 0x30,
	0x10, 0xc5, 0x95, 0xaf, 0x6d, 0xbe, 0x32, 0x48, 0x15, 0x58, 0xfc, 0x89, 0xba, 0x40, 0x21, 0xab,
	0x6e, 0x48, 0x45, 0x39, 0x01, 0xec, 0x58, 0x51, 0x02, 0x2b, 0x84, 0x84, 0x9c, 0x74, 0x00, 0x0b,
	0x12, 0x9b, 0xb1, 0x53, 0x15, 0x2e, 0xc0, 0x8d, 0x38, 0x1f, 0xb2, 0x63, 0x42, 0xa9, 0x04, 0x62,
	0xc1, 0xee, 0xbd, 0xc9, 0xe4, 0x37, 0x9e, 0x67, 0x43, 0x72, 0x47, 0xaa, 0x18, 0xab, 0x7c, 0xcc,
	0x6b, 0x73, 0x2f, 0x49, 0xbc, 0xe0, 0xa7, 0xa2, 0x54, 0x91, 0x34, 0x92, 0x0d, 0xca, 0x5c, 0x2e,
	0xd2, 0xb6, 0x9c, 0x5c, 0xc3, 0xc6, 0xf1, 0x87, 0xc9, 0xf0, 0xa9, 0x46, 0x6d, 0xd8, 0x16, 0xf4,
	0x8c, 0x7c, 0xc0, 0x2a, 0x0a, 0xe2, 0x60, 0xb4, 0x96, 0x35, 0x86, 0xed, 0x40, 0xa8, 0xb1, 0x20,
	0x34, 0xd1, 0x3f, 0x57, 0xf6, 0x8e, 0x0d, 0xa1, 0xaf, 0xb1, 0xa8, 0x49, 0x98, 0xe7, 0xa8, 0x13,
	0x07, 0xa3, 0x7e, 0xd6, 0xfa, 0xe4, 0x0a, 0x06, 0x2d, 0xfd, 0x84, 0x6b, 0x51, 0xfc, 0x21, 0xfb,
	0x10, 0x7a, 0x53, 0x4e, 0xbc, 0x64, 0x0c, 0xba, 0x15, 0x2f, 0xd1, 0x13, 0x9d, 0xb6, 0x63, 0xe6,
	0xfc, 0xb1, 0x46, 0xcf, 0x6b, 0x4c, 0xf2, 0x1a, 0xc0, 0xe6, 0xd2, 0xb6, 0x5a, 0xc9, 0x4a, 0x23,
	0xdb, 0x85, 0xff, 0xb5, 0x46, 0xba, 0x11, 0x33, 0x8f, 0x08, 0xad, 0x3d, 0x9d, 0x59, 0xb0, 0x55,
	0x9e, 0xe1, 0x34, 0x3b, 0x80, 0x50, 0xd9, 0xa9, 0x3a, 0xea, 0xc4, 0x9d, 0xd1, 0xfa, 0x64, 0x3b,
	0xfd, 0x1a, 0x68, 0xea, 0xce, 0x94, 0xf9, 0x26, 0xbb, 0x98, 0x22, 0xbc, 0x15, 0x8b, 0xa8, 0xdb,
	0xa0, 0x1b, 0x37, 0x79, 0x0b, 0x96, 0x72, 0xbf, 0x40, 0x9a, 0x8b, 0x02, 0xd9, 0x39, 0xc0, 0x94,
	0x93, 0xc6, 0x4b, 0x97, 0x49, 0xbc, 0x4a, 0x5e, 0xbd, 0xa7, 0xe1, 0xfe, 0x0f, 0x1d, 0x7e, 0xb7,
	0x33, 0x8f, 0x6c, 0xc2, 0xdf, 0xfb, 0xf6, 0x07, 0xf7, 0xfd, 0x17, 0xc0, 0x3c, 0x74, 0xcf, 0xe8,
	0xe8, 0x3d, 0x00, 0x00, 0xff, 0xff, 0xec, 0xa0, 0x4e, 0x98, 0x6c, 0x02, 0x00, 0x00,
}
