package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/service"
)

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

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Auth:        NewAuthHandleGrpc(service.Auth),
		Role:        NewRoleHandleGrpc(service.Role),
		User:        NewUserHandleGrpc(service.User),
		Card:        NewCardHandleGrpc(service.Card),
		Merchant:    NewMerchantHandleGrpc(service.Merchant),
		Transaction: NewTransactionHandleGrpc(service.Transaction),
		Saldo:       NewSaldoHandleGrpc(service.Saldo),
		Topup:       NewTopupHandleGrpc(service.Topup),
		Transfer:    NewTransferHandleGrpc(service.Transfer),
		Withdraw:    NewWithdrawHandleGrpc(service.Withdraw),
	}
}
