package card_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidCardID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid card ID", http.StatusBadRequest)
	}

	ErrApiInvalidUserID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid user ID", http.StatusBadRequest)
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

	ErrApiFailedFindAllCards = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all cards", http.StatusInternalServerError)
	}
	ErrApiFailedFindByIdCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch card by ID", http.StatusInternalServerError)
	}

	ErrApiFailedDashboardCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch card dashboard", http.StatusInternalServerError)
	}
	ErrApiFailedDashboardCardByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch dashboard by card number", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyBalance = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly balance", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyBalance = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly balance", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTopupAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly topup amount", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly topup amount", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyWithdrawAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly withdraw amount", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly withdraw amount", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTransactionAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transaction amount", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransactionAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transaction amount", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTransferSenderAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer sender amount", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferSenderAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer sender amount", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferReceiverAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer receiver amount", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferReceiverAmount = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer receiver amount", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyBalanceByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly balance by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyBalanceByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly balance by card", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTopupAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly topup amount by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTopupAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly topup amount by card", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyWithdrawAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly withdraw amount by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyWithdrawAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly withdraw amount by card", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTransactionAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transaction amount by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransactionAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transaction amount by card", http.StatusInternalServerError)
	}

	ErrApiFailedFindMonthlyTransferSenderAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer sender amount by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferSenderAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer sender amount by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindMonthlyTransferReceiverAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch monthly transfer receiver amount by card", http.StatusInternalServerError)
	}
	ErrApiFailedFindYearlyTransferReceiverAmountByCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch yearly transfer receiver amount by card", http.StatusInternalServerError)
	}

	ErrApiFailedFindByUserIdCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch card by user ID", http.StatusInternalServerError)
	}
	ErrApiFailedFindByActiveCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch active cards", http.StatusInternalServerError)
	}
	ErrApiFailedFindByTrashedCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch trashed cards", http.StatusInternalServerError)
	}

	ErrApiFailedFindByCardNumber = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch card by card number", http.StatusInternalServerError)
	}

	ErrApiFailedCreateCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create card", http.StatusInternalServerError)
	}
	ErrApiFailedUpdateCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update card", http.StatusInternalServerError)
	}

	ErrApiValidateCreateCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create Card request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update Card request", http.StatusBadRequest)
	}

	ErrApiBindCreateCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create Card request", http.StatusBadRequest)
	}

	ErrApiBindUpdateCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update Card request", http.StatusBadRequest)
	}

	ErrApiFailedTrashCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash card", http.StatusInternalServerError)
	}
	ErrApiFailedRestoreCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore card", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteCardPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete card", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllCard = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all cards", http.StatusInternalServerError)
	}
	ErrApiFailedDeleteAllCardPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all cards", http.StatusInternalServerError)
	}
)
