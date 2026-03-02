package apimapper

type ResponseApiMapper struct {
	AuthResponseMapper        AuthResponseMapper
	CardResponseMapper        CardResponseMapper
	RoleResponseMapper        RoleResponseMapper
	SaldoResponseMapper       SaldoResponseMapper
	TransactionResponseMapper TransactionResponseMapper
	TransferResponseMapper    TransferResponseMapper
	TopupResponseMapper       TopupResponseMapper
	WithdrawResponseMapper    WithdrawResponseMapper
	UserResponseMapper        UserResponseMapper
	MerchantResponseMapper    MerchantResponseMapper
}

func NewResponseApiMapper() *ResponseApiMapper {
	return &ResponseApiMapper{
		AuthResponseMapper:        NewAuthResponseMapper(),
		CardResponseMapper:        NewCardResponseMapper(),
		SaldoResponseMapper:       NewSaldoResponseMapper(),
		TransactionResponseMapper: NewTransactionResponseMapper(),
		TransferResponseMapper:    NewTransferResponseMapper(),
		TopupResponseMapper:       NewTopupResponseMapper(),
		WithdrawResponseMapper:    NewWithdrawResponseMapper(),
		UserResponseMapper:        NewUserResponseMapper(),
		RoleResponseMapper:        NewRoleResponseMapper(),
		MerchantResponseMapper:    NewMerchantResponseMapper(),
	}
}
