package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type AuthService interface {
	Register(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	Login(request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse)
	RefreshToken(token string) (*response.TokenResponse, *response.ErrorResponse)
	GetMe(token string) (*response.UserResponse, *response.ErrorResponse)
}

type RoleService interface {
	FindAll(req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	FindByActiveRole(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashedRole(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	FindByUserId(id int) ([]*response.RoleResponse, *response.ErrorResponse)
	CreateRole(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	TrashedRole(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	RestoreRole(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	DeleteRolePermanent(role_id int) (bool, *response.ErrorResponse)

	RestoreAllRole() (bool, *response.ErrorResponse)
	DeleteAllRolePermanent() (bool, *response.ErrorResponse)
}

type CardService interface {
	FindAll(req *requests.FindAllCards) ([]*response.CardResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(card_id int) (*response.CardResponse, *response.ErrorResponse)
	FindByUserID(userID int) (*response.CardResponse, *response.ErrorResponse)
	FindByCardNumber(card_number string) (*response.CardResponse, *response.ErrorResponse)

	DashboardCard() (*response.DashboardCard, *response.ErrorResponse)
	DashboardCardCardNumber(cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse)

	FindMonthlyBalance(year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)
	FindYearlyBalance(year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)
	FindMonthlyTopupAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTopupAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyWithdrawAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyWithdrawAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransactionAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransactionAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransferAmountSender(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountSender(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransferAmountReceiver(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountReceiver(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	FindMonthlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)
	FindYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)
	FindMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	CreateCard(request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse)
	UpdateCard(request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse)
	TrashedCard(cardId int) (*response.CardResponse, *response.ErrorResponse)
	RestoreCard(cardId int) (*response.CardResponse, *response.ErrorResponse)
	DeleteCardPermanent(cardId int) (bool, *response.ErrorResponse)

	RestoreAllCard() (bool, *response.ErrorResponse)
	DeleteAllCardPermanent() (bool, *response.ErrorResponse)
}

type MerchantService interface {
	FindAll(req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)

	FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)

	FindMonthlyPaymentMethodsMerchant(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodMerchant(year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountMerchant(year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountMerchant(year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)

	FindMonthlyTotalAmountMerchant(year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)

	FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
	FindMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
	FindMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)

	FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse)
	FindMonthlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
	FindMonthlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)

	FindByActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByApiKey(api_key string) (*response.MerchantResponse, *response.ErrorResponse)
	FindByMerchantUserId(user_id int) ([]*response.MerchantResponse, *response.ErrorResponse)
	CreateMerchant(request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	TrashedMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)
	RestoreMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)
	DeleteMerchantPermanent(merchant_id int) (bool, *response.ErrorResponse)

	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type SaldoService interface {
	FindAll(req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, *response.ErrorResponse)
	FindById(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)

	FindMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse)
	FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse)
	FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse)
	FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse)

	FindByCardNumber(card_number string) (*response.SaldoResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse)
	CreateSaldo(request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)
	UpdateSaldo(request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)
	TrashSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	RestoreSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	DeleteSaldoPermanent(saldo_id int) (bool, *response.ErrorResponse)

	RestoreAllSaldo() (bool, *response.ErrorResponse)
	DeleteAllSaldoPermanent() (bool, *response.ErrorResponse)
}

type TopupService interface {
	FindAll(req *requests.FindAllTopups) ([]*response.TopupResponse, *int, *response.ErrorResponse)
	FindAllByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, *response.ErrorResponse)

	FindById(topupID int) (*response.TopupResponse, *response.ErrorResponse)

	FindMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTopupStatusSuccess(year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTopupStatusFailed(year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTopupMethods(year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	FindYearlyTopupMethods(year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyTopupAmounts(year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTopupAmounts(year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)

	FindMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	FindYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)

	FindByActive(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse)
	CreateTopup(request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)
	UpdateTopup(request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)
	TrashedTopup(topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	RestoreTopup(topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse)
	DeleteTopupPermanent(topup_id int) (bool, *response.ErrorResponse)

	RestoreAllTopup() (bool, *response.ErrorResponse)
	DeleteAllTopupPermanent() (bool, *response.ErrorResponse)
}

type TransactionService interface {
	FindAll(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindAllByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse)

	FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse)

	FindMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransactionStatusSuccess(year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransactionStatusFailed(year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyPaymentMethods(year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	FindYearlyPaymentMethods(year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
	FindMonthlyAmounts(year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	FindYearlyAmounts(year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)

	FindMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	FindYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
	FindMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	FindYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)

	FindByActive(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindTransactionByMerchantId(merchant_id int) ([]*response.TransactionResponse, *response.ErrorResponse)
	Create(apiKey string, request *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	Update(apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	TrashedTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse)
	RestoreTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse)
	DeleteTransactionPermanent(transaction_id int) (bool, *response.ErrorResponse)

	RestoreAllTransaction() (bool, *response.ErrorResponse)
	DeleteAllTransactionPermanent() (bool, *response.ErrorResponse)
}

type TransferService interface {
	FindAll(req *requests.FindAllTranfers) ([]*response.TransferResponse, *int, *response.ErrorResponse)
	FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse)

	FindMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
	FindMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)

	FindByActive(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse)
	FindTransferByTransferFrom(transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse)
	FindTransferByTransferTo(transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse)
	CreateTransaction(request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse)
	UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse)
	TrashedTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse)
	RestoreTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse)
	DeleteTransferPermanent(transfer_id int) (bool, *response.ErrorResponse)

	RestoreAllTransfer() (bool, *response.ErrorResponse)
	DeleteAllTransferPermanent() (bool, *response.ErrorResponse)
}

type UserService interface {
	FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse)
	FindByID(id int) (*response.UserResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	TrashedUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	RestoreUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	DeleteUserPermanent(user_id int) (bool, *response.ErrorResponse)

	RestoreAllUser() (bool, *response.ErrorResponse)
	DeleteAllUserPermanent() (bool, *response.ErrorResponse)
}

type WithdrawService interface {
	FindAll(req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)
	FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse)

	FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse)

	FindMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyWithdrawStatusSuccess(year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyWithdrawStatusFailed(year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	FindYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
	FindMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	FindYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)

	FindByActive(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse)
	Create(request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)
	Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)
	TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)
	RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)
	DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse)

	RestoreAllWithdraw() (bool, *response.ErrorResponse)
	DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse)
}
