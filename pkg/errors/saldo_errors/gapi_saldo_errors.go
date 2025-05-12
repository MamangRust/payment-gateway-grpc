package saldo_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcSaldoNotFound          = response.NewGrpcError("saldo", "Saldo not found", int(codes.NotFound))
	ErrGrpcSaldoInvalidID         = response.NewGrpcError("saldo", "Invalid Saldo ID", int(codes.InvalidArgument))
	ErrGrpcSaldoInvalidCardNumber = response.NewGrpcError("saldo", "Invalid Saldo Card Number", int(codes.InvalidArgument))
	ErrGrpcSaldoInvalidMonth      = response.NewGrpcError("saldo", "Invalid Saldo Month", int(codes.InvalidArgument))
	ErrGrpcSaldoInvalidYear       = response.NewGrpcError("saldo", "Invalid Saldo Year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllSaldo     = response.NewGrpcError("saldo", "Failed to fetch all saldo", int(codes.Internal))
	ErrGrpcFailedFindByIdSaldo    = response.NewGrpcError("saldo", "Failed to fetch saldo by ID", int(codes.Internal))
	ErrGrpcFailedFindByCardNumber = response.NewGrpcError("saldo", "Failed to fetch saldo by card number", int(codes.Internal))
	ErrGrpcFailedFindByActive     = response.NewGrpcError("saldo", "Failed to fetch active saldo", int(codes.Internal))
	ErrGrpcFailedFindByTrashed    = response.NewGrpcError("saldo", "Failed to fetch trashed saldo", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTotalSaldoBalance = response.NewGrpcError("saldo", "Failed to fetch monthly total saldo balance", int(codes.Internal))
	ErrGrpcFailedFindYearTotalSaldoBalance    = response.NewGrpcError("saldo", "Failed to fetch yearly total saldo balance", int(codes.Internal))
	ErrGrpcFailedFindMonthlySaldoBalances     = response.NewGrpcError("saldo", "Failed to fetch monthly saldo balances", int(codes.Internal))
	ErrGrpcFailedFindYearlySaldoBalances      = response.NewGrpcError("saldo", "Failed to fetch yearly saldo balances", int(codes.Internal))

	ErrGrpcFailedCreateSaldo   = response.NewGrpcError("saldo", "Failed to create saldo", int(codes.Internal))
	ErrGrpcFailedUpdateSaldo   = response.NewGrpcError("saldo", "Failed to update saldo", int(codes.Internal))
	ErrGrpcValidateCreateSaldo = response.NewGrpcError("saldo", "Invalid input for create saldo", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSaldo = response.NewGrpcError("saldo", "Invalid input for update saldo", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedSaldo         = response.NewGrpcError("saldo", "Failed to move saldo to trash", int(codes.Internal))
	ErrGrpcFailedRestoreSaldo         = response.NewGrpcError("saldo", "Failed to restore saldo", int(codes.Internal))
	ErrGrpcFailedDeleteSaldoPermanent = response.NewGrpcError("saldo", "Failed to permanently delete saldo", int(codes.Internal))

	ErrGrpcFailedRestoreAllSaldo         = response.NewGrpcError("saldo", "Failed to restore all saldo", int(codes.Internal))
	ErrGrpcFailedDeleteAllSaldoPermanent = response.NewGrpcError("saldo", "Failed to permanently delete all saldo", int(codes.Internal))
)
