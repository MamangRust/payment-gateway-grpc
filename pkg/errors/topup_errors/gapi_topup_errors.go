package topup_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTopupNotFound     = response.NewGrpcError("topup", "Topup not found", int(codes.NotFound))
	ErrGrpcTopupInvalidID    = response.NewGrpcError("topup", "Invalid Topup ID", int(codes.InvalidArgument))
	ErrGrpcTopupInvalidMonth = response.NewGrpcError("month", "Invalid Topup Month", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber = response.NewGrpcError("card_id", "Invalid card number", int(codes.InvalidArgument))
	ErrGrpcTopupInvalidYear  = response.NewGrpcError("year", "Invalid Topup Year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllTopup             = response.NewGrpcError("topup", "Failed to fetch all topups", int(codes.Internal))
	ErrGrpcFailedFindAllTopupByCardNumber = response.NewGrpcError("topup", "Failed to fetch topups by card number", int(codes.Internal))
	ErrGrpcFailedFindByIdTopup            = response.NewGrpcError("topup", "Failed to fetch topup by ID", int(codes.Internal))
	ErrGrpcFailedFindByCardNumberTopup    = response.NewGrpcError("topup", "Failed to fetch topup by card number", int(codes.Internal))
	ErrGrpcFailedFindByActiveTopup        = response.NewGrpcError("topup", "Failed to fetch active topups", int(codes.Internal))
	ErrGrpcFailedFindByTrashedTopup       = response.NewGrpcError("topup", "Failed to fetch trashed topups", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTopupStatusSuccess             = response.NewGrpcError("topup", "Failed to fetch monthly successful topups", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupStatusSuccess              = response.NewGrpcError("topup", "Failed to fetch yearly successful topups", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTopupStatusFailed              = response.NewGrpcError("topup", "Failed to fetch monthly failed topups", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupStatusFailed               = response.NewGrpcError("topup", "Failed to fetch yearly failed topups", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTopupStatusSuccessByCardNumber = response.NewGrpcError("topup", "Failed to fetch monthly successful topups by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupStatusSuccessByCardNumber  = response.NewGrpcError("topup", "Failed to fetch yearly successful topups by card number", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTopupStatusFailedByCardNumber  = response.NewGrpcError("topup", "Failed to fetch monthly failed topups by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupStatusFailedByCardNumber   = response.NewGrpcError("topup", "Failed to fetch yearly failed topups by card number", int(codes.Internal))

	ErrGrpcFailedFindMonthlyTopupMethods             = response.NewGrpcError("topup", "Failed to fetch monthly topup methods", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupMethods              = response.NewGrpcError("topup", "Failed to fetch yearly topup methods", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTopupAmounts             = response.NewGrpcError("topup", "Failed to fetch monthly topup amounts", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupAmounts              = response.NewGrpcError("topup", "Failed to fetch yearly topup amounts", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTopupMethodsByCardNumber = response.NewGrpcError("topup", "Failed to fetch monthly topup methods by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupMethodsByCardNumber  = response.NewGrpcError("topup", "Failed to fetch yearly topup methods by card number", int(codes.Internal))
	ErrGrpcFailedFindMonthlyTopupAmountsByCardNumber = response.NewGrpcError("topup", "Failed to fetch monthly topup amounts by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyTopupAmountsByCardNumber  = response.NewGrpcError("topup", "Failed to fetch yearly topup amounts by card number", int(codes.Internal))

	ErrGrpcFailedCreateTopup   = response.NewGrpcError("topup", "Failed to create topup", int(codes.Internal))
	ErrGrpcFailedUpdateTopup   = response.NewGrpcError("topup", "Failed to update topup", int(codes.Internal))
	ErrGrpcValidateCreateTopup = response.NewGrpcError("topup", "Invalid input for create topup", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTopup = response.NewGrpcError("topup", "Invalid input for update topup", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedTopup         = response.NewGrpcError("topup", "Failed to move topup to trash", int(codes.Internal))
	ErrGrpcFailedRestoreTopup         = response.NewGrpcError("topup", "Failed to restore topup", int(codes.Internal))
	ErrGrpcFailedDeleteTopupPermanent = response.NewGrpcError("topup", "Failed to permanently delete topup", int(codes.Internal))

	ErrGrpcFailedRestoreAllTopup         = response.NewGrpcError("topup", "Failed to restore all topups", int(codes.Internal))
	ErrGrpcFailedDeleteAllTopupPermanent = response.NewGrpcError("topup", "Failed to permanently delete all topups", int(codes.Internal))
)
