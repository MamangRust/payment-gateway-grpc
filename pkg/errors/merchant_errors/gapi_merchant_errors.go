package merchant_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcMerchantNotFound      = response.NewGrpcError("merchant", "Merchant not found", int(codes.NotFound))
	ErrGrpcMerchantInvalidID     = response.NewGrpcError("merchant", "Invalid Merchant ID", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidUserID = response.NewGrpcError("merchant", "Invalid Merchant User ID", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidApiKey = response.NewGrpcError("merchant", "Invalid Merchant Api Key", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidMonth  = response.NewGrpcError("month", "Invalid Merchant Month", int(codes.InvalidArgument))
	ErrGrpcMerchantInvalidYear   = response.NewGrpcError("year", "Invalid Merchant Year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllMerchant       = response.NewGrpcError("merchant", "Failed to fetch all merchants", int(codes.Internal))
	ErrGrpcFailedFindByIdMerchant      = response.NewGrpcError("merchant", "Failed to fetch merchant by ID", int(codes.Internal))
	ErrGrpcFailedFindByApiKey          = response.NewGrpcError("merchant", "Failed to fetch merchant by API key", int(codes.Internal))
	ErrGrpcFailedFindByMerchantUserId  = response.NewGrpcError("merchant", "Failed to fetch merchant by user ID", int(codes.Internal))
	ErrGrpcFailedFindByActiveMerchant  = response.NewGrpcError("merchant", "Failed to fetch active merchants", int(codes.Internal))
	ErrGrpcFailedFindByTrashedMerchant = response.NewGrpcError("merchant", "Failed to fetch trashed merchants", int(codes.Internal))

	ErrGrpcFailedFindAllTransactionMerchant   = response.NewGrpcError("merchant", "Failed to fetch all merchant transactions", int(codes.Internal))
	ErrGrpcFailedFindAllTransactionByMerchant = response.NewGrpcError("merchant", "Failed to fetch transactions by merchant", int(codes.Internal))
	ErrGrpcFailedFindAllTransactionByApikey   = response.NewGrpcError("merchant", "Failed to fetch transactions by API key", int(codes.Internal))

	ErrGrpcFailedFindMonthlyPaymentMethodMerchant     = response.NewGrpcError("merchant", "Failed to fetch monthly payment methods", int(codes.Internal))
	ErrGrpcFailedFindYearlyPaymentMethodMerchant      = response.NewGrpcError("merchant", "Failed to fetch yearly payment methods", int(codes.Internal))
	ErrGrpcFailedFindMonthlyPaymentMethodByMerchantId = response.NewGrpcError("merchant", "Failed to fetch monthly payment methods by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindYearlyPaymentMethodByMerchantId  = response.NewGrpcError("merchant", "Failed to fetch yearly payment methods by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindMonthlyPaymentMethodByApikey     = response.NewGrpcError("merchant", "Failed to fetch monthly payment methods by API key", int(codes.Internal))
	ErrGrpcFailedFindYearlyPaymentMethodByApikey      = response.NewGrpcError("merchant", "Failed to fetch yearly payment methods by API key", int(codes.Internal))

	ErrGrpcFailedFindMonthlyAmountMerchant     = response.NewGrpcError("merchant", "Failed to fetch monthly amount", int(codes.Internal))
	ErrGrpcFailedFindYearlyAmountMerchant      = response.NewGrpcError("merchant", "Failed to fetch yearly amount", int(codes.Internal))
	ErrGrpcFailedFindMonthlyAmountByMerchantId = response.NewGrpcError("merchant", "Failed to fetch monthly amount by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindYearlyAmountByMerchantId  = response.NewGrpcError("merchant", "Failed to fetch yearly amount by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindMonthlyAmountByApikey     = response.NewGrpcError("merchant", "Failed to fetch monthly amount by API key", int(codes.Internal))
	ErrGrpcFailedFindYearlyAmountByApikey      = response.NewGrpcError("merchant", "Failed to fetch yearly amount by API key", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTotalAmountMerchant     = response.NewGrpcError("merchant", "Failed to fetch monthly total amount", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalAmountMerchant      = response.NewGrpcError("merchant", "Failed to fetch yearly total amount", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTotalAmountByMerchantId = response.NewGrpcError("merchant", "Failed to fetch monthly total amount by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalAmountByMerchantId  = response.NewGrpcError("merchant", "Failed to fetch yearly total amount by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTotalAmountByApikey     = response.NewGrpcError("merchant", "Failed to fetch monthly total amount by API key", int(codes.Internal))
	ErrGrpcFailedFindYearlyTotalAmountByApikey      = response.NewGrpcError("merchant", "Failed to fetch yearly total amount by API key", int(codes.Internal))

	ErrGrpcFailedCreateMerchant   = response.NewGrpcError("merchant", "Failed to create merchant", int(codes.Internal))
	ErrGrpcFailedUpdateMerchant   = response.NewGrpcError("merchant", "Failed to update merchant", int(codes.Internal))
	ErrGrpcValidateCreateMerchant = response.NewGrpcError("merchant", "Invalid input for create merchant", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateMerchant = response.NewGrpcError("merchant", "Invalid input for update merchant", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedMerchant         = response.NewGrpcError("merchant", "Failed to move merchant to trash", int(codes.Internal))
	ErrGrpcFailedRestoreMerchant         = response.NewGrpcError("merchant", "Failed to restore merchant", int(codes.Internal))
	ErrGrpcFailedDeleteMerchantPermanent = response.NewGrpcError("merchant", "Failed to permanently delete merchant", int(codes.Internal))

	ErrGrpcFailedRestoreAllMerchant         = response.NewGrpcError("merchant", "Failed to restore all merchants", int(codes.Internal))
	ErrGrpcFailedDeleteAllMerchantPermanent = response.NewGrpcError("merchant", "Failed to permanently delete all merchants", int(codes.Internal))
)
