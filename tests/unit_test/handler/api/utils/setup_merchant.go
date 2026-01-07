package utils_test

import (
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	mock_apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api/mocks"
	mock_pb "MamangRust/paymentgatewaygrpc/internal/pb/mocks"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

type merchantTestSuite struct {
	Ctrl               *gomock.Controller
	MockMerchantClient *mock_pb.MockMerchantServiceClient
	MockLogger         *mock_logger.MockLoggerInterface
	E                  *echo.Echo
	Rec                *httptest.ResponseRecorder
	Handler            *api.MerchantHandleApi
}

func SetupTestMerchant(t *testing.T) *merchantTestSuite {
	ctrl := gomock.NewController(t)
	mockMerchantClient := mock_pb.NewMockMerchantServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockMerchantResponseMapper(ctrl)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("userID", 1)
			return next(c)
		}
	})

	handler := api.NewHandlerMerchant(mockMerchantClient, e, mockLogger, mockMapper)

	return &merchantTestSuite{
		Ctrl:               ctrl,
		MockMerchantClient: mockMerchantClient,
		MockLogger:         mockLogger,
		E:                  e,
		Rec:                httptest.NewRecorder(),
		Handler:            handler,
	}
}
