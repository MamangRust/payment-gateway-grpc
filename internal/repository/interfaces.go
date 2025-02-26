package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	FindAllUsers(search string, page, pageSize int) ([]*record.UserRecord, int, error)
	FindById(user_id int) (*record.UserRecord, error)
	FindByEmail(email string) (*record.UserRecord, error)
	FindByActive(search string, page, pageSize int) ([]*record.UserRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.UserRecord, int, error)
	CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error)
	TrashedUser(user_id int) (*record.UserRecord, error)
	RestoreUser(user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(user_id int) (bool, error)
	RestoreAllUser() (bool, error)
	DeleteAllUserPermanent() (bool, error)
}

type RoleRepository interface {
	FindAllRoles(page int, pageSize int, search string) ([]*record.RoleRecord, int, error)
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
	FindByUserId(user_id int) ([]*record.RoleRecord, error)
	FindByActiveRole(page int, pageSize int, search string) ([]*record.RoleRecord, int, error)
	FindByTrashedRole(page int, pageSize int, search string) ([]*record.RoleRecord, int, error)
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
	FindAllCards(search string, page, pageSize int) ([]*record.CardRecord, int, error)
	FindById(card_id int) (*record.CardRecord, error)
	FindCardByUserId(user_id int) (*record.CardRecord, error)
	FindByActive(search string, page, pageSize int) ([]*record.CardRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.CardRecord, int, error)
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

	GetMonthlyBalancesByCardNumber(card_number string, year int) ([]*record.CardMonthBalance, error)
	GetYearlyBalanceByCardNumber(card_number string, year int) ([]*record.CardYearlyBalance, error)
	GetMonthlyTopupAmountByCardNumber(cardNumber string, year int) ([]*record.CardMonthAmount, error)
	GetYearlyTopupAmountByCardNumber(cardNumber string, year int) ([]*record.CardYearAmount, error)
	GetMonthlyWithdrawAmountByCardNumber(cardNumber string, year int) ([]*record.CardMonthAmount, error)
	GetYearlyWithdrawAmountByCardNumber(cardNumber string, year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransactionAmountByCardNumber(cardNumber string, year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransactionAmountByCardNumber(cardNumber string, year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountBySender(cardNumber string, year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountBySender(cardNumber string, year int) ([]*record.CardYearAmount, error)
	GetMonthlyTransferAmountByReceiver(cardNumber string, year int) ([]*record.CardMonthAmount, error)
	GetYearlyTransferAmountByReceiver(cardNumber string, year int) ([]*record.CardYearAmount, error)

	CreateCard(request *requests.CreateCardRequest) (*record.CardRecord, error)
	UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error)
	TrashedCard(cardId int) (*record.CardRecord, error)
	RestoreCard(cardId int) (*record.CardRecord, error)
	DeleteCardPermanent(card_id int) (bool, error)
	RestoreAllCard() (bool, error)
	DeleteAllCardPermanent() (bool, error)
}

type MerchantRepository interface {
	FindAllMerchants(search string, page, pageSize int) ([]*record.MerchantRecord, int, error)
	FindById(merchant_id int) (*record.MerchantRecord, error)
	GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error)

	FindAllTransactions(search string, page, pageSize int) ([]*record.MerchantTransactionsRecord, int, error)
	GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error)

	FindAllTransactionsByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.MerchantTransactionsRecord, int, error)
	GetMonthlyPaymentMethodByMerchants(merchantID int, year int) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodByMerchants(merchantID int, year int) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountByMerchants(merchantID int, year int) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountByMerchants(merchantID int, year int) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountByMerchants(merchantID int, year int) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountByMerchants(merchantID int, year int) ([]*record.MerchantYearlyTotalAmount, error)

	FindAllTransactionsByApikey(api_key string, search string, page, pageSize int) ([]*record.MerchantTransactionsRecord, int, error)
	GetMonthlyPaymentMethodByApikey(api_key string, year int) ([]*record.MerchantMonthlyPaymentMethod, error)
	GetYearlyPaymentMethodByApikey(api_key string, year int) ([]*record.MerchantYearlyPaymentMethod, error)
	GetMonthlyAmountByApikey(api_key string, year int) ([]*record.MerchantMonthlyAmount, error)
	GetYearlyAmountByApikey(api_key string, year int) ([]*record.MerchantYearlyAmount, error)
	GetMonthlyTotalAmountByApikey(api_key string, year int) ([]*record.MerchantMonthlyTotalAmount, error)
	GetYearlyTotalAmountByApikey(api_key string, year int) ([]*record.MerchantYearlyTotalAmount, error)

	FindByActive(search string, page, pageSize int) ([]*record.MerchantRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.MerchantRecord, int, error)
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
	FindAllSaldos(search string, page, pageSize int) ([]*record.SaldoRecord, int, error)
	FindById(saldo_id int) (*record.SaldoRecord, error)

	GetMonthlyTotalSaldoBalance(year int, month int) ([]*record.SaldoMonthTotalBalance, error)
	GetYearTotalSaldoBalance(year int) ([]*record.SaldoYearTotalBalance, error)
	GetMonthlySaldoBalances(year int) ([]*record.SaldoMonthSaldoBalance, error)
	GetYearlySaldoBalances(year int) ([]*record.SaldoYearSaldoBalance, error)

	FindByCardNumber(card_number string) (*record.SaldoRecord, error)
	FindByActive(search string, page, pageSize int) ([]*record.SaldoRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.SaldoRecord, int, error)
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
	FindAllTopups(search string, page, pageSize int) ([]*record.TopupRecord, int, error)
	FindAllTopupByCardNumber(card_number string, search string, page, pageSize int) ([]*record.TopupRecord, int, error)

	FindById(topup_id int) (*record.TopupRecord, error)

	GetMonthTopupStatusSuccess(year int, month int) ([]*record.TopupRecordMonthStatusSuccess, error)
	GetYearlyTopupStatusSuccess(year int) ([]*record.TopupRecordYearStatusSuccess, error)

	GetMonthTopupStatusFailed(year int, month int) ([]*record.TopupRecordMonthStatusFailed, error)
	GetYearlyTopupStatusFailed(year int) ([]*record.TopupRecordYearStatusFailed, error)

	GetMonthTopupStatusSuccessByCardNumber(card_number string, year int, month int) ([]*record.TopupRecordMonthStatusSuccess, error)
	GetYearlyTopupStatusSuccessByCardNumber(card_number string, year int) ([]*record.TopupRecordYearStatusSuccess, error)

	GetMonthTopupStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.TopupRecordMonthStatusFailed, error)
	GetYearlyTopupStatusFailedByCardNumber(card_number string, year int) ([]*record.TopupRecordYearStatusFailed, error)

	GetMonthlyTopupMethods(year int) ([]*record.TopupMonthMethod, error)
	GetYearlyTopupMethods(year int) ([]*record.TopupYearlyMethod, error)
	GetMonthlyTopupAmounts(year int) ([]*record.TopupMonthAmount, error)
	GetYearlyTopupAmounts(year int) ([]*record.TopupYearlyAmount, error)

	GetMonthlyTopupMethodsByCardNumber(card_number string, year int) ([]*record.TopupMonthMethod, error)
	GetYearlyTopupMethodsByCardNumber(card_number string, year int) ([]*record.TopupYearlyMethod, error)
	GetMonthlyTopupAmountsByCardNumber(card_number string, year int) ([]*record.TopupMonthAmount, error)
	GetYearlyTopupAmountsByCardNumber(card_number string, year int) ([]*record.TopupYearlyAmount, error)

	FindByActive(search string, page, pageSize int) ([]*record.TopupRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.TopupRecord, int, error)

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
	FindAllTransactions(search string, page, pageSize int) ([]*record.TransactionRecord, int, error)
	FindAllTransactionByCardNumber(card_number string, search string, page, pageSize int) ([]*record.TransactionRecord, int, error)

	FindById(transaction_id int) (*record.TransactionRecord, error)

	GetMonthTransactionStatusSuccess(year int, month int) ([]*record.TransactionRecordMonthStatusSuccess, error)
	GetYearlyTransactionStatusSuccess(year int) ([]*record.TransactionRecordYearStatusSuccess, error)
	GetMonthTransactionStatusFailed(year int, month int) ([]*record.TransactionRecordMonthStatusFailed, error)
	GetYearlyTransactionStatusFailed(year int) ([]*record.TransactionRecordYearStatusFailed, error)

	GetMonthTransactionStatusSuccessByCardNumber(card_number string, year int, month int) ([]*record.TransactionRecordMonthStatusSuccess, error)
	GetYearlyTransactionStatusSuccessByCardNumber(card_number string, year int) ([]*record.TransactionRecordYearStatusSuccess, error)
	GetMonthTransactionStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.TransactionRecordMonthStatusFailed, error)
	GetYearlyTransactionStatusFailedByCardNumber(card_number string, year int) ([]*record.TransactionRecordYearStatusFailed, error)

	GetMonthlyPaymentMethods(year int) ([]*record.TransactionMonthMethod, error)
	GetYearlyPaymentMethods(year int) ([]*record.TransactionYearMethod, error)
	GetMonthlyAmounts(year int) ([]*record.TransactionMonthAmount, error)
	GetYearlyAmounts(year int) ([]*record.TransactionYearlyAmount, error)

	GetMonthlyPaymentMethodsByCardNumber(card_number string, year int) ([]*record.TransactionMonthMethod, error)
	GetYearlyPaymentMethodsByCardNumber(card_number string, year int) ([]*record.TransactionYearMethod, error)
	GetMonthlyAmountsByCardNumber(card_number string, year int) ([]*record.TransactionMonthAmount, error)
	GetYearlyAmountsByCardNumber(card_number string, year int) ([]*record.TransactionYearlyAmount, error)

	FindByActive(search string, page, pageSize int) ([]*record.TransactionRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.TransactionRecord, int, error)

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
	FindAll(search string, page, pageSize int) ([]*record.TransferRecord, int, error)
	FindById(id int) (*record.TransferRecord, error)

	GetMonthTransferStatusSuccess(year int, month int) ([]*record.TransferRecordMonthStatusSuccess, error)
	GetYearlyTransferStatusSuccess(year int) ([]*record.TransferRecordYearStatusSuccess, error)
	GetMonthTransferStatusFailed(year int, month int) ([]*record.TransferRecordMonthStatusFailed, error)
	GetYearlyTransferStatusFailed(year int) ([]*record.TransferRecordYearStatusFailed, error)

	GetMonthTransferStatusSuccessByCardNumber(cardNumber string, year int, month int) ([]*record.TransferRecordMonthStatusSuccess, error)
	GetYearlyTransferStatusSuccessByCardNumber(cardNumber string, year int) ([]*record.TransferRecordYearStatusSuccess, error)
	GetMonthTransferStatusFailedByCardNumber(cardNumber string, year int, month int) ([]*record.TransferRecordMonthStatusFailed, error)
	GetYearlyTransferStatusFailedByCardNumber(cardNumber string, year int) ([]*record.TransferRecordYearStatusFailed, error)

	GetMonthlyTransferAmounts(year int) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmounts(year int) ([]*record.TransferYearAmount, error)
	GetMonthlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*record.TransferMonthAmount, error)
	GetMonthlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*record.TransferMonthAmount, error)
	GetYearlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*record.TransferYearAmount, error)
	GetYearlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*record.TransferYearAmount, error)

	FindByActive(search string, page, pageSize int) ([]*record.TransferRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.TransferRecord, int, error)
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
	FindAll(search string, page, pageSize int) ([]*record.WithdrawRecord, int, error)
	FindAllByCardNumber(card_number string, search string, page, pageSize int) ([]*record.WithdrawRecord, int, error)
	FindById(id int) (*record.WithdrawRecord, error)

	GetMonthWithdrawStatusSuccess(year int, month int) ([]*record.WithdrawRecordMonthStatusSuccess, error)
	GetYearlyWithdrawStatusSuccess(year int) ([]*record.WithdrawRecordYearStatusSuccess, error)
	GetMonthWithdrawStatusFailed(year int, month int) ([]*record.WithdrawRecordMonthStatusFailed, error)
	GetYearlyWithdrawStatusFailed(year int) ([]*record.WithdrawRecordYearStatusFailed, error)

	GetMonthWithdrawStatusSuccessByCardNumber(card_number string, year int, month int) ([]*record.WithdrawRecordMonthStatusSuccess, error)
	GetYearlyWithdrawStatusSuccessByCardNumber(card_number string, year int) ([]*record.WithdrawRecordYearStatusSuccess, error)
	GetMonthWithdrawStatusFailedByCardNumber(card_number string, year int, month int) ([]*record.WithdrawRecordMonthStatusFailed, error)
	GetYearlyWithdrawStatusFailedByCardNumber(card_number string, year int) ([]*record.WithdrawRecordYearStatusFailed, error)

	GetMonthlyWithdraws(year int) ([]*record.WithdrawMonthlyAmount, error)
	GetYearlyWithdraws(year int) ([]*record.WithdrawYearlyAmount, error)
	GetMonthlyWithdrawsByCardNumber(cardNumber string, year int) ([]*record.WithdrawMonthlyAmount, error)
	GetYearlyWithdrawsByCardNumber(cardNumber string, year int) ([]*record.WithdrawYearlyAmount, error)

	FindByActive(search string, page, pageSize int) ([]*record.WithdrawRecord, int, error)
	FindByTrashed(search string, page, pageSize int) ([]*record.WithdrawRecord, int, error)

	CreateWithdraw(request *requests.CreateWithdrawRequest) (*record.WithdrawRecord, error)
	UpdateWithdraw(request *requests.UpdateWithdrawRequest) (*record.WithdrawRecord, error)
	UpdateWithdrawStatus(request *requests.UpdateWithdrawStatus) (*record.WithdrawRecord, error)

	TrashedWithdraw(WithdrawID int) (*record.WithdrawRecord, error)
	RestoreWithdraw(WithdrawID int) (*record.WithdrawRecord, error)
	DeleteWithdrawPermanent(WithdrawID int) (bool, error)

	RestoreAllWithdraw() (bool, error)
	DeleteAllWithdrawPermanent() (bool, error)
}
