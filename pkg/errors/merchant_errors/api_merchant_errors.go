package merchant_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidCardID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid card ID", http.StatusBadRequest)
	}

	ErrApiInvalidMerchantID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid merchant ID", http.StatusBadRequest)
	}

	ErrApiInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid month value", http.StatusBadRequest)
	}

	ErrApiInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid year value", http.StatusBadRequest)
	}

	ErrApiInvalidUserID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid user ID", http.StatusBadRequest)
	}

	ErrApiInvalidApiKey = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid API key", http.StatusUnauthorized)
	}

	ErrApiFailedFindAllMerchants = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all merchants", http.StatusInternalServerError)
	}
	ErrApiFailedFindAllMerchantsActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all merchants active", http.StatusInternalServerError)
	}
	ErrApiFailedFindAllMerchantsTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all merchants trashed", http.StatusInternalServerError)
	}

	ErrApiFailedFindByUserId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch merchant by user ID", http.StatusInternalServerError)
	}
	ErrApiFailedFindByIdMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch merchant by ID", http.StatusInternalServerError)
	}
	ErrApiFailedFindByApiKeyMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch merchant by api key", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllTransactions = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all transactions", http.StatusInternalServerError)
	}
	ErrApiFailedFindAllTransactionByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transactions by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindAllTransactionByApikey = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch transactions by API key", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyPaymentMethodsMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly payment methods by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyPaymentMethodMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly payment methods by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyAmountMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly amount by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyAmountMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly amount by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTotalAmountMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly total amount by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTotalAmountMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly total amount by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyPaymentMethodByMerchants = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly payment methods by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyPaymentMethodByMerchants = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly payment methods by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyAmountByMerchants = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly amounts by merchant", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyAmountByMerchants = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly amounts by merchant", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyPaymentMethodByApikeys = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly payment methods by API key", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyPaymentMethodByApikeys = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly payment methods by API key", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyAmountByApikeys = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly amounts by API key", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyAmountByApikeys = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly amounts by API key", http.StatusInternalServerError)
	}

	ErrApiFailedCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create merchant", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update merchant", http.StatusInternalServerError)
	}

	ErrApiValidateCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create merchant request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update merchant request", http.StatusBadRequest)
	}

	ErrApiBindCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create merchant request", http.StatusBadRequest)
	}

	ErrApiBindUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update merchant request", http.StatusBadRequest)
	}

	ErrApiFailedTrashMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash merchant", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore merchant", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteMerchantPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete merchant", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all merchants", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllMerchantPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all merchants", http.StatusInternalServerError)
	}
)
