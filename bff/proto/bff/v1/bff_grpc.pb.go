// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v4.24.3
// source: bff/v1/bff.proto

package bffv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ScanTasksService_CreateTaskFromFile_FullMethodName  = "/bff.v1.ScanTasksService/CreateTaskFromFile"
	ScanTasksService_UploadOriginalVideo_FullMethodName = "/bff.v1.ScanTasksService/UploadOriginalVideo"
	ScanTasksService_GetTasksPreview_FullMethodName     = "/bff.v1.ScanTasksService/GetTasksPreview"
	ScanTasksService_GetTask_FullMethodName             = "/bff.v1.ScanTasksService/GetTask"
)

// ScanTasksServiceClient is the client API for ScanTasksService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScanTasksServiceClient interface {
	CreateTaskFromFile(ctx context.Context, in *CreateTaskFromFileRequest, opts ...grpc.CallOption) (*CreateTaskFromFileResponse, error)
	UploadOriginalVideo(ctx context.Context, in *UploadOriginalVideoRequest, opts ...grpc.CallOption) (*UploadOriginalVideoResponse, error)
	GetTasksPreview(ctx context.Context, in *GetTasksPreviewRequest, opts ...grpc.CallOption) (*GetTasksPreviewResponse, error)
	GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error)
}

type scanTasksServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScanTasksServiceClient(cc grpc.ClientConnInterface) ScanTasksServiceClient {
	return &scanTasksServiceClient{cc}
}

func (c *scanTasksServiceClient) CreateTaskFromFile(ctx context.Context, in *CreateTaskFromFileRequest, opts ...grpc.CallOption) (*CreateTaskFromFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateTaskFromFileResponse)
	err := c.cc.Invoke(ctx, ScanTasksService_CreateTaskFromFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scanTasksServiceClient) UploadOriginalVideo(ctx context.Context, in *UploadOriginalVideoRequest, opts ...grpc.CallOption) (*UploadOriginalVideoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadOriginalVideoResponse)
	err := c.cc.Invoke(ctx, ScanTasksService_UploadOriginalVideo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scanTasksServiceClient) GetTasksPreview(ctx context.Context, in *GetTasksPreviewRequest, opts ...grpc.CallOption) (*GetTasksPreviewResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTasksPreviewResponse)
	err := c.cc.Invoke(ctx, ScanTasksService_GetTasksPreview_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scanTasksServiceClient) GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTaskResponse)
	err := c.cc.Invoke(ctx, ScanTasksService_GetTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScanTasksServiceServer is the server API for ScanTasksService service.
// All implementations must embed UnimplementedScanTasksServiceServer
// for forward compatibility
type ScanTasksServiceServer interface {
	CreateTaskFromFile(context.Context, *CreateTaskFromFileRequest) (*CreateTaskFromFileResponse, error)
	UploadOriginalVideo(context.Context, *UploadOriginalVideoRequest) (*UploadOriginalVideoResponse, error)
	GetTasksPreview(context.Context, *GetTasksPreviewRequest) (*GetTasksPreviewResponse, error)
	GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error)
	mustEmbedUnimplementedScanTasksServiceServer()
}

// UnimplementedScanTasksServiceServer must be embedded to have forward compatible implementations.
type UnimplementedScanTasksServiceServer struct {
}

func (UnimplementedScanTasksServiceServer) CreateTaskFromFile(context.Context, *CreateTaskFromFileRequest) (*CreateTaskFromFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTaskFromFile not implemented")
}
func (UnimplementedScanTasksServiceServer) UploadOriginalVideo(context.Context, *UploadOriginalVideoRequest) (*UploadOriginalVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadOriginalVideo not implemented")
}
func (UnimplementedScanTasksServiceServer) GetTasksPreview(context.Context, *GetTasksPreviewRequest) (*GetTasksPreviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTasksPreview not implemented")
}
func (UnimplementedScanTasksServiceServer) GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTask not implemented")
}
func (UnimplementedScanTasksServiceServer) mustEmbedUnimplementedScanTasksServiceServer() {}

// UnsafeScanTasksServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScanTasksServiceServer will
// result in compilation errors.
type UnsafeScanTasksServiceServer interface {
	mustEmbedUnimplementedScanTasksServiceServer()
}

func RegisterScanTasksServiceServer(s grpc.ServiceRegistrar, srv ScanTasksServiceServer) {
	s.RegisterService(&ScanTasksService_ServiceDesc, srv)
}

func _ScanTasksService_CreateTaskFromFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskFromFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScanTasksServiceServer).CreateTaskFromFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScanTasksService_CreateTaskFromFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanTasksServiceServer).CreateTaskFromFile(ctx, req.(*CreateTaskFromFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScanTasksService_UploadOriginalVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadOriginalVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScanTasksServiceServer).UploadOriginalVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScanTasksService_UploadOriginalVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanTasksServiceServer).UploadOriginalVideo(ctx, req.(*UploadOriginalVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScanTasksService_GetTasksPreview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTasksPreviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScanTasksServiceServer).GetTasksPreview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScanTasksService_GetTasksPreview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanTasksServiceServer).GetTasksPreview(ctx, req.(*GetTasksPreviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScanTasksService_GetTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScanTasksServiceServer).GetTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScanTasksService_GetTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanTasksServiceServer).GetTask(ctx, req.(*GetTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ScanTasksService_ServiceDesc is the grpc.ServiceDesc for ScanTasksService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScanTasksService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bff.v1.ScanTasksService",
	HandlerType: (*ScanTasksServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTaskFromFile",
			Handler:    _ScanTasksService_CreateTaskFromFile_Handler,
		},
		{
			MethodName: "UploadOriginalVideo",
			Handler:    _ScanTasksService_UploadOriginalVideo_Handler,
		},
		{
			MethodName: "GetTasksPreview",
			Handler:    _ScanTasksService_GetTasksPreview_Handler,
		},
		{
			MethodName: "GetTask",
			Handler:    _ScanTasksService_GetTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bff/v1/bff.proto",
}
