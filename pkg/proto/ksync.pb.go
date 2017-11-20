// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/ksync.proto

/*
Package proto_ksync is a generated protocol buffer package.

It is generated from these files:
	proto/ksync.proto
	proto/radar.proto

It has these top-level messages:
	SpecList
	Spec
	SpecDetails
	ServiceList
	Service
	RemoteContainer
	ContainerPath
	BasePath
	Error
	VersionInfo
*/
package proto_ksync

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SpecList struct {
	Items map[string]*Spec `protobuf:"bytes,1,rep,name=items" json:"items,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *SpecList) Reset()                    { *m = SpecList{} }
func (m *SpecList) String() string            { return proto.CompactTextString(m) }
func (*SpecList) ProtoMessage()               {}
func (*SpecList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SpecList) GetItems() map[string]*Spec {
	if m != nil {
		return m.Items
	}
	return nil
}

type Spec struct {
	Details  *SpecDetails `protobuf:"bytes,1,opt,name=details" json:"details,omitempty"`
	Services *ServiceList `protobuf:"bytes,2,opt,name=services" json:"services,omitempty"`
	Status   string       `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
}

func (m *Spec) Reset()                    { *m = Spec{} }
func (m *Spec) String() string            { return proto.CompactTextString(m) }
func (*Spec) ProtoMessage()               {}
func (*Spec) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Spec) GetDetails() *SpecDetails {
	if m != nil {
		return m.Details
	}
	return nil
}

func (m *Spec) GetServices() *ServiceList {
	if m != nil {
		return m.Services
	}
	return nil
}

func (m *Spec) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type SpecDetails struct {
	Name          string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ContainerName string `protobuf:"bytes,2,opt,name=container_name,json=containerName" json:"container_name,omitempty"`
	PodName       string `protobuf:"bytes,3,opt,name=pod_name,json=podName" json:"pod_name,omitempty"`
	Selector      string `protobuf:"bytes,4,opt,name=selector" json:"selector,omitempty"`
	Namespace     string `protobuf:"bytes,5,opt,name=namespace" json:"namespace,omitempty"`
	LocalPath     string `protobuf:"bytes,6,opt,name=local_path,json=localPath" json:"local_path,omitempty"`
	RemotePath    string `protobuf:"bytes,7,opt,name=remote_path,json=remotePath" json:"remote_path,omitempty"`
	Reload        bool   `protobuf:"varint,8,opt,name=reload" json:"reload,omitempty"`
}

func (m *SpecDetails) Reset()                    { *m = SpecDetails{} }
func (m *SpecDetails) String() string            { return proto.CompactTextString(m) }
func (*SpecDetails) ProtoMessage()               {}
func (*SpecDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SpecDetails) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SpecDetails) GetContainerName() string {
	if m != nil {
		return m.ContainerName
	}
	return ""
}

func (m *SpecDetails) GetPodName() string {
	if m != nil {
		return m.PodName
	}
	return ""
}

func (m *SpecDetails) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

func (m *SpecDetails) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *SpecDetails) GetLocalPath() string {
	if m != nil {
		return m.LocalPath
	}
	return ""
}

func (m *SpecDetails) GetRemotePath() string {
	if m != nil {
		return m.RemotePath
	}
	return ""
}

func (m *SpecDetails) GetReload() bool {
	if m != nil {
		return m.Reload
	}
	return false
}

type ServiceList struct {
	Items []*Service `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
}

func (m *ServiceList) Reset()                    { *m = ServiceList{} }
func (m *ServiceList) String() string            { return proto.CompactTextString(m) }
func (*ServiceList) ProtoMessage()               {}
func (*ServiceList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ServiceList) GetItems() []*Service {
	if m != nil {
		return m.Items
	}
	return nil
}

type Service struct {
	SpecDetails     *SpecDetails     `protobuf:"bytes,1,opt,name=spec_details,json=specDetails" json:"spec_details,omitempty"`
	RemoteContainer *RemoteContainer `protobuf:"bytes,2,opt,name=remote_container,json=remoteContainer" json:"remote_container,omitempty"`
	Status          string           `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
}

func (m *Service) Reset()                    { *m = Service{} }
func (m *Service) String() string            { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()               {}
func (*Service) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Service) GetSpecDetails() *SpecDetails {
	if m != nil {
		return m.SpecDetails
	}
	return nil
}

func (m *Service) GetRemoteContainer() *RemoteContainer {
	if m != nil {
		return m.RemoteContainer
	}
	return nil
}

func (m *Service) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type RemoteContainer struct {
	Id            string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ContainerName string `protobuf:"bytes,2,opt,name=container_name,json=containerName" json:"container_name,omitempty"`
	NodeName      string `protobuf:"bytes,3,opt,name=node_name,json=nodeName" json:"node_name,omitempty"`
	PodName       string `protobuf:"bytes,4,opt,name=pod_name,json=podName" json:"pod_name,omitempty"`
}

func (m *RemoteContainer) Reset()                    { *m = RemoteContainer{} }
func (m *RemoteContainer) String() string            { return proto.CompactTextString(m) }
func (*RemoteContainer) ProtoMessage()               {}
func (*RemoteContainer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RemoteContainer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RemoteContainer) GetContainerName() string {
	if m != nil {
		return m.ContainerName
	}
	return ""
}

func (m *RemoteContainer) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *RemoteContainer) GetPodName() string {
	if m != nil {
		return m.PodName
	}
	return ""
}

func init() {
	proto.RegisterType((*SpecList)(nil), "proto.ksync.SpecList")
	proto.RegisterType((*Spec)(nil), "proto.ksync.Spec")
	proto.RegisterType((*SpecDetails)(nil), "proto.ksync.SpecDetails")
	proto.RegisterType((*ServiceList)(nil), "proto.ksync.ServiceList")
	proto.RegisterType((*Service)(nil), "proto.ksync.Service")
	proto.RegisterType((*RemoteContainer)(nil), "proto.ksync.RemoteContainer")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Ksync service

type KsyncClient interface {
	GetSpecList(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*SpecList, error)
}

type ksyncClient struct {
	cc *grpc.ClientConn
}

func NewKsyncClient(cc *grpc.ClientConn) KsyncClient {
	return &ksyncClient{cc}
}

func (c *ksyncClient) GetSpecList(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*SpecList, error) {
	out := new(SpecList)
	err := grpc.Invoke(ctx, "/proto.ksync.Ksync/GetSpecList", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Ksync service

type KsyncServer interface {
	GetSpecList(context.Context, *google_protobuf.Empty) (*SpecList, error)
}

func RegisterKsyncServer(s *grpc.Server, srv KsyncServer) {
	s.RegisterService(&_Ksync_serviceDesc, srv)
}

func _Ksync_GetSpecList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KsyncServer).GetSpecList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ksync.Ksync/GetSpecList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KsyncServer).GetSpecList(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ksync_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ksync.Ksync",
	HandlerType: (*KsyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSpecList",
			Handler:    _Ksync_GetSpecList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ksync.proto",
}

func init() { proto.RegisterFile("proto/ksync.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 495 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4b, 0x8b, 0xd4, 0x40,
	0x10, 0xb6, 0xe7, 0x99, 0xa9, 0xe8, 0x3e, 0x1a, 0x1d, 0xda, 0xd9, 0x15, 0x43, 0x40, 0x1c, 0x3c,
	0x64, 0x20, 0x8a, 0xf8, 0x00, 0x2f, 0xba, 0x0c, 0xb2, 0x22, 0x12, 0x7f, 0xc0, 0xd0, 0x9b, 0x94,
	0xbb, 0x61, 0x33, 0xe9, 0x90, 0xee, 0x59, 0x98, 0x9b, 0x47, 0x6f, 0xfe, 0x0c, 0x7f, 0xa1, 0x77,
	0xe9, 0xc7, 0x64, 0x33, 0x3b, 0x2b, 0xec, 0x29, 0xdd, 0xdf, 0xf7, 0x55, 0x75, 0x55, 0x7d, 0x29,
	0x38, 0xac, 0x6a, 0xa1, 0xc4, 0xec, 0x52, 0xae, 0xcb, 0x34, 0x32, 0x67, 0xea, 0x9b, 0x4f, 0x64,
	0xa0, 0xc9, 0xd1, 0xb9, 0x10, 0xe7, 0x05, 0xce, 0x0c, 0x76, 0xb6, 0xfa, 0x31, 0xc3, 0x65, 0xa5,
	0xd6, 0x56, 0x19, 0xfe, 0x26, 0xe0, 0x7d, 0xaf, 0x30, 0xfd, 0x92, 0x4b, 0x45, 0x5f, 0x43, 0x3f,
	0x57, 0xb8, 0x94, 0x8c, 0x04, 0xdd, 0xa9, 0x1f, 0x07, 0x51, 0x2b, 0x4d, 0xb4, 0x51, 0x45, 0x9f,
	0xb5, 0xe4, 0xa4, 0x54, 0xf5, 0x3a, 0xb1, 0xf2, 0xc9, 0x29, 0xc0, 0x35, 0x48, 0x0f, 0xa0, 0x7b,
	0x89, 0x6b, 0x46, 0x02, 0x32, 0x1d, 0x25, 0xfa, 0x48, 0x9f, 0x43, 0xff, 0x8a, 0x17, 0x2b, 0x64,
	0x9d, 0x80, 0x4c, 0xfd, 0xf8, 0x70, 0x27, 0x6f, 0x62, 0xf9, 0x77, 0x9d, 0x37, 0x24, 0xfc, 0x45,
	0xa0, 0xa7, 0x31, 0x1a, 0xc3, 0x30, 0x43, 0xc5, 0xf3, 0x42, 0x9a, 0x5c, 0x7e, 0xcc, 0x76, 0xe2,
	0x3e, 0x59, 0x3e, 0xd9, 0x08, 0xe9, 0x2b, 0xf0, 0x24, 0xd6, 0x57, 0x79, 0x8a, 0xd2, 0x3d, 0x76,
	0x23, 0xc8, 0x92, 0xba, 0x8f, 0xa4, 0x51, 0xd2, 0x31, 0x0c, 0xa4, 0xe2, 0x6a, 0x25, 0x59, 0xd7,
	0x14, 0xed, 0x6e, 0xe1, 0x5f, 0x02, 0x7e, 0xeb, 0x19, 0x4a, 0xa1, 0x57, 0xf2, 0x25, 0xba, 0xd6,
	0xcc, 0x99, 0x3e, 0x83, 0xbd, 0x54, 0x94, 0x8a, 0xe7, 0x25, 0xd6, 0x0b, 0xc3, 0x76, 0x0c, 0xfb,
	0xa0, 0x41, 0xbf, 0x6a, 0xd9, 0x63, 0xf0, 0x2a, 0x91, 0x59, 0x81, 0x7d, 0x64, 0x58, 0x89, 0xcc,
	0x50, 0x13, 0x5d, 0x73, 0x81, 0xa9, 0x12, 0x35, 0xeb, 0x19, 0xaa, 0xb9, 0xd3, 0x63, 0x18, 0xe9,
	0x10, 0x59, 0xf1, 0x14, 0x59, 0xdf, 0x90, 0xd7, 0x00, 0x7d, 0x02, 0x50, 0x88, 0x94, 0x17, 0x8b,
	0x8a, 0xab, 0x0b, 0x36, 0xb0, 0xb4, 0x41, 0xbe, 0x71, 0x75, 0x41, 0x9f, 0x82, 0x5f, 0xe3, 0x52,
	0x28, 0xb4, 0xfc, 0xd0, 0xf0, 0x60, 0x21, 0x23, 0x18, 0xc3, 0xa0, 0xc6, 0x42, 0xf0, 0x8c, 0x79,
	0x01, 0x99, 0x7a, 0x89, 0xbb, 0x85, 0x6f, 0xc1, 0x6f, 0x0d, 0x8a, 0xbe, 0xd8, 0xfe, 0x2d, 0x1e,
	0xde, 0x36, 0x51, 0xf7, 0x2b, 0x84, 0x7f, 0x08, 0x0c, 0x1d, 0x44, 0xdf, 0xc3, 0x7d, 0x59, 0x61,
	0xba, 0xb8, 0xab, 0x8b, 0xbe, 0x6c, 0xcd, 0x7a, 0x0e, 0x07, 0xae, 0xf8, 0x66, 0x90, 0xce, 0xd1,
	0xe3, 0xad, 0x04, 0x89, 0x11, 0x7d, 0xdc, 0x68, 0x92, 0xfd, 0x7a, 0x1b, 0xf8, 0xaf, 0xb9, 0x3f,
	0x09, 0xec, 0xdf, 0x08, 0xa6, 0x7b, 0xd0, 0xc9, 0x33, 0x67, 0x6f, 0x27, 0xcf, 0xee, 0x6a, 0xee,
	0x11, 0x8c, 0x4a, 0x91, 0x61, 0xdb, 0x5d, 0x4f, 0x03, 0x3b, 0xce, 0xf7, 0xb6, 0x9c, 0x8f, 0xe7,
	0xd0, 0x3f, 0xd5, 0x4d, 0xd0, 0x0f, 0xe0, 0xcf, 0x51, 0x35, 0x7b, 0x38, 0x8e, 0xec, 0xca, 0x46,
	0x9b, 0x95, 0x8d, 0x4e, 0xf4, 0xca, 0x4e, 0x1e, 0xdd, 0xba, 0x90, 0xe1, 0xbd, 0xb3, 0x81, 0xc1,
	0x5f, 0xfe, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xb2, 0x00, 0xd9, 0x4e, 0x0b, 0x04, 0x00, 0x00,
}
