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
	FindAll(page int, pageSize int, search string) ([]*response.RoleResponse, int, *response.ErrorResponse)
	FindByActiveRole(page int, pageSize int, search string) ([]*response.RoleResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashedRole(page int, pageSize int, search string) ([]*response.RoleResponseDeleteAt, int, *response.ErrorResponse)
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
	FindAll(page int, pageSize int, search string) ([]*response.CardResponse, int, *response.ErrorResponse)
	FindById(card_id int) (*response.CardResponse, *response.ErrorResponse)
	FindByUserID(userID int) (*response.CardResponse, *response.ErrorResponse)
	FindByActive(page int, pageSize int, search string) ([]*response.CardResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.CardResponseDeleteAt, int, *response.ErrorResponse)
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

	FindMonthlyBalanceByCardNumber(card_number string, year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)
	FindYearlyBalanceByCardNumber(card_number string, year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)
	FindMonthlyTopupAmountByCardNumber(cardNumber string, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTopupAmountByCardNumber(cardNumber string, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyWithdrawAmountByCardNumber(cardNumber string, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyWithdrawAmountByCardNumber(cardNumber string, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransactionAmountByCardNumber(cardNumber string, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransactionAmountByCardNumber(cardNumber string, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransferAmountBySender(cardNumber string, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountBySender(cardNumber string, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
	FindMonthlyTransferAmountByReceiver(cardNumber string, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)
	FindYearlyTransferAmountByReceiver(cardNumber string, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	CreateCard(request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse)
	UpdateCard(request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse)
	TrashedCard(cardId int) (*response.CardResponse, *response.ErrorResponse)
	RestoreCard(cardId int) (*response.CardResponse, *response.ErrorResponse)
	DeleteCardPermanent(cardId int) (bool, *response.ErrorResponse)

	RestoreAllCard() (bool, *response.ErrorResponse)
	DeleteAllCardPermanent() (bool, *response.ErrorResponse)
}

type MerchantService interface {
	FindAll(page int, pageSize int, search string) ([]*response.MerchantResponse, int, *response.ErrorResponse)
	FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse)

	FindAllTransactions(page int, pageSize int, search string) ([]*response.MerchantTransactionResponse, int, *response.ErrorResponse)

	FindMonthlyPaymentMethodsMerchant(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodMerchant(year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountMerchant(year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountMerchant(year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)

	FindMonthlyTotalAmountMerchant(year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)

	FindAllTransactionsByMerchant(merchant_id int, page int, pageSize int, search string) ([]*response.MerchantTransactionResponse, int, *response.ErrorResponse)
	FindMonthlyPaymentMethodByMerchants(merchantID int, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodByMerchants(merchantID int, year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
	FindMonthlyTotalAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountByMerchants(merchantID int, year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)

	FindAllTransactionsByApikey(api_key string, page int, pageSize int, search string) ([]*response.MerchantTransactionResponse, int, *response.ErrorResponse)
	FindMonthlyPaymentMethodByApikeys(api_key string, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)
	FindYearlyPaymentMethodByApikeys(api_key string, year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
	FindMonthlyAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)
	FindYearlyAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
	FindMonthlyTotalAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)
	FindYearlyTotalAmountByApikeys(api_key string, year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)

	FindByActive(page int, pageSize int, search string) ([]*response.MerchantResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.MerchantResponseDeleteAt, int, *response.ErrorResponse)
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
	FindAll(page int, pageSize int, search string) ([]*response.SaldoResponse, int, *response.ErrorResponse)
	FindById(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)

	FindMonthlyTotalSaldoBalance(year int, month int) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse)
	FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse)
	FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse)
	FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse)

	FindByCardNumber(card_number string) (*response.SaldoResponse, *response.ErrorResponse)
	FindByActive(page int, pageSize int, search string) ([]*response.SaldoResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.SaldoResponseDeleteAt, int, *response.ErrorResponse)
	CreateSaldo(request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)
	UpdateSaldo(request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse)
	TrashSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	RestoreSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse)
	DeleteSaldoPermanent(saldo_id int) (bool, *response.ErrorResponse)

	RestoreAllSaldo() (bool, *response.ErrorResponse)
	DeleteAllSaldoPermanent() (bool, *response.ErrorResponse)
}

type TopupService interface {
	FindAll(page int, pageSize int, search string) ([]*response.TopupResponse, int, *response.ErrorResponse)
	FindAllByCardNumber(card_number string, page int, pageSize int, search string) ([]*response.TopupResponse, int, *response.ErrorResponse)

	FindById(topupID int) (*response.TopupResponse, *response.ErrorResponse)

	FindMonthTopupStatusSuccess(year int, month int) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTopupStatusSuccess(year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTopupStatusFailed(year int, month int) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTopupStatusFailed(year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthTopupStatusSuccessByCardNumber(card_number string, year int, month int) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTopupStatusSuccessByCardNumber(card_number string, year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTopupStatusFailedByCardNumber(card_number string, year int, month int) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTopupStatusFailedByCardNumber(card_number string, year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTopupMethods(year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	FindYearlyTopupMethods(year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyTopupAmounts(year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTopupAmounts(year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)

	FindMonthlyTopupMethodsByCardNumber(cardNumber string, year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)
	FindYearlyTopupMethodsByCardNumber(cardNumber string, year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
	FindMonthlyTopupAmountsByCardNumber(cardNumber string, year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTopupAmountsByCardNumber(cardNumber string, year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)

	FindByActive(page int, pageSize int, search string) ([]*response.TopupResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.TopupResponseDeleteAt, int, *response.ErrorResponse)
	CreateTopup(request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)
	UpdateTopup(request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse)
	TrashedTopup(topup_id int) (*response.TopupResponse, *response.ErrorResponse)
	RestoreTopup(topup_id int) (*response.TopupResponse, *response.ErrorResponse)
	DeleteTopupPermanent(topup_id int) (bool, *response.ErrorResponse)

	RestoreAllTopup() (bool, *response.ErrorResponse)
	DeleteAllTopupPermanent() (bool, *response.ErrorResponse)
}

type TransactionService interface {
	FindAll(page int, pageSize int, search string) ([]*response.TransactionResponse, int, *response.ErrorResponse)
	FindAllByCardNumber(card_number string, page int, pageSize int, search string) ([]*response.TransactionResponse, int, *response.ErrorResponse)

	FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse)

	FindMonthTransactionStatusSuccess(year int, month int) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransactionStatusSuccess(year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransactionStatusFailed(year int, month int) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransactionStatusFailed(year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthTransactionStatusSuccessByCardNumber(card_number string, year int, month int) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransactionStatusSuccessByCardNumber(card_number string, year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransactionStatusFailedByCardNumber(card_number string, year int, month int) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransactionStatusFailedByCardNumber(card_number string, year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyPaymentMethods(year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	FindYearlyPaymentMethods(year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
	FindMonthlyAmounts(year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	FindYearlyAmounts(year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)

	FindMonthlyPaymentMethodsByCardNumber(cardNumber string, year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)
	FindYearlyPaymentMethodsByCardNumber(cardNumber string, year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
	FindMonthlyAmountsByCardNumber(cardNumber string, year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)
	FindYearlyAmountsByCardNumber(cardNumber string, year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)

	FindByActive(page int, pageSize int, search string) ([]*response.TransactionResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.TransactionResponseDeleteAt, int, *response.ErrorResponse)
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
	FindAll(page int, pageSize int, search string) ([]*response.TransferResponse, int, *response.ErrorResponse)
	FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse)

	FindMonthTransferStatusSuccess(year int, month int) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransferStatusFailed(year int, month int) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthTransferStatusSuccessByCardNumber(cardNumber string, year int, month int) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyTransferStatusSuccessByCardNumber(cardNumber string, year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthTransferStatusFailedByCardNumber(cardNumber string, year int, month int) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyTransferStatusFailedByCardNumber(cardNumber string, year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
	FindMonthlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindMonthlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmountsBySenderCardNumber(cardNumber string, year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
	FindYearlyTransferAmountsByReceiverCardNumber(cardNumber string, year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)

	FindByActive(page int, pageSize int, search string) ([]*response.TransferResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.TransferResponseDeleteAt, int, *response.ErrorResponse)
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
	FindAll(page int, pageSize int, search string) ([]*response.UserResponse, int, *response.ErrorResponse)
	FindByID(id int) (*response.UserResponse, *response.ErrorResponse)
	FindByActive(page int, pageSize int, search string) ([]*response.UserResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.UserResponseDeleteAt, int, *response.ErrorResponse)
	CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)
	TrashedUser(user_id int) (*response.UserResponse, *response.ErrorResponse)
	RestoreUser(user_id int) (*response.UserResponse, *response.ErrorResponse)
	DeleteUserPermanent(user_id int) (bool, *response.ErrorResponse)

	RestoreAllUser() (bool, *response.ErrorResponse)
	DeleteAllUserPermanent() (bool, *response.ErrorResponse)
}

type WithdrawService interface {
	FindAll(page int, pageSize int, search string) ([]*response.WithdrawResponse, int, *response.ErrorResponse)
	FindAllByCardNumber(card_number string, page int, pageSize int, search string) ([]*response.WithdrawResponse, int, *response.ErrorResponse)

	FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse)

	FindMonthWithdrawStatusSuccess(year int, month int) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyWithdrawStatusSuccess(year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthWithdrawStatusFailed(year int, month int) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyWithdrawStatusFailed(year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthWithdrawStatusSuccessByCardNumber(card_number string, year int, month int) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)
	FindYearlyWithdrawStatusSuccessByCardNumber(card_number string, year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)
	FindMonthWithdrawStatusFailedByCardNumber(card_number string, year int, month int) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)
	FindYearlyWithdrawStatusFailedByCardNumber(card_number string, year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)

	FindMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	FindYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
	FindMonthlyWithdrawsByCardNumber(cardNumber string, year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)
	FindYearlyWithdrawsByCardNumber(cardNumber string, year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)

	FindByActive(page int, pageSize int, search string) ([]*response.WithdrawResponseDeleteAt, int, *response.ErrorResponse)
	FindByTrashed(page int, pageSize int, search string) ([]*response.WithdrawResponseDeleteAt, int, *response.ErrorResponse)
	Create(request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)
	Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse)
	TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)
	RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse)
	DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse)

	RestoreAllWithdraw() (bool, *response.ErrorResponse)
	DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse)
}
