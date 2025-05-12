package repository

import (
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	DB "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
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

type Deps struct {
	DB           *DB.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		User:         NewUserRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
		Role:         NewRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.RoleRecordMapper),
		UserRole:     NewUserRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRoleRecordMapper),
		RefreshToken: NewRefreshTokenRepository(deps.DB, deps.Ctx, deps.MapperRecord.RefreshTokenRecordMapper),
		Saldo:        NewSaldoRepository(deps.DB, deps.Ctx, deps.MapperRecord.SaldoRecordMapper),
		Topup:        NewTopupRepository(deps.DB, deps.Ctx, deps.MapperRecord.TopupRecordMapper),
		Withdraw:     NewWithdrawRepository(deps.DB, deps.Ctx, deps.MapperRecord.WithdrawRecordMapper),
		Transfer:     NewTransferRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransferRecordMapper),
		Merchant:     NewMerchantRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		Card:         NewCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		Transaction:  NewTransactionRepository(deps.DB, deps.Ctx, deps.MapperRecord.TransactionRecordMapper),
	}
}
