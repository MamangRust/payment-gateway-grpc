package withdraw_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcWithdrawNotFound              = errors.NewGrpcError("Withdraw not found", int(codes.NotFound))
	ErrGrpcWithdrawInvalidID             = errors.NewGrpcError("Invalid Withdraw ID", int(codes.InvalidArgument))
	ErrGrpcInvalidUserID                 = errors.NewGrpcError("Invalid user ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber             = errors.NewGrpcError("Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth                  = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear                   = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))
	ErrGrpcValidateCreateWithdrawRequest = errors.NewGrpcError("Invalid input for create withdraw", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateWithdrawRequest = errors.NewGrpcError("Invalid input for update withdraw", int(codes.InvalidArgument))
)
