// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -destination=mocks/protomapper.go
//

// Package mock_protomapper is a generated GoMock package.
package mock_protomapper

import (
	response "MamangRust/paymentgatewaygrpc/internal/domain/response"
	pb "MamangRust/paymentgatewaygrpc/internal/pb"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthProtoMapper is a mock of AuthProtoMapper interface.
type MockAuthProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockAuthProtoMapperMockRecorder
	isgomock struct{}
}

// MockAuthProtoMapperMockRecorder is the mock recorder for MockAuthProtoMapper.
type MockAuthProtoMapperMockRecorder struct {
	mock *MockAuthProtoMapper
}

// NewMockAuthProtoMapper creates a new mock instance.
func NewMockAuthProtoMapper(ctrl *gomock.Controller) *MockAuthProtoMapper {
	mock := &MockAuthProtoMapper{ctrl: ctrl}
	mock.recorder = &MockAuthProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthProtoMapper) EXPECT() *MockAuthProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseLogin mocks base method.
func (m *MockAuthProtoMapper) ToResponseLogin(token string) *pb.ApiResponseLogin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseLogin", token)
	ret0, _ := ret[0].(*pb.ApiResponseLogin)
	return ret0
}

// ToResponseLogin indicates an expected call of ToResponseLogin.
func (mr *MockAuthProtoMapperMockRecorder) ToResponseLogin(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseLogin", reflect.TypeOf((*MockAuthProtoMapper)(nil).ToResponseLogin), token)
}

// ToResponseRegister mocks base method.
func (m *MockAuthProtoMapper) ToResponseRegister(response *response.UserResponse) *pb.ApiResponseRegister {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseRegister", response)
	ret0, _ := ret[0].(*pb.ApiResponseRegister)
	return ret0
}

// ToResponseRegister indicates an expected call of ToResponseRegister.
func (mr *MockAuthProtoMapperMockRecorder) ToResponseRegister(response any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseRegister", reflect.TypeOf((*MockAuthProtoMapper)(nil).ToResponseRegister), response)
}

// MockCardProtoMapper is a mock of CardProtoMapper interface.
type MockCardProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockCardProtoMapperMockRecorder
	isgomock struct{}
}

// MockCardProtoMapperMockRecorder is the mock recorder for MockCardProtoMapper.
type MockCardProtoMapperMockRecorder struct {
	mock *MockCardProtoMapper
}

// NewMockCardProtoMapper creates a new mock instance.
func NewMockCardProtoMapper(ctrl *gomock.Controller) *MockCardProtoMapper {
	mock := &MockCardProtoMapper{ctrl: ctrl}
	mock.recorder = &MockCardProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCardProtoMapper) EXPECT() *MockCardProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseCard mocks base method.
func (m *MockCardProtoMapper) ToResponseCard(card *response.CardResponse) *pb.CardResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseCard", card)
	ret0, _ := ret[0].(*pb.CardResponse)
	return ret0
}

// ToResponseCard indicates an expected call of ToResponseCard.
func (mr *MockCardProtoMapperMockRecorder) ToResponseCard(card any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseCard", reflect.TypeOf((*MockCardProtoMapper)(nil).ToResponseCard), card)
}

// ToResponsesCard mocks base method.
func (m *MockCardProtoMapper) ToResponsesCard(cards []*response.CardResponse) []*pb.CardResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesCard", cards)
	ret0, _ := ret[0].([]*pb.CardResponse)
	return ret0
}

// ToResponsesCard indicates an expected call of ToResponsesCard.
func (mr *MockCardProtoMapperMockRecorder) ToResponsesCard(cards any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesCard", reflect.TypeOf((*MockCardProtoMapper)(nil).ToResponsesCard), cards)
}

// MockMerchantProtoMapper is a mock of MerchantProtoMapper interface.
type MockMerchantProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockMerchantProtoMapperMockRecorder
	isgomock struct{}
}

// MockMerchantProtoMapperMockRecorder is the mock recorder for MockMerchantProtoMapper.
type MockMerchantProtoMapperMockRecorder struct {
	mock *MockMerchantProtoMapper
}

// NewMockMerchantProtoMapper creates a new mock instance.
func NewMockMerchantProtoMapper(ctrl *gomock.Controller) *MockMerchantProtoMapper {
	mock := &MockMerchantProtoMapper{ctrl: ctrl}
	mock.recorder = &MockMerchantProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerchantProtoMapper) EXPECT() *MockMerchantProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseMerchant mocks base method.
func (m *MockMerchantProtoMapper) ToResponseMerchant(merchant *response.MerchantResponse) *pb.MerchantResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseMerchant", merchant)
	ret0, _ := ret[0].(*pb.MerchantResponse)
	return ret0
}

// ToResponseMerchant indicates an expected call of ToResponseMerchant.
func (mr *MockMerchantProtoMapperMockRecorder) ToResponseMerchant(merchant any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseMerchant", reflect.TypeOf((*MockMerchantProtoMapper)(nil).ToResponseMerchant), merchant)
}

// ToResponsesMerchant mocks base method.
func (m *MockMerchantProtoMapper) ToResponsesMerchant(merchants []*response.MerchantResponse) []*pb.MerchantResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesMerchant", merchants)
	ret0, _ := ret[0].([]*pb.MerchantResponse)
	return ret0
}

// ToResponsesMerchant indicates an expected call of ToResponsesMerchant.
func (mr *MockMerchantProtoMapperMockRecorder) ToResponsesMerchant(merchants any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesMerchant", reflect.TypeOf((*MockMerchantProtoMapper)(nil).ToResponsesMerchant), merchants)
}

// MockSaldoProtoMapper is a mock of SaldoProtoMapper interface.
type MockSaldoProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockSaldoProtoMapperMockRecorder
	isgomock struct{}
}

// MockSaldoProtoMapperMockRecorder is the mock recorder for MockSaldoProtoMapper.
type MockSaldoProtoMapperMockRecorder struct {
	mock *MockSaldoProtoMapper
}

// NewMockSaldoProtoMapper creates a new mock instance.
func NewMockSaldoProtoMapper(ctrl *gomock.Controller) *MockSaldoProtoMapper {
	mock := &MockSaldoProtoMapper{ctrl: ctrl}
	mock.recorder = &MockSaldoProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSaldoProtoMapper) EXPECT() *MockSaldoProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseSaldo mocks base method.
func (m *MockSaldoProtoMapper) ToResponseSaldo(saldo *response.SaldoResponse) *pb.SaldoResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseSaldo", saldo)
	ret0, _ := ret[0].(*pb.SaldoResponse)
	return ret0
}

// ToResponseSaldo indicates an expected call of ToResponseSaldo.
func (mr *MockSaldoProtoMapperMockRecorder) ToResponseSaldo(saldo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseSaldo", reflect.TypeOf((*MockSaldoProtoMapper)(nil).ToResponseSaldo), saldo)
}

// ToResponsesSaldo mocks base method.
func (m *MockSaldoProtoMapper) ToResponsesSaldo(saldos []*response.SaldoResponse) []*pb.SaldoResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesSaldo", saldos)
	ret0, _ := ret[0].([]*pb.SaldoResponse)
	return ret0
}

// ToResponsesSaldo indicates an expected call of ToResponsesSaldo.
func (mr *MockSaldoProtoMapperMockRecorder) ToResponsesSaldo(saldos any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesSaldo", reflect.TypeOf((*MockSaldoProtoMapper)(nil).ToResponsesSaldo), saldos)
}

// MockTopupProtoMapper is a mock of TopupProtoMapper interface.
type MockTopupProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockTopupProtoMapperMockRecorder
	isgomock struct{}
}

// MockTopupProtoMapperMockRecorder is the mock recorder for MockTopupProtoMapper.
type MockTopupProtoMapperMockRecorder struct {
	mock *MockTopupProtoMapper
}

// NewMockTopupProtoMapper creates a new mock instance.
func NewMockTopupProtoMapper(ctrl *gomock.Controller) *MockTopupProtoMapper {
	mock := &MockTopupProtoMapper{ctrl: ctrl}
	mock.recorder = &MockTopupProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopupProtoMapper) EXPECT() *MockTopupProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseTopup mocks base method.
func (m *MockTopupProtoMapper) ToResponseTopup(topup *response.TopupResponse) *pb.TopupResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseTopup", topup)
	ret0, _ := ret[0].(*pb.TopupResponse)
	return ret0
}

// ToResponseTopup indicates an expected call of ToResponseTopup.
func (mr *MockTopupProtoMapperMockRecorder) ToResponseTopup(topup any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseTopup", reflect.TypeOf((*MockTopupProtoMapper)(nil).ToResponseTopup), topup)
}

// ToResponsesTopup mocks base method.
func (m *MockTopupProtoMapper) ToResponsesTopup(topups []*response.TopupResponse) []*pb.TopupResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesTopup", topups)
	ret0, _ := ret[0].([]*pb.TopupResponse)
	return ret0
}

// ToResponsesTopup indicates an expected call of ToResponsesTopup.
func (mr *MockTopupProtoMapperMockRecorder) ToResponsesTopup(topups any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesTopup", reflect.TypeOf((*MockTopupProtoMapper)(nil).ToResponsesTopup), topups)
}

// MockTransactionProtoMapper is a mock of TransactionProtoMapper interface.
type MockTransactionProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionProtoMapperMockRecorder
	isgomock struct{}
}

// MockTransactionProtoMapperMockRecorder is the mock recorder for MockTransactionProtoMapper.
type MockTransactionProtoMapperMockRecorder struct {
	mock *MockTransactionProtoMapper
}

// NewMockTransactionProtoMapper creates a new mock instance.
func NewMockTransactionProtoMapper(ctrl *gomock.Controller) *MockTransactionProtoMapper {
	mock := &MockTransactionProtoMapper{ctrl: ctrl}
	mock.recorder = &MockTransactionProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionProtoMapper) EXPECT() *MockTransactionProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseTransaction mocks base method.
func (m *MockTransactionProtoMapper) ToResponseTransaction(transaction *response.TransactionResponse) *pb.TransactionResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseTransaction", transaction)
	ret0, _ := ret[0].(*pb.TransactionResponse)
	return ret0
}

// ToResponseTransaction indicates an expected call of ToResponseTransaction.
func (mr *MockTransactionProtoMapperMockRecorder) ToResponseTransaction(transaction any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseTransaction", reflect.TypeOf((*MockTransactionProtoMapper)(nil).ToResponseTransaction), transaction)
}

// ToResponsesTransaction mocks base method.
func (m *MockTransactionProtoMapper) ToResponsesTransaction(transactions []*response.TransactionResponse) []*pb.TransactionResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesTransaction", transactions)
	ret0, _ := ret[0].([]*pb.TransactionResponse)
	return ret0
}

// ToResponsesTransaction indicates an expected call of ToResponsesTransaction.
func (mr *MockTransactionProtoMapperMockRecorder) ToResponsesTransaction(transactions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesTransaction", reflect.TypeOf((*MockTransactionProtoMapper)(nil).ToResponsesTransaction), transactions)
}

// MockTransferProtoMapper is a mock of TransferProtoMapper interface.
type MockTransferProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockTransferProtoMapperMockRecorder
	isgomock struct{}
}

// MockTransferProtoMapperMockRecorder is the mock recorder for MockTransferProtoMapper.
type MockTransferProtoMapperMockRecorder struct {
	mock *MockTransferProtoMapper
}

// NewMockTransferProtoMapper creates a new mock instance.
func NewMockTransferProtoMapper(ctrl *gomock.Controller) *MockTransferProtoMapper {
	mock := &MockTransferProtoMapper{ctrl: ctrl}
	mock.recorder = &MockTransferProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferProtoMapper) EXPECT() *MockTransferProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseTransfer mocks base method.
func (m *MockTransferProtoMapper) ToResponseTransfer(transfer *response.TransferResponse) *pb.TransferResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseTransfer", transfer)
	ret0, _ := ret[0].(*pb.TransferResponse)
	return ret0
}

// ToResponseTransfer indicates an expected call of ToResponseTransfer.
func (mr *MockTransferProtoMapperMockRecorder) ToResponseTransfer(transfer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseTransfer", reflect.TypeOf((*MockTransferProtoMapper)(nil).ToResponseTransfer), transfer)
}

// ToResponsesTransfer mocks base method.
func (m *MockTransferProtoMapper) ToResponsesTransfer(transfers []*response.TransferResponse) []*pb.TransferResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesTransfer", transfers)
	ret0, _ := ret[0].([]*pb.TransferResponse)
	return ret0
}

// ToResponsesTransfer indicates an expected call of ToResponsesTransfer.
func (mr *MockTransferProtoMapperMockRecorder) ToResponsesTransfer(transfers any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesTransfer", reflect.TypeOf((*MockTransferProtoMapper)(nil).ToResponsesTransfer), transfers)
}

// MockUserProtoMapper is a mock of UserProtoMapper interface.
type MockUserProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockUserProtoMapperMockRecorder
	isgomock struct{}
}

// MockUserProtoMapperMockRecorder is the mock recorder for MockUserProtoMapper.
type MockUserProtoMapperMockRecorder struct {
	mock *MockUserProtoMapper
}

// NewMockUserProtoMapper creates a new mock instance.
func NewMockUserProtoMapper(ctrl *gomock.Controller) *MockUserProtoMapper {
	mock := &MockUserProtoMapper{ctrl: ctrl}
	mock.recorder = &MockUserProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserProtoMapper) EXPECT() *MockUserProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseUser mocks base method.
func (m *MockUserProtoMapper) ToResponseUser(user *response.UserResponse) *pb.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseUser", user)
	ret0, _ := ret[0].(*pb.UserResponse)
	return ret0
}

// ToResponseUser indicates an expected call of ToResponseUser.
func (mr *MockUserProtoMapperMockRecorder) ToResponseUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseUser", reflect.TypeOf((*MockUserProtoMapper)(nil).ToResponseUser), user)
}

// ToResponsesUser mocks base method.
func (m *MockUserProtoMapper) ToResponsesUser(users []*response.UserResponse) []*pb.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesUser", users)
	ret0, _ := ret[0].([]*pb.UserResponse)
	return ret0
}

// ToResponsesUser indicates an expected call of ToResponsesUser.
func (mr *MockUserProtoMapperMockRecorder) ToResponsesUser(users any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesUser", reflect.TypeOf((*MockUserProtoMapper)(nil).ToResponsesUser), users)
}

// MockWithdrawalProtoMapper is a mock of WithdrawalProtoMapper interface.
type MockWithdrawalProtoMapper struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawalProtoMapperMockRecorder
	isgomock struct{}
}

// MockWithdrawalProtoMapperMockRecorder is the mock recorder for MockWithdrawalProtoMapper.
type MockWithdrawalProtoMapperMockRecorder struct {
	mock *MockWithdrawalProtoMapper
}

// NewMockWithdrawalProtoMapper creates a new mock instance.
func NewMockWithdrawalProtoMapper(ctrl *gomock.Controller) *MockWithdrawalProtoMapper {
	mock := &MockWithdrawalProtoMapper{ctrl: ctrl}
	mock.recorder = &MockWithdrawalProtoMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawalProtoMapper) EXPECT() *MockWithdrawalProtoMapperMockRecorder {
	return m.recorder
}

// ToResponseWithdrawal mocks base method.
func (m *MockWithdrawalProtoMapper) ToResponseWithdrawal(withdrawal *response.WithdrawResponse) *pb.WithdrawResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponseWithdrawal", withdrawal)
	ret0, _ := ret[0].(*pb.WithdrawResponse)
	return ret0
}

// ToResponseWithdrawal indicates an expected call of ToResponseWithdrawal.
func (mr *MockWithdrawalProtoMapperMockRecorder) ToResponseWithdrawal(withdrawal any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponseWithdrawal", reflect.TypeOf((*MockWithdrawalProtoMapper)(nil).ToResponseWithdrawal), withdrawal)
}

// ToResponsesWithdrawal mocks base method.
func (m *MockWithdrawalProtoMapper) ToResponsesWithdrawal(withdrawals []*response.WithdrawResponse) []*pb.WithdrawResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToResponsesWithdrawal", withdrawals)
	ret0, _ := ret[0].([]*pb.WithdrawResponse)
	return ret0
}

// ToResponsesWithdrawal indicates an expected call of ToResponsesWithdrawal.
func (mr *MockWithdrawalProtoMapperMockRecorder) ToResponsesWithdrawal(withdrawals any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToResponsesWithdrawal", reflect.TypeOf((*MockWithdrawalProtoMapper)(nil).ToResponsesWithdrawal), withdrawals)
}
