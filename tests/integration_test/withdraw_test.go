package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerTestSuite) TestFindAllWithdraw() {
	s.Run("success find all withdraw", func() {
		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		expectedResponse := &pb.ApiResponsePaginationWithdraw{
			Status:  "success",
			Message: "Successfully fetch withdraws",
			Pagination: &pb.PaginationMeta{
				TotalRecords: 1,
				TotalPages:   1,
			},
			Data: []*pb.WithdrawResponse{
				{
					WithdrawId: 1,
					CardNumber: "1234567890123456",
				},
			},
		}

		res, err := s.withdrawClient.FindAllWithdraw(s.ctx, findAllWithdrawRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
	s.Run("failure find all withdraw", func() {
		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindAllWithdraw(s.ctx, findAllWithdrawRequest)

		s.Error(err)
		s.Nil(res)
	})
	s.Run("empty find all withdraw", func() {
		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "nonexistent",
		}

		expectedResponse := &pb.ApiResponsePaginationWithdraw{
			Status:  "success",
			Message: "No withdraws found",
			Pagination: &pb.PaginationMeta{
				TotalRecords: 0,
				TotalPages:   0,
			},
			Data: nil,
		}

		res, err := s.withdrawClient.FindAllWithdraw(s.ctx, findAllWithdrawRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, 0)
	})
}

func (s *ServerTestSuite) TestFindByIdWithdraw() {
	s.Run("success find by id withdraw", func() {
		withdrawId := int32(1)

		expectedResponse := &pb.ApiResponseWithdraw{
			Status:  "success",
			Message: "Withdraw found successfully",
			Data: &pb.WithdrawResponse{
				WithdrawId: withdrawId,
				CardNumber: "1234567890123456",
				CreatedAt:  "2022-01-01T00:00:00Z",
				UpdatedAt:  "2022-01-01T00:00:00Z",
			},
		}

		res, err := s.withdrawClient.FindByIdWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawId})

		s.NoError(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure find by id withdraw", func() {
		withdrawId := int32(1)

		res, err := s.withdrawClient.FindByIdWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawId})

		s.Error(err)
		s.Nil(res)
	})

	s.Run("invalid id find by id withdraw", func() {
		withdrawId := int32(9999)

		res, err := s.withdrawClient.FindByIdWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: withdrawId})

		s.Error(err)
		s.Nil(res)
	})
}

func (s *ServerTestSuite) TestFindByCardNumberWithdraw() {
	s.Run("success find by card number withdraw", func() {
		req := &pb.FindByCardNumberRequest{
			CardNumber: "1234567890123456",
		}

		expectedResponse := &pb.ApiResponsesWithdraw{
			Status:  "success",
			Message: "Withdraw found successfully",
			Data: []*pb.WithdrawResponse{
				{
					WithdrawId: 1,
					CardNumber: "1234567890123456",
				},
				{
					WithdrawId: 2,
					CardNumber: "1234567890123456",
				},
			},
		}

		res, err := s.withdrawClient.FindByCardNumber(s.ctx, req)

		s.NoError(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure find by card number withdraw", func() {
		req := &pb.FindByCardNumberRequest{
			CardNumber: "1234567890123456",
		}

		res, err := s.withdrawClient.FindByCardNumber(s.ctx, req)

		s.Error(err)
		s.Nil(res)

	})
}

func (s *ServerTestSuite) TestFindActiveWithdraw() {
	s.Run("success find active withdraw", func() {
		expectedResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{
			Status:  "success",
			Message: "Active withdraws retrieved successfully",
			Pagination: &pb.PaginationMeta{
				TotalRecords: 2,
				TotalPages:   1,
			},
			Data: []*pb.WithdrawResponseDeleteAt{
				{
					WithdrawId: 1,
					CardNumber: "1234567890123456",
					CreatedAt:  "2022-01-01T00:00:00Z",
					UpdatedAt:  "2022-01-01T00:00:00Z",
				},
				{
					WithdrawId: 2,
					CardNumber: "9876543210987654",
					CreatedAt:  "2022-01-02T00:00:00Z",
					UpdatedAt:  "2022-01-02T00:00:00Z",
				},
			},
		}

		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindByActive(s.ctx, findAllWithdrawRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure find active withdraw", func() {
		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindByActive(s.ctx, findAllWithdrawRequest)

		s.Error(err)
		s.Nil(res)
	})

	s.Run("empty find active withdraw", func() {
		expectedResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{
			Status:  "success",
			Message: "No active withdraws found",
			Pagination: &pb.PaginationMeta{
				TotalRecords: 0,
				TotalPages:   0,
			},
			Data: []*pb.WithdrawResponseDeleteAt{},
		}

		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindByActive(s.ctx, findAllWithdrawRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, 0)
	})
}

func (s *ServerTestSuite) TestFindTrashedWithdraw() {
	s.Run("success find trashed withdraw", func() {
		expectedResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{
			Status:  "success",
			Message: "Trashed withdraws retrieved successfully",
			Pagination: &pb.PaginationMeta{
				TotalRecords: 2,
				TotalPages:   1,
			},
			Data: []*pb.WithdrawResponseDeleteAt{
				{
					WithdrawId: 1,
					CardNumber: "1234567890123456",
					CreatedAt:  "2022-01-01T00:00:00Z",
					UpdatedAt:  "2022-01-01T00:00:00Z",
				},
				{
					WithdrawId: 2,
					CardNumber: "9876543210987654",
					CreatedAt:  "2022-01-02T00:00:00Z",
					UpdatedAt:  "2022-01-02T00:00:00Z",
				},
			},
		}

		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindByTrashed(s.ctx, findAllWithdrawRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure find trashed withdraw", func() {
		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindByTrashed(s.ctx, findAllWithdrawRequest)

		s.Error(err)
		s.Nil(res)
	})

	s.Run("empty find trashed withdraw", func() {
		expectedResponse := &pb.ApiResponsePaginationWithdrawDeleteAt{
			Status:  "success",
			Message: "No trashed withdraws found",
			Pagination: &pb.PaginationMeta{
				TotalRecords: 0,
				TotalPages:   0,
			},
			Data: []*pb.WithdrawResponseDeleteAt{},
		}

		findAllWithdrawRequest := &pb.FindAllWithdrawRequest{
			Page:     1,
			PageSize: 10,
			Search:   "example",
		}

		res, err := s.withdrawClient.FindByTrashed(s.ctx, findAllWithdrawRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, 0)
	})
}

func (s *ServerTestSuite) TestCreateWithdraw() {
	s.Run("success create withdraw", func() {
		req := &pb.CreateWithdrawRequest{
			CardNumber:     "123456789",
			WithdrawAmount: 1000,
			WithdrawTime:   timestamppb.Now(),
		}

		expectedRespons := &pb.ApiResponseWithdraw{
			Status:  "success",
			Message: "Withdraw created successfully",
			Data: &pb.WithdrawResponse{
				WithdrawId: 1,
				CardNumber: "1234567890123456",
				CreatedAt:  "2022-01-01T00:00:00Z",
				UpdatedAt:  "2022-01-01T00:00:00Z",
			},
		}

		res, err := s.withdrawClient.CreateWithdraw(s.ctx, req)

		s.NoError(err)
		s.Equal(expectedRespons, res)
		s.Equal(expectedRespons.Status, res.Status)
		s.Equal(expectedRespons.Message, res.Message)
		s.Equal(expectedRespons.Data.WithdrawId, res.Data.WithdrawId)
		s.Equal(expectedRespons.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedRespons.Data.CreatedAt, res.Data.CreatedAt)
		s.Equal(expectedRespons.Data.UpdatedAt, res.Data.UpdatedAt)
	})

	s.Run("failure create withdraw", func() {
		req := &pb.CreateWithdrawRequest{
			CardNumber:     "123456789",
			WithdrawAmount: 1000,
			WithdrawTime:   timestamppb.Now(),
		}

		res, err := s.withdrawClient.CreateWithdraw(s.ctx, req)

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "internal server error")
	})

	s.Run("validation create withdraw", func() {
		req := &pb.CreateWithdrawRequest{
			CardNumber:     "",
			WithdrawAmount: 0,
			WithdrawTime:   timestamppb.Now(),
		}

		res, err := s.withdrawClient.CreateWithdraw(s.ctx, req)

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "validation error: invalid input")
	})
}

func (s *ServerTestSuite) TestUpdateWithdraw() {
	s.Run("success update withdraw", func() {
		req := &pb.UpdateWithdrawRequest{
			WithdrawId:     1,
			CardNumber:     "987654321",
			WithdrawAmount: 2000,
			WithdrawTime:   timestamppb.Now(),
		}

		expectedResponse := &pb.ApiResponseWithdraw{
			Status:  "success",
			Message: "Withdraw updated successfully",
			Data: &pb.WithdrawResponse{
				WithdrawId:     1,
				CardNumber:     "987654321",
				WithdrawAmount: 2000,
				CreatedAt:      "2022-01-01T00:00:00Z",
				UpdatedAt:      "2022-01-02T00:00:00Z",
			},
		}

		res, err := s.withdrawClient.UpdateWithdraw(s.ctx, req)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.WithdrawId, res.Data.WithdrawId)
		s.Equal(expectedResponse.Data.CardNumber, res.Data.CardNumber)
		s.Equal(expectedResponse.Data.WithdrawAmount, res.Data.WithdrawAmount)
		s.Equal(expectedResponse.Data.UpdatedAt, res.Data.UpdatedAt)
	})

	s.Run("failure update withdraw", func() {
		req := &pb.UpdateWithdrawRequest{
			WithdrawId:     1,
			CardNumber:     "987654321",
			WithdrawAmount: 2000,
			WithdrawTime:   timestamppb.Now(),
		}

		res, err := s.withdrawClient.UpdateWithdraw(s.ctx, req)

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "internal server error")
	})

	s.Run("validation update withdraw", func() {
		req := &pb.UpdateWithdrawRequest{
			WithdrawId:     0,
			CardNumber:     "",
			WithdrawAmount: 0,
			WithdrawTime:   nil,
		}

		res, err := s.withdrawClient.UpdateWithdraw(s.ctx, req)

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "validation error: invalid input")
	})
}

func (s *ServerTestSuite) TestTrashedWithdraw() {
	s.Run("success trashed withdraw", func() {
		expected := &pb.ApiResponseWithdraw{
			Status:  "success",
			Message: "Successfully trashed withdraw",
			Data: &pb.WithdrawResponse{
				WithdrawId: 1,
				CardNumber: "1234567890123456",
				CreatedAt:  "2022-01-01T00:00:00Z",
			},
		}

		res, err := s.withdrawClient.TrashedWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: 1})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expected.Status, res.Status)
		s.Equal(expected.Message, res.Message)
		s.Equal(expected.Data.WithdrawId, res.Data.WithdrawId)
		s.Equal(expected.Data.CardNumber, res.Data.CardNumber)
	})

	s.Run("failure trashed withdraw", func() {
		res, err := s.withdrawClient.TrashedWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: 2})

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "internal server error")
	})

	s.Run("invalid id trashed withdraw", func() {
		invalidID := int32(0)

		res, err := s.withdrawClient.TrashedWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: invalidID})

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "validation error: invalid withdraw ID")
	})
}

func (s *ServerTestSuite) TestRestoreWithdraw() {
	s.Run("success restore withdraw", func() {
		expected := &pb.ApiResponseWithdraw{
			Status:  "success",
			Message: "Successfully restored withdraw",
			Data: &pb.WithdrawResponse{
				WithdrawId: 1,
				CardNumber: "1234567890123456",
				CreatedAt:  "2022-01-01T00:00:00Z",
			},
		}

		res, err := s.withdrawClient.RestoreWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: 1})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expected.Status, res.Status)
		s.Equal(expected.Message, res.Message)
		s.Equal(expected.Data.WithdrawId, res.Data.WithdrawId)
		s.Equal(expected.Data.CardNumber, res.Data.CardNumber)
	})
	s.Run("failure restore withdraw", func() {
		res, err := s.withdrawClient.RestoreWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: 2})

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "internal server error")
	})

	s.Run("invalid id restore withdraw", func() {
		invalidID := int32(0)

		res, err := s.withdrawClient.RestoreWithdraw(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: invalidID})

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "validation error: invalid withdraw ID")
	})
}

func (s *ServerTestSuite) TestDeleteWithdraw() {
	s.Run("success delete permanent withdraw", func() {
		expected := &pb.ApiResponseWithdrawDelete{
			Status:  "success",
			Message: "Successfully deleted withdraw",
		}

		res, err := s.withdrawClient.DeleteWithdrawPermanent(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: 1})

		s.NoError(err)
		s.NotNil(res)
		s.Equal(expected.Status, res.Status)
		s.Equal(expected.Message, res.Message)
	})

	s.Run("failure delete permanent withdraw", func() {
		res, err := s.withdrawClient.DeleteWithdrawPermanent(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: 2})

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "internal server error")
	})

	s.Run("invalid id delete permanent withdraw", func() {
		invalidID := int32(0)

		res, err := s.withdrawClient.DeleteWithdrawPermanent(s.ctx, &pb.FindByIdWithdrawRequest{WithdrawId: invalidID})

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "validation error: invalid withdraw ID")
	})
}
