package responseservice

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/interfaces.go
type CardResponseMapper interface {
	ToCardResponse(card *record.CardRecord) *response.CardResponse
	ToCardsResponse(cards []*record.CardRecord) []*response.CardResponse

	ToCardResponseDeleteAt(card *record.CardRecord) *response.CardResponseDeleteAt
	ToCardsResponseDeleteAt(cards []*record.CardRecord) []*response.CardResponseDeleteAt

	ToGetMonthlyBalance(card *record.CardMonthBalance) *response.CardResponseMonthBalance
	ToGetMonthlyBalances(cards []*record.CardMonthBalance) []*response.CardResponseMonthBalance

	ToGetYearlyBalance(card *record.CardYearlyBalance) *response.CardResponseYearlyBalance
	ToGetYearlyBalances(cards []*record.CardYearlyBalance) []*response.CardResponseYearlyBalance

	ToGetMonthlyAmount(card *record.CardMonthAmount) *response.CardResponseMonthAmount
	ToGetMonthlyAmounts(cards []*record.CardMonthAmount) []*response.CardResponseMonthAmount
	ToGetYearlyAmount(card *record.CardYearAmount) *response.CardResponseYearAmount
	ToGetYearlyAmounts(cards []*record.CardYearAmount) []*response.CardResponseYearAmount
}

type UserResponseMapper interface {
	ToUserResponse(user *record.UserRecord) *response.UserResponse
	ToUsersResponse(users []*record.UserRecord) []*response.UserResponse

	ToUserResponseDeleteAt(user *record.UserRecord) *response.UserResponseDeleteAt
	ToUsersResponseDeleteAt(users []*record.UserRecord) []*response.UserResponseDeleteAt
}

type RoleResponseMapper interface {
	ToRoleResponse(role *record.RoleRecord) *response.RoleResponse
	ToRolesResponse(roles []*record.RoleRecord) []*response.RoleResponse

	ToRoleResponseDeleteAt(role *record.RoleRecord) *response.RoleResponseDeleteAt
	ToRolesResponseDeleteAt(roles []*record.RoleRecord) []*response.RoleResponseDeleteAt
}

type RefreshTokenResponseMapper interface {
	ToRefreshTokenResponse(refresh *record.RefreshTokenRecord) *response.RefreshTokenResponse
	ToRefreshTokenResponses(refreshs []*record.RefreshTokenRecord) []*response.RefreshTokenResponse
}

type SaldoResponseMapper interface {
	ToSaldoResponse(saldo *record.SaldoRecord) *response.SaldoResponse
	ToSaldoResponses(saldos []*record.SaldoRecord) []*response.SaldoResponse

	ToSaldoMonthTotalBalanceResponse(ss *record.SaldoMonthTotalBalance) *response.SaldoMonthTotalBalanceResponse
	ToSaldoMonthTotalBalanceResponses(ss []*record.SaldoMonthTotalBalance) []*response.SaldoMonthTotalBalanceResponse

	ToSaldoYearTotalBalanceResponse(ss *record.SaldoYearTotalBalance) *response.SaldoYearTotalBalanceResponse
	ToSaldoYearTotalBalanceResponses(ss []*record.SaldoYearTotalBalance) []*response.SaldoYearTotalBalanceResponse

	ToSaldoMonthBalanceResponse(ss *record.SaldoMonthSaldoBalance) *response.SaldoMonthBalanceResponse
	ToSaldoMonthBalanceResponses(ss []*record.SaldoMonthSaldoBalance) []*response.SaldoMonthBalanceResponse

	ToSaldoYearBalanceResponse(ss *record.SaldoYearSaldoBalance) *response.SaldoYearBalanceResponse
	ToSaldoYearBalanceResponses(ss []*record.SaldoYearSaldoBalance) []*response.SaldoYearBalanceResponse

	ToSaldoResponseDeleteAt(saldo *record.SaldoRecord) *response.SaldoResponseDeleteAt
	ToSaldoResponsesDeleteAt(saldos []*record.SaldoRecord) []*response.SaldoResponseDeleteAt
}

type TopupResponseMapper interface {
	ToTopupResponse(topup *record.TopupRecord) *response.TopupResponse
	ToTopupResponses(topups []*record.TopupRecord) []*response.TopupResponse

	ToTopupResponseMonthStatusSuccess(s *record.TopupRecordMonthStatusSuccess) *response.TopupResponseMonthStatusSuccess
	ToTopupResponsesMonthStatusSuccess(topups []*record.TopupRecordMonthStatusSuccess) []*response.TopupResponseMonthStatusSuccess
	ToTopupResponseYearStatusSuccess(s *record.TopupRecordYearStatusSuccess) *response.TopupResponseYearStatusSuccess
	ToTopupResponsesYearStatusSuccess(topups []*record.TopupRecordYearStatusSuccess) []*response.TopupResponseYearStatusSuccess

	ToTopupResponseMonthStatusFailed(s *record.TopupRecordMonthStatusFailed) *response.TopupResponseMonthStatusFailed
	ToTopupResponsesMonthStatusFailed(topups []*record.TopupRecordMonthStatusFailed) []*response.TopupResponseMonthStatusFailed
	ToTopupResponseYearStatusFailed(s *record.TopupRecordYearStatusFailed) *response.TopupResponseYearStatusFailed
	ToTopupResponsesYearStatusFailed(topups []*record.TopupRecordYearStatusFailed) []*response.TopupResponseYearStatusFailed

	ToTopupMonthlyMethodResponse(s *record.TopupMonthMethod) *response.TopupMonthMethodResponse
	ToTopupMonthlyMethodResponses(s []*record.TopupMonthMethod) []*response.TopupMonthMethodResponse
	ToTopupYearlyMethodResponse(s *record.TopupYearlyMethod) *response.TopupYearlyMethodResponse
	ToTopupYearlyMethodResponses(s []*record.TopupYearlyMethod) []*response.TopupYearlyMethodResponse

	ToTopupMonthlyAmountResponse(s *record.TopupMonthAmount) *response.TopupMonthAmountResponse
	ToTopupMonthlyAmountResponses(s []*record.TopupMonthAmount) []*response.TopupMonthAmountResponse
	ToTopupYearlyAmountResponse(s *record.TopupYearlyAmount) *response.TopupYearlyAmountResponse
	ToTopupYearlyAmountResponses(s []*record.TopupYearlyAmount) []*response.TopupYearlyAmountResponse

	ToTopupResponseDeleteAt(topup *record.TopupRecord) *response.TopupResponseDeleteAt
	ToTopupResponsesDeleteAt(topups []*record.TopupRecord) []*response.TopupResponseDeleteAt
}

type TransactionResponseMapper interface {
	ToTransactionResponse(transaction *record.TransactionRecord) *response.TransactionResponse
	ToTransactionsResponse(transactions []*record.TransactionRecord) []*response.TransactionResponse

	ToTransactionResponseMonthStatusSuccess(s *record.TransactionRecordMonthStatusSuccess) *response.TransactionResponseMonthStatusSuccess
	ToTransactionResponsesMonthStatusSuccess(Transactions []*record.TransactionRecordMonthStatusSuccess) []*response.TransactionResponseMonthStatusSuccess
	ToTransactionResponseYearStatusSuccess(s *record.TransactionRecordYearStatusSuccess) *response.TransactionResponseYearStatusSuccess
	ToTransactionResponsesYearStatusSuccess(Transactions []*record.TransactionRecordYearStatusSuccess) []*response.TransactionResponseYearStatusSuccess

	ToTransactionResponseMonthStatusFailed(s *record.TransactionRecordMonthStatusFailed) *response.TransactionResponseMonthStatusFailed
	ToTransactionResponsesMonthStatusFailed(Transactions []*record.TransactionRecordMonthStatusFailed) []*response.TransactionResponseMonthStatusFailed
	ToTransactionResponseYearStatusFailed(s *record.TransactionRecordYearStatusFailed) *response.TransactionResponseYearStatusFailed
	ToTransactionResponsesYearStatusFailed(Transactions []*record.TransactionRecordYearStatusFailed) []*response.TransactionResponseYearStatusFailed

	ToTransactionMonthlyMethodResponse(s *record.TransactionMonthMethod) *response.TransactionMonthMethodResponse
	ToTransactionMonthlyMethodResponses(s []*record.TransactionMonthMethod) []*response.TransactionMonthMethodResponse
	ToTransactionYearlyMethodResponse(s *record.TransactionYearMethod) *response.TransactionYearMethodResponse
	ToTransactionYearlyMethodResponses(s []*record.TransactionYearMethod) []*response.TransactionYearMethodResponse

	ToTransactionMonthlyAmountResponse(s *record.TransactionMonthAmount) *response.TransactionMonthAmountResponse
	ToTransactionMonthlyAmountResponses(s []*record.TransactionMonthAmount) []*response.TransactionMonthAmountResponse
	ToTransactionYearlyAmountResponse(s *record.TransactionYearlyAmount) *response.TransactionYearlyAmountResponse
	ToTransactionYearlyAmountResponses(s []*record.TransactionYearlyAmount) []*response.TransactionYearlyAmountResponse

	ToTransactionResponseDeleteAt(transaction *record.TransactionRecord) *response.TransactionResponseDeleteAt
	ToTransactionsResponseDeleteAt(transactions []*record.TransactionRecord) []*response.TransactionResponseDeleteAt
}

type TransferResponseMapper interface {
	ToTransferResponse(transfer *record.TransferRecord) *response.TransferResponse
	ToTransfersResponse(transfers []*record.TransferRecord) []*response.TransferResponse

	ToTransferResponseMonthStatusSuccess(s *record.TransferRecordMonthStatusSuccess) *response.TransferResponseMonthStatusSuccess
	ToTransferResponsesMonthStatusSuccess(Transfers []*record.TransferRecordMonthStatusSuccess) []*response.TransferResponseMonthStatusSuccess
	ToTransferResponseYearStatusSuccess(s *record.TransferRecordYearStatusSuccess) *response.TransferResponseYearStatusSuccess
	ToTransferResponsesYearStatusSuccess(Transfers []*record.TransferRecordYearStatusSuccess) []*response.TransferResponseYearStatusSuccess

	ToTransferResponseMonthStatusFailed(s *record.TransferRecordMonthStatusFailed) *response.TransferResponseMonthStatusFailed
	ToTransferResponsesMonthStatusFailed(Transfers []*record.TransferRecordMonthStatusFailed) []*response.TransferResponseMonthStatusFailed
	ToTransferResponseYearStatusFailed(s *record.TransferRecordYearStatusFailed) *response.TransferResponseYearStatusFailed
	ToTransferResponsesYearStatusFailed(Transfers []*record.TransferRecordYearStatusFailed) []*response.TransferResponseYearStatusFailed

	ToTransferResponseMonthAmount(s *record.TransferMonthAmount) *response.TransferMonthAmountResponse
	ToTransferResponsesMonthAmount(s []*record.TransferMonthAmount) []*response.TransferMonthAmountResponse

	ToTransferResponseYearAmount(s *record.TransferYearAmount) *response.TransferYearAmountResponse
	ToTransferResponsesYearAmount(s []*record.TransferYearAmount) []*response.TransferYearAmountResponse

	ToTransferResponseDeleteAt(transfer *record.TransferRecord) *response.TransferResponseDeleteAt
	ToTransfersResponseDeleteAt(transfers []*record.TransferRecord) []*response.TransferResponseDeleteAt
}

type WithdrawResponseMapper interface {
	ToWithdrawResponse(withdraw *record.WithdrawRecord) *response.WithdrawResponse
	ToWithdrawsResponse(withdraws []*record.WithdrawRecord) []*response.WithdrawResponse

	ToWithdrawResponseMonthStatusSuccess(s *record.WithdrawRecordMonthStatusSuccess) *response.WithdrawResponseMonthStatusSuccess
	ToWithdrawResponsesMonthStatusSuccess(Withdraws []*record.WithdrawRecordMonthStatusSuccess) []*response.WithdrawResponseMonthStatusSuccess
	ToWithdrawResponseYearStatusSuccess(s *record.WithdrawRecordYearStatusSuccess) *response.WithdrawResponseYearStatusSuccess
	ToWithdrawResponsesYearStatusSuccess(Withdraws []*record.WithdrawRecordYearStatusSuccess) []*response.WithdrawResponseYearStatusSuccess

	ToWithdrawResponseMonthStatusFailed(s *record.WithdrawRecordMonthStatusFailed) *response.WithdrawResponseMonthStatusFailed
	ToWithdrawResponsesMonthStatusFailed(Withdraws []*record.WithdrawRecordMonthStatusFailed) []*response.WithdrawResponseMonthStatusFailed
	ToWithdrawResponseYearStatusFailed(s *record.WithdrawRecordYearStatusFailed) *response.WithdrawResponseYearStatusFailed
	ToWithdrawResponsesYearStatusFailed(Withdraws []*record.WithdrawRecordYearStatusFailed) []*response.WithdrawResponseYearStatusFailed

	ToWithdrawAmountMonthlyResponse(s *record.WithdrawMonthlyAmount) *response.WithdrawMonthlyAmountResponse
	ToWithdrawsAmountMonthlyResponses(s []*record.WithdrawMonthlyAmount) []*response.WithdrawMonthlyAmountResponse

	ToWithdrawAmountYearlyResponse(s *record.WithdrawYearlyAmount) *response.WithdrawYearlyAmountResponse
	ToWithdrawsAmountYearlyResponses(s []*record.WithdrawYearlyAmount) []*response.WithdrawYearlyAmountResponse

	ToWithdrawResponseDeleteAt(withdraw *record.WithdrawRecord) *response.WithdrawResponseDeleteAt
	ToWithdrawsResponseDeleteAt(withdraws []*record.WithdrawRecord) []*response.WithdrawResponseDeleteAt
}

type MerchantResponseMapper interface {
	ToMerchantResponse(merchant *record.MerchantRecord) *response.MerchantResponse
	ToMerchantsResponse(merchants []*record.MerchantRecord) []*response.MerchantResponse

	ToMerchantMonthlyTotalAmount(ms *record.MerchantMonthlyTotalAmount) *response.MerchantResponseMonthlyTotalAmount
	ToMerchantMonthlyTotalAmounts(ms []*record.MerchantMonthlyTotalAmount) []*response.MerchantResponseMonthlyTotalAmount
	ToMerchantYearlyTotalAmount(ms *record.MerchantYearlyTotalAmount) *response.MerchantResponseYearlyTotalAmount
	ToMerchantYearlyTotalAmounts(ms []*record.MerchantYearlyTotalAmount) []*response.MerchantResponseYearlyTotalAmount

	ToMerchantTransactionResponse(merchant *record.MerchantTransactionsRecord) *response.MerchantTransactionResponse
	ToMerchantsTransactionResponse(merchants []*record.MerchantTransactionsRecord) []*response.MerchantTransactionResponse

	ToMerchantMonthlyPaymentMethod(ms *record.MerchantMonthlyPaymentMethod) *response.MerchantResponseMonthlyPaymentMethod
	ToMerchantMonthlyPaymentMethods(ms []*record.MerchantMonthlyPaymentMethod) []*response.MerchantResponseMonthlyPaymentMethod
	ToMerchantYearlyPaymentMethod(ms *record.MerchantYearlyPaymentMethod) *response.MerchantResponseYearlyPaymentMethod
	ToMerchantYearlyPaymentMethods(ms []*record.MerchantYearlyPaymentMethod) []*response.MerchantResponseYearlyPaymentMethod

	ToMerchantMonthlyAmount(ms *record.MerchantMonthlyAmount) *response.MerchantResponseMonthlyAmount
	ToMerchantMonthlyAmounts(ms []*record.MerchantMonthlyAmount) []*response.MerchantResponseMonthlyAmount
	ToMerchantYearlyAmount(ms *record.MerchantYearlyAmount) *response.MerchantResponseYearlyAmount
	ToMerchantYearlyAmounts(ms []*record.MerchantYearlyAmount) []*response.MerchantResponseYearlyAmount

	ToMerchantResponseDeleteAt(merchant *record.MerchantRecord) *response.MerchantResponseDeleteAt
	ToMerchantsResponseDeleteAt(merchants []*record.MerchantRecord) []*response.MerchantResponseDeleteAt
}
