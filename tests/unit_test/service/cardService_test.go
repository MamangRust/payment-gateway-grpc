package test

import (
	"errors"
	"testing"

	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type CardServiceSuite struct {
	suite.Suite
	mockCardRepo *mocks.MockCardRepository
	mockUserRepo *mocks.MockUserRepository
	mockLogger   *mock_logger.MockLoggerInterface
	mockMapper   *mock_responseservice.MockCardResponseMapper
	mockCtrl     *gomock.Controller
}

func (suite *CardServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockCardRepo = mocks.NewMockCardRepository(suite.mockCtrl)
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockCardResponseMapper(suite.mockCtrl)
}

func (suite *CardServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *CardServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllCards{Search: "debit", Page: 1, PageSize: 10}
	total := 1
	mockCards := []*record.CardRecord{{ID: 1, CardNumber: "1111-xxxx-xxxx-1111"}}
	expectedResponse := []*response.CardResponse{{ID: 1, CardNumber: "1111-xxxx-xxxx-1111"}}

	suite.mockLogger.EXPECT().Debug("Fetching card records", gomock.Any())
	suite.mockCardRepo.EXPECT().FindAllCards(req).Return(mockCards, &total, nil)
	suite.mockMapper.EXPECT().ToCardsResponse(mockCards).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched card records", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *CardServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllCards{Search: "debit", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching card records", gomock.Any())
	suite.mockCardRepo.EXPECT().FindAllCards(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch card", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(card_errors.ErrFailedFindAllCards, err)
}

func (suite *CardServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllCards{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockCards := []*record.CardRecord{{ID: 1, DeletedAt: nil}}
	expectedResponse := []*response.CardResponseDeleteAt{{ID: 1, DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active card records", gomock.Any())
	suite.mockCardRepo.EXPECT().FindByActive(req).Return(mockCards, &total, nil)
	suite.mockMapper.EXPECT().ToCardsResponseDeleteAt(mockCards).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active card records", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *CardServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllCards{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active card records", gomock.Any())
	suite.mockCardRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch active card records", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(card_errors.ErrFailedFindActiveCards, err)
}

func (suite *CardServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllCards{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockCards := []*record.CardRecord{{ID: 2, DeletedAt: &trashedTime}}
	expectedResponse := []*response.CardResponseDeleteAt{{ID: 2, DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed card records", gomock.Any())
	suite.mockCardRepo.EXPECT().FindByTrashed(req).Return(mockCards, &total, nil)
	suite.mockMapper.EXPECT().ToCardsResponseDeleteAt(mockCards).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed card records", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *CardServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllCards{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed card records", gomock.Any())
	suite.mockCardRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed card records", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(card_errors.ErrFailedFindTrashedCards, err)
}

func (suite *CardServiceSuite) TestFindById_Success() {
	cardID := 1
	mockCard := &record.CardRecord{ID: cardID, CardNumber: "1111-xxxx-xxxx-1111"}
	expectedResponse := &response.CardResponse{ID: cardID, CardNumber: "1111-xxxx-xxxx-1111"}

	suite.mockLogger.EXPECT().Debug("Fetching card by ID", zap.Int("card_id", cardID))
	suite.mockCardRepo.EXPECT().FindById(cardID).Return(mockCard, nil)
	suite.mockMapper.EXPECT().ToCardResponse(mockCard).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched card", zap.Int("card_id", cardID))

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(cardID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *CardServiceSuite) TestFindById_NotFound() {
	cardID := 99
	repoError := errors.New("card not found")

	suite.mockLogger.EXPECT().Debug("Fetching card by ID", zap.Int("card_id", cardID))
	suite.mockCardRepo.EXPECT().FindById(cardID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve Card details", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(cardID)

	suite.Nil(result)
	suite.Equal(card_errors.ErrFailedFindById, err)
}

func (suite *CardServiceSuite) TestFindByUserID_Success() {
	userID := 123
	mockCard := &record.CardRecord{ID: 1, UserID: userID}
	expectedResponse := &response.CardResponse{ID: 1, UserID: userID}

	suite.mockLogger.EXPECT().Debug("Fetching card by user ID", zap.Int("user_id", userID))
	suite.mockCardRepo.EXPECT().FindCardByUserId(userID).Return(mockCard, nil)
	suite.mockMapper.EXPECT().ToCardResponse(mockCard).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched card records by user ID", zap.Int("user_id", userID))

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByUserID(userID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *CardServiceSuite) TestFindByUserID_NotFound() {
	userID := 999
	repoError := errors.New("card not found")

	suite.mockLogger.EXPECT().Debug("Fetching card by user ID", zap.Int("user_id", userID))
	suite.mockCardRepo.EXPECT().FindCardByUserId(userID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve Card details by user", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByUserID(userID)

	suite.Nil(result)
	suite.Equal(card_errors.ErrFailedFindByUserID, err)
}

func (suite *CardServiceSuite) TestFindByCardNumber_Success() {
	cardNumber := "1234-5678-9012-3456"
	mockCard := &record.CardRecord{CardNumber: cardNumber}
	expectedResponse := &response.CardResponse{CardNumber: cardNumber}

	suite.mockLogger.EXPECT().Debug("Fetching card record by card number", zap.String("card_number", cardNumber))
	suite.mockCardRepo.EXPECT().FindCardByCardNumber(cardNumber).Return(mockCard, nil)
	suite.mockMapper.EXPECT().ToCardResponse(mockCard).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched card record by card number", zap.String("card_number", cardNumber))

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByCardNumber(cardNumber)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *CardServiceSuite) TestFindByCardNumber_NotFound() {
	cardNumber := "not-found-xxxx"
	repoError := errors.New("card not found")

	suite.mockLogger.EXPECT().Debug("Fetching card record by card number", zap.String("card_number", cardNumber))
	suite.mockCardRepo.EXPECT().FindCardByCardNumber(cardNumber).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve find card", gomock.Any())

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindByCardNumber(cardNumber)

	suite.Nil(result)
	suite.Equal(card_errors.ErrCardNotFoundRes, err)
}

func (suite *CardServiceSuite) TestTrashedCard_Success() {
	cardID := 1
	trashedCard := &record.CardRecord{ID: cardID}
	expectedResponse := &response.CardResponseDeleteAt{ID: cardID}

	suite.mockLogger.EXPECT().Debug("Trashing card", zap.Int("card_id", cardID))
	suite.mockCardRepo.EXPECT().TrashedCard(cardID).Return(trashedCard, nil)
	suite.mockMapper.EXPECT().ToCardResponseDeleteAt(trashedCard).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully trashed card", zap.Int("card_id", cardID))

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedCard(cardID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *CardServiceSuite) TestTrashedCard_Failure() {
	cardID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Trashing card", zap.Int("card_id", cardID))
	suite.mockLogger.EXPECT().Error("Failed to move card to trash", gomock.Any())
	suite.mockCardRepo.EXPECT().TrashedCard(cardID).Return(nil, dbError)

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedCard(cardID)

	suite.Nil(result)
	suite.Equal(card_errors.ErrFailedTrashCard, err)
}

func (suite *CardServiceSuite) TestRestoreCard_Success() {
	cardID := 1
	restoredCard := &record.CardRecord{ID: cardID, DeletedAt: nil}
	expectedResponse := &response.CardResponseDeleteAt{ID: cardID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Restoring card", zap.Int("card_id", cardID))
	suite.mockCardRepo.EXPECT().RestoreCard(cardID).Return(restoredCard, nil)
	suite.mockMapper.EXPECT().ToCardResponseDeleteAt(restoredCard).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully restored card", zap.Int("card_id", cardID))

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreCard(cardID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *CardServiceSuite) TestRestoreCard_Failure() {
	cardID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Restoring card", zap.Int("card_id", cardID))
	suite.mockLogger.EXPECT().Error("Failed to restore cashier from trash", gomock.Any())
	suite.mockCardRepo.EXPECT().RestoreCard(cardID).Return(nil, dbError)

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreCard(cardID)

	suite.Nil(result)
	suite.Equal(card_errors.ErrFailedRestoreCard, err)
}

func (suite *CardServiceSuite) TestDeleteCardPermanent_Success() {
	cardID := 1

	suite.mockLogger.EXPECT().Debug("Permanently deleting card", zap.Int("card_id", cardID))
	suite.mockCardRepo.EXPECT().DeleteCardPermanent(cardID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted card permanently", zap.Int("card_id", cardID))

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteCardPermanent(cardID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *CardServiceSuite) TestDeleteCardPermanent_Failure() {
	cardID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting card", zap.Int("card_id", cardID))
	suite.mockLogger.EXPECT().Error("Failed to permanently delete card", gomock.Any())
	suite.mockCardRepo.EXPECT().DeleteCardPermanent(cardID).Return(false, dbError)

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteCardPermanent(cardID)

	suite.False(result)
	suite.Equal(card_errors.ErrFailedDeleteCard, err)
}

func (suite *CardServiceSuite) TestRestoreAllCard_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all cards")
	suite.mockCardRepo.EXPECT().RestoreAllCard().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all cards")

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllCard()

	suite.Nil(err)
	suite.True(result)
}

func (suite *CardServiceSuite) TestRestoreAllCard_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all cards")
	suite.mockLogger.EXPECT().Error("Failed to restore all trashed cards", gomock.Any())
	suite.mockCardRepo.EXPECT().RestoreAllCard().Return(false, dbError)

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllCard()

	suite.False(result)
	suite.Equal(card_errors.ErrFailedRestoreAllCards, err)
}

func (suite *CardServiceSuite) TestDeleteAllCardPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all cards")
	suite.mockCardRepo.EXPECT().DeleteAllCardPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all cards permanently")

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllCardPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *CardServiceSuite) TestDeleteAllCardPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all cards")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all trashed card", gomock.Any())
	suite.mockCardRepo.EXPECT().DeleteAllCardPermanent().Return(false, dbError)

	svc := service.NewCardService(suite.mockCardRepo, suite.mockUserRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllCardPermanent()

	suite.False(result)
	suite.Equal(card_errors.ErrFailedDeleteAllCards, err)
}

func TestCardServiceSuite(t *testing.T) {
	suite.Run(t, new(CardServiceSuite))
}
