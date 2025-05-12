package topup_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrTopupNotFoundRes                = response.NewErrorResponse("Topup not found", http.StatusNotFound)
	ErrFailedFindAllTopups             = response.NewErrorResponse("Failed to fetch Topups", http.StatusInternalServerError)
	ErrFailedFindAllTopupsByCardNumber = response.NewErrorResponse("Failed to fetch Topups by card number", http.StatusInternalServerError)
	ErrFailedFindTopupById             = response.NewErrorResponse("Failed to find Topup by ID", http.StatusInternalServerError)
	ErrFailedFindActiveTopups          = response.NewErrorResponse("Failed to fetch active Topups", http.StatusInternalServerError)
	ErrFailedFindTrashedTopups         = response.NewErrorResponse("Failed to fetch trashed Topups", http.StatusInternalServerError)

	ErrFailedFindMonthTopupStatusSuccess        = response.NewErrorResponse("Failed to get monthly topup success status", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusSuccess       = response.NewErrorResponse("Failed to get yearly topup success status", http.StatusInternalServerError)
	ErrFailedFindMonthTopupStatusFailed         = response.NewErrorResponse("Failed to get monthly topup failed status", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusFailed        = response.NewErrorResponse("Failed to get yearly topup failed status", http.StatusInternalServerError)
	ErrFailedFindMonthTopupStatusSuccessByCard  = response.NewErrorResponse("Failed to get monthly topup success status by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusSuccessByCard = response.NewErrorResponse("Failed to get yearly topup success status by card", http.StatusInternalServerError)
	ErrFailedFindMonthTopupStatusFailedByCard   = response.NewErrorResponse("Failed to get monthly topup failed status by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupStatusFailedByCard  = response.NewErrorResponse("Failed to get yearly topup failed status by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyTopupMethods       = response.NewErrorResponse("Failed to get monthly topup methods", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupMethods        = response.NewErrorResponse("Failed to get yearly topup methods", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupMethodsByCard = response.NewErrorResponse("Failed to get monthly topup methods by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupMethodsByCard  = response.NewErrorResponse("Failed to get yearly topup methods by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyTopupAmounts       = response.NewErrorResponse("Failed to get monthly topup amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmounts        = response.NewErrorResponse("Failed to get yearly topup amounts", http.StatusInternalServerError)
	ErrFailedFindMonthlyTopupAmountsByCard = response.NewErrorResponse("Failed to get monthly topup amounts by card", http.StatusInternalServerError)
	ErrFailedFindYearlyTopupAmountsByCard  = response.NewErrorResponse("Failed to get yearly topup amounts by card", http.StatusInternalServerError)

	ErrFailedCreateTopup = response.NewErrorResponse("Failed to create Topup", http.StatusInternalServerError)
	ErrFailedUpdateTopup = response.NewErrorResponse("Failed to update Topup", http.StatusInternalServerError)

	ErrFailedTrashTopup   = response.NewErrorResponse("Failed to trash Topup", http.StatusInternalServerError)
	ErrFailedRestoreTopup = response.NewErrorResponse("Failed to restore Topup", http.StatusInternalServerError)
	ErrFailedDeleteTopup  = response.NewErrorResponse("Failed to delete Topup permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllTopups = response.NewErrorResponse("Failed to restore all Topups", http.StatusInternalServerError)
	ErrFailedDeleteAllTopups  = response.NewErrorResponse("Failed to delete all Topups permanently", http.StatusInternalServerError)
)
