package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"context"
)

func (s *ServerTestSuite) TestFindAllTopup() {
	s.Run("success find all topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTopup{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data: []*pb.TopupResponse{
				{
					Id:          1,
					CardNumber:  "1234567890",
					TopupAmount: 10000,
				},
				{
					Id:          2,
					CardNumber:  "0987654321",
					TopupAmount: 20000,
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  1,
			},
		}

		res, err := s.topupClient.FindAllTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		s.Equal(expectedResponse.Pagination.TotalPages, res.Pagination.TotalPages)
	})

	s.Run("failure find all topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTopup{
			Status:  "error",
			Message: "Failed to fetch topup data",
		}

		res, err := s.topupClient.FindAllTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("empty topup data", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "nonexistent",
		}

		expectedResponse := &pb.ApiResponsePaginationTopup{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data:    []*pb.TopupResponse{},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  0,
			},
		}

		res, err := s.topupClient.FindAllTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(0, len(res.Data))
		s.Equal(expectedResponse.Pagination.TotalPages, res.Pagination.TotalPages)
	})
}

func (s *ServerTestSuite) TestFindByIdTopup() {
	s.Run("success find topup by id", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data: &pb.TopupResponse{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		}

		res, err := s.topupClient.FindByIdTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.TopupAmount, res.Data.TopupAmount)
	})

	s.Run("failure find topup by id", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 99,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		}

		res, err := s.topupClient.FindByIdTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("invalid id find topup by id", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 0,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Invalid topup ID",
		}

		res, err := s.topupClient.FindByIdTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestFindByCardNumberTopup() {
	s.Run("success find topup by card number", func() {
		req := &pb.FindByCardNumberTopupRequest{
			CardNumber: "1234567890",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data: &pb.TopupResponse{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		}

		res, err := s.topupClient.FindByCardNumberTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.TopupAmount, res.Data.TopupAmount)
	})

	s.Run("failure find topup by card number", func() {
		req := &pb.FindByCardNumberTopupRequest{
			CardNumber: "9876543210",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to retrieve topup data",
		}

		res, err := s.topupClient.FindByCardNumberTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestFindByActiveTopup() {
	s.Run("success find active topup", func() {
		expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data: []*pb.TopupResponseDeleteAt{
				{
					Id:          1,
					CardNumber:  "1234567890",
					TopupAmount: 10000,
				},
				{
					Id:          2,
					CardNumber:  "0987654321",
					TopupAmount: 20000,
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  1,
			},
		}

		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.topupClient.FindByActive(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
		s.Equal(expectedResponse.Data[0].CardNumber, res.Data[0].CardNumber)
		s.Equal(expectedResponse.Data[0].TopupAmount, res.Data[0].TopupAmount)
	})

	s.Run("failure find active topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data:    []*pb.TopupResponseDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  1,
			},
		}

		res, err := s.topupClient.FindByActive(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("empty find active topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data:    []*pb.TopupResponseDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  1,
			},
		}

		res, err := s.topupClient.FindByActive(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByTrashedTopup() {
	s.Run("success find trashed topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data: []*pb.TopupResponseDeleteAt{
				{
					Id:          1,
					CardNumber:  "1234567890",
					TopupAmount: 10000,
				},
				{
					Id:          2,
					CardNumber:  "0987654321",
					TopupAmount: 20000,
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  1,
			},
		}

		res, err := s.topupClient.FindByTrashed(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
		s.Equal(expectedResponse.Data[0].CardNumber, res.Data[0].CardNumber)
		s.Equal(expectedResponse.Data[0].TopupAmount, res.Data[0].TopupAmount)
	})

	s.Run("failure find trashed topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationTopupDeleteAt{
			Status:  "success",
			Message: "Topup data retrieved successfully",
			Data:    []*pb.TopupResponseDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage: 1,
				PageSize:    10,
				TotalPages:  1,
			},
		}

		res, err := s.topupClient.FindByTrashed(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("empty find trashed topup", func() {
		req := &pb.FindAllTopupRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsesTopup{
			Status:  "success",
			Message: "No trashed topup data found",
			Data:    []*pb.TopupResponse{},
		}

		res, err := s.topupClient.FindByTrashed(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestCreateTopup() {
	s.Run("success create topup", func() {
		req := &pb.CreateTopupRequest{
			CardNumber:  "1234567890",
			TopupAmount: 10000,
			TopupNo:     "1234567890",
			TopupMethod: "mandiri",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "success",
			Message: "Topup created successfully",
			Data: &pb.TopupResponse{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		}

		res, err := s.topupClient.CreateTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.TopupAmount, res.Data.TopupAmount)
	})

	s.Run("failure create topup", func() {
		req := &pb.CreateTopupRequest{
			CardNumber:  "1234567890",
			TopupAmount: 10000,
			TopupNo:     "1234567890",
			TopupMethod: "mandiri",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to create topup",
			Data:    nil,
		}

		res, err := s.topupClient.CreateTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data, res.Data)
	})

	s.Run("validation create topup", func() {
		req := &pb.CreateTopupRequest{
			CardNumber:  "1234567890",
			TopupAmount: 10000,
			TopupNo:     "1234567890",
			TopupMethod: "mandirii",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to create topup",
			Data:    nil,
		}

		res, err := s.topupClient.CreateTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data, res.Data)
	})
}

func (s *ServerTestSuite) TestUpdateTopup() {
	s.Run("success update topup", func() {
		req := &pb.UpdateTopupRequest{
			TopupId:     1,
			CardNumber:  "1234567890",
			TopupAmount: 15000,
			TopupMethod: "mandiri",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "success",
			Message: "Topup updated successfully",
			Data: &pb.TopupResponse{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 15000,
			},
		}

		res, err := s.topupClient.UpdateTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.TopupAmount, res.Data.TopupAmount)
	})

	s.Run("failure update topup", func() {
		req := &pb.UpdateTopupRequest{
			TopupId:     1,
			CardNumber:  "1234567890",
			TopupAmount: 15000,
			TopupMethod: "mandiri",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to update topup",
			Data:    nil,
		}

		res, err := s.topupClient.UpdateTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data, res.Data)
	})

	s.Run("validation update topup", func() {
		req := &pb.UpdateTopupRequest{
			TopupId:     1,
			CardNumber:  "1234567890",
			TopupAmount: 15000,

			TopupMethod: "mandirii",
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to update topup",
			Data:    nil,
		}

		res, err := s.topupClient.UpdateTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data, res.Data)
	})
}

func (s *ServerTestSuite) TestTrashTopup() {
	s.Run("success trash topup", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "success",
			Message: "Topup deleted successfully",
			Data: &pb.TopupResponse{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		}

		res, err := s.topupClient.TrashedTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.TopupAmount, res.Data.TopupAmount)
	})

	s.Run("failure trash topup", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Failed to delete topup",
		}

		res, err := s.topupClient.TrashedTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("Invalid topup id", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: -1,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "error",
			Message: "Invalid topup id",
		}

		res, err := s.topupClient.TrashedTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestRestoreTopup() {
	s.Run("success restore topup", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopup{
			Status:  "success",
			Message: "Topup restored successfully",
			Data: &pb.TopupResponse{
				Id:          1,
				CardNumber:  "1234567890",
				TopupAmount: 10000,
			},
		}

		res, err := s.topupClient.RestoreTopup(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.TopupAmount, res.Data.TopupAmount)
	})

	s.Run("failure restore topup", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopupDeleteAt{
			Status:  "error",
			Message: "Failed to restore topup",
		}

		res, err := s.topupClient.RestoreTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("Invalid topup id", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: -1,
		}

		expectedResponse := &pb.ApiResponseTopupDeleteAt{
			Status:  "error",
			Message: "Invalid topup id",
		}

		res, err := s.topupClient.RestoreTopup(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestDeleteTopupPermanent() {
	s.Run("success delete topup", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopupDelete{
			Status:  "success",
			Message: "Topup deleted successfully",
		}

		res, err := s.topupClient.DeleteTopupPermanent(context.Background(), req)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure delete topup", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopupDelete{
			Status:  "error",
			Message: "Failed to delete topup",
		}

		res, err := s.topupClient.DeleteTopupPermanent(context.Background(), req)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("Invalid ID", func() {
		req := &pb.FindByIdTopupRequest{
			TopupId: 1,
		}

		expectedResponse := &pb.ApiResponseTopupDelete{
			Status:  "error",
			Message: "Bad Request: Invalid ID",
		}

		res, err := s.topupClient.DeleteTopupPermanent(context.Background(), req)

		s.Error(err)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}
