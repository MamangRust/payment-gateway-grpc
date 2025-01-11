package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"context"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleGrpc struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
	mapping     protomapper.UserProtoMapper
}

func NewUserHandleGrpc(user service.UserService, mapper protomapper.UserProtoMapper) *userHandleGrpc {
	return &userHandleGrpc{userService: user, mapping: mapper}
}

func (s *userHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userService.FindAll(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch users: ",
		})
	}

	so := s.mapping.ToResponsesUser(users)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}

	return &pb.ApiResponsePaginationUser{
		Status:     "success",
		Message:    "Successfully fetched users",
		Data:       so,
		Pagination: paginationMeta,
	}, nil
}

func (s *userHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid user id",
		})
	}

	user, err := s.userService.FindByID(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch user: " + err.Message,
		})
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    s.mapping.ToResponseUser(user),
	}, nil

}

func (s *userHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userService.FindByActive(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch active users: " + err.Message,
		})
	}

	so := s.mapping.ToResponsesUserDeleteAt(users)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}

	return &pb.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active users",
		Data:       so,
		Pagination: paginationMeta,
	}, nil
}

func (s *userHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userService.FindByTrashed(page, pageSize, search)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch trashed users: " + err.Message,
		})
	}

	so := s.mapping.ToResponsesUserDeleteAt(users)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}

	return &pb.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched users",
		Data:       so,
		Pagination: paginationMeta,
	}, nil
}

func (s *userHandleGrpc) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create user: ",
		})
	}

	user, err := s.userService.CreateUser(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to create user: ",
		})
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    s.mapping.ToResponseUser(user),
	}, nil
}

func (s *userHandleGrpc) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid user id",
		})
	}

	req := &requests.UpdateUserRequest{
		UserID:          int(request.GetId()),
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update user: ",
		})
	}

	user, err := s.userService.UpdateUser(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to update user: " + err.Message,
		})
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user",
		Data:    s.mapping.ToResponseUser(user),
	}, nil
}

func (s *userHandleGrpc) TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid user id",
		})
	}

	user, err := s.userService.TrashedUser(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed user: " + err.Message,
		})
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully trashed user",
		Data:    s.mapping.ToResponseUser(user),
	}, nil
}

func (s *userHandleGrpc) RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid user id",
		})
	}

	user, err := s.userService.RestoreUser(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore user: " + err.Message,
		})
	}

	return &pb.ApiResponseUser{
		Status:  "success",
		Message: "Successfully restored user",
		Data:    s.mapping.ToResponseUser(user),
	}, nil
}

func (s *userHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	if request.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Invalid user id",
		})
	}

	_, err := s.userService.DeleteUserPermanent(int(request.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete user permanently: " + err.Message,
		})
	}

	return &pb.ApiResponseUserDelete{
		Status:  "success",
		Message: "Successfully deleted user permanently",
	}, nil
}

func (s *userHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.userService.RestoreAllUser()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all user: ",
		})
	}

	return &pb.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully restore all user",
	}, nil
}

func (s *userHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.userService.DeleteAllUserPermanent()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", &pb.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete user permanent: ",
		})
	}

	return &pb.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully delete user permanent",
	}, nil
}
