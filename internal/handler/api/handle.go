package api

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	auth_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/auth"
	card_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/card"
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/merchant"
	role_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/role"
	saldo_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/saldo"
	topup_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/topup"
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transaction"
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transfer"
	user_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/user"
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/withdraw"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"

	"github.com/labstack/echo/v4"

	"google.golang.org/grpc"
)

type Deps struct {
	Conn    *grpc.ClientConn
	Token   auth.TokenManager
	E       *echo.Echo
	Logger  logger.LoggerInterface
	Mapping *apimapper.ResponseApiMapper
	Cache   *cache.CacheStore
}

func NewHandler(deps Deps) {
	observability, _ := observability.NewObservability("client", deps.Logger)

	apiHandler := errors.NewApiHandler(observability, deps.Logger)

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

	clientAuth := pb.NewAuthServiceClient(deps.Conn)
	clientRole := pb.NewRoleServiceClient(deps.Conn)
	clientCard := pb.NewCardServiceClient(deps.Conn)
	clientMerchant := pb.NewMerchantServiceClient(deps.Conn)
	clientUser := pb.NewUserServiceClient(deps.Conn)
	clientSaldo := pb.NewSaldoServiceClient(deps.Conn)
	clientTopup := pb.NewTopupServiceClient(deps.Conn)
	clientTransaction := pb.NewTransactionServiceClient(deps.Conn)
	clientTransfer := pb.NewTransferServiceClient(deps.Conn)
	clientWithdraw := pb.NewWithdrawServiceClient(deps.Conn)

	NewHandlerAuth(
		deps.E,
		clientAuth,
		deps.Logger,
		deps.Mapping.AuthResponseMapper,
		apiHandler,
		cacheAuth,
	)

	NewHandlerRole(
		clientRole,
		deps.E,
		deps.Logger,
		deps.Mapping.RoleResponseMapper,
		apiHandler,
		cacheRole,
	)

	NewHandlerUser(
		clientUser,
		deps.E,
		deps.Logger,
		deps.Mapping.UserResponseMapper,
		apiHandler,
		cacheUser,
	)

	NewHandlerCard(
		clientCard,
		deps.E,
		deps.Logger,
		apiHandler,
		deps.Mapping.CardResponseMapper,
		cacheCard,
	)

	NewHandlerMerchant(
		clientMerchant,
		deps.E,
		deps.Logger,
		apiHandler,
		deps.Mapping.MerchantResponseMapper,
		cacheMerchant,
	)

	NewHandlerTransaction(
		clientTransaction,
		clientMerchant,
		deps.E,
		deps.Logger,
		deps.Mapping.TransactionResponseMapper,
		apiHandler,
		cacheTransaction,
	)

	NewHandlerSaldo(
		clientSaldo,
		deps.E,
		deps.Logger,
		deps.Mapping.SaldoResponseMapper,
		apiHandler,
		cacheSaldo,
	)

	NewHandlerTopup(
		clientTopup,
		deps.E,
		deps.Logger,
		deps.Mapping.TopupResponseMapper,
		apiHandler,
		cacheTopup,
	)

	NewHandlerTransfer(
		clientTransfer,
		deps.E,
		deps.Logger,
		deps.Mapping.TransferResponseMapper,
		apiHandler,
		cacheTransfer,
	)

	NewHandlerWithdraw(
		clientWithdraw,
		deps.E,
		deps.Logger,
		deps.Mapping.WithdrawResponseMapper,
		apiHandler,
		cacheWithdraw,
	)
}
