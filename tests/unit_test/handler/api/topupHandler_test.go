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

type TopupHandlerTestSuite struct {
	suite.Suite
	Ctrl               *gomock.Controller
	MockMerchantClient *mock_pb.MockMerchantServiceClient
	MockTopupClient    *mock_pb.MockTopupServiceClient
	MockLogger         *mock_logger.MockLoggerInterface
	MockMapper         *mock_apimapper.MockTopupResponseMapper
	E                  *echo.Echo
	Handler            *api.TopupHandleApi
}

func (suite *TopupHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockTopupClient = mock_pb.NewMockTopupServiceClient(suite.Ctrl)
	suite.MockMerchantClient = mock_pb.NewMockMerchantServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockTopupResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerTopup(suite.MockTopupClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *TopupHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *TopupHandlerTestSuite) TestFindAllTopup_Success() {
	grpcResponse := &pb.ApiResponsePaginationTopup{
		Status:  "success",
		Message: "Topups retrieved successfully",
		Data: []*pb.TopupResponse{
			{Id: 1, TopupNo: "TP123", CardNumber: "1234", TopupAmount: 100000},
			{Id: 2, TopupNo: "TP124", CardNumber: "5678", TopupAmount: 200000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationTopup{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.TopupResponse{{ID: 1, TopupNo: "TP123", CardNumber: "1234", TopupAmount: 100000}, {ID: 2, TopupNo: "TP124", CardNumber: "5678", TopupAmount: 200000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockTopupClient.EXPECT().FindAllTopup(gomock.Any(), &pb.FindAllTopupRequest{Page: 1, PageSize: 10, Search: "test"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTopup(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/topups?page=1&page_size=10&search=test", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *TopupHandlerTestSuite) TestFindAllTopup_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().FindAllTopup(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve topup data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/topups", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestFindAllTopupByCardNumber_Success() {
	cardNumber := "1234567890123456"
	grpcResponse := &pb.ApiResponsePaginationTopup{
		Status:  "success",
		Message: "Topups for card retrieved successfully",
		Data: []*pb.TopupResponse{
			{Id: 1, TopupNo: "TP123", CardNumber: cardNumber, TopupAmount: 100000},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 1, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationTopup{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.TopupResponse{{ID: 1, TopupNo: "TP123", CardNumber: cardNumber, TopupAmount: 100000}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 1, TotalPages: 1},
	}

	suite.MockTopupClient.EXPECT().FindAllTopupByCardNumber(gomock.Any(), &pb.FindAllTopupByCardNumberRequest{CardNumber: cardNumber, Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTopup(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/topups/card-number/%s?page=1&page_size=10", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindAllByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TopupHandlerTestSuite) TestFindAllTopupByCardNumber_Failure() {
	cardNumber := "1234567890123456"
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().FindAllTopupByCardNumber(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve topup data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/topups/card-number/%s", cardNumber), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("card_number")
	c.SetParamValues(cardNumber)

	err := suite.Handler.FindAllByCardNumber(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestFindByIdTopup_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTopup{Status: "success", Data: &pb.TopupResponse{Id: int32(id), TopupNo: "TP123", CardNumber: "1234", TopupAmount: 100000}}
	expectedApiResponse := &response.ApiResponseTopup{Status: "success", Data: &response.TopupResponse{ID: id, TopupNo: "TP123", CardNumber: "1234", TopupAmount: 100000}}

	suite.MockTopupClient.EXPECT().FindByIdTopup(gomock.Any(), &pb.FindByIdTopupRequest{TopupId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopup(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/topups/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *TopupHandlerTestSuite) TestFindByIdTopup_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/topups/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.FindById(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestFindByActiveTopup_Success() {
	grpcResponse := &pb.ApiResponsePaginationTopupDeleteAt{Status: "success", Data: []*pb.TopupResponseDeleteAt{{Id: 1, TopupNo: "TP123"}}}
	expectedApiResponse := &response.ApiResponsePaginationTopupDeleteAt{Status: "success", Data: []*response.TopupResponseDeleteAt{{ID: 1, TopupNo: "TP123"}}}

	suite.MockTopupClient.EXPECT().FindByActive(gomock.Any(), &pb.FindAllTopupRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTopupDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/topups/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTopupDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TopupHandlerTestSuite) TestFindByActiveTopup_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().FindByActive(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve topup data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/topups/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestFindByTrashedTopup_Success() {
	grpcResponse := &pb.ApiResponsePaginationTopupDeleteAt{Status: "success", Data: []*pb.TopupResponseDeleteAt{{Id: 2, TopupNo: "TP999"}}}
	expectedApiResponse := &response.ApiResponsePaginationTopupDeleteAt{Status: "success", Data: []*response.TopupResponseDeleteAt{{ID: 2, TopupNo: "TP999"}}}

	suite.MockTopupClient.EXPECT().FindByTrashed(gomock.Any(), &pb.FindAllTopupRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationTopupDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/topups/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTopupDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *TopupHandlerTestSuite) TestFindByTrashedTopup_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().FindByTrashed(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve topup data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/topups/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestCreateTopup_Success() {
	apiKey := "test-api-key"

	requestBody := requests.CreateTopupRequest{
		CardNumber:  "1234567890123456",
		TopupAmount: 150000,
		TopupMethod: "alfamart",
	}

	grpcRequest := &pb.CreateTopupRequest{
		CardNumber:  requestBody.CardNumber,
		TopupAmount: int32(requestBody.TopupAmount),
		TopupMethod: requestBody.TopupMethod,
	}

	grpcResponse := &pb.ApiResponseTopup{
		Status: "success",
		Data: &pb.TopupResponse{
			Id:          3,
			TopupNo:     "TP125",
			CardNumber:  requestBody.CardNumber,
			TopupAmount: int32(requestBody.TopupAmount),
		},
	}

	expectedApiResponse := &response.ApiResponseTopup{
		Status: "success",
		Data: &response.TopupResponse{
			ID:          3,
			TopupNo:     "TP125",
			CardNumber:  requestBody.CardNumber,
			TopupAmount: requestBody.TopupAmount,
		},
	}

	suite.MockTopupClient.EXPECT().CreateTopup(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopup(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/topups/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("apiKey", apiKey)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("TP125", resp.Data.TopupNo)
}

func (suite *TopupHandlerTestSuite) TestCreateTopup_Failure() {
	apiKey := "test-api-key"

	requestBody := requests.CreateTopupRequest{
		CardNumber:  "1234567890123456",
		TopupAmount: 150000,
		TopupMethod: "alfamart",
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().CreateTopup(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create topup", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/topups/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("apiKey", apiKey)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestCreateTopup_ValidationError() {

	invalidRequestBody := requests.CreateTopupRequest{
		CardNumber:  "",
		TopupAmount: 40000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/topups/create", bytes.NewReader(bodyBytes))
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

func (suite *TopupHandlerTestSuite) TestUpdateTopup_Success() {
	id := 1

	requestBody := requests.UpdateTopupRequest{
		TopupID:     &id,
		CardNumber:  "1234567890123456",
		TopupAmount: 200000,
		TopupMethod: "alfamart",
	}

	grpcRequest := &pb.UpdateTopupRequest{
		TopupId:     int32(id),
		CardNumber:  requestBody.CardNumber,
		TopupAmount: int32(requestBody.TopupAmount),
		TopupMethod: requestBody.TopupMethod,
	}

	grpcResponse := &pb.ApiResponseTopup{
		Status: "success",
		Data: &pb.TopupResponse{
			Id:          int32(id),
			TopupAmount: int32(requestBody.TopupAmount),
		},
	}
	expectedApiResponse := &response.ApiResponseTopup{
		Status: "success",
		Data: &response.TopupResponse{
			ID:          id,
			TopupAmount: requestBody.TopupAmount,
		},
	}

	suite.MockTopupClient.EXPECT().UpdateTopup(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopup(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/topups/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTopup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(200000), resp.Data.TopupAmount)
}

func (suite *TopupHandlerTestSuite) TestUpdateTopup_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/topups/update/abc", nil)
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

func (suite *TopupHandlerTestSuite) TestUpdateTopup_ValidationError() {
	id := 1
	apiKey := "test-api-key"

	invalidRequestBody := requests.UpdateTopupRequest{
		CardNumber:  "",
		TopupAmount: 40000,
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/topups/update/%d", id), bytes.NewReader(bodyBytes))
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

func (suite *TopupHandlerTestSuite) TestTrashTopup_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTopupDeleteAt{Status: "success", Data: &pb.TopupResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseTopupDeleteAt{Status: "success", Data: &response.TopupResponseDeleteAt{ID: id}}

	suite.MockTopupClient.EXPECT().TrashedTopup(gomock.Any(), &pb.FindByIdTopupRequest{TopupId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopupDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashTopup(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestTrashTopup_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().TrashedTopup(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trashed topup", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashTopup(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestRestoreTopup_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTopupDeleteAt{Status: "success", Data: &pb.TopupResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseTopupDeleteAt{Status: "success", Data: &response.TopupResponseDeleteAt{ID: id}}

	suite.MockTopupClient.EXPECT().RestoreTopup(gomock.Any(), &pb.FindByIdTopupRequest{TopupId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopupDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreTopup(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestRestoreTopup_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().RestoreTopup(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore topup", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreTopup(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestDeleteTopupPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseTopupDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseTopupDelete{Status: "success"}

	suite.MockTopupClient.EXPECT().DeleteTopupPermanent(gomock.Any(), &pb.FindByIdTopupRequest{TopupId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopupDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/topups/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteTopupPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestDeleteTopupPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().DeleteTopupPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete topup", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/topups/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteTopupPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestRestoreAllTopup_Success() {
	grpcResponse := &pb.ApiResponseTopupAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseTopupAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully restored all topup").Times(1)

	suite.MockTopupClient.EXPECT().RestoreAllTopup(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopupAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllTopup(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestRestoreAllTopup_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().RestoreAllTopup(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all topup", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllTopup(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestDeleteAllTopupPermanent_Success() {
	grpcResponse := &pb.ApiResponseTopupAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseTopupAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully deleted all topup permanently").Times(1)

	suite.MockTopupClient.EXPECT().DeleteAllTopupPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseTopupAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/trashed/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllTopupPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *TopupHandlerTestSuite) TestDeleteAllTopupPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockTopupClient.EXPECT().DeleteAllTopupPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all topup", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/topups/trashed/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllTopupPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestTopupHandlerSuite(t *testing.T) {
	suite.Run(t, new(TopupHandlerTestSuite))
}
