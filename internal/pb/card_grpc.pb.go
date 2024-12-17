// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: card.proto

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
	CardService_FindAllCard_FullMethodName         = "/pb.CardService/FindAllCard"
	CardService_FindByIdCard_FullMethodName        = "/pb.CardService/FindByIdCard"
	CardService_FindByUserIdCard_FullMethodName    = "/pb.CardService/FindByUserIdCard"
	CardService_FindByActiveCard_FullMethodName    = "/pb.CardService/FindByActiveCard"
	CardService_FindByTrashedCard_FullMethodName   = "/pb.CardService/FindByTrashedCard"
	CardService_FindByCardNumber_FullMethodName    = "/pb.CardService/FindByCardNumber"
	CardService_CreateCard_FullMethodName          = "/pb.CardService/CreateCard"
	CardService_UpdateCard_FullMethodName          = "/pb.CardService/UpdateCard"
	CardService_TrashedCard_FullMethodName         = "/pb.CardService/TrashedCard"
	CardService_RestoreCard_FullMethodName         = "/pb.CardService/RestoreCard"
	CardService_DeleteCardPermanent_FullMethodName = "/pb.CardService/DeleteCardPermanent"
)

// CardServiceClient is the client API for CardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CardServiceClient interface {
	FindAllCard(ctx context.Context, in *FindAllCardRequest, opts ...grpc.CallOption) (*ApiResponsePaginationCard, error)
	FindByIdCard(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	FindByUserIdCard(ctx context.Context, in *FindByUserIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	FindByActiveCard(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponseCards, error)
	FindByTrashedCard(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponseCards, error)
	FindByCardNumber(ctx context.Context, in *FindByCardNumberRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	CreateCard(ctx context.Context, in *CreateCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	UpdateCard(ctx context.Context, in *UpdateCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	TrashedCard(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	RestoreCard(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error)
	DeleteCardPermanent(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCardDelete, error)
}

type cardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCardServiceClient(cc grpc.ClientConnInterface) CardServiceClient {
	return &cardServiceClient{cc}
}

func (c *cardServiceClient) FindAllCard(ctx context.Context, in *FindAllCardRequest, opts ...grpc.CallOption) (*ApiResponsePaginationCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponsePaginationCard)
	err := c.cc.Invoke(ctx, CardService_FindAllCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) FindByIdCard(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_FindByIdCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) FindByUserIdCard(ctx context.Context, in *FindByUserIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_FindByUserIdCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) FindByActiveCard(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponseCards, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCards)
	err := c.cc.Invoke(ctx, CardService_FindByActiveCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) FindByTrashedCard(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ApiResponseCards, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCards)
	err := c.cc.Invoke(ctx, CardService_FindByTrashedCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) FindByCardNumber(ctx context.Context, in *FindByCardNumberRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_FindByCardNumber_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) CreateCard(ctx context.Context, in *CreateCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_CreateCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) UpdateCard(ctx context.Context, in *UpdateCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_UpdateCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) TrashedCard(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_TrashedCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) RestoreCard(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCard)
	err := c.cc.Invoke(ctx, CardService_RestoreCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cardServiceClient) DeleteCardPermanent(ctx context.Context, in *FindByIdCardRequest, opts ...grpc.CallOption) (*ApiResponseCardDelete, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApiResponseCardDelete)
	err := c.cc.Invoke(ctx, CardService_DeleteCardPermanent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CardServiceServer is the server API for CardService service.
// All implementations must embed UnimplementedCardServiceServer
// for forward compatibility.
type CardServiceServer interface {
	FindAllCard(context.Context, *FindAllCardRequest) (*ApiResponsePaginationCard, error)
	FindByIdCard(context.Context, *FindByIdCardRequest) (*ApiResponseCard, error)
	FindByUserIdCard(context.Context, *FindByUserIdCardRequest) (*ApiResponseCard, error)
	FindByActiveCard(context.Context, *emptypb.Empty) (*ApiResponseCards, error)
	FindByTrashedCard(context.Context, *emptypb.Empty) (*ApiResponseCards, error)
	FindByCardNumber(context.Context, *FindByCardNumberRequest) (*ApiResponseCard, error)
	CreateCard(context.Context, *CreateCardRequest) (*ApiResponseCard, error)
	UpdateCard(context.Context, *UpdateCardRequest) (*ApiResponseCard, error)
	TrashedCard(context.Context, *FindByIdCardRequest) (*ApiResponseCard, error)
	RestoreCard(context.Context, *FindByIdCardRequest) (*ApiResponseCard, error)
	DeleteCardPermanent(context.Context, *FindByIdCardRequest) (*ApiResponseCardDelete, error)
	mustEmbedUnimplementedCardServiceServer()
}

// UnimplementedCardServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCardServiceServer struct{}

func (UnimplementedCardServiceServer) FindAllCard(context.Context, *FindAllCardRequest) (*ApiResponsePaginationCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllCard not implemented")
}
func (UnimplementedCardServiceServer) FindByIdCard(context.Context, *FindByIdCardRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByIdCard not implemented")
}
func (UnimplementedCardServiceServer) FindByUserIdCard(context.Context, *FindByUserIdCardRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByUserIdCard not implemented")
}
func (UnimplementedCardServiceServer) FindByActiveCard(context.Context, *emptypb.Empty) (*ApiResponseCards, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByActiveCard not implemented")
}
func (UnimplementedCardServiceServer) FindByTrashedCard(context.Context, *emptypb.Empty) (*ApiResponseCards, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByTrashedCard not implemented")
}
func (UnimplementedCardServiceServer) FindByCardNumber(context.Context, *FindByCardNumberRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByCardNumber not implemented")
}
func (UnimplementedCardServiceServer) CreateCard(context.Context, *CreateCardRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCard not implemented")
}
func (UnimplementedCardServiceServer) UpdateCard(context.Context, *UpdateCardRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCard not implemented")
}
func (UnimplementedCardServiceServer) TrashedCard(context.Context, *FindByIdCardRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TrashedCard not implemented")
}
func (UnimplementedCardServiceServer) RestoreCard(context.Context, *FindByIdCardRequest) (*ApiResponseCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestoreCard not implemented")
}
func (UnimplementedCardServiceServer) DeleteCardPermanent(context.Context, *FindByIdCardRequest) (*ApiResponseCardDelete, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCardPermanent not implemented")
}
func (UnimplementedCardServiceServer) mustEmbedUnimplementedCardServiceServer() {}
func (UnimplementedCardServiceServer) testEmbeddedByValue()                     {}

// UnsafeCardServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CardServiceServer will
// result in compilation errors.
type UnsafeCardServiceServer interface {
	mustEmbedUnimplementedCardServiceServer()
}

func RegisterCardServiceServer(s grpc.ServiceRegistrar, srv CardServiceServer) {
	// If the following call pancis, it indicates UnimplementedCardServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CardService_ServiceDesc, srv)
}

func _CardService_FindAllCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAllCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).FindAllCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_FindAllCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).FindAllCard(ctx, req.(*FindAllCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_FindByIdCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).FindByIdCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_FindByIdCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).FindByIdCard(ctx, req.(*FindByIdCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_FindByUserIdCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByUserIdCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).FindByUserIdCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_FindByUserIdCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).FindByUserIdCard(ctx, req.(*FindByUserIdCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_FindByActiveCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).FindByActiveCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_FindByActiveCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).FindByActiveCard(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_FindByTrashedCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).FindByTrashedCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_FindByTrashedCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).FindByTrashedCard(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_FindByCardNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByCardNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).FindByCardNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_FindByCardNumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).FindByCardNumber(ctx, req.(*FindByCardNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_CreateCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).CreateCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_CreateCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).CreateCard(ctx, req.(*CreateCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_UpdateCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).UpdateCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_UpdateCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).UpdateCard(ctx, req.(*UpdateCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_TrashedCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).TrashedCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_TrashedCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).TrashedCard(ctx, req.(*FindByIdCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_RestoreCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).RestoreCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_RestoreCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).RestoreCard(ctx, req.(*FindByIdCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CardService_DeleteCardPermanent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindByIdCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CardServiceServer).DeleteCardPermanent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CardService_DeleteCardPermanent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CardServiceServer).DeleteCardPermanent(ctx, req.(*FindByIdCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CardService_ServiceDesc is the grpc.ServiceDesc for CardService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.CardService",
	HandlerType: (*CardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAllCard",
			Handler:    _CardService_FindAllCard_Handler,
		},
		{
			MethodName: "FindByIdCard",
			Handler:    _CardService_FindByIdCard_Handler,
		},
		{
			MethodName: "FindByUserIdCard",
			Handler:    _CardService_FindByUserIdCard_Handler,
		},
		{
			MethodName: "FindByActiveCard",
			Handler:    _CardService_FindByActiveCard_Handler,
		},
		{
			MethodName: "FindByTrashedCard",
			Handler:    _CardService_FindByTrashedCard_Handler,
		},
		{
			MethodName: "FindByCardNumber",
			Handler:    _CardService_FindByCardNumber_Handler,
		},
		{
			MethodName: "CreateCard",
			Handler:    _CardService_CreateCard_Handler,
		},
		{
			MethodName: "UpdateCard",
			Handler:    _CardService_UpdateCard_Handler,
		},
		{
			MethodName: "TrashedCard",
			Handler:    _CardService_TrashedCard_Handler,
		},
		{
			MethodName: "RestoreCard",
			Handler:    _CardService_RestoreCard_Handler,
		},
		{
			MethodName: "DeleteCardPermanent",
			Handler:    _CardService_DeleteCardPermanent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "card.proto",
}
