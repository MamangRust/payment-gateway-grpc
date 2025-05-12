package transaction_errors

import "errors"

var (
	ErrFindAllTransactionsFailed          = errors.New("failed to find all transactions")
	ErrFindActiveTransactionsFailed       = errors.New("failed to find active transactions")
	ErrFindTrashedTransactionsFailed      = errors.New("failed to find trashed transactions")
	ErrFindTransactionsByCardNumberFailed = errors.New("failed to find transactions by card number")
	ErrFindTransactionByIdFailed          = errors.New("failed to find transaction by ID")

	ErrGetMonthTransactionStatusSuccessFailed        = errors.New("failed to get monthly transaction status success")
	ErrGetYearlyTransactionStatusSuccessFailed       = errors.New("failed to get yearly transaction status success")
	ErrGetMonthTransactionStatusSuccessByCardFailed  = errors.New("failed to get monthly transaction status success by card number")
	ErrGetYearlyTransactionStatusSuccessByCardFailed = errors.New("failed to get yearly transaction status success by card number")

	ErrGetMonthTransactionStatusFailedFailed        = errors.New("failed to get monthly transaction status failed")
	ErrGetYearlyTransactionStatusFailedFailed       = errors.New("failed to get yearly transaction status failed")
	ErrGetMonthTransactionStatusFailedByCardFailed  = errors.New("failed to get monthly transaction status failed by card number")
	ErrGetYearlyTransactionStatusFailedByCardFailed = errors.New("failed to get yearly transaction status failed by card number")

	ErrGetMonthlyPaymentMethodsFailed       = errors.New("failed to get monthly payment methods")
	ErrGetYearlyPaymentMethodsFailed        = errors.New("failed to get yearly payment methods")
	ErrGetMonthlyAmountsFailed              = errors.New("failed to get monthly amounts")
	ErrGetYearlyAmountsFailed               = errors.New("failed to get yearly amounts")
	ErrGetMonthlyPaymentMethodsByCardFailed = errors.New("failed to get monthly payment methods by card number")
	ErrGetYearlyPaymentMethodsByCardFailed  = errors.New("failed to get yearly payment methods by card number")
	ErrGetMonthlyAmountsByCardFailed        = errors.New("failed to get monthly amounts by card number")
	ErrGetYearlyAmountsByCardFailed         = errors.New("failed to get yearly amounts by card number")

	ErrFindTransactionByMerchantIdFailed = errors.New("failed to find transaction by merchant ID")

	ErrCreateTransactionFailed       = errors.New("failed to create transaction")
	ErrUpdateTransactionFailed       = errors.New("failed to update transaction")
	ErrUpdateTransactionStatusFailed = errors.New("failed to update transaction status")

	ErrTrashedTransactionFailed             = errors.New("failed to soft-delete (trash) transaction")
	ErrRestoreTransactionFailed             = errors.New("failed to restore transaction")
	ErrDeleteTransactionPermanentFailed     = errors.New("failed to permanently delete transaction")
	ErrRestoreAllTransactionsFailed         = errors.New("failed to restore all transactions")
	ErrDeleteAllTransactionsPermanentFailed = errors.New("failed to permanently delete all transactions")
)
