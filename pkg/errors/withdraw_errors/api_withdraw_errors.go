package withdraw_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiWithdrawInvalidID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid Withdraw ID", http.StatusBadRequest)
	}

	ErrApiWithdrawInvalidUserID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid Withdraw Merchant ID", http.StatusBadRequest)
	}

	ErrApiInvalidCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid card number", http.StatusBadRequest)
	}

	ErrApiInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid month", http.StatusBadRequest)
	}

	ErrApiInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "year", "Invalid year", http.StatusBadRequest)
	}

	ErrApiWithdrawNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Withdraw not found", http.StatusNotFound)
	}

	ErrApiFailedFindAllWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedFindAllWithdrawByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch withdraws by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindByIdWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch withdraw by ID", http.StatusInternalServerError)
	}
	ErrApiFailedFindByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch withdraws using card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindByActiveWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedFindByTrashedWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed withdraws", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyWithdrawStatusSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawStatusSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyWithdrawStatusFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawStatusFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyWithdrawStatusSuccessCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly successful withdraws by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawStatusSuccessCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly successful withdraws by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyWithdrawStatusFailedCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly failed withdraws by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawStatusFailedCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly failed withdraws by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyWithdraws = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly withdraw amounts", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdraws = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly withdraw amounts", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyWithdrawsByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly withdraw amounts by card number", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawsByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly withdraw amounts by card number", http.StatusInternalServerError)
	}

	ErrApiFailedCreateWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create withdraw", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update withdraw", http.StatusInternalServerError)
	}

	ErrApiBindCreateWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create withdraw request", http.StatusBadRequest)
	}

	ErrApiBindUpdateWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update withdraw request", http.StatusBadRequest)
	}

	ErrApiValidateCreateWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create withdraw request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update withdraw request", http.StatusBadRequest)
	}

	ErrApiFailedTrashedWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to move withdraw to trash", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore withdraw", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteWithdrawPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete withdraw", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllWithdraw = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all withdraws", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllWithdrawPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all withdraws", http.StatusInternalServerError)
	}
)
