package recordmapper

type RecordMapper struct {
	UserRecordMapper         UserRecordMapping
	RoleRecordMapper         RoleRecordMapping
	UserRoleRecordMapper     UserRoleRecordMapping
	RefreshTokenRecordMapper RefreshTokenRecordMapping
	SaldoRecordMapper        SaldoRecordMapping
	TopupRecordMapper        TopupRecordMapping
	TransferRecordMapper     TransferRecordMapping
	WithdrawRecordMapper     WithdrawRecordMapping
	CardRecordMapper         CardRecordMapping
	TransactionRecordMapper  TransactionRecordMapping
	MerchantRecordMapper     MerchantRecordMapping
}

func NewRecordMapper() *RecordMapper {
	return &RecordMapper{
		UserRecordMapper:         NewUserRecordMapper(),
		RoleRecordMapper:         NewRoleRecordMapper(),
		UserRoleRecordMapper:     NewUserRoleRecordMapper(),
		RefreshTokenRecordMapper: NewRefreshTokenRecordMapper(),
		SaldoRecordMapper:        NewSaldoRecordMapper(),
		TopupRecordMapper:        NewTopupRecordMapper(),
		TransferRecordMapper:     NewTransferRecordMapper(),
		WithdrawRecordMapper:     NewWithdrawRecordMapper(),
		CardRecordMapper:         NewCardRecordMapper(),
		TransactionRecordMapper:  NewTransactionRecordMapper(),
		MerchantRecordMapper:     NewMerchantRecordMapper(),
	}
}
