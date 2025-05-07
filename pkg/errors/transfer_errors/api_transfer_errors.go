package transfer_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiTransferInvalidID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid Transfer ID", http.StatusBadRequest)
	}

	ErrApiTransferInvalidMerchantID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid Transfer Merchant ID", http.StatusBadRequest)
	}

	ErrApiInvalidCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid card number", http.StatusBadRequest)
	}
	ErrApiInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid month", http.StatusBadRequest)
	}

	ErrApiInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid year", http.StatusBadRequest)
	}

	ErrApiFailedFindAllTransfers = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all transfers", http.StatusInternalServerError)
	}
	ErrApiFailedFindByIdTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transfer by ID", http.StatusInternalServerError)
	}
	ErrApiFailedFindByTransferFrom = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transfers by transfer_from", http.StatusInternalServerError)
	}
	ErrApiFailedFindByTransferTo = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transfers by transfer_to", http.StatusInternalServerError)
	}
	ErrApiFailedFindByActiveTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active transfers", http.StatusInternalServerError)
	}
	ErrApiFailedFindByTrashedTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed transfers", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTransferStatusSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful transfers", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferStatusSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful transfers", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferStatusFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed transfers", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferStatusFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed transfers", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferStatusSuccessByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful transfers by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferStatusSuccessByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful transfers by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferStatusFailedByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed transfers by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferStatusFailedByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed transfers by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTransferAmounts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer amounts", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferAmounts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer amounts", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferAmountsBySenderCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer amounts by sender card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferAmountsByReceiverCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer amounts by receiver card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferAmountsBySenderCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer amounts by sender card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferAmountsByReceiverCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer amounts by receiver card number", http.StatusInternalServerError)
	}

	ErrApiFailedCreateTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create transfer", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update transfer", http.StatusInternalServerError)
	}

	ErrApiBindCreateTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create transfer request", http.StatusBadRequest)
	}

	ErrApiBindUpdateTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update transfer request", http.StatusBadRequest)
	}

	ErrApiValidateCreateTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create transfer request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update transfer request", http.StatusBadRequest)
	}

	ErrApiFailedTrashedTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash transfer", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore transfer", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteTransferPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete transfer", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllTransfer = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all transfers", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllTransferPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all transfers", http.StatusInternalServerError)
	}
)
