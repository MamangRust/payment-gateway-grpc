package transaction_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTransactionNotFound          = response.NewGrpcError("transaction", "Transaction not found", int(codes.NotFound))
	ErrGrpcTransactionInvalidID         = response.NewGrpcError("transaction", "Invalid Transaction ID", int(codes.InvalidArgument))
	ErrGrpcTransactionInvalidMerchantID = response.NewGrpcError("transaction", "Invalid Transaction Merchant ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber            = response.NewGrpcError("card_id", "Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth                 = response.NewGrpcError("month", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear                  = response.NewGrpcError("year", "Invalid year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllTransaction             = response.NewGrpcError("transaction", "Failed to fetch all transactions", int(codes.Internal))
	ErrGrpcFailedFindAllTransactionByCardNumber = response.NewGrpcError("transaction", "Failed to fetch transactions by card number", int(codes.Internal))
	ErrGrpcFailedFindByIdTransaction            = response.NewGrpcError("transaction", "Failed to fetch transaction by ID", int(codes.Internal))
	ErrGrpcFailedFindTransactionByMerchantId    = response.NewGrpcError("transaction", "Failed to fetch transactions by merchant ID", int(codes.Internal))
	ErrGrpcFailedFindByActiveTransaction        = response.NewGrpcError("transaction", "Failed to fetch active transactions", int(codes.Internal))
	ErrGrpcFailedFindByTrashedTransaction       = response.NewGrpcError("transaction", "Failed to fetch trashed transactions", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTransactionStatusSuccess             = response.NewGrpcError("transaction", "Failed to fetch monthly successful transactions", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransactionStatusSuccess              = response.NewGrpcError("transaction", "Failed to fetch yearly successful transactions", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransactionStatusFailed              = response.NewGrpcError("transaction", "Failed to fetch monthly failed transactions", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransactionStatusFailed               = response.NewGrpcError("transaction", "Failed to fetch yearly failed transactions", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransactionStatusSuccessByCardNumber = response.NewGrpcError("transaction", "Failed to fetch monthly successful transactions by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransactionStatusSuccessByCardNumber  = response.NewGrpcError("transaction", "Failed to fetch yearly successful transactions by card number", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransactionStatusFailedByCardNumber  = response.NewGrpcError("transaction", "Failed to fetch monthly failed transactions by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransactionStatusFailedByCardNumber   = response.NewGrpcError("transaction", "Failed to fetch yearly failed transactions by card number", int(codes.Internal))

	ErrGrpcFailedFindMonthlyPaymentMethods             = response.NewGrpcError("transaction", "Failed to fetch monthly payment methods", int(codes.Internal))
	ErrGrpcFailedFindYearlyPaymentMethods              = response.NewGrpcError("transaction", "Failed to fetch yearly payment methods", int(codes.Internal))
	ErrGrpcFailedFindMonthlyAmounts                    = response.NewGrpcError("transaction", "Failed to fetch monthly transaction amounts", int(codes.Internal))
	ErrGrpcFailedFindYearlyAmounts                     = response.NewGrpcError("transaction", "Failed to fetch yearly transaction amounts", int(codes.Internal))
	ErrGrpcFailedFindMonthlyPaymentMethodsByCardNumber = response.NewGrpcError("transaction", "Failed to fetch monthly payment methods by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyPaymentMethodsByCardNumber  = response.NewGrpcError("transaction", "Failed to fetch yearly payment methods by card number", int(codes.Internal))
	ErrGrpcFailedFindMonthlyAmountsByCardNumber        = response.NewGrpcError("transaction", "Failed to fetch monthly amounts by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyAmountsByCardNumber         = response.NewGrpcError("transaction", "Failed to fetch yearly amounts by card number", int(codes.Internal))

	ErrGrpcFailedCreateTransaction          = response.NewGrpcError("transaction", "Failed to create transaction", int(codes.Internal))
	ErrGrpcFailedUpdateTransaction          = response.NewGrpcError("transaction", "Failed to update transaction", int(codes.Internal))
	ErrGrpcValidateCreateTransactionRequest = response.NewGrpcError("transaction", "Invalid input for create card", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransactionRequest = response.NewGrpcError("transaction", "Invalid input for update card", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedTransaction         = response.NewGrpcError("transaction", "Failed to move transaction to trash", int(codes.Internal))
	ErrGrpcFailedRestoreTransaction         = response.NewGrpcError("transaction", "Failed to restore transaction", int(codes.Internal))
	ErrGrpcFailedDeleteTransactionPermanent = response.NewGrpcError("transaction", "Failed to permanently delete transaction", int(codes.Internal))

	ErrGrpcFailedRestoreAllTransaction         = response.NewGrpcError("transaction", "Failed to restore all transactions", int(codes.Internal))
	ErrGrpcFailedDeleteAllTransactionPermanent = response.NewGrpcError("transaction", "Failed to permanently delete all transactions", int(codes.Internal))
)
