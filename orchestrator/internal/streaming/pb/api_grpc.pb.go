// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: streaming/api.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Orchestrator_Process_FullMethodName   = "/streaming.Orchestrator/Process"
	Orchestrator_GetResult_FullMethodName = "/streaming.Orchestrator/GetResult"
	Orchestrator_Cancel_FullMethodName    = "/streaming.Orchestrator/Cancel"
)

// OrchestratorClient is the client API for Orchestrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrchestratorClient interface {
	Process(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error)
	GetResult(ctx context.Context, in *GetResultReq, opts ...grpc.CallOption) (*Response, error)
	Cancel(ctx context.Context, in *CancelReq, opts ...grpc.CallOption) (*CancelResp, error)
}

type orchestratorClient struct {
	cc grpc.ClientConnInterface
}

func NewOrchestratorClient(cc grpc.ClientConnInterface) OrchestratorClient {
	return &orchestratorClient{cc}
}

func (c *orchestratorClient) Process(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Orchestrator_Process_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) GetResult(ctx context.Context, in *GetResultReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Orchestrator_GetResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) Cancel(ctx context.Context, in *CancelReq, opts ...grpc.CallOption) (*CancelResp, error) {
	out := new(CancelResp)
	err := c.cc.Invoke(ctx, Orchestrator_Cancel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrchestratorServer is the server API for Orchestrator service.
// All implementations must embed UnimplementedOrchestratorServer
// for forward compatibility
type OrchestratorServer interface {
	Process(context.Context, *Query) (*Response, error)
	GetResult(context.Context, *GetResultReq) (*Response, error)
	Cancel(context.Context, *CancelReq) (*CancelResp, error)
	mustEmbedUnimplementedOrchestratorServer()
}

// UnimplementedOrchestratorServer must be embedded to have forward compatible implementations.
type UnimplementedOrchestratorServer struct {
}

func (UnimplementedOrchestratorServer) Process(context.Context, *Query) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Process not implemented")
}
func (UnimplementedOrchestratorServer) GetResult(context.Context, *GetResultReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResult not implemented")
}
func (UnimplementedOrchestratorServer) Cancel(context.Context, *CancelReq) (*CancelResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cancel not implemented")
}
func (UnimplementedOrchestratorServer) mustEmbedUnimplementedOrchestratorServer() {}

// UnsafeOrchestratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrchestratorServer will
// result in compilation errors.
type UnsafeOrchestratorServer interface {
	mustEmbedUnimplementedOrchestratorServer()
}

func RegisterOrchestratorServer(s grpc.ServiceRegistrar, srv OrchestratorServer) {
	s.RegisterService(&Orchestrator_ServiceDesc, srv)
}

func _Orchestrator_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Orchestrator_Process_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).Process(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_GetResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResultReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).GetResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Orchestrator_GetResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).GetResult(ctx, req.(*GetResultReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_Cancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).Cancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Orchestrator_Cancel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).Cancel(ctx, req.(*CancelReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Orchestrator_ServiceDesc is the grpc.ServiceDesc for Orchestrator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Orchestrator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "streaming.Orchestrator",
	HandlerType: (*OrchestratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _Orchestrator_Process_Handler,
		},
		{
			MethodName: "GetResult",
			Handler:    _Orchestrator_GetResult_Handler,
		},
		{
			MethodName: "Cancel",
			Handler:    _Orchestrator_Cancel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "streaming/api.proto",
}
