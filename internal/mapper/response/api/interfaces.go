package apimapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/interfaces.go
type AuthResponseMapper interface {
	ToResponseLogin(res *pb.ApiResponseLogin) *response.ApiResponseLogin
	ToResponseRegister(res *pb.ApiResponseRegister) *response.ApiResponseRegister
	ToResponseRefreshToken(res *pb.ApiResponseRefreshToken) *response.ApiResponseRefreshToken
	ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe
}

type RoleResponseMapper interface {
	ToApiResponseRoleAll(pbResponse *pb.ApiResponseRoleAll) *response.ApiResponseRoleAll
	ToApiResponseRoleDelete(pbResponse *pb.ApiResponseRoleDelete) *response.ApiResponseRoleDelete
	ToApiResponseRole(pbResponse *pb.ApiResponseRole) *response.ApiResponseRole
	ToApiResponsesRole(pbResponse *pb.ApiResponsesRole) *response.ApiResponsesRole
	ToApiResponsePaginationRole(pbResponse *pb.ApiResponsePaginationRole) *response.ApiResponsePaginationRole
	ToApiResponsePaginationRoleDeleteAt(pbResponse *pb.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt
}

type CardResponseMapper interface {
	ToApiResponseCard(card *pb.ApiResponseCard) *response.ApiResponseCard
	ToApiResponsesCard(cards *pb.ApiResponsePaginationCard) *response.ApiResponsePaginationCard
	ToApiResponseCardAll(card *pb.ApiResponseCardAll) *response.ApiResponseCardAll
	ToApiResponseCardDeleteAt(card *pb.ApiResponseCardDelete) *response.ApiResponseCardDelete
	ToApiResponsesCardDeletedAt(cards *pb.ApiResponsePaginationCardDeleteAt) *response.ApiResponsePaginationCardDeleteAt
	ToApiResponseDashboardCard(dash *pb.ApiResponseDashboardCard) *response.ApiResponseDashboardCard
	ToApiResponseDashboardCardCardNumber(dash *pb.ApiResponseDashboardCardNumber) *response.ApiResponseDashboardCardNumber
	ToApiResponseMonthlyBalances(cards *pb.ApiResponseMonthlyBalance) *response.ApiResponseMonthlyBalance
	ToApiResponseYearlyBalances(cards *pb.ApiResponseYearlyBalance) *response.ApiResponseYearlyBalance
	ToApiResponseMonthlyAmounts(cards *pb.ApiResponseMonthlyAmount) *response.ApiResponseMonthlyAmount
	ToApiResponseYearlyAmounts(cards *pb.ApiResponseYearlyAmount) *response.ApiResponseYearlyAmount
}

type MerchantResponseMapper interface {
	ToApiResponseMerchant(merchants *pb.ApiResponseMerchant) *response.ApiResponseMerchant
	ToApiResponseMerchants(merchants *pb.ApiResponsesMerchant) *response.ApiResponsesMerchant
	ToApiResponsesMerchant(merchants *pb.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant
	ToApiResponsesMerchantDeleteAt(merchants *pb.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
	ToApiResponseMerchantAll(card *pb.ApiResponseMerchantAll) *response.ApiResponseMerchantAll
	ToApiResponseMerchantDeleteAt(card *pb.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete
	ToApiResponseMerchantsTransactionResponse(merchants *pb.ApiResponsePaginationMerchantTransaction) *response.ApiResponsePaginationMerchantTransaction
	ToApiResponseMonthlyPaymentMethods(ms *pb.ApiResponseMerchantMonthlyPaymentMethod) *response.ApiResponseMerchantMonthlyPaymentMethod
	ToApiResponseYearlyPaymentMethods(ms *pb.ApiResponseMerchantYearlyPaymentMethod) *response.ApiResponseMerchantYearlyPaymentMethod
	ToApiResponseMonthlyAmounts(ms *pb.ApiResponseMerchantMonthlyAmount) *response.ApiResponseMerchantMonthlyAmount
	ToApiResponseYearlyAmounts(ms *pb.ApiResponseMerchantYearlyAmount) *response.ApiResponseMerchantYearlyAmount
	ToApiResponseMonthlyTotalAmounts(ms *pb.ApiResponseMerchantMonthlyTotalAmount) *response.ApiResponseMerchantMonthlyTotalAmount
	ToApiResponseYearlyTotalAmounts(ms *pb.ApiResponseMerchantYearlyTotalAmount) *response.ApiResponseMerchantYearlyTotalAmount
}

type SaldoResponseMapper interface {
	ToApiResponseSaldo(pbResponse *pb.ApiResponseSaldo) *response.ApiResponseSaldo
	ToApiResponsesSaldo(pbResponse *pb.ApiResponsesSaldo) *response.ApiResponsesSaldo
	ToApiResponseSaldoDelete(pbResponse *pb.ApiResponseSaldoDelete) *response.ApiResponseSaldoDelete
	ToApiResponseSaldoAll(pbResponse *pb.ApiResponseSaldoAll) *response.ApiResponseSaldoAll
	ToApiResponseMonthTotalSaldo(pbResponse *pb.ApiResponseMonthTotalSaldo) *response.ApiResponseMonthTotalSaldo
	ToApiResponseYearTotalSaldo(pbResponse *pb.ApiResponseYearTotalSaldo) *response.ApiResponseYearTotalSaldo
	ToApiResponseMonthSaldoBalances(pbResponse *pb.ApiResponseMonthSaldoBalances) *response.ApiResponseMonthSaldoBalances
	ToApiResponseYearSaldoBalances(pbResponse *pb.ApiResponseYearSaldoBalances) *response.ApiResponseYearSaldoBalances
	ToApiResponsePaginationSaldo(pbResponse *pb.ApiResponsePaginationSaldo) *response.ApiResponsePaginationSaldo
	ToApiResponsePaginationSaldoDeleteAt(pbResponse *pb.ApiResponsePaginationSaldoDeleteAt) *response.ApiResponsePaginationSaldoDeleteAt
}

type TopupResponseMapper interface {
	ToApiResponseTopup(s *pb.ApiResponseTopup) *response.ApiResponseTopup
	ToApiResponseTopupAll(s *pb.ApiResponseTopupAll) *response.ApiResponseTopupAll
	ToApiResponseTopupDelete(s *pb.ApiResponseTopupDelete) *response.ApiResponseTopupDelete
	ToApiResponsePaginationTopup(s *pb.ApiResponsePaginationTopup) *response.ApiResponsePaginationTopup
	ToApiResponsePaginationTopupDeleteAt(s *pb.ApiResponsePaginationTopupDeleteAt) *response.ApiResponsePaginationTopupDeleteAt
	ToApiResponseTopupMonthStatusSuccess(s *pb.ApiResponseTopupMonthStatusSuccess) *response.ApiResponseTopupMonthStatusSuccess
	ToApiResponseTopupYearStatusSuccess(s *pb.ApiResponseTopupYearStatusSuccess) *response.ApiResponseTopupYearStatusSuccess
	ToApiResponseTopupMonthStatusFailed(s *pb.ApiResponseTopupMonthStatusFailed) *response.ApiResponseTopupMonthStatusFailed
	ToApiResponseTopupYearStatusFailed(s *pb.ApiResponseTopupYearStatusFailed) *response.ApiResponseTopupYearStatusFailed
	ToApiResponseTopupMonthMethod(s *pb.ApiResponseTopupMonthMethod) *response.ApiResponseTopupMonthMethod
	ToApiResponseTopupYearMethod(s *pb.ApiResponseTopupYearMethod) *response.ApiResponseTopupYearMethod
	ToApiResponseTopupMonthAmount(s *pb.ApiResponseTopupMonthAmount) *response.ApiResponseTopupMonthAmount
	ToApiResponseTopupYearAmount(s *pb.ApiResponseTopupYearAmount) *response.ApiResponseTopupYearAmount
}

type TransactionResponseMapper interface {
	ToApiResponseTransactionMonthStatusSuccess(pbResponse *pb.ApiResponseTransactionMonthStatusSuccess) *response.ApiResponseTransactionMonthStatusSuccess
	ToApiResponseTransactionYearStatusSuccess(pbResponse *pb.ApiResponseTransactionYearStatusSuccess) *response.ApiResponseTransactionYearStatusSuccess
	ToApiResponseTransactionMonthStatusFailed(pbResponse *pb.ApiResponseTransactionMonthStatusFailed) *response.ApiResponseTransactionMonthStatusFailed
	ToApiResponseTransactionYearStatusFailed(pbResponse *pb.ApiResponseTransactionYearStatusFailed) *response.ApiResponseTransactionYearStatusFailed
	ToApiResponseTransactionMonthMethod(pbResponse *pb.ApiResponseTransactionMonthMethod) *response.ApiResponseTransactionMonthMethod
	ToApiResponseTransactionYearMethod(pbResponse *pb.ApiResponseTransactionYearMethod) *response.ApiResponseTransactionYearMethod
	ToApiResponseTransactionMonthAmount(pbResponse *pb.ApiResponseTransactionMonthAmount) *response.ApiResponseTransactionMonthAmount
	ToApiResponseTransactionYearAmount(pbResponse *pb.ApiResponseTransactionYearAmount) *response.ApiResponseTransactionYearAmount
	ToApiResponseTransaction(pbResponse *pb.ApiResponseTransaction) *response.ApiResponseTransaction
	ToApiResponseTransactions(pbResponse *pb.ApiResponseTransactions) *response.ApiResponseTransactions
	ToApiResponseTransactionDelete(pbResponse *pb.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete
	ToApiResponseTransactionAll(pbResponse *pb.ApiResponseTransactionAll) *response.ApiResponseTransactionAll
	ToApiResponsePaginationTransaction(pbResponse *pb.ApiResponsePaginationTransaction) *response.ApiResponsePaginationTransaction
	ToApiResponsePaginationTransactionDeleteAt(pbResponse *pb.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt
}

type TransferResponseMapper interface {
	ToApiResponseTransferMonthStatusSuccess(pbResponse *pb.ApiResponseTransferMonthStatusSuccess) *response.ApiResponseTransferMonthStatusSuccess
	ToApiResponseTransferYearStatusSuccess(pbResponse *pb.ApiResponseTransferYearStatusSuccess) *response.ApiResponseTransferYearStatusSuccess
	ToApiResponseTransferMonthStatusFailed(pbResponse *pb.ApiResponseTransferMonthStatusFailed) *response.ApiResponseTransferMonthStatusFailed
	ToApiResponseTransferYearStatusFailed(pbResponse *pb.ApiResponseTransferYearStatusFailed) *response.ApiResponseTransferYearStatusFailed
	ToApiResponseTransferMonthAmount(pbResponse *pb.ApiResponseTransferMonthAmount) *response.ApiResponseTransferMonthAmount
	ToApiResponseTransferYearAmount(pbResponse *pb.ApiResponseTransferYearAmount) *response.ApiResponseTransferYearAmount
	ToApiResponseTransfer(pbResponse *pb.ApiResponseTransfer) *response.ApiResponseTransfer
	ToApiResponseTransfers(pbResponse *pb.ApiResponseTransfers) *response.ApiResponseTransfers
	ToApiResponseTransferDelete(pbResponse *pb.ApiResponseTransferDelete) *response.ApiResponseTransferDelete
	ToApiResponseTransferAll(pbResponse *pb.ApiResponseTransferAll) *response.ApiResponseTransferAll
	ToApiResponsePaginationTransfer(pbResponse *pb.ApiResponsePaginationTransfer) *response.ApiResponsePaginationTransfer
	ToApiResponsePaginationTransferDeleteAt(pbResponse *pb.ApiResponsePaginationTransferDeleteAt) *response.ApiResponsePaginationTransferDeleteAt
}

type UserResponseMapper interface {
	ToApiResponseUser(pbResponse *pb.ApiResponseUser) *response.ApiResponseUser
	ToApiResponsesUser(pbResponse *pb.ApiResponsesUser) *response.ApiResponsesUser
	ToApiResponseUserDelete(pbResponse *pb.ApiResponseUserDelete) *response.ApiResponseUserDelete
	ToApiResponseUserAll(pbResponse *pb.ApiResponseUserAll) *response.ApiResponseUserAll
	ToApiResponsePaginationUserDeleteAt(pbResponse *pb.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt
	ToApiResponsePaginationUser(pbResponse *pb.ApiResponsePaginationUser) *response.ApiResponsePaginationUser
}

type WithdrawResponseMapper interface {
	ToApiResponseWithdraw(pbResponse *pb.ApiResponseWithdraw) *response.ApiResponseWithdraw
	ToApiResponsesWithdraw(pbResponse *pb.ApiResponsesWithdraw) *response.ApiResponsesWithdraw
	ToApiResponseWithdrawDelete(pbResponse *pb.ApiResponseWithdrawDelete) *response.ApiResponseWithdrawDelete
	ToApiResponseWithdrawAll(pbResponse *pb.ApiResponseWithdrawAll) *response.ApiResponseWithdrawAll
	ToApiResponsePaginationWithdraw(pbResponse *pb.ApiResponsePaginationWithdraw) *response.ApiResponsePaginationWithdraw
	ToApiResponsePaginationWithdrawDeleteAt(pbResponse *pb.ApiResponsePaginationWithdrawDeleteAt) *response.ApiResponsePaginationWithdrawDeleteAt
	ToApiResponseWithdrawMonthStatusSuccess(pbResponse *pb.ApiResponseWithdrawMonthStatusSuccess) *response.ApiResponseWithdrawMonthStatusSuccess
	ToApiResponseWithdrawYearStatusSuccess(pbResponse *pb.ApiResponseWithdrawYearStatusSuccess) *response.ApiResponseWithdrawYearStatusSuccess
	ToApiResponseWithdrawMonthStatusFailed(pbResponse *pb.ApiResponseWithdrawMonthStatusFailed) *response.ApiResponseWithdrawMonthStatusFailed
	ToApiResponseWithdrawYearStatusFailed(pbResponse *pb.ApiResponseWithdrawYearStatusFailed) *response.ApiResponseWithdrawYearStatusFailed
	ToApiResponseWithdrawMonthAmount(pbResponse *pb.ApiResponseWithdrawMonthAmount) *response.ApiResponseWithdrawMonthAmount
	ToApiResponseWithdrawYearAmount(pbResponse *pb.ApiResponseWithdrawYearAmount) *response.ApiResponseWithdrawYearAmount
}
