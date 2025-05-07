package transfer_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllTransfers = response.NewErrorResponse("Failed to fetch all transfers", http.StatusInternalServerError)
	ErrTransferNotFound       = response.NewErrorResponse("Transfer not found", http.StatusNotFound)

	ErrFailedFindMonthTransferStatusSuccess = response.NewErrorResponse("Failed to fetch monthly successful transfers", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusSuccess  = response.NewErrorResponse("Failed to fetch yearly successful transfers", http.StatusInternalServerError)
	ErrFailedFindMonthTransferStatusFailed  = response.NewErrorResponse("Failed to fetch monthly failed transfers", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusFailed   = response.NewErrorResponse("Failed to fetch yearly failed transfers", http.StatusInternalServerError)

	ErrFailedFindMonthTransferStatusSuccessByCard = response.NewErrorResponse("Failed to fetch monthly successful transfers by card", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusSuccessByCard  = response.NewErrorResponse("Failed to fetch yearly successful transfers by card", http.StatusInternalServerError)
	ErrFailedFindMonthTransferStatusFailedByCard  = response.NewErrorResponse("Failed to fetch monthly failed transfers by card", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusFailedByCard   = response.NewErrorResponse("Failed to fetch yearly failed transfers by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyTransferAmounts               = response.NewErrorResponse("Failed to fetch monthly transfer amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmounts                = response.NewErrorResponse("Failed to fetch yearly transfer amounts", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountsBySenderCard   = response.NewErrorResponse("Failed to fetch monthly transfer amounts by sender card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountsByReceiverCard = response.NewErrorResponse("Failed to fetch monthly transfer amounts by receiver card", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountsBySenderCard    = response.NewErrorResponse("Failed to fetch yearly transfer amounts by sender card", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountsByReceiverCard  = response.NewErrorResponse("Failed to fetch yearly transfer amounts by receiver card", http.StatusInternalServerError)

	ErrFailedFindActiveTransfers  = response.NewErrorResponse("Failed to fetch active transfers", http.StatusInternalServerError)
	ErrFailedFindTrashedTransfers = response.NewErrorResponse("Failed to fetch trashed transfers", http.StatusInternalServerError)

	ErrFailedFindTransfersBySender   = response.NewErrorResponse("Failed to fetch transfers by sender", http.StatusInternalServerError)
	ErrFailedFindTransfersByReceiver = response.NewErrorResponse("Failed to fetch transfers by receiver", http.StatusInternalServerError)

	ErrFailedCreateTransfer = response.NewErrorResponse("Failed to create transfer", http.StatusInternalServerError)
	ErrFailedUpdateTransfer = response.NewErrorResponse("Failed to update transfer", http.StatusInternalServerError)

	ErrFailedTrashedTransfer             = response.NewErrorResponse("Failed to trash transfer", http.StatusInternalServerError)
	ErrFailedRestoreTransfer             = response.NewErrorResponse("Failed to restore transfer", http.StatusInternalServerError)
	ErrFailedDeleteTransferPermanent     = response.NewErrorResponse("Failed to permanently delete transfer", http.StatusInternalServerError)
	ErrFailedRestoreAllTransfers         = response.NewErrorResponse("Failed to restore all transfers", http.StatusInternalServerError)
	ErrFailedDeleteAllTransfersPermanent = response.NewErrorResponse("Failed to permanently delete all transfers", http.StatusInternalServerError)
)
