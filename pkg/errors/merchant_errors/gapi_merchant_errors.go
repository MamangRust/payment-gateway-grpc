package merchant_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantNotFound       = errors.NewGrpcError("Merchant not found", int(codes.NotFound))
	ErrGrpcMerchantInvalidID      = errors.NewGrpcError("Invalid Merchant ID", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidUserID  = errors.NewGrpcError("Invalid Merchant User ID", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidApiKey  = errors.NewGrpcError("Invalid Merchant Api Key", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidMonth   = errors.NewGrpcError("Invalid Merchant Month", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidYear    = errors.NewGrpcError("Invalid Merchant Year", int(codes.InvalidArgument))
	ErrGrpcValidateCreateMerchant = errors.NewGrpcError("Invalid input for create merchant", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchant = errors.NewGrpcError("Invalid input for update merchant", int(codes.InvalidArgument))
)
