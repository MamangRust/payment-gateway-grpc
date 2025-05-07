package transfer_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTransferNotFound  = response.NewGrpcError("transfer", "Transfer not found", int(codes.NotFound))
	ErrGrpcTransferInvalidID = response.NewGrpcError("transfer", "Invalid Transfer ID", int(codes.InvalidArgument))
	ErrGrpcInvalidUserID     = response.NewGrpcError("card_id", "Invalid user ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber = response.NewGrpcError("card_id", "Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth      = response.NewGrpcError("month", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear       = response.NewGrpcError("year", "Invalid year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllTransfer            = response.NewGrpcError("transfer", "Failed to fetch all transfers", int(codes.Internal))
	ErrGrpcFailedFindByIdTransfer           = response.NewGrpcError("transfer", "Failed to fetch transfer by ID", int(codes.Internal))
	ErrGrpcFailedFindTransferByTransferFrom = response.NewGrpcError("transfer", "Failed to fetch transfers by sender", int(codes.Internal))
	ErrGrpcFailedFindTransferByTransferTo   = response.NewGrpcError("transfer", "Failed to fetch transfers by receiver", int(codes.Internal))
	ErrGrpcFailedFindByActiveTransfer       = response.NewGrpcError("transfer", "Failed to fetch active transfers", int(codes.Internal))
	ErrGrpcFailedFindByTrashedTransfer      = response.NewGrpcError("transfer", "Failed to fetch trashed transfers", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTransferStatusSuccess             = response.NewGrpcError("transfer", "Failed to fetch monthly successful transfers", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferStatusSuccess              = response.NewGrpcError("transfer", "Failed to fetch yearly successful transfers", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransferStatusFailed              = response.NewGrpcError("transfer", "Failed to fetch monthly failed transfers", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferStatusFailed               = response.NewGrpcError("transfer", "Failed to fetch yearly failed transfers", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransferStatusSuccessByCardNumber = response.NewGrpcError("transfer", "Failed to fetch monthly successful transfers by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferStatusSuccessByCardNumber  = response.NewGrpcError("transfer", "Failed to fetch yearly successful transfers by card number", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransferStatusFailedByCardNumber  = response.NewGrpcError("transfer", "Failed to fetch monthly failed transfers by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferStatusFailedByCardNumber   = response.NewGrpcError("transfer", "Failed to fetch yearly failed transfers by card number", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTransferAmounts                     = response.NewGrpcError("transfer", "Failed to fetch monthly transfer amounts", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferAmounts                      = response.NewGrpcError("transfer", "Failed to fetch yearly transfer amounts", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransferAmountsBySenderCardNumber   = response.NewGrpcError("transfer", "Failed to fetch monthly transfer amounts by sender", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTransferAmountsByReceiverCardNumber = response.NewGrpcError("transfer", "Failed to fetch monthly transfer amounts by receiver", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferAmountsBySenderCardNumber    = response.NewGrpcError("transfer", "Failed to fetch yearly transfer amounts by sender", int(codes.Internal))
	ErrGrpcFailedFindYearlyTransferAmountsByReceiverCardNumber  = response.NewGrpcError("transfer", "Failed to fetch yearly transfer amounts by receiver", int(codes.Internal))

	ErrGrpcFailedCreateTransfer          = response.NewGrpcError("transfer", "Failed to create transfer", int(codes.Internal))
	ErrGrpcFailedUpdateTransfer          = response.NewGrpcError("transfer", "Failed to update transfer", int(codes.Internal))
	ErrGrpcValidateCreateTransferRequest = response.NewGrpcError("transfer", "Invalid input for create transfer", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransferRequest = response.NewGrpcError("transfer", "Invalid input for update transfer", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedTransfer         = response.NewGrpcError("transfer", "Failed to move transfer to trash", int(codes.Internal))
	ErrGrpcFailedRestoreTransfer         = response.NewGrpcError("transfer", "Failed to restore transfer", int(codes.Internal))
	ErrGrpcFailedDeleteTransferPermanent = response.NewGrpcError("transfer", "Failed to permanently delete transfer", int(codes.Internal))

	ErrGrpcFailedRestoreAllTransfer         = response.NewGrpcError("transfer", "Failed to restore all transfers", int(codes.Internal))
	ErrGrpcFailedDeleteAllTransferPermanent = response.NewGrpcError("transfer", "Failed to permanently delete all transfers", int(codes.Internal))
)
