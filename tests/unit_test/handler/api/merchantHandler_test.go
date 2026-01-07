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

type MerchantHandlerTestSuite struct {
	suite.Suite
	Ctrl                  *gomock.Controller
	MockMerchantClient    *mock_pb.MockMerchantServiceClient
	MockTransactionClient *mock_pb.MockTransactionServiceClient
	MockLogger            *mock_logger.MockLoggerInterface
	MockMapper            *mock_apimapper.MockMerchantResponseMapper
	E                     *echo.Echo
	Handler               *api.MerchantHandleApi
}

func (suite *MerchantHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockMerchantClient = mock_pb.NewMockMerchantServiceClient(suite.Ctrl)
	suite.MockTransactionClient = mock_pb.NewMockTransactionServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockMerchantResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerMerchant(suite.MockMerchantClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *MerchantHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *MerchantHandlerTestSuite) TestFindAllMerchant_Success() {
	grpcResponse := &pb.ApiResponsePaginationMerchant{
		Status:  "success",
		Message: "Merchants retrieved successfully",
		Data: []*pb.MerchantResponse{
			{Id: 1, Name: "Merchant 1", ApiKey: "api-key-1"},
			{Id: 2, Name: "Merchant 2", ApiKey: "api-key-2"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationMerchant{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.MerchantResponse{{ID: 1, Name: "Merchant 1", ApiKey: "api-key-1"}, {ID: 2, Name: "Merchant 2", ApiKey: "api-key-2"}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockMerchantClient.EXPECT().FindAllMerchant(gomock.Any(), &pb.FindAllMerchantRequest{Page: 1, PageSize: 10, Search: "test"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesMerchant(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants?page=1&page_size=10&search=test", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *MerchantHandlerTestSuite) TestFindAllMerchant_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().FindAllMerchant(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve merchant data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestFindByIdMerchant_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseMerchant{Status: "success", Data: &pb.MerchantResponse{Id: int32(id), Name: "Merchant 1", ApiKey: "api-key-1"}}
	expectedApiResponse := &response.ApiResponseMerchant{Status: "success", Data: &response.MerchantResponse{ID: id, Name: "Merchant 1", ApiKey: "api-key-1"}}

	suite.MockMerchantClient.EXPECT().FindByIdMerchant(gomock.Any(), &pb.FindByIdMerchantRequest{MerchantId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchant(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *MerchantHandlerTestSuite) TestFindByIdMerchant_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/merchants/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.FindById(c)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestFindByApiKey_Success() {
	apiKey := "test-api-key-123"

	grpcResponse := &pb.ApiResponseMerchant{
		Status: "success",
		Data: &pb.MerchantResponse{
			Id:     1,
			Name:   "Merchant 1",
			ApiKey: apiKey,
		},
	}

	expectedApiResponse := &response.ApiResponseMerchant{
		Status: "success",
		Data: &response.MerchantResponse{
			ID:     1,
			Name:   "Merchant 1",
			ApiKey: apiKey,
		},
	}

	suite.MockMerchantClient.EXPECT().
		FindByApiKey(
			gomock.Any(),
			&pb.FindByApiKeyRequest{ApiKey: apiKey},
		).
		Return(grpcResponse, nil).
		Times(1)

	suite.MockMapper.EXPECT().
		ToApiResponseMerchant(grpcResponse).
		Return(expectedApiResponse).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/merchants/api-key?api_key="+apiKey,
		nil,
	)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByApiKey(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(apiKey, resp.Data.ApiKey)
}

func (suite *MerchantHandlerTestSuite) TestFindByApiKey_Failure() {
	apiKey := "test-api-key-123"
	grpcError := errors.New("gRPC service unavailable")

	suite.MockMerchantClient.EXPECT().
		FindByApiKey(
			gomock.Any(),
			&pb.FindByApiKeyRequest{ApiKey: apiKey},
		).
		Return(nil, grpcError).
		Times(1)

	suite.MockLogger.EXPECT().
		Debug("Failed to retrieve merchant data", gomock.Any()).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/merchants/api-key?api_key="+apiKey,
		nil,
	)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByApiKey(c)

	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestFindByMerchantUserId_Success() {
	userId := int32(1)

	grpcResponse := &pb.ApiResponsesMerchant{
		Status: "success",
		Data: []*pb.MerchantResponse{
			{Id: 1, Name: "Merchant 1", UserId: userId},
		},
	}

	expectedApiResponse := &response.ApiResponsesMerchant{
		Status: "success",
		Data: []*response.MerchantResponse{
			{ID: 1, Name: "Merchant 1", UserID: int(userId)},
		},
	}

	suite.MockMerchantClient.EXPECT().
		FindByMerchantUserId(
			gomock.Any(),
			&pb.FindByMerchantUserIdRequest{UserId: userId},
		).
		Return(grpcResponse, nil).
		Times(1)

	suite.MockMapper.EXPECT().
		ToApiResponseMerchants(grpcResponse).
		Return(expectedApiResponse).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/merchants/merchant-user",
		nil,
	)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("user_id", userId)

	err := suite.Handler.FindByMerchantUserId(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsesMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *MerchantHandlerTestSuite) TestFindByMerchantUserId_Failure() {
	userId := int32(1)
	grpcError := errors.New("gRPC service unavailable")

	suite.MockMerchantClient.EXPECT().
		FindByMerchantUserId(
			gomock.Any(),
			&pb.FindByMerchantUserIdRequest{UserId: userId},
		).
		Return(nil, grpcError).
		Times(1)

	suite.MockLogger.EXPECT().
		Debug("Failed to retrieve merchant data", gomock.Any()).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/merchants/merchant-user",
		nil,
	)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.Set("user_id", userId)

	err := suite.Handler.FindByMerchantUserId(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestFindByActiveMerchant_Success() {
	grpcResponse := &pb.ApiResponsePaginationMerchantDeleteAt{Status: "success", Data: []*pb.MerchantResponseDeleteAt{{Id: 1, Name: "Merchant 1"}}}
	expectedApiResponse := &response.ApiResponsePaginationMerchantDeleteAt{Status: "success", Data: []*response.MerchantResponseDeleteAt{{ID: 1, Name: "Merchant 1"}}}

	suite.MockMerchantClient.EXPECT().FindByActive(gomock.Any(), &pb.FindAllMerchantRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesMerchantDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationMerchantDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *MerchantHandlerTestSuite) TestFindByActiveMerchant_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().FindByActive(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve merchant data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestFindByTrashedMerchant_Success() {
	grpcResponse := &pb.ApiResponsePaginationMerchantDeleteAt{Status: "success", Data: []*pb.MerchantResponseDeleteAt{{Id: 2, Name: "Merchant 2"}}}
	expectedApiResponse := &response.ApiResponsePaginationMerchantDeleteAt{Status: "success", Data: []*response.MerchantResponseDeleteAt{{ID: 2, Name: "Merchant 2"}}}

	suite.MockMerchantClient.EXPECT().FindByTrashed(gomock.Any(), &pb.FindAllMerchantRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesMerchantDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationMerchantDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *MerchantHandlerTestSuite) TestFindByTrashedMerchant_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().FindByTrashed(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve merchant data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/merchants/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestCreateMerchant_Success() {
	userId := 1
	requestBody := requests.CreateMerchantRequest{
		Name:   "Test Merchant",
		UserID: userId,
	}

	grpcRequest := &pb.CreateMerchantRequest{
		Name:   requestBody.Name,
		UserId: int32(requestBody.UserID),
	}

	grpcResponse := &pb.ApiResponseMerchant{
		Status: "success",
		Data: &pb.MerchantResponse{
			Id:     3,
			Name:   requestBody.Name,
			UserId: int32(userId),
		},
	}

	expectedApiResponse := &response.ApiResponseMerchant{
		Status: "success",
		Data: &response.MerchantResponse{
			ID:     3,
			Name:   requestBody.Name,
			UserID: userId,
		},
	}

	suite.MockMerchantClient.EXPECT().CreateMerchant(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchant(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/merchants/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(3, resp.Data.ID)
}

func (suite *MerchantHandlerTestSuite) TestCreateMerchant_Failure() {
	userId := 1
	requestBody := requests.CreateMerchantRequest{
		Name:   "Test Merchant",
		UserID: userId,
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().CreateMerchant(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create merchant", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/merchants/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestCreateMerchant_ValidationError() {
	invalidRequestBody := requests.CreateMerchantRequest{
		Name: "",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/merchants/create", bytes.NewReader(bodyBytes))
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

func (suite *MerchantHandlerTestSuite) TestUpdateMerchant_Success() {
	id := 1
	userId := 1
	requestBody := requests.UpdateMerchantRequest{
		Name:   "Updated Merchant",
		UserID: userId,
		Status: "active",
	}

	grpcRequest := &pb.UpdateMerchantRequest{
		MerchantId: int32(id),
		Name:       requestBody.Name,
		UserId:     int32(requestBody.UserID),
		Status:     requestBody.Status,
	}

	grpcResponse := &pb.ApiResponseMerchant{
		Status: "success",
		Data: &pb.MerchantResponse{
			Id:     int32(id),
			Name:   requestBody.Name,
			Status: requestBody.Status,
		},
	}
	expectedApiResponse := &response.ApiResponseMerchant{
		Status: "success",
		Data: &response.MerchantResponse{
			ID:     id,
			Name:   requestBody.Name,
			Status: requestBody.Status,
		},
	}

	suite.MockMerchantClient.EXPECT().UpdateMerchant(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchant(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchants/updates/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseMerchant
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("Updated Merchant", resp.Data.Name)
}

func (suite *MerchantHandlerTestSuite) TestUpdateMerchant_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/merchants/updates/abc", nil)
	// suite.MockLogger.EXPECT().
	// 	Debug("Bad Request", gomock.Any()).
	// 	Times(1)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestUpdateMerchant_ValidationError() {
	id := 1
	invalidRequestBody := requests.UpdateMerchantRequest{
		Name: "",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchants/updates/%d", id), bytes.NewReader(bodyBytes))
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

func (suite *MerchantHandlerTestSuite) TestTrashedMerchant_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseMerchantDeleteAt{Status: "success", Data: &pb.MerchantResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseMerchantDeleteAt{Status: "success", Data: &response.MerchantResponseDeleteAt{ID: id}}

	suite.MockMerchantClient.EXPECT().TrashedMerchant(gomock.Any(), &pb.FindByIdMerchantRequest{MerchantId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchantDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedMerchant(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestTrashedMerchant_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().TrashedMerchant(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trashed merchant", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedMerchant(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestRestoreMerchant_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseMerchantDeleteAt{Status: "success", Data: &pb.MerchantResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseMerchantDeleteAt{Status: "success", Data: &response.MerchantResponseDeleteAt{ID: id}}

	suite.MockMerchantClient.EXPECT().RestoreMerchant(gomock.Any(), &pb.FindByIdMerchantRequest{MerchantId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchantDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreMerchant(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestRestoreMerchant_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().RestoreMerchant(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore merchant", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreMerchant(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestDeleteMerchantPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseMerchantDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseMerchantDelete{Status: "success"}

	suite.MockMerchantClient.EXPECT().DeleteMerchantPermanent(gomock.Any(), &pb.FindByIdMerchantRequest{MerchantId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchantDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/merchants/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Delete(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestDeleteMerchantPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().DeleteMerchantPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete merchant", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/merchants/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Delete(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestRestoreAllMerchant_Success() {
	grpcResponse := &pb.ApiResponseMerchantAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseMerchantAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully restored all merchant").Times(1)

	suite.MockMerchantClient.EXPECT().RestoreAllMerchant(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchantAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllMerchant(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestRestoreAllMerchant_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().RestoreAllMerchant(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all merchant", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllMerchant(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestDeleteAllMerchantPermanent_Success() {
	grpcResponse := &pb.ApiResponseMerchantAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseMerchantAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully deleted all merchant permanently").Times(1)

	suite.MockMerchantClient.EXPECT().DeleteAllMerchantPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseMerchantAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllMerchantPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MerchantHandlerTestSuite) TestDeleteAllMerchantPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockMerchantClient.EXPECT().DeleteAllMerchantPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all merchant", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/merchants/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllMerchantPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestMerchantHandlerSuite(t *testing.T) {
	suite.Run(t, new(MerchantHandlerTestSuite))
}
