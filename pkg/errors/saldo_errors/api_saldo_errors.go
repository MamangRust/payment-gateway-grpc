package saldo_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidSaldoID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid saldo ID", http.StatusBadRequest)
	}
	ErrApiInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid month value", http.StatusBadRequest)
	}

	ErrApiInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid year value", http.StatusBadRequest)
	}

	ErrApiInvalidCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid card-number value", http.StatusBadRequest)
	}

	ErrApiFailedFindAllSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all saldos", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllSaldoTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all saldos trashed", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllSaldoActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all saldos active", http.StatusInternalServerError)
	}

	ErrApiFailedFindByIdSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch saldo by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTotalSaldoBalance = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly total saldo balance", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearTotalSaldoBalance = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly total saldo balance", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlySaldoBalances = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly saldo balances", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlySaldoBalances = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly saldo balances", http.StatusInternalServerError)
	}

	ErrApiFailedFindByCardNumberSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch saldo by card number", http.StatusInternalServerError)
	}

	ErrApiFailedCreateSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create saldo", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update saldo", http.StatusInternalServerError)
	}

	ErrApiValidateCreateSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create Saldo request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update Saldo request", http.StatusBadRequest)
	}

	ErrApiBindCreateSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create Saldo request", http.StatusBadRequest)
	}

	ErrApiBindUpdateSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update Saldo request", http.StatusBadRequest)
	}

	ErrApiFailedTrashSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash saldo", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore saldo", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteSaldoPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete saldo", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllSaldo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all saldos", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllSaldoPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all saldos", http.StatusInternalServerError)
	}
)
