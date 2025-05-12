package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"context"
)

func (s *ServerTestSuite) TestFindAllMerchant() {
	s.Run("success find all merchant", func() {
		findAllRequest := &pb.FindAllMerchantRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationMerchant{
			Status:  "success",
			Message: "Merchants retrieved successfully",
			Data: []*pb.MerchantResponse{
				{Id: 1, Name: "Merchant 1"},
				{Id: 2, Name: "Merchant 2"},
			},
		}
		res, err := s.merchantClient.FindAllMerchant(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
		s.Equal(expectedResponse.Data[0].Name, res.Data[0].Name)
		s.Equal(expectedResponse.Data[1].Id, res.Data[1].Id)
		s.Equal(expectedResponse.Data[1].Name, res.Data[1].Name)
	})

	s.Run("failure find all merchant", func() {
		findAllRequest := &pb.FindAllMerchantRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationMerchant{
			Status:  "error",
			Message: "Failed to fetch merchant records",
		}
		res, err := s.merchantClient.FindAllMerchant(context.Background(), findAllRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("empty find all merchant", func() {
		findAllRequest := &pb.FindAllMerchantRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationMerchant{
			Status:  "success",
			Message: "Merchants retrieved successfully",
			Data:    []*pb.MerchantResponse{},
		}
		res, err := s.merchantClient.FindAllMerchant(context.Background(), findAllRequest)

		s.NoError(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestFindByIdMerchant() {
	s.Run("success find by id merchant", func() {
		findByIdRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant retrieved successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "Merchant 1",
			},
		}
		res, err := s.merchantClient.FindByIdMerchant(context.Background(), findByIdRequest)

		s.NoError(err)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.Name, res.Data.Name)
	})

	s.Run("failure find by id merchant", func() {
		findByIdRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Merchant not found",
		}
		res, err := s.merchantClient.FindByIdMerchant(context.Background(), findByIdRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid find by id merchant", func() {
		findByIdRequest := &pb.FindByIdMerchantRequest{
			MerchantId: -1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Merchant not found",
		}
		res, err := s.merchantClient.FindByIdMerchant(context.Background(), findByIdRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestFindByApikey() {
	s.Run("success find by api key", func() {
		findByApiKeyRequest := &pb.FindByApiKeyRequest{
			ApiKey: "my_api_key",
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant retrieved successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "Merchant 1",
			},
		}
		res, err := s.merchantClient.FindByApiKey(context.Background(), findByApiKeyRequest)

		s.NoError(err)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.Name, res.Data.Name)
	})

	s.Run("failure find by api key", func() {
		findByApiKeyRequest := &pb.FindByApiKeyRequest{
			ApiKey: "my_api_key",
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Merchant not found",
		}
		res, err := s.merchantClient.FindByApiKey(context.Background(), findByApiKeyRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)

	})
}

func (s *ServerTestSuite) TestFindByMerchantUserId() {
	s.Run("success find by merchant user id", func() {
		findByMerchantUserIdRequest := &pb.FindByMerchantUserIdRequest{
			UserId: 1,
		}

		expectedResponse := &pb.ApiResponsesMerchant{
			Status:  "success",
			Message: "Merchant retrieved successfully",
			Data: []*pb.MerchantResponse{
				{
					Id:   1,
					Name: "Merchant 1",
				},
			},
		}
		res, err := s.merchantClient.FindByMerchantUserId(context.Background(), findByMerchantUserIdRequest)

		s.NoError(err)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
		s.Equal(expectedResponse.Data[0].Name, res.Data[0].Name)

	})

	s.Run("failure find by merchant user id", func() {
		findByMerchantUserIdRequest := &pb.FindByMerchantUserIdRequest{
			UserId: 1,
		}

		expectedResponse := &pb.ApiResponsesMerchant{
			Status:  "error",
			Message: "Merchant not found",
		}
		res, err := s.merchantClient.FindByMerchantUserId(context.Background(), findByMerchantUserIdRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid find by merchant user id", func() {
		findByMerchantUserIdRequest := &pb.FindByMerchantUserIdRequest{
			UserId: -1,
		}

		expectedResponse := &pb.ApiResponsesMerchant{
			Status:  "error",
			Message: "Merchant not found",
		}
		res, err := s.merchantClient.FindByMerchantUserId(context.Background(), findByMerchantUserIdRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestFindByActiveMerchant() {
	s.Run("success find by active merchant", func() {
		findAllRequest := &pb.FindAllMerchantRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.merchantClient.FindByActive(context.Background(), findAllRequest)

		expectedResponse := &pb.ApiResponsePaginationMerchantDeleteAt{
			Status:  "success",
			Message: "Merchants retrieved successfully",
			Data: []*pb.MerchantResponseDeleteAt{
				{Id: 1, Name: "Merchant 1"},
				{Id: 2, Name: "Merchant 2"},
			},
		}

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)

		s.Equal(expectedResponse.Data[0].Id, res.Data[0].Id)
		s.Equal(expectedResponse.Data[0].Name, res.Data[0].Name)
	})

	s.Run("failure find by active merchant", func() {
		findAllRequest := &pb.FindAllMerchantRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.merchantClient.FindByActive(context.Background(), findAllRequest)

		expectedResponse := &pb.ApiResponsesMerchant{
			Status:  "error",
			Message: "Failed to retrieve merchant data: ",
		}
		s.Error(err)
		s.Equal(expectedResponse, res)

	})

	s.Run("empty find by active merchant", func() {
		findAllRequest := &pb.FindAllMerchantRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.merchantClient.FindByActive(context.Background(), findAllRequest)

		expectedResponse := &pb.ApiResponsesMerchant{
			Status:  "error",
			Message: "No active merchants found",
		}
		s.Error(err)
		s.Equal(expectedResponse, res)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(0, len(res.Data))

	})
}

func (s *ServerTestSuite) TestCreateMerchant() {
	s.Run("success create merchant", func() {
		createMerchantRequest := &pb.CreateMerchantRequest{
			Name:   "Merchant One",
			UserId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant created successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "Merchant One",
			},
		}
		res, err := s.merchantClient.CreateMerchant(context.Background(), createMerchantRequest)

		s.NoError(err)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)

		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.Name, res.Data.Name)
	})

	s.Run("failure create merchant", func() {
		createMerchantRequest := &pb.CreateMerchantRequest{
			Name:   "Merchant One",
			UserId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to create merchant",
		}
		res, err := s.merchantClient.CreateMerchant(context.Background(), createMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid create merchant", func() {
		createMerchantRequest := &pb.CreateMerchantRequest{
			Name:   "Merchant One",
			UserId: -1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to create merchant",
		}
		res, err := s.merchantClient.CreateMerchant(context.Background(), createMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestUpdateMerchant() {
	s.Run("success update merchant", func() {
		updateMerchantRequest := &pb.UpdateMerchantRequest{
			MerchantId: 1,
			Name:       "Updated Merchant Name",
			UserId:     1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant updated successfully",
			Data: &pb.MerchantResponse{
				Id:   1,
				Name: "Updated Merchant Name",
			},
		}
		res, err := s.merchantClient.UpdateMerchant(context.Background(), updateMerchantRequest)

		s.NoError(err)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data.Id, res.Data.Id)
		s.Equal(expectedResponse.Data.Name, res.Data.Name)
	})

	s.Run("failure update merchant", func() {
		updateMerchantRequest := &pb.UpdateMerchantRequest{
			MerchantId: 1,
			Name:       "Updated Merchant Name",
			UserId:     1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to update merchant",
		}
		res, err := s.merchantClient.UpdateMerchant(context.Background(), updateMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid update merchant", func() {
		updateMerchantRequest := &pb.UpdateMerchantRequest{
			MerchantId: -1,
			Name:       "Updated Merchant Name",
			UserId:     1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to update merchant",
		}
		res, err := s.merchantClient.UpdateMerchant(context.Background(), updateMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestTrashedMerchant() {
	s.Run("success trashed merchant", func() {
		trashedMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant trashed successfully",
		}
		res, err := s.merchantClient.TrashedMerchant(context.Background(), trashedMerchantRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure trashed merchant", func() {
		trashedMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to trash merchant",
		}
		res, err := s.merchantClient.TrashedMerchant(context.Background(), trashedMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid trashed merchant", func() {
		trashedMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: -1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to trash merchant",
		}
		res, err := s.merchantClient.TrashedMerchant(context.Background(), trashedMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestRestoreMerchant() {
	s.Run("success restore merchant", func() {
		restoreMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant restored successfully",
		}
		res, err := s.merchantClient.RestoreMerchant(context.Background(), restoreMerchantRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure restore merchant", func() {
		restoreMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to restore merchant",
		}
		res, err := s.merchantClient.RestoreMerchant(context.Background(), restoreMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid restore merchant", func() {
		restoreMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: -1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to restore merchant",
		}
		res, err := s.merchantClient.RestoreMerchant(context.Background(), restoreMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})
}

func (s *ServerTestSuite) TestDeleteMerchant() {
	s.Run("success delete merchant", func() {
		deleteMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "success",
			Message: "Merchant deleted successfully",
		}
		res, err := s.merchantClient.DeleteMerchantPermanent(context.Background(), deleteMerchantRequest)

		s.NoError(err)
		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("failure delete merchant", func() {
		deleteMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: 1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to delete merchant",
		}
		res, err := s.merchantClient.DeleteMerchantPermanent(context.Background(), deleteMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

	s.Run("invalid delete merchant", func() {
		deleteMerchantRequest := &pb.FindByIdMerchantRequest{
			MerchantId: -1,
		}

		expectedResponse := &pb.ApiResponseMerchant{
			Status:  "error",
			Message: "Failed to delete merchant",
		}
		res, err := s.merchantClient.DeleteMerchantPermanent(context.Background(), deleteMerchantRequest)

		s.Error(err)
		s.Equal(expectedResponse, res)
	})

}
