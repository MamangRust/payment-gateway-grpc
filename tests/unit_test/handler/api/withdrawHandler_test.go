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

type WithdrawHandlerTestSuite struct {
	suite.Suite
	Ctrl               *gomock.Controller
	MockWithdrawClient *mock_pb.MockWithdrawServiceClient
	MockLogger         *mock_logger.MockLoggerInterface
	MockMapper         *mock_apimapper.MockWithdrawResponseMapper
	E                  *echo.Echo
	Handler            *api.WithdrawHandleApi
}

func (suite *WithdrawHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockWithdrawClient = mock_pb.NewMockWithdrawServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockWithdrawResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerWithdraw(suite.MockWithdrawClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *WithdrawHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *WithdrawHandlerTestSuite) TestFindAll_Success() {
	grpcResponse := &pb.ApiResponsePaginationWithdraw{
		Status:  "success",
		Message: "Withdraws retrieved successfully",
		Data: []*pb.WithdrawResponse{
			{WithdrawId: 1, WithdrawNo: "WD123", CardNumber: "1234", WithdrawAmount: 100000},
			{WithdrawId: 2, WithdrawNo: "WD124", CardNumber: "5678", WithdrawAmount: 200000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationWithdraw{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.WithdrawResponse{{ID: 1, WithdrawNo: "WD123", CardNumber: "1234", WithdrawAmount: 100000}, {ID: 2, WithdrawNo: "WD124", CardNumber: "5678", WithdrawAmount: 200000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockWithdrawClient.EXPECT().FindAllWithdraw(gomock.Any(), &pb.FindAllWithdrawRequest{Page: 1, PageSize: 10, Search: "test"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationWithdraw(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw?page=1&page_size=10&search=test", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *WithdrawHandlerTestSuite) TestFindAll_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().FindAllWithdraw(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve withdraw data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestFindAllByCardNumber_Success() {
	cardNumber := "1234567890123456"
	grpcResponse := &pb.ApiResponsePaginationWithdraw{
		Status:  "success",
		Message: "Withdraws for card retrieved successfully",
		Data: []*pb.WithdrawResponse{
			{WithdrawId: 1, WithdrawNo: "WD123", CardNumber: cardNumber, WithdrawAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 1, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationWithdraw{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.WithdrawResponse{{ID: 1, WithdrawNo: "WD123", CardNumber: cardNumber, WithdrawAmount: 100000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 1, TotalPages: 1},
	}

	suite.MockWithdrawClient.EXPECT().FindAllWithdrawByCardNumber(gomock.Any(), &pb.FindAllWithdrawByCardNumberRequest{CardNumber: cardNumber, Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationWithdraw(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/withdraw/card-number/%s?page=1&page_size=10", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindAllByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *WithdrawHandlerTestSuite) TestFindAllByCardNumber_Failure() {
	cardNumber := "1234567890123456"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().FindAllWithdrawByCardNumber(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve withdraw data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/withdraw/card-number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindAllByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestFindById_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseWithdraw{Status: "success", Data: &pb.WithdrawResponse{WithdrawId: int32(id), WithdrawNo: "WD123", CardNumber: "1234", WithdrawAmount: 100000}}
	expectedApiResponse := &response.ApiResponseWithdraw{Status: "success", Data: &response.WithdrawResponse{ID: id, WithdrawNo: "WD123", CardNumber: "1234", WithdrawAmount: 100000}}

	suite.MockWithdrawClient.EXPECT().FindByIdWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{WithdrawId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdraw(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *WithdrawHandlerTestSuite) TestFindById_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	suite.MockLogger.EXPECT().
		Debug("Invalid withdraw ID", gomock.Any()).
		Times(1)

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestFindByActive_Success() {
	grpcResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{Status: "success", Data: []*pb.WithdrawResponseDeleteAt{{WithdrawId: 1, WithdrawNo: "WD123"}}}
	expectedApiResponse := &response.ApiResponsePaginationWithdrawDeleteAt{Status: "success", Data: []*response.WithdrawResponseDeleteAt{{ID: 1, WithdrawNo: "WD123"}}}

	suite.MockWithdrawClient.EXPECT().FindByActive(gomock.Any(), &pb.FindAllWithdrawRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationWithdrawDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationWithdrawDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *WithdrawHandlerTestSuite) TestFindByActive_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().FindByActive(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve withdraw data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestFindByTrashed_Success() {
	grpcResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{Status: "success", Data: []*pb.WithdrawResponseDeleteAt{{WithdrawId: 2, WithdrawNo: "WD999"}}}
	expectedApiResponse := &response.ApiResponsePaginationWithdrawDeleteAt{Status: "success", Data: []*response.WithdrawResponseDeleteAt{{ID: 2, WithdrawNo: "WD999"}}}

	suite.MockWithdrawClient.EXPECT().FindByTrashed(gomock.Any(), &pb.FindAllWithdrawRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationWithdrawDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationWithdrawDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *WithdrawHandlerTestSuite) TestFindByTrashed_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().FindByTrashed(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve withdraw data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/withdraw/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestCreate_Success() {
	withdrawTime := time.Now()
	requestBody := requests.CreateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 150000,
		WithdrawTime:   withdrawTime,
	}
	grpcRequest := &pb.CreateWithdrawRequest{
		CardNumber:     requestBody.CardNumber,
		WithdrawAmount: int32(requestBody.WithdrawAmount),
		WithdrawTime:   timestamppb.New(withdrawTime),
	}
	grpcResponse := &pb.ApiResponseWithdraw{Status: "success", Data: &pb.WithdrawResponse{WithdrawId: 3, WithdrawNo: "WD125", CardNumber: requestBody.CardNumber, WithdrawAmount: int32(requestBody.WithdrawAmount)}}
	expectedApiResponse := &response.ApiResponseWithdraw{Status: "success", Data: &response.WithdrawResponse{ID: 3, WithdrawNo: "WD125", CardNumber: requestBody.CardNumber, WithdrawAmount: requestBody.WithdrawAmount}}

	suite.MockWithdrawClient.EXPECT().CreateWithdraw(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdraw(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("WD125", resp.Data.WithdrawNo)
}

func (suite *WithdrawHandlerTestSuite) TestCreate_Failure() {
	requestBody := requests.CreateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 150000,
		WithdrawTime:   time.Now(),
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().CreateWithdraw(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create withdraw", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestCreate_ValidationError() {
	invalidRequestBody := requests.CreateWithdrawRequest{
		CardNumber:     "",
		WithdrawAmount: 40000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/withdraw/create",
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug(gomock.Any()).
		Times(1)

	err := suite.Handler.Create(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestUpdate_Success() {
	id := 1
	withdrawTime := time.Now()
	requestBody := requests.UpdateWithdrawRequest{
		CardNumber:     "1234567890123456",
		WithdrawAmount: 200000,
		WithdrawTime:   withdrawTime,
	}
	grpcRequest := &pb.UpdateWithdrawRequest{
		WithdrawId:     int32(id),
		CardNumber:     requestBody.CardNumber,
		WithdrawAmount: int32(requestBody.WithdrawAmount),
		WithdrawTime:   timestamppb.New(withdrawTime),
	}
	grpcResponse := &pb.ApiResponseWithdraw{Status: "success", Data: &pb.WithdrawResponse{WithdrawId: int32(id), WithdrawAmount: int32(requestBody.WithdrawAmount)}}
	expectedApiResponse := &response.ApiResponseWithdraw{Status: "success", Data: &response.WithdrawResponse{ID: id, WithdrawAmount: requestBody.WithdrawAmount}}

	suite.MockWithdrawClient.EXPECT().UpdateWithdraw(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdraw(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/withdraw/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseWithdraw
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(200000), resp.Data.WithdrawAmount)
}

func (suite *WithdrawHandlerTestSuite) TestUpdate_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/update/abc", nil)
	suite.MockLogger.EXPECT().
		Debug("Invalid withdraw ID", gomock.Any()).
		Times(1)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestUpdate_ValidationError() {
	id := 1
	invalidRequestBody := requests.UpdateWithdrawRequest{
		CardNumber:     "",
		WithdrawAmount: 40000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/api/withdraw/update/%d", id),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	suite.MockLogger.EXPECT().
		Debug(gomock.Any()).
		Times(1)

	err := suite.Handler.Update(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestTrashWithdraw_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseWithdrawDeleteAt{Status: "success", Data: &pb.WithdrawResponseDeleteAt{WithdrawId: int32(id)}}
	expectedApiResponse := &response.ApiResponseWithdrawDeleteAt{Status: "success", Data: &response.WithdrawResponseDeleteAt{ID: id}}

	suite.MockWithdrawClient.EXPECT().TrashedWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{WithdrawId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdrawDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashWithdraw(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestTrashWithdraw_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().TrashedWithdraw(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trash withdraw", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashWithdraw(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestRestoreWithdraw_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseWithdrawDeleteAt{Status: "success", Data: &pb.WithdrawResponseDeleteAt{WithdrawId: int32(id)}}
	expectedApiResponse := &response.ApiResponseWithdrawDeleteAt{Status: "success", Data: &response.WithdrawResponseDeleteAt{ID: id}}

	suite.MockWithdrawClient.EXPECT().RestoreWithdraw(gomock.Any(), &pb.FindByIdWithdrawRequest{WithdrawId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdrawDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreWithdraw(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestRestoreWithdraw_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().RestoreWithdraw(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore withdraw", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreWithdraw(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestDeleteWithdrawPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseWithdrawDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseWithdrawDelete{Status: "success"}

	suite.MockWithdrawClient.EXPECT().DeleteWithdrawPermanent(gomock.Any(), &pb.FindByIdWithdrawRequest{WithdrawId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdrawDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteWithdrawPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestDeleteWithdrawPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().DeleteWithdrawPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete withdraw permanently", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteWithdrawPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestRestoreAllWithdraw_Success() {
	grpcResponse := &pb.ApiResponseWithdrawAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseWithdrawAll{Status: "success"}
	suite.MockLogger.EXPECT().Debug("Successfully restored all withdraw").Times(1)

	suite.MockWithdrawClient.EXPECT().RestoreAllWithdraw(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdrawAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllWithdraw(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestRestoreAllWithdraw_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().RestoreAllWithdraw(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all withdraw", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/withdraw/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllWithdraw(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestDeleteAllWithdrawPermanent_Success() {
	grpcResponse := &pb.ApiResponseWithdrawAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseWithdrawAll{Status: "success"}
	suite.MockLogger.EXPECT().Debug("Successfully deleted all withdraw permanently").Times(1)

	suite.MockWithdrawClient.EXPECT().DeleteAllWithdrawPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseWithdrawAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllWithdrawPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *WithdrawHandlerTestSuite) TestDeleteAllWithdrawPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockWithdrawClient.EXPECT().DeleteAllWithdrawPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all withdraw", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/withdraw/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllWithdrawPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestWithdrawHandlerSuite(t *testing.T) {
	suite.Run(t, new(WithdrawHandlerTestSuite))
}
