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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TransferHandlerTestSuite struct {
	suite.Suite
	Ctrl               *gomock.Controller
	MockTransferClient *mock_pb.MockTransferServiceClient
	MockLogger         *mock_logger.MockLoggerInterface
	E                  *echo.Echo
	Mapper             *mock_apimapper.MockTransferResponseMapper
	Handler            *api.TransferHandleApi
}

func (suite *TransferHandlerTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockTransferClient = mock_pb.NewMockTransferServiceClient(suite.Ctrl)
	suite.MockLogger = mock_logger.NewMockLoggerInterface(suite.Ctrl)
	suite.E = echo.New()
	suite.Mapper = mock_apimapper.NewMockTransferResponseMapper(suite.Ctrl)
	suite.Handler = api.NewHandlerTransfer(suite.MockTransferClient, suite.E, suite.MockLogger, suite.Mapper)
}

func (suite *TransferHandlerTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *TransferHandlerTestSuite) TestFindAllTransfer_Success() {
	grpcResp := &pb.ApiResponsePaginationTransfer{
		Status:  "success",
		Message: "Successfully fetch transfers",
		Pagination: &pb.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
		},
		Data: []*pb.TransferResponse{
			{
				Id:             1,
				TransferFrom:   "test",
				TransferTo:     "test",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
				CreatedAt:      "2022-01-01 00:00:00",
				UpdatedAt:      "2022-01-01 00:00:00",
			},
		},
	}

	apiResp := &response.ApiResponsePaginationTransfer{
		Status:  "success",
		Message: "Successfully fetch transfers",
		Pagination: &response.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
		},
		Data: []*response.TransferResponse{
			{
				ID:             1,
				TransferFrom:   "test",
				TransferTo:     "test",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
				CreatedAt:      "2022-01-01 00:00:00",
				UpdatedAt:      "2022-01-01 00:00:00",
			},
		},
	}

	suite.MockTransferClient.EXPECT().
		FindAllTransfer(
			gomock.Any(),
			&pb.FindAllTransferRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponsePaginationTransfer(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/transfer?page=1&page_size=10",
		nil,
	)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Equal("Successfully fetch transfers", resp.Message)
	suite.Len(resp.Data, 1)
	suite.Equal(1, resp.Data[0].ID)
}

func (suite *TransferHandlerTestSuite) TestFindAllTransfer_Failure() {
	grpcError := errors.New("some internal error")

	suite.MockTransferClient.EXPECT().
		FindAllTransfer(
			gomock.Any(),
			&pb.FindAllTransferRequest{Page: 1, PageSize: 10, Search: ""},
		).
		Return(nil, grpcError)

	suite.MockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)

	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("error", resp.Status)
	suite.Contains(resp.Message, "Failed to fetch all transfers")
}

func (suite *TransferHandlerTestSuite) TestFindAllTransfer_Empty() {
	grpcResp := &pb.ApiResponsePaginationTransfer{
		Status:  "success",
		Message: "No transfers found",
		Data:    []*pb.TransferResponse{},
		Pagination: &pb.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   0,
			TotalRecords: 0,
		},
	}

	apiResp := &response.ApiResponsePaginationTransfer{
		Status:  "success",
		Message: "No transfers found",
		Data:    []*response.TransferResponse{},
		Pagination: &response.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   0,
			TotalRecords: 0,
		},
	}

	suite.MockTransferClient.EXPECT().
		FindAllTransfer(
			gomock.Any(),
			&pb.FindAllTransferRequest{
				Page:     1,
				PageSize: 10,
				Search:   "",
			},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponsePaginationTransfer(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindAll(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Equal("No transfers found", resp.Message)
	suite.Len(resp.Data, 0)
}

func (suite *TransferHandlerTestSuite) TestFindByIdTransfer_Success() {
	transferID := 1

	grpcResp := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: &pb.TransferResponse{
			Id:             int32(transferID),
			TransferFrom:   "test",
			TransferTo:     "test",
			TransferAmount: 10000,
			TransferTime:   "2022-01-01 00:00:00",
			CreatedAt:      "2022-01-01 00:00:00",
			UpdatedAt:      "2022-01-01 00:00:00",
		},
	}

	apiResp := &response.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: &response.TransferResponse{
			ID:             transferID,
			TransferFrom:   "test",
			TransferTo:     "test",
			TransferAmount: 10000,
			TransferTime:   "2022-01-01 00:00:00",
			CreatedAt:      "2022-01-01 00:00:00",
			UpdatedAt:      "2022-01-01 00:00:00",
		},
	}

	suite.MockTransferClient.EXPECT().
		FindByIdTransfer(
			gomock.Any(),
			&pb.FindByIdTransferRequest{TransferId: int32(transferID)},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransfer(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.NotNil(resp.Data)
	suite.Equal(transferID, resp.Data.ID)
}

func (suite *TransferHandlerTestSuite) TestFindByIdTransfer_Failure() {
	transferID := 1
	expectedGRPCRequest := &pb.FindByIdTransferRequest{TransferId: int32(transferID)}
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().FindByIdTransfer(gomock.Any(), expectedGRPCRequest).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/1", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("error", resp.Status)
	suite.Contains(resp.Message, "Failed to fetch transfer by ID")
}

func (suite *TransferHandlerTestSuite) TestFindByIdTransfer_InvalidID() {
	req := httptest.NewRequest(http.MethodGet, "/api/transfer/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	suite.MockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	err := suite.Handler.FindById(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("error", resp.Status)
	suite.Equal("Invalid Transfer ID", resp.Message)
}

func (suite *TransferHandlerTestSuite) TestFindByTransferByTransferFrom_Success() {
	transferFrom := "test_user"

	grpcResp := &pb.ApiResponseTransfers{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*pb.TransferResponse{
			{
				Id:             1,
				TransferFrom:   transferFrom,
				TransferTo:     "test_to",
				TransferAmount: 10000,
			},
		},
	}

	apiResp := &response.ApiResponseTransfers{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*response.TransferResponse{
			{
				ID:             1,
				TransferFrom:   transferFrom,
				TransferTo:     "test_to",
				TransferAmount: 10000,
			},
		},
	}

	suite.MockTransferClient.EXPECT().
		FindTransferByTransferFrom(
			gomock.Any(),
			&pb.FindTransferByTransferFromRequest{TransferFrom: transferFrom},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransfers(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/from/test_user", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("transfer_from")
	c.SetParamValues(transferFrom)

	err := suite.Handler.FindByTransferByTransferFrom(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransfers
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 1)
	suite.Equal(transferFrom, resp.Data[0].TransferFrom)
}

func (suite *TransferHandlerTestSuite) TestFindByTransferByTransferFrom_Failure() {
	transferFrom := "test_user"
	expectedGRPCRequest := &pb.FindTransferByTransferFromRequest{TransferFrom: transferFrom}
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().FindTransferByTransferFrom(gomock.Any(), expectedGRPCRequest).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/from/test_user", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("transfer_from")
	c.SetParamValues("test_user")

	err := suite.Handler.FindByTransferByTransferFrom(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("error", resp.Status)
	suite.Contains(resp.Message, "Failed to fetch transfers by transfer_from")
}

func (suite *TransferHandlerTestSuite) TestFindByTransferByTransferTo_Success() {
	transferTo := "test_to"

	grpcResp := &pb.ApiResponseTransfers{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*pb.TransferResponse{
			{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     transferTo,
				TransferAmount: 10000,
			},
		},
	}

	apiResp := &response.ApiResponseTransfers{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*response.TransferResponse{
			{
				ID:             1,
				TransferFrom:   "test_from",
				TransferTo:     transferTo,
				TransferAmount: 10000,
			},
		},
	}

	suite.MockTransferClient.EXPECT().
		FindTransferByTransferTo(
			gomock.Any(),
			&pb.FindTransferByTransferToRequest{TransferTo: transferTo},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransfers(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/to/test_to", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("transfer_to")
	c.SetParamValues(transferTo)

	err := suite.Handler.FindByTransferByTransferTo(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransfers
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 1)
	suite.Equal(transferTo, resp.Data[0].TransferTo)
}

func (suite *TransferHandlerTestSuite) TestFindByTransferByTransferTo_Failure() {
	transferTo := "test_to"
	expectedGRPCRequest := &pb.FindTransferByTransferToRequest{TransferTo: transferTo}
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().FindTransferByTransferTo(gomock.Any(), expectedGRPCRequest).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/to/test_to", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("transfer_to")
	c.SetParamValues("test_to")

	err := suite.Handler.FindByTransferByTransferTo(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("error", resp.Status)
	suite.Contains(resp.Message, "Failed to fetch transfers by transfer_to")
}

func (suite *TransferHandlerTestSuite) TestFindByActiveTransfer_Success() {
	reqPB := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	grpcResp := &pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*pb.TransferResponseDeleteAt{
			{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
			},
		},
	}

	apiResp := &response.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*response.TransferResponseDeleteAt{
			{
				ID:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
			},
		},
	}

	suite.MockTransferClient.EXPECT().
		FindByActiveTransfer(gomock.Any(), reqPB).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponsePaginationTransferDeleteAt(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActiveTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransferDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 1)
}

func (suite *TransferHandlerTestSuite) TestFindByActiveTransfer_Failure() {
	request := &pb.FindAllTransferRequest{Search: "", Page: 1, PageSize: 10}
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().FindByActiveTransfer(gomock.Any(), request).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActiveTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Failed to fetch active transfers")
}

func (suite *TransferHandlerTestSuite) TestFindByActiveTransfer_Empty() {
	reqPB := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	grpcResp := &pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "No active transfers found",
		Data:    []*pb.TransferResponseDeleteAt{},
	}

	apiResp := &response.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "No active transfers found",
		Data:    []*response.TransferResponseDeleteAt{},
	}

	suite.MockTransferClient.EXPECT().
		FindByActiveTransfer(gomock.Any(), reqPB).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponsePaginationTransferDeleteAt(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/active", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByActiveTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransferDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Empty(resp.Data)
}

func (suite *TransferHandlerTestSuite) TestFindByTrashedTransfer_Success() {
	reqPB := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	grpcResp := &pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*pb.TransferResponseDeleteAt{
			{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
			},
		},
	}

	apiResp := &response.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "Transfer retrieved successfully",
		Data: []*response.TransferResponseDeleteAt{
			{
				ID:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
			},
		},
	}

	suite.MockTransferClient.EXPECT().
		FindByTrashedTransfer(gomock.Any(), reqPB).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponsePaginationTransferDeleteAt(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashedTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransferDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Len(resp.Data, 1)
}

func (suite *TransferHandlerTestSuite) TestFindByTrashedTransfer_Failure() {
	request := &pb.FindAllTransferRequest{Search: "", Page: 1, PageSize: 10}
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().FindByTrashedTransfer(gomock.Any(), request).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to retrieve transfer data: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashedTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Failed to fetch trashed transfers")
}

func (suite *TransferHandlerTestSuite) TestFindByTrashedTransfer_Empty() {
	reqPB := &pb.FindAllTransferRequest{
		Search:   "",
		Page:     1,
		PageSize: 10,
	}

	grpcResp := &pb.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "No trashed transfers found",
		Data:    []*pb.TransferResponseDeleteAt{},
	}

	apiResp := &response.ApiResponsePaginationTransferDeleteAt{
		Status:  "success",
		Message: "No trashed transfers found",
		Data:    []*response.TransferResponseDeleteAt{},
	}

	suite.MockTransferClient.EXPECT().
		FindByTrashedTransfer(gomock.Any(), reqPB).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponsePaginationTransferDeleteAt(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/api/transfer/trashed", nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err := suite.Handler.FindByTrashedTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponsePaginationTransferDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Empty(resp.Data)
}

func (suite *TransferHandlerTestSuite) TestCreateTransfer_Success() {
	body := requests.CreateTransferRequest{
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	grpcReq := &pb.CreateTransferRequest{
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	}

	grpcResp := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer created successfully",
		Data: &pb.TransferResponse{
			Id:             1,
			TransferFrom:   body.TransferFrom,
			TransferTo:     body.TransferTo,
			TransferAmount: int32(body.TransferAmount),
		},
	}

	apiResp := &response.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer created successfully",
		Data: &response.TransferResponse{
			ID:             1,
			TransferFrom:   body.TransferFrom,
			TransferTo:     body.TransferTo,
			TransferAmount: body.TransferAmount,
		},
	}

	suite.MockTransferClient.EXPECT().
		CreateTransfer(gomock.Any(), grpcReq).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransfer(grpcResp).
		Return(apiResp).
		Times(1)

	requestBodyBytes, err := json.Marshal(body)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err = suite.Handler.CreateTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Equal(int(1), resp.Data.ID)
}

func (suite *TransferHandlerTestSuite) TestCreateTransfer_Failure() {
	body := requests.CreateTransferRequest{TransferFrom: "test_from", TransferTo: "test_to", TransferAmount: 500000}
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().CreateTransfer(gomock.Any(), gomock.Any()).Return(nil, grpcError)
	suite.MockLogger.EXPECT().Debug("Failed to create transfer: ", gomock.Any()).Times(1)

	requestBodyBytes, err := json.Marshal(body)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err = suite.Handler.CreateTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Failed to create transfer")
}

func (suite *TransferHandlerTestSuite) TestCreateTransfer_ValidationError() {
	invalidBody := requests.CreateTransferRequest{TransferFrom: "test_from", TransferTo: "test_to"}

	suite.MockTransferClient.EXPECT().CreateTransfer(gomock.Any(), gomock.Any()).Times(0)
	suite.MockLogger.EXPECT().Debug("Validation Error: ", gomock.Any()).Times(1)

	requestBodyBytes, err := json.Marshal(invalidBody)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/api/transfer", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	err = suite.Handler.CreateTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "validation failed: invalid create transfer request")
}

func (suite *TransferHandlerTestSuite) TestUpdateTransfer_Success() {
	id := 1

	body := requests.UpdateTransferRequest{
		TransferID:     &id,
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	expectedGRPCRequest := &pb.UpdateTransferRequest{
		TransferId:     int32(id),
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	}

	expectedGRPCResponse := &pb.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer updated successfully",
		Data: &pb.TransferResponse{
			Id:             int32(id),
			TransferFrom:   body.TransferFrom,
			TransferTo:     body.TransferTo,
			TransferAmount: int32(body.TransferAmount),
		},
	}

	expectedAPIResponse := &response.ApiResponseTransfer{
		Status:  "success",
		Message: "Transfer updated successfully",
		Data: &response.TransferResponse{
			ID:             id,
			TransferFrom:   body.TransferFrom,
			TransferTo:     body.TransferTo,
			TransferAmount: body.TransferAmount,
		},
	}

	suite.MockTransferClient.EXPECT().
		UpdateTransfer(gomock.Any(), expectedGRPCRequest).
		Return(expectedGRPCResponse, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransfer(expectedGRPCResponse).
		Return(expectedAPIResponse).
		Times(1)

	requestBodyBytes, err := json.Marshal(body)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/api/transfer/%d", id),
		bytes.NewBuffer(requestBodyBytes),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err = suite.Handler.UpdateTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransfer
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
}

func (suite *TransferHandlerTestSuite) TestUpdateTransfer_Failure() {
	id := 1
	body := requests.UpdateTransferRequest{
		TransferID:     &id,
		TransferFrom:   "test_from",
		TransferTo:     "test_to",
		TransferAmount: 500000,
	}

	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().
		UpdateTransfer(gomock.Any(), gomock.Any()).
		Return(nil, grpcError).
		Times(1)

	suite.MockLogger.EXPECT().
		Debug("Failed to update transfer: ", gomock.Any()).
		Times(1)

	requestBodyBytes, err := json.Marshal(body)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/api/transfer/%d", id),
		bytes.NewBuffer(requestBodyBytes),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err = suite.Handler.UpdateTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TransferHandlerTestSuite) TestUpdateTransfer_ValidationError() {
	id := 1
	invalidBody := requests.UpdateTransferRequest{
		TransferID:   &id,
		TransferFrom: "test_from",
		TransferTo:   "test_to",
	}

	suite.MockTransferClient.EXPECT().
		UpdateTransfer(gomock.Any(), gomock.Any()).
		Times(0)

	suite.MockLogger.EXPECT().
		Debug("Validation Error: ", gomock.Any()).
		Times(1)

	requestBodyBytes, err := json.Marshal(invalidBody)
	suite.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/api/transfer/%d", id),
		bytes.NewBuffer(requestBodyBytes),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err = suite.Handler.UpdateTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TransferHandlerTestSuite) TestTrashTransfer_Success() {
	id := 1

	grpcResp := &pb.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Transfer trashed successfully",
		Data:    &pb.TransferResponseDeleteAt{Id: int32(id)},
	}

	apiResp := &response.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Transfer trashed successfully",
		Data:    &response.TransferResponseDeleteAt{ID: id},
	}

	suite.MockTransferClient.EXPECT().
		TrashedTransfer(
			gomock.Any(),
			&pb.FindByIdTransferRequest{TransferId: int32(id)},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransferDeleteAt(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/api/transfer/%d/trash", id),
		nil,
	)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.TrashTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransferDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
}

func (suite *TransferHandlerTestSuite) TestTrashTransfer_Failure() {
	id := 1
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().
		TrashedTransfer(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(nil, grpcError)

	suite.MockLogger.EXPECT().Debug("Failed to trash transfer: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%d/trash", id), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	err := suite.Handler.TrashTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Failed to trash transfer")
}

func (suite *TransferHandlerTestSuite) TestTrashTransfer_InvalidID() {
	invalidID := "invalid_id"

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%s/trash", invalidID), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	suite.MockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	err := suite.Handler.TrashTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Invalid Transfer ID")
}

func (suite *TransferHandlerTestSuite) TestRestoreTransfer_Success() {
	id := 1

	grpcResp := &pb.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Transfer restored successfully",
		Data:    &pb.TransferResponseDeleteAt{Id: int32(id)},
	}

	apiResp := &response.ApiResponseTransferDeleteAt{
		Status:  "success",
		Message: "Transfer restored successfully",
		Data:    &response.TransferResponseDeleteAt{ID: id},
	}

	suite.MockTransferClient.EXPECT().
		RestoreTransfer(
			gomock.Any(),
			&pb.FindByIdTransferRequest{TransferId: int32(id)},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransferDeleteAt(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/api/transfer/%d/restore", id),
		nil,
	)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	err := suite.Handler.RestoreTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransferDeleteAt
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Equal("success", resp.Status)
}

func (suite *TransferHandlerTestSuite) TestRestoreTransfer_Failure() {
	id := 1
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().
		RestoreTransfer(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(nil, grpcError)

	suite.MockLogger.EXPECT().Debug("Failed to restore transfer: ", gomock.Any()).Times(1)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%d/restore", id), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	err := suite.Handler.RestoreTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Failed to restore transfer")
}

func (suite *TransferHandlerTestSuite) TestRestoreTransfer_InvalidID() {
	invalidID := "invalid_id"

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/transfer/%s/restore", invalidID), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	suite.MockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	err := suite.Handler.RestoreTransfer(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Invalid Transfer ID")
}

func (suite *TransferHandlerTestSuite) TestDeleteTransferPermanent_Success() {
	id := 1

	grpcResp := &pb.ApiResponseTransferDelete{
		Status:  "success",
		Message: "Transfer deleted permanently",
	}

	apiResp := &response.ApiResponseTransferDelete{
		Status:  "success",
		Message: "Transfer deleted permanently",
	}

	suite.MockTransferClient.EXPECT().
		DeleteTransferPermanent(
			gomock.Any(),
			&pb.FindByIdTransferRequest{TransferId: int32(id)},
		).
		Return(grpcResp, nil).
		Times(1)

	suite.Mapper.EXPECT().
		ToApiResponseTransferDelete(grpcResp).
		Return(apiResp).
		Times(1)

	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/api/transfer/%d/delete", id),
		nil,
	)

	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	err := suite.Handler.DeleteTransferPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	var resp response.ApiResponseTransferDelete
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)

	suite.Equal("success", resp.Status)
	suite.Equal("Transfer deleted permanently", resp.Message)
}

func (suite *TransferHandlerTestSuite) TestDeleteTransferPermanent_Failure() {
	id := 1
	grpcError := errors.New("internal server error")

	suite.MockTransferClient.EXPECT().
		DeleteTransferPermanent(gomock.Any(), &pb.FindByIdTransferRequest{TransferId: int32(id)}).
		Return(nil, grpcError)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transfer/%d/delete", id), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))

	err := suite.Handler.DeleteTransferPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Failed to permanently delete transfer")
}

func (suite *TransferHandlerTestSuite) TestDeleteTransferPermanent_InvalidID() {
	invalidID := "invalid_id"

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/transfer/%s/delete", invalidID), nil)
	rec := httptest.NewRecorder()
	c := suite.E.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(invalidID)

	suite.MockLogger.EXPECT().Debug("Bad Request: Invalid ID", gomock.Any()).Times(1)

	err := suite.Handler.DeleteTransferPermanent(c)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, rec.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.NoError(err)
	suite.Contains(resp.Message, "Invalid Transfer ID")
}

func TestTransferHandlerSuite(t *testing.T) {
	suite.Run(t, new(TransferHandlerTestSuite))
}
