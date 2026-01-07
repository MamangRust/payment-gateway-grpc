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

type UserHandlerTestSuite struct {
	suite.Suite
	Ctrl           *gomock.Controller
	MockUserClient *mock_pb.MockUserServiceClient
	MockLogger     *mock_logger.MockLoggerInterface
	MockMapper     *mock_apimapper.MockUserResponseMapper
	E              *echo.Echo
	Handler        *api.UserHandleApi
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockUserClient = mock_pb.NewMockUserServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.MockMapper = mock_apimapper.NewMockUserResponseMapper(suite.Ctrl)
	suite.E = echo.New()
	suite.Handler = api.NewHandlerUser(suite.MockUserClient, suite.E, suite.MockLogger, suite.MockMapper)
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *UserHandlerTestSuite) TestFindAllUser_Success() {
	grpcResponse := &pb.ApiResponsePaginationUser{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data: []*pb.UserResponse{
			{Id: 1, Firstname: "John", Lastname: "Doe", Email: "john.doe@example.com"},
			{Id: 2, Firstname: "Jane", Lastname: "Doe", Email: "jane.doe@example.com"},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}
	expectedApiResponse := &response.ApiResponsePaginationUser{
		Status:     grpcResponse.Status,
		Message:    grpcResponse.Message,
		Data:       []*response.UserResponse{{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, {ID: 2, FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com"}},
		Pagination: &response.PaginationMeta{CurrentPage: 1, PageSize: 2, TotalPages: 1},
	}

	suite.MockUserClient.EXPECT().FindAll(gomock.Any(), &pb.FindAllUserRequest{Page: 1, PageSize: 10, Search: "john"}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationUser(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/user?page=1&page_size=10&search=john", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAllUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 2)
}

func (suite *UserHandlerTestSuite) TestFindAllUser_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve user data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAllUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestFindByIdUser_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseUser{Status: "success", Data: &pb.UserResponse{Id: int32(id), Firstname: "John", Lastname: "Doe", Email: "john.doe@example.com"}}
	expectedApiResponse := &response.ApiResponseUser{Status: "success", Data: &response.UserResponse{ID: id, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}}

	suite.MockUserClient.EXPECT().FindById(gomock.Any(), &pb.FindByIdUserRequest{Id: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUser(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/user/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal(int(id), resp.Data.ID)
}

func (suite *UserHandlerTestSuite) TestFindByIdUser_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/user/abc", nil)
	rec := httptest.NewRecorder()

	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	suite.MockLogger.EXPECT().
		Debug("Invalid user ID", gomock.Any()).
		Times(1)

	err := suite.Handler.FindById(c)

	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *UserHandlerTestSuite) TestFindByActiveUser_Success() {
	grpcResponse := &pb.ApiResponsePaginationUserDeleteAt{Status: "success", Data: []*pb.UserResponseDeleteAt{{Id: 1, Firstname: "John"}}}
	expectedApiResponse := &response.ApiResponsePaginationUserDeleteAt{Status: "success", Data: []*response.UserResponseDeleteAt{{ID: 1, FirstName: "John"}}}

	suite.MockUserClient.EXPECT().FindByActive(gomock.Any(), &pb.FindAllUserRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationUserDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/user/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationUserDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *UserHandlerTestSuite) TestFindByActiveUser_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().FindByActive(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve user data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/user/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActive(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestFindByTrashedUser_Success() {
	grpcResponse := &pb.ApiResponsePaginationUserDeleteAt{Status: "success", Data: []*pb.UserResponseDeleteAt{{Id: 2, Firstname: "Old User"}}}
	expectedApiResponse := &response.ApiResponsePaginationUserDeleteAt{Status: "success", Data: []*response.UserResponseDeleteAt{{ID: 2, FirstName: "Old User"}}}

	suite.MockUserClient.EXPECT().FindByTrashed(gomock.Any(), &pb.FindAllUserRequest{Page: 1, PageSize: 10, Search: ""}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponsePaginationUserDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodGet, "/api/user/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationUserDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Len(resp.Data, 1)
}

func (suite *UserHandlerTestSuite) TestFindByTrashedUser_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().FindByTrashed(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve user data", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/user/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashed(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestCreateUser_Success() {
	requestBody := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	grpcRequest := &pb.CreateUserRequest{
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}
	grpcResponse := &pb.ApiResponseUser{Status: "success", Data: &pb.UserResponse{Id: 3, Firstname: "John", Lastname: "Doe", Email: "john.doe@example.com"}}
	expectedApiResponse := &response.ApiResponseUser{Status: "success", Data: &response.UserResponse{ID: 3, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}}

	suite.MockUserClient.EXPECT().Create(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUser(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("John", resp.Data.FirstName)
}

func (suite *UserHandlerTestSuite) TestCreateUser_Failure() {
	requestBody := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create user", zap.Error(grpcError)).Times(1)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.Create(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestCreateUser_ValidationError() {
	invalidRequestBody := requests.CreateUserRequest{
		FirstName: "",
		Email:     "invalid-email",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewReader(bodyBytes))
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

func (suite *UserHandlerTestSuite) TestUpdateUser_Success() {
	id := 1
	requestBody := requests.UpdateUserRequest{
		FirstName:       "John",
		LastName:        "Smith",
		Email:           "john.smith@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}
	grpcRequest := &pb.UpdateUserRequest{
		Id:              int32(id),
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}
	grpcResponse := &pb.ApiResponseUser{Status: "success", Data: &pb.UserResponse{Id: int32(id), Firstname: "John", Lastname: "Smith", Email: "john.smith@example.com"}}
	expectedApiResponse := &response.ApiResponseUser{Status: "success", Data: &response.UserResponse{ID: id, FirstName: "John", LastName: "Smith", Email: "john.smith@example.com"}}

	suite.MockUserClient.EXPECT().Update(gomock.Any(), grpcRequest).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUser(grpcResponse).Return(expectedApiResponse)

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user/update/%d", id), bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("Smith", resp.Data.LastName)
}

func (suite *UserHandlerTestSuite) TestUpdateUser_InvalidID() {
	req := httptest.NewRequest(http.MethodPost, "/api/user/update/abc", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := suite.Handler.Update(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *UserHandlerTestSuite) TestUpdateUser_ValidationError() {
	id := 1
	invalidRequestBody := requests.UpdateUserRequest{
		FirstName: "",
		Email:     "invalid-email",
	}

	bodyBytes, _ := json.Marshal(invalidRequestBody)
	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/api/user/update/%d", id),
		bytes.NewReader(bodyBytes),
	)
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

func (suite *UserHandlerTestSuite) TestTrashedUser_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseUserDeleteAt{Status: "success", Data: &pb.UserResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseUserDeleteAt{Status: "success", Data: &response.UserResponseDeleteAt{ID: id}}

	suite.MockUserClient.EXPECT().TrashedUser(gomock.Any(), &pb.FindByIdUserRequest{Id: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUserDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/user/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *UserHandlerTestSuite) TestTrashedUser_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().TrashedUser(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to trashed user", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/user/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.TrashedUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestRestoreUser_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseUserDeleteAt{Status: "success", Data: &pb.UserResponseDeleteAt{Id: int32(id)}}
	expectedApiResponse := &response.ApiResponseUserDeleteAt{Status: "success", Data: &response.UserResponseDeleteAt{ID: id}}

	suite.MockUserClient.EXPECT().RestoreUser(gomock.Any(), &pb.FindByIdUserRequest{Id: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUserDeleteAt(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/user/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *UserHandlerTestSuite) TestRestoreUser_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().RestoreUser(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to restore user", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/user/restore/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.RestoreUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestDeleteUserPermanent_Success() {
	id := 1
	grpcResponse := &pb.ApiResponseUserDelete{Status: "success"}
	expectedApiResponse := &response.ApiResponseUserDelete{Status: "success"}

	suite.MockUserClient.EXPECT().DeleteUserPermanent(gomock.Any(), &pb.FindByIdUserRequest{Id: int32(id)}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUserDelete(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteUserPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *UserHandlerTestSuite) TestDeleteUserPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().DeleteUserPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to delete user", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/permanent/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.DeleteUserPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestRestoreAllUser_Success() {
	grpcResponse := &pb.ApiResponseUserAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseUserAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully restored all user").Times(1)

	suite.MockUserClient.EXPECT().RestoreAllUser(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUserAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/user/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *UserHandlerTestSuite) TestRestoreAllUser_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().RestoreAllUser(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to restore all user", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/user/restore/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.RestoreAllUser(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *UserHandlerTestSuite) TestDeleteAllUserPermanent_Success() {
	grpcResponse := &pb.ApiResponseUserAll{Status: "success"}
	expectedApiResponse := &response.ApiResponseUserAll{Status: "success"}

	suite.MockLogger.EXPECT().Debug("Successfully deleted all user permanently").Times(1)

	suite.MockUserClient.EXPECT().DeleteAllUserPermanent(gomock.Any(), &emptypb.Empty{}).Return(grpcResponse, nil)
	suite.MockMapper.EXPECT().ToApiResponseUserAll(grpcResponse).Return(expectedApiResponse)

	req := httptest.NewRequest(http.MethodPost, "/api/user/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllUserPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *UserHandlerTestSuite) TestDeleteAllUserPermanent_Failure() {
	grpcError := errors.New("gRPC service unavailable")
	suite.MockUserClient.EXPECT().DeleteAllUserPermanent(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Error("Failed to permanently delete all user", zap.Error(grpcError)).Times(1)

	req := httptest.NewRequest(http.MethodPost, "/api/user/permanent/all", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.DeleteAllUserPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
