package card_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcInvalidCardID     = response.NewGrpcError("card_id", "Invalid card ID", int(codes.InvalidArgument))
	ErrGrpcInvalidUserID     = response.NewGrpcError("card_id", "Invalid user ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber = response.NewGrpcError("card_id", "Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth      = response.NewGrpcError("month", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear       = response.NewGrpcError("year", "Invalid year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllCard       = response.NewGrpcError("card", "Failed to fetch cards", int(codes.Internal))
	ErrGrpcCardNotFound            = response.NewGrpcError("card", "Card not found", int(codes.NotFound))
	ErrGrpcFailedFindByIdCard      = response.NewGrpcError("card", "Failed to fetch card by ID", int(codes.Internal))
	ErrGrpcFailedFindByUserIdCard  = response.NewGrpcError("card", "Failed to fetch card by user ID", int(codes.Internal))
	ErrGrpcFailedFindByActiveCard  = response.NewGrpcError("card", "Failed to fetch active cards", int(codes.Internal))
	ErrGrpcFailedFindByTrashedCard = response.NewGrpcError("card", "Failed to fetch trashed cards", int(codes.Internal))
	ErrGrpcFailedFindByCardNumber  = response.NewGrpcError("card", "Failed to fetch card by card number", int(codes.Internal))

	ErrGrpcFailedDashboardCard       = response.NewGrpcError("card", "Failed to fetch card dashboard data", int(codes.Internal))
	ErrGrpcFailedDashboardCardNumber = response.NewGrpcError("card", "Failed to fetch dashboard data by card number", int(codes.Internal))

	ErrGrpcFailedMonthlyBalance       = response.NewGrpcError("card", "Failed to fetch monthly balance", int(codes.Internal))
	ErrGrpcFailedYearlyBalance        = response.NewGrpcError("card", "Failed to fetch yearly balance", int(codes.Internal))
	ErrGrpcFailedMonthlyBalanceByCard = response.NewGrpcError("card", "Failed to fetch monthly balance by card number", int(codes.Internal))
	ErrGrpcFailedYearlyBalanceByCard  = response.NewGrpcError("card", "Failed to fetch yearly balance by card number", int(codes.Internal))

	ErrGrpcFailedMonthlyTopupAmount = response.NewGrpcError("card", "Failed to fetch monthly topup amount", int(codes.Internal))
	ErrGrpcFailedYearlyTopupAmount  = response.NewGrpcError("card", "Failed to fetch yearly topup amount", int(codes.Internal))
	ErrGrpcFailedMonthlyTopupByCard = response.NewGrpcError("card", "Failed to fetch monthly topup amount by card number", int(codes.Internal))
	ErrGrpcFailedYearlyTopupByCard  = response.NewGrpcError("card", "Failed to fetch yearly topup amount by card number", int(codes.Internal))

	ErrGrpcFailedMonthlyWithdrawAmount = response.NewGrpcError("card", "Failed to fetch monthly withdraw amount", int(codes.Internal))
	ErrGrpcFailedYearlyWithdrawAmount  = response.NewGrpcError("card", "Failed to fetch yearly withdraw amount", int(codes.Internal))
	ErrGrpcFailedMonthlyWithdrawByCard = response.NewGrpcError("card", "Failed to fetch monthly withdraw amount by card number", int(codes.Internal))
	ErrGrpcFailedYearlyWithdrawByCard  = response.NewGrpcError("card", "Failed to fetch yearly withdraw amount by card number", int(codes.Internal))

	ErrGrpcFailedMonthlyTransactionAmount = response.NewGrpcError("card", "Failed to fetch monthly transaction amount", int(codes.Internal))
	ErrGrpcFailedYearlyTransactionAmount  = response.NewGrpcError("card", "Failed to fetch yearly transaction amount", int(codes.Internal))
	ErrGrpcFailedMonthlyTransactionByCard = response.NewGrpcError("card", "Failed to fetch monthly transaction amount by card number", int(codes.Internal))
	ErrGrpcFailedYearlyTransactionByCard  = response.NewGrpcError("card", "Failed to fetch yearly transaction amount by card number", int(codes.Internal))

	ErrGrpcFailedMonthlyTransferSender       = response.NewGrpcError("card", "Failed to fetch monthly transfer sender amount", int(codes.Internal))
	ErrGrpcFailedYearlyTransferSender        = response.NewGrpcError("card", "Failed to fetch yearly transfer sender amount", int(codes.Internal))
	ErrGrpcFailedMonthlyTransferSenderByCard = response.NewGrpcError("card", "Failed to fetch monthly transfer sender amount by card number", int(codes.Internal))
	ErrGrpcFailedYearlyTransferSenderByCard  = response.NewGrpcError("card", "Failed to fetch yearly transfer sender amount by card number", int(codes.Internal))

	ErrGrpcFailedMonthlyTransferReceiver       = response.NewGrpcError("card", "Failed to fetch monthly transfer receiver amount", int(codes.Internal))
	ErrGrpcFailedYearlyTransferReceiver        = response.NewGrpcError("card", "Failed to fetch yearly transfer receiver amount", int(codes.Internal))
	ErrGrpcFailedMonthlyTransferReceiverByCard = response.NewGrpcError("card", "Failed to fetch monthly transfer receiver amount by card number", int(codes.Internal))
	ErrGrpcFailedYearlyTransferReceiverByCard  = response.NewGrpcError("card", "Failed to fetch yearly transfer receiver amount by card number", int(codes.Internal))

	ErrGrpcFailedCreateCard          = response.NewGrpcError("card", "Failed to create card", int(codes.Internal))
	ErrGrpcFailedUpdateCard          = response.NewGrpcError("card", "Failed to update card", int(codes.Internal))
	ErrGrpcValidateCreateCardRequest = response.NewGrpcError("card", "Invalid input for create card", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCardRequest = response.NewGrpcError("card", "Invalid input for update card", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedCard         = response.NewGrpcError("card", "Failed to move card to trash", int(codes.Internal))
	ErrGrpcFailedRestoreCard         = response.NewGrpcError("card", "Failed to restore card", int(codes.Internal))
	ErrGrpcFailedDeleteCardPermanent = response.NewGrpcError("card", "Failed to permanently delete card", int(codes.Internal))

	ErrGrpcFailedRestoreAllCards         = response.NewGrpcError("card", "Failed to restore all cards", int(codes.Internal))
	ErrGrpcFailedDeleteAllCardsPermanent = response.NewGrpcError("card", "Failed to permanently delete all cards", int(codes.Internal))
)
