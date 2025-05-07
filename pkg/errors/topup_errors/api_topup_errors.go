package topup_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidTopupID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid topup ID", http.StatusBadRequest)
	}

	ErrApiInvalidCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid card number", http.StatusBadRequest)
	}

	ErrApiInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid month value", http.StatusBadRequest)
	}

	ErrApiInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid year value", http.StatusBadRequest)
	}

	ErrApiFailedFindAllTopups = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all topups", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllTopupsTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all topups trashed", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllTopupsActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all topups active", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllByCardNumberTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch topups by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindByIdTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch topup by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTopupStatusSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful topups", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupStatusSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful topups", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTopupStatusFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed topups", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupStatusFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed topups", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTopupStatusSuccessByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful topups by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupStatusSuccessByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful topups by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTopupStatusFailedByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed topups by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupStatusFailedByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed topups by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTopupMethods = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly topup methods", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupMethods = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly topup methods", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTopupMethodsByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly topup methods by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupMethodsByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly topup methods by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTopupAmounts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly topup amounts", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupAmounts = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly topup amounts", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTopupAmountsByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly topup amounts by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupAmountsByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly topup amounts by card number", http.StatusInternalServerError)
	}

	ErrApiFailedCreateTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create topup", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update topup", http.StatusInternalServerError)
	}

	ErrApiValidateCreateTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create topup request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update topup request", http.StatusBadRequest)
	}

	ErrApiBindCreateTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create topup request", http.StatusBadRequest)
	}

	ErrApiBindUpdateTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update topup request", http.StatusBadRequest)
	}

	ErrApiFailedTrashTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash topup", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore topup", http.StatusInternalServerError)
	}
	ErrApiFailedDeletePermanentTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete topup", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllTopup = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all topups", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllTopupPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all topups", http.StatusInternalServerError)
	}
)
