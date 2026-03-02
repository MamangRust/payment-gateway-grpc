package card_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidCardID     = errors.NewGrpcError("Invalid card ID", int(codes.InvalidArgument))
	ErrGrpcInvalidUserID     = errors.NewGrpcError("Invalid user ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber = errors.NewGrpcError("Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth      = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear       = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))

	ErrGrpcValidateCreateCardRequest = errors.NewGrpcError("Invalid input for create card", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCardRequest = errors.NewGrpcError("Invalid input for update card", int(codes.InvalidArgument))
)
