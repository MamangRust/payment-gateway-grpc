package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	FindAllUsers(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByActive(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByTrashed(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindById(user_id int) (*record.UserRecord, error)
	FindByEmail(email string) (*record.UserRecord, error)
	CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error)
	TrashedUser(user_id int) (*record.UserRecord, error)
	RestoreUser(user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(user_id int) (bool, error)
	RestoreAllUser() (bool, error)
	DeleteAllUserPermanent() (bool, error)
}

type RoleRepository interface {
	FindAllRoles(req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)
	FindByActiveRole(req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)
	FindByTrashedRole(req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
	FindByUserId(user_id int) ([]*record.RoleRecord, error)
	CreateRole(request *requests.CreateRoleRequest) (*record.RoleRecord, error)
	UpdateRole(request *requests.UpdateRoleRequest) (*record.RoleRecord, error)
	TrashedRole(role_id int) (*record.RoleRecord, error)
	RestoreRole(role_id int) (*record.RoleRecord, error)
	DeleteRolePermanent(role_id int) (bool, error)
	RestoreAllRole() (bool, error)
	DeleteAllRolePermanent() (bool, error)
}

type RefreshTokenRepository interface {
	FindByToken(token string) (*record.RefreshTokenRecord, error)
	FindByUserId(user_id int) (*record.RefreshTokenRecord, error)
	CreateRefreshToken(req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error)
	UpdateRefreshToken(req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByUserId(user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error)
	RemoveRoleFromUser(req *requests.RemoveUserRoleRequest) error
}

type CardRepository interface {
	FindAllCards(req *requests.FindAllCards) ([]*record.CardRecord, *int, error)
	FindByActive(req *requests.FindAllCards) ([]*record.CardRecord, *int, error)
	FindByTrashed(req *requests.FindAllCards) ([]*record.CardRecord, *int, error)
	FindById(card_id int) (*record.CardRecord, error)
	FindCardByUserId(user_id int) (*record.CardRecord, error)
	FindCardByCardNumber(card_number string) (*record.CardRecord, error)

	GetTotalBalances() (*int64, error)
	GetTotalTopAmount() (*int64, error)
	GetTotalWithdrawAmount() (*int64, error)
	GetTotalTransactionAmount() (*int64, error)
	GetTotalTransferAmount() (*int64, error)

	GetTotalBalanceByCardNumber(cardNumber string) (*int64, error)
	GetTotalTopupAmountByCardNumber(cardNumber string) (*int64, error)
	GetTotalWithdrawAmountByCardNumber(cardNumber string) (*int64, error)
	GetTotalTransactionAmountByCardNumber(cardNumber string) (*int64, error)
	GetTotalTransferAmountBySender(senderCardNumber string) (*int64, error)
	GetTotalTransferAmountByReceiver(receiverCardNumber string) (*int64, error)

	GetMonthlyBalance(year int) ([]*record.CardMonthBalance, error)
	GetYearlyBalance(year int) ([]*record.CardYearlyBalance, error)
	GetMonthlyTopupAmount(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTopupAmount(year int) ([]*record.CardYearAmount, error)
	GetMonthlyWithdrawAmount(year int) ([]*record.CardMonthAmount, error)
	GetYearlyWithdrawAmount(year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransactionAmount(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransactionAmount(year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountSender(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountSender(year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountReceiver(year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountReceiver(year int) ([]*record.CardYearAmount, error)

	GetMonthlyBalancesByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthBalance, error)
	GetYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearlyBalance, error)
	GetMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)

	CreateCard(request *requests.CreateCardRequest) (*record.CardRecord, error)
	UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error)
	TrashedCard(cardId int) (*record.CardRecord, error)
	RestoreCard(cardId int) (*record.CardRecord, error)
	DeleteCardPermanent(card_id int) (bool, error)
	RestoreAllCard() (bool, error)
	DeleteAllCardPermanent() (bool, error)
}

type MerchantRepository interface {
	FindAllMerchants(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
	FindByActive(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error)

	FindById(merchant_id int) (*record.MerchantRecord, error)
	GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error)

	FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error)
	GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error)

	FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error)
	GetMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error)

	FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error)
	GetMonthlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error)

	FindByApiKey(api_key string) (*record.MerchantRecord, error)
	FindByName(name string) (*record.MerchantRecord, error)
	FindByMerchantUserId(user_id int) ([]*record.MerchantRecord, error)

	CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchantStatus(request *requests.UpdateMerchantStatus) (*record.MerchantRecord, error)

	TrashedMerchant(merchantId int) (*record.MerchantRecord, error)
	RestoreMerchant(merchant_id int) (*record.MerchantRecord, error)
	DeleteMerchantPermanent(merchant_id int) (bool, error)

	RestoreAllMerchant() (bool, error)
	DeleteAllMerchantPermanent() (bool, error)
}

type SaldoRepository interface {
	FindAllSaldos(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)
	FindByActive(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)
	FindByTrashed(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error)
	FindById(saldo_id int) (*record.SaldoRecord, error)

	GetMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*record.SaldoMonthTotalBalance, error)
	GetYearTotalSaldoBalance(year int) ([]*record.SaldoYearTotalBalance, error)
	GetMonthlySaldoBalances(year int) ([]*record.SaldoMonthSaldoBalance, error)
	GetYearlySaldoBalances(year int) ([]*record.SaldoYearSaldoBalance, error)

	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
	CreateSaldo(request *requests.CreateSaldoRequest) (*record.SaldoRecord, error)
	UpdateSaldo(request *requests.UpdateSaldoRequest) (*record.SaldoRecord, error)
	UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error)
	UpdateSaldoWithdraw(request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error)
	TrashedSaldo(saldoID int) (*record.SaldoRecord, error)
	RestoreSaldo(saldoID int) (*record.SaldoRecord, error)
	DeleteSaldoPermanent(saldo_id int) (bool, error)

	RestoreAllSaldo() (bool, error)
	DeleteAllSaldoPermanent() (bool, error)
}

type TopupRepository interface {
	FindAllTopups(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)
	FindByActive(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)
	FindByTrashed(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error)
	FindAllTopupByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*record.TopupRecord, *int, error)

	FindById(topup_id int) (*record.TopupRecord, error)

	GetMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusSuccess, error)
	GetYearlyTopupStatusSuccess(year int) ([]*record.TopupRecordYearStatusSuccess, error)

	GetMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusFailed, error)
	GetYearlyTopupStatusFailed(year int) ([]*record.TopupRecordYearStatusFailed, error)

	GetMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusSuccess, error)
	GetYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusSuccess, error)

	GetMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusFailed, error)
	GetYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusFailed, error)

	GetMonthlyTopupMethods(year int) ([]*record.TopupMonthMethod, error)
	GetYearlyTopupMethods(year int) ([]*record.TopupYearlyMethod, error)
	GetMonthlyTopupAmounts(year int) ([]*record.TopupMonthAmount, error)
	GetYearlyTopupAmounts(year int) ([]*record.TopupYearlyAmount, error)

	GetMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthMethod, error)
	GetYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyMethod, error)
	GetMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupMonthAmount, error)
	GetYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*record.TopupYearlyAmount, error)

	CreateTopup(request *requests.CreateTopupRequest) (*record.TopupRecord, error)
	UpdateTopup(request *requests.UpdateTopupRequest) (*record.TopupRecord, error)

	UpdateTopupAmount(request *requests.UpdateTopupAmount) (*record.TopupRecord, error)
	UpdateTopupStatus(request *requests.UpdateTopupStatus) (*record.TopupRecord, error)

	TrashedTopup(topup_id int) (*record.TopupRecord, error)
	RestoreTopup(topup_id int) (*record.TopupRecord, error)
	DeleteTopupPermanent(topup_id int) (bool, error)

	RestoreAllTopup() (bool, error)
	DeleteAllTopupPermanent() (bool, error)
}

type TransactionRepository interface {
	FindAllTransactions(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)
	FindByActive(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)
	FindByTrashed(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error)
	FindAllTransactionByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*record.TransactionRecord, *int, error)

	FindById(transaction_id int) (*record.TransactionRecord, error)

	GetMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusSuccess, error)
	GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error)
	GetMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusFailed, error)
	GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error)

	GetMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusSuccess, error)
	GetYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusSuccess, error)
	GetMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusFailed, error)
	GetYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusFailed, error)

	GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error)
	GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error)
	GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error)
	GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error)

	GetMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthMethod, error)
	GetYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error)
	GetMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error)
	GetYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error)

	FindTransactionByMerchantId(merchant_id int) ([]*record.TransactionRecord, error)

	CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransactionStatus(request *requests.UpdateTransactionStatus) (*record.TransactionRecord, error)
	TrashedTransaction(transaction_id int) (*record.TransactionRecord, error)
	RestoreTransaction(topup_id int) (*record.TransactionRecord, error)
	DeleteTransactionPermanent(topup_id int) (bool, error)

	RestoreAllTransaction() (bool, error)
	DeleteAllTransactionPermanent() (bool, error)
}

type TransferRepository interface {
	FindAll(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error)
	FindByActive(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error)
	FindByTrashed(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error)

	FindById(id int) (*record.TransferRecord, error)

	GetMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error)
	GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error)
	GetMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error)
	GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error)

	GetMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error)
	GetYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error)
	GetMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error)
	GetYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error)

	GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error)
	GetMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error)
	GetMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error)

	FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error)
	FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error)

	CreateTransfer(request *requests.CreateTransferRequest) (*record.TransferRecord, error)
	UpdateTransfer(request *requests.UpdateTransferRequest) (*record.TransferRecord, error)
	UpdateTransferAmount(request *requests.UpdateTransferAmountRequest) (*record.TransferRecord, error)
	UpdateTransferStatus(request *requests.UpdateTransferStatus) (*record.TransferRecord, error)

	TrashedTransfer(transfer_id int) (*record.TransferRecord, error)
	RestoreTransfer(transfer_id int) (*record.TransferRecord, error)
	DeleteTransferPermanent(topup_id int) (bool, error)

	RestoreAllTransfer() (bool, error)
	DeleteAllTransferPermanent() (bool, error)
}

type WithdrawRepository interface {
	FindAll(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)
	FindByActive(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)
	FindByTrashed(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error)
	FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*record.WithdrawRecord, *int, error)

	FindById(id int) (*record.WithdrawRecord, error)

	GetMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusSuccess, error)
	GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error)
	GetMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusFailed, error)
	GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error)

	GetMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusSuccess, error)
	GetYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusSuccess, error)
	GetMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusFailed, error)
	GetYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusFailed, error)

	GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error)
	GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error)
	GetMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawMonthlyAmount, error)
	GetYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error)

	CreateWithdraw(request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error)
	UpdateWithdraw(request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error)
	UpdateWithdrawStatus(request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error)

	TrashedWithdraw(WithdrawID int) (*record.WithdrawRecord, error)
	RestoreWithdraw(WithdrawID int) (*record.WithdrawRecord, error)
	DeleteWithdrawPermanent(WithdrawID int) (bool, error)

	RestoreAllWithdraw() (bool, error)
	DeleteAllWithdrawPermanent() (bool, error)
}
