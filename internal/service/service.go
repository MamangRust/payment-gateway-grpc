package service

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	auth_cache "MamangRust/paymentgatewaygrpc/internal/cache/auth"
	card_cache "MamangRust/paymentgatewaygrpc/internal/cache/card"
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/merchant"
	role_cache "MamangRust/paymentgatewaygrpc/internal/cache/role"
	saldo_cache "MamangRust/paymentgatewaygrpc/internal/cache/saldo"
	topup_cache "MamangRust/paymentgatewaygrpc/internal/cache/topup"
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/transaction"
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/transfer"
	user_cache "MamangRust/paymentgatewaygrpc/internal/cache/user"
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/withdraw"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
)

type Service struct {
	Auth        AuthService
	User        UserService
	Role        RoleService
	Saldo       SaldoService
	Topup       TopupService
	Transfer    TransferService
	Withdraw    WithdrawService
	Card        CardService
	Merchant    MerchantService
	Transaction TransactionService
}

type Deps struct {
	Repositories *repository.Repositories
	Token        auth.TokenManager
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
	Cache        *cache.CacheStore
}

func NewService(deps Deps) *Service {
	observability, _ := observability.NewObservability("grpc-server", deps.Logger)

	cacheAuth := auth_cache.NewMencache(deps.Cache)
	cacheUser := user_cache.NewUserMencache(deps.Cache)
	cacheRole := role_cache.NewRoleMencache(deps.Cache)
	cacheSaldo := saldo_cache.NewSaldoMencache(deps.Cache)
	cacheTopup := topup_cache.NewTopupMencache(deps.Cache)
	cacheTransaction := transaction_cache.NewTransactionMencache(deps.Cache)
	cacheTransfer := transfer_cache.NewTransferMencache(deps.Cache)
	cacheWithdraw := withdraw_cache.NewWithdrawMencache(deps.Cache)
	cacheMerchant := merchant_cache.NewMerchantMencache(deps.Cache)
	cacheCard := card_cache.NewCardMencache(deps.Cache)

	return &Service{
		Auth: NewAuthService(AuthServiceDeps{
			AuthRepo:         deps.Repositories.User,
			RefreshTokenRepo: deps.Repositories.RefreshToken,
			RoleRepo:         deps.Repositories.Role,
			UserRoleRepo:     deps.Repositories.UserRole,
			Hash:             deps.Hash,
			Token:            deps.Token,
			Logger:           deps.Logger,
			Tracer:           observability,
			CacheIdentity:    cacheAuth.IdentityCache,
			CacheLogin:       cacheAuth.LoginCache,
		}),
		User: NewUserService(UserServiceDeps{
			UserRepo:      deps.Repositories.User,
			Logger:        deps.Logger,
			Observability: observability,
			Hashing:       deps.Hash,
			Cache:         cacheUser,
		}),
		Role: NewRoleService(RoleServiceDeps{
			RoleRepo:      deps.Repositories.Role,
			Observability: observability,
			Logger:        deps.Logger,
			Cache:         cacheRole,
		}),
		Saldo: NewSaldoService(SaldoServiceDeps{
			SaldoRepo:     deps.Repositories.Saldo,
			CardRepo:      deps.Repositories.Card,
			Logger:        deps.Logger,
			Observability: observability,
			Cache:         cacheSaldo,
		}),
		Topup: NewTopupService(TopupServiceDeps{
			CardRepo:      deps.Repositories.Card,
			TopupRepo:     deps.Repositories.Topup,
			SaldoRepo:     deps.Repositories.Saldo,
			Logger:        deps.Logger,
			Observability: observability,
			Cache:         cacheTopup,
		}),
		Transfer: NewTransferService(TransferServiceDeps{
			UserRepo:      deps.Repositories.User,
			CardRepo:      deps.Repositories.Card,
			SaldoRepo:     deps.Repositories.Saldo,
			TransferRepo:  deps.Repositories.Transfer,
			Logger:        deps.Logger,
			Observability: observability,
			Cache:         cacheTransfer,
		}),
		Withdraw: NewWithdrawService(WithdrawServiceDeps{
			UserRepo:      deps.Repositories.User,
			SaldoRepo:     deps.Repositories.Saldo,
			WithdrawRepo:  deps.Repositories.Withdraw,
			Logger:        deps.Logger,
			Observability: observability,
			Cache:         cacheWithdraw,
		}),
		Card: NewCardService(CardServiceDeps{
			CardRepo:      deps.Repositories.Card,
			UserRepo:      deps.Repositories.User,
			Logger:        deps.Logger,
			Observability: observability,
			cache:         cacheCard,
		}),
		Merchant: NewMerchantService(MerchantServiceDeps{
			UserRepo:      deps.Repositories.User,
			MerchantRepo:  deps.Repositories.Merchant,
			Logger:        deps.Logger,
			Observability: observability,
			Cache:         cacheMerchant,
		}),
		Transaction: NewTransactionService(TransactionServiceDeps{
			MerchantRepo:    deps.Repositories.Merchant,
			CardRepo:        deps.Repositories.Card,
			SaldoRepo:       deps.Repositories.Saldo,
			TransactionRepo: deps.Repositories.Transaction,
			Logger:          deps.Logger,
			Cache:           cacheTransaction,
		}),
	}
}
