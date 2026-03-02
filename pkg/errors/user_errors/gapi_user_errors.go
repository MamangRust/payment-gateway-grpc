package user_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcUserInvalidId = errors.NewGrpcError("Invalid User ID", int(codes.NotFound))

	ErrGrpcValidateCreateUser = errors.NewGrpcError("validation failed: invalid create User request", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateUser = errors.NewGrpcError("validation failed: invalid update User request", int(codes.InvalidArgument))
)
