package transaction_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrFailedFindAllTransactions = response.NewErrorResponse("Failed to fetch all transactions", http.StatusInternalServerError)
	ErrFailedFindAllByCardNumber = response.NewErrorResponse("Failed to fetch transactions by card number", http.StatusInternalServerError)
	ErrTransactionNotFound       = response.NewErrorResponse("Transaction not found", http.StatusNotFound)

	ErrFailedFindMonthTransactionSuccess = response.NewErrorResponse("Failed to fetch monthly successful transactions", http.StatusInternalServerError)
	ErrFailedFindYearTransactionSuccess  = response.NewErrorResponse("Failed to fetch yearly successful transactions", http.StatusInternalServerError)
	ErrFailedFindMonthTransactionFailed  = response.NewErrorResponse("Failed to fetch monthly failed transactions", http.StatusInternalServerError)
	ErrFailedFindYearTransactionFailed   = response.NewErrorResponse("Failed to fetch yearly failed transactions", http.StatusInternalServerError)

	ErrFailedFindMonthTransactionSuccessByCard = response.NewErrorResponse("Failed to fetch monthly successful transactions by card", http.StatusInternalServerError)
	ErrFailedFindYearTransactionSuccessByCard  = response.NewErrorResponse("Failed to fetch yearly successful transactions by card", http.StatusInternalServerError)
	ErrFailedFindMonthTransactionFailedByCard  = response.NewErrorResponse("Failed to fetch monthly failed transactions by card", http.StatusInternalServerError)
	ErrFailedFindYearTransactionFailedByCard   = response.NewErrorResponse("Failed to fetch yearly failed transactions by card", http.StatusInternalServerError)

	ErrFailedFindMonthlyPaymentMethods = response.NewErrorResponse("Failed to fetch monthly payment methods", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethods  = response.NewErrorResponse("Failed to fetch yearly payment methods", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmounts        = response.NewErrorResponse("Failed to fetch monthly amounts", http.StatusInternalServerError)
	ErrFailedFindYearlyAmounts         = response.NewErrorResponse("Failed to fetch yearly amounts", http.StatusInternalServerError)

	ErrFailedFindMonthlyPaymentMethodsByCard = response.NewErrorResponse("Failed to fetch monthly payment methods by card", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodsByCard  = response.NewErrorResponse("Failed to fetch yearly payment methods by card", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountsByCard        = response.NewErrorResponse("Failed to fetch monthly amounts by card", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountsByCard         = response.NewErrorResponse("Failed to fetch yearly amounts by card", http.StatusInternalServerError)

	ErrFailedFindByActiveTransactions  = response.NewErrorResponse("Failed to fetch active transactions", http.StatusInternalServerError)
	ErrFailedFindByTrashedTransactions = response.NewErrorResponse("Failed to fetch trashed transactions", http.StatusInternalServerError)
	ErrFailedFindByMerchantID          = response.NewErrorResponse("Failed to fetch transactions by merchant ID", http.StatusInternalServerError)

	ErrFailedCreateTransaction = response.NewErrorResponse("Failed to create transaction", http.StatusInternalServerError)
	ErrFailedUpdateTransaction = response.NewErrorResponse("Failed to update transaction", http.StatusInternalServerError)

	ErrFailedTrashedTransaction         = response.NewErrorResponse("Failed to trash transaction", http.StatusInternalServerError)
	ErrFailedRestoreTransaction         = response.NewErrorResponse("Failed to restore transaction", http.StatusInternalServerError)
	ErrFailedDeleteTransactionPermanent = response.NewErrorResponse("Failed to permanently delete transaction", http.StatusInternalServerError)

	ErrFailedRestoreAllTransactions         = response.NewErrorResponse("Failed to restore all transactions", http.StatusInternalServerError)
	ErrFailedDeleteAllTransactionsPermanent = response.NewErrorResponse("Failed to permanently delete all transactions", http.StatusInternalServerError)
)
