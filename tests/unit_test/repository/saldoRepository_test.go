package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllSaldos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldos := []*record.SaldoRecord{
		{
			ID:             1,
			CardNumber:     "1234",
			TotalBalance:   1000,
			WithdrawAmount: 500,
			WithdrawTime:   "2024-12-24T10:00:00Z",
			CreatedAt:      "2024-12-24T10:00:00Z",
			UpdatedAt:      "2024-12-24T10:00:00Z",
			DeletedAt:      nil,
		},
		{
			ID:             2,
			CardNumber:     "5678",
			TotalBalance:   2000,
			WithdrawAmount: 700,
			WithdrawTime:   "2024-12-24T11:00:00Z",
			CreatedAt:      "2024-12-24T11:00:00Z",
			UpdatedAt:      "2024-12-24T11:00:00Z",
			DeletedAt:      nil,
		},
	}

	mockRepo.EXPECT().FindAllSaldos("", 1, 10).Return(saldos, 2, nil)

	result, total, err := mockRepo.FindAllSaldos("", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, saldos, result)
	assert.Equal(t, 2, total)
}

func TestFindAllSaldos_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	mockRepo.EXPECT().FindAllSaldos("", 1, 10).Return(nil, 0, fmt.Errorf("database error"))

	result, total, err := mockRepo.FindAllSaldos("", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.EqualError(t, err, "database error")
}

func TestFindAllSaldos_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	mockRepo.EXPECT().FindAllSaldos("", 1, 10).Return([]*record.SaldoRecord{}, 1, nil)

	result, total, err := mockRepo.FindAllSaldos("", 1, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 1, total)
}

func TestFindByIdSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoId := 1

	expectedCard := &record.SaldoRecord{
		ID:             1,
		CardNumber:     "1234",
		TotalBalance:   1000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().FindById(saldoId).Return(expectedCard, nil)

	result, err := mockRepo.FindById(saldoId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard, result)
}

func TestFindByIdSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoId := 1

	mockRepo.EXPECT().FindById(saldoId).Return(nil, errors.New("saldo not found"))

	result, err := mockRepo.FindById(saldoId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "saldo not found")
}

func TestFindByCardNumberSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	cardNumber := "hesoyam"

	expectedCard := &record.SaldoRecord{
		ID:             1,
		CardNumber:     "hesoyam",
		TotalBalance:   1000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().FindByCardNumber(cardNumber).Return(expectedCard, nil)

	result, err := mockRepo.FindByCardNumber(cardNumber)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard, result)
}

func TestFindByIdCardNumberSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoNumber := "hesoyam"

	mockRepo.EXPECT().FindByCardNumber(saldoNumber).Return(nil, errors.New("saldo not found"))

	result, err := mockRepo.FindByCardNumber(saldoNumber)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "saldo not found")
}

func TestFindByActiveSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	activeSaldos := []*record.SaldoRecord{
		{
			ID:             1,
			CardNumber:     "1234",
			TotalBalance:   1000,
			WithdrawAmount: 500,
			WithdrawTime:   "2024-12-24T10:00:00Z",
			CreatedAt:      "2024-12-24T10:00:00Z",
			UpdatedAt:      "2024-12-24T10:00:00Z",
			DeletedAt:      nil,
		},
	}
	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(activeSaldos, 1, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, activeSaldos, result)
}

func TestFindByActiveSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to fetch active saldos"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.Contains(t, err.Error(), "failed to fetch active saldos")
}

func TestFindByActiveSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return([]*record.SaldoRecord{}, 0, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.Empty(t, result)
}

func TestFindByTrashedSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	trashedSaldos := []*record.SaldoRecord{
		{
			ID:             2,
			CardNumber:     "5678",
			TotalBalance:   2000,
			WithdrawAmount: 700,
			WithdrawTime:   "2024-12-24T11:00:00Z",
			CreatedAt:      "2024-12-24T11:00:00Z",
			UpdatedAt:      "2024-12-24T11:00:00Z",
			DeletedAt:      utils.PtrString("2024-12-24T11:00:00Z"),
		},
	}
	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(trashedSaldos, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.Equal(t, trashedSaldos, result)
}

func TestFindByTrashedSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, fmt.Errorf("failed to fetch trashed saldos"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Error(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch trashed saldos")
}

func TestFindByTrashedSaldo_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return([]*record.SaldoRecord{}, 0, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCreateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	request := requests.CreateSaldoRequest{
		CardNumber:   "1234",
		TotalBalance: 1000,
	}

	expectedSaldo := &record.SaldoRecord{
		ID:             1,
		CardNumber:     request.CardNumber,
		TotalBalance:   request.TotalBalance,
		WithdrawAmount: 0,
		WithdrawTime:   "",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-24T10:00:00Z",
	}

	mockRepo.EXPECT().CreateSaldo(&request).Return(expectedSaldo, nil)

	result, err := mockRepo.CreateSaldo(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSaldo, result)
}

func TestCreateSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	request := requests.CreateSaldoRequest{
		CardNumber:   "1234",
		TotalBalance: 1000,
	}

	mockRepo.EXPECT().CreateSaldo(&request).Return(nil, fmt.Errorf("failed to create saldo"))

	result, err := mockRepo.CreateSaldo(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create saldo")
}

func TestCreateSaldo_ValidationError(t *testing.T) {
	request := requests.CreateSaldoRequest{
		CardNumber:   "",
		TotalBalance: 0,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TotalBalance' failed on the 'required' tag")
}

func TestUpdateSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	request := requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "5678",
		TotalBalance: 2000,
	}

	updatedSaldo := &record.SaldoRecord{
		ID:             1,
		CardNumber:     "5678",
		TotalBalance:   2000,
		WithdrawAmount: 500,
		WithdrawTime:   "2024-12-24T10:00:00Z",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-25T10:00:00Z",
	}

	mockRepo.EXPECT().UpdateSaldo(&request).Return(updatedSaldo, nil)

	result, err := mockRepo.UpdateSaldo(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedSaldo, result)
}

func TestUpdateSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	request := requests.UpdateSaldoRequest{
		SaldoID:      1,
		CardNumber:   "5678",
		TotalBalance: 2000,
	}

	mockRepo.EXPECT().UpdateSaldo(&request).Return(nil, fmt.Errorf("saldo not found"))

	result, err := mockRepo.UpdateSaldo(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "saldo not found")
}

func TestUpdateSaldo_ValidationError(t *testing.T) {
	request := requests.UpdateSaldoRequest{
		SaldoID:      0,
		CardNumber:   "",
		TotalBalance: 0,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'SaldoID' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TotalBalance' failed on the 'required' tag")
}

func TestUpdateSaldoBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	request := requests.UpdateSaldoBalance{
		CardNumber:   "1234",
		TotalBalance: 60000,
	}

	updatedSaldo := &record.SaldoRecord{
		ID:             1,
		CardNumber:     request.CardNumber,
		TotalBalance:   request.TotalBalance,
		WithdrawAmount: 0,
		WithdrawTime:   "",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      time.Now().Format(time.RFC3339),
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().UpdateSaldoBalance(&request).Return(
		updatedSaldo,
		nil,
	)

	result, err := mockRepo.UpdateSaldoBalance(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedSaldo.TotalBalance, result.TotalBalance)
	assert.Equal(t, updatedSaldo.CardNumber, result.CardNumber)
}

func TestUpdateSaldoBalance_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	request := requests.UpdateSaldoBalance{
		CardNumber:   "1234",
		TotalBalance: 100000,
	}

	mockRepo.EXPECT().UpdateSaldoBalance(&request).Return(nil, fmt.Errorf("saldo not found"))

	result, err := mockRepo.UpdateSaldoBalance(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "saldo not found")
}

func TestUpdateSaldoBalance_ValidationError(t *testing.T) {
	request := requests.UpdateSaldoBalance{
		CardNumber:   "",
		TotalBalance: 40000,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TotalBalance' failed on the 'min' tag")
}

func TestUpdateSaldoWithdraw_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	now := time.Now()
	withdrawAmount := 500
	request := requests.UpdateSaldoWithdraw{
		CardNumber:     "1234",
		TotalBalance:   500000,
		WithdrawAmount: &withdrawAmount,
		WithdrawTime:   &now,
	}

	updatedSaldo := &record.SaldoRecord{
		ID:             1,
		CardNumber:     "1234",
		TotalBalance:   500,
		WithdrawAmount: 500,
		WithdrawTime:   now.Format(time.RFC3339),
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      time.Now().Format(time.RFC3339),
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().UpdateSaldoWithdraw(&request).Return(
		updatedSaldo,
		nil,
	)

	result, err := mockRepo.UpdateSaldoWithdraw(&request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedSaldo.TotalBalance, result.TotalBalance)
	assert.Equal(t, updatedSaldo.WithdrawAmount, result.WithdrawAmount)
}

func TestUpdateSaldoWithdraw_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	now := time.Now()
	withdrawAmount := 500
	request := requests.UpdateSaldoWithdraw{
		CardNumber:     "1234",
		TotalBalance:   500000,
		WithdrawAmount: &withdrawAmount,
		WithdrawTime:   &now,
	}

	mockRepo.EXPECT().UpdateSaldoWithdraw(&request).Return(nil, fmt.Errorf("failed to update saldo"))

	result, err := mockRepo.UpdateSaldoWithdraw(&request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update saldo")
}

func TestUpdateSaldoWithdraw_ValidationError(t *testing.T) {
	now := time.Now()
	invalidWithdrawAmount := -100
	request := requests.UpdateSaldoWithdraw{
		CardNumber:     "",
		TotalBalance:   40000,
		WithdrawAmount: &invalidWithdrawAmount,
		WithdrawTime:   &now,
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'CardNumber' failed on the 'required' tag")
	assert.Contains(t, err.Error(), "Field validation for 'TotalBalance' failed on the 'min' tag")
	assert.Contains(t, err.Error(), "Field validation for 'WithdrawAmount' failed on the 'gte' tag")
}

func TestTrashedSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoID := 1
	trashedSaldo := &record.SaldoRecord{
		ID:             saldoID,
		CardNumber:     "1234",
		TotalBalance:   500000,
		WithdrawAmount: 0,
		WithdrawTime:   "",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-25T10:00:00Z",
		DeletedAt:      utils.PtrString("2024-12-25T10:00:00Z"),
	}

	mockRepo.EXPECT().TrashedSaldo(saldoID).Return(trashedSaldo, nil)

	result, err := mockRepo.TrashedSaldo(saldoID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, trashedSaldo.ID, result.ID)
	assert.NotNil(t, result.DeletedAt)
}

func TestTrashedSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoID := 1

	mockRepo.EXPECT().TrashedSaldo(saldoID).Return(nil, fmt.Errorf("failed to trash saldo"))

	result, err := mockRepo.TrashedSaldo(saldoID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to trash saldo")
}

func TestRestoreSaldo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoID := 1
	restoredSaldo := &record.SaldoRecord{
		ID:             saldoID,
		CardNumber:     "1234",
		TotalBalance:   500000,
		WithdrawAmount: 0,
		WithdrawTime:   "",
		CreatedAt:      "2024-12-24T10:00:00Z",
		UpdatedAt:      "2024-12-25T10:00:00Z",
		DeletedAt:      nil,
	}

	mockRepo.EXPECT().RestoreSaldo(saldoID).Return(restoredSaldo, nil)

	result, err := mockRepo.RestoreSaldo(saldoID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, restoredSaldo.ID, result.ID)
	assert.Nil(t, result.DeletedAt)
}

func TestRestoreSaldo_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoID := 1

	mockRepo.EXPECT().RestoreSaldo(saldoID).Return(nil, fmt.Errorf("failed to restore saldo"))

	result, err := mockRepo.RestoreSaldo(saldoID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to restore saldo")
}

func TestDeleteSaldoPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoID := 1

	mockRepo.EXPECT().DeleteSaldoPermanent(saldoID).Return(nil)

	err := mockRepo.DeleteSaldoPermanent(saldoID)

	assert.NoError(t, err)
}

func TestDeleteSaldoPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSaldoRepository(ctrl)

	saldoID := 1

	mockRepo.EXPECT().DeleteSaldoPermanent(saldoID).Return(fmt.Errorf("failed to delete saldo permanently"))

	err := mockRepo.DeleteSaldoPermanent(saldoID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete saldo permanently")
}
