package merchant_errors

import "errors"

var (
	ErrFindAllMerchantsFailed     = errors.New("failed to find all merchants")
	ErrFindActiveMerchantsFailed  = errors.New("failed to find active merchants")
	ErrFindTrashedMerchantsFailed = errors.New("failed to find trashed merchants")
	ErrFindMerchantByIdFailed     = errors.New("failed to find merchant by ID")
	ErrFindMerchantByApiKeyFailed = errors.New("failed to find merchant by API key")
	ErrFindMerchantByNameFailed   = errors.New("failed to find merchant by name")
	ErrFindMerchantByUserIdFailed = errors.New("failed to find merchant by user ID")

	ErrGetMonthlyTotalAmountMerchantFailed    = errors.New("failed to get monthly total amount of merchant")
	ErrGetYearlyTotalAmountMerchantFailed     = errors.New("failed to get yearly total amount of merchant")
	ErrGetMonthlyTotalAmountByMerchantsFailed = errors.New("failed to get monthly total amount by merchants")
	ErrGetYearlyTotalAmountByMerchantsFailed  = errors.New("failed to get yearly total amount by merchants")
	ErrGetMonthlyTotalAmountByApikeyFailed    = errors.New("failed to get monthly total amount by API key")
	ErrGetYearlyTotalAmountByApikeyFailed     = errors.New("failed to get yearly total amount by API key")

	ErrGetMonthlyAmountMerchantFailed    = errors.New("failed to get monthly amount of merchant")
	ErrGetYearlyAmountMerchantFailed     = errors.New("failed to get yearly amount of merchant")
	ErrGetMonthlyAmountByMerchantsFailed = errors.New("failed to get monthly amount by merchants")
	ErrGetYearlyAmountByMerchantsFailed  = errors.New("failed to get yearly amount by merchants")
	ErrGetMonthlyAmountByApikeyFailed    = errors.New("failed to get monthly amount by API key")
	ErrGetYearlyAmountByApikeyFailed     = errors.New("failed to get yearly amount by API key")

	ErrGetMonthlyPaymentMethodsMerchantFailed   = errors.New("failed to get monthly payment methods of merchant")
	ErrGetYearlyPaymentMethodMerchantFailed     = errors.New("failed to get yearly payment method of merchant")
	ErrGetMonthlyPaymentMethodByMerchantsFailed = errors.New("failed to get monthly payment method by merchants")
	ErrGetYearlyPaymentMethodByMerchantsFailed  = errors.New("failed to get yearly payment method by merchants")
	ErrGetMonthlyPaymentMethodByApikeyFailed    = errors.New("failed to get monthly payment method by API key")
	ErrGetYearlyPaymentMethodByApikeyFailed     = errors.New("failed to get yearly payment method by API key")

	ErrFindAllTransactionsFailed           = errors.New("failed to find all merchant transactions")
	ErrFindAllTransactionsByMerchantFailed = errors.New("failed to find merchant transactions by merchant ID")
	ErrFindAllTransactionsByApiKeyFailed   = errors.New("failed to find merchant transactions by API key")

	ErrCreateMerchantFailed       = errors.New("failed to create merchant")
	ErrUpdateMerchantFailed       = errors.New("failed to update merchant")
	ErrUpdateMerchantStatusFailed = errors.New("failed to update merchant status")

	ErrTrashedMerchantFailed            = errors.New("failed to soft-delete (trash) merchant")
	ErrRestoreMerchantFailed            = errors.New("failed to restore merchant")
	ErrDeleteMerchantPermanentFailed    = errors.New("failed to permanently delete merchant")
	ErrRestoreAllMerchantFailed         = errors.New("failed to restore all merchants")
	ErrDeleteAllMerchantPermanentFailed = errors.New("failed to permanently delete all merchants")
)
