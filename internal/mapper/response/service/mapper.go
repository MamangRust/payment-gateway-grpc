package responseservice

type ResponseServiceMapper struct {
	CardResponseMapper         CardResponseMapper
	RoleResponseMapper         RoleResponseMapper
	RefreshTokenResponseMapper RefreshTokenResponseMapper
	SaldoResponseMapper        SaldoResponseMapper
	TransactionResponseMapper  TransactionResponseMapper
	TransferResponseMapper     TransferResponseMapper
	TopupResponseMapper        TopupResponseMapper
	WithdrawResponseMapper     WithdrawResponseMapper
	UserResponseMapper         UserResponseMapper
	MerchantResponseMapper     MerchantResponseMapper
}

func NewResponseServiceMapper() *ResponseServiceMapper {
	return &ResponseServiceMapper{
		CardResponseMapper:         NewCardResponseMapper(),
		SaldoResponseMapper:        NewSaldoResponseMapper(),
		TransactionResponseMapper:  NewTransactionResponseMapper(),
		TransferResponseMapper:     NewTransferResponseMapper(),
		TopupResponseMapper:        NewTopupResponseMapper(),
		WithdrawResponseMapper:     NewWithdrawResponseMapper(),
		UserResponseMapper:         NewUserResponseMapper(),
		RefreshTokenResponseMapper: NewRefreshTokenResponseMapper(),
		RoleResponseMapper:         NewRoleResponseMapper(),
		MerchantResponseMapper:     NewMerchantResponseMapper(),
	}
}
