package transaction_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidTransactionApiKey = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid transaction api-key", http.StatusBadRequest)
	}

	ErrApiInvalidTransactionID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid transaction ID", http.StatusBadRequest)
	}

	ErrApiInvalidTransactionCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid transaction Card Number", http.StatusBadRequest)
	}

	ErrApiInvalidTransactionMerchantID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid transaction merchant ID", http.StatusBadRequest)
	}

	ErrApiInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid month value", http.StatusBadRequest)
	}

	ErrApiInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid year value", http.StatusBadRequest)
	}

	ErrApiFailedFindAllTransactions = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all transactions", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllTransactionsActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all transactions active", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllTransactionsTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all transactions active", http.StatusInternalServerError)
	}

	ErrApiFailedFindByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transactions by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transaction by ID", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlySuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful transactions", http.StatusInternalServerError)
	}

	ErrApiFailedYearlySuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful transactions", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlyFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed transactions", http.StatusInternalServerError)
	}

	ErrApiFailedYearlyFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed transactions", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlySuccessByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful transactions by card number", http.StatusInternalServerError)
	}

	ErrApiFailedYearlySuccessByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful transactions by card number", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlyFailedByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed transactions by card number", http.StatusInternalServerError)
	}

	ErrApiFailedYearlyFailedByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed transactions by card number", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlyMethods = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly payment methods", http.StatusInternalServerError)
	}

	ErrApiFailedYearlyMethods = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly payment methods", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlyAmounts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly amounts", http.StatusInternalServerError)
	}

	ErrApiFailedYearlyAmounts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly amounts", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlyMethodsByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly payment methods by card number", http.StatusInternalServerError)
	}

	ErrApiFailedYearlyMethodsByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly payment methods by card number", http.StatusInternalServerError)
	}

	ErrApiFailedMonthlyAmountsByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly amounts by card number", http.StatusInternalServerError)
	}

	ErrApiFailedYearlyAmountsByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly amounts by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindByMerchantID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transactions by merchant ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active transactions", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed transactions", http.StatusInternalServerError)
	}

	ErrApiFailedCreateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create transaction", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update transaction", http.StatusInternalServerError)
	}

	ErrApiValidateCreateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create transaction request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update transaction request", http.StatusBadRequest)
	}

	ErrApiBindCreateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create transaction request", http.StatusBadRequest)
	}

	ErrApiBindUpdateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update transaction request", http.StatusBadRequest)
	}

	ErrApiFailedRestoreTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore transaction", http.StatusInternalServerError)
	}

	ErrApiFailedTrashTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to move transaction to trash", http.StatusInternalServerError)
	}

	ErrApiFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete transaction", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllTransactions = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all transactions", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all trashed transactions", http.StatusInternalServerError)
	}
)
