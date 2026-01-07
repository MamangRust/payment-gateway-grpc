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

type CardHandlerTestSuite struct {
	suite.Suite
	Ctrl               *gomock.Controller
	MockMerchantClient *mock_pb.MockMerchantServiceClient
	MockCardClient     *mock_pb.MockCardServiceClient
	MockLogger         *mock_logger.MockLoggerInterface
	MockMapper         *mock_apimapper.MockCardResponseMapper
	E                  *echo.Echo
	Handler            *api.CardHandleApi
}

func (suite *CardHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockCardClient = mock_pb.NewMockCardServiceClient(suite.Ctrl)
	suite.MockMerchantClient = mock_pb.NewMockMerchantServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockCardResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerCard(suite.MockCardClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *CardHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *CardHandlerTestSuite) TestFindAllCard_Success() {
	grpcResponse := &pb.ApiResponsePaginationCard{
		Status:  "success",
		Message: "Cards retrieved successfully",
		Data: []*pb.CardResponse{
			{Id: 1, UserId: 1, CardNumber: "1234567890123456", CardType: "credit", CardProvider: "visa"},
			{Id: 2, UserId: 2, CardNumber: "9876543210987654", CardType: "debit", CardProvider: "mastercard"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationCard{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.CardResponse{{ID: 1, UserID: 1, CardNumber: "1234567890123456", CardType: "credit", CardProvider: "visa"}, {ID: 2, UserID: 2, CardNumber: "9876543210987654", CardType: "debit", CardProvider: "mastercard"}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockCardClient.EXPECT().FindAllCard(gomock.Any(), &pb.FindAllCardRequest{Page: 1, PageSize: 10, Search: "test"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesCard(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/card?page=1&page_size=10&search=test", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *CardHandlerTestSuite) TestFindAllCard_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().FindAllCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch card records", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/card", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestFindByIdCard_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseCard{Status: "success", Data: &pb.CardResponse{Id: int32(id), UserId: 1, CardNumber: "1234567890123456", CardType: "credit", CardProvider: "visa"}}
	expectedApiResponse := &response.ApiResponseCard{Status: "success", Data: &response.CardResponse{ID: id, UserID: 1, CardNumber: "1234567890123456", CardType: "credit", CardProvider: "visa"}}

	suite.MockCardClient.EXPECT().FindByIdCard(gomock.Any(), &pb.FindByIdCardRequest{CardId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCard(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/card/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *CardHandlerTestSuite) TestFindByIdCard_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/card/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	suite.MockLogger.EXPECT().
		Debug("Invalid card ID", gomock.Any()).
		Times(1)

	err := suite.Handler.FindById(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CardHandlerTestSuite) TestFindByUserID_Success() {
	userId := "1"

	grpcResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Cards for user retrieved successfully",
		Data: &pb.CardResponse{
			Id:           1,
			UserId:       1,
			CardNumber:   "1234567890123456",
			CardType:     "credit",
			CardProvider: "visa",
		},
	}

	expectedApiResponse := &response.ApiResponseCard{
		Status:  grpcResponse.Status,
		Message: grpcResponse.Message,
		Data: &response.CardResponse{
			ID:           1,
			UserID:       1,
			CardNumber:   "1234567890123456",
			CardType:     "credit",
			CardProvider: "visa",
		},
	}

	suite.MockCardClient.EXPECT().
		FindByUserIdCard(
			gomock.Any(),
			&pb.FindByUserIdCardRequest{UserId: 1},
		).
		Return(grpcResponse, nil).
		Times(1)

	suite.MockMapper.EXPECT().
		ToApiResponseCard(grpcResponse).
		Return(expectedApiResponse).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/card/user",
		nil,
	)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("userID", userId)

	err := suite.Handler.FindByUserID(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
}

func (suite *CardHandlerTestSuite) TestFindByUserID_Failure() {
	userId := "1"
	grpcError := errors.New("gRPC service unavailable")

	suite.MockCardClient.EXPECT().
		FindByUserIdCard(
			gomock.Any(),
			&pb.FindByUserIdCardRequest{UserId: 1},
		).
		Return(nil, grpcError).
		Times(1)

	suite.MockLogger.EXPECT().
		Debug("Failed to fetch card record", gomock.Any()).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/card/user",
		nil,
	)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("userID", userId)

	err := suite.Handler.FindByUserID(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestFindByActiveCard_Success() {
	grpcResponse := &pb.ApiResponsePaginationCardDeleteAt{Status: "success", Data: []*pb.CardResponseDeleteAt{{Id: 1, UserId: 1, CardNumber: "1234567890123456", CardType: "credit", CardProvider: "visa"}}}
	expectedApiResponse := &response.ApiResponsePaginationCardDeleteAt{Status: "success", Data: []*response.CardResponseDeleteAt{{ID: 1, UserID: 1, CardNumber: "1234567890123456", CardType: "credit", CardProvider: "visa"}}}

	suite.MockCardClient.EXPECT().FindByActiveCard(gomock.Any(), &pb.FindAllCardRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesCardDeletedAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/card/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationCardDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *CardHandlerTestSuite) TestFindByActiveCard_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().FindByActiveCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch card record", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/card/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestFindByTrashedCard_Success() {
	grpcResponse := &pb.ApiResponsePaginationCardDeleteAt{Status: "success", Data: []*pb.CardResponseDeleteAt{{Id: 2, UserId: 2, CardNumber: "9876543210987654", CardType: "debit", CardProvider: "mastercard"}}}
	expectedApiResponse := &response.ApiResponsePaginationCardDeleteAt{Status: "success", Data: []*response.CardResponseDeleteAt{{ID: 2, UserID: 2, CardNumber: "9876543210987654", CardType: "debit", CardProvider: "mastercard"}}}

	suite.MockCardClient.EXPECT().FindByTrashedCard(gomock.Any(), &pb.FindAllCardRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesCardDeletedAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/card/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationCardDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *CardHandlerTestSuite) TestFindByTrashedCard_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().FindByTrashedCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch card record", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/card/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestFindByCardNumber_Success() {
	cardNumber := "1234567890123456"
	grpcResponse := &pb.ApiResponseCard{Status: "success", Data: &pb.CardResponse{Id: 1, UserId: 1, CardNumber: cardNumber, CardType: "credit", CardProvider: "visa"}}
	expectedApiResponse := &response.ApiResponseCard{Status: "success", Data: &response.CardResponse{ID: 1, UserID: 1, CardNumber: cardNumber, CardType: "credit", CardProvider: "visa"}}

	suite.MockCardClient.EXPECT().FindByCardNumber(gomock.Any(), &pb.FindByCardNumberRequest{CardNumber: cardNumber}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCard(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/card_number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(cardNumber, resp.Data.CardNumber)
}

func (suite *CardHandlerTestSuite) TestFindByCardNumber_Failure() {
	cardNumber := "1234567890123456"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().FindByCardNumber(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch card record", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/card_number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestCreateCard_Success() {
	expireDate := time.Now().AddDate(2, 0, 0)
	requestBody := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "visa",
	}

	grpcRequest := &pb.CreateCardRequest{
		UserId:       int32(requestBody.UserID),
		CardType:     requestBody.CardType,
		ExpireDate:   timestamppb.New(expireDate),
		Cvv:          requestBody.CVV,
		CardProvider: requestBody.CardProvider,
	}

	grpcResponse := &pb.ApiResponseCard{
		Status: "success",
		Data: &pb.CardResponse{
			Id:           3,
			UserId:       int32(requestBody.UserID),
			CardNumber:   "1234567890123456",
			CardType:     requestBody.CardType,
			ExpireDate:   expireDate.GoString(),
			Cvv:          requestBody.CVV,
			CardProvider: requestBody.CardProvider,
		},
	}

	expectedApiResponse := &response.ApiResponseCard{
		Status: "success",
		Data: &response.CardResponse{
			ID:           3,
			UserID:       requestBody.UserID,
			CardNumber:   "1234567890123456",
			CardType:     requestBody.CardType,
			ExpireDate:   expireDate.Format("2006-01-02"),
			CVV:          requestBody.CVV,
			CardProvider: requestBody.CardProvider,
		},
	}

	suite.MockCardClient.EXPECT().CreateCard(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCard(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.CreateCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(3, resp.Data.ID)
}

func (suite *CardHandlerTestSuite) TestCreateCard_Failure() {
	expireDate := time.Now().AddDate(2, 0, 0)
	requestBody := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "visa",
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().CreateCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create card", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.CreateCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestCreateCard_ValidationError() {
	invalidRequestBody := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "123",
		CardProvider: "visa",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug("Validation Error: ", gomock.Any()).
		Times(1)

	err := suite.Handler.CreateCard(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CardHandlerTestSuite) TestUpdateCard_Success() {
	id := 1
	expireDate := time.Now().AddDate(2, 0, 0)
	requestBody := requests.UpdateCardRequest{
		CardID:       id,
		UserID:       1,
		CardType:     "debit",
		ExpireDate:   expireDate,
		CVV:          "456",
		CardProvider: "mastercard",
	}

	grpcRequest := &pb.UpdateCardRequest{
		CardId:       int32(requestBody.CardID),
		UserId:       int32(requestBody.UserID),
		CardType:     requestBody.CardType,
		ExpireDate:   timestamppb.New(expireDate),
		Cvv:          requestBody.CVV,
		CardProvider: requestBody.CardProvider,
	}

	grpcResponse := &pb.ApiResponseCard{
		Status: "success",
		Data: &pb.CardResponse{
			Id:           int32(id),
			UserId:       int32(requestBody.UserID),
			CardNumber:   "1234567890123456",
			CardType:     requestBody.CardType,
			ExpireDate:   expireDate.GoString(),
			Cvv:          requestBody.CVV,
			CardProvider: requestBody.CardProvider,
		},
	}
	expectedApiResponse := &response.ApiResponseCard{
		Status: "success",
		Data: &response.CardResponse{
			ID:           id,
			UserID:       requestBody.UserID,
			CardNumber:   "1234567890123456",
			CardType:     requestBody.CardType,
			ExpireDate:   expireDate.Format("2006-01-02"),
			CVV:          requestBody.CVV,
			CardProvider: requestBody.CardProvider,
		},
	}

	suite.MockCardClient.EXPECT().UpdateCard(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCard(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.UpdateCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("debit", resp.Data.CardType)
}

func (suite *CardHandlerTestSuite) TestUpdateCard_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/card/update/abc", nil)
	suite.MockLogger.EXPECT().
		Debug("Bad Request: Invalid ID", gomock.Any()).
		Times(1)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.UpdateCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CardHandlerTestSuite) TestUpdateCard_ValidationError() {
	id := 1
	invalidRequestBody := requests.UpdateCardRequest{
		CardID:       id,
		UserID:       1,
		CardType:     "",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "456",
		CardProvider: "mastercard",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	suite.MockLogger.EXPECT().
		Debug("Validation Error: ", gomock.Any()).
		Times(1)

	err := suite.Handler.UpdateCard(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CardHandlerTestSuite) TestTrashedCard_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseCardDeleteAt{Status: "success", Data: &pb.CardResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseCardDeleteAt{Status: "success", Data: &response.CardResponseDeleteAt{ID: id}}

	suite.MockCardClient.EXPECT().TrashedCard(gomock.Any(), &pb.FindByIdCardRequest{CardId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCardDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/card/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CardHandlerTestSuite) TestTrashedCard_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().TrashedCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trashed card", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/card/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestRestoreCard_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseCardDeleteAt{Status: "success", Data: &pb.CardResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseCardDeleteAt{Status: "success", Data: &response.CardResponseDeleteAt{ID: id}}

	suite.MockCardClient.EXPECT().RestoreCard(gomock.Any(), &pb.FindByIdCardRequest{CardId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCardDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/card/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CardHandlerTestSuite) TestRestoreCard_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().RestoreCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore card", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/card/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestDeleteCardPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseCardDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseCardDelete{Status: "success"}

	suite.MockCardClient.EXPECT().DeleteCardPermanent(gomock.Any(), &pb.FindByIdCardRequest{CardId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCardDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/card/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteCardPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CardHandlerTestSuite) TestDeleteCardPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().DeleteCardPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete card", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/card/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteCardPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestRestoreAllCard_Success() {
	grpcResponse := &pb.ApiResponseCardAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseCardAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully restored all cards").Times(1)

	suite.MockCardClient.EXPECT().RestoreAllCard(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCardAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/card/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CardHandlerTestSuite) TestRestoreAllCard_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().RestoreAllCard(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all cards", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/card/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllCard(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CardHandlerTestSuite) TestDeleteAllCardPermanent_Success() {
	grpcResponse := &pb.ApiResponseCardAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseCardAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully deleted all cards permanently").Times(1)

	suite.MockCardClient.EXPECT().DeleteAllCardPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseCardAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/card/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllCardPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CardHandlerTestSuite) TestDeleteAllCardPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockCardClient.EXPECT().DeleteAllCardPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all cards", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/card/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllCardPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestCardHandlerSuite(t *testing.T) {
	suite.Run(t, new(CardHandlerTestSuite))
}
