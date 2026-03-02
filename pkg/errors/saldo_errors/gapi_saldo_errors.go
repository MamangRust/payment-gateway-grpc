package saldo_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcSaldoNotFound          = errors.NewGrpcError("Saldo not found", int(codes.NotFound))
	ErrGrpcSaldoInvalidID         = errors.NewGrpcError("Invalid Saldo ID", int(codes.InvalidArgument))
	ErrGrpcSaldoInvalidCardNumber = errors.NewGrpcError("Invalid Saldo Card Number", int(codes.InvalidArgument))
	ErrGrpcSaldoInvalidMonth      = errors.NewGrpcError("Invalid Saldo Month", int(codes.InvalidArgument))
	ErrGrpcSaldoInvalidYear       = errors.NewGrpcError("Invalid Saldo Year", int(codes.InvalidArgument))
	ErrGrpcValidateCreateSaldo    = errors.NewGrpcError("Invalid input for create saldo", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSaldo    = errors.NewGrpcError("Invalid input for update saldo", int(codes.InvalidArgument))
)
