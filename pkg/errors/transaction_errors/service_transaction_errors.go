package transaction_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrFailedFindAllTransactions = errors.NewErrorResponse("Failed to fetch all transactions", http.StatusInternalServerError)
	ErrFailedFindAllByCardNumber = errors.NewErrorResponse("Failed to fetch transactions by card number", http.StatusInternalServerError)
	ErrTransactionNotFound       = errors.NewErrorResponse("Transaction not found", http.StatusNotFound)

	ErrFailedFindMonthTransactionSuccess = errors.NewErrorResponse("Failed to fetch monthly successful transactions", http.StatusInternalServerError)
	ErrFailedFindYearTransactionSuccess  = errors.NewErrorResponse("Failed to fetch yearly successful transactions", http.StatusInternalServerError)
	ErrFailedFindMonthTransactionFailed  = errors.NewErrorResponse("Failed to fetch monthly failed transactions", http.StatusInternalServerError)
	ErrFailedFindYearTransactionFailed   = errors.NewErrorResponse("Failed to fetch yearly failed transactions", http.StatusInternalServerError)

	ErrFailedFindMonthTransactionSuccessByCard = errors.NewErrorResponse("Failed to fetch monthly successful transactions by card", http.StatusInternalServerError)
	ErrFailedFindYearTransactionSuccessByCard  = errors.NewErrorResponse("Failed to fetch yearly successful transactions by card", http.StatusInternalServerError)
	ErrFailedFindMonthTransactionFailedByCard  = errors.NewErrorResponse("Failed to fetch monthly failed transactions by card", http.StatusInternalServerError)
	ErrFailedFindYearTransactionFailedByCard   = errors.NewErrorResponse("Failed to fetch yearly failed transactions by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyPaymentMethods = errors.NewErrorResponse("Failed to fetch monthly payment methods", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethods  = errors.NewErrorResponse("Failed to fetch yearly payment methods", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmounts        = errors.NewErrorResponse("Failed to fetch monthly amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyAmounts         = errors.NewErrorResponse("Failed to fetch yearly amounts", http.StatusInternalServerError)

	ErrFailedFindMonthlyPaymentMethodsByCard = errors.NewErrorResponse("Failed to fetch monthly payment methods by card", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodsByCard  = errors.NewErrorResponse("Failed to fetch yearly payment methods by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountsByCard        = errors.NewErrorResponse("Failed to fetch monthly amounts by card", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountsByCard         = errors.NewErrorResponse("Failed to fetch yearly amounts by card", http.StatusInternalServerError)

	ErrFailedFindByActiveTransactions  = errors.NewErrorResponse("Failed to fetch active transactions", http.StatusInternalServerError)
	ErrFailedFindByTrashedTransactions = errors.NewErrorResponse("Failed to fetch trashed transactions", http.StatusInternalServerError)
	ErrFailedFindByMerchantID          = errors.NewErrorResponse("Failed to fetch transactions by merchant ID", http.StatusInternalServerError)

	ErrFailedCreateTransaction = errors.NewErrorResponse("Failed to create transaction", http.StatusInternalServerError)
	ErrFailedUpdateTransaction = errors.NewErrorResponse("Failed to update transaction", http.StatusInternalServerError)

	ErrFailedTrashedTransaction         = errors.NewErrorResponse("Failed to trash transaction", http.StatusInternalServerError)
	ErrFailedRestoreTransaction         = errors.NewErrorResponse("Failed to restore transaction", http.StatusInternalServerError)
	ErrFailedDeleteTransactionPermanent = errors.NewErrorResponse("Failed to permanently delete transaction", http.StatusInternalServerError)

	ErrFailedRestoreAllTransactions         = errors.NewErrorResponse("Failed to restore all transactions", http.StatusInternalServerError)
	ErrFailedDeleteAllTransactionsPermanent = errors.NewErrorResponse("Failed to permanently delete all transactions", http.StatusInternalServerError)
)
