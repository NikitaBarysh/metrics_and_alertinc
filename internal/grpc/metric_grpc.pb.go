// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: proto/metric.proto

package grpc

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SendMetric_Update_FullMethodName = "/grpc.SendMetric/Update"
)

// SendMetricClient is the client API for SendMetric service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SendMetricClient interface {
	Update(ctx context.Context, in *UpdateMetric, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type sendMetricClient struct {
	cc grpc.ClientConnInterface
}

func NewSendMetricClient(cc grpc.ClientConnInterface) SendMetricClient {
	return &sendMetricClient{cc}
}

func (c *sendMetricClient) Update(ctx context.Context, in *UpdateMetric, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SendMetric_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendMetricServer is the server API for SendMetric service.
// All implementations must embed UnimplementedSendMetricServer
// for forward compatibility
type SendMetricServer interface {
	Update(context.Context, *UpdateMetric) (*emptypb.Empty, error)
	mustEmbedUnimplementedSendMetricServer()
}

// UnimplementedSendMetricServer must be embedded to have forward compatible implementations.
type UnimplementedSendMetricServer struct {
}

func (UnimplementedSendMetricServer) Update(context.Context, *UpdateMetric) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedSendMetricServer) mustEmbedUnimplementedSendMetricServer() {}

// UnsafeSendMetricServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SendMetricServer will
// result in compilation errors.
type UnsafeSendMetricServer interface {
	mustEmbedUnimplementedSendMetricServer()
}

func RegisterSendMetricServer(s grpc.ServiceRegistrar, srv SendMetricServer) {
	s.RegisterService(&SendMetric_ServiceDesc, srv)
}

func _SendMetric_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetric)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendMetricServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SendMetric_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendMetricServer).Update(ctx, req.(*UpdateMetric))
	}
	return interceptor(ctx, in, info, handler)
}

// SendMetric_ServiceDesc is the grpc.ServiceDesc for SendMetric service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SendMetric_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.SendMetric",
	HandlerType: (*SendMetricServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _SendMetric_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metric.proto",
}