package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandleGrpc interface {
	pb.AuthServiceServer
	LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error)
	RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error)
}

type RoleHandleGrpc interface {
	pb.RoleServiceServer

	FindAllRole(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRole, error)
	FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error)
	FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error)
	FindByActive(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error)
	CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.ApiResponseRole, error)
	UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error)
	TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error)
	RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error)
	DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error)
	RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error)
	DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error)
}

type UserHandleGrpc interface {
	pb.UserServiceServer

	FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error)
	FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error)
	FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error)
	Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error)
	Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error)
	TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error)
	RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error)
	DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error)

	RestoreAllUser(context.Context, *emptypb.Empty) (*pb.ApiResponseUserAll, error)
	DeleteAllUserPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseUserAll, error)
}

type CardHandleGrpc interface {
	pb.CardServiceServer

	FindAllCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCard, error)
	FindByIdCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error)
	FindByUserIdCard(ctx context.Context, req *pb.FindByUserIdCardRequest) (*pb.ApiResponseCard, error)
	FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseCard, error)

	DashboardCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseDashboardCard, error)
	DashboardCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseDashboardCardNumber, error)

	FindMonthlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseMonthlyBalance, error)
	FindYearlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseYearlyBalance, error)
	FindMonthlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error)

	FindMonthlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseMonthlyBalance, error)
	FindYearlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseYearlyBalance, error)
	FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)

	FindMonthlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)
	FindMonthlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error)
	FindYearlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error)

	FindByActiveCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error)
	FindByTrashedCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error)
	CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.ApiResponseCard, error)
	UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.ApiResponseCard, error)
	TrashedCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error)
	RestoreCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error)
	DeleteCardPermanent(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDelete, error)

	RestoreAllCard(context.Context, *emptypb.Empty) (*pb.ApiResponseCardAll, error)
	DeleteAllCardPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseCardAll, error)
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer

	FindAllMerchant(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error)
	FindByIdMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error)

	FindMonthlyPaymentMethodsMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error)
	FindYearlyPaymentMethodMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyPaymentMethod, error)
	FindMonthlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantMonthlyAmount, error)
	FindYearlyAmountMerchant(ctx context.Context, req *pb.FindYearMerchant) (*pb.ApiResponseMerchantYearlyAmount, error)

	FindAllTransactionByMerchant(ctx context.Context, req *pb.FindAllMerchantTransaction) (*pb.ApiResponsePaginationMerchantTransaction, error)
	FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error)
	FindYearlyPaymentMethodByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyPaymentMethod, error)
	FindMonthlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyAmount, error)
	FindYearlyAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyAmount, error)
	FindMonthlyTotalAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantMonthlyTotalAmount, error)
	FindYearlyTotalAmountByMerchants(ctx context.Context, req *pb.FindYearMerchantById) (*pb.ApiResponseMerchantYearlyTotalAmount, error)

	FindAllTransactionByApikey(ctx context.Context, req *pb.FindAllMerchantApikey) (*pb.ApiResponsePaginationMerchantTransaction, error)
	FindMonthlyPaymentMethodByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyPaymentMethod, error)
	FindYearlyPaymentMethodByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyPaymentMethod, error)
	FindMonthlyAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyAmount, error)
	FindYearlyAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyAmount, error)
	FindMonthlyTotalAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantMonthlyTotalAmount, error)
	FindYearlyTotalAmountByApikey(ctx context.Context, req *pb.FindYearMerchantByApikey) (*pb.ApiResponseMerchantYearlyTotalAmount, error)

	FindByApiKey(ctx context.Context, req *pb.FindByApiKeyRequest) (*pb.ApiResponseMerchant, error)

	FindByMerchantUserId(ctx context.Context, req *pb.FindByMerchantUserIdRequest) (*pb.ApiResponsesMerchant, error)
	FindByActive(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error)
	CreateMerchant(ctx context.Context, req *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error)
	UpdateMerchant(ctx context.Context, req *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error)
	TrashedMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error)
	RestoreMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error)
	DeleteMerchant(ctx context.Context, req *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error)

	RestoreAllMerchant(context.Context, *emptypb.Empty) (*pb.ApiResponseMerchantAll, error)
	DeleteAllMerchantPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseMerchantAll, error)
}

type SaldoHandleGrpc interface {
	pb.SaldoServiceServer

	FindAllSaldo(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldo, error)
	FindByIdSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error)

	FindMonthlyTotalSaldoBalance(ctx context.Context, req *pb.FindMonthlySaldoTotalBalance) (*pb.ApiResponseMonthTotalSaldo, error)
	FindYearTotalSaldoBalance(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearTotalSaldo, error)

	FindMonthlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseMonthSaldoBalances, error)
	FindYearlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearSaldoBalances, error)

	FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseSaldo, error)
	FindByActive(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error)
	CreateSaldo(ctx context.Context, req *pb.CreateSaldoRequest) (*pb.ApiResponseSaldo, error)
	UpdateSaldo(ctx context.Context, req *pb.UpdateSaldoRequest) (*pb.ApiResponseSaldo, error)
	TrashedSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error)
	RestoreSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error)
	DeleteSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDelete, error)

	RestoreAllSaldo(context.Context, *emptypb.Empty) (*pb.ApiResponseSaldoAll, error)
	DeleteAllSaldoPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseSaldoAll, error)
}

type TopupHandleGrpc interface {
	pb.TopupServiceServer

	FindAllTopup(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopup, error)
	FindAllTopupByCardNumber(ctx context.Context, req *pb.FindAllTopupByCardNumberRequest) (*pb.ApiResponsePaginationTopup, error)
	FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error)

	FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error)
	FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error)
	FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error)
	FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error)

	FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error)
	FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error)
	FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error)
	FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error)

	FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error)
	FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error)

	FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error)
	FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error)

	FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error)
	FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error)

	FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error)
	FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error)

	FindByActive(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error)
	CreateTopup(ctx context.Context, req *pb.CreateTopupRequest) (*pb.ApiResponseTopup, error)
	UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error)
	TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error)
	RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error)
	DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error)

	RestoreAllTopup(context.Context, *emptypb.Empty) (*pb.ApiResponseTopupAll, error)
	DeleteAllTopupPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseTopupAll, error)
}

type TransactionHandleGrpc interface {
	pb.TransactionServiceServer

	FindAllTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error)
	FindByIdTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error)

	FindMonthlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthMethod, error)
	FindYearlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearMethod, error)
	FindMonthlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthAmount, error)
	FindYearlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearAmount, error)

	FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthMethod, error)
	FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearMethod, error)
	FindMonthlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthAmount, error)
	FindYearlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearAmount, error)

	FindTransactionByMerchantIdRequest(ctx context.Context, request *pb.FindTransactionByMerchantIdRequest) (*pb.ApiResponseTransactions, error)
	FindByActiveTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error)
	FindByTrashedTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error)
	CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error)
	UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error)
	TrashedTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error)
	RestoreTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error)
	DeleteTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error)

	RestoreAllTransaction(context.Context, *emptypb.Empty) (*pb.ApiResponseTransactionAll, error)
	DeleteAllTransactionPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseTransactionAll, error)
}

type TransferHandleGrpc interface {
	pb.TransferServiceServer

	FindAllTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransfer, error)
	FindByIdTransfer(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error)

	FindMonthlyTransferStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusSuccess, error)
	FindYearlyTransferStatusSuccess(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusSuccess, error)
	FindMonthlyTransferStatusFailed(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusFailed, error)
	FindYearlyTransferStatusFailed(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusFailed, error)

	FindMonthlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusSuccess, error)
	FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusSuccess, error)
	FindMonthlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusFailed, error)
	FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusFailed, error)

	FindMonthlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferMonthAmount, error)
	FindYearlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearAmount, error)

	FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error)
	FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error)
	FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error)
	FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error)

	FindByTransferByTransferFrom(ctx context.Context, request *pb.FindTransferByTransferFromRequest) (*pb.ApiResponseTransfers, error)
	FindByTransferByTransferTo(ctx context.Context, request *pb.FindTransferByTransferToRequest) (*pb.ApiResponseTransfers, error)
	FindByActiveTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error)
	FindByTrashedTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error)
	CreateTransfer(ctx context.Context, req *pb.CreateTransferRequest) (*pb.ApiResponseTransfer, error)
	UpdateTransfer(ctx context.Context, req *pb.UpdateTransferRequest) (*pb.ApiResponseTransfer, error)
	TrashedTransfer(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error)
	RestoreTransfer(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error)
	DeleteTransferPermanent(ctx context.Context, req *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDelete, error)

	RestoreAllTransfer(context.Context, *emptypb.Empty) (*pb.ApiResponseTransferAll, error)
	DeleteAllTransferPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseTransferAll, error)
}

type WithdrawHandleGrpc interface {
	pb.WithdrawServiceServer

	FindAllWithdraw(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdraw, error)
	FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error)

	FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error)
	FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error)
	FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawMonthAmount, error)
	FindYearlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawYearAmount, error)

	FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponsesWithdraw, error)
	FindByActive(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error)
	FindByTrashed(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error)
	CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error)

	RestoreAllWithdraw(context.Context, *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error)
	DeleteAllWithdrawPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error)
}
