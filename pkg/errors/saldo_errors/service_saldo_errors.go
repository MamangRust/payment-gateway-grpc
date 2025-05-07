package saldo_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllSaldos = response.NewErrorResponse("Failed to fetch saldos", http.StatusInternalServerError)
	ErrFailedSaldoNotFound = response.NewErrorResponse("Saldo not found", http.StatusNotFound)

	ErrFailedFindMonthlyTotalSaldoBalance = response.NewErrorResponse("Failed to fetch monthly total saldo balance", http.StatusInternalServerError)
	ErrFailedFindYearTotalSaldoBalance    = response.NewErrorResponse("Failed to fetch yearly total saldo balance", http.StatusInternalServerError)
	ErrFailedFindMonthlySaldoBalances     = response.NewErrorResponse("Failed to fetch monthly saldo balances", http.StatusInternalServerError)
	ErrFailedFindYearlySaldoBalances      = response.NewErrorResponse("Failed to fetch yearly saldo balances", http.StatusInternalServerError)

	ErrFailedFindSaldoByCardNumber = response.NewErrorResponse("Failed to find saldo by card number", http.StatusInternalServerError)
	ErrFailedFindActiveSaldos      = response.NewErrorResponse("Failed to fetch active saldos", http.StatusInternalServerError)
	ErrFailedFindTrashedSaldos     = response.NewErrorResponse("Failed to fetch trashed saldos", http.StatusInternalServerError)

	ErrFailedCreateSaldo = response.NewErrorResponse("Failed to create saldo", http.StatusInternalServerError)
	ErrFailedUpdateSaldo = response.NewErrorResponse("Failed to update saldo", http.StatusInternalServerError)

	ErrFailedTrashSaldo              = response.NewErrorResponse("Failed to trash saldo", http.StatusInternalServerError)
	ErrFailedRestoreSaldo            = response.NewErrorResponse("Failed to restore saldo", http.StatusInternalServerError)
	ErrFailedDeleteSaldoPermanent    = response.NewErrorResponse("Failed to permanently delete saldo", http.StatusInternalServerError)
	ErrFailedRestoreAllSaldo         = response.NewErrorResponse("Failed to restore all saldos", http.StatusInternalServerError)
	ErrFailedDeleteAllSaldoPermanent = response.NewErrorResponse("Failed to permanently delete all saldos", http.StatusInternalServerError)
)
