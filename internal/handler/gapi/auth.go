package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"context"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	mapping     protomapper.AuthProtoMapper
}

func NewAuthHandleGrpc(auth service.AuthService, mapping protomapper.AuthProtoMapper) *authHandleGrpc {
	return &authHandleGrpc{authService: auth, mapping: mapping}
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.authService.Login(request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseLogin("success", "Login successful", res), nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	res, err := s.authService.RefreshToken(req.RefreshToken)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseRefreshToken("success", "Refresh token successful", res), nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	res, err := s.authService.GetMe(req.AccessToken)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseGetMe("success", "Refresh token successful", res), nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	request := &requests.CreateUserRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.authService.Register(request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseRegister("success", "Registration successful", res), nil
}
