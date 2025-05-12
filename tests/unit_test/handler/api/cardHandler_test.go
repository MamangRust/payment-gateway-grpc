package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	mock_apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindAllCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockCards := &pb.ApiResponsePaginationCard{
		Status:  "success",
		Message: "Fetched cards successfully",
		Data: []*pb.CardResponse{
			{Id: 1, CardNumber: "Card 1"},
			{Id: 2, CardNumber: "Card 2"},
		},
	}
	mockCardClient.EXPECT().FindAllCard(gomock.Any(), gomock.Any()).Return(mockCards, nil).Times(1)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/card/?page=1&page_size=10", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid_token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response pb.ApiResponsePaginationCard

	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Data, 2)
}

func TestFindAllCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/card/?page=1&page_size=10", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer invalid_token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCardClient.EXPECT().FindAllCard(gomock.Any(), gomock.Any()).Return(nil, echo.ErrUnauthorized).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindAll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Unauthorized", resp.Message)
}

func TestFindAllCard_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockCards := &pb.ApiResponsePaginationCard{
		Status:  "success",
		Message: "No cards found",
		Data:    []*pb.CardResponse{},
	}
	mockCardClient.EXPECT().FindAllCard(gomock.Any(), gomock.Any()).Return(mockCards, nil).Times(1)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/api/card/?page=1&page_size=10", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid_token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response pb.ApiResponsePaginationCard
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Data, 0)
}

func TestFindByIdCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	expectedCard := &pb.CardResponse{
		Id:         int32(cardID),
		CardNumber: "1234567890123456",
	}

	mockResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    expectedCard,
	}

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockCardClient.EXPECT().FindByIdCard(gomock.Any(), &pb.FindByIdCardRequest{CardId: int32(cardID)}).Return(mockResponse, nil).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    pb.CardResponse `json:"data"`
	}

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully fetched card record", resp.Message)
	assert.Equal(t, expectedCard.Id, resp.Data.Id)
	assert.Equal(t, expectedCard.CardNumber, resp.Data.CardNumber)
}

func TestFindByIdCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	mockResponse := &pb.ApiResponseCard{
		Status:  "error",
		Message: "Failed to fetch card record",
	}

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockCardClient.EXPECT().FindByIdCard(gomock.Any(), &pb.FindByIdCardRequest{CardId: int32(cardID)}).Return(mockResponse, nil).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to fetch card record", resp.Message)
}

func TestFindByIdCard_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})
	req := httptest.NewRequest(http.MethodGet, "/api/card/invalid_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid_id")

	mockLogger.EXPECT().Debug("Invalid card ID", gomock.Any()).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid card ID", resp.Message)
}

func TestFindByUserID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	userID := int32(42)
	token := "mocked_token_string"

	expectedCard := &pb.CardResponse{
		Id:         1,
		CardNumber: "1234567890123456",
	}
	mockResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    expectedCard,
	}

	mockCardClient.EXPECT().
		FindByUserIdCard(
			gomock.Any(),
			&pb.FindByUserIdCardRequest{UserId: userID},
		).
		Return(mockResponse, nil)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := jwt.MapClaims{"user_id": 1}
			token := &jwt.Token{Claims: claims}
			c.Set("user", token)
			c.Set("userID", 1)
			return next(c)
		}
	})
	req := httptest.NewRequest(http.MethodGet, "/api/card/user", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", userID)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)
	err := handler.FindByUserID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    pb.CardResponse `json:"data"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully fetched card record", resp.Message)
	assert.Equal(t, expectedCard.Id, resp.Data.Id)
	assert.Equal(t, expectedCard.CardNumber, resp.Data.CardNumber)
}

func TestFindByUserID_Failure_InvalidUserIDType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	invalidUserID := "not_a_number"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", invalidUserID)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByUserID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to parse UserID", resp.Message)
}

func TestFindByUserID_Failure_CardFetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	userID := int32(42)
	token := "mocked_token_string"

	mockLogger.EXPECT().Debug("Failed to fetch card record", gomock.Any()).Times(1)

	mockCardClient.EXPECT().
		FindByUserIdCard(
			gomock.Any(),
			&pb.FindByUserIdCardRequest{UserId: userID},
		).
		Return(nil, fmt.Errorf("internal server error"))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/user", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", userID)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByUserID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to fetch card record: internal server error", resp.Message)
}

func TestFindByActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	expectedCard := []*pb.CardResponseDeleteAt{
		{
			Id:         1,
			CardNumber: "1234567890123456",
		},
		{
			Id:         2,
			CardNumber: "9876543210987654",
		},
	}

	request := &pb.FindAllCardRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockResponse := &pb.ApiResponsePaginationCardDeleteAt{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    expectedCard,
	}

	mockCardClient.EXPECT().
		FindByActiveCard(context.Background(), request).
		Return(mockResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCards
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully fetched card record", resp.Message)
	assert.Equal(t, expectedCard[0].Id, resp.Data[0].Id)
	assert.Equal(t, expectedCard[0].CardNumber, resp.Data[0].CardNumber)
}

func TestFindByActive_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockResponse := &pb.ApiResponsePaginationCardDeleteAt{
		Status:  "success",
		Message: "No active cards found",
		Data:    []*pb.CardResponseDeleteAt{},
		Pagination: &pb.PaginationMeta{
			CurrentPage:  1,
			PageSize:     10,
			TotalPages:   1,
			TotalRecords: 0,
		},
	}

	request := &pb.FindAllCardRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockCardClient.EXPECT().
		FindByActiveCard(context.Background(), request).
		Return(mockResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCards
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No active cards found", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestFindByActive_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockLogger.EXPECT().Debug("Failed to fetch card record", gomock.Any()).Times(1)

	request := &pb.FindAllCardRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockCardClient.EXPECT().
		FindByActiveCard(gomock.Any(), request).
		Return(nil, fmt.Errorf("internal server error"))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByActive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to fetch card record: internal server error", resp.Message)
}

func TestFindByTrashed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	expectedCard := []*pb.CardResponseDeleteAt{
		{
			Id:         1,
			CardNumber: "1234567890123456",
		},
		{
			Id:         2,
			CardNumber: "9876543210987654",
		},
	}
	request := &pb.FindAllCardRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockResponse := &pb.ApiResponsePaginationCardDeleteAt{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    expectedCard,
	}

	mockCardClient.EXPECT().
		FindByTrashedCard(context.Background(), request).
		Return(mockResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCards
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully fetched card record", resp.Message)
	assert.Equal(t, expectedCard[0].Id, resp.Data[0].Id)
	assert.Equal(t, expectedCard[0].CardNumber, resp.Data[0].CardNumber)
}

func TestFindByTrashed_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	request := &pb.FindAllCardRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockResponse := &pb.ApiResponsePaginationCardDeleteAt{
		Status:  "success",
		Message: "No trashed cards found",
		Data:    []*pb.CardResponseDeleteAt{},
	}

	mockCardClient.EXPECT().
		FindByTrashedCard(context.Background(), request).
		Return(mockResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCards
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No trashed cards found", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestFindByTrashed_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockLogger.EXPECT().Debug("Failed to fetch card record", gomock.Any()).Times(1)

	mockCardClient.EXPECT().
		FindByTrashedCard(gomock.Any(), &emptypb.Empty{}).
		Return(nil, fmt.Errorf("internal server error"))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/card/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByTrashed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to fetch card record: internal server error", resp.Message)

}

func TestFindByCardNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardNumber := "1234567890123456"

	expectedCard := &pb.CardResponse{
		Id:         1,
		CardNumber: "1234567890123456",
	}

	mockResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    expectedCard,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	mockCardClient.EXPECT().FindByCardNumber(gomock.Any(), &pb.FindByCardNumberRequest{CardNumber: cardNumber}).Return(mockResponse, nil).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    pb.CardResponse `json:"data"`
	}

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully fetched card record", resp.Message)
	assert.Equal(t, expectedCard.Id, resp.Data.Id)
	assert.Equal(t, expectedCard.CardNumber, resp.Data.CardNumber)
}

func TestFindByCardNumberCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardNumber := "1234567890123456"

	mockError := status.Errorf(codes.NotFound, "Card with number %s not found", cardNumber)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/card/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	mockCardClient.EXPECT().FindByCardNumber(gomock.Any(), &pb.FindByCardNumberRequest{CardNumber: cardNumber}).Return(nil, mockError).Times(1)

	mockLogger.EXPECT().Debug("Failed to fetch card record", gomock.Any()).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.FindByCardNumber(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to fetch card record: ", resp.Message)
}

func TestCreateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	body := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	}

	expectedResponse := &pb.CardResponse{
		Id:           1,
		UserId:       1,
		CardType:     "credit",
		ExpireDate:   body.ExpireDate.String(),
		Cvv:          "123",
		CardProvider: "Visa",
	}

	mockResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully created card record",
		Data:    expectedResponse,
	}

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCardClient.EXPECT().
		CreateCard(gomock.Any(), &pb.CreateCardRequest{
			UserId:       int32(body.UserID),
			CardType:     body.CardType,
			ExpireDate:   timestamppb.New(body.ExpireDate),
			Cvv:          body.CVV,
			CardProvider: body.CardProvider,
		}).
		Return(mockResponse, nil).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.CreateCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, mockResponse.Data.Id, resp.Data.Id)
	assert.Equal(t, mockResponse.Data.CardType, resp.Data.CardType)
	assert.Equal(t, mockResponse.Data.CardProvider, resp.Data.CardProvider)
}

func TestCreateCard_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	body := requests.CreateCardRequest{
		UserID:       0,
		CardType:     "",
		ExpireDate:   time.Now().AddDate(0, 0, -1),
		CVV:          "",
		CardProvider: "",
	}

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger.EXPECT().
		Debug("Validation Error: ", gomock.Any()).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.CreateCard(c)

	fmt.Println(err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error: ")
}

func TestCreateCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	body := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   time.Now().AddDate(1, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	}

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/card/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCardClient.EXPECT().
		CreateCard(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to create card", gomock.Any()).Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.CreateCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to create card")
}

func TestUpdateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	body := requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	cardID := 1

	expectedResponse := &pb.CardResponse{
		Id:           int32(cardID),
		UserId:       int32(cardID),
		CardNumber:   "1234567890123456",
		CardType:     "debit",
		ExpireDate:   body.ExpireDate.String(),
		Cvv:          "456",
		CardProvider: "MasterCard",
	}

	mockResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully updated card record",
		Data:    expectedResponse,
	}

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/card/update/%d", cardID), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockCardClient.EXPECT().
		UpdateCard(gomock.Any(), &pb.UpdateCardRequest{
			CardId:       int32(cardID),
			UserId:       int32(body.UserID),
			CardType:     body.CardType,
			ExpireDate:   timestamppb.New(body.ExpireDate),
			Cvv:          body.CVV,
			CardProvider: body.CardProvider,
		}).
		Return(mockResponse, nil).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.UpdateCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, mockResponse.Data.Id, resp.Data.Id)
	assert.Equal(t, mockResponse.Data.CardType, resp.Data.CardType)
	assert.Equal(t, mockResponse.Data.CardProvider, resp.Data.CardProvider)
}

func TestUpdateCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	body := requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(2, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	cardID := 1

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/card/update/%d", cardID), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockLogger.EXPECT().Debug("Failed to update card", gomock.Any()).Times(1)

	mockCardClient.EXPECT().
		UpdateCard(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("internal server error")).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.UpdateCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to update card: ")
}

func TestUpdateCard_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	body := requests.UpdateCardRequest{
		CardType: "",
	}

	cardID := 1

	e := echo.New()
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/card/update/%d", cardID), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockLogger.EXPECT().Debug(gomock.Any()).Times(1)

	handler := api.NewHandlerCard(nil, e, mockLogger, mockMapper)

	err := handler.UpdateCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error")
}

func TestTrashedCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	expectedResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully trashed card record",
		Data: &pb.CardResponse{
			Id:           int32(cardID),
			CardType:     "debit",
			ExpireDate:   "2025-12-31",
			Cvv:          "123",
			CardProvider: "Visa",
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/card/trash/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockCardClient.EXPECT().
		TrashedCard(gomock.Any(), &pb.FindByIdCardRequest{
			CardId: int32(cardID),
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.TrashedCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully trashed card record", resp.Message)
	assert.Equal(t, expectedResponse.Data.Id, resp.Data.Id)
	assert.Equal(t, expectedResponse.Data.CardType, resp.Data.CardType)
}

func TestTrashedCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	mockLogger.EXPECT().Debug("Failed to trashed card", gomock.Any()).Times(1)

	mockCardClient.EXPECT().
		TrashedCard(gomock.Any(), &pb.FindByIdCardRequest{
			CardId: int32(cardID),
		}).
		Return(nil, fmt.Errorf("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/card/trash/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.TrashedCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to trashed card: internal server error", resp.Message)
}

func TestTrashedCard_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/card/trash/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	handler := api.NewHandlerCard(nil, e, mockLogger, mockMapper)

	err := handler.TrashedCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestRestoreCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	expectedResponse := &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully restored card record",
		Data: &pb.CardResponse{
			Id:           int32(cardID),
			CardType:     "debit",
			ExpireDate:   "2025-12-31",
			Cvv:          "123",
			CardProvider: "Visa",
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card/restore/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockCardClient.EXPECT().
		RestoreCard(gomock.Any(), &pb.FindByIdCardRequest{
			CardId: int32(cardID),
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.RestoreCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully restored card record", resp.Message)
	assert.Equal(t, expectedResponse.Data.Id, resp.Data.Id)
	assert.Equal(t, expectedResponse.Data.CardType, resp.Data.CardType)
}

func TestRestoreCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	mockLogger.EXPECT().Debug("Failed to restore card", gomock.Any()).Times(1)

	mockCardClient.EXPECT().
		RestoreCard(gomock.Any(), &pb.FindByIdCardRequest{
			CardId: int32(cardID),
		}).
		Return(nil, fmt.Errorf("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/card/restore/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.RestoreCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to restore card: internal server error", resp.Message)
}

func TestRestoreCard_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/card/restore/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	handler := api.NewHandlerCard(nil, e, mockLogger, nil)

	err := handler.RestoreCard(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}

func TestDeleteCardPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	expectedResponse := &pb.ApiResponseCardDelete{
		Status:  "success",
		Message: "Successfully deleted card record",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/card/delete/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	mockCardClient.EXPECT().
		DeleteCardPermanent(gomock.Any(), &pb.FindByIdCardRequest{
			CardId: int32(cardID),
		}).
		Return(expectedResponse, nil).
		Times(1)

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.DeleteCardPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseCard
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Successfully deleted card record", resp.Message)
	assert.Nil(t, resp.Data)
}

func TestDeleteCardPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	cardID := 1

	mockLogger.EXPECT().Debug("Failed to delete card", gomock.Any()).Times(1)

	mockCardClient.EXPECT().
		DeleteCardPermanent(gomock.Any(), &pb.FindByIdCardRequest{
			CardId: int32(cardID),
		}).
		Return(nil, fmt.Errorf("internal server error")).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/card/delete/%d", cardID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", cardID))

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	err := handler.DeleteCardPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to delete card: internal server error", resp.Message)
}

func TestDeleteCardPermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	mockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/card/delete/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	handler := api.NewHandlerCard(nil, e, mockLogger, mockMapper)

	err := handler.DeleteCardPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Bad Request: Invalid ID", resp.Message)
}
