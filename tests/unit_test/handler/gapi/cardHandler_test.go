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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindAllCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindAllCardRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockCards := []*response.CardResponse{
		{ID: 1, CardNumber: "1222"},
		{ID: 2, CardNumber: "1222"},
	}

	mockCardService.EXPECT().FindAll(1, 10, "test").Return(mockCards, 2, nil)
	mockProtoMapper.EXPECT().ToResponsesCard(mockCards).Return([]*pb.CardResponse{
		{Id: 1, CardNumber: "1222"},
		{Id: 2, CardNumber: "1222"},
	})

	res, err := cardHandler.FindAllCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card records", res.GetMessage())
	assert.Equal(t, int32(1), res.GetPagination().GetCurrentPage())
	assert.Equal(t, int32(10), res.GetPagination().GetPageSize())
	assert.Equal(t, int32(2), res.GetPagination().GetTotalRecords())
	assert.Equal(t, 2, len(res.GetData()))
}

func TestFindAllCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindAllCardRequest{
		Page:     1,
		PageSize: 10,
		Search:   "test",
	}

	mockCardService.EXPECT().FindAll(1, 10, "test").Return(nil, 0, &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to fetch card records",
	})

	res, err := cardHandler.FindAllCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to fetch card records")
}

func TestFindAllCard_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindAllCardRequest{
		Page:     1,
		PageSize: 10,
		Search:   "empty",
	}

	mockCardService.EXPECT().FindAll(1, 10, "empty").Return([]*response.CardResponse{}, 0, nil)
	mockProtoMapper.EXPECT().ToResponsesCard([]*response.CardResponse{}).Return([]*pb.CardResponse{})

	res, err := cardHandler.FindAllCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card records", res.GetMessage())
	assert.Equal(t, int32(0), res.GetPagination().GetTotalRecords())
	assert.Equal(t, 0, len(res.GetData()))
}

func TestFindByIdCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByIdCardRequest{CardId: 1}

	mockCard := &response.CardResponse{ID: 1, CardNumber: "Card 1"}
	mockCardService.EXPECT().FindById(1).Return(mockCard, nil)
	mockProtoMapper.EXPECT().ToResponseCard(mockCard).Return(&pb.CardResponse{
		Id:         1,
		CardNumber: "Card 1",
	})

	res, err := cardHandler.FindByIdCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.Equal(t, "Card 1", res.GetData().GetCardNumber())
}

func TestFindByIdCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByIdCardRequest{CardId: 1}

	mockCardService.EXPECT().FindById(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Card not found",
	})

	res, err := cardHandler.FindByIdCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Card not found")
}

func TestFindByIdCard_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByIdCardRequest{CardId: 0}

	res, err := cardHandler.FindByIdCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid card id")
}

func TestFindByUserIdCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByUserIdCardRequest{UserId: 1}

	mockCard := &response.CardResponse{ID: 1, CardNumber: "Card 1"}
	mockCardService.EXPECT().FindByUserID(1).Return(mockCard, nil)
	mockProtoMapper.EXPECT().ToResponseCard(mockCard).Return(&pb.CardResponse{
		Id:         1,
		CardNumber: "Card 1",
	})

	res, err := cardHandler.FindByUserIdCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.Equal(t, "Card 1", res.GetData().GetCardNumber())
}

func TestFindByUserIdCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByUserIdCardRequest{UserId: 1}

	mockCardService.EXPECT().FindByUserID(1).Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Card not found",
	})

	res, err := cardHandler.FindByUserIdCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Card not found")
}

func TestFindByUserIdCard_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByUserIdCardRequest{UserId: 0}

	res, err := cardHandler.FindByUserIdCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Invalid user id")
}

func TestFindByActiveCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	activeCards := []*response.CardResponseDeleteAt{
		{ID: 1, CardNumber: "Active Card 1"},
		{ID: 2, CardNumber: "Active Card 2"},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockCardService.EXPECT().FindByActive(pageSize, page, search).Return(activeCards, expected, nil)
	mockProtoMapper.EXPECT().ToResponsesCardDeletedAt(activeCards).Return([]*pb.CardResponseDeleteAt{
		{Id: 1, CardNumber: "Active Card 1"},
		{Id: 2, CardNumber: "Active Card 2"},
	})

	res, err := cardHandler.FindByActiveCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
	assert.Equal(t, "Active Card 1", res.GetData()[0].GetCardNumber())
	assert.Equal(t, "Active Card 2", res.GetData()[1].GetCardNumber())
}

func TestFindByActiveCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockCardService.EXPECT().FindByActive(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Card not found",
	})

	res, err := cardHandler.FindByActiveCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Card not found")
}

func TestFindByActiveCard_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockCardService.EXPECT().FindByActive(pageSize, page, search).Return([]*response.CardResponseDeleteAt{}, expected, nil)
	mockProtoMapper.EXPECT().ToResponsesCardDeletedAt([]*response.CardResponseDeleteAt{}).Return([]*pb.CardResponseDeleteAt{})

	res, err := cardHandler.FindByActiveCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByTrashedCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	trashedCards := []*response.CardResponseDeleteAt{
		{ID: 1, CardNumber: "Trashed Card 1"},
		{ID: 2, CardNumber: "Trashed Card 2"},
	}

	page := 1
	pageSize := 1
	search := ""
	expected := 2

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockCardService.EXPECT().FindByTrashed(pageSize, page, search).Return(trashedCards, expected, nil)
	mockProtoMapper.EXPECT().ToResponsesCardDeletedAt(trashedCards).Return([]*pb.CardResponseDeleteAt{
		{Id: 1, CardNumber: "Trashed Card 1"},
		{Id: 2, CardNumber: "Trashed Card 2"},
	})

	res, err := cardHandler.FindByTrashedCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.Len(t, res.GetData(), 2)
	assert.Equal(t, "Trashed Card 1", res.GetData()[0].GetCardNumber())
	assert.Equal(t, "Trashed Card 2", res.GetData()[1].GetCardNumber())
}

func TestFindByTrashedCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	page := 1
	pageSize := 1
	search := ""
	expected := 0

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockCardService.EXPECT().FindByTrashed(pageSize, page, search).Return(nil, expected, &response.ErrorResponse{
		Status:  "error",
		Message: "Card not found",
	})

	res, err := cardHandler.FindByTrashedCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Card not found")
}

func TestFindByTrashedCard_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	page := 1
	pageSize := 1
	search := ""

	expected := 0

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	mockCardService.EXPECT().FindByTrashed(pageSize, page, search).Return([]*response.CardResponseDeleteAt{}, expected, nil)
	mockProtoMapper.EXPECT().ToResponsesCardDeletedAt([]*response.CardResponseDeleteAt{}).Return([]*pb.CardResponseDeleteAt{})

	res, err := cardHandler.FindByTrashedCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.Len(t, res.GetData(), 0)
}

func TestFindByCardNumber_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	card := response.CardResponse{
		ID:         1,
		CardNumber: "1234567890123456",
	}

	req := &pb.FindByCardNumberRequest{CardNumber: "1234567890123456"}

	mockCardService.EXPECT().FindByCardNumber("1234567890123456").Return(&card, nil)
	mockProtoMapper.EXPECT().ToResponseCard(&card).Return(&pb.CardResponse{
		Id:         1,
		CardNumber: "1234567890123456",
	})

	res, err := cardHandler.FindByCardNumber(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.GetStatus())
	assert.Equal(t, "Successfully fetched card record", res.GetMessage())
	assert.NotNil(t, res.GetData())
	assert.Equal(t, "1234567890123456", res.GetData().GetCardNumber())
}

func TestFindByCardNumber_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByCardNumberRequest{CardNumber: "1234567890123456"}

	mockCardService.EXPECT().FindByCardNumber("1234567890123456").Return(nil, &response.ErrorResponse{
		Status:  "error",
		Message: "Card not found",
	})

	res, err := cardHandler.FindByCardNumber(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Card not found")
}

func TestCreateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	now := time.Now().UTC()
	expireDate := now.AddDate(5, 0, 0)
	expireDateProto := timestamppb.New(expireDate)

	req := &pb.CreateCardRequest{
		UserId:       1,
		CardType:     "credit",
		ExpireDate:   expireDateProto,
		Cvv:          "123",
		CardProvider: "mandiri",
	}

	mockCreateCardRequest := &requests.CreateCardRequest{
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockCard := &response.CardResponse{
		ID:           1,
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate.String(),
		CVV:          "123",
		CardProvider: "mandiri",
		CreatedAt:    now.String(),
		UpdatedAt:    now.String(),
	}

	mockCardService.EXPECT().
		CreateCard(mockCreateCardRequest).
		Return(mockCard, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseCard(mockCard).
		Return(&pb.CardResponse{
			Id:           1,
			UserId:       1,
			CardType:     "credit",
			ExpireDate:   expireDateProto.String(),
			Cvv:          "123",
			CardProvider: "mandiri",
			CreatedAt:    timestamppb.New(now).String(),
			UpdatedAt:    timestamppb.New(now).String(),
		})

	response, err := cardHandler.CreateCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully created card", response.GetMessage())

}

func TestCreateCard_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	reqMock := requests.CreateCardRequest{
		UserID:       1,
		CardType:     "s",
		ExpireDate:   time.Now(),
		CVV:          "123",
		CardProvider: "mandiri",
	}

	req := &pb.CreateCardRequest{
		UserId:       1,
		CardType:     "s",
		ExpireDate:   timestamppb.Now(),
		Cvv:          "",
		CardProvider: "mandiri",
	}

	if err := reqMock.Validate(); err != nil {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "card type must be credit or debit")
	}

	response, err := cardHandler.CreateCard(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
}

func TestCreateCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	now := time.Now().UTC()
	expireDate := now.AddDate(5, 0, 0)
	expireDateProto := timestamppb.New(expireDate)

	req := &pb.CreateCardRequest{
		UserId:       1,
		CardType:     "credit",
		ExpireDate:   expireDateProto,
		Cvv:          "123",
		CardProvider: "mandiri",
	}

	mockCreateCardRequest := &requests.CreateCardRequest{
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to create card",
	}
	mockCardService.EXPECT().
		CreateCard(mockCreateCardRequest).
		Return(nil, mockError).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseCard(gomock.Any()).
		Times(0)

	response, err := cardHandler.CreateCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())

}

func TestUpdateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	now := time.Now().UTC()
	expireDate := now.AddDate(5, 0, 0)
	expireDateProto := timestamppb.New(expireDate)

	req := &pb.UpdateCardRequest{
		CardId:       1,
		UserId:       1,
		CardType:     "credit",
		ExpireDate:   expireDateProto,
		Cvv:          "123",
		CardProvider: "mandiri",
	}

	mockUpdateCardRequest := &requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockCard := &response.CardResponse{
		ID:           1,
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate.String(),
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockCardService.EXPECT().
		UpdateCard(mockUpdateCardRequest).
		Return(mockCard, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseCard(mockCard).
		Times(1)

	response, err := cardHandler.UpdateCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully updated card", response.GetMessage())

}

func TestUpdateCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	now := time.Now().UTC()
	expireDate := now.AddDate(5, 0, 0)
	expireDateProto := timestamppb.New(expireDate)

	req := &pb.UpdateCardRequest{
		CardId:       1,
		UserId:       1,
		CardType:     "credit",
		ExpireDate:   expireDateProto,
		Cvv:          "123",
		CardProvider: "mandiri",
	}

	mockUpdateCardRequest := &requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to update card",
	}
	mockCardService.EXPECT().
		UpdateCard(mockUpdateCardRequest).
		Return(nil, mockError).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseCard(gomock.Any()).
		Times(0)

	response, err := cardHandler.UpdateCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
}

func TestUpdateCard_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	now := time.Now().UTC()
	expireDate := now.AddDate(5, 0, 0)
	expireDateProto := timestamppb.New(expireDate)

	req := &pb.UpdateCardRequest{
		CardId:       1,
		UserId:       1,
		CardType:     "s",
		ExpireDate:   expireDateProto,
		Cvv:          "123",
		CardProvider: "mandiri",
	}

	mockUpdateCardRequest := &requests.UpdateCardRequest{
		CardID:       1,
		UserID:       1,
		CardType:     "s",
		ExpireDate:   expireDate,
		CVV:          "123",
		CardProvider: "mandiri",
	}

	if err := mockUpdateCardRequest.Validate(); err != nil {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "card type must be credit or debit")
	}

	response, err := cardHandler.UpdateCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)
}

func TestTrashedCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByIdCardRequest{
		CardId: 1,
	}

	mockCardResponse := &response.CardResponse{
		ID:           1,
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   "2029-12-26",
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockCardService.EXPECT().
		TrashedCard(1).
		Return(mockCardResponse, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseCard(mockCardResponse).
		Return(&pb.CardResponse{
			Id:           1,
			UserId:       1,
			CardType:     "credit",
			ExpireDate:   "2029-12-26",
			Cvv:          "123",
			CardProvider: "mandiri",
		}).Times(1)

	response, err := cardHandler.TrashedCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully trashed card", response.GetMessage())
	assert.NotNil(t, response.GetData())
}

func TestTrashedCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 1,
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to trashed card",
	}

	mockCardService.EXPECT().
		TrashedCard(1).
		Return(nil, mockError).
		Times(1)

	response, err := cardHandler.TrashedCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

}

func TestTrashedCard_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 0,
	}

	mockCardService.EXPECT().
		TrashedCard(0).
		Return(nil, &response.ErrorResponse{Status: "error", Message: "Invalid card id"}).
		Times(1)

	response, err := cardHandler.TrashedCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Invalid Id")
}

func TestRestoreCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	mockProtoMapper := mock_protomapper.NewMockCardProtoMapper(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, mockProtoMapper)

	req := &pb.FindByIdCardRequest{
		CardId: 1,
	}

	mockCardResponse := &response.CardResponse{
		ID:           1,
		UserID:       1,
		CardType:     "credit",
		ExpireDate:   "2029-12-26",
		CVV:          "123",
		CardProvider: "mandiri",
	}

	mockCardService.EXPECT().
		RestoreCard(1).
		Return(mockCardResponse, nil).
		Times(1)

	mockProtoMapper.EXPECT().
		ToResponseCard(mockCardResponse).
		Return(&pb.CardResponse{
			Id:           1,
			UserId:       1,
			CardType:     "credit",
			ExpireDate:   "2029-12-26",
			Cvv:          "123",
			CardProvider: "mandiri",
		}).Times(1)

	response, err := cardHandler.RestoreCard(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully restored card", response.GetMessage())
	assert.NotNil(t, response.GetData())
}

func TestRestoreCard_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 1,
	}

	mockError := &response.ErrorResponse{
		Status:  "error",
		Message: "Failed to restore card",
	}

	mockCardService.EXPECT().
		RestoreCard(1).
		Return(nil, mockError).
		Times(1)

	response, err := cardHandler.RestoreCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Failed to restore card")
}

func TestRestoreCard_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 0,
	}

	response, err := cardHandler.RestoreCard(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())

	assert.Contains(t, err.Error(), "Failed to restore card: ")
}

func TestDeleteCardPermanent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 1,
	}

	mockResponse := &pb.ApiResponseCardDelete{
		Status:  "success",
		Message: "Successfully deleted card",
	}

	mockCardService.EXPECT().
		DeleteCardPermanent(1).
		Return(mockResponse, nil).
		Times(1)

	response, err := cardHandler.DeleteCardPermanent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.GetStatus())
	assert.Equal(t, "Successfully deleted card", response.GetMessage())
}

func TestDeleteCardPermanent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 1,
	}

	mockCardService.EXPECT().
		DeleteCardPermanent(1).
		Return(nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete card",
		}).
		Times(1)

	response, err := cardHandler.DeleteCardPermanent(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "Failed to delete card")
}

func TestDeleteCardPermanent_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := mock_service.NewMockCardService(ctrl)
	cardHandler := gapi.NewCardHandleGrpc(mockCardService, nil)

	req := &pb.FindByIdCardRequest{
		CardId: 0,
	}

	response, err := cardHandler.DeleteCardPermanent(context.Background(), req)

	assert.Nil(t, response)
	assert.Error(t, err)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	assert.Contains(t, err.Error(), "Failed to delete card: ")
}
