package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	mock_apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionHandlerTestSuite struct {
	suite.Suite
	Ctrl                  *gomock.Controller
	MockMerchantClient    *mock_pb.MockMerchantServiceClient
	MockTransactionClient *mock_pb.MockTransactionServiceClient
	MockLogger            *mock_logger.MockLoggerInterface
	MockMapper            *mock_apimapper.MockTransactionResponseMapper
	E                     *echo.Echo
	Handler               *api.TransactionHandleApi
}

func (suite *TransactionHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockTransactionClient = mock_pb.NewMockTransactionServiceClient(suite.Ctrl)
	suite.MockMerchantClient = mock_pb.NewMockMerchantServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockTransactionResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerTransaction(suite.MockTransactionClient, suite.MockMerchantClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *TransactionHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *TransactionHandlerTestSuite) TestFindAllTransaction_Success() {
	grpcResponse := &pb.ApiResponsePaginationTransaction{
		Status:  "success",
		Message: "Transactions retrieved successfully",
		Data: []*pb.TransactionResponse{
			{Id: 1, TransactionNo: "TX123", CardNumber: "1234", Amount: 100000},
			{Id: 2, TransactionNo: "TX124", CardNumber: "5678", Amount: 200000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationTransaction{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.TransactionResponse{{ID: 1, TransactionNo: "TX123", CardNumber: "1234", Amount: 100000}, {ID: 2, TransactionNo: "TX124", CardNumber: "5678", Amount: 200000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockTransactionClient.EXPECT().FindAllTransaction(gomock.Any(), &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: "test"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTransaction(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions?page=1&page_size=10&search=test", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *TransactionHandlerTestSuite) TestFindAllTransaction_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().FindAllTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transaction data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestFindAllTransactionByCardNumber_Success() {
	cardNumber := "1234567890123456"
	grpcResponse := &pb.ApiResponsePaginationTransaction{
		Status:  "success",
		Message: "Transactions for card retrieved successfully",
		Data: []*pb.TransactionResponse{
			{Id: 1, TransactionNo: "TX123", CardNumber: cardNumber, Amount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 1, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationTransaction{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.TransactionResponse{{ID: 1, TransactionNo: "TX123", CardNumber: cardNumber, Amount: 100000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 1, TotalPages: 1},
	}

	suite.MockTransactionClient.EXPECT().FindAllTransactionByCardNumber(gomock.Any(), &pb.FindAllTransactionCardNumberRequest{CardNumber: cardNumber, Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTransaction(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transactions/card-number/%s?page=1&page_size=10", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindAllTransactionByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TransactionHandlerTestSuite) TestFindAllTransactionByCardNumber_Failure() {
	cardNumber := "1234567890123456"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().FindAllTransactionByCardNumber(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transaction data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/transactions/card-number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindAllTransactionByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestFindByIdTransaction_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTransaction{Status: "success", Data: &pb.TransactionResponse{Id: int32(id), TransactionNo: "TX123", CardNumber: "1234", Amount: 100000}}
	expectedApiResponse := &response.ApiResponseTransaction{Status: "success", Data: &response.TransactionResponse{ID: id, TransactionNo: "TX123", CardNumber: "1234", Amount: 100000}}

	suite.MockTransactionClient.EXPECT().FindByIdTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{TransactionId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransaction(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *TransactionHandlerTestSuite) TestFindByIdTransaction_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/transactions/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	suite.MockLogger.EXPECT().
		Debug("Invalid transaction ID", gomock.Any()).
		Times(1)

	err := suite.Handler.FindById(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestFindTransactionByMerchantId_Success() {
	merchantId := 1

	grpcResponse := &pb.ApiResponseTransactions{
		Status: "success",
		Data: []*pb.TransactionResponse{
			{Id: 1, TransactionNo: "TX123", MerchantId: int32(merchantId)},
		},
	}

	expectedApiResponse := &response.ApiResponseTransactions{
		Status: "success",
		Data: []*response.TransactionResponse{
			{ID: 1, TransactionNo: "TX123", MerchantID: merchantId},
		},
	}

	suite.MockTransactionClient.EXPECT().
		FindTransactionByMerchantId(
			gomock.Any(),
			&pb.FindTransactionByMerchantIdRequest{MerchantId: int32(merchantId)},
		).
		Return(grpcResponse, nil)

	suite.MockMapper.EXPECT().
		ToApiResponseTransactions(grpcResponse).
		Return(expectedApiResponse)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/transactions/merchant/1",
		nil,
	)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("merchant_id")
	c.SetParamValues("1")

	err := suite.Handler.FindByTransactionMerchantId(c)

	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransactions
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TransactionHandlerTestSuite) TestFindTransactionByMerchantId_Failure() {
	merchantId := 1
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().FindTransactionByMerchantId(gomock.Any(), &pb.FindTransactionByMerchantIdRequest{MerchantId: int32(merchantId)}).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transaction data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/merchant/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("merchant_id")
	c.SetParamValues("1")

	err := suite.Handler.FindByTransactionMerchantId(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestFindByActiveTransaction_Success() {
	grpcResponse := &pb.ApiResponsePaginationTransactionDeleteAt{Status: "success", Data: []*pb.TransactionResponseDeleteAt{{Id: 1, TransactionNo: "TX123"}}}
	expectedApiResponse := &response.ApiResponsePaginationTransactionDeleteAt{Status: "success", Data: []*response.TransactionResponseDeleteAt{{ID: 1, TransactionNo: "TX123"}}}

	suite.MockTransactionClient.EXPECT().FindByActiveTransaction(gomock.Any(), &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTransactionDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActiveTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransactionDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TransactionHandlerTestSuite) TestFindByActiveTransaction_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().FindByActiveTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transaction data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActiveTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestFindByTrashedTransaction_Success() {
	grpcResponse := &pb.ApiResponsePaginationTransactionDeleteAt{Status: "success", Data: []*pb.TransactionResponseDeleteAt{{Id: 2, TransactionNo: "TX999"}}}
	expectedApiResponse := &response.ApiResponsePaginationTransactionDeleteAt{Status: "success", Data: []*response.TransactionResponseDeleteAt{{ID: 2, TransactionNo: "TX999"}}}

	suite.MockTransactionClient.EXPECT().FindByTrashedTransaction(gomock.Any(), &pb.FindAllTransactionRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTransactionDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashedTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransactionDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TransactionHandlerTestSuite) TestFindByTrashedTransaction_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().FindByTrashedTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transaction data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashedTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestCreateTransaction_Success() {
	transactionTime := time.Now()
	id := 1
	apiKey := "test-api-key"

	requestBody := requests.CreateTransactionRequest{
		CardNumber:      "1234567890123456",
		Amount:          150000,
		PaymentMethod:   "alfamart",
		MerchantID:      &id,
		TransactionTime: transactionTime,
	}

	grpcRequest := &pb.CreateTransactionRequest{
		ApiKey:          apiKey,
		CardNumber:      requestBody.CardNumber,
		Amount:          int32(requestBody.Amount),
		PaymentMethod:   requestBody.PaymentMethod,
		MerchantId:      int32(*requestBody.MerchantID),
		TransactionTime: timestamppb.New(transactionTime),
	}

	grpcResponse := &pb.ApiResponseTransaction{
		Status: "success",
		Data: &pb.TransactionResponse{
			Id:            3,
			TransactionNo: "TX125",
			CardNumber:    requestBody.CardNumber,
			Amount:        int32(requestBody.Amount),
		},
	}

	expectedApiResponse := &response.ApiResponseTransaction{
		Status: "success",
		Data: &response.TransactionResponse{
			ID:            3,
			TransactionNo: "TX125",
			CardNumber:    requestBody.CardNumber,
			Amount:        requestBody.Amount,
		},
	}

	suite.MockTransactionClient.EXPECT().CreateTransaction(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransaction(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("apiKey", apiKey)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("TX125", resp.Data.TransactionNo)
}

func (suite *TransactionHandlerTestSuite) TestCreateTransaction_Failure() {
	id := 1
	apiKey := "test-api-key"

	requestBody := requests.CreateTransactionRequest{
		CardNumber:      "1234567890123456",
		Amount:          150000,
		PaymentMethod:   "alfamart",
		MerchantID:      &id,
		TransactionTime: time.Now(),
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create transaction", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("apiKey", apiKey)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestCreateTransaction_ValidationError() {
	apiKey := "test-api-key"

	invalidRequestBody := requests.CreateTransactionRequest{
		CardNumber: "",
		Amount:     40000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("apiKey", apiKey)

	suite.MockLogger.EXPECT().
		Debug("Validation Error", gomock.Any()).
		Times(1)

	err := suite.Handler.Create(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestUpdateTransaction_Success() {
	id := 1
	apiKey := "test-api-key"
	transactionTime := time.Now()

	requestBody := requests.UpdateTransactionRequest{
		CardNumber:      "1234567890123456",
		Amount:          200000,
		PaymentMethod:   "alfamart",
		MerchantID:      &id,
		TransactionTime: transactionTime,
	}

	grpcRequest := &pb.UpdateTransactionRequest{
		TransactionId:   int32(id),
		CardNumber:      requestBody.CardNumber,
		ApiKey:          apiKey,
		Amount:          int32(requestBody.Amount),
		PaymentMethod:   requestBody.PaymentMethod,
		MerchantId:      int32(*requestBody.MerchantID),
		TransactionTime: timestamppb.New(transactionTime),
	}

	grpcResponse := &pb.ApiResponseTransaction{
		Status: "success",
		Data: &pb.TransactionResponse{
			Id:     int32(id),
			Amount: int32(requestBody.Amount),
		},
	}
	expectedApiResponse := &response.ApiResponseTransaction{
		Status: "success",
		Data: &response.TransactionResponse{
			ID:     id,
			Amount: requestBody.Amount,
		},
	}

	suite.MockTransactionClient.EXPECT().UpdateTransaction(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransaction(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transactions/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	c.Set("apiKey", apiKey)

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransaction
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(200000), resp.Data.Amount)
}

func (suite *TransactionHandlerTestSuite) TestUpdateTransaction_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/update/abc", nil)
	suite.MockLogger.EXPECT().
		Debug("Bad Request", gomock.Any()).
		Times(1)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestUpdateTransaction_ValidationError() {
	id := 1
	apiKey := "test-api-key"

	invalidRequestBody := requests.UpdateTransactionRequest{
		CardNumber: "",
		Amount:     40000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/transactions/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	c.Set("apiKey", apiKey)

	suite.MockLogger.EXPECT().
		Debug("Validation Error", gomock.Any()).
		Times(1)

	err := suite.Handler.Update(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestTrashedTransaction_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTransactionDeleteAt{Status: "success", Data: &pb.TransactionResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseTransactionDeleteAt{Status: "success", Data: &response.TransactionResponseDeleteAt{ID: id}}

	suite.MockTransactionClient.EXPECT().TrashedTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{TransactionId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransactionDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestTrashedTransaction_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().TrashedTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trashed transaction", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestRestoreTransaction_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTransactionDeleteAt{Status: "success", Data: &pb.TransactionResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseTransactionDeleteAt{Status: "success", Data: &response.TransactionResponseDeleteAt{ID: id}}

	suite.MockTransactionClient.EXPECT().RestoreTransaction(gomock.Any(), &pb.FindByIdTransactionRequest{TransactionId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransactionDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestRestoreTransaction_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().RestoreTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore transaction", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestDeleteTransactionPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTransactionDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseTransactionDelete{Status: "success"}

	suite.MockTransactionClient.EXPECT().DeleteTransactionPermanent(gomock.Any(), &pb.FindByIdTransactionRequest{TransactionId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransactionDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeletePermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestDeleteTransactionPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().DeleteTransactionPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete transaction", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeletePermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestRestoreAllTransaction_Success() {
	grpcResponse := &pb.ApiResponseTransactionAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseTransactionAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully restored all transaction").Times(1)

	suite.MockTransactionClient.EXPECT().RestoreAllTransaction(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransactionAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestRestoreAllTransaction_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().RestoreAllTransaction(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all transaction", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllTransaction(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestDeleteAllTransactionPermanent_Success() {
	grpcResponse := &pb.ApiResponseTransactionAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseTransactionAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully deleted all transaction permanently").Times(1)

	suite.MockTransactionClient.EXPECT().DeleteAllTransactionPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTransactionAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllTransactionPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TransactionHandlerTestSuite) TestDeleteAllTransactionPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTransactionClient.EXPECT().DeleteAllTransactionPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all transaction", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/transactions/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllTransactionPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestTransactionHandlerSuite(t *testing.T) {
	suite.Run(t, new(TransactionHandlerTestSuite))
}
