package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersWithPaginationRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetActiveUsersWithPaginationRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetTrashedUsersWithPaginationRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
	FindByEmail(ctx context.Context, email string) (*db.GetUserByEmailRow, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*db.GetUserByEmailWithPasswordRow, error)
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.UpdateUserRow, error)
	TrashedUser(ctx context.Context, user_id int) (*db.User, error)
	RestoreUser(ctx context.Context, user_id int) (*db.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleRepository interface {
	FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, error)
	FindById(ctx context.Context, role_id int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
	FindByUserId(ctx context.Context, user_id int) ([]*db.Role, error)
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	TrashedRole(ctx context.Context, role_id int) (*db.Role, error)
	RestoreRole(ctx context.Context, role_id int) (*db.Role, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*db.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*db.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*db.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*db.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

type CardRepository interface {
	FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, error)
	FindById(ctx context.Context, card_id int) (*db.GetCardByIDRow, error)
	FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error)
	FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error)

	GetTotalBalances(ctx context.Context) (*int64, error)
	GetTotalTopAmount(ctx context.Context) (*int64, error)
	GetTotalWithdrawAmount(ctx context.Context) (*int64, error)
	GetTotalTransactionAmount(ctx context.Context) (*int64, error)
	GetTotalTransferAmount(ctx context.Context) (*int64, error)

	GetTotalBalanceByCardNumber(ctx context.Context, cardNumber string) (*int64, error)
	GetTotalTopupAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error)
	GetTotalWithdrawAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error)
	GetTotalTransactionAmountByCardNumber(ctx context.Context, cardNumber string) (*int64, error)
	GetTotalTransferAmountBySender(ctx context.Context, senderCardNumber string) (*int64, error)
	GetTotalTransferAmountByReceiver(ctx context.Context, receiverCardNumber string) (*int64, error)

	GetMonthlyBalance(ctx context.Context, year int) ([]*db.GetMonthlyBalancesRow, error)
	GetYearlyBalance(ctx context.Context, year int) ([]*db.GetYearlyBalancesRow, error)
	GetMonthlyTopupAmount(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountRow, error)
	GetYearlyTopupAmount(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountRow, error)
	GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawAmountRow, error)
	GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*db.GetYearlyWithdrawAmountRow, error)
	GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*db.GetMonthlyTransactionAmountRow, error)
	GetYearlyTransactionAmount(ctx context.Context, year int) ([]*db.GetYearlyTransactionAmountRow, error)
	GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountSenderRow, error)
	GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountSenderRow, error)
	GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountReceiverRow, error)
	GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountReceiverRow, error)

	GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyBalancesByCardNumberRow, error)
	GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyBalancesByCardNumberRow, error)
	GetMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTopupAmountByCardNumberRow, error)
	GetYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTopupAmountByCardNumberRow, error)
	GetMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyWithdrawAmountByCardNumberRow, error)
	GetYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyWithdrawAmountByCardNumberRow, error)
	GetMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransactionAmountByCardNumberRow, error)
	GetYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransactionAmountByCardNumberRow, error)
	GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountBySenderRow, error)
	GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountBySenderRow, error)
	GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetMonthlyTransferAmountByReceiverRow, error)
	GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*db.GetYearlyTransferAmountByReceiverRow, error)

	CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error)
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error)
	TrashedCard(ctx context.Context, cardId int) (*db.Card, error)
	RestoreCard(ctx context.Context, cardId int) (*db.Card, error)
	DeleteCardPermanent(ctx context.Context, cardId int) (bool, error)
	RestoreAllCard(ctx context.Context) (bool, error)
	DeleteAllCardPermanent(ctx context.Context) (bool, error)
}

type MerchantRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetActiveMerchantsRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetTrashedMerchantsRow, error)

	FindById(ctx context.Context, merchant_id int) (*db.GetMerchantByIDRow, error)
	FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*db.FindAllTransactionsRow, error)
	FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*db.FindAllTransactionsByMerchantRow, error)
	FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*db.FindAllTransactionsByApikeyRow, error)

	GetMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyTotalAmountMerchantRow, error)
	GetYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyTotalAmountMerchantRow, error)
	GetMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsMerchantRow, error)
	GetYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodMerchantRow, error)

	GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyAmountMerchantRow, error)
	GetYearlyAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyAmountMerchantRow, error)

	GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetMonthlyPaymentMethodByMerchantsRow, error)
	GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetYearlyPaymentMethodByMerchantsRow, error)
	GetMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetMonthlyAmountByMerchantsRow, error)
	GetYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetYearlyAmountByMerchantsRow, error)
	GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetMonthlyTotalAmountByMerchantRow, error)
	GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetYearlyTotalAmountByMerchantRow, error)

	GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetMonthlyPaymentMethodByApikeyRow, error)
	GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetYearlyPaymentMethodByApikeyRow, error)
	GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetMonthlyAmountByApikeyRow, error)
	GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetYearlyAmountByApikeyRow, error)
	GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetMonthlyTotalAmountByApikeyRow, error)
	GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetYearlyTotalAmountByApikeyRow, error)

	FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error)
	FindByName(ctx context.Context, name string) (*db.GetMerchantByNameRow, error)
	FindByMerchantUserId(ctx context.Context, user_id int) ([]*db.GetMerchantsByUserIDRow, error)

	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatus) (*db.UpdateMerchantStatusRow, error)

	TrashedMerchant(ctx context.Context, merchantId int) (*db.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantId int) (*db.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantId int) (bool, error)

	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type SaldoRepository interface {
	FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, error)
	FindById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, error)

	GetMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*db.GetMonthlyTotalSaldoBalanceRow, error)
	GetYearTotalSaldoBalance(ctx context.Context, year int) ([]*db.GetYearlyTotalSaldoBalancesRow, error)
	GetMonthlySaldoBalances(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, error)
	GetYearlySaldoBalances(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, error)

	FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error)
	CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*db.CreateSaldoRow, error)
	UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*db.UpdateSaldoRow, error)
	UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error)
	UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*db.UpdateSaldoWithdrawRow, error)
	TrashedSaldo(ctx context.Context, saldoID int) (*db.Saldo, error)
	RestoreSaldo(ctx context.Context, saldoID int) (*db.Saldo, error)
	DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error)

	RestoreAllSaldo(ctx context.Context) (bool, error)
	DeleteAllSaldoPermanent(ctx context.Context) (bool, error)
}

type TopupRepository interface {
	FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTopupsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetActiveTopupsRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTrashedTopupsRow, error)
	FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*db.GetTopupsByCardNumberRow, error)

	FindById(ctx context.Context, topup_id int) (*db.GetTopupByIDRow, error)

	GetMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusSuccessRow, error)
	GetYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusSuccessRow, error)

	GetMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*db.GetMonthTopupStatusFailedRow, error)
	GetYearlyTopupStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTopupStatusFailedRow, error)

	GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusSuccessCardNumberRow, error)
	GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusSuccessCardNumberRow, error)

	GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*db.GetMonthTopupStatusFailedCardNumberRow, error)
	GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*db.GetYearlyTopupStatusFailedCardNumberRow, error)

	GetMonthlyTopupMethods(ctx context.Context, year int) ([]*db.GetMonthlyTopupMethodsRow, error)
	GetYearlyTopupMethods(ctx context.Context, year int) ([]*db.GetYearlyTopupMethodsRow, error)
	GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTopupAmountsRow, error)
	GetYearlyTopupAmounts(ctx context.Context, year int) ([]*db.GetYearlyTopupAmountsRow, error)

	GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupMethodsByCardNumberRow, error)
	GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupMethodsByCardNumberRow, error)
	GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetMonthlyTopupAmountsByCardNumberRow, error)
	GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*db.GetYearlyTopupAmountsByCardNumberRow, error)

	CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*db.CreateTopupRow, error)
	UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*db.UpdateTopupRow, error)

	UpdateTopupAmount(ctx context.Context, request *requests.UpdateTopupAmount) (*db.UpdateTopupAmountRow, error)
	UpdateTopupStatus(ctx context.Context, request *requests.UpdateTopupStatus) (*db.UpdateTopupStatusRow, error)

	TrashedTopup(ctx context.Context, topup_id int) (*db.Topup, error)
	RestoreTopup(ctx context.Context, topup_id int) (*db.Topup, error)
	DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error)

	RestoreAllTopup(ctx context.Context) (bool, error)
	DeleteAllTopupPermanent(ctx context.Context) (bool, error)
}

type TransactionRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTransactionsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetActiveTransactionsRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*db.GetTrashedTransactionsRow, error)
	FindAllTransactionByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*db.GetTransactionsByCardNumberRow, error)
	FindById(ctx context.Context, transaction_id int) (*db.GetTransactionByIDRow, error)

	GetMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusSuccessRow, error)
	GetYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusSuccessRow, error)
	GetMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*db.GetMonthTransactionStatusFailedRow, error)
	GetYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionStatusFailedRow, error)

	GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusSuccessCardNumberRow, error)
	GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusSuccessCardNumberRow, error)
	GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*db.GetMonthTransactionStatusFailedCardNumberRow, error)
	GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*db.GetYearlyTransactionStatusFailedCardNumberRow, error)

	GetMonthlyPaymentMethods(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsRow, error)
	GetYearlyPaymentMethods(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodsRow, error)
	GetMonthlyAmounts(ctx context.Context, year int) ([]*db.GetMonthlyAmountsRow, error)
	GetYearlyAmounts(ctx context.Context, year int) ([]*db.GetYearlyAmountsRow, error)

	GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyPaymentMethodsByCardNumberRow, error)
	GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyPaymentMethodsByCardNumberRow, error)
	GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetMonthlyAmountsByCardNumberRow, error)
	GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*db.GetYearlyAmountsByCardNumberRow, error)

	FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*db.GetTransactionsByMerchantIDRow, error)

	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error)
	UpdateTransactionStatus(ctx context.Context, request *requests.UpdateTransactionStatus) (*db.UpdateTransactionStatusRow, error)
	TrashedTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error)
	RestoreTransaction(ctx context.Context, topup_id int) (*db.Transaction, error)
	DeleteTransactionPermanent(ctx context.Context, topup_id int) (bool, error)

	RestoreAllTransaction(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}

type TransferRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTransfersRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetActiveTransfersRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTranfers) ([]*db.GetTrashedTransfersRow, error)

	FindById(ctx context.Context, id int) (*db.GetTransferByIDRow, error)

	GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusSuccessRow, error)
	GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusSuccessRow, error)
	GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*db.GetMonthTransferStatusFailedRow, error)
	GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyTransferStatusFailedRow, error)

	GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusSuccessCardNumberRow, error)
	GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusSuccessCardNumberRow, error)
	GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*db.GetMonthTransferStatusFailedCardNumberRow, error)
	GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*db.GetYearlyTransferStatusFailedCardNumberRow, error)

	GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*db.GetMonthlyTransferAmountsRow, error)
	GetYearlyTransferAmounts(ctx context.Context, year int) ([]*db.GetYearlyTransferAmountsRow, error)
	GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsBySenderCardNumberRow, error)
	GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsBySenderCardNumberRow, error)
	GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetMonthlyTransferAmountsByReceiverCardNumberRow, error)
	GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*db.GetYearlyTransferAmountsByReceiverCardNumberRow, error)

	FindTransferByTransferFrom(ctx context.Context, transferFrom string) ([]*db.GetTransfersBySourceCardRow, error)
	FindTransferByTransferTo(ctx context.Context, transferTo string) ([]*db.GetTransfersByDestinationCardRow, error)

	CreateTransfer(ctx context.Context, request *requests.CreateTransferRequest) (*db.CreateTransferRow, error)
	UpdateTransfer(ctx context.Context, request *requests.UpdateTransferRequest) (*db.UpdateTransferRow, error)
	UpdateTransferAmount(ctx context.Context, request *requests.UpdateTransferAmountRequest) (*db.UpdateTransferAmountRow, error)
	UpdateTransferStatus(ctx context.Context, request *requests.UpdateTransferStatus) (*db.UpdateTransferStatusRow, error)

	TrashedTransfer(ctx context.Context, transferID int) (*db.Transfer, error)
	RestoreTransfer(ctx context.Context, transferID int) (*db.Transfer, error)
	DeleteTransferPermanent(ctx context.Context, transferID int) (bool, error)

	RestoreAllTransfer(ctx context.Context) (bool, error)
	DeleteAllTransferPermanent(ctx context.Context) (bool, error)
}

type WithdrawRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetWithdrawsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetActiveWithdrawsRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetTrashedWithdrawsRow, error)
	FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*db.GetWithdrawsByCardNumberRow, error)
	FindById(ctx context.Context, id int) (*db.GetWithdrawByIDRow, error)

	GetMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusSuccessRow, error)
	GetYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusSuccessRow, error)
	GetMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*db.GetMonthWithdrawStatusFailedRow, error)
	GetYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*db.GetYearlyWithdrawStatusFailedRow, error)

	GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusSuccessCardNumberRow, error)
	GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusSuccessCardNumberRow, error)
	GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*db.GetMonthWithdrawStatusFailedCardNumberRow, error)
	GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*db.GetYearlyWithdrawStatusFailedCardNumberRow, error)

	GetMonthlyWithdraws(ctx context.Context, year int) ([]*db.GetMonthlyWithdrawsRow, error)
	GetYearlyWithdraws(ctx context.Context, year int) ([]*db.GetYearlyWithdrawsRow, error)
	GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetMonthlyWithdrawsByCardNumberRow, error)
	GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*db.GetYearlyWithdrawsByCardNumberRow, error)

	CreateWithdraw(ctx context.Context, request *requests.CreateWithdrawRequest) (*db.CreateWithdrawRow, error)
	UpdateWithdraw(ctx context.Context, request *requests.UpdateWithdrawRequest) (*db.UpdateWithdrawRow, error)
	UpdateWithdrawStatus(ctx context.Context, request *requests.UpdateWithdrawStatus) (*db.UpdateWithdrawStatusRow, error)

	TrashedWithdraw(ctx context.Context, withdrawID int) (*db.Withdraw, error)
	RestoreWithdraw(ctx context.Context, withdrawID int) (*db.Withdraw, error)
	DeleteWithdrawPermanent(ctx context.Context, withdrawID int) (bool, error)

	RestoreAllWithdraw(ctx context.Context) (bool, error)
	DeleteAllWithdrawPermanent(ctx context.Context) (bool, error)
}
