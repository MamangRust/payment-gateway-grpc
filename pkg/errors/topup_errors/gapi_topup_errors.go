package topup_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTopupNotFound       = errors.NewGrpcError("Topup not found", int(codes.NotFound))
	ErrGrpcTopupInvalidID      = errors.NewGrpcError("Invalid Topup ID", int(codes.InvalidArgument))
	ErrGrpcTopupInvalidMonth   = errors.NewGrpcError("Invalid Topup Month", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber   = errors.NewGrpcError("Invalid card number", int(codes.InvalidArgument))
	ErrGrpcTopupInvalidYear    = errors.NewGrpcError("Invalid Topup Year", int(codes.InvalidArgument))
	ErrGrpcValidateCreateTopup = errors.NewGrpcError("Invalid input for create topup", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTopup = errors.NewGrpcError("Invalid input for update topup", int(codes.InvalidArgument))
)
