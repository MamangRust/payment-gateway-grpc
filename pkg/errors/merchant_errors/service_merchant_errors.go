package merchant_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"net/http"
)

var (
	ErrMerchantNotFoundRes        = response.NewErrorResponse("Merchant not found", http.StatusNotFound)
	ErrFailedFindAllMerchants     = response.NewErrorResponse("Failed to fetch Merchants", http.StatusInternalServerError)
	ErrFailedFindActiveMerchants  = response.NewErrorResponse("Failed to fetch active Merchants", http.StatusInternalServerError)
	ErrFailedFindTrashedMerchants = response.NewErrorResponse("Failed to fetch trashed Merchants", http.StatusInternalServerError)
	ErrFailedFindMerchantById     = response.NewErrorResponse("Failed to find Merchant by ID", http.StatusInternalServerError)
	ErrFailedFindByApiKey         = response.NewErrorResponse("Failed to find Merchant by API key", http.StatusInternalServerError)
	ErrFailedFindByMerchantUserId = response.NewErrorResponse("Failed to find Merchant by User ID", http.StatusInternalServerError)

	ErrFailedFindAllTransactions           = response.NewErrorResponse("Failed to fetch Merchant transactions", http.StatusInternalServerError)
	ErrFailedFindAllTransactionsByMerchant = response.NewErrorResponse("Failed to fetch transactions by Merchant", http.StatusInternalServerError)
	ErrFailedFindAllTransactionsByApikey   = response.NewErrorResponse("Failed to fetch transactions by API key", http.StatusInternalServerError)

	ErrFailedFindMonthlyPaymentMethodsMerchant   = response.NewErrorResponse("Failed to get monthly payment methods", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodMerchant     = response.NewErrorResponse("Failed to get yearly payment method", http.StatusInternalServerError)
	ErrFailedFindMonthlyPaymentMethodByMerchants = response.NewErrorResponse("Failed to get monthly payment methods by Merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodByMerchants  = response.NewErrorResponse("Failed to get yearly payment method by Merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyPaymentMethodByApikeys   = response.NewErrorResponse("Failed to get monthly payment methods by API key", http.StatusInternalServerError)
	ErrFailedFindYearlyPaymentMethodByApikeys    = response.NewErrorResponse("Failed to get yearly payment method by API key", http.StatusInternalServerError)

	ErrFailedFindMonthlyAmountMerchant    = response.NewErrorResponse("Failed to get monthly amount", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountMerchant     = response.NewErrorResponse("Failed to get yearly amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountByMerchants = response.NewErrorResponse("Failed to get monthly amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountByMerchants  = response.NewErrorResponse("Failed to get yearly amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyAmountByApikeys   = response.NewErrorResponse("Failed to get monthly amount by API key", http.StatusInternalServerError)
	ErrFailedFindYearlyAmountByApikeys    = response.NewErrorResponse("Failed to get yearly amount by API key", http.StatusInternalServerError)

	ErrFailedFindMonthlyTotalAmountMerchant    = response.NewErrorResponse("Failed to get monthly total amount", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalAmountMerchant     = response.NewErrorResponse("Failed to get yearly total amount", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalAmountByMerchants = response.NewErrorResponse("Failed to get monthly total amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalAmountByMerchants  = response.NewErrorResponse("Failed to get yearly total amount by Merchant", http.StatusInternalServerError)
	ErrFailedFindMonthlyTotalAmountByApikeys   = response.NewErrorResponse("Failed to get monthly total amount by API key", http.StatusInternalServerError)
	ErrFailedFindYearlyTotalAmountByApikeys    = response.NewErrorResponse("Failed to get yearly total amount by API key", http.StatusInternalServerError)

	ErrFailedCreateMerchant = response.NewErrorResponse("Failed to create Merchant", http.StatusInternalServerError)
	ErrFailedUpdateMerchant = response.NewErrorResponse("Failed to update Merchant", http.StatusInternalServerError)

	ErrFailedTrashMerchant   = response.NewErrorResponse("Failed to trash Merchant", http.StatusInternalServerError)
	ErrFailedRestoreMerchant = response.NewErrorResponse("Failed to restore Merchant", http.StatusInternalServerError)
	ErrFailedDeleteMerchant  = response.NewErrorResponse("Failed to delete Merchant permanently", http.StatusInternalServerError)

	ErrFailedRestoreAllMerchants = response.NewErrorResponse("Failed to restore all Merchants", http.StatusInternalServerError)
	ErrFailedDeleteAllMerchants  = response.NewErrorResponse("Failed to delete all Merchants permanently", http.StatusInternalServerError)
)
