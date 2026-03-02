package merchant_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrMerchantNotFoundRes        = errors.NewErrorResponse("Merchant not found", http.StatusNotFound)
	ErrFailedFindAllMerchants     = errors.NewErrorResponse("Failed to fetch Merchants", http.StatusInternalServerError)
	ErrFailedFindActiveMerchants  = errors.NewErrorResponse("Failed to fetch active Merchants", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchants = errors.NewErrorResponse("Failed to fetch trashed Merchants", http.StatusInternalServerError)
	ErrFailedFindMerchantById     = errors.NewErrorResponse("Failed to find Merchant by ID", http.StatusInternalServerError)
	ErrFailedFindByApiKey         = errors.NewErrorResponse("Failed to find Merchant by API key", http.StatusInternalServerError)
	ErrFailedFindByMerchantUserId = errors.NewErrorResponse("Failed to find Merchant by User ID", http.StatusInternalServerError)

	ErrFailedFindAllTransactions           = errors.NewErrorResponse("Failed to fetch Merchant transactions", http.StatusInternalServerError)
	ErrFailedFindAllTransactionsByMerchant = errors.NewErrorResponse("Failed to fetch transactions by Merchant", http.StatusInternalServerError)
	ErrFailedFindAllTransactionsByApikey   = errors.NewErrorResponse("Failed to fetch transactions by API key", http.StatusInternalServerError)

	ErrFailedFindMonthlyPaymentMethodsMerchant   = errors.NewErrorResponse("Failed to get monthly payment methods", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodMerchant     = errors.NewErrorResponse("Failed to get yearly payment method", http.StatusInternalServerError)
	ErrFailedFindMonthlyPaymentMethodByMerchants = errors.NewErrorResponse("Failed to get monthly payment methods by Merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodByMerchants  = errors.NewErrorResponse("Failed to get yearly payment method by Merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyPaymentMethodByApikeys   = errors.NewErrorResponse("Failed to get monthly payment methods by API key", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodByApikeys    = errors.NewErrorResponse("Failed to get yearly payment method by API key", http.StatusInternalServerError)

	ErrFailedFindMonthlyAmountMerchant    = errors.NewErrorResponse("Failed to get monthly amount", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountMerchant     = errors.NewErrorResponse("Failed to get yearly amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountByMerchants = errors.NewErrorResponse("Failed to get monthly amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountByMerchants  = errors.NewErrorResponse("Failed to get yearly amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountByApikeys   = errors.NewErrorResponse("Failed to get monthly amount by API key", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountByApikeys    = errors.NewErrorResponse("Failed to get yearly amount by API key", http.StatusInternalServerError)

	ErrFailedFindMonthlyTotalAmountMerchant    = errors.NewErrorResponse("Failed to get monthly total amount", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalAmountMerchant     = errors.NewErrorResponse("Failed to get yearly total amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalAmountByMerchants = errors.NewErrorResponse("Failed to get monthly total amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalAmountByMerchants  = errors.NewErrorResponse("Failed to get yearly total amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalAmountByApikeys   = errors.NewErrorResponse("Failed to get monthly total amount by API key", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalAmountByApikeys    = errors.NewErrorResponse("Failed to get yearly total amount by API key", http.StatusInternalServerError)

	ErrFailedCreateMerchant = errors.NewErrorResponse("Failed to create Merchant", http.StatusInternalServerError)
	ErrFailedUpdateMerchant = errors.NewErrorResponse("Failed to update Merchant", http.StatusInternalServerError)

	ErrFailedTrashMerchant   = errors.NewErrorResponse("Failed to trash Merchant", http.StatusInternalServerError)
	ErrFailedRestoreMerchant = errors.NewErrorResponse("Failed to restore Merchant", http.StatusInternalServerError)
	ErrFailedDeleteMerchant  = errors.NewErrorResponse("Failed to delete Merchant permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllMerchants = errors.NewErrorResponse("Failed to restore all Merchants", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchants  = errors.NewErrorResponse("Failed to delete all Merchants permanently", http.StatusInternalServerError)
)
