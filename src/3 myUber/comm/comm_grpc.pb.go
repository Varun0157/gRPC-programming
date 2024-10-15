// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: comm/comm.proto

package comm

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RiderService_RequestRide_FullMethodName = "/comm.RiderService/RequestRide"
	RiderService_GetStatus_FullMethodName   = "/comm.RiderService/GetStatus"
)

// RiderServiceClient is the client API for RiderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Rider APIs
type RiderServiceClient interface {
	RequestRide(ctx context.Context, in *RideRequest, opts ...grpc.CallOption) (*RideResponse, error)
	GetStatus(ctx context.Context, in *RideStatusRequest, opts ...grpc.CallOption) (*RideStatusResponse, error)
}

type riderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRiderServiceClient(cc grpc.ClientConnInterface) RiderServiceClient {
	return &riderServiceClient{cc}
}

func (c *riderServiceClient) RequestRide(ctx context.Context, in *RideRequest, opts ...grpc.CallOption) (*RideResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RideResponse)
	err := c.cc.Invoke(ctx, RiderService_RequestRide_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *riderServiceClient) GetStatus(ctx context.Context, in *RideStatusRequest, opts ...grpc.CallOption) (*RideStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RideStatusResponse)
	err := c.cc.Invoke(ctx, RiderService_GetStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RiderServiceServer is the server API for RiderService service.
// All implementations must embed UnimplementedRiderServiceServer
// for forward compatibility.
//
// Rider APIs
type RiderServiceServer interface {
	RequestRide(context.Context, *RideRequest) (*RideResponse, error)
	GetStatus(context.Context, *RideStatusRequest) (*RideStatusResponse, error)
	mustEmbedUnimplementedRiderServiceServer()
}

// UnimplementedRiderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRiderServiceServer struct{}

func (UnimplementedRiderServiceServer) RequestRide(context.Context, *RideRequest) (*RideResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestRide not implemented")
}
func (UnimplementedRiderServiceServer) GetStatus(context.Context, *RideStatusRequest) (*RideStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedRiderServiceServer) mustEmbedUnimplementedRiderServiceServer() {}
func (UnimplementedRiderServiceServer) testEmbeddedByValue()                      {}

// UnsafeRiderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RiderServiceServer will
// result in compilation errors.
type UnsafeRiderServiceServer interface {
	mustEmbedUnimplementedRiderServiceServer()
}

func RegisterRiderServiceServer(s grpc.ServiceRegistrar, srv RiderServiceServer) {
	// If the following call pancis, it indicates UnimplementedRiderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RiderService_ServiceDesc, srv)
}

func _RiderService_RequestRide_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RideRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RiderServiceServer).RequestRide(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RiderService_RequestRide_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RiderServiceServer).RequestRide(ctx, req.(*RideRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RiderService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RideStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RiderServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RiderService_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RiderServiceServer).GetStatus(ctx, req.(*RideStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RiderService_ServiceDesc is the grpc.ServiceDesc for RiderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RiderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "comm.RiderService",
	HandlerType: (*RiderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestRide",
			Handler:    _RiderService_RequestRide_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _RiderService_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comm/comm.proto",
}

const (
	DriverService_AssignDriver_FullMethodName        = "/comm.DriverService/AssignDriver"
	DriverService_CompleteRideRequest_FullMethodName = "/comm.DriverService/CompleteRideRequest"
	DriverService_AcceptRideRequest_FullMethodName   = "/comm.DriverService/AcceptRideRequest"
	DriverService_RejectRideRequest_FullMethodName   = "/comm.DriverService/RejectRideRequest"
)

// DriverServiceClient is the client API for DriverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Driver APIs
type DriverServiceClient interface {
	AssignDriver(ctx context.Context, in *DriverAssignmentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DriverAssignmentResponse], error)
	CompleteRideRequest(ctx context.Context, in *DriverCompleteRequest, opts ...grpc.CallOption) (*DriverCompleteResponse, error)
	AcceptRideRequest(ctx context.Context, in *DriverAcceptRequest, opts ...grpc.CallOption) (*DriverAcceptResponse, error)
	RejectRideRequest(ctx context.Context, in *DriverRejectRequest, opts ...grpc.CallOption) (*DriverRejectResponse, error)
}

type driverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverServiceClient(cc grpc.ClientConnInterface) DriverServiceClient {
	return &driverServiceClient{cc}
}

func (c *driverServiceClient) AssignDriver(ctx context.Context, in *DriverAssignmentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DriverAssignmentResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DriverService_ServiceDesc.Streams[0], DriverService_AssignDriver_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DriverAssignmentRequest, DriverAssignmentResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DriverService_AssignDriverClient = grpc.ServerStreamingClient[DriverAssignmentResponse]

func (c *driverServiceClient) CompleteRideRequest(ctx context.Context, in *DriverCompleteRequest, opts ...grpc.CallOption) (*DriverCompleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DriverCompleteResponse)
	err := c.cc.Invoke(ctx, DriverService_CompleteRideRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverServiceClient) AcceptRideRequest(ctx context.Context, in *DriverAcceptRequest, opts ...grpc.CallOption) (*DriverAcceptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DriverAcceptResponse)
	err := c.cc.Invoke(ctx, DriverService_AcceptRideRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverServiceClient) RejectRideRequest(ctx context.Context, in *DriverRejectRequest, opts ...grpc.CallOption) (*DriverRejectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DriverRejectResponse)
	err := c.cc.Invoke(ctx, DriverService_RejectRideRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DriverServiceServer is the server API for DriverService service.
// All implementations must embed UnimplementedDriverServiceServer
// for forward compatibility.
//
// Driver APIs
type DriverServiceServer interface {
	AssignDriver(*DriverAssignmentRequest, grpc.ServerStreamingServer[DriverAssignmentResponse]) error
	CompleteRideRequest(context.Context, *DriverCompleteRequest) (*DriverCompleteResponse, error)
	AcceptRideRequest(context.Context, *DriverAcceptRequest) (*DriverAcceptResponse, error)
	RejectRideRequest(context.Context, *DriverRejectRequest) (*DriverRejectResponse, error)
	mustEmbedUnimplementedDriverServiceServer()
}

// UnimplementedDriverServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDriverServiceServer struct{}

func (UnimplementedDriverServiceServer) AssignDriver(*DriverAssignmentRequest, grpc.ServerStreamingServer[DriverAssignmentResponse]) error {
	return status.Errorf(codes.Unimplemented, "method AssignDriver not implemented")
}
func (UnimplementedDriverServiceServer) CompleteRideRequest(context.Context, *DriverCompleteRequest) (*DriverCompleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteRideRequest not implemented")
}
func (UnimplementedDriverServiceServer) AcceptRideRequest(context.Context, *DriverAcceptRequest) (*DriverAcceptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptRideRequest not implemented")
}
func (UnimplementedDriverServiceServer) RejectRideRequest(context.Context, *DriverRejectRequest) (*DriverRejectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectRideRequest not implemented")
}
func (UnimplementedDriverServiceServer) mustEmbedUnimplementedDriverServiceServer() {}
func (UnimplementedDriverServiceServer) testEmbeddedByValue()                       {}

// UnsafeDriverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DriverServiceServer will
// result in compilation errors.
type UnsafeDriverServiceServer interface {
	mustEmbedUnimplementedDriverServiceServer()
}

func RegisterDriverServiceServer(s grpc.ServiceRegistrar, srv DriverServiceServer) {
	// If the following call pancis, it indicates UnimplementedDriverServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DriverService_ServiceDesc, srv)
}

func _DriverService_AssignDriver_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DriverAssignmentRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DriverServiceServer).AssignDriver(m, &grpc.GenericServerStream[DriverAssignmentRequest, DriverAssignmentResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DriverService_AssignDriverServer = grpc.ServerStreamingServer[DriverAssignmentResponse]

func _DriverService_CompleteRideRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DriverCompleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).CompleteRideRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DriverService_CompleteRideRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).CompleteRideRequest(ctx, req.(*DriverCompleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverService_AcceptRideRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DriverAcceptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).AcceptRideRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DriverService_AcceptRideRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).AcceptRideRequest(ctx, req.(*DriverAcceptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverService_RejectRideRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DriverRejectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServiceServer).RejectRideRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DriverService_RejectRideRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServiceServer).RejectRideRequest(ctx, req.(*DriverRejectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DriverService_ServiceDesc is the grpc.ServiceDesc for DriverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DriverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "comm.DriverService",
	HandlerType: (*DriverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CompleteRideRequest",
			Handler:    _DriverService_CompleteRideRequest_Handler,
		},
		{
			MethodName: "AcceptRideRequest",
			Handler:    _DriverService_AcceptRideRequest_Handler,
		},
		{
			MethodName: "RejectRideRequest",
			Handler:    _DriverService_RejectRideRequest_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AssignDriver",
			Handler:       _DriverService_AssignDriver_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "comm/comm.proto",
}