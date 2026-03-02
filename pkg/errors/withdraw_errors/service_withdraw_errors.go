package withdraw_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllWithdraws       = errors.NewErrorResponse("Failed to fetch all withdraws", http.StatusInternalServerError)
	ErrWithdrawNotFound             = errors.NewErrorResponse("Withdraw not found", http.StatusNotFound)
	ErrFailedFindAllWithdrawsByCard = errors.NewErrorResponse("Failed to fetch all withdraws by card number", http.StatusInternalServerError)

	ErrFailedFindMonthWithdrawStatusSuccess = errors.NewErrorResponse("Failed to fetch monthly successful withdraws", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusSuccess  = errors.NewErrorResponse("Failed to fetch yearly successful withdraws", http.StatusInternalServerError)
	ErrFailedFindMonthWithdrawStatusFailed  = errors.NewErrorResponse("Failed to fetch monthly failed withdraws", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusFailed   = errors.NewErrorResponse("Failed to fetch yearly failed withdraws", http.StatusInternalServerError)

	ErrFailedFindMonthWithdrawStatusSuccessByCard = errors.NewErrorResponse("Failed to fetch monthly successful withdraws by card", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusSuccessByCard  = errors.NewErrorResponse("Failed to fetch yearly successful withdraws by card", http.StatusInternalServerError)
	ErrFailedFindMonthWithdrawStatusFailedByCard  = errors.NewErrorResponse("Failed to fetch monthly failed withdraws by card", http.StatusInternalServerError)
	ErrFailedFindYearWithdrawStatusFailedByCard   = errors.NewErrorResponse("Failed to fetch yearly failed withdraws by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyWithdraws             = errors.NewErrorResponse("Failed to fetch monthly withdraw amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdraws              = errors.NewErrorResponse("Failed to fetch yearly withdraw amounts", http.StatusInternalServerError)
	ErrFailedFindMonthlyWithdrawsByCardNumber = errors.NewErrorResponse("Failed to fetch monthly withdraw amounts by card", http.StatusInternalServerError)
	ErrFailedFindYearlyWithdrawsByCardNumber  = errors.NewErrorResponse("Failed to fetch yearly withdraw amounts by card", http.StatusInternalServerError)

	ErrFailedFindActiveWithdraws  = errors.NewErrorResponse("Failed to fetch active withdraws", http.StatusInternalServerError)
	ErrFailedFindTrashedWithdraws = errors.NewErrorResponse("Failed to fetch trashed withdraws", http.StatusInternalServerError)

	ErrFailedCreateWithdraw = errors.NewErrorResponse("Failed to create withdraw", http.StatusInternalServerError)
	ErrFailedUpdateWithdraw = errors.NewErrorResponse("Failed to update withdraw", http.StatusInternalServerError)

	ErrFailedTrashedWithdraw            = errors.NewErrorResponse("Failed to trash withdraw", http.StatusInternalServerError)
	ErrFailedRestoreWithdraw            = errors.NewErrorResponse("Failed to restore withdraw", http.StatusInternalServerError)
	ErrFailedDeleteWithdrawPermanent    = errors.NewErrorResponse("Failed to permanently delete withdraw", http.StatusInternalServerError)
	ErrFailedRestoreAllWithdraw         = errors.NewErrorResponse("Failed to restore all withdraws", http.StatusInternalServerError)
	ErrFailedDeleteAllWithdrawPermanent = errors.NewErrorResponse("Failed to permanently delete all withdraws", http.StatusInternalServerError)
)
