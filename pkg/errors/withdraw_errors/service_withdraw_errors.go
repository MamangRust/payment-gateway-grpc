package withdraw_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllWithdraws       = response.NewErrorResponse("Failed to fetch all withdraws", http.StatusInternalServerError)
	ErrWithdrawNotFound             = response.NewErrorResponse("Withdraw not found", http.StatusNotFound)
	ErrFailedFindAllWithdrawsByCard = response.NewErrorResponse("Failed to fetch all withdraws by card number", http.StatusInternalServerError)

	ErrFailedFindMonthWithdrawStatusSuccess = response.NewErrorResponse("Failed to fetch monthly successful withdraws", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusSuccess  = response.NewErrorResponse("Failed to fetch yearly successful withdraws", http.StatusInternalServerError)
	ErrFailedFindMonthWithdrawStatusFailed  = response.NewErrorResponse("Failed to fetch monthly failed withdraws", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusFailed   = response.NewErrorResponse("Failed to fetch yearly failed withdraws", http.StatusInternalServerError)

	ErrFailedFindMonthWithdrawStatusSuccessByCard = response.NewErrorResponse("Failed to fetch monthly successful withdraws by card", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusSuccessByCard  = response.NewErrorResponse("Failed to fetch yearly successful withdraws by card", http.StatusInternalServerError)
	ErrFailedFindMonthWithdrawStatusFailedByCard  = response.NewErrorResponse("Failed to fetch monthly failed withdraws by card", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusFailedByCard   = response.NewErrorResponse("Failed to fetch yearly failed withdraws by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyWithdraws             = response.NewErrorResponse("Failed to fetch monthly withdraw amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdraws              = response.NewErrorResponse("Failed to fetch yearly withdraw amounts", http.StatusInternalServerError)
	ErrFailedFindMonthlyWithdrawsByCardNumber = response.NewErrorResponse("Failed to fetch monthly withdraw amounts by card", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdrawsByCardNumber  = response.NewErrorResponse("Failed to fetch yearly withdraw amounts by card", http.StatusInternalServerError)

	ErrFailedFindActiveWithdraws  = response.NewErrorResponse("Failed to fetch active withdraws", http.StatusInternalServerError)
	ErrFailedFindTrashedWithdraws = response.NewErrorResponse("Failed to fetch trashed withdraws", http.StatusInternalServerError)

	ErrFailedCreateWithdraw = response.NewErrorResponse("Failed to create withdraw", http.StatusInternalServerError)
	ErrFailedUpdateWithdraw = response.NewErrorResponse("Failed to update withdraw", http.StatusInternalServerError)

	ErrFailedTrashedWithdraw            = response.NewErrorResponse("Failed to trash withdraw", http.StatusInternalServerError)
	ErrFailedRestoreWithdraw            = response.NewErrorResponse("Failed to restore withdraw", http.StatusInternalServerError)
	ErrFailedDeleteWithdrawPermanent    = response.NewErrorResponse("Failed to permanently delete withdraw", http.StatusInternalServerError)
	ErrFailedRestoreAllWithdraw         = response.NewErrorResponse("Failed to restore all withdraws", http.StatusInternalServerError)
	ErrFailedDeleteAllWithdrawPermanent = response.NewErrorResponse("Failed to permanently delete all withdraws", http.StatusInternalServerError)
)
