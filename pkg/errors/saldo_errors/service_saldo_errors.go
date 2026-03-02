package saldo_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllSaldos = errors.NewErrorResponse("Failed to fetch saldos", http.StatusInternalServerError)
	ErrFailedSaldoNotFound = errors.NewErrorResponse("Saldo not found", http.StatusNotFound)

	ErrFailedFindMonthlyTotalSaldoBalance = errors.NewErrorResponse("Failed to fetch monthly total saldo balance", http.StatusInternalServerError)
	ErrFailedFindYearTotalSaldoBalance    = errors.NewErrorResponse("Failed to fetch yearly total saldo balance", http.StatusInternalServerError)
	ErrFailedFindMonthlySaldoBalances     = errors.NewErrorResponse("Failed to fetch monthly saldo balances", http.StatusInternalServerError)
	ErrFailedFindYearlySaldoBalances      = errors.NewErrorResponse("Failed to fetch yearly saldo balances", http.StatusInternalServerError)

	ErrFailedFindSaldoByCardNumber = errors.NewErrorResponse("Failed to find saldo by card number", http.StatusInternalServerError)
	ErrFailedFindActiveSaldos      = errors.NewErrorResponse("Failed to fetch active saldos", http.StatusInternalServerError)
	ErrFailedFindTrashedSaldos     = errors.NewErrorResponse("Failed to fetch trashed saldos", http.StatusInternalServerError)

	ErrFailedCreateSaldo = errors.NewErrorResponse("Failed to create saldo", http.StatusInternalServerError)
	ErrFailedUpdateSaldo = errors.NewErrorResponse("Failed to update saldo", http.StatusInternalServerError)

	ErrFailedInsufficientBalance = errors.NewErrorResponse("Failed to Insuffient balance", http.StatusBadRequest)

	ErrFailedTrashSaldo              = errors.NewErrorResponse("Failed to trash saldo", http.StatusInternalServerError)
	ErrFailedRestoreSaldo            = errors.NewErrorResponse("Failed to restore saldo", http.StatusInternalServerError)
	ErrFailedDeleteSaldoPermanent    = errors.NewErrorResponse("Failed to permanently delete saldo", http.StatusInternalServerError)
	ErrFailedRestoreAllSaldo         = errors.NewErrorResponse("Failed to restore all saldos", http.StatusInternalServerError)
	ErrFailedDeleteAllSaldoPermanent = errors.NewErrorResponse("Failed to permanently delete all saldos", http.StatusInternalServerError)
)
