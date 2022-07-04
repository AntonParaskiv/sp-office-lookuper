// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

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

// TransferBoxApiClient is the client API for TransferBoxApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransferBoxApiClient interface {
	GetSortpointId(ctx context.Context, in *GetSortpointIdRequest, opts ...grpc.CallOption) (*GetSortpointIdResponse, error)
}

type transferBoxApiClient struct {
	cc grpc.ClientConnInterface
}

func NewTransferBoxApiClient(cc grpc.ClientConnInterface) TransferBoxApiClient {
	return &transferBoxApiClient{cc}
}

func (c *transferBoxApiClient) GetSortpointId(ctx context.Context, in *GetSortpointIdRequest, opts ...grpc.CallOption) (*GetSortpointIdResponse, error) {
	out := new(GetSortpointIdResponse)
	err := c.cc.Invoke(ctx, "/TransferBoxApi/GetSortpointId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransferBoxApiServer is the server API for TransferBoxApi service.
// All implementations must embed UnimplementedTransferBoxApiServer
// for forward compatibility
type TransferBoxApiServer interface {
	GetSortpointId(context.Context, *GetSortpointIdRequest) (*GetSortpointIdResponse, error)
	mustEmbedUnimplementedTransferBoxApiServer()
}

// UnimplementedTransferBoxApiServer must be embedded to have forward compatible implementations.
type UnimplementedTransferBoxApiServer struct {
}

func (UnimplementedTransferBoxApiServer) GetSortpointId(context.Context, *GetSortpointIdRequest) (*GetSortpointIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSortpointId not implemented")
}
func (UnimplementedTransferBoxApiServer) mustEmbedUnimplementedTransferBoxApiServer() {}

// UnsafeTransferBoxApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransferBoxApiServer will
// result in compilation errors.
type UnsafeTransferBoxApiServer interface {
	mustEmbedUnimplementedTransferBoxApiServer()
}

func RegisterTransferBoxApiServer(s grpc.ServiceRegistrar, srv TransferBoxApiServer) {
	s.RegisterService(&TransferBoxApi_ServiceDesc, srv)
}

func _TransferBoxApi_GetSortpointId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSortpointIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransferBoxApiServer).GetSortpointId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TransferBoxApi/GetSortpointId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransferBoxApiServer).GetSortpointId(ctx, req.(*GetSortpointIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransferBoxApi_ServiceDesc is the grpc.ServiceDesc for TransferBoxApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransferBoxApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TransferBoxApi",
	HandlerType: (*TransferBoxApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSortpointId",
			Handler:    _TransferBoxApi_GetSortpointId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/protobuf/transferBoxApi.proto",
}
