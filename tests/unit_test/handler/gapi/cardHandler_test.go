package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	mock_protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto/mocks"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	mock_service "MamangRust/paymentgatewaygrpc/internal/service/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardHandleGrpcTestSuite struct {
	suite.Suite
	Ctrl            *gomock.Controller
	MockCardService *mock_service.MockCardService
	MockProtoMapper *mock_protomapper.MockCardProtoMapper
	Handler         gapi.CardHandleGrpc
}

func (suite *CardHandleGrpcTestSuite) SetupTest() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockCardService = mock_service.NewMockCardService(suite.Ctrl)
	suite.MockProtoMapper = mock_protomapper.NewMockCardProtoMapper(suite.Ctrl)
	suite.Handler = gapi.NewCardHandleGrpc(suite.MockCardService, suite.MockProtoMapper)
}

func (suite *CardHandleGrpcTestSuite) TearDownTest() {
	suite.Ctrl.Finish()
}

func (suite *CardHandleGrpcTestSuite) TestFindAllCard_Success() {
	req := &pb.FindAllCardRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockCards := []*response.CardResponse{
		{ID: 1, CardNumber: "1234567890123456", CardType: "Credit", UserID: 1},
		{ID: 2, CardNumber: "6543210987654321", CardType: "Debit", UserID: 2},
	}
	mockProtoCards := []*pb.CardResponse{
		{Id: 1, CardNumber: "1234567890123456", CardType: "Credit", UserId: 1},
		{Id: 2, CardNumber: "6543210987654321", CardType: "Debit", UserId: 2},
	}

	totalRecords := 2
	suite.MockCardService.EXPECT().
		FindAll(gomock.Eq(&requests.FindAllCards{Page: 1, PageSize: 10, Search: "test"})).
		Return(mockCards, &totalRecords, nil)

	suite.MockProtoMapper.EXPECT().
		ToProtoResponsePaginationCard(gomock.Any(), "success", "Successfully fetched card records", mockCards).
		Return(&pb.ApiResponsePaginationCard{
			Status:  "success",
			Message: "Successfully fetched card records",
			Data:    mockProtoCards,
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   1,
				TotalRecords: 2,
			},
		})

	res, err := suite.Handler.FindAllCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("Successfully fetched card records", res.GetMessage())
	suite.Equal(int32(2), res.GetPagination().GetTotalRecords())
	suite.Equal(2, len(res.GetData()))
}

func (suite *CardHandleGrpcTestSuite) TestFindAllCard_Failure() {
	req := &pb.FindAllCardRequest{Page: 1, PageSize: 10, Search: "test"}
	serviceError := &response.ErrorResponse{Status: "error", Message: "Failed to fetch cards"}

	totalRecords := 0
	suite.MockCardService.EXPECT().FindAll(gomock.Any()).Return(nil, &totalRecords, serviceError)

	res, _ := suite.Handler.FindAllCard(context.Background(), req)

	suite.Nil(res)
}

func (suite *CardHandleGrpcTestSuite) TestFindByIdCard_Success() {
	req := &pb.FindByIdCardRequest{CardId: 1}
	mockCard := &response.CardResponse{ID: 1, CardNumber: "1234567890123456", CardType: "Credit", UserID: 1}
	mockProtoCard := &pb.CardResponse{Id: 1, CardNumber: "1234567890123456", CardType: "Credit", UserId: 1}

	suite.MockCardService.EXPECT().FindById(1).Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCard("success", "Successfully fetched card record", mockCard).Return(&pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    mockProtoCard,
	})

	res, err := suite.Handler.FindByIdCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("1234567890123456", res.GetData().GetCardNumber())
}

func (suite *CardHandleGrpcTestSuite) TestFindByIdCard_InvalidId() {
	req := &pb.FindByIdCardRequest{CardId: 0}

	res, err := suite.Handler.FindByIdCard(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid card ID")
}

func (suite *CardHandleGrpcTestSuite) TestFindByUserIdCard_Success() {
	req := &pb.FindByUserIdCardRequest{UserId: 1}
	mockCard := &response.CardResponse{ID: 1, CardNumber: "1234567890123456", CardType: "Credit", UserID: 1}
	mockProtoCard := &pb.CardResponse{Id: 1, CardNumber: "1234567890123456", CardType: "Credit", UserId: 1}

	suite.MockCardService.EXPECT().FindByUserID(1).Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCard("success", "Successfully fetched card record", mockCard).Return(&pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    mockProtoCard,
	})

	res, err := suite.Handler.FindByUserIdCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("1234567890123456", res.GetData().GetCardNumber())
}

func (suite *CardHandleGrpcTestSuite) TestFindByCardNumber_Success() {
	req := &pb.FindByCardNumberRequest{CardNumber: "1234567890123456"}
	mockCard := &response.CardResponse{ID: 1, CardNumber: "1234567890123456", CardType: "Credit", UserID: 1}
	mockProtoCard := &pb.CardResponse{Id: 1, CardNumber: "1234567890123456", CardType: "Credit", UserId: 1}

	suite.MockCardService.EXPECT().FindByCardNumber("1234567890123456").Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCard("success", "Successfully fetched card record", mockCard).Return(&pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    mockProtoCard,
	})

	res, err := suite.Handler.FindByCardNumber(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("1234567890123456", res.GetData().GetCardNumber())
}

func (suite *CardHandleGrpcTestSuite) TestFindByActiveCard_Success() {
	req := &pb.FindAllCardRequest{Page: 1, PageSize: 10, Search: ""}
	activeCards := []*response.CardResponseDeleteAt{
		{ID: 1, CardNumber: "1234567890123456", CardType: "Credit", UserID: 1},
	}
	totalRecords := 1
	suite.MockCardService.EXPECT().FindByActive(gomock.Any()).Return(activeCards, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationCardDeletedAt(gomock.Any(), "success", "Successfully fetched card record", activeCards).Return(&pb.ApiResponsePaginationCardDeleteAt{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data: []*pb.CardResponseDeleteAt{
			{Id: 1, CardNumber: "1234567890123456", CardType: "Credit", UserId: 1},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByActiveCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *CardHandleGrpcTestSuite) TestFindByTrashedCard_Success() {
	req := &pb.FindAllCardRequest{Page: 1, PageSize: 10, Search: ""}
	trashedCards := []*response.CardResponseDeleteAt{
		{ID: 1, CardNumber: "1234567890123456", CardType: "Credit", UserID: 1},
	}
	totalRecords := 1
	suite.MockCardService.EXPECT().FindByTrashed(gomock.Any()).Return(trashedCards, &totalRecords, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponsePaginationCardDeletedAt(gomock.Any(), "success", "Successfully fetched card record", trashedCards).Return(&pb.ApiResponsePaginationCardDeleteAt{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data: []*pb.CardResponseDeleteAt{
			{Id: 1, CardNumber: "1234567890123456", CardType: "Credit", UserId: 1},
		},
		Pagination: &pb.PaginationMeta{CurrentPage: 1, PageSize: 10, TotalPages: 1, TotalRecords: 1},
	})

	res, err := suite.Handler.FindByTrashedCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Len(res.GetData(), 1)
}

func (suite *CardHandleGrpcTestSuite) TestCreateCard_Success() {
	expireDate := time.Now().AddDate(2, 0, 0)
	req := &pb.CreateCardRequest{
		UserId:       1,
		CardType:     "credit",
		ExpireDate:   timestamppb.New(expireDate),
		Cvv:          "123",
		CardProvider: "alfamart",
	}

	mockCard := &response.CardResponse{
		ID:           1,
		CardNumber:   "1234567890123456",
		CardType:     "credit",
		UserID:       1,
		ExpireDate:   expireDate.Format(time.RFC3339),
		CVV:          "123",
		CardProvider: "alfamart",
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	mockProtoCard := &pb.CardResponse{
		Id:           1,
		CardNumber:   "1234567890123456",
		CardType:     "credit",
		UserId:       1,
		ExpireDate:   expireDate.Format(time.RFC3339),
		Cvv:          "123",
		CardProvider: "alfamart",
	}

	suite.MockCardService.EXPECT().CreateCard(gomock.Any()).Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCard("success", "Successfully created card", mockCard).Return(&pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully created card",
		Data:    mockProtoCard,
	})

	res, err := suite.Handler.CreateCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("credit", res.GetData().GetCardType())
}

func (suite *CardHandleGrpcTestSuite) TestCreateCard_ValidationError() {
	req := &pb.CreateCardRequest{
		UserId:       0,
		CardType:     "",
		ExpireDate:   nil,
		Cvv:          "",
		CardProvider: "",
	}

	res, err := suite.Handler.CreateCard(context.Background(), req)

	suite.Nil(res)
	suite.Error(err)
	statusErr, ok := status.FromError(err)
	suite.True(ok)
	suite.Equal(codes.InvalidArgument, statusErr.Code())
	suite.Contains(statusErr.Message(), "Invalid input for create card")
}

func (suite *CardHandleGrpcTestSuite) TestUpdateCard_Success() {
	cardId := 1
	expireDate := time.Now().AddDate(2, 0, 0)
	req := &pb.UpdateCardRequest{
		CardId:       int32(cardId),
		UserId:       1,
		CardType:     "debit",
		ExpireDate:   timestamppb.New(expireDate),
		Cvv:          "456",
		CardProvider: "alfamart",
	}

	mockCard := &response.CardResponse{
		ID:           1,
		CardNumber:   "1234567890123456",
		CardType:     "debit",
		UserID:       1,
		ExpireDate:   expireDate.Format(time.RFC3339),
		CVV:          "456",
		CardProvider: "alfamart",
	}

	mockProtoCard := &pb.CardResponse{
		Id:           1,
		CardNumber:   "1234567890123456",
		CardType:     "debit",
		UserId:       1,
		ExpireDate:   expireDate.Format(time.RFC3339),
		Cvv:          "456",
		CardProvider: "alfamart",
	}

	suite.MockCardService.EXPECT().UpdateCard(gomock.Any()).Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCard("success", "Successfully updated card", mockCard).Return(&pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully updated card",
		Data:    mockProtoCard,
	})

	res, err := suite.Handler.UpdateCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
	suite.Equal("debit", res.GetData().GetCardType())
}

func (suite *CardHandleGrpcTestSuite) TestTrashedCard_Success() {
	req := &pb.FindByIdCardRequest{CardId: 1}
	mockCard := &response.CardResponseDeleteAt{ID: 1, CardNumber: "1234567890123456"}

	suite.MockCardService.EXPECT().TrashedCard(1).Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCardDeleteAt("success", "Successfully trashed card", mockCard).Return(&pb.ApiResponseCardDeleteAt{
		Status:  "success",
		Message: "Successfully trashed card",
		Data:    &pb.CardResponseDeleteAt{Id: 1, CardNumber: "1234567890123456"},
	})

	res, err := suite.Handler.TrashedCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *CardHandleGrpcTestSuite) TestRestoreCard_Success() {
	req := &pb.FindByIdCardRequest{CardId: 1}
	mockCard := &response.CardResponseDeleteAt{ID: 1, CardNumber: "1234567890123456"}

	suite.MockCardService.EXPECT().RestoreCard(1).Return(mockCard, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCardDeleteAt("success", "Successfully restored card", mockCard).Return(&pb.ApiResponseCardDeleteAt{
		Status:  "success",
		Message: "Successfully restored card",
		Data:    &pb.CardResponseDeleteAt{Id: 1, CardNumber: "1234567890123456"},
	})

	res, err := suite.Handler.RestoreCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *CardHandleGrpcTestSuite) TestDeleteCardPermanent_Success() {
	req := &pb.FindByIdCardRequest{CardId: 1}

	suite.MockCardService.EXPECT().DeleteCardPermanent(1).Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCardDelete("success", "Successfully deleted card").Return(&pb.ApiResponseCardDelete{
		Status:  "success",
		Message: "Successfully deleted card",
	})

	res, err := suite.Handler.DeleteCardPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *CardHandleGrpcTestSuite) TestRestoreAllCard_Success() {
	req := &emptypb.Empty{}

	suite.MockCardService.EXPECT().RestoreAllCard().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCardAll("success", "Successfully restore card").Return(&pb.ApiResponseCardAll{
		Status:  "success",
		Message: "Successfully restore card",
	})

	res, err := suite.Handler.RestoreAllCard(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func (suite *CardHandleGrpcTestSuite) TestDeleteAllCardPermanent_Success() {
	req := &emptypb.Empty{}

	suite.MockCardService.EXPECT().DeleteAllCardPermanent().Return(true, nil)
	suite.MockProtoMapper.EXPECT().ToProtoResponseCardAll("success", "Successfully delete card permanent").Return(&pb.ApiResponseCardAll{
		Status:  "success",
		Message: "Successfully delete card permanent",
	})

	res, err := suite.Handler.DeleteAllCardPermanent(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal("success", res.GetStatus())
}

func TestCardHandleGrpcSuite(t *testing.T) {
	suite.Run(t, new(CardHandleGrpcTestSuite))
}
