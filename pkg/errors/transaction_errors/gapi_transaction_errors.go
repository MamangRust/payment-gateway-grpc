package transaction_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTransactionNotFound              = errors.NewGrpcError("Transaction not found", int(codes.NotFound))
	ErrGrpcTransactionInvalidID             = errors.NewGrpcError("Invalid Transaction ID", int(codes.InvalidArgument))
	ErrGrpcTransactionInvalidMerchantID     = errors.NewGrpcError("Invalid Transaction Merchant ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber                = errors.NewGrpcError("Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth                     = errors.NewGrpcError("Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear                      = errors.NewGrpcError("Invalid year", int(codes.InvalidArgument))
	ErrGrpcValidateCreateTransactionRequest = errors.NewGrpcError("Invalid input for create transaction", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransactionRequest = errors.NewGrpcError("Invalid input for update transaction", int(codes.InvalidArgument))
)
