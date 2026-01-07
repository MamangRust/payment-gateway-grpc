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

type RoleHandlerTestSuite struct {
	suite.Suite
	Ctrl           *gomock.Controller
	MockRoleClient *mock_pb.MockRoleServiceClient
	MockLogger     *mock_logger.MockLoggerInterface
	MockMapper     *mock_apimapper.MockRoleResponseMapper
	E              *echo.Echo
	Handler        *api.RoleHandleApi
}

func (suite *RoleHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockRoleClient = mock_pb.NewMockRoleServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockRoleResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerRole(suite.MockRoleClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *RoleHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *RoleHandlerTestSuite) TestFindAll_Success() {
	grpcResponse := &pb.ApiResponsePaginationRole{
		Status:  "success",
		Message: "Roles retrieved successfully",
		Data: []*pb.RoleResponse{
			{Id: 1, Name: "Admin"},
			{Id: 2, Name: "User"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationRole{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.RoleResponse{{ID: 1, Name: "Admin"}, {ID: 2, Name: "User"}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockRoleClient.EXPECT().FindAllRole(gomock.Any(), &pb.FindAllRoleRequest{Page: 1, PageSize: 10, Search: "admin"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationRole(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/role?page=1&page_size=10&search=admin", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationRole
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *RoleHandlerTestSuite) TestFindAll_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().FindAllRole(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch role records", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/role", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestFindById_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseRole{Status: "success", Data: &pb.RoleResponse{Id: int32(id), Name: "Admin"}}
	expectedApiResponse := &response.ApiResponseRole{Status: "success", Data: &response.RoleResponse{ID: id, Name: "Admin"}}

	suite.MockRoleClient.EXPECT().FindByIdRole(gomock.Any(), &pb.FindByIdRoleRequest{RoleId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRole(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/role/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseRole
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *RoleHandlerTestSuite) TestFindById_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/role/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestFindByUserId_Success() {
	userID := 1
	grpcResponse := &pb.ApiResponsesRole{Status: "success", Data: []*pb.RoleResponse{{Id: 1, Name: "Admin"}}}
	expectedApiResponse := &response.ApiResponsesRole{Status: "success", Data: []*response.RoleResponse{{ID: 1, Name: "Admin"}}}

	suite.MockRoleClient.EXPECT().FindByUserId(gomock.Any(), &pb.FindByIdUserRoleRequest{UserId: int32(userID)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsesRole(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/role/user/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("1")

	err := suite.Handler.FindByUserId(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsesRole
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *RoleHandlerTestSuite) TestFindByUserId_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/role/user/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("abc")

	err := suite.Handler.FindByUserId(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestFindByActive_Success() {
	grpcResponse := &pb.ApiResponsePaginationRoleDeleteAt{Status: "success", Data: []*pb.RoleResponseDeleteAt{{Id: 1, Name: "Admin"}}}
	expectedApiResponse := &response.ApiResponsePaginationRoleDeleteAt{Status: "success", Data: []*response.RoleResponseDeleteAt{{ID: 1, Name: "Admin"}}}

	suite.MockRoleClient.EXPECT().FindByActive(gomock.Any(), &pb.FindAllRoleRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationRoleDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/role/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationRoleDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *RoleHandlerTestSuite) TestFindByActive_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().FindByActive(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch active roles", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/role/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestFindByTrashed_Success() {
	grpcResponse := &pb.ApiResponsePaginationRoleDeleteAt{Status: "success", Data: []*pb.RoleResponseDeleteAt{{Id: 2, Name: "Old Role"}}}
	expectedApiResponse := &response.ApiResponsePaginationRoleDeleteAt{Status: "success", Data: []*response.RoleResponseDeleteAt{{ID: 2, Name: "Old Role"}}}

	suite.MockRoleClient.EXPECT().FindByTrashed(gomock.Any(), &pb.FindAllRoleRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationRoleDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/role/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationRoleDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *RoleHandlerTestSuite) TestFindByTrashed_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().FindByTrashed(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to fetch trashed roles", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/role/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestCreate_Success() {
	requestBody := requests.CreateRoleRequest{Name: "Editor"}
	grpcRequest := &pb.CreateRoleRequest{Name: "Editor"}
	grpcResponse := &pb.ApiResponseRole{Status: "success", Data: &pb.RoleResponse{Id: 3, Name: "Editor"}}
	expectedApiResponse := &response.ApiResponseRole{Status: "success", Data: &response.RoleResponse{ID: 3, Name: "Editor"}}

	suite.MockRoleClient.EXPECT().CreateRole(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRole(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseRole
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("Editor", resp.Data.Name)
}

func (suite *RoleHandlerTestSuite) TestCreate_Failure() {
	requestBody := requests.CreateRoleRequest{Name: "Editor"}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create role", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestCreate_ValidationError() {
	invalidRequestBody := requests.CreateRoleRequest{Name: ""}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestUpdate_Success() {
	id := 1
	requestBody := requests.UpdateRoleRequest{Name: "Super Admin"}
	grpcRequest := &pb.UpdateRoleRequest{Id: int32(id), Name: "Super Admin"}
	grpcResponse := &pb.ApiResponseRole{Status: "success", Data: &pb.RoleResponse{Id: int32(id), Name: "Super Admin"}}
	expectedApiResponse := &response.ApiResponseRole{Status: "success", Data: &response.RoleResponse{ID: id, Name: "Super Admin"}}

	suite.MockRoleClient.EXPECT().UpdateRole(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRole(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseRole
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("Super Admin", resp.Data.Name)
}

func (suite *RoleHandlerTestSuite) TestUpdate_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/role/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestUpdate_ValidationError() {
	id := 1
	invalidRequestBody := requests.UpdateRoleRequest{Name: ""}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestTrash_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseRoleDeleteAt{Status: "success", Data: &pb.RoleResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseRoleDeleteAt{Status: "success", Data: &response.RoleResponseDeleteAt{ID: id}}

	suite.MockRoleClient.EXPECT().TrashedRole(gomock.Any(), &pb.FindByIdRoleRequest{RoleId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRoleDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/role/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Trashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestTrash_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().TrashedRole(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trash role", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/role/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Trashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestRestore_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseRoleDeleteAt{Status: "success", Data: &pb.RoleResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseRoleDeleteAt{Status: "success", Data: &response.RoleResponseDeleteAt{ID: id}}

	suite.MockRoleClient.EXPECT().RestoreRole(gomock.Any(), &pb.FindByIdRoleRequest{RoleId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRoleDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/role/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Restore(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestRestore_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().RestoreRole(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore role", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/role/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.Restore(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestDeletePermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseRoleDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseRoleDelete{Status: "success"}

	suite.MockRoleClient.EXPECT().DeleteRolePermanent(gomock.Any(), &pb.FindByIdRoleRequest{RoleId: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRoleDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/role/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeletePermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestDeletePermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().DeleteRolePermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete role permanently", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/role/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeletePermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestRestoreAll_Success() {
	grpcResponse := &pb.ApiResponseRoleAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseRoleAll{Status: "success"}

	suite.MockRoleClient.EXPECT().RestoreAllRole(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRoleAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/role/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestRestoreAll_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().RestoreAllRole(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore all roles", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/role/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestDeleteAllPermanent_Success() {
	grpcResponse := &pb.ApiResponseRoleAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseRoleAll{Status: "success"}

	suite.MockRoleClient.EXPECT().DeleteAllRolePermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseRoleAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/role/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *RoleHandlerTestSuite) TestDeleteAllPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockRoleClient.EXPECT().DeleteAllRolePermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete all roles permanently", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/role/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestRoleHandlerSuite(t *testing.T) {
	suite.Run(t, new(RoleHandlerTestSuite))
}
