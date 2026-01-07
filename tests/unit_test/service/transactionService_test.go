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
	"MamangRust/paymentgatewaygrpc/pkg/errors/transaction_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type TransactionServiceSuite struct {
	suite.Suite
	mockMerchantRepo    *mocks.MockMerchantRepository
	mockCardRepo        *mocks.MockCardRepository
	mockSaldoRepo       *mocks.MockSaldoRepository
	mockTransactionRepo *mocks.MockTransactionRepository
	mockLogger          *mock_logger.MockLoggerInterface
	mockMapper          *mock_responseservice.MockTransactionResponseMapper
	mockCtrl            *gomock.Controller
}

func (suite *TransactionServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockMerchantRepo = mocks.NewMockMerchantRepository(suite.mockCtrl)
	suite.mockCardRepo = mocks.NewMockCardRepository(suite.mockCtrl)
	suite.mockSaldoRepo = mocks.NewMockSaldoRepository(suite.mockCtrl)
	suite.mockTransactionRepo = mocks.NewMockTransactionRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockTransactionResponseMapper(suite.mockCtrl)
}

func (suite *TransactionServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TransactionServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllTransactions{Search: "success", Page: 1, PageSize: 10}
	total := 1
	mockTransactions := []*record.TransactionRecord{{ID: 1, TransactionNo: "TX001"}}
	expectedResponse := []*response.TransactionResponse{{ID: 1, TransactionNo: "TX001"}}

	suite.mockLogger.EXPECT().Debug("Fetching transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindAllTransactions(req).Return(mockTransactions, &total, nil)
	suite.mockMapper.EXPECT().ToTransactionsResponse(mockTransactions).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransactionServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllTransactions{Search: "success", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindAllTransactions(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transaction_errors.ErrFailedFindAllTransactions, err)
}

func (suite *TransactionServiceSuite) TestFindAllByCardNumber_Success() {
	req := &requests.FindAllTransactionCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	total := 1
	mockTransactions := []*record.TransactionRecord{{ID: 2, CardNumber: "1111-xxxx"}}
	expectedResponse := []*response.TransactionResponse{{ID: 2, CardNumber: "1111-xxxx"}}

	suite.mockLogger.EXPECT().Debug("Fetching transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindAllTransactionByCardNumber(req).Return(mockTransactions, &total, nil)
	suite.mockMapper.EXPECT().ToTransactionsResponse(mockTransactions).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAllByCardNumber(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransactionServiceSuite) TestFindAllByCardNumber_Failure() {
	req := &requests.FindAllTransactionCardNumber{CardNumber: "1111-xxxx", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindAllTransactionByCardNumber(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindAllByCardNumber(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transaction_errors.ErrFailedFindAllByCardNumber, err)
}

func (suite *TransactionServiceSuite) TestFindById_Success() {
	transactionID := 1
	mockTransaction := &record.TransactionRecord{ID: transactionID, TransactionNo: "TX001"}
	expectedResponse := &response.TransactionResponse{ID: transactionID, TransactionNo: "TX001"}

	suite.mockLogger.EXPECT().Debug("Fetching transaction by ID", zap.Int("transaction_id", transactionID))
	suite.mockTransactionRepo.EXPECT().FindById(transactionID).Return(mockTransaction, nil)
	suite.mockMapper.EXPECT().ToTransactionResponse(mockTransaction).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(transactionID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransactionServiceSuite) TestFindById_NotFound() {
	transactionID := 99
	repoError := errors.New("transaction not found")

	suite.mockLogger.EXPECT().Debug("Fetching transaction by ID", zap.Int("transaction_id", transactionID))
	suite.mockTransactionRepo.EXPECT().FindById(transactionID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("failed to find transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindById(transactionID)

	suite.Nil(result)
	suite.Equal(transaction_errors.ErrTransactionNotFound, err)
}

func (suite *TransactionServiceSuite) TestFindByActive_Success() {
	req := &requests.FindAllTransactions{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockTransactions := []*record.TransactionRecord{{ID: 1, TransactionNo: "TX001", DeletedAt: nil}}
	expectedResponse := []*response.TransactionResponseDeleteAt{{ID: 1, TransactionNo: "TX001", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindByActive(req).Return(mockTransactions, &total, nil)
	suite.mockMapper.EXPECT().ToTransactionsResponseDeleteAt(mockTransactions).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransactionServiceSuite) TestFindByActive_Failure() {
	req := &requests.FindAllTransactions{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindByActive(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch active transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByActive(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transaction_errors.ErrFailedFindByActiveTransactions, err)
}

func (suite *TransactionServiceSuite) TestFindByTrashed_Success() {
	req := &requests.FindAllTransactions{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockTransactions := []*record.TransactionRecord{{ID: 2, TransactionNo: "TX002", DeletedAt: &trashedTime}}
	expectedResponse := []*response.TransactionResponseDeleteAt{{ID: 2, TransactionNo: "TX002", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindByTrashed(req).Return(mockTransactions, &total, nil)
	suite.mockMapper.EXPECT().ToTransactionsResponseDeleteAt(mockTransactions).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *TransactionServiceSuite) TestFindByTrashed_Failure() {
	req := &requests.FindAllTransactions{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().FindByTrashed(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed transaction", gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := svc.FindByTrashed(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(transaction_errors.ErrFailedFindByTrashedTransactions, err)
}

func (suite *TransactionServiceSuite) TestFindTransactionByMerchantId_Success() {
	merchantID := 5
	mockTransactions := []*record.TransactionRecord{{ID: 1, MerchantID: merchantID}}
	expectedResponse := []*response.TransactionResponse{{ID: 1, MerchantID: merchantID}}

	suite.mockLogger.EXPECT().Debug("Starting FindTransactionByMerchantId process", zap.Int("merchant_id", merchantID))
	suite.mockTransactionRepo.EXPECT().FindTransactionByMerchantId(merchantID).Return(mockTransactions, nil)
	suite.mockMapper.EXPECT().ToTransactionsResponse(mockTransactions).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched transaction by merchant ID", zap.Int("merchant_id", merchantID))

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindTransactionByMerchantId(merchantID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransactionServiceSuite) TestFindTransactionByMerchantId_Failure() {
	merchantID := 99
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Starting FindTransactionByMerchantId process", zap.Int("merchant_id", merchantID))
	suite.mockTransactionRepo.EXPECT().FindTransactionByMerchantId(merchantID).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch transaction by merchant ID", gomock.Any(), gomock.Any())

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.FindTransactionByMerchantId(merchantID)

	suite.Nil(result)
	suite.Equal(transaction_errors.ErrFailedFindByMerchantID, err)
}

func (suite *TransactionServiceSuite) TestTrashedTransaction_Success() {
	transactionID := 1
	trashedTransaction := &record.TransactionRecord{ID: transactionID}
	expectedResponse := &response.TransactionResponseDeleteAt{ID: transactionID}

	suite.mockLogger.EXPECT().Debug("Starting TrashedTransaction process", zap.Int("transaction_id", transactionID))
	suite.mockTransactionRepo.EXPECT().TrashedTransaction(transactionID).Return(trashedTransaction, nil)
	suite.mockMapper.EXPECT().ToTransactionResponseDeleteAt(trashedTransaction).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully trashed transaction", zap.Int("transaction_id", transactionID))

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedTransaction(transactionID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransactionServiceSuite) TestTrashedTransaction_Failure() {
	transactionID := 1
	dbError := errors.New("failed to trash")

	suite.mockLogger.EXPECT().Debug("Starting TrashedTransaction process", zap.Int("transaction_id", transactionID))
	suite.mockLogger.EXPECT().Error("Failed to move transaction to trash", gomock.Any())
	suite.mockTransactionRepo.EXPECT().TrashedTransaction(transactionID).Return(nil, dbError)

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.TrashedTransaction(transactionID)

	suite.Nil(result)
	suite.Equal(transaction_errors.ErrFailedTrashedTransaction, err)
}

func (suite *TransactionServiceSuite) TestRestoreTransaction_Success() {
	transactionID := 1
	restoredTransaction := &record.TransactionRecord{ID: transactionID, DeletedAt: nil}
	expectedResponse := &response.TransactionResponseDeleteAt{ID: transactionID, DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Starting RestoreTransaction process", zap.Int("transaction_id", transactionID))
	suite.mockTransactionRepo.EXPECT().RestoreTransaction(transactionID).Return(restoredTransaction, nil)
	suite.mockMapper.EXPECT().ToTransactionResponseDeleteAt(restoredTransaction).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully restored transaction", zap.Int("transaction_id", transactionID))

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreTransaction(transactionID)

	suite.Nil(err)
	suite.Equal(expectedResponse, result)
}

func (suite *TransactionServiceSuite) TestRestoreTransaction_Failure() {
	transactionID := 1
	dbError := errors.New("failed to restore")

	suite.mockLogger.EXPECT().Debug("Starting RestoreTransaction process", zap.Int("transaction_id", transactionID))
	suite.mockLogger.EXPECT().Error("Failed to restore transaction from trash", gomock.Any())
	suite.mockTransactionRepo.EXPECT().RestoreTransaction(transactionID).Return(nil, dbError)

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreTransaction(transactionID)

	suite.Nil(result)
	suite.Equal(transaction_errors.ErrFailedRestoreTransaction, err)
}

func (suite *TransactionServiceSuite) TestDeleteTransactionPermanent_Success() {
	transactionID := 1

	suite.mockLogger.EXPECT().Debug("Starting DeleteTransactionPermanent process", zap.Int("transaction_id", transactionID))
	suite.mockTransactionRepo.EXPECT().DeleteTransactionPermanent(transactionID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully permanently deleted transaction", zap.Int("transaction_id", transactionID))

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteTransactionPermanent(transactionID)

	suite.Nil(err)
	suite.True(result)
}

func (suite *TransactionServiceSuite) TestDeleteTransactionPermanent_Failure() {
	transactionID := 1
	dbError := errors.New("failed to delete permanently")

	suite.mockLogger.EXPECT().Debug("Starting DeleteTransactionPermanent process", zap.Int("transaction_id", transactionID))
	suite.mockLogger.EXPECT().Error("Failed to permanently delete transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().DeleteTransactionPermanent(transactionID).Return(false, dbError)

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteTransactionPermanent(transactionID)

	suite.False(result)
	suite.Equal(transaction_errors.ErrFailedDeleteTransactionPermanent, err)
}

func (suite *TransactionServiceSuite) TestRestoreAllTransaction_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all transactions")
	suite.mockTransactionRepo.EXPECT().RestoreAllTransaction().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all transactions")

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllTransaction()

	suite.Nil(err)
	suite.True(result)
}

func (suite *TransactionServiceSuite) TestRestoreAllTransaction_Failure() {
	dbError := errors.New("failed to restore all")

	suite.mockLogger.EXPECT().Debug("Restoring all transactions")
	suite.mockLogger.EXPECT().Error("Failed to restore all transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().RestoreAllTransaction().Return(false, dbError)

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.RestoreAllTransaction()

	suite.False(result)
	suite.Equal(transaction_errors.ErrFailedRestoreAllTransactions, err)
}

func (suite *TransactionServiceSuite) TestDeleteAllTransactionPermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all transactions")
	suite.mockTransactionRepo.EXPECT().DeleteAllTransactionPermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all transactions permanently")

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllTransactionPermanent()

	suite.Nil(err)
	suite.True(result)
}

func (suite *TransactionServiceSuite) TestDeleteAllTransactionPermanent_Failure() {
	dbError := errors.New("failed to delete all permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all transactions")
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all transaction", gomock.Any())
	suite.mockTransactionRepo.EXPECT().DeleteAllTransactionPermanent().Return(false, dbError)

	svc := service.NewTransactionService(suite.mockMerchantRepo, suite.mockCardRepo, suite.mockSaldoRepo, suite.mockTransactionRepo, suite.mockLogger, suite.mockMapper)
	result, err := svc.DeleteAllTransactionPermanent()

	suite.False(result)
	suite.Equal(transaction_errors.ErrFailedDeleteAllTransactionsPermanent, err)
}

func TestTransactionServiceSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceSuite))
}
