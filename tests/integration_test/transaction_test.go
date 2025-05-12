package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerTestSuite) TestFindAllTransaction() {
	s.Run("success find all transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransaction{
			Status:  "success",
			Message: "Transactions retrieved successfully",
			Data: []*pb.TransactionResponse{
				{
					Id:         1,
					CardNumber: "1234567890123456",
				},
				{
					Id:         2,
					CardNumber: "1234567890123457",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    2,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindAllTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, 2)
		s.Equal(expectedResponse.Pagination.CurrentPage, res.Pagination.CurrentPage)
	})

	s.Run("failure find all transaction", func() {
		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		_, err := s.transactionClient.FindAllTransaction(context.Background(), req)

		s.Error(err)
	})

	s.Run("empty find all transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransaction{
			Status:  "success",
			Message: "No transactions found",
			Data:    nil,
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  0,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindAllTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
		s.Equal(expectedResponse.Pagination.CurrentPage, res.Pagination.CurrentPage)
	})
}

func (s *ServerTestSuite) TestFindByIdTransaction() {
	s.Run("success find by id transaction", func() {
		id := 1
		expectedResponse := &pb.ApiResponseTransaction{
			Status:  "success",
			Message: "Transaction retrieved successfully",
			Data: &pb.TransactionResponse{
				Id:         int32(id),
				CardNumber: "1234567890123456",
			},
		}

		req := &pb.FindByIdTransactionRequest{
			TransactionId: int32(id),
		}

		res, err := s.transactionClient.FindByIdTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("failure find by id transaction", func() {
		id := 999
		req := &pb.FindByIdTransactionRequest{
			TransactionId: int32(id),
		}

		res, err := s.transactionClient.FindByIdTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
	})

	s.Run("invalid id find by id transaction", func() {
		id := -1

		req := &pb.FindByIdTransactionRequest{
			TransactionId: int32(id),
		}

		res, err := s.transactionClient.FindByIdTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
	})
}

func (s *ServerTestSuite) TestFindByCardNumberTransaction() {
	s.Run("success find by card number transaction", func() {
		cardNumber := "1234567890123456"

		expectedResponse := &pb.ApiResponseTransactions{
			Status:  "success",
			Message: "Transaction retrieved successfully",
			Data: []*pb.TransactionResponse{
				{
					Id:         1,
					CardNumber: cardNumber,
				},
				{
					Id:         2,
					CardNumber: cardNumber,
				},
			},
		}

		req := &pb.FindByCardNumberTransactionRequest{
			CardNumber: cardNumber,
		}

		res, err := s.transactionClient.FindByCardNumberTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		for i, data := range expectedResponse.Data {
			s.Equal(data.Id, res.Data[i].Id)
			s.Equal(data.CardNumber, res.Data[i].CardNumber)
		}
	})

	s.Run("failure find by card number transaction", func() {
		cardNumber := "invalid_card_number"

		expectedResponse := &pb.ApiResponseTransactions{
			Status:  "error",
			Message: "Transaction not found",
			Data:    nil,
		}

		req := &pb.FindByCardNumberTransactionRequest{
			CardNumber: cardNumber,
		}

		res, err := s.transactionClient.FindByCardNumberTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res.Data)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestFindTransactionMerchantId() {
	s.Run("success find transaction by merchant id", func() {
		merchantId := 12345

		expectedResponse := &pb.ApiResponseTransactions{
			Status:  "success",
			Message: "Transaction retrieved successfully",
			Data: []*pb.TransactionResponse{
				{
					Id:         1,
					CardNumber: "1234567890123456",
				},
				{
					Id:         2,
					CardNumber: "1234567890123457",
				},
			},
		}

		req := &pb.FindTransactionByMerchantIdRequest{
			MerchantId: int32(merchantId),
		}

		res, err := s.transactionClient.FindTransactionByMerchantId(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		for i, data := range expectedResponse.Data {
			s.Equal(data.Id, res.Data[i].Id)
			s.Equal(data.CardNumber, res.Data[i].CardNumber)
		}
	})

	s.Run("failure find transaction by merchant id", func() {
		merchantId := 99999

		expectedResponse := &pb.ApiResponseTransactions{
			Status:  "error",
			Message: "Transaction not found",
			Data:    nil,
		}

		req := &pb.FindTransactionByMerchantIdRequest{
			MerchantId: int32(merchantId),
		}

		res, err := s.transactionClient.FindTransactionByMerchantId(context.Background(), req)

		s.Error(err)
		s.Nil(res.Data)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("invalid id find transaction by merchant id", func() {
		merchantId := -1

		expectedResponse := &pb.ApiResponseTransactions{
			Status:  "error",
			Message: "Invalid merchant ID",
			Data:    nil,
		}

		req := &pb.FindTransactionByMerchantIdRequest{
			MerchantId: int32(merchantId),
		}

		res, err := s.transactionClient.FindTransactionByMerchantId(context.Background(), req)

		s.Error(err)
		s.Nil(res.Data)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestFindByActiveTransaction() {
	s.Run("success find by active transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
			Status:  "success",
			Message: "Transactions retrieved successfully",
			Data: []*pb.TransactionResponseDeleteAt{
				{
					Id:         1,
					CardNumber: "1234567890123456",
				},
				{
					Id:         2,
					CardNumber: "1234567890123457",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    2,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindByActiveTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		for i, data := range expectedResponse.Data {
			s.Equal(data.Id, res.Data[i].Id)
			s.Equal(data.CardNumber, res.Data[i].CardNumber)
		}
	})

	s.Run("failure find by active transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
			Status:  "success",
			Message: "Transactions retrieved successfully",
			Data:    []*pb.TransactionResponseDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    2,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindByActiveTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res.Data)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("empty find by active transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
			Status:  "success",
			Message: "Transactions retrieved successfully",
			Data: []*pb.TransactionResponseDeleteAt{
				{
					Id:         1,
					CardNumber: "1234567890123456",
				},
				{
					Id:         2,
					CardNumber: "1234567890123457",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    2,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindByActiveTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByTrashTransaction() {
	s.Run("success find by trash transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
			Status:  "success",
			Message: "Transactions retrieved successfully",
			Data: []*pb.TransactionResponseDeleteAt{
				{
					Id:         1,
					CardNumber: "1234567890123456",
				},
				{
					Id:         2,
					CardNumber: "1234567890123457",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    2,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindByTrashedTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, len(expectedResponse.Data))
		for i, data := range expectedResponse.Data {
			s.Equal(data.Id, res.Data[i].Id)
			s.Equal(data.CardNumber, res.Data[i].CardNumber)
		}
	})

	s.Run("failure find by trash transaction", func() {
		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindByTrashedTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), "Internal server error")
	})

	s.Run("empty find by trash transaction", func() {
		expectedResponse := &pb.ApiResponsePaginationTransactionDeleteAt{
			Status:  "success",
			Message: "Transactions retrieved successfully",
			Data:    []*pb.TransactionResponseDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    2,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTransactionRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.transactionClient.FindByTrashedTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestCreateTransaction() {
	s.Run("success create transaction", func() {
		req := &pb.CreateTransactionRequest{
			CardNumber:      "1234567890123456",
			Amount:          500000,
			PaymentMethod:   "mandiri",
			MerchantId:      1,
			TransactionTime: timestamppb.Now(),
		}

		expectedResponse := &pb.ApiResponseTransaction{
			Status:  "success",
			Message: "Transaction created successfully",
			Data: &pb.TransactionResponse{
				Id:         1,
				CardNumber: "1234567890123456",
			},
		}

		res, err := s.transactionClient.CreateTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("failure create transaction", func() {
		req := &pb.CreateTransactionRequest{
			CardNumber:      "1234567890123456",
			Amount:          500000,
			PaymentMethod:   "mandiri",
			MerchantId:      99999,
			TransactionTime: timestamppb.Now(),
		}

		expectedErrorMessage := "merchant not found"

		res, err := s.transactionClient.CreateTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("validation create transaction", func() {
		req := &pb.CreateTransactionRequest{
			CardNumber:      "",
			Amount:          -500000,
			PaymentMethod:   "invalid_method",
			MerchantId:      1,
			TransactionTime: timestamppb.Now(),
		}

		expectedErrorMessage := "validation error"

		res, err := s.transactionClient.CreateTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestUpdateTransaction() {
	s.Run("success update transaction", func() {
		req := &pb.UpdateTransactionRequest{
			TransactionId:   1,
			CardNumber:      "1234567890123456",
			Amount:          600000,
			PaymentMethod:   "mandiri",
			TransactionTime: timestamppb.Now(),
		}

		expectedResponse := &pb.ApiResponseTransaction{
			Status:  "success",
			Message: "Transaction updated successfully",
			Data: &pb.TransactionResponse{
				Id:         1,
				CardNumber: "1234567890123456",
				Amount:     600000,
			},
		}

		res, err := s.transactionClient.UpdateTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.Amount, res.Data.Amount)
	})

	s.Run("failure update transaction", func() {
		req := &pb.UpdateTransactionRequest{
			TransactionId:   99999,
			CardNumber:      "1234567890123456",
			Amount:          600000,
			PaymentMethod:   "mandiri",
			TransactionTime: timestamppb.Now(),
		}

		expectedErrorMessage := "transaction not found"

		res, err := s.transactionClient.UpdateTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("validation update transaction", func() {
		req := &pb.UpdateTransactionRequest{
			TransactionId:   1,
			CardNumber:      "",
			Amount:          -1000,
			PaymentMethod:   "invalid_method",
			TransactionTime: timestamppb.Now(),
		}

		expectedErrorMessage := "validation error"

		res, err := s.transactionClient.UpdateTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestTrashedTransaction() {
	s.Run("success trashed transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}

		expectedResponse := &pb.ApiResponseTransaction{
			Status:  "success",
			Message: "Transaction trashed successfully",
			Data: &pb.TransactionResponse{
				Id:         1,
				CardNumber: "1234567890123456",
			},
		}

		res, err := s.transactionClient.TrashedTransaction(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("failure trashed transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: 99999,
		}

		expectedErrorMessage := "transaction not found"

		res, err := s.transactionClient.TrashedTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("invalid id trashed transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: -1,
		}

		expectedErrorMessage := "invalid transaction ID"

		res, err := s.transactionClient.TrashedTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestRestoreTransaction() {
	// Test case: Success
	s.Run("success restore transaction", func() {
		// Persiapkan data
		req := &pb.FindByIdTransactionRequest{
			TransactionId: 1, // ID valid dari transaksi yang ingin di-restore
		}

		expectedResponse := &pb.ApiResponseTransaction{
			Status:  "success",
			Message: "Transaction restored successfully",
			Data: &pb.TransactionResponse{
				Id:         1,
				CardNumber: "1234567890123456",
			},
		}

		// Eksekusi API langsung (tanpa mock untuk integration test)
		res, err := s.transactionClient.RestoreTransaction(context.Background(), req)

		// Validasi hasil
		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
	})

	// Test case: Failure (e.g., transaction not found)
	s.Run("failure restore transaction", func() {
		// Persiapkan data untuk transaksi yang tidak ditemukan
		req := &pb.FindByIdTransactionRequest{
			TransactionId: 99999, // ID yang tidak ada
		}

		expectedErrorMessage := "transaction not found"

		// Eksekusi API
		res, err := s.transactionClient.RestoreTransaction(context.Background(), req)

		// Validasi hasil
		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	// Test case: Invalid ID
	s.Run("invalid id restore transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: -1, // ID negatif atau tidak valid
		}

		expectedErrorMessage := "invalid transaction ID"

		// Eksekusi API
		res, err := s.transactionClient.RestoreTransaction(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestDeletePermanentTransaction() {
	s.Run("success delete permanent transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: 1,
		}

		expectedResponse := &pb.ApiResponseTransactionDelete{
			Status:  "success",
			Message: "Transaction deleted permanently",
		}

		res, err := s.transactionClient.DeleteTransactionPermanent(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure delete permanent transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: 99999,
		}

		expectedErrorMessage := "transaction not found"

		res, err := s.transactionClient.DeleteTransactionPermanent(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("invalid id delete permanent transaction", func() {
		req := &pb.FindByIdTransactionRequest{
			TransactionId: -1,
		}

		expectedErrorMessage := "invalid transaction ID"

		res, err := s.transactionClient.DeleteTransactionPermanent(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}
