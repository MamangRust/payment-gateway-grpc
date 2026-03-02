package repository

import (
	DB "MamangRust/paymentgatewaygrpc/pkg/database/schema"
)

type Repositories struct {
	User         UserRepository
	Saldo        SaldoRepository
	Role         RoleRepository
	UserRole     UserRoleRepository
	RefreshToken RefreshTokenRepository
	Topup        TopupRepository
	Withdraw     WithdrawRepository
	Transfer     TransferRepository
	Merchant     MerchantRepository
	Card         CardRepository
	Transaction  TransactionRepository
}

func NewRepositories(db *DB.Queries) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Role:         NewRoleRepository(db),
		UserRole:     NewUserRoleRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
		Saldo:        NewSaldoRepository(db),
		Topup:        NewTopupRepository(db),
		Withdraw:     NewWithdrawRepository(db),
		Transfer:     NewTransferRepository(db),
		Merchant:     NewMerchantRepository(db),
		Card:         NewCardRepository(db),
		Transaction:  NewTransactionRepository(db),
	}
}
