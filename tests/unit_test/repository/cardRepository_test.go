package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/tests/utils"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllCards_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	cards := []*record.CardRecord{
		{
			ID:           1,
			UserID:       1,
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
		},
		{
			ID:           2,
			UserID:       2,
			CardType:     "Credit",
			ExpireDate:   "2024-06-30",
			CVV:          "456",
			CardProvider: "MasterCard",
		},
	}
	page := 1
	pageSize := 10
	search := ""

	mockRepo.EXPECT().FindAllCards(search, page, pageSize).Return(cards, 2, nil)

	result, total, err := mockRepo.FindAllCards(search, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, cards, result)
	assert.Equal(t, 2, total)
}

func TestFindAllCards_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	mockRepo.EXPECT().FindAllCards("", 1, 10).Return(nil, 0, errors.New("database error"))

	result, total, err := mockRepo.FindAllCards("", 1, 10)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	assert.EqualError(t, err, "database error")
}

func TestFindAllCards_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	mockRepo.EXPECT().FindAllCards("", 1, 10).Return([]*record.CardRecord{}, 0, nil)

	result, total, err := mockRepo.FindAllCards("", 1, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, total)
}

func TestFindById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	card := &record.CardRecord{
		ID:           1,
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
	}

	mockRepo.EXPECT().FindById(1).Return(card, nil)

	result, err := mockRepo.FindById(1)
	assert.NoError(t, err)
	assert.Equal(t, card, result)
}

func TestFindById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	mockRepo.EXPECT().FindById(1).Return(nil, errors.New("card not found"))

	result, err := mockRepo.FindById(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "card not found")
}

func TestFindCardByUserId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	userID := 1
	expectedCard := &record.CardRecord{
		ID:           1,
		UserID:       userID,
		CardType:     "Credit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
	}

	mockRepo.EXPECT().FindCardByUserId(userID).Return(expectedCard, nil)

	result, err := mockRepo.FindCardByUserId(userID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard, result)
}

func TestFindCardByUserId_Failure_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	userID := 1

	mockRepo.EXPECT().FindCardByUserId(userID).Return(nil, errors.New("card not found"))

	result, err := mockRepo.FindCardByUserId(userID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "card not found")
}

func TestFindByActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	expectedCards := []*record.CardRecord{
		{
			ID:           1,
			UserID:       1,
			CardType:     "Debit",
			ExpireDate:   "2025-12-31",
			CVV:          "123",
			CardProvider: "Visa",
		},
		{
			ID:           2,
			UserID:       2,
			CardType:     "Credit",
			ExpireDate:   "2026-01-15",
			CVV:          "456",
			CardProvider: "MasterCard",
		},
	}
	page := 1
	pageSize := 10
	search := ""
	expected := 2

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(expectedCards, 2, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCards, result)
}

func TestFindByActive_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return(nil, 0, errors.New("database error"))

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expected, totalRecord)
	assert.EqualError(t, err, "database error")
}

func TestFindByActive_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByActive(search, page, pageSize).Return([]*record.CardRecord{}, 0, nil)

	result, totalRecord, err := mockRepo.FindByActive(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestFindByTrashed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	expectedCards := []*record.CardRecord{
		{
			ID:           3,
			UserID:       3,
			CardType:     "Debit",
			ExpireDate:   "2024-11-30",
			CVV:          "789",
			CardProvider: "Visa",
		},
	}
	page := 1
	pageSize := 10
	search := ""
	expected := 1

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(expectedCards, 1, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, expected, totalRecord)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCards, result)
}

func TestFindByTrashed_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return(nil, 0, errors.New("database error"))

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

func TestFindByTrashed_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	page := 1
	pageSize := 10
	search := ""
	expected := 0

	mockRepo.EXPECT().FindByTrashed(search, page, pageSize).Return([]*record.CardRecord{}, 0, nil)

	result, totalRecord, err := mockRepo.FindByTrashed(search, page, pageSize)

	assert.Equal(t, expected, totalRecord)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestCreateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	request := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "Credit",
		ExpireDate:   time.Now().AddDate(3, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}
	card := &record.CardRecord{
		ID:           1,
		UserID:       request.UserID,
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate.String(),
		CVV:          request.CVV,
		CardProvider: request.CardProvider,
	}

	mockRepo.EXPECT().CreateCard(&request).Return(card, nil)

	result, err := mockRepo.CreateCard(&request)
	assert.NoError(t, err)
	assert.Equal(t, card, result)
}

func TestCreateCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	request := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "Credit",
		ExpireDate:   time.Now().AddDate(3, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	mockRepo.EXPECT().CreateCard(&request).Return(nil, errors.New("failed to create card"))

	result, err := mockRepo.CreateCard(&request)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to create card")
}

func TestCreateCard_ValidationError(t *testing.T) {
	request := requests.CreateCardRequest{
		UserID:       0,
		CardType:     "Credit",
		ExpireDate:   time.Now().AddDate(-1, 0, 0),
		CVV:          "12",
		CardProvider: "",
	}

	err := request.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "card type must be credit or debit")
}

func TestUpdateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	request := requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "Credit",
		ExpireDate:   time.Now().AddDate(3, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}
	updatedCard := &record.CardRecord{
		ID:           request.CardID,
		UserID:       request.UserID,
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate.String(),
		CVV:          request.CVV,
		CardProvider: request.CardProvider,
	}

	mockRepo.EXPECT().UpdateCard(&request).Return(updatedCard, nil)

	result, err := mockRepo.UpdateCard(&request)
	assert.NoError(t, err)
	assert.Equal(t, updatedCard, result)
}

func TestUpdateCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	request := requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "Credit",
		ExpireDate:   time.Now().AddDate(3, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	mockRepo.EXPECT().UpdateCard(&request).Return(nil, errors.New("failed to update card"))

	result, err := mockRepo.UpdateCard(&request)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to update card")

}

func TestUpdateCard_ValidationError(t *testing.T) {
	request := requests.UpdateCardRequest{
		CardID:       0,
		UserID:       0,
		CardType:     "Credit",
		ExpireDate:   time.Now().AddDate(-1, 0, 0),
		CVV:          "12",
		CardProvider: "",
	}

	err := request.Validate()

	assert.Error(t, err)

	assert.Contains(t, err.Error(), "card type must be credit or debit")

}

func TestTrashedCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	cardID := 1
	expectedCard := &record.CardRecord{
		ID:           cardID,
		UserID:       1,
		CardType:     "Debit",
		ExpireDate:   "2025-12-31",
		CVV:          "123",
		CardProvider: "Visa",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
		DeletedAt:    utils.PtrString("2024-12-21T09:00:00Z"),
	}

	mockRepo.EXPECT().TrashedCard(cardID).Return(expectedCard, nil)

	result, err := mockRepo.TrashedCard(cardID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard, result)
}

func TestTrashedCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	cardID := 1

	mockRepo.EXPECT().TrashedCard(cardID).Return(nil, errors.New("failed to trash card"))

	result, err := mockRepo.TrashedCard(cardID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to trash card")
}

func TestRestoreCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	cardID := 1
	expectedCard := &record.CardRecord{
		ID:           cardID,
		UserID:       1,
		CardType:     "Credit",
		ExpireDate:   "2026-01-15",
		CVV:          "456",
		CardProvider: "MasterCard",
		CreatedAt:    "2024-12-21T09:00:00Z",
		UpdatedAt:    "2024-12-21T09:00:00Z",
		DeletedAt:    nil,
	}

	mockRepo.EXPECT().RestoreCard(cardID).Return(expectedCard, nil)

	result, err := mockRepo.RestoreCard(cardID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard, result)
}

func TestDeleteCardPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	mockRepo.EXPECT().DeleteCardPermanent(1).Return(nil)

	err := mockRepo.DeleteCardPermanent(1)
	assert.NoError(t, err)
}

func TestDeleteCardPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCardRepository(ctrl)

	mockRepo.EXPECT().DeleteCardPermanent(1).Return(errors.New("delete failed"))

	err := mockRepo.DeleteCardPermanent(1)
	assert.Error(t, err)
	assert.EqualError(t, err, "delete failed")
}
