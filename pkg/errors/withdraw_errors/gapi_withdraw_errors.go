package withdraw_errors

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"

	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcWithdrawNotFound  = response.NewGrpcError("withdraw", "Withdraw not found", int(codes.NotFound))
	ErrGrpcWithdrawInvalidID = response.NewGrpcError("withdraw", "Invalid Withdraw ID", int(codes.InvalidArgument))
	ErrGrpcInvalidUserID     = response.NewGrpcError("card_id", "Invalid user ID", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber = response.NewGrpcError("card_id", "Invalid card number", int(codes.InvalidArgument))
	ErrGrpcInvalidMonth      = response.NewGrpcError("month", "Invalid month", int(codes.InvalidArgument))
	ErrGrpcInvalidYear       = response.NewGrpcError("year", "Invalid year", int(codes.InvalidArgument))

	ErrGrpcFailedFindAllWithdraw             = response.NewGrpcError("withdraw", "Failed to fetch all withdraws", int(codes.Internal))
	ErrGrpcFailedFindAllWithdrawByCardNumber = response.NewGrpcError("withdraw", "Failed to fetch withdraws by card number", int(codes.Internal))
	ErrGrpcFailedFindByIdWithdraw            = response.NewGrpcError("withdraw", "Failed to fetch withdraw by ID", int(codes.Internal))
	ErrGrpcFailedFindByCardNumber            = response.NewGrpcError("withdraw", "Failed to fetch withdraws using card number", int(codes.Internal))
	ErrGrpcFailedFindByActiveWithdraw        = response.NewGrpcError("withdraw", "Failed to fetch active withdraws", int(codes.Internal))
	ErrGrpcFailedFindByTrashedWithdraw       = response.NewGrpcError("withdraw", "Failed to fetch trashed withdraws", int(codes.Internal))

	ErrGrpcFailedFindMonthlyWithdrawStatusSuccess           = response.NewGrpcError("withdraw", "Failed to fetch monthly successful withdraws", int(codes.Internal))
	ErrGrpcFailedFindYearlyWithdrawStatusSuccess            = response.NewGrpcError("withdraw", "Failed to fetch yearly successful withdraws", int(codes.Internal))
	ErrGrpcFailedFindMonthlyWithdrawStatusFailed            = response.NewGrpcError("withdraw", "Failed to fetch monthly failed withdraws", int(codes.Internal))
	ErrGrpcFailedFindYearlyWithdrawStatusFailed             = response.NewGrpcError("withdraw", "Failed to fetch yearly failed withdraws", int(codes.Internal))
	ErrGrpcFailedFindMonthlyWithdrawStatusSuccessCardNumber = response.NewGrpcError("withdraw", "Failed to fetch monthly successful withdraws by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyWithdrawStatusSuccessCardNumber  = response.NewGrpcError("withdraw", "Failed to fetch yearly successful withdraws by card number", int(codes.Internal))
	ErrGrpcFailedFindMonthlyWithdrawStatusFailedCardNumber  = response.NewGrpcError("withdraw", "Failed to fetch monthly failed withdraws by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyWithdrawStatusFailedCardNumber   = response.NewGrpcError("withdraw", "Failed to fetch yearly failed withdraws by card number", int(codes.Internal))

	ErrGrpcFailedFindMonthlyWithdraws             = response.NewGrpcError("withdraw", "Failed to fetch monthly withdraw amounts", int(codes.Internal))
	ErrGrpcFailedFindYearlyWithdraws              = response.NewGrpcError("withdraw", "Failed to fetch yearly withdraw amounts", int(codes.Internal))
	ErrGrpcFailedFindMonthlyWithdrawsByCardNumber = response.NewGrpcError("withdraw", "Failed to fetch monthly withdraw amounts by card number", int(codes.Internal))
	ErrGrpcFailedFindYearlyWithdrawsByCardNumber  = response.NewGrpcError("withdraw", "Failed to fetch yearly withdraw amounts by card number", int(codes.Internal))

	ErrGrpcFailedCreateWithdraw          = response.NewGrpcError("withdraw", "Failed to create withdraw", int(codes.Internal))
	ErrGrpcFailedUpdateWithdraw          = response.NewGrpcError("withdraw", "Failed to update withdraw", int(codes.Internal))
	ErrGrpcValidateCreateWithdrawRequest = response.NewGrpcError("withdraw", "Invalid input for create withdraw", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateWithdrawRequest = response.NewGrpcError("withdraw", "Invalid input for update withdraw", int(codes.InvalidArgument))

	ErrGrpcFailedTrashedWithdraw         = response.NewGrpcError("withdraw", "Failed to move withdraw to trash", int(codes.Internal))
	ErrGrpcFailedRestoreWithdraw         = response.NewGrpcError("withdraw", "Failed to restore withdraw", int(codes.Internal))
	ErrGrpcFailedDeleteWithdrawPermanent = response.NewGrpcError("withdraw", "Failed to permanently delete withdraw", int(codes.Internal))

	ErrGrpcFailedRestoreAllWithdraw         = response.NewGrpcError("withdraw", "Failed to restore all withdraws", int(codes.Internal))
	ErrGrpcFailedDeleteAllWithdrawPermanent = response.NewGrpcError("withdraw", "Failed to permanently delete all withdraws", int(codes.Internal))
)
