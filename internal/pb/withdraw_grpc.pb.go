// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: withdraw.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	WithdrawService_FindAllWithdraw_FullMethodName         = "/pb.WithdrawService/FindAllWithdraw"
	WithdrawService_FindByIdWithdraw_FullMethodName        = "/pb.WithdrawService/FindByIdWithdraw"
	WithdrawService_FindByCardNumber_FullMethodName        = "/pb.WithdrawService/FindByCardNumber"
	WithdrawService_FindByActive_FullMethodName            = "/pb.WithdrawService/FindByActive"
	WithdrawService_FindByTrashed_FullMethodName           = "/pb.WithdrawService/FindByTrashed"
	WithdrawService_CreateWithdraw_FullMethodName          = "/pb.WithdrawService/CreateWithdraw"
	WithdrawService_UpdateWithdraw_FullMethodName          = "/pb.WithdrawService/UpdateWithdraw"
	WithdrawService_TrashedWithdraw_FullMethodName         = "/pb.WithdrawService/TrashedWithdraw"
	WithdrawService_RestoreWithdraw_FullMethodName         = "/pb.WithdrawService/RestoreWithdraw"
	WithdrawService_DeleteWithdrawPermanent_FullMethodName = "/pb.WithdrawService/DeleteWithdrawPermanent"
)

// WithdrawServiceClient is the client API for WithdrawService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WithdrawServiceClient interface {
	FindAllWithdraw(ctx context.Context, in *FindAllWithdrawRequest, opts ...grpc.CallOption) (*ApiResponsePaginationWithdraw, error)
	FindByIdWithdraw(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error)
	FindByCardNumber(ctx context.Context, in *FindByCardNumberRequest, opts ...grpc.CallOption) (*ApiResponsesWithdraw, error)
	FindByActive(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponsesWithdraw, error)
	FindByTrashed(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponsesWithdraw, error)
	CreateWithdraw(ctx context.Context, in *CreateWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error)
	UpdateWithdraw(ctx context.Context, in *UpdateWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error)
	TrashedWithdraw(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error)
	RestoreWithdraw(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdrawDelete, error)
}

type withdrawServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWithdrawServiceClient(cc grpc.ClientConnInterface) WithdrawServiceClient {
	return &withdrawServiceClient{cc}
}

func (c *withdrawServiceClient) FindAllWithdraw(ctx context.Context, in *FindAllWithdrawRequest, opts ...grpc.CallOption) (*ApiResponsePaginationWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponsePaginationWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_FindAllWithdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) FindByIdWithdraw(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_FindByIdWithdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) FindByCardNumber(ctx context.Context, in *FindByCardNumberRequest, opts ...grpc.CallOption) (*ApiResponsesWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponsesWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_FindByCardNumber_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) FindByActive(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponsesWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponsesWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_FindByActive_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) FindByTrashed(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponsesWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponsesWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_FindByTrashed_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) CreateWithdraw(ctx context.Context, in *CreateWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_CreateWithdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) UpdateWithdraw(ctx context.Context, in *UpdateWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_UpdateWithdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) TrashedWithdraw(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_TrashedWithdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) RestoreWithdraw(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdraw, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseWithdraw)
	err := c.cc.Invoke(ctx, WithdrawService_RestoreWithdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawServiceClient) DeleteWithdrawPermanent(ctx context.Context, in *FindByIdWithdrawRequest, opts ...grpc.CallOption) (*ApiResponseWithdrawDelete, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseWithdrawDelete)
	err := c.cc.Invoke(ctx, WithdrawService_DeleteWithdrawPermanent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WithdrawServiceServer is the server API for WithdrawService service.
// All implementations must embed UnimplementedWithdrawServiceServer
// for forward compatibility.
type WithdrawServiceServer interface {
	FindAllWithdraw(context.Context, *FindAllWithdrawRequest) (*ApiResponsePaginationWithdraw, error)
	FindByIdWithdraw(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdraw, error)
	FindByCardNumber(context.Context, *FindByCardNumberRequest) (*ApiResponsesWithdraw, error)
	FindByActive(context.Context, *emptypb.Empty) (*ApiResponsesWithdraw, error)
	FindByTrashed(context.Context, *emptypb.Empty) (*ApiResponsesWithdraw, error)
	CreateWithdraw(context.Context, *CreateWithdrawRequest) (*ApiResponseWithdraw, error)
	UpdateWithdraw(context.Context, *UpdateWithdrawRequest) (*ApiResponseWithdraw, error)
	TrashedWithdraw(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdraw, error)
	RestoreWithdraw(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdraw, error)
	DeleteWithdrawPermanent(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdrawDelete, error)
	mustEmbedUnimplementedWithdrawServiceServer()
}

// UnimplementedWithdrawServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedWithdrawServiceServer struct{}

func (UnimplementedWithdrawServiceServer) FindAllWithdraw(context.Context, *FindAllWithdrawRequest) (*ApiResponsePaginationWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllWithdraw not implemented")
}
func (UnimplementedWithdrawServiceServer) FindByIdWithdraw(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByIdWithdraw not implemented")
}
func (UnimplementedWithdrawServiceServer) FindByCardNumber(context.Context, *FindByCardNumberRequest) (*ApiResponsesWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByCardNumber not implemented")
}
func (UnimplementedWithdrawServiceServer) FindByActive(context.Context, *emptypb.Empty) (*ApiResponsesWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByActive not implemented")
}
func (UnimplementedWithdrawServiceServer) FindByTrashed(context.Context, *emptypb.Empty) (*ApiResponsesWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByTrashed not implemented")
}
func (UnimplementedWithdrawServiceServer) CreateWithdraw(context.Context, *CreateWithdrawRequest) (*ApiResponseWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWithdraw not implemented")
}
func (UnimplementedWithdrawServiceServer) UpdateWithdraw(context.Context, *UpdateWithdrawRequest) (*ApiResponseWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateWithdraw not implemented")
}
func (UnimplementedWithdrawServiceServer) TrashedWithdraw(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TrashedWithdraw not implemented")
}
func (UnimplementedWithdrawServiceServer) RestoreWithdraw(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdraw, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestoreWithdraw not implemented")
}
func (UnimplementedWithdrawServiceServer) DeleteWithdrawPermanent(context.Context, *FindByIdWithdrawRequest) (*ApiResponseWithdrawDelete, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWithdrawPermanent not implemented")
}
func (UnimplementedWithdrawServiceServer) mustEmbedUnimplementedWithdrawServiceServer() {}
func (UnimplementedWithdrawServiceServer) testEmbeddedByValue()                         {}

// UnsafeWithdrawServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WithdrawServiceServer will
// result in compilation errors.
type UnsafeWithdrawServiceServer interface {
	mustEmbedUnimplementedWithdrawServiceServer()
}

func RegisterWithdrawServiceServer(s grpc.ServiceRegistrar, srv WithdrawServiceServer) {
	// If the following call pancis, it indicates UnimplementedWithdrawServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&WithdrawService_ServiceDesc, srv)
}

func _WithdrawService_FindAllWithdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAllWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).FindAllWithdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_FindAllWithdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).FindAllWithdraw(ctx, req.(*FindAllWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_FindByIdWithdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).FindByIdWithdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_FindByIdWithdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).FindByIdWithdraw(ctx, req.(*FindByIdWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_FindByCardNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByCardNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).FindByCardNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_FindByCardNumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).FindByCardNumber(ctx, req.(*FindByCardNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_FindByActive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).FindByActive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_FindByActive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).FindByActive(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_FindByTrashed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).FindByTrashed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_FindByTrashed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).FindByTrashed(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_CreateWithdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).CreateWithdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_CreateWithdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).CreateWithdraw(ctx, req.(*CreateWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_UpdateWithdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).UpdateWithdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_UpdateWithdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).UpdateWithdraw(ctx, req.(*UpdateWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_TrashedWithdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).TrashedWithdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_TrashedWithdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).TrashedWithdraw(ctx, req.(*FindByIdWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_RestoreWithdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).RestoreWithdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_RestoreWithdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).RestoreWithdraw(ctx, req.(*FindByIdWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WithdrawService_DeleteWithdrawPermanent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawServiceServer).DeleteWithdrawPermanent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WithdrawService_DeleteWithdrawPermanent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawServiceServer).DeleteWithdrawPermanent(ctx, req.(*FindByIdWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WithdrawService_ServiceDesc is the grpc.ServiceDesc for WithdrawService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WithdrawService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.WithdrawService",
	HandlerType: (*WithdrawServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAllWithdraw",
			Handler:    _WithdrawService_FindAllWithdraw_Handler,
		},
		{
			MethodName: "FindByIdWithdraw",
			Handler:    _WithdrawService_FindByIdWithdraw_Handler,
		},
		{
			MethodName: "FindByCardNumber",
			Handler:    _WithdrawService_FindByCardNumber_Handler,
		},
		{
			MethodName: "FindByActive",
			Handler:    _WithdrawService_FindByActive_Handler,
		},
		{
			MethodName: "FindByTrashed",
			Handler:    _WithdrawService_FindByTrashed_Handler,
		},
		{
			MethodName: "CreateWithdraw",
			Handler:    _WithdrawService_CreateWithdraw_Handler,
		},
		{
			MethodName: "UpdateWithdraw",
			Handler:    _WithdrawService_UpdateWithdraw_Handler,
		},
		{
			MethodName: "TrashedWithdraw",
			Handler:    _WithdrawService_TrashedWithdraw_Handler,
		},
		{
			MethodName: "RestoreWithdraw",
			Handler:    _WithdrawService_RestoreWithdraw_Handler,
		},
		{
			MethodName: "DeleteWithdrawPermanent",
			Handler:    _WithdrawService_DeleteWithdrawPermanent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "withdraw.proto",
}
