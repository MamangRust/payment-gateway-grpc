package topup_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrTopupNotFoundRes                = errors.NewErrorResponse("Topup not found", http.StatusNotFound)
	ErrFailedFindAllTopups             = errors.NewErrorResponse("Failed to fetch Topups", http.StatusInternalServerError)
	ErrFailedFindAllTopupsByCardNumber = errors.NewErrorResponse("Failed to fetch Topups by card number", http.StatusInternalServerError)
	ErrFailedFindTopupById             = errors.NewErrorResponse("Failed to find Topup by ID", http.StatusInternalServerError)
	ErrFailedFindActiveTopups          = errors.NewErrorResponse("Failed to fetch active Topups", http.StatusInternalServerError)
	ErrFailedFindTrashedTopups         = errors.NewErrorResponse("Failed to fetch trashed Topups", http.StatusInternalServerError)

	ErrFailedFindMonthTopupStatusSuccess        = errors.NewErrorResponse("Failed to get monthly topup success status", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusSuccess       = errors.NewErrorResponse("Failed to get yearly topup success status", http.StatusInternalServerError)
	ErrFailedFindMonthTopupStatusFailed         = errors.NewErrorResponse("Failed to get monthly topup failed status", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusFailed        = errors.NewErrorResponse("Failed to get yearly topup failed status", http.StatusInternalServerError)
	ErrFailedFindMonthTopupStatusSuccessByCard  = errors.NewErrorResponse("Failed to get monthly topup success status by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusSuccessByCard = errors.NewErrorResponse("Failed to get yearly topup success status by card", http.StatusInternalServerError)
	ErrFailedFindMonthTopupStatusFailedByCard   = errors.NewErrorResponse("Failed to get monthly topup failed status by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusFailedByCard  = errors.NewErrorResponse("Failed to get yearly topup failed status by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyTopupMethods       = errors.NewErrorResponse("Failed to get monthly topup methods", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupMethods        = errors.NewErrorResponse("Failed to get yearly topup methods", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupMethodsByCard = errors.NewErrorResponse("Failed to get monthly topup methods by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupMethodsByCard  = errors.NewErrorResponse("Failed to get yearly topup methods by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyTopupAmounts       = errors.NewErrorResponse("Failed to get monthly topup amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmounts        = errors.NewErrorResponse("Failed to get yearly topup amounts", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupAmountsByCard = errors.NewErrorResponse("Failed to get monthly topup amounts by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmountsByCard  = errors.NewErrorResponse("Failed to get yearly topup amounts by card", http.StatusInternalServerError)

	ErrFailedCreateTopup = errors.NewErrorResponse("Failed to create Topup", http.StatusInternalServerError)
	ErrFailedUpdateTopup = errors.NewErrorResponse("Failed to update Topup", http.StatusInternalServerError)

	ErrFailedTrashTopup   = errors.NewErrorResponse("Failed to trash Topup", http.StatusInternalServerError)
	ErrFailedRestoreTopup = errors.NewErrorResponse("Failed to restore Topup", http.StatusInternalServerError)
	ErrFailedDeleteTopup  = errors.NewErrorResponse("Failed to delete Topup permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllTopups = errors.NewErrorResponse("Failed to restore all Topups", http.StatusInternalServerError)
	ErrFailedDeleteAllTopups  = errors.NewErrorResponse("Failed to delete all Topups permanently", http.StatusInternalServerError)
)
