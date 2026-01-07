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
)

type SaldoHandlerTestSuite struct {
	suite.Suite
	Ctrl               *gomock.Controller
	MockMerchantClient *mock_pb.MockMerchantServiceClient
	MockSaldoClient    *mock_pb.MockSaldoServiceClient
	MockLogger         *mock_logger.MockLoggerInterface
	MockMapper         *mock_apimapper.MockSaldoResponseMapper
	E                  *echo.Echo
	Handler            *api.SaldoHandleApi
}

func (suite *SaldoHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockSaldoClient = mock_pb.NewMockSaldoServiceClient(suite.Ctrl)
	suite.MockMerchantClient = mock_pb.NewMockMerchantServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockSaldoResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerSaldo(suite.MockSaldoClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *SaldoHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *SaldoHandlerTestSuite) TestFindAllSaldo_Success() {
	grpcResponse := &pb.ApiResponsePaginationSaldo{
		Status:  "success",
		Message: "Saldos retrieved successfully",
		Data: []*pb.SaldoResponse{
			{SaldoId: 1, CardNumber: "1234", TotalBalance: 100000},
			{SaldoId: 2, CardNumber: "5678", TotalBalance: 200000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationSaldo{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.SaldoResponse{{ID: 1, CardNumber: "1234", TotalBalance: 100000}, {ID: 2, CardNumber: "5678", TotalBalance: 200000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockSaldoClient.EXPECT().FindAllSaldo(gomock.Any(), &pb.FindAllSaldoRequest{Page: 1, PageSize: 10, Search: "test"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationSaldo(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos?page=1&page_size=10&search=test", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *SaldoHandlerTestSuite) TestFindAllSaldo_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().FindAllSaldo(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve saldo data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestFindByIdSaldo_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseSaldo{Status: "success", Data: &pb.SaldoResponse{SaldoId: int32(id), CardNumber: "1234", TotalBalance: 100000}}
	expectedApiResponse := &response.ApiResponseSaldo{Status: "success", Data: &response.SaldoResponse{ID: id, CardNumber: "1234", TotalBalance: 100000}}

	suite.MockSaldoClient.EXPECT().FindByIdSaldo(gomock.Any(), &pb.FindByIdSaldoRequest{SaldoId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldo(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *SaldoHandlerTestSuite) TestFindByIdSaldo_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/saldos/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	suite.MockLogger.EXPECT().
		Debug("Invalid saldo ID", gomock.Any()).
		Times(1)

	err := suite.Handler.FindById(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestFindByCardNumber_Success() {
	cardNumber := "1234567890123456"
	grpcResponse := &pb.ApiResponseSaldo{
		Status: "success",
		Data:   &pb.SaldoResponse{SaldoId: 1, CardNumber: cardNumber, TotalBalance: 100000},
	}
	expectedApiResponse := &response.ApiResponseSaldo{
		Status: "success",
		Data:   &response.SaldoResponse{ID: 1, CardNumber: cardNumber, TotalBalance: 100000},
	}

	suite.MockSaldoClient.EXPECT().FindByCardNumber(gomock.Any(), &pb.FindByCardNumberRequest{CardNumber: cardNumber}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldo(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldos/card_number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
}

func (suite *SaldoHandlerTestSuite) TestFindByCardNumber_Failure() {
	cardNumber := "1234567890123456"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().FindByCardNumber(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve saldo data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/saldos/card_number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestFindByActiveSaldo_Success() {
	grpcResponse := &pb.ApiResponsePaginationSaldoDeleteAt{Status: "success", Data: []*pb.SaldoResponseDeleteAt{{SaldoId: 1, CardNumber: "1234"}}}
	expectedApiResponse := &response.ApiResponsePaginationSaldoDeleteAt{Status: "success", Data: []*response.SaldoResponseDeleteAt{{ID: 1, CardNumber: "1234"}}}

	suite.MockSaldoClient.EXPECT().FindByActive(gomock.Any(), &pb.FindAllSaldoRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationSaldoDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationSaldoDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *SaldoHandlerTestSuite) TestFindByActiveSaldo_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().FindByActive(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve saldo data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestFindByTrashedSaldo_Success() {
	grpcResponse := &pb.ApiResponsePaginationSaldoDeleteAt{Status: "success", Data: []*pb.SaldoResponseDeleteAt{{SaldoId: 2, CardNumber: "5678"}}}
	expectedApiResponse := &response.ApiResponsePaginationSaldoDeleteAt{Status: "success", Data: []*response.SaldoResponseDeleteAt{{ID: 2, CardNumber: "5678"}}}

	suite.MockSaldoClient.EXPECT().FindByTrashed(gomock.Any(), &pb.FindAllSaldoRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationSaldoDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationSaldoDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *SaldoHandlerTestSuite) TestFindByTrashedSaldo_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().FindByTrashed(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve saldo data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/saldos/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestCreateSaldo_Success() {
	requestBody := requests.CreateSaldoRequest{
		CardNumber:   "1234567890123456",
		TotalBalance: 150000,
	}

	grpcRequest := &pb.CreateSaldoRequest{
		CardNumber:   requestBody.CardNumber,
		TotalBalance: int32(requestBody.TotalBalance),
	}

	grpcResponse := &pb.ApiResponseSaldo{
		Status: "success",
		Data: &pb.SaldoResponse{
			SaldoId:      3,
			CardNumber:   requestBody.CardNumber,
			TotalBalance: int32(requestBody.TotalBalance),
		},
	}

	expectedApiResponse := &response.ApiResponseSaldo{
		Status: "success",
		Data: &response.SaldoResponse{
			ID:           3,
			CardNumber:   requestBody.CardNumber,
			TotalBalance: requestBody.TotalBalance,
		},
	}

	suite.MockSaldoClient.EXPECT().CreateSaldo(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldo(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/saldos/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(3, resp.Data.ID)
}

func (suite *SaldoHandlerTestSuite) TestCreateSaldo_Failure() {
	requestBody := requests.CreateSaldoRequest{
		CardNumber:   "1234567890123456",
		TotalBalance: 150000,
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().CreateSaldo(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create saldo", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/saldos/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestCreateSaldo_ValidationError() {
	invalidRequestBody := requests.CreateSaldoRequest{
		CardNumber:   "", // Empty card number should fail validation
		TotalBalance: 150000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/saldos/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	suite.MockLogger.EXPECT().
		Debug("Validation Error", gomock.Any()).
		Times(1)

	err := suite.Handler.Create(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestUpdateSaldo_Success() {
	id := 1
	requestBody := requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   "1234567890123456",
		TotalBalance: 200000,
	}

	grpcRequest := &pb.UpdateSaldoRequest{
		SaldoId:      int32(id),
		CardNumber:   requestBody.CardNumber,
		TotalBalance: int32(requestBody.TotalBalance),
	}

	grpcResponse := &pb.ApiResponseSaldo{
		Status: "success",
		Data: &pb.SaldoResponse{
			SaldoId:      int32(id),
			TotalBalance: int32(requestBody.TotalBalance),
		},
	}
	expectedApiResponse := &response.ApiResponseSaldo{
		Status: "success",
		Data: &response.SaldoResponse{
			ID:           id,
			TotalBalance: requestBody.TotalBalance,
		},
	}

	suite.MockSaldoClient.EXPECT().UpdateSaldo(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldo(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/saldos/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseSaldo
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(200000), resp.Data.TotalBalance)
}

func (suite *SaldoHandlerTestSuite) TestUpdateSaldo_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/saldos/update/abc", nil)
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

func (suite *SaldoHandlerTestSuite) TestUpdateSaldo_ValidationError() {
	id := 1
	invalidRequestBody := requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   "", // Empty card number should fail validation
		TotalBalance: 200000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/saldos/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	suite.MockLogger.EXPECT().
		Debug("Validation Error", gomock.Any()).
		Times(1)

	err := suite.Handler.Update(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestTrashSaldo_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseSaldoDeleteAt{Status: "success", Data: &pb.SaldoResponseDeleteAt{SaldoId: int32(id)}}
	expectedApiResponse := &response.ApiResponseSaldoDeleteAt{Status: "success", Data: &response.SaldoResponseDeleteAt{ID: id}}

	suite.MockSaldoClient.EXPECT().TrashedSaldo(gomock.Any(), &pb.FindByIdSaldoRequest{SaldoId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldoDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashSaldo(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestTrashSaldo_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().TrashedSaldo(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trashed saldo", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashSaldo(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestRestoreSaldo_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseSaldoDeleteAt{Status: "success", Data: &pb.SaldoResponseDeleteAt{SaldoId: int32(id)}}
	expectedApiResponse := &response.ApiResponseSaldoDeleteAt{Status: "success", Data: &response.SaldoResponseDeleteAt{ID: id}}

	suite.MockSaldoClient.EXPECT().RestoreSaldo(gomock.Any(), &pb.FindByIdSaldoRequest{SaldoId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldoDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreSaldo(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestRestoreSaldo_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().RestoreSaldo(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore saldo", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreSaldo(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestDeleteSaldoPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseSaldoDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseSaldoDelete{Status: "success"}

	suite.MockSaldoClient.EXPECT().DeleteSaldoPermanent(gomock.Any(), &pb.FindByIdSaldoRequest{SaldoId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldoDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/saldos/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Delete(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestDeleteSaldoPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().DeleteSaldoPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete saldo", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/saldos/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Delete(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestRestoreAllSaldo_Success() {
	grpcResponse := &pb.ApiResponseSaldoAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseSaldoAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully restored all saldo").Times(1)

	suite.MockSaldoClient.EXPECT().RestoreAllSaldo(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldoAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllSaldo(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestRestoreAllSaldo_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().RestoreAllSaldo(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all saldo", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllSaldo(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestDeleteAllSaldoPermanent_Success() {
	grpcResponse := &pb.ApiResponseSaldoAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseSaldoAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully deleted all saldo permanently").Times(1)

	suite.MockSaldoClient.EXPECT().DeleteAllSaldoPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseSaldoAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllSaldoPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *SaldoHandlerTestSuite) TestDeleteAllSaldoPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockSaldoClient.EXPECT().DeleteAllSaldoPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all saldo", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/saldos/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllSaldoPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestSaldoHandlerSuite(t *testing.T) {
	suite.Run(t, new(SaldoHandlerTestSuite))
}
