package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type AuthService interface {
	Register(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, error)
	RefreshToken(ctx context.Context, token string) (*response.TokenResponse, error)
	GetMe(ctx context.Context, token string) (*db.GetUserByIDRow, error)
}

type UserService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersWithPaginationRow, *int, error)
	FindByID(ctx context.Context, id int) (*db.GetUserByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetActiveUsersWithPaginationRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetTrashedUsersWithPaginationRow, *int, error)
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.UpdateUserRow, error)
	TrashedUser(ctx context.Context, user_id int) (*db.User, error)
	RestoreUser(ctx context.Context, user_id int) (*db.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)

	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleService interface {
	FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, *int, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, *int, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, *int, error)
	FindById(ctx context.Context, role_id int) (*db.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*db.Role, error)
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	TrashedRole(ctx context.Context, role_id int) (*db.Role, error)
	RestoreRole(ctx context.Context, role_id int) (*db.Role, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)

	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}

type CardService interface {
	FindAll(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, *int, error)
	FindById(ctx context.Context, card_id int) (*db.GetCardByIDRow, error)
	FindByUserID(ctx context.Context, userID int) (*db.GetCardByUserIDRow, error)
	FindByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error)

	DashboardCard(ctx context.Context) (*response.DashboardCard, error)
	DashboardCardCardNumber(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, error)

	FindMonthlyBalance(ctx context.Context, year int) ([]*db.GetMonthlyBalancesRow, error)
	FindYearlyBalance(ctx context.Context, year int) ([]*db.GetYearlyBalancesRow, error)
	FindMonthlyTopupAmount(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountRow, error)
	FindYearlyTopupAmount(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountRow, error)
	FindMonthlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawAmountRow, error)
	FindYearlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetYearlyWithdrawAmountRow, error)
	FindMonthlyTransactionAmount(ctx context.Context, year int) ([]*db.GetMonthlyTransactionAmountRow, error)
	FindYearlyTransactionAmount(ctx context.Context, year int) ([]*db.GetYearlyTransactionAmountRow, error)
	FindMonthlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountSenderRow, error)
	FindYearlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountSenderRow, error)
	FindMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountReceiverRow, error)
	FindYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountReceiverRow, error)

	FindMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyBalancesByCardNumberRow, error)
	FindYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyBalancesByCardNumberRow, error)
	FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTopupAmountByCardNumberRow, error)
	FindYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTopupAmountByCardNumberRow, error)
	FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyWithdrawAmountByCardNumberRow, error)
	FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyWithdrawAmountByCardNumberRow, error)
	FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransactionAmountByCardNumberRow, error)
	FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransactionAmountByCardNumberRow, error)
	FindMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountBySenderRow, error)
	FindYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountBySenderRow, error)
	FindMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountByReceiverRow, error)
	FindYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountByReceiverRow, error)

	CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error)
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error)
	TrashedCard(ctx context.Context, cardId int) (*db.Card, error)
	RestoreCard(ctx context.Context, cardId int) (*db.Card, error)
	DeleteCardPermanent(ctx context.Context, cardId int) (bool, error)

	RestoreAllCard(ctx context.Context) (bool, error)
	DeleteAllCardPermanent(ctx context.Context) (bool, error)
}

type MerchantService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error)
	FindById(ctx context.Context, merchant_id int) (*db.GetMerchantByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetActiveMerchantsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetTrashedMerchantsRow, *int, error)
	FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error)
	FindByMerchantUserId(ctx context.Context, user_id int) ([]*db.GetMerchantsByUserIDRow, error)

	FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*db.FindAllTransactionsRow, *int, error)
	FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*db.FindAllTransactionsByApikeyRow, *int, error)
	FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*db.FindAllTransactionsByMerchantRow, *int, error)

	FindMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyTotalAmountMerchantRow, error)
	FindYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyTotalAmountMerchantRow, error)
	FindMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsMerchantRow, error)
	FindYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodMerchantRow, error)
	FindMonthlyAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyAmountMerchantRow, error)
	FindYearlyAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyAmountMerchantRow, error)

	FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetMonthlyPaymentMethodByMerchantsRow, error)
	FindYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetYearlyPaymentMethodByMerchantsRow, error)
	FindMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetMonthlyAmountByMerchantsRow, error)
	FindYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetYearlyAmountByMerchantsRow, error)
	FindMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetMonthlyTotalAmountByMerchantRow, error)
	FindYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetYearlyTotalAmountByMerchantRow, error)

	FindMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetMonthlyPaymentMethodByApikeyRow, error)
	FindYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetYearlyPaymentMethodByApikeyRow, error)
	FindMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetMonthlyAmountByApikeyRow, error)
	FindYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetYearlyAmountByApikeyRow, error)
	FindMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetMonthlyTotalAmountByApikeyRow, error)
	FindYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetYearlyTotalAmountByApikeyRow, error)

	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error)
	TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error)
	RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)

	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type SaldoService interface {
	FindAll(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, *int, error)
	FindById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, error)
	FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error)
	FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, *int, error)

	FindMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*db.GetMonthlyTotalSaldoBalanceRow, error)
	FindYearTotalSaldoBalance(ctx context.Context, year int) ([]*db.GetYearlyTotalSaldoBalancesRow, error)
	FindMonthlySaldoBalances(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, error)
	FindYearlySaldoBalances(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, error)

	CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*db.CreateSaldoRow, error)
	UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*db.UpdateSaldoRow, error)
	TrashSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error)
	RestoreSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error)
	DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error)

	RestoreAllSaldo(ctx context.Context) (bool, error)
	DeleteAllSaldoPermanent(ctx context.Context) (bool, error)
}

type TopupService interface {
	FindAll(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTopupsRow, *int, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*db.GetTopupsByCardNumberRow, *int, error)
	FindById(ctx context.Context, topupID int) (*db.GetTopupByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetActiveTopupsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTrashedTopupsRow, *int, error)

	FindMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusSuccessRow, error)
	FindYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusSuccessRow, error)
	FindMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusFailedRow, error)
	FindYearlyTopupStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusFailedRow, error)

	FindMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusSuccessCardNumberRow, error)
	FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusSuccessCardNumberRow, error)
	FindMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusFailedCardNumberRow, error)
	FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusFailedCardNumberRow, error)

	FindMonthlyTopupMethods(ctx context.Context, year int) ([]*db.GetMonthlyTopupMethodsRow, error)
	FindYearlyTopupMethods(ctx context.Context, year int) ([]*db.GetYearlyTopupMethodsRow, error)
	FindMonthlyTopupAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountsRow, error)
	FindYearlyTopupAmounts(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountsRow, error)

	FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupMethodsByCardNumberRow, error)
	FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupMethodsByCardNumberRow, error)
	FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupAmountsByCardNumberRow, error)
	FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupAmountsByCardNumberRow, error)

	CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*db.CreateTopupRow, error)
	UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*db.UpdateTopupRow, error)
	TrashedTopup(ctx context.Context, topup_id int) (*db.Topup, error)
	RestoreTopup(ctx context.Context, topup_id int) (*db.Topup, error)
	DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error)

	RestoreAllTopup(ctx context.Context) (bool, error)
	DeleteAllTopupPermanent(ctx context.Context) (bool, error)
}

type TransactionService interface {
	FindAll(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTransactionsRow, *int, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*db.GetTransactionsByCardNumberRow, *int, error)
	FindById(ctx context.Context, transactionID int) (*db.GetTransactionByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetActiveTransactionsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTrashedTransactionsRow, *int, error)
	FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*db.GetTransactionsByMerchantIDRow, error)

	FindMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusSuccessRow, error)
	FindYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusSuccessRow, error)
	FindMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusFailedRow, error)
	FindYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusFailedRow, error)

	FindMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusSuccessCardNumberRow, error)
	FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusSuccessCardNumberRow, error)
	FindMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusFailedCardNumberRow, error)
	FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusFailedCardNumberRow, error)

	FindMonthlyPaymentMethods(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsRow, error)
	FindYearlyPaymentMethods(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodsRow, error)
	FindMonthlyAmounts(ctx context.Context, year int) ([]*db.GetMonthlyAmountsRow, error)
	FindYearlyAmounts(ctx context.Context, year int) ([]*db.GetYearlyAmountsRow, error)

	FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyPaymentMethodsByCardNumberRow, error)
	FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyPaymentMethodsByCardNumberRow, error)
	FindMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyAmountsByCardNumberRow, error)
	FindYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyAmountsByCardNumberRow, error)

	Create(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error)
	Update(ctx context.Context, apiKey string, request *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error)
	TrashedTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error)
	DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error)

	RestoreAllTransaction(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}

type TransferService interface {
	FindAll(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTransfersRow, *int, error)
	FindById(ctx context.Context, transferId int) (*db.GetTransferByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetActiveTransfersRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTrashedTransfersRow, *int, error)
	FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*db.GetTransfersBySourceCardRow, error)
	FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*db.GetTransfersByDestinationCardRow, error)

	FindMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusSuccessRow, error)
	FindYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusSuccessRow, error)
	FindMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusFailedRow, error)
	FindYearlyTransferStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusFailedRow, error)
	FindMonthlyTransferAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountsRow, error)
	FindYearlyTransferAmounts(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountsRow, error)

	FindMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusSuccessCardNumberRow, error)
	FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusSuccessCardNumberRow, error)
	FindMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusFailedCardNumberRow, error)
	FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusFailedCardNumberRow, error)

	FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsBySenderCardNumberRow, error)
	FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsBySenderCardNumberRow, error)
	FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsByReceiverCardNumberRow, error)
	FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsByReceiverCardNumberRow, error)

	CreateTransaction(ctx context.Context, request *requests.CreateTransferRequest) (*db.CreateTransferRow, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransferRequest) (*db.UpdateTransferRow, error)
	TrashedTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error)
	RestoreTransfer(ctx context.Context, transfer_id int) (*db.Transfer, error)
	DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, error)

	RestoreAllTransfer(ctx context.Context) (bool, error)
	DeleteAllTransferPermanent(ctx context.Context) (bool, error)
}

type WithdrawService interface {
	FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetWithdrawsRow, *int, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*db.GetWithdrawsByCardNumberRow, *int, error)
	FindById(ctx context.Context, withdrawID int) (*db.GetWithdrawByIDRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetActiveWithdrawsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetTrashedWithdrawsRow, *int, error)

	FindMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusSuccessRow, error)
	FindYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusSuccessRow, error)
	FindMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusFailedRow, error)
	FindYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusFailedRow, error)
	FindMonthlyWithdraws(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawsRow, error)
	FindYearlyWithdraws(ctx context.Context, year int) ([]*db.GetYearlyWithdrawsRow, error)

	FindMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusSuccessCardNumberRow, error)
	FindYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusSuccessCardNumberRow, error)
	FindMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusFailedCardNumberRow, error)
	FindYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusFailedCardNumberRow, error)

	FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, error)
	FindYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, error)

	Create(ctx context.Context, request *requests.CreateWithdrawRequest) (*db.CreateWithdrawRow, error)
	Update(ctx context.Context, request *requests.UpdateWithdrawRequest) (*db.UpdateWithdrawRow, error)
	TrashedWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error)
	RestoreWithdraw(ctx context.Context, withdraw_id int) (*db.Withdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, error)

	RestoreAllWithdraw(ctx context.Context) (bool, error)
	DeleteAllWithdrawPermanent(ctx context.Context) (bool, error)
}
