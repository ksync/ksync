// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/ksync.proto

package proto_ksync

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"

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
	Items                map[string]*Spec `protobuf:"bytes,1,rep,name=items" json:"items,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SpecList) Reset()         { *m = SpecList{} }
func (m *SpecList) String() string { return proto.CompactTextString(m) }
func (*SpecList) ProtoMessage()    {}
func (*SpecList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{0}
}
func (m *SpecList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpecList.Unmarshal(m, b)
}
func (m *SpecList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpecList.Marshal(b, m, deterministic)
}
func (dst *SpecList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpecList.Merge(dst, src)
}
func (m *SpecList) XXX_Size() int {
	return xxx_messageInfo_SpecList.Size(m)
}
func (m *SpecList) XXX_DiscardUnknown() {
	xxx_messageInfo_SpecList.DiscardUnknown(m)
}

var xxx_messageInfo_SpecList proto.InternalMessageInfo

func (m *SpecList) GetItems() map[string]*Spec {
	if m != nil {
		return m.Items
	}
	return nil
}

type Spec struct {
	Details              *SpecDetails `protobuf:"bytes,1,opt,name=details" json:"details,omitempty"`
	Services             *ServiceList `protobuf:"bytes,2,opt,name=services" json:"services,omitempty"`
	Status               string       `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Spec) Reset()         { *m = Spec{} }
func (m *Spec) String() string { return proto.CompactTextString(m) }
func (*Spec) ProtoMessage()    {}
func (*Spec) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{1}
}
func (m *Spec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Spec.Unmarshal(m, b)
}
func (m *Spec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Spec.Marshal(b, m, deterministic)
}
func (dst *Spec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Spec.Merge(dst, src)
}
func (m *Spec) XXX_Size() int {
	return xxx_messageInfo_Spec.Size(m)
}
func (m *Spec) XXX_DiscardUnknown() {
	xxx_messageInfo_Spec.DiscardUnknown(m)
}

var xxx_messageInfo_Spec proto.InternalMessageInfo

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
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ContainerName        string   `protobuf:"bytes,2,opt,name=container_name,json=containerName" json:"container_name,omitempty"`
	PodName              string   `protobuf:"bytes,3,opt,name=pod_name,json=podName" json:"pod_name,omitempty"`
	Selector             string   `protobuf:"bytes,4,opt,name=selector" json:"selector,omitempty"`
	Namespace            string   `protobuf:"bytes,5,opt,name=namespace" json:"namespace,omitempty"`
	LocalPath            string   `protobuf:"bytes,6,opt,name=local_path,json=localPath" json:"local_path,omitempty"`
	RemotePath           string   `protobuf:"bytes,7,opt,name=remote_path,json=remotePath" json:"remote_path,omitempty"`
	Reload               bool     `protobuf:"varint,8,opt,name=reload" json:"reload,omitempty"`
	LocalReadOnly        bool     `protobuf:"varint,9,opt,name=local_read_only,json=localReadOnly" json:"local_read_only,omitempty"`
	RemoteReadOnly       bool     `protobuf:"varint,10,opt,name=remote_read_only,json=remoteReadOnly" json:"remote_read_only,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpecDetails) Reset()         { *m = SpecDetails{} }
func (m *SpecDetails) String() string { return proto.CompactTextString(m) }
func (*SpecDetails) ProtoMessage()    {}
func (*SpecDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{2}
}
func (m *SpecDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpecDetails.Unmarshal(m, b)
}
func (m *SpecDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpecDetails.Marshal(b, m, deterministic)
}
func (dst *SpecDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpecDetails.Merge(dst, src)
}
func (m *SpecDetails) XXX_Size() int {
	return xxx_messageInfo_SpecDetails.Size(m)
}
func (m *SpecDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_SpecDetails.DiscardUnknown(m)
}

var xxx_messageInfo_SpecDetails proto.InternalMessageInfo

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

func (m *SpecDetails) GetLocalReadOnly() bool {
	if m != nil {
		return m.LocalReadOnly
	}
	return false
}

func (m *SpecDetails) GetRemoteReadOnly() bool {
	if m != nil {
		return m.RemoteReadOnly
	}
	return false
}

type ServiceList struct {
	Items                []*Service `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ServiceList) Reset()         { *m = ServiceList{} }
func (m *ServiceList) String() string { return proto.CompactTextString(m) }
func (*ServiceList) ProtoMessage()    {}
func (*ServiceList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{3}
}
func (m *ServiceList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceList.Unmarshal(m, b)
}
func (m *ServiceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceList.Marshal(b, m, deterministic)
}
func (dst *ServiceList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceList.Merge(dst, src)
}
func (m *ServiceList) XXX_Size() int {
	return xxx_messageInfo_ServiceList.Size(m)
}
func (m *ServiceList) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceList.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceList proto.InternalMessageInfo

func (m *ServiceList) GetItems() []*Service {
	if m != nil {
		return m.Items
	}
	return nil
}

type Service struct {
	SpecDetails          *SpecDetails     `protobuf:"bytes,1,opt,name=spec_details,json=specDetails" json:"spec_details,omitempty"`
	RemoteContainer      *RemoteContainer `protobuf:"bytes,2,opt,name=remote_container,json=remoteContainer" json:"remote_container,omitempty"`
	Status               string           `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}
func (*Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{4}
}
func (m *Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service.Unmarshal(m, b)
}
func (m *Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service.Marshal(b, m, deterministic)
}
func (dst *Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service.Merge(dst, src)
}
func (m *Service) XXX_Size() int {
	return xxx_messageInfo_Service.Size(m)
}
func (m *Service) XXX_DiscardUnknown() {
	xxx_messageInfo_Service.DiscardUnknown(m)
}

var xxx_messageInfo_Service proto.InternalMessageInfo

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
	Id                   string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ContainerName        string   `protobuf:"bytes,2,opt,name=container_name,json=containerName" json:"container_name,omitempty"`
	NodeName             string   `protobuf:"bytes,3,opt,name=node_name,json=nodeName" json:"node_name,omitempty"`
	PodName              string   `protobuf:"bytes,4,opt,name=pod_name,json=podName" json:"pod_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoteContainer) Reset()         { *m = RemoteContainer{} }
func (m *RemoteContainer) String() string { return proto.CompactTextString(m) }
func (*RemoteContainer) ProtoMessage()    {}
func (*RemoteContainer) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{5}
}
func (m *RemoteContainer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoteContainer.Unmarshal(m, b)
}
func (m *RemoteContainer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoteContainer.Marshal(b, m, deterministic)
}
func (dst *RemoteContainer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteContainer.Merge(dst, src)
}
func (m *RemoteContainer) XXX_Size() int {
	return xxx_messageInfo_RemoteContainer.Size(m)
}
func (m *RemoteContainer) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteContainer.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteContainer proto.InternalMessageInfo

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

type Alive struct {
	Alive                bool     `protobuf:"varint,1,opt,name=alive" json:"alive,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Alive) Reset()         { *m = Alive{} }
func (m *Alive) String() string { return proto.CompactTextString(m) }
func (*Alive) ProtoMessage()    {}
func (*Alive) Descriptor() ([]byte, []int) {
	return fileDescriptor_ksync_f806e6cd344008f8, []int{6}
}
func (m *Alive) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Alive.Unmarshal(m, b)
}
func (m *Alive) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Alive.Marshal(b, m, deterministic)
}
func (dst *Alive) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Alive.Merge(dst, src)
}
func (m *Alive) XXX_Size() int {
	return xxx_messageInfo_Alive.Size(m)
}
func (m *Alive) XXX_DiscardUnknown() {
	xxx_messageInfo_Alive.DiscardUnknown(m)
}

var xxx_messageInfo_Alive proto.InternalMessageInfo

func (m *Alive) GetAlive() bool {
	if m != nil {
		return m.Alive
	}
	return false
}

func init() {
	proto.RegisterType((*SpecList)(nil), "proto.ksync.SpecList")
	proto.RegisterMapType((map[string]*Spec)(nil), "proto.ksync.SpecList.ItemsEntry")
	proto.RegisterType((*Spec)(nil), "proto.ksync.Spec")
	proto.RegisterType((*SpecDetails)(nil), "proto.ksync.SpecDetails")
	proto.RegisterType((*ServiceList)(nil), "proto.ksync.ServiceList")
	proto.RegisterType((*Service)(nil), "proto.ksync.Service")
	proto.RegisterType((*RemoteContainer)(nil), "proto.ksync.RemoteContainer")
	proto.RegisterType((*Alive)(nil), "proto.ksync.Alive")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// KsyncClient is the client API for Ksync service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KsyncClient interface {
	GetSpecList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*SpecList, error)
	Restart(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Error, error)
	RestartSyncthing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Error, error)
	IsAlive(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Alive, error)
}

type ksyncClient struct {
	cc *grpc.ClientConn
}

func NewKsyncClient(cc *grpc.ClientConn) KsyncClient {
	return &ksyncClient{cc}
}

func (c *ksyncClient) GetSpecList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*SpecList, error) {
	out := new(SpecList)
	err := c.cc.Invoke(ctx, "/proto.ksync.Ksync/GetSpecList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ksyncClient) Restart(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/proto.ksync.Ksync/Restart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ksyncClient) RestartSyncthing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/proto.ksync.Ksync/RestartSyncthing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ksyncClient) IsAlive(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Alive, error) {
	out := new(Alive)
	err := c.cc.Invoke(ctx, "/proto.ksync.Ksync/IsAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Ksync service

type KsyncServer interface {
	GetSpecList(context.Context, *empty.Empty) (*SpecList, error)
	Restart(context.Context, *empty.Empty) (*Error, error)
	RestartSyncthing(context.Context, *empty.Empty) (*Error, error)
	IsAlive(context.Context, *empty.Empty) (*Alive, error)
}

func RegisterKsyncServer(s *grpc.Server, srv KsyncServer) {
	s.RegisterService(&_Ksync_serviceDesc, srv)
}

func _Ksync_GetSpecList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
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
		return srv.(KsyncServer).GetSpecList(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ksync_Restart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KsyncServer).Restart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ksync.Ksync/Restart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KsyncServer).Restart(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ksync_RestartSyncthing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KsyncServer).RestartSyncthing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ksync.Ksync/RestartSyncthing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KsyncServer).RestartSyncthing(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ksync_IsAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KsyncServer).IsAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ksync.Ksync/IsAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KsyncServer).IsAlive(ctx, req.(*empty.Empty))
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
		{
			MethodName: "Restart",
			Handler:    _Ksync_Restart_Handler,
		},
		{
			MethodName: "RestartSyncthing",
			Handler:    _Ksync_RestartSyncthing_Handler,
		},
		{
			MethodName: "IsAlive",
			Handler:    _Ksync_IsAlive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ksync.proto",
}

func init() { proto.RegisterFile("proto/ksync.proto", fileDescriptor_ksync_f806e6cd344008f8) }

var fileDescriptor_ksync_f806e6cd344008f8 = []byte{
	// 605 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x6e, 0xd3, 0x3e,
	0x14, 0x5f, 0xb2, 0xa6, 0x49, 0x4f, 0xfe, 0xfb, 0xb2, 0xf6, 0x9f, 0x42, 0xb7, 0x89, 0x2a, 0x12,
	0x50, 0x71, 0x91, 0x4a, 0x05, 0x31, 0x3e, 0x24, 0x04, 0x82, 0x6a, 0x9a, 0x86, 0x00, 0x65, 0x0f,
	0x50, 0x79, 0xc9, 0xa1, 0x8d, 0x96, 0xc6, 0x91, 0xed, 0x56, 0xca, 0x1d, 0x97, 0xdc, 0xf1, 0x18,
	0x3c, 0x0a, 0x0f, 0xc5, 0x0d, 0x8a, 0x9d, 0xa6, 0xed, 0xba, 0x49, 0xe5, 0xaa, 0xf6, 0xef, 0xe3,
	0x1c, 0xbb, 0xfe, 0xe5, 0xc0, 0x41, 0xce, 0x99, 0x64, 0xbd, 0x1b, 0x51, 0x64, 0x51, 0xa0, 0xd6,
	0xc4, 0x55, 0x3f, 0x81, 0x82, 0xda, 0xc7, 0x23, 0xc6, 0x46, 0x29, 0xf6, 0x14, 0x76, 0x3d, 0xfd,
	0xd6, 0xc3, 0x49, 0x2e, 0x0b, 0xad, 0x6c, 0x57, 0x66, 0x4e, 0x63, 0xca, 0x35, 0xe4, 0xff, 0x34,
	0xc0, 0xb9, 0xca, 0x31, 0xfa, 0x94, 0x08, 0x49, 0x5e, 0x80, 0x95, 0x48, 0x9c, 0x08, 0xcf, 0xe8,
	0x6c, 0x77, 0xdd, 0x7e, 0x27, 0x58, 0xaa, 0x1c, 0xcc, 0x55, 0xc1, 0x45, 0x29, 0x19, 0x64, 0x92,
	0x17, 0xa1, 0x96, 0xb7, 0x2f, 0x01, 0x16, 0x20, 0xd9, 0x87, 0xed, 0x1b, 0x2c, 0x3c, 0xa3, 0x63,
	0x74, 0x5b, 0x61, 0xb9, 0x24, 0x4f, 0xc0, 0x9a, 0xd1, 0x74, 0x8a, 0x9e, 0xd9, 0x31, 0xba, 0x6e,
	0xff, 0x60, 0xad, 0x6e, 0xa8, 0xf9, 0xd7, 0xe6, 0x4b, 0xc3, 0xff, 0x61, 0x40, 0xa3, 0xc4, 0x48,
	0x1f, 0xec, 0x18, 0x25, 0x4d, 0x52, 0xa1, 0x6a, 0xb9, 0x7d, 0x6f, 0xcd, 0xf7, 0x51, 0xf3, 0xe1,
	0x5c, 0x48, 0x9e, 0x83, 0x23, 0x90, 0xcf, 0x92, 0x08, 0x45, 0xd5, 0xec, 0x96, 0x49, 0x93, 0xe5,
	0x3d, 0xc2, 0x5a, 0x49, 0x8e, 0xa0, 0x29, 0x24, 0x95, 0x53, 0xe1, 0x6d, 0xab, 0x43, 0x57, 0x3b,
	0xff, 0xb7, 0x09, 0xee, 0x52, 0x1b, 0x42, 0xa0, 0x91, 0xd1, 0x09, 0x56, 0x57, 0x53, 0x6b, 0xf2,
	0x08, 0x76, 0x23, 0x96, 0x49, 0x9a, 0x64, 0xc8, 0x87, 0x8a, 0x35, 0x15, 0xbb, 0x53, 0xa3, 0x9f,
	0x4b, 0xd9, 0x03, 0x70, 0x72, 0x16, 0x6b, 0x81, 0x6e, 0x62, 0xe7, 0x2c, 0x56, 0x54, 0xbb, 0x3c,
	0x73, 0x8a, 0x91, 0x64, 0xdc, 0x6b, 0x28, 0xaa, 0xde, 0x93, 0x13, 0x68, 0x95, 0x16, 0x91, 0xd3,
	0x08, 0x3d, 0x4b, 0x91, 0x0b, 0x80, 0x9c, 0x02, 0xa4, 0x2c, 0xa2, 0xe9, 0x30, 0xa7, 0x72, 0xec,
	0x35, 0x35, 0xad, 0x90, 0xaf, 0x54, 0x8e, 0xc9, 0x43, 0x70, 0x39, 0x4e, 0x98, 0x44, 0xcd, 0xdb,
	0x8a, 0x07, 0x0d, 0x29, 0xc1, 0x11, 0x34, 0x39, 0xa6, 0x8c, 0xc6, 0x9e, 0xd3, 0x31, 0xba, 0x4e,
	0x58, 0xed, 0xc8, 0x63, 0xd8, 0xd3, 0x75, 0x39, 0xd2, 0x78, 0xc8, 0xb2, 0xb4, 0xf0, 0x5a, 0x4a,
	0xb0, 0xa3, 0xe0, 0x10, 0x69, 0xfc, 0x25, 0x4b, 0x0b, 0xd2, 0x85, 0xfd, 0xaa, 0xc1, 0x42, 0x08,
	0x4a, 0xb8, 0xab, 0xf1, 0xb9, 0xd2, 0x7f, 0x05, 0xee, 0xd2, 0x5f, 0x4f, 0x9e, 0xae, 0x06, 0xed,
	0xf0, 0xae, 0x37, 0xaa, 0xc2, 0xe5, 0xff, 0x32, 0xc0, 0xae, 0x20, 0xf2, 0x06, 0xfe, 0x13, 0x39,
	0x46, 0xc3, 0x4d, 0x73, 0xe1, 0x8a, 0xa5, 0xd7, 0x3b, 0xaf, 0x4f, 0x5b, 0x3f, 0x4d, 0x95, 0x91,
	0x93, 0x95, 0x02, 0xa1, 0x12, 0x7d, 0x98, 0x6b, 0xc2, 0x3d, 0xbe, 0x0a, 0xdc, 0x1b, 0x97, 0xef,
	0x06, 0xec, 0xdd, 0x32, 0x93, 0x5d, 0x30, 0x93, 0xb8, 0x0a, 0x8c, 0x99, 0xc4, 0x9b, 0xc6, 0xe5,
	0x18, 0x5a, 0x19, 0x8b, 0x71, 0x39, 0x2f, 0x4e, 0x09, 0xac, 0x65, 0xa9, 0xb1, 0x92, 0x25, 0xff,
	0x14, 0xac, 0xf7, 0x69, 0x32, 0x43, 0x72, 0x08, 0x16, 0x2d, 0x17, 0xaa, 0xb5, 0x13, 0xea, 0x4d,
	0xff, 0x8f, 0x01, 0xd6, 0x65, 0x79, 0x49, 0xf2, 0x16, 0xdc, 0x73, 0x94, 0xf5, 0x97, 0x7f, 0x14,
	0xe8, 0xb9, 0x11, 0xcc, 0xe7, 0x46, 0x30, 0x28, 0xe7, 0x46, 0xfb, 0xff, 0x3b, 0x47, 0x80, 0xbf,
	0x45, 0xce, 0xc0, 0x0e, 0x51, 0x48, 0xca, 0xef, 0xf7, 0x92, 0x15, 0xef, 0x80, 0x73, 0xc6, 0xfd,
	0x2d, 0xf2, 0x0e, 0xf6, 0x2b, 0xe3, 0x55, 0x91, 0x45, 0x72, 0x9c, 0x64, 0xa3, 0x7f, 0xac, 0x70,
	0x06, 0xf6, 0x85, 0xd0, 0xb7, 0xdc, 0xcc, 0xa8, 0xb4, 0xfe, 0xd6, 0x75, 0x53, 0x81, 0xcf, 0xfe,
	0x06, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x39, 0xb4, 0xdd, 0x44, 0x05, 0x00, 0x00,
}
