// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -destination=mocks/interfaces.go
//

// Package mock_responsemapper is a generated GoMock package.
package mock_responsemapper

import (
	record "MamangRust/paymentgatewaygrpc/internal/domain/record"
	response "MamangRust/paymentgatewaygrpc/internal/domain/response"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCardResponseMapper is a mock of CardResponseMapper interface.
type MockCardResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockCardResponseMapperMockRecorder
	isgomock struct{}
}

// MockCardResponseMapperMockRecorder is the mock recorder for MockCardResponseMapper.
type MockCardResponseMapperMockRecorder struct {
	mock *MockCardResponseMapper
}

// NewMockCardResponseMapper creates a new mock instance.
func NewMockCardResponseMapper(ctrl *gomock.Controller) *MockCardResponseMapper {
	mock := &MockCardResponseMapper{ctrl: ctrl}
	mock.recorder = &MockCardResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCardResponseMapper) EXPECT() *MockCardResponseMapperMockRecorder {
	return m.recorder
}

// ToCardResponse mocks base method.
func (m *MockCardResponseMapper) ToCardResponse(card *record.CardRecord) *response.CardResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToCardResponse", card)
	ret0, _ := ret[0].(*response.CardResponse)
	return ret0
}

// ToCardResponse indicates an expected call of ToCardResponse.
func (mr *MockCardResponseMapperMockRecorder) ToCardResponse(card any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToCardResponse", reflect.TypeOf((*MockCardResponseMapper)(nil).ToCardResponse), card)
}

// ToCardsResponse mocks base method.
func (m *MockCardResponseMapper) ToCardsResponse(cards []*record.CardRecord) []*response.CardResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToCardsResponse", cards)
	ret0, _ := ret[0].([]*response.CardResponse)
	return ret0
}

// ToCardsResponse indicates an expected call of ToCardsResponse.
func (mr *MockCardResponseMapperMockRecorder) ToCardsResponse(cards any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToCardsResponse", reflect.TypeOf((*MockCardResponseMapper)(nil).ToCardsResponse), cards)
}

// MockUserResponseMapper is a mock of UserResponseMapper interface.
type MockUserResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockUserResponseMapperMockRecorder
	isgomock struct{}
}

// MockUserResponseMapperMockRecorder is the mock recorder for MockUserResponseMapper.
type MockUserResponseMapperMockRecorder struct {
	mock *MockUserResponseMapper
}

// NewMockUserResponseMapper creates a new mock instance.
func NewMockUserResponseMapper(ctrl *gomock.Controller) *MockUserResponseMapper {
	mock := &MockUserResponseMapper{ctrl: ctrl}
	mock.recorder = &MockUserResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserResponseMapper) EXPECT() *MockUserResponseMapperMockRecorder {
	return m.recorder
}

// ToUserResponse mocks base method.
func (m *MockUserResponseMapper) ToUserResponse(user *record.UserRecord) *response.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToUserResponse", user)
	ret0, _ := ret[0].(*response.UserResponse)
	return ret0
}

// ToUserResponse indicates an expected call of ToUserResponse.
func (mr *MockUserResponseMapperMockRecorder) ToUserResponse(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToUserResponse", reflect.TypeOf((*MockUserResponseMapper)(nil).ToUserResponse), user)
}

// ToUsersResponse mocks base method.
func (m *MockUserResponseMapper) ToUsersResponse(users []*record.UserRecord) []*response.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToUsersResponse", users)
	ret0, _ := ret[0].([]*response.UserResponse)
	return ret0
}

// ToUsersResponse indicates an expected call of ToUsersResponse.
func (mr *MockUserResponseMapperMockRecorder) ToUsersResponse(users any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToUsersResponse", reflect.TypeOf((*MockUserResponseMapper)(nil).ToUsersResponse), users)
}

// MockSaldoResponseMapper is a mock of SaldoResponseMapper interface.
type MockSaldoResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockSaldoResponseMapperMockRecorder
	isgomock struct{}
}

// MockSaldoResponseMapperMockRecorder is the mock recorder for MockSaldoResponseMapper.
type MockSaldoResponseMapperMockRecorder struct {
	mock *MockSaldoResponseMapper
}

// NewMockSaldoResponseMapper creates a new mock instance.
func NewMockSaldoResponseMapper(ctrl *gomock.Controller) *MockSaldoResponseMapper {
	mock := &MockSaldoResponseMapper{ctrl: ctrl}
	mock.recorder = &MockSaldoResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSaldoResponseMapper) EXPECT() *MockSaldoResponseMapperMockRecorder {
	return m.recorder
}

// ToSaldoResponse mocks base method.
func (m *MockSaldoResponseMapper) ToSaldoResponse(saldo *record.SaldoRecord) *response.SaldoResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToSaldoResponse", saldo)
	ret0, _ := ret[0].(*response.SaldoResponse)
	return ret0
}

// ToSaldoResponse indicates an expected call of ToSaldoResponse.
func (mr *MockSaldoResponseMapperMockRecorder) ToSaldoResponse(saldo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToSaldoResponse", reflect.TypeOf((*MockSaldoResponseMapper)(nil).ToSaldoResponse), saldo)
}

// ToSaldoResponses mocks base method.
func (m *MockSaldoResponseMapper) ToSaldoResponses(saldos []*record.SaldoRecord) []*response.SaldoResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToSaldoResponses", saldos)
	ret0, _ := ret[0].([]*response.SaldoResponse)
	return ret0
}

// ToSaldoResponses indicates an expected call of ToSaldoResponses.
func (mr *MockSaldoResponseMapperMockRecorder) ToSaldoResponses(saldos any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToSaldoResponses", reflect.TypeOf((*MockSaldoResponseMapper)(nil).ToSaldoResponses), saldos)
}

// MockTopupResponseMapper is a mock of TopupResponseMapper interface.
type MockTopupResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockTopupResponseMapperMockRecorder
	isgomock struct{}
}

// MockTopupResponseMapperMockRecorder is the mock recorder for MockTopupResponseMapper.
type MockTopupResponseMapperMockRecorder struct {
	mock *MockTopupResponseMapper
}

// NewMockTopupResponseMapper creates a new mock instance.
func NewMockTopupResponseMapper(ctrl *gomock.Controller) *MockTopupResponseMapper {
	mock := &MockTopupResponseMapper{ctrl: ctrl}
	mock.recorder = &MockTopupResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopupResponseMapper) EXPECT() *MockTopupResponseMapperMockRecorder {
	return m.recorder
}

// ToTopupResponse mocks base method.
func (m *MockTopupResponseMapper) ToTopupResponse(topup *record.TopupRecord) *response.TopupResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTopupResponse", topup)
	ret0, _ := ret[0].(*response.TopupResponse)
	return ret0
}

// ToTopupResponse indicates an expected call of ToTopupResponse.
func (mr *MockTopupResponseMapperMockRecorder) ToTopupResponse(topup any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTopupResponse", reflect.TypeOf((*MockTopupResponseMapper)(nil).ToTopupResponse), topup)
}

// ToTopupResponses mocks base method.
func (m *MockTopupResponseMapper) ToTopupResponses(topups []*record.TopupRecord) []*response.TopupResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTopupResponses", topups)
	ret0, _ := ret[0].([]*response.TopupResponse)
	return ret0
}

// ToTopupResponses indicates an expected call of ToTopupResponses.
func (mr *MockTopupResponseMapperMockRecorder) ToTopupResponses(topups any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTopupResponses", reflect.TypeOf((*MockTopupResponseMapper)(nil).ToTopupResponses), topups)
}

// MockTransactionResponseMapper is a mock of TransactionResponseMapper interface.
type MockTransactionResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionResponseMapperMockRecorder
	isgomock struct{}
}

// MockTransactionResponseMapperMockRecorder is the mock recorder for MockTransactionResponseMapper.
type MockTransactionResponseMapperMockRecorder struct {
	mock *MockTransactionResponseMapper
}

// NewMockTransactionResponseMapper creates a new mock instance.
func NewMockTransactionResponseMapper(ctrl *gomock.Controller) *MockTransactionResponseMapper {
	mock := &MockTransactionResponseMapper{ctrl: ctrl}
	mock.recorder = &MockTransactionResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionResponseMapper) EXPECT() *MockTransactionResponseMapperMockRecorder {
	return m.recorder
}

// ToTransactionResponse mocks base method.
func (m *MockTransactionResponseMapper) ToTransactionResponse(transaction *record.TransactionRecord) *response.TransactionResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTransactionResponse", transaction)
	ret0, _ := ret[0].(*response.TransactionResponse)
	return ret0
}

// ToTransactionResponse indicates an expected call of ToTransactionResponse.
func (mr *MockTransactionResponseMapperMockRecorder) ToTransactionResponse(transaction any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTransactionResponse", reflect.TypeOf((*MockTransactionResponseMapper)(nil).ToTransactionResponse), transaction)
}

// ToTransactionsResponse mocks base method.
func (m *MockTransactionResponseMapper) ToTransactionsResponse(transactions []*record.TransactionRecord) []*response.TransactionResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTransactionsResponse", transactions)
	ret0, _ := ret[0].([]*response.TransactionResponse)
	return ret0
}

// ToTransactionsResponse indicates an expected call of ToTransactionsResponse.
func (mr *MockTransactionResponseMapperMockRecorder) ToTransactionsResponse(transactions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTransactionsResponse", reflect.TypeOf((*MockTransactionResponseMapper)(nil).ToTransactionsResponse), transactions)
}

// MockTransferResponseMapper is a mock of TransferResponseMapper interface.
type MockTransferResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockTransferResponseMapperMockRecorder
	isgomock struct{}
}

// MockTransferResponseMapperMockRecorder is the mock recorder for MockTransferResponseMapper.
type MockTransferResponseMapperMockRecorder struct {
	mock *MockTransferResponseMapper
}

// NewMockTransferResponseMapper creates a new mock instance.
func NewMockTransferResponseMapper(ctrl *gomock.Controller) *MockTransferResponseMapper {
	mock := &MockTransferResponseMapper{ctrl: ctrl}
	mock.recorder = &MockTransferResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferResponseMapper) EXPECT() *MockTransferResponseMapperMockRecorder {
	return m.recorder
}

// ToTransferResponse mocks base method.
func (m *MockTransferResponseMapper) ToTransferResponse(transfer *record.TransferRecord) *response.TransferResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTransferResponse", transfer)
	ret0, _ := ret[0].(*response.TransferResponse)
	return ret0
}

// ToTransferResponse indicates an expected call of ToTransferResponse.
func (mr *MockTransferResponseMapperMockRecorder) ToTransferResponse(transfer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTransferResponse", reflect.TypeOf((*MockTransferResponseMapper)(nil).ToTransferResponse), transfer)
}

// ToTransfersResponse mocks base method.
func (m *MockTransferResponseMapper) ToTransfersResponse(transfers []*record.TransferRecord) []*response.TransferResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTransfersResponse", transfers)
	ret0, _ := ret[0].([]*response.TransferResponse)
	return ret0
}

// ToTransfersResponse indicates an expected call of ToTransfersResponse.
func (mr *MockTransferResponseMapperMockRecorder) ToTransfersResponse(transfers any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTransfersResponse", reflect.TypeOf((*MockTransferResponseMapper)(nil).ToTransfersResponse), transfers)
}

// MockWithdrawResponseMapper is a mock of WithdrawResponseMapper interface.
type MockWithdrawResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawResponseMapperMockRecorder
	isgomock struct{}
}

// MockWithdrawResponseMapperMockRecorder is the mock recorder for MockWithdrawResponseMapper.
type MockWithdrawResponseMapperMockRecorder struct {
	mock *MockWithdrawResponseMapper
}

// NewMockWithdrawResponseMapper creates a new mock instance.
func NewMockWithdrawResponseMapper(ctrl *gomock.Controller) *MockWithdrawResponseMapper {
	mock := &MockWithdrawResponseMapper{ctrl: ctrl}
	mock.recorder = &MockWithdrawResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawResponseMapper) EXPECT() *MockWithdrawResponseMapperMockRecorder {
	return m.recorder
}

// ToWithdrawResponse mocks base method.
func (m *MockWithdrawResponseMapper) ToWithdrawResponse(withdraw *record.WithdrawRecord) *response.WithdrawResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToWithdrawResponse", withdraw)
	ret0, _ := ret[0].(*response.WithdrawResponse)
	return ret0
}

// ToWithdrawResponse indicates an expected call of ToWithdrawResponse.
func (mr *MockWithdrawResponseMapperMockRecorder) ToWithdrawResponse(withdraw any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToWithdrawResponse", reflect.TypeOf((*MockWithdrawResponseMapper)(nil).ToWithdrawResponse), withdraw)
}

// ToWithdrawsResponse mocks base method.
func (m *MockWithdrawResponseMapper) ToWithdrawsResponse(withdraws []*record.WithdrawRecord) []*response.WithdrawResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToWithdrawsResponse", withdraws)
	ret0, _ := ret[0].([]*response.WithdrawResponse)
	return ret0
}

// ToWithdrawsResponse indicates an expected call of ToWithdrawsResponse.
func (mr *MockWithdrawResponseMapperMockRecorder) ToWithdrawsResponse(withdraws any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToWithdrawsResponse", reflect.TypeOf((*MockWithdrawResponseMapper)(nil).ToWithdrawsResponse), withdraws)
}

// MockMerchantResponseMapper is a mock of MerchantResponseMapper interface.
type MockMerchantResponseMapper struct {
	ctrl     *gomock.Controller
	recorder *MockMerchantResponseMapperMockRecorder
	isgomock struct{}
}

// MockMerchantResponseMapperMockRecorder is the mock recorder for MockMerchantResponseMapper.
type MockMerchantResponseMapperMockRecorder struct {
	mock *MockMerchantResponseMapper
}

// NewMockMerchantResponseMapper creates a new mock instance.
func NewMockMerchantResponseMapper(ctrl *gomock.Controller) *MockMerchantResponseMapper {
	mock := &MockMerchantResponseMapper{ctrl: ctrl}
	mock.recorder = &MockMerchantResponseMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerchantResponseMapper) EXPECT() *MockMerchantResponseMapperMockRecorder {
	return m.recorder
}

// ToMerchantResponse mocks base method.
func (m *MockMerchantResponseMapper) ToMerchantResponse(merchant *record.MerchantRecord) *response.MerchantResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToMerchantResponse", merchant)
	ret0, _ := ret[0].(*response.MerchantResponse)
	return ret0
}

// ToMerchantResponse indicates an expected call of ToMerchantResponse.
func (mr *MockMerchantResponseMapperMockRecorder) ToMerchantResponse(merchant any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToMerchantResponse", reflect.TypeOf((*MockMerchantResponseMapper)(nil).ToMerchantResponse), merchant)
}

// ToMerchantsResponse mocks base method.
func (m *MockMerchantResponseMapper) ToMerchantsResponse(merchants []*record.MerchantRecord) []*response.MerchantResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToMerchantsResponse", merchants)
	ret0, _ := ret[0].([]*response.MerchantResponse)
	return ret0
}

// ToMerchantsResponse indicates an expected call of ToMerchantsResponse.
func (mr *MockMerchantResponseMapperMockRecorder) ToMerchantsResponse(merchants any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToMerchantsResponse", reflect.TypeOf((*MockMerchantResponseMapper)(nil).ToMerchantsResponse), merchants)
}
