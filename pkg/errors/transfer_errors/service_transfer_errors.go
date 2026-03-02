package transfer_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllTransfers = errors.NewErrorResponse("Failed to fetch all transfers", http.StatusInternalServerError)
	ErrTransferNotFound       = errors.NewErrorResponse("Transfer not found", http.StatusNotFound)

	ErrFailedFindMonthTransferStatusSuccess = errors.NewErrorResponse("Failed to fetch monthly successful transfers", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusSuccess  = errors.NewErrorResponse("Failed to fetch yearly successful transfers", http.StatusInternalServerError)
	ErrFailedFindMonthTransferStatusFailed  = errors.NewErrorResponse("Failed to fetch monthly failed transfers", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusFailed   = errors.NewErrorResponse("Failed to fetch yearly failed transfers", http.StatusInternalServerError)

	ErrFailedFindMonthTransferStatusSuccessByCard = errors.NewErrorResponse("Failed to fetch monthly successful transfers by card", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusSuccessByCard  = errors.NewErrorResponse("Failed to fetch yearly successful transfers by card", http.StatusInternalServerError)
	ErrFailedFindMonthTransferStatusFailedByCard  = errors.NewErrorResponse("Failed to fetch monthly failed transfers by card", http.StatusInternalServerError)
	ErrFailedFindYearTransferStatusFailedByCard   = errors.NewErrorResponse("Failed to fetch yearly failed transfers by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyTransferAmounts               = errors.NewErrorResponse("Failed to fetch monthly transfer amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmounts                = errors.NewErrorResponse("Failed to fetch yearly transfer amounts", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountsBySenderCard   = errors.NewErrorResponse("Failed to fetch monthly transfer amounts by sender card", http.StatusInternalServerError)
	ErrFailedFindMonthlyTransferAmountsByReceiverCard = errors.NewErrorResponse("Failed to fetch monthly transfer amounts by receiver card", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountsBySenderCard    = errors.NewErrorResponse("Failed to fetch yearly transfer amounts by sender card", http.StatusInternalServerError)
	ErrFailedFindYearlyTransferAmountsByReceiverCard  = errors.NewErrorResponse("Failed to fetch yearly transfer amounts by receiver card", http.StatusInternalServerError)

	ErrFailedFindActiveTransfers  = errors.NewErrorResponse("Failed to fetch active transfers", http.StatusInternalServerError)
	ErrFailedFindTrashedTransfers = errors.NewErrorResponse("Failed to fetch trashed transfers", http.StatusInternalServerError)

	ErrFailedFindTransfersBySender   = errors.NewErrorResponse("Failed to fetch transfers by sender", http.StatusInternalServerError)
	ErrFailedFindTransfersByReceiver = errors.NewErrorResponse("Failed to fetch transfers by receiver", http.StatusInternalServerError)

	ErrFailedCreateTransfer = errors.NewErrorResponse("Failed to create transfer", http.StatusInternalServerError)
	ErrFailedUpdateTransfer = errors.NewErrorResponse("Failed to update transfer", http.StatusInternalServerError)

	ErrFailedTrashedTransfer             = errors.NewErrorResponse("Failed to trash transfer", http.StatusInternalServerError)
	ErrFailedRestoreTransfer             = errors.NewErrorResponse("Failed to restore transfer", http.StatusInternalServerError)
	ErrFailedDeleteTransferPermanent     = errors.NewErrorResponse("Failed to permanently delete transfer", http.StatusInternalServerError)
	ErrFailedRestoreAllTransfers         = errors.NewErrorResponse("Failed to restore all transfers", http.StatusInternalServerError)
	ErrFailedDeleteAllTransfersPermanent = errors.NewErrorResponse("Failed to permanently delete all transfers", http.StatusInternalServerError)
)
