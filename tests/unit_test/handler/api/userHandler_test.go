package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindAllUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	expectedResponse := &pb.ApiResponsePaginationUser{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data: []*pb.UserResponse{
			{
				Id:        1,
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john.doe@example.com",
			},
			{
				Id:        2,
				Firstname: "Jane",
				Lastname:  "Doe",
				Email:     "jane.doe@example.com",
			},
		},
		Pagination: &pb.PaginationMeta{
			CurrentPage: 1,
			PageSize:    2,
			TotalPages:  1,
		},
	}

	mockUserClient.EXPECT().
		FindAll(gomock.Any(), &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}).
		Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.FindAllUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Users retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)
}

func TestFindAllUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockUserClient.EXPECT().
		FindAll(gomock.Any(), &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}).
		Return(nil, fmt.Errorf("internal server error"))

	mockLogger.EXPECT().Debug("Failed to retrieve user data", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.FindAllUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve user data: ", resp.Message)
}

func TestFindAllUser_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	mockUserClient.EXPECT().
		FindAll(gomock.Any(), &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}).
		Return(&pb.ApiResponsePaginationUser{
			Status:     "success",
			Message:    "No users found",
			Data:       []*pb.UserResponse{},
			Pagination: &pb.PaginationMeta{},
		}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.FindAllUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "No users found", resp.Message)
	assert.Len(t, resp.Data, 0)
}

func TestFindByIdUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockUserClient.EXPECT().
		FindById(gomock.Any(), &pb.FindByIdUserRequest{
			Id: int32(id),
		}).
		Return(&pb.ApiResponseUser{
			Status:  "success",
			Message: "User retrieved successfully",
			Data: &pb.UserResponse{
				Id:        int32(id),
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john.doe@example.com",
			},
		}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)
	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User retrieved successfully", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestFindByIdUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockUserClient.EXPECT().
		FindById(gomock.Any(), &pb.FindByIdUserRequest{
			Id: int32(id),
		}).
		Return(nil, fmt.Errorf("gRPC service unavailable")).Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve user data", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)
	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Failed to retrieve user data")
}

func TestFindByIdUser_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	id := 1

	mockUserClient.EXPECT().
		FindById(gomock.Any(), &pb.FindByIdUserRequest{
			Id: int32(id),
		}).
		Return(&pb.ApiResponseUser{
			Status:  "success",
			Message: "User not found",
			Data:    nil,
		}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)
	err := handler.FindById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User not found", resp.Message)
}

func TestFindByActiveUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllUserRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockUserClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(&pb.ApiResponsePaginationUserDeleteAt{
			Status:  "success",
			Message: "Users retrieved successfully",
			Data: []*pb.UserResponseWithDeleteAt{
				{
					Id:        1,
					Firstname: "John",
					Lastname:  "Doe",
					CreatedAt: "2023-08-01T10:00:00Z",
					UpdatedAt: "2023-08-01T10:00:00Z",
				},
				{
					Id:        2,
					Firstname: "Jane",
					Lastname:  "Doe",
					CreatedAt: "2023-08-01T10:00:00Z",
					UpdatedAt: "2023-08-01T10:00:00Z",
				},
			},
		}, nil).
		Times(1)

	e := echo.New()
	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/users/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.FindByActive(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsesUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Users retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)

	users := resp.Data

	assert.Equal(t, "John", users[0].Firstname)
	assert.Equal(t, "Jane", users[1].Firstname)
}

func TestFindByActiveUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllUserRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockUserClient.EXPECT().
		FindByActive(gomock.Any(), request).
		Return(nil, fmt.Errorf("failed to retrieve user data")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve user data", gomock.Any()).Times(1)

	e := echo.New()
	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/users/active", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.FindByActive(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve user data: ", resp.Message)
}

func TestFindByTrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllUserRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockUserClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(&pb.ApiResponsePaginationUserDeleteAt{
			Status:  "success",
			Message: "Trashed users retrieved successfully",
			Data: []*pb.UserResponseWithDeleteAt{
				{
					Id:        1,
					Firstname: "John",
					Lastname:  "Doe",
					CreatedAt: "2023-08-01T10:00:00Z",
					UpdatedAt: "2023-08-01T10:00:00Z",
				},
				{
					Id:        2,
					Firstname: "Jane",
					Lastname:  "Doe",
					CreatedAt: "2023-08-01T10:00:00Z",
					UpdatedAt: "2023-08-01T10:00:00Z",
				},
			},
		}, nil).
		Times(1)

	e := echo.New()
	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/users/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.FindByTrashed(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponsePaginationUserDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Trashed users retrieved successfully", resp.Message)
	assert.Len(t, resp.Data, 2)

	users := resp.Data

	assert.Equal(t, "John", users[0].Firstname)
	assert.Equal(t, "Jane", users[1].Firstname)
}

func TestFindByTrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	request := &pb.FindAllUserRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	mockUserClient.EXPECT().
		FindByTrashed(gomock.Any(), request).
		Return(nil, fmt.Errorf("failed to retrieve trashed user data")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to retrieve user data", gomock.Any()).Times(1)

	e := echo.New()
	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/users/trashed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.FindByTrashed(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to retrieve user data: ", resp.Message)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedGRPCRequest := &pb.CreateUserRequest{
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}

	expectedGRPCResponse := &pb.UserResponse{
		Id:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
	}

	expectedAPIResponse := &pb.ApiResponseUser{
		Status:  "success",
		Message: "User created successfully",
		Data:    expectedGRPCResponse,
	}

	mockUserClient.EXPECT().
		Create(gomock.Any(), expectedGRPCRequest).
		Return(expectedAPIResponse, nil).
		Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err = handler.Create(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User created successfully", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, expectedGRPCResponse, resp.Data)
}

func TestCreateUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	expectedGRPCRequest := &pb.CreateUserRequest{
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}

	mockUserClient.EXPECT().
		Create(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to create user", gomock.Any()).Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err = handler.Create(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to create user: ", resp.Message)
}

func TestCreateUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.CreateUserRequest{
		FirstName: "",
		LastName:  "",
		Email:     "invalid-email",
	}

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockLogger.EXPECT().Debug("Validation Error", gomock.Any()).Times(1)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err = handler.Create(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error:")
}

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	expectedGRPCRequest := &pb.UpdateUserRequest{
		Id:              int32(requestBody.UserID),
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}

	expectedGRPCResponse := &pb.ApiResponseUser{
		Status:  "success",
		Message: "User updated successfully",
		Data: &pb.UserResponse{
			Id:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	mockUserClient.EXPECT().
		Update(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPut, "/api/users/update/1", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err = handler.Update(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User updated successfully", resp.Message)
	assert.Equal(t, "John", resp.Data.Firstname)
	assert.Equal(t, "Doe", resp.Data.Lastname)
}

func TestUpdateUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.UpdateUserRequest{
		UserID:          1,
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	expectedGRPCRequest := &pb.UpdateUserRequest{
		Id:              int32(requestBody.UserID),
		Firstname:       requestBody.FirstName,
		Lastname:        requestBody.LastName,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		ConfirmPassword: requestBody.ConfirmPassword,
	}

	mockUserClient.EXPECT().
		Update(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to update user", gomock.Any()).Times(1)

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPut, "/api/users/update/1", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err = handler.Update(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to update user: ", resp.Message)
}

func TestUpdateUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	requestBody := requests.UpdateUserRequest{
		UserID:    1,
		FirstName: "",
		LastName:  "",
		Email:     "invalid-email",
	}

	e := echo.New()
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPut, "/api/users/update/1", bytes.NewReader(requestBodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockLogger.EXPECT().Debug("Validation Error", gomock.Any()).Times(1)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err = handler.Update(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Message, "Validation Error:")
}

func TestTrashedUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := 1
	expectedGRPCRequest := &pb.FindByIdUserRequest{
		Id: int32(userID),
	}

	expectedGRPCResponse := &pb.ApiResponseUser{
		Status:  "success",
		Message: "User trashed successfully",
		Data: &pb.UserResponse{
			Id:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	mockUserClient.EXPECT().
		TrashedUser(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.TrashedUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User trashed successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "John", resp.Data.Firstname)
}

func TestTrashedUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := 1
	expectedGRPCRequest := &pb.FindByIdUserRequest{
		Id: int32(userID),
	}

	mockUserClient.EXPECT().
		TrashedUser(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to trashed user", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/trashed/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.TrashedUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to trashed user: ", resp.Message)
}

func TestTrashedUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/trashed/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Invalid user ID", gomock.Any()).Times(1)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.TrashedUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid user ID", resp.Message)
}

func TestRestoreUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := 1
	expectedGRPCRequest := &pb.FindByIdUserRequest{
		Id: int32(userID),
	}

	expectedGRPCResponse := &pb.ApiResponseUser{
		Status:  "success",
		Message: "User restored successfully",
		Data: &pb.UserResponse{
			Id:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	mockUserClient.EXPECT().
		RestoreUser(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.RestoreUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User restored successfully", resp.Message)
	assert.Equal(t, int32(1), resp.Data.Id)
	assert.Equal(t, "John", resp.Data.Firstname)
}

func TestRestoreUser_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := 1
	expectedGRPCRequest := &pb.FindByIdUserRequest{
		Id: int32(userID),
	}

	mockUserClient.EXPECT().
		RestoreUser(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to restore user", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/restore/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.RestoreUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to restore user: ", resp.Message)
}

func TestRestoreUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/restore/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Invalid user ID", gomock.Any()).Times(1)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.RestoreUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid user ID", resp.Message)
}

func TestDeleteUserPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := 1
	expectedGRPCRequest := &pb.FindByIdUserRequest{
		Id: int32(userID),
	}

	expectedGRPCResponse := &pb.ApiResponseUserDelete{
		Status:  "success",
		Message: "User deleted permanently",
	}

	mockUserClient.EXPECT().
		DeleteUserPermanent(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/users/delete-permanent/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.DeleteUserPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.ApiResponseUserDelete
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "User deleted permanently", resp.Message)
}

func TestDeleteUserPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	userID := 1
	expectedGRPCRequest := &pb.FindByIdUserRequest{
		Id: int32(userID),
	}

	mockUserClient.EXPECT().
		DeleteUserPermanent(gomock.Any(), expectedGRPCRequest).
		Return(nil, errors.New("internal server error")).
		Times(1)

	mockLogger.EXPECT().Debug("Failed to delete user", gomock.Any()).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/users/delete-permanent/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.DeleteUserPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Failed to delete user: ", resp.Message)
}

func TestDeleteUserPermanent_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mock_pb.NewMockUserServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/users/delete-permanent/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	mockLogger.EXPECT().Debug("Invalid user ID", gomock.Any()).Times(1)

	handler := api.NewHandlerUser(mockUserClient, e, mockLogger)

	err := handler.DeleteUserPermanent(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "Invalid user ID", resp.Message)
}
