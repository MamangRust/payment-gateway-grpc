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

type testSuiteCard struct {
	E              *echo.Echo
	Rec            *httptest.ResponseRecorder
	MockCardClient *mock_pb.MockCardServiceClient
	MockLogger     *mock_logger.MockLoggerInterface
	MockMapper     *mock_apimapper.MockCardResponseMapper
	Handler        *api.CardHandleApi
	Ctrl           *gomock.Controller
}

func SetupTestCard(t *testing.T) *testSuiteCard {
	ctrl := gomock.NewController(t)
	mockCardClient := mock_pb.NewMockCardServiceClient(ctrl)
	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockMapper := mock_apimapper.NewMockCardResponseMapper(ctrl)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("userID", 1)
			return next(c)
		}
	})

	handler := api.NewHandlerCard(mockCardClient, e, mockLogger, mockMapper)

	return &testSuiteCard{
		E:              e,
		Rec:            httptest.NewRecorder(),
		MockCardClient: mockCardClient,
		MockLogger:     mockLogger,
		MockMapper:     mockMapper,
		Handler:        handler,
		Ctrl:           ctrl,
	}
}
