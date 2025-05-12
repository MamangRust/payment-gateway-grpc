package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerTestSuite) TestFindAllTransfer() {
	s.Run("success find all transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTransfer{
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

		res, err := s.transferClient.FindAllTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Pagination, res.Pagination)
		s.Len(res.Data, len(expectedResponse.Data))
		s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
		s.Equal(expectedResponse.Data[0].TransferFrom, res.Data[0].TransferFrom)
	})

	s.Run("failure find all transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedErrorMessage := "internal server error"

		res, err := s.transferClient.FindAllTransfer(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("empty find all transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTransfer{
			Status:  "success",
			Message: "No transfers found",
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				TotalPages:   0,
				TotalRecords: 0,
			},
			Data: []*pb.TransferResponse{},
		}

		res, err := s.transferClient.FindAllTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Pagination, res.Pagination)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByIdTransfer() {
	s.Run("success find by id transfer", func() {
		transferID := int32(1)

		req := &pb.FindByIdTransferRequest{
			TransferId: transferID,
		}

		expectedResponse := &pb.ApiResponseTransfer{
			Status:  "success",
			Message: "Transfer retrieved successfully",
			Data: &pb.TransferResponse{
				Id:             transferID,
				TransferFrom:   "test",
				TransferTo:     "test",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
				CreatedAt:      "2022-01-01 00:00:00",
				UpdatedAt:      "2022-01-01 00:00:00",
			},
		}

		res, err := s.transferClient.FindByIdTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data, res.Data)
	})

	s.Run("failure find by id transfer", func() {
		transferID := int32(999)

		req := &pb.FindByIdTransferRequest{
			TransferId: transferID,
		}

		expectedErrorMessage := "transfer not found"

		res, err := s.transferClient.FindByIdTransfer(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("invalid id transfer", func() {
		transferID := int32(-1)

		req := &pb.FindByIdTransferRequest{
			TransferId: transferID,
		}

		expectedErrorMessage := "invalid transfer id"

		res, err := s.transferClient.FindByIdTransfer(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestFindByTransferByTransferFrom() {
	s.Run("success find transfer by transfer from", func() {
		transferFrom := "test_user"

		req := &pb.FindTransferByTransferFromRequest{
			TransferFrom: transferFrom,
		}

		expectedResponse := &pb.ApiResponseTransfers{
			Status:  "success",
			Message: "Transfer retrieved successfully",
			Data: []*pb.TransferResponse{
				{
					Id:             1,
					TransferFrom:   transferFrom,
					TransferTo:     "test_to",
					TransferAmount: 10000,
					TransferTime:   "2022-01-01 00:00:00",
				},
			},
		}

		res, err := s.transferClient.FindTransferByTransferFrom(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		if len(res.Data) > 0 {
			s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
			s.Equal(expectedResponse.Data[0].TransferFrom, res.Data[0].TransferFrom)
			s.Equal(expectedResponse.Data[0].TransferTo, res.Data[0].TransferTo)
			s.Equal(expectedResponse.Data[0].TransferAmount, res.Data[0].TransferAmount)
			s.Equal(expectedResponse.Data[0].TransferTime, res.Data[0].TransferTime)
		}
	})

	s.Run("failure find transfer by transfer from", func() {
		transferFrom := "non_existent_user"

		req := &pb.FindTransferByTransferFromRequest{
			TransferFrom: transferFrom,
		}

		expectedErrorMessage := "transfer data not found"

		res, err := s.transferClient.FindTransferByTransferFrom(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestFindByTransferByTransferTo() {
	s.Run("success find transfer by transfer to", func() {
		transferTo := "test_to"

		req := &pb.FindTransferByTransferToRequest{
			TransferTo: transferTo,
		}

		expectedResponse := &pb.ApiResponseTransfers{
			Status:  "success",
			Message: "Transfer retrieved successfully",
			Data: []*pb.TransferResponse{
				{
					Id:             1,
					TransferFrom:   "test_from",
					TransferTo:     transferTo,
					TransferAmount: 10000,
					TransferTime:   "2022-01-01 00:00:00",
				},
			},
		}

		res, err := s.transferClient.FindTransferByTransferTo(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		if len(res.Data) > 0 {
			s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
			s.Equal(expectedResponse.Data[0].TransferFrom, res.Data[0].TransferFrom)
			s.Equal(expectedResponse.Data[0].TransferTo, res.Data[0].TransferTo)
			s.Equal(expectedResponse.Data[0].TransferAmount, res.Data[0].TransferAmount)
			s.Equal(expectedResponse.Data[0].TransferTime, res.Data[0].TransferTime)
		}
	})

	s.Run("failure find transfer by transfer to", func() {
		transferTo := "non_existent_to"

		req := &pb.FindTransferByTransferToRequest{
			TransferTo: transferTo,
		}

		expectedErrorMessage := "transfer data not found"

		res, err := s.transferClient.FindTransferByTransferTo(context.Background(), req)

		s.Error(err)
		s.Nil(res)
		s.Contains(err.Error(), expectedErrorMessage)
	})
}

func (s *ServerTestSuite) TestFindByActiveTransfer() {
	s.Run("success find active transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponseTransfers{
			Status:  "success",
			Message: "Transfer retrieved successfully",
			Data: []*pb.TransferResponse{
				{
					Id:             1,
					TransferFrom:   "test_from",
					TransferTo:     "test_to",
					TransferAmount: 10000,
					TransferTime:   "2022-01-01 00:00:00",
				},
			},
		}

		res, err := s.transferClient.FindByActiveTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		if len(res.Data) > 0 {
			s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
			s.Equal(expectedResponse.Data[0].TransferFrom, res.Data[0].TransferFrom)
			s.Equal(expectedResponse.Data[0].TransferTo, res.Data[0].TransferTo)
			s.Equal(expectedResponse.Data[0].TransferAmount, res.Data[0].TransferAmount)
			s.Equal(expectedResponse.Data[0].TransferTime, res.Data[0].TransferTime)
		}
	})

	s.Run("failure find active transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedErrorMessage := "internal server error"

		_, err := s.transferClient.FindByActiveTransfer(context.Background(), req)

		s.Error(err)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("empty find active transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponseTransfers{
			Status:  "success",
			Message: "No active transfers found",
			Data:    []*pb.TransferResponse{},
		}

		res, err := s.transferClient.FindByActiveTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByTrashedTransfer() {
	s.Run("success find trashed transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTransferDeleteAt{
			Status:  "success",
			Message: "Successfully fetch transfers",
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				TotalPages:   1,
				TotalRecords: 1,
			},
			Data: []*pb.TransferResponseDeleteAt{
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

		res, err := s.transferClient.FindByTrashedTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		if len(res.Data) > 0 {
			s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
			s.Equal(expectedResponse.Data[0].TransferFrom, res.Data[0].TransferFrom)
			s.Equal(expectedResponse.Data[0].TransferTo, res.Data[0].TransferTo)
			s.Equal(expectedResponse.Data[0].TransferAmount, res.Data[0].TransferAmount)
			s.Equal(expectedResponse.Data[0].TransferTime, res.Data[0].TransferTime)
		}
	})

	s.Run("failure find trashed transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedErrorMessage := "internal server error"

		_, err := s.transferClient.FindByTrashedTransfer(context.Background(), req)

		s.Error(err)
		s.Contains(err.Error(), expectedErrorMessage)
	})

	s.Run("empty find trashed transfer", func() {
		req := &pb.FindAllTransferRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTransferDeleteAt{
			Status:  "success",
			Message: "Successfully fetch transfers",
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				TotalPages:   1,
				TotalRecords: 1,
			},
			Data: []*pb.TransferResponseDeleteAt{},
		}

		res, err := s.transferClient.FindByTrashedTransfer(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestCreateTransfer() {
	s.Run("success create transfer", func() {
		request := &pb.CreateTransferRequest{
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 10000,
		}

		expectedResponse := &pb.ApiResponseTransfer{
			Status:  "success",
			Message: "Transfer created successfully",
			Data: &pb.TransferResponse{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 10000,
				TransferTime:   "2022-01-01 00:00:00",
			},
		}

		res, err := s.transferClient.CreateTransfer(context.Background(), request)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.TransferFrom, res.Data.TransferFrom)
		s.Equal(expectedResponse.Data.TransferTo, res.Data.TransferTo)
		s.Equal(expectedResponse.Data.TransferAmount, res.Data.TransferAmount)
		s.Equal(expectedResponse.Data.TransferTime, res.Data.TransferTime)
	})

	s.Run("failure create transfer", func() {
		request := &pb.CreateTransferRequest{
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 10000,
		}

		res, err := s.transferClient.CreateTransfer(context.Background(), request)

		s.Error(err)
		s.Nil(res)
	})

	s.Run("validation error create transfer", func() {
		request := &pb.CreateTransferRequest{
			TransferFrom:   "",
			TransferTo:     "test_to",
			TransferAmount: 10000,
		}

		res, err := s.transferClient.CreateTransfer(context.Background(), request)

		s.Error(err)
		s.Nil(res)
	})
}

func (s *ServerTestSuite) TestUpdateTransfer() {
	s.Run("success update transfer", func() {
		request := &pb.UpdateTransferRequest{
			TransferId:     1,
			TransferFrom:   "updated_from",
			TransferTo:     "updated_to",
			TransferAmount: 15000,
		}

		expectedResponse := &pb.ApiResponseTransfer{
			Status:  "success",
			Message: "Transfer updated successfully",
			Data: &pb.TransferResponse{
				Id:             1,
				TransferFrom:   "updated_from",
				TransferTo:     "updated_to",
				TransferAmount: 15000,
				TransferTime:   "2022-01-01 00:00:00",
			},
		}

		res, err := s.transferClient.UpdateTransfer(context.Background(), request)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.TransferFrom, res.Data.TransferFrom)
		s.Equal(expectedResponse.Data.TransferTo, res.Data.TransferTo)
		s.Equal(expectedResponse.Data.TransferAmount, res.Data.TransferAmount)
		s.Equal(expectedResponse.Data.TransferTime, res.Data.TransferTime)
	})

	s.Run("failure update transfer", func() {
		request := &pb.UpdateTransferRequest{
			TransferId:     99,
			TransferFrom:   "test_from",
			TransferTo:     "test_to",
			TransferAmount: 50000,
		}

		expectedError := status.Error(codes.NotFound, "Transfer not found")

		res, err := s.transferClient.UpdateTransfer(context.Background(), request)

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedError.Error(), err.Error())
	})

	s.Run("validation error update transfer", func() {
		request := &pb.UpdateTransferRequest{
			TransferId:     1,
			TransferFrom:   "",
			TransferTo:     "test_to",
			TransferAmount: 5000,
		}

		expectedValidationError := status.Error(codes.InvalidArgument, "TransferFrom is required")

		res, err := s.transferClient.UpdateTransfer(context.Background(), request)

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedValidationError.Error(), err.Error())
	})
}

func (s *ServerTestSuite) TestTrashedTransfer() {
	s.Run("success trashed transfer", func() {
		id := 1
		expectedResponse := &pb.ApiResponseTransfer{
			Status:  "success",
			Message: "Transfer trashed successfully",
			Data: &pb.TransferResponse{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 500000,
				TransferTime:   "2022-01-01 00:00:00",
			},
		}

		request := &pb.FindByIdTransferRequest{TransferId: int32(id)}
		res, err := s.transferClient.TrashedTransfer(context.Background(), request)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.TransferFrom, res.Data.TransferFrom)
		s.Equal(expectedResponse.Data.TransferTo, res.Data.TransferTo)
		s.Equal(expectedResponse.Data.TransferAmount, res.Data.TransferAmount)
		s.Equal(expectedResponse.Data.TransferTime, res.Data.TransferTime)
	})

	s.Run("failure trashed transfer", func() {
		id := 99
		expectedError := status.Error(codes.NotFound, "Transfer not found")

		request := &pb.FindByIdTransferRequest{TransferId: int32(id)}
		res, err := s.transferClient.TrashedTransfer(context.Background(), request)

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedError.Error(), err.Error())
	})

	s.Run("invalid id trashed transfer", func() {
		expectedInvalidIDError := status.Error(codes.InvalidArgument, "Invalid Transfer ID")

		res, err := s.transferClient.TrashedTransfer(context.Background(), &pb.FindByIdTransferRequest{TransferId: 0})

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedInvalidIDError.Error(), err.Error())
	})
}

func (s *ServerTestSuite) TestRestoreTransfer() {
	s.Run("success restore transfer", func() {
		id := 1
		expectedResponse := &pb.ApiResponseTransfer{
			Status:  "success",
			Message: "Transfer restored successfully",
			Data: &pb.TransferResponse{
				Id:             1,
				TransferFrom:   "test_from",
				TransferTo:     "test_to",
				TransferAmount: 500000,
				TransferTime:   "2022-01-01 00:00:00",
			},
		}

		request := &pb.FindByIdTransferRequest{TransferId: int32(id)}
		res, err := s.transferClient.RestoreTransfer(context.Background(), request)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.TransferFrom, res.Data.TransferFrom)
		s.Equal(expectedResponse.Data.TransferTo, res.Data.TransferTo)
		s.Equal(expectedResponse.Data.TransferAmount, res.Data.TransferAmount)
		s.Equal(expectedResponse.Data.TransferTime, res.Data.TransferTime)
	})

	s.Run("failure restore transfer", func() {
		id := 99
		expectedError := status.Error(codes.NotFound, "Transfer not found")

		request := &pb.FindByIdTransferRequest{TransferId: int32(id)}
		res, err := s.transferClient.RestoreTransfer(context.Background(), request)

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedError.Error(), err.Error())
	})

	s.Run("invalid id restore transfer", func() {
		expectedInvalidIDError := status.Error(codes.InvalidArgument, "Invalid Transfer ID")

		res, err := s.transferClient.RestoreTransfer(context.Background(), &pb.FindByIdTransferRequest{TransferId: 0})

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedInvalidIDError.Error(), err.Error())
	})
}

func (s *ServerTestSuite) TestDeleteTransfer() {
	s.Run("success delete transfer", func() {
		id := 1
		expectedResponse := &pb.ApiResponseTransferDelete{
			Status:  "success",
			Message: "Transfer deleted successfully",
		}

		request := &pb.FindByIdTransferRequest{TransferId: int32(id)}
		res, err := s.transferClient.DeleteTransferPermanent(context.Background(), request)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure delete transfer", func() {
		id := 99
		expectedError := status.Error(codes.NotFound, "Transfer not found")

		request := &pb.FindByIdTransferRequest{TransferId: int32(id)}
		res, err := s.transferClient.DeleteTransferPermanent(context.Background(), request)

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedError.Error(), err.Error())
	})

	s.Run("invalid id delete transfer", func() {
		expectedInvalidIDError := status.Error(codes.InvalidArgument, "Invalid Transfer ID")

		res, err := s.transferClient.DeleteTransferPermanent(context.Background(), &pb.FindByIdTransferRequest{TransferId: 0})

		s.Error(err)
		s.Nil(res)
		s.Equal(expectedInvalidIDError.Error(), err.Error())
	})
}
