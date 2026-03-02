package role_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcRoleNotFound  = errors.NewGrpcError("Role not found", int(codes.NotFound))
	ErrGrpcRoleInvalidId = errors.NewGrpcError("Invalid Role ID", int(codes.NotFound))

	ErrGrpcFailedFindAll     = errors.NewGrpcError("Failed to fetch Roles", int(codes.Internal))
	ErrGrpcFailedFindActive  = errors.NewGrpcError("Failed to fetch active Roles", int(codes.Internal))
	ErrGrpcFailedFindTrashed = errors.NewGrpcError("Failed to fetch trashed Roles", int(codes.Internal))

	ErrGrpcFailedCreateRole   = errors.NewGrpcError("Failed to create Role", int(codes.Internal))
	ErrGrpcFailedUpdateRole   = errors.NewGrpcError("Failed to update Role", int(codes.Internal))
	ErrGrpcValidateCreateRole = errors.NewGrpcError("validation failed: invalid create Role request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateRole = errors.NewGrpcError("validation failed: invalid update Role request", int(codes.InvalidArgument))
)
