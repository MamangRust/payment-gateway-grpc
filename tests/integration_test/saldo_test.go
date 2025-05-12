package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"context"
)

func (s *ServerTestSuite) TestFindAllSaldo() {
	s.Run("success find all saldo", func() {
		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationSaldo{
			Status:  "success",
			Message: "Saldo data retrieved successfully",
			Data: []*pb.SaldoResponse{
				{
					SaldoId:      1,
					TotalBalance: 10000,
				},
				{
					SaldoId:      2,
					TotalBalance: 20000,
				},
			},
		}

		res, err := s.saldoClient.FindAllSaldo(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, 2)
		s.Equal(expectedResponse.Data[0].SaldoId, res.Data[0].SaldoId)
		s.Equal(expectedResponse.Data[0].TotalBalance, res.Data[0].TotalBalance)
	})

	s.Run("failure find all saldo", func() {
		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationSaldo{
			Status:  "error",
			Message: "Failed to retrieve saldo data",
		}

		res, err := s.saldoClient.FindAllSaldo(context.Background(), findAllRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("empty find all saldo", func() {
		page := 1
		pageSize := 10
		search := "notfound"

		expectedResponse := &pb.ApiResponsePaginationSaldo{
			Status:  "success",
			Message: "Saldo data retrieved successfully",
			Data:    []*pb.SaldoResponse{},
		}

		res, err := s.saldoClient.FindAllSaldo(context.Background(), &pb.FindAllSaldoRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		})

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Len(res.Data, 0)
	})
}

func (s *ServerTestSuite) TestFindByIdSaldo() {
	s.Run("success find saldo by id", func() {
		id := int32(1)

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo data retrieved successfully",
			Data: &pb.SaldoResponse{
				SaldoId:      id,
				TotalBalance: 10000,
			},
		}

		res, err := s.saldoClient.FindByIdSaldo(context.Background(), &pb.FindByIdSaldoRequest{
			SaldoId: id,
		})

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.SaldoId, res.Data.SaldoId)
		s.Equal(expectedResponse.Data.TotalBalance, res.Data.TotalBalance)
	})

	s.Run("failure find saldo by id", func() {
		id := int32(2)

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Saldo not found",
		}

		res, err := s.saldoClient.FindByIdSaldo(context.Background(), &pb.FindByIdSaldoRequest{
			SaldoId: id,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})

	s.Run("invalid id for find saldo by id", func() {
		id := int32(-1)

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Invalid Saldo ID",
		}

		res, err := s.saldoClient.FindByIdSaldo(context.Background(), &pb.FindByIdSaldoRequest{
			SaldoId: id,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByCardNumberSaldo() {
	s.Run("success find saldo by card number", func() {
		cardNumber := "1234567890"

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo data retrieved successfully",
			Data: &pb.SaldoResponse{
				SaldoId:      1,
				TotalBalance: 10000,
			},
		}

		res, err := s.saldoClient.FindByCardNumber(context.Background(), &pb.FindByCardNumberRequest{
			CardNumber: cardNumber,
		})

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.SaldoId, res.Data.SaldoId)
		s.Equal(expectedResponse.Data.TotalBalance, res.Data.TotalBalance)
	})

	s.Run("failure find saldo by card number", func() {
		cardNumber := "0987654321"

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Saldo not found",
		}

		res, err := s.saldoClient.FindByCardNumber(context.Background(), &pb.FindByCardNumberRequest{
			CardNumber: cardNumber,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByActiveSaldo() {
	s.Run("success find active saldo", func() {
		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationSaldoDeleteAt{
			Status:  "success",
			Message: "Saldo data retrieved successfully",
			Data: []*pb.SaldoResponseDeleteAt{
				{
					SaldoId:      1,
					TotalBalance: 10000,
				},
				{
					SaldoId:      2,
					TotalBalance: 20000,
				},
			},
		}

		res, err := s.saldoClient.FindByActive(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		s.Equal(expectedResponse.Data[0].SaldoId, res.Data[0].SaldoId)
		s.Equal(expectedResponse.Data[0].TotalBalance, res.Data[0].TotalBalance)
	})

	s.Run("failure find active saldo", func() {
		expectedResponse := &pb.ApiResponsesSaldo{
			Status:  "error",
			Message: "Failed to retrieve active saldo",
		}

		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.saldoClient.FindByActive(context.Background(), findAllRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})

	s.Run("empty active saldo", func() {
		expectedResponse := &pb.ApiResponsesSaldo{
			Status:  "success",
			Message: "No active saldo available",
			Data:    []*pb.SaldoResponse{},
		}

		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.saldoClient.FindByActive(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestFindByTrashedSaldo() {
	s.Run("success find trashed saldo", func() {
		expectedResponse := &pb.ApiResponsePaginationSaldoDeleteAt{
			Status:  "success",
			Message: "Saldo data retrieved successfully",
			Data: []*pb.SaldoResponseDeleteAt{
				{
					SaldoId:      1,
					TotalBalance: 10000,
				},
				{
					SaldoId:      2,
					TotalBalance: 20000,
				},
			},
		}

		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.saldoClient.FindByTrashed(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(len(expectedResponse.Data), len(res.Data))
		s.Equal(expectedResponse.Data[0].SaldoId, res.Data[0].SaldoId)
		s.Equal(expectedResponse.Data[0].TotalBalance, res.Data[0].TotalBalance)
	})

	s.Run("failure find trashed saldo", func() {
		expectedResponse := &pb.ApiResponsesSaldo{
			Status:  "error",
			Message: "Failed to retrieve trashed saldo data",
		}

		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.saldoClient.FindByTrashed(context.Background(), findAllRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})

	s.Run("empty trashed saldo", func() {
		expectedResponse := &pb.ApiResponsesSaldo{
			Status:  "success",
			Message: "No trashed saldo available",
			Data:    []*pb.SaldoResponse{},
		}

		findAllRequest := &pb.FindAllSaldoRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.saldoClient.FindByTrashed(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})
}

func (s *ServerTestSuite) TestCreateSaldo() {
	s.Run("success create saldo", func() {
		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo data created successfully",
			Data: &pb.SaldoResponse{
				SaldoId:      1,
				TotalBalance: 10000,
			},
		}

		res, err := s.saldoClient.CreateSaldo(context.Background(), &pb.CreateSaldoRequest{
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		})

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.SaldoId, res.Data.SaldoId)
		s.Equal(expectedResponse.Data.TotalBalance, res.Data.TotalBalance)
	})

	s.Run("failure create saldo", func() {
		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to create saldo data",
		}

		res, err := s.saldoClient.CreateSaldo(context.Background(), &pb.CreateSaldoRequest{
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})

	s.Run("validation create saldo", func() {
		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to create saldo data",
		}

		res, err := s.saldoClient.CreateSaldo(context.Background(), &pb.CreateSaldoRequest{
			CardNumber:   "1234567890",
			TotalBalance: 0,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})
}

func (s *ServerTestSuite) TestUpdateSaldo() {
	s.Run("success update saldo", func() {
		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo data updated successfully",
			Data: &pb.SaldoResponse{
				SaldoId:      1,
				TotalBalance: 10000,
			},
		}

		res, err := s.saldoClient.UpdateSaldo(context.Background(), &pb.UpdateSaldoRequest{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		})

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.SaldoId, res.Data.SaldoId)
		s.Equal(expectedResponse.Data.TotalBalance, res.Data.TotalBalance)
	})

	s.Run("failure update saldo", func() {
		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to update saldo data",
		}

		res, err := s.saldoClient.UpdateSaldo(context.Background(), &pb.UpdateSaldoRequest{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 10000,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})

	s.Run("validation update saldo", func() {
		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to update saldo data",
		}

		res, err := s.saldoClient.UpdateSaldo(context.Background(), &pb.UpdateSaldoRequest{
			SaldoId:      1,
			CardNumber:   "1234567890",
			TotalBalance: 0,
		})

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Nil(res.Data)
	})
}

func (s *ServerTestSuite) TestTrashSaldo() {
	s.Run("success trashed saldo", func() {
		trashedSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 1,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo trashed successfully",
		}

		res, err := s.saldoClient.TrashedSaldo(context.Background(), trashedSaldoRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure trashed saldo", func() {
		trashedSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 1,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to trashed saldo data",
		}

		res, err := s.saldoClient.TrashedSaldo(context.Background(), trashedSaldoRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("invalid id", func() {
		trashedSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 0,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to trashed saldo data",
		}

		res, err := s.saldoClient.TrashedSaldo(context.Background(), trashedSaldoRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestRestoreSaldo() {
	s.Run("success restore saldo", func() {
		restoreSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 1,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo restored successfully",
		}

		res, err := s.saldoClient.RestoreSaldo(context.Background(), restoreSaldoRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure restore saldo", func() {
		restoreSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 1,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to restore saldo data",
		}

		res, err := s.saldoClient.RestoreSaldo(context.Background(), restoreSaldoRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("invalid id", func() {
		restoreSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 0,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to restore saldo data",
		}

		res, err := s.saldoClient.RestoreSaldo(context.Background(), restoreSaldoRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestDeletePermanentSaldo() {
	s.Run("success delete saldo permanently", func() {
		deleteSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 1,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "success",
			Message: "Saldo deleted permanently",
		}

		res, err := s.saldoClient.DeleteSaldoPermanent(context.Background(), deleteSaldoRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure delete saldo permanently", func() {
		deleteSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 1,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to delete saldo data permanently",
		}

		res, err := s.saldoClient.DeleteSaldoPermanent(context.Background(), deleteSaldoRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("invalid id", func() {
		deleteSaldoRequest := &pb.FindByIdSaldoRequest{
			SaldoId: 0,
		}

		expectedResponse := &pb.ApiResponseSaldo{
			Status:  "error",
			Message: "Failed to delete saldo data permanently",
		}

		res, err := s.saldoClient.DeleteSaldoPermanent(context.Background(), deleteSaldoRequest)

		s.Error(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})
}
