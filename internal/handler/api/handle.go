package api

import (
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	"MamangRust/paymentgatewaygrpc/pkg/logger"

	"github.com/labstack/echo/v4"

	"google.golang.org/grpc"
)

type Deps struct {
	Conn    *grpc.ClientConn
	Token   auth.TokenManager
	E       *echo.Echo
	Logger  logger.LoggerInterface
	Mapping apimapper.ResponseApiMapper
}

func NewHandler(deps Deps) {

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

	NewHandlerAuth(clientAuth, deps.E, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(clientRole, deps.E, deps.Logger, deps.Mapping.RoleResponseMapper)
	NewHandlerUser(clientUser, deps.E, deps.Logger, deps.Mapping.UserResponseMapper)
	NewHandlerCard(clientCard, deps.E, deps.Logger, deps.Mapping.CardResponseMapper)
	NewHandlerMerchant(clientMerchant, deps.E, deps.Logger, deps.Mapping.MerchantResponseMapper)
	NewHandlerTransaction(clientTransaction, clientMerchant, deps.E, deps.Logger, deps.Mapping.TransactionResponseMapper)
	NewHandlerSaldo(clientSaldo, deps.E, deps.Logger, deps.Mapping.SaldoResponseMapper)
	NewHandlerTopup(clientTopup, deps.E, deps.Logger, deps.Mapping.TopupResponseMapper)
	NewHandlerTransfer(clientTransfer, deps.E, deps.Logger, deps.Mapping.TransferResponseMapper)
	NewHandlerWithdraw(clientWithdraw, deps.E, deps.Logger, deps.Mapping.WithdrawResponseMapper)
}
