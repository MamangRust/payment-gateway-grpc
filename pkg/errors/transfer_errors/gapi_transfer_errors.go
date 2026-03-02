package transfer_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTransferNotFound              = errors.NewGrpcError("Transfer not found", int(codes.NotFound))
	ErrGrpcTransferInvalidID             = errors.NewGrpcError("Invalid Transfer ID", int(codes.InvalidArgument))
	ErrGrpcInvalidUserID                 = errors.NewGrpcError("Invalid user ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber             = errors.NewGrpcError("Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth                  = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear                   = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))
	ErrGrpcValidateCreateTransferRequest = errors.NewGrpcError("Invalid input for create transfer", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransferRequest = errors.NewGrpcError("Invalid input for update transfer", int(codes.InvalidArgument))
)
