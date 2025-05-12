package gapi

import (
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	"MamangRust/paymentgatewaygrpc/internal/service"
)

type Deps struct {
	Service service.Service
	Mapper  protomapper.ProtoMapper
}

type Handler struct {
	Auth        AuthHandleGrpc
	Role        RoleHandleGrpc
	User        UserHandleGrpc
	Card        CardHandleGrpc
	Merchant    MerchantHandleGrpc
	Transaction TransactionHandleGrpc
	Saldo       SaldoHandleGrpc
	Topup       TopupHandleGrpc
	Transfer    TransferHandleGrpc
	Withdraw    WithdrawHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Auth:        NewAuthHandleGrpc(deps.Service.Auth, deps.Mapper.AuthProtoMapper),
		Role:        NewRoleHandleGrpc(deps.Service.Role, deps.Mapper.RoleProtoMapper),
		User:        NewUserHandleGrpc(deps.Service.User, deps.Mapper.UserProtoMapper),
		Card:        NewCardHandleGrpc(deps.Service.Card, deps.Mapper.CardProtoMapper),
		Merchant:    NewMerchantHandleGrpc(deps.Service.Merchant, deps.Mapper.MerchantProtoMapper),
		Transaction: NewTransactionHandleGrpc(deps.Service.Transaction, deps.Mapper.TransactionProtoMapper),
		Saldo:       NewSaldoHandleGrpc(deps.Service.Saldo, deps.Mapper.SaldoProtoMapper),
		Topup:       NewTopupHandleGrpc(deps.Service.Topup, deps.Mapper.TopupProtoMapper),
		Transfer:    NewTransferHandleGrpc(deps.Service.Transfer, deps.Mapper.TransferProtoMapper),
		Withdraw:    NewWithdrawHandleGrpc(deps.Service.Withdraw, deps.Mapper.WithdrawalProtoMapper),
	}
}
