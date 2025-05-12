package transfer_errors

import "errors"

var (
	ErrFindAllTransfersFailed     = errors.New("failed to find all transfers")
	ErrFindActiveTransfersFailed  = errors.New("failed to find active transfers")
	ErrFindTrashedTransfersFailed = errors.New("failed to find trashed transfers")
	ErrFindTransferByIdFailed     = errors.New("failed to find transfer by ID")

	ErrGetMonthTransferStatusSuccessFailed        = errors.New("failed to get monthly transfer status success")
	ErrGetYearlyTransferStatusSuccessFailed       = errors.New("failed to get yearly transfer status success")
	ErrGetMonthTransferStatusSuccessByCardFailed  = errors.New("failed to get monthly transfer status success by card number")
	ErrGetYearlyTransferStatusSuccessByCardFailed = errors.New("failed to get yearly transfer status success by card number")

	ErrGetMonthTransferStatusFailedFailed        = errors.New("failed to get monthly transfer status failed")
	ErrGetYearlyTransferStatusFailedFailed       = errors.New("failed to get yearly transfer status failed")
	ErrGetMonthTransferStatusFailedByCardFailed  = errors.New("failed to get monthly transfer status failed by card number")
	ErrGetYearlyTransferStatusFailedByCardFailed = errors.New("failed to get yearly transfer status failed by card number")

	ErrGetMonthlyTransferAmountsFailed               = errors.New("failed to get monthly transfer amounts")
	ErrGetYearlyTransferAmountsFailed                = errors.New("failed to get yearly transfer amounts")
	ErrGetMonthlyTransferAmountsBySenderCardFailed   = errors.New("failed to get monthly transfer amounts by sender card number")
	ErrGetYearlyTransferAmountsBySenderCardFailed    = errors.New("failed to get yearly transfer amounts by sender card number")
	ErrGetMonthlyTransferAmountsByReceiverCardFailed = errors.New("failed to get monthly transfer amounts by receiver card number")
	ErrGetYearlyTransferAmountsByReceiverCardFailed  = errors.New("failed to get yearly transfer amounts by receiver card number")

	ErrFindTransferByTransferFromFailed = errors.New("failed to find transfer by transfer from")
	ErrFindTransferByTransferToFailed   = errors.New("failed to find transfer by transfer to")

	ErrCreateTransferFailed       = errors.New("failed to create transfer")
	ErrUpdateTransferFailed       = errors.New("failed to update transfer")
	ErrUpdateTransferAmountFailed = errors.New("failed to update transfer amount")
	ErrUpdateTransferStatusFailed = errors.New("failed to update transfer status")

	ErrTrashedTransferFailed             = errors.New("failed to soft-delete (trash) transfer")
	ErrRestoreTransferFailed             = errors.New("failed to restore transfer")
	ErrDeleteTransferPermanentFailed     = errors.New("failed to permanently delete transfer")
	ErrRestoreAllTransfersFailed         = errors.New("failed to restore all transfers")
	ErrDeleteAllTransfersPermanentFailed = errors.New("failed to permanently delete all transfers")
)
