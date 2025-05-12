package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerTestSuite) TestFindAllCard(t *testing.T) {
	s.Run("Success Find All Cards", func() {
		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationCard{
			Status:  "success",
			Message: "Fetched cards successfully",
			Data: []*pb.CardResponse{
				{Id: 1, CardNumber: "Card 1"},
				{Id: 2, CardNumber: "Card 2"},
			},
		}

		res, err := s.cardClient.FindAllCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, len(expectedResponse.Data))

		for i := range expectedResponse.Data {
			s.Equal(expectedResponse.Data[i].Id, res.Data[i].Id)
			s.Equal(expectedResponse.Data[i].CardNumber, res.Data[i].CardNumber)
		}
	})

	s.Run("Empty Find All Cards", func() {
		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedResponse := &pb.ApiResponsePaginationCard{
			Status:  "success",
			Message: "No cards found",
			Data:    []*pb.CardResponse{},
		}

		res, err := s.cardClient.FindAllCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, len(expectedResponse.Data))
	})

	s.Run("Error Find All Cards", func() {
		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedResponse := &pb.ApiResponsePaginationCard{
			Status:  "error",
			Message: "Failed to fetch card records",
			Data:    []*pb.CardResponse{},
		}

		res, err := s.cardClient.FindAllCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, len(expectedResponse.Data))
	})

}

func (s *ServerTestSuite) TestFindByIdCard() {
	s.Run("Success Find Card by ID", func() {
		cardID := int32(1)
		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully fetched card record",
			Data: &pb.CardResponse{
				Id:         cardID,
				CardNumber: "1234567890123456",
			},
		}

		res, err := s.cardClient.FindByIdCard(s.ctx, &pb.FindByIdCardRequest{CardId: cardID})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("Failure - Card Not Found", func() {
		cardID := int32(999)
		expectedError := "card not found"

		res, err := s.cardClient.FindByIdCard(s.ctx, &pb.FindByIdCardRequest{CardId: cardID})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

	s.Run("Failure - Invalid Card ID", func() {
		cardID := int32(-1)
		expectedError := "invalid card ID"

		res, err := s.cardClient.FindByIdCard(s.ctx, &pb.FindByIdCardRequest{CardId: cardID})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})
}

func (s *ServerTestSuite) TestFindByUserIDCard() {
	s.Run("Success Find Card by UserID", func() {
		userID := int32(42)
		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully fetched card record",
			Data: &pb.CardResponse{
				Id:         1,
				CardNumber: "1234567890123456",
			},
		}

		res, err := s.cardClient.FindByUserIdCard(s.ctx, &pb.FindByUserIdCardRequest{UserId: userID})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("Failure - Invalid UserID", func() {
		userID := int32(-1)
		expectedError := "invalid user ID"

		res, err := s.cardClient.FindByUserIdCard(s.ctx, &pb.FindByUserIdCardRequest{UserId: userID})

		s.Error(err)

		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

	s.Run("Failure - Card Fetch Error", func() {
		userID := int32(42)
		expectedError := "internal server error"

		res, err := s.cardClient.FindByUserIdCard(s.ctx, &pb.FindByUserIdCardRequest{UserId: userID})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})
}

func (s *ServerTestSuite) TestFindByActiveCard() {
	s.Run("Success Find Active Card", func() {
		expectedResponse := &pb.ApiResponsePaginationCardDeleteAt{
			Status:  "success",
			Message: "Fetched cards successfully",
			Data: []*pb.CardResponseDeleteAt{
				{Id: 1, CardNumber: "Card 1"},
				{Id: 2, CardNumber: "Card 2"},
			},
		}

		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.cardClient.FindByActiveCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)

		s.Len(res.Data, len(expectedResponse.Data))
	})

	s.Run("Failure - Card Fetch Error", func() {
		expectedError := "internal server error"

		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.cardClient.FindByActiveCard(s.ctx, findAllCardsRequest)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

	s.Run("Failure - No Active Cards", func() {
		expectedResponse := &pb.ApiResponseCards{
			Status:  "success",
			Message: "No active cards found",
			Data:    []*pb.CardResponse{},
		}

		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.cardClient.FindByActiveCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, len(expectedResponse.Data))
	})
}

func (s *ServerTestSuite) TestfindByTrashedCard() {
	s.Run("Success Find Trashed Card", func() {
		expectedResponse := &pb.ApiResponsePaginationCardDeleteAt{
			Status:  "success",
			Message: "Fetched cards successfully",
			Data: []*pb.CardResponseDeleteAt{
				{Id: 1, CardNumber: "Card 1"},
				{Id: 2, CardNumber: "Card 2"},
			},
		}

		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.cardClient.FindByTrashedCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)

		s.Len(res.Data, len(expectedResponse.Data))
	})

	s.Run("Failure - Card Fetch Error", func() {
		expectedError := "internal server error"

		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.cardClient.FindByTrashedCard(s.ctx, findAllCardsRequest)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

	s.Run("Failure - No Trashed Cards", func() {
		expectedResponse := &pb.ApiResponseCards{
			Status:  "success",
			Message: "No trashed cards found",
			Data:    []*pb.CardResponse{},
		}

		findAllCardsRequest := &pb.FindAllCardRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.cardClient.FindByTrashedCard(s.ctx, findAllCardsRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, len(expectedResponse.Data))
	})
}

func (s *ServerTestSuite) TestFindByCardNumberCard() {
	s.Run("Success Find Card by Card Number", func() {
		cardNumber := "1234567890123456"
		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully fetched card record",
			Data: &pb.CardResponse{
				Id:         1,
				CardNumber: cardNumber,
			},
		}

		res, err := s.cardClient.FindByCardNumber(s.ctx, &pb.FindByCardNumberRequest{CardNumber: cardNumber})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("Failure - Card Not Found", func() {
		cardNumber := "9999999999999999"
		expectedError := "card not found"

		res, err := s.cardClient.FindByCardNumber(s.ctx, &pb.FindByCardNumberRequest{CardNumber: cardNumber})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

	s.Run("Failure - Invalid Card Number", func() {
		cardNumber := "1234"
		expectedError := "invalid card number"

		res, err := s.cardClient.FindByCardNumber(s.ctx, &pb.FindByCardNumberRequest{CardNumber: cardNumber})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

}

func (s *ServerTestSuite) TestCreateCard() {
	s.Run("Success Create Card", func() {
		cardNumber := "1234567890123456"
		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully created card",
			Data: &pb.CardResponse{
				Id:         1,
				CardNumber: cardNumber,
			},
		}

		res, err := s.cardClient.CreateCard(s.ctx, &pb.CreateCardRequest{
			UserId:       1,
			CardType:     "credit",
			ExpireDate:   timestamppb.Now(),
			Cvv:          "123",
			CardProvider: "mandiri",
		})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("Failure - Card Not Found", func() {
		expectedError := "card not found"

		res, err := s.cardClient.CreateCard(s.ctx, &pb.CreateCardRequest{
			UserId:       1,
			CardType:     "credit",
			ExpireDate:   timestamppb.Now(),
			Cvv:          "123",
			CardProvider: "mandiri",
		})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})

	s.Run("Failure - Validation Error", func() {
		expectedError := "card type must be credit or debit"

		res, err := s.cardClient.CreateCard(s.ctx, &pb.CreateCardRequest{
			UserId:       1,
			CardType:     "invalid",
			ExpireDate:   timestamppb.Now(),
			Cvv:          "123",
			CardProvider: "mandiri",
		})

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedError)
	})
}

func (s *ServerTestSuite) TestUpdateCard() {
	s.Run("Success Update Card", func() {
		cardID := 1
		reqBody := requests.UpdateCardRequest{
			CardID:       cardID,
			UserID:       1,
			CardType:     "debit",
			ExpireDate:   time.Now().AddDate(2, 0, 0),
			CVV:          "456",
			CardProvider: "MasterCard",
		}

		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully updated card record",
			Data: &pb.CardResponse{
				Id:           int32(cardID),
				UserId:       int32(reqBody.UserID),
				CardNumber:   "1234567890123456",
				CardType:     reqBody.CardType,
				ExpireDate:   reqBody.ExpireDate.String(),
				Cvv:          reqBody.CVV,
				CardProvider: reqBody.CardProvider,
			},
		}

		req := &pb.UpdateCardRequest{
			CardId:       int32(cardID),
			UserId:       int32(reqBody.UserID),
			CardType:     reqBody.CardType,
			ExpireDate:   timestamppb.New(reqBody.ExpireDate),
			Cvv:          reqBody.CVV,
			CardProvider: reqBody.CardProvider,
		}

		res, err := s.cardClient.UpdateCard(s.ctx, req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardType, res.Data.CardType)
		s.Equal(expectedResponse.Data.CardProvider, res.Data.CardProvider)
	})

	s.Run("Failure Update Card", func() {
		cardID := int32(1)
		req := &pb.UpdateCardRequest{
			CardId:       cardID,
			UserId:       1,
			CardType:     "debit",
			ExpireDate:   timestamppb.New(time.Now().AddDate(2, 0, 0)),
			Cvv:          "456",
			CardProvider: "MasterCard",
		}

		_, err := s.cardClient.UpdateCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "internal server error")
	})

	s.Run("Validation Error", func() {
		cardID := int32(1)
		req := &pb.UpdateCardRequest{
			CardId:       cardID,
			UserId:       1,
			CardType:     "",
			ExpireDate:   timestamppb.New(time.Now().AddDate(2, 0, 0)),
			Cvv:          "456",
			CardProvider: "MasterCard",
		}

		_, err := s.cardClient.UpdateCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "Validation Error")
	})

	s.Run("Invalid Card ID", func() {
		req := &pb.UpdateCardRequest{
			CardId:       0,
			UserId:       1,
			CardType:     "debit",
			ExpireDate:   timestamppb.New(time.Now().AddDate(2, 0, 0)),
			Cvv:          "456",
			CardProvider: "MasterCard",
		}

		_, err := s.cardClient.UpdateCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "Invalid Card ID")
	})
}

func (s *ServerTestSuite) TestTrashedCard() {
	s.Run("Success Trashed Card", func() {
		cardID := int32(1)

		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully trashed card record",
			Data: &pb.CardResponse{
				Id:           cardID,
				CardType:     "debit",
				ExpireDate:   "2025-12-31",
				Cvv:          "123",
				CardProvider: "Visa",
			},
		}

		req := &pb.FindByIdCardRequest{CardId: cardID}

		res, err := s.cardClient.TrashedCard(s.ctx, req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardType, res.Data.CardType)
		s.Equal(expectedResponse.Data.CardProvider, res.Data.CardProvider)
	})

	s.Run("Failure Trashed Card", func() {
		cardID := int32(1)
		req := &pb.FindByIdCardRequest{CardId: cardID}

		_, err := s.cardClient.TrashedCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "internal server error")
	})

	s.Run("Invalid Card ID", func() {

		req := &pb.FindByIdCardRequest{
			CardId: 0,
		}

		_, err := s.cardClient.TrashedCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "Invalid Card ID")
	})
}

func (s *ServerTestSuite) TestRestoreCard() {
	s.Run("Success Restore Card", func() {
		cardID := int32(1)

		expectedResponse := &pb.ApiResponseCard{
			Status:  "success",
			Message: "Successfully restored card record",
			Data: &pb.CardResponse{
				Id:           cardID,
				CardType:     "debit",
				ExpireDate:   "2025-12-31",
				Cvv:          "123",
				CardProvider: "Visa",
			},
		}

		req := &pb.FindByIdCardRequest{CardId: cardID}

		res, err := s.cardClient.RestoreCard(s.ctx, req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardType, res.Data.CardType)
		s.Equal(expectedResponse.Data.CardProvider, res.Data.CardProvider)
	})

	s.Run("Failure Restore Card", func() {
		cardID := int32(1)
		req := &pb.FindByIdCardRequest{CardId: cardID}

		_, err := s.cardClient.RestoreCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "internal server error")
	})

	s.Run("Invalid Card ID", func() {
		req := &pb.FindByIdCardRequest{
			CardId: 0,
		}

		_, err := s.cardClient.RestoreCard(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "Invalid Card ID")
	})
}

func (s *ServerTestSuite) TestDeletePermanentCard() {
	s.Run("Success Delete Permanent Card", func() {
		cardID := int32(1)

		expectedResponse := &pb.ApiResponseCardDelete{
			Status:  "success",
			Message: "Successfully deleted card record permanently",
		}

		req := &pb.FindByIdCardRequest{CardId: cardID}

		res, err := s.cardClient.DeleteCardPermanent(s.ctx, req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("Failure Delete Permanent Card", func() {
		cardID := int32(1)
		req := &pb.FindByIdCardRequest{CardId: cardID}

		_, err := s.cardClient.DeleteCardPermanent(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "internal server error")
	})

	s.Run("Invalid Card ID", func() {
		req := &pb.FindByIdCardRequest{
			CardId: 0,
		}

		_, err := s.cardClient.DeleteCardPermanent(s.ctx, req)

		s.Error(err)
		s.Contains(err.Error(), "Invalid Card ID")
	})
}
