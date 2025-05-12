package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerTestSuite) TestFindAllUser() {
	s.Run("Success Find All User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		findAllUserResponse := &pb.ApiResponsePaginationUser{
			Status:  "success",
			Message: "Successfully fetched users",
			Data: []*pb.UserResponse{
				{
					Id:        1,
					Firstname: "John",
					Lastname:  "Doe",
					Email:     "test@example.com",
					CreatedAt: "2024-12-30 03:29:39",
					UpdatedAt: "2024-12-30 03:29:39",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 1,
			},
		}

		res, err := s.userClient.FindAll(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(findAllUserResponse.Status, res.Status)
		s.Equal(findAllUserResponse.Message, res.Message)

		for i := range findAllUserResponse.Data {
			s.Equal(findAllUserResponse.Data[i].Id, res.Data[i].Id)
			s.Equal(findAllUserResponse.Data[i].Firstname, res.Data[i].Firstname)
			s.Equal(findAllUserResponse.Data[i].Lastname, res.Data[i].Lastname)
			s.Equal(findAllUserResponse.Data[i].Email, res.Data[i].Email)

			s.NotEmpty(res.Data[i].CreatedAt)
			s.NotEmpty(res.Data[i].UpdatedAt)
		}

		s.Equal(findAllUserResponse.Pagination.CurrentPage, res.Pagination.CurrentPage)
		s.Equal(findAllUserResponse.Pagination.PageSize, res.Pagination.PageSize)
		s.Equal(findAllUserResponse.Pagination.TotalPages, res.Pagination.TotalPages)
		s.Equal(findAllUserResponse.Pagination.TotalRecords, res.Pagination.TotalRecords)
	})

	s.Run("Empty Find All User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		findAllUserResponse := &pb.ApiResponsePaginationUser{
			Status:  "success",
			Message: "No users found",
			Data:    []*pb.UserResponse{},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 0,
			},
		}

		res, err := s.userClient.FindAll(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.Nil(res)
		s.Equal(findAllUserResponse.Status, res.Status)
		s.Equal(findAllUserResponse.Message, res.Message)
		s.Empty(res.Data)
		s.Equal(findAllUserResponse.Pagination.CurrentPage, res.Pagination.CurrentPage)
		s.Equal(findAllUserResponse.Pagination.PageSize, res.Pagination.PageSize)
		s.Equal(findAllUserResponse.Pagination.TotalPages, res.Pagination.TotalPages)
		s.Equal(findAllUserResponse.Pagination.TotalRecords, res.Pagination.TotalRecords)
	})

	s.Run("Failure Find All User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		findAllUserResponse := &pb.ApiResponsePaginationUser{
			Status:  "error",
			Message: "Failed to fetch users",
			Data:    []*pb.UserResponse{},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 0,
			},
		}

		res, err := s.userClient.FindAll(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(findAllUserResponse.Status, res.Status)
		s.Equal(findAllUserResponse.Message, res.Message)

	})

}

func (s *ServerTestSuite) TestFindByIdUser() {
	s.Run("Success Find By Id User", func() {

		findByIdUserRequest := &pb.FindByIdUserRequest{
			Id: 1,
		}

		findByIdUserResponse := &pb.ApiResponseUser{
			Status:  "success",
			Message: "Successfully fetched user",
			Data: &pb.UserResponse{
				Id:        1,
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "test@example.com",
			},
		}

		res, err := s.userClient.FindById(s.ctx, findByIdUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(findByIdUserResponse.Status, res.Status)
		s.Equal(findByIdUserResponse.Message, res.Message)
	})

	s.Run("Failure Find By Id User", func() {
		findByIdUserRequest := &pb.FindByIdUserRequest{
			Id: 1,
		}

		findByIdUserResponse := &pb.ApiResponseUser{
			Status:  "error",
			Message: "Failed to fetch user",
			Data:    nil,
		}

		res, err := s.userClient.FindById(s.ctx, findByIdUserRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(findByIdUserResponse.Status, res.Status)
		s.Equal(findByIdUserResponse.Message, res.Message)
	})
}

func (s *ServerTestSuite) TestActiveUser() {
	s.Run("Success Active User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}
		expectedResponse := &pb.ApiResponsePaginationUserDeleteAt{
			Status:  "success",
			Message: "Successfully fetched users",
			Data: []*pb.UserResponseWithDeleteAt{
				{
					Id:        1,
					Firstname: "John",
					Lastname:  "Doe",
					Email:     "test@example.com",
					CreatedAt: "2024-12-30 03:29:39",
					UpdatedAt: "2024-12-30 03:29:39",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 1,
			},
		}

		res, err := s.userClient.FindByActive(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
	})

	s.Run("Empty Active User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationUserDeleteAt{
			Status:  "success",
			Message: "Successfully fetched users",
			Data:    []*pb.UserResponseWithDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 1,
			},
		}

		res, err := s.userClient.FindByActive(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})

	s.Run("Failure Active User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}
		expectedError := status.Error(codes.Internal, "internal server error")

		res, err := s.userClient.FindByActive(s.ctx, findAllUserRequest)

		s.Error(err)
		s.NotNil(res)
		s.Equal(expectedError, err)
	})

}

func (s *ServerTestSuite) TestTrashedUser() {
	s.Run("Success Trashed User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedResponse := &pb.ApiResponsePaginationUserDeleteAt{
			Status:  "success",
			Message: "Successfully fetched users",
			Data: []*pb.UserResponseWithDeleteAt{
				{
					Id:        1,
					Firstname: "John",
					Lastname:  "Doe",
					Email:     "test@example.com",
					CreatedAt: "2024-12-30 03:29:39",
					UpdatedAt: "2024-12-30 03:29:39",
				},
			},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 1,
			},
		}

		res, err := s.userClient.FindByTrashed(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Equal(expectedResponse.Data, res.Data)
	})

	s.Run("Empty Trashed User", func() {
		expectedResponse := &pb.ApiResponsePaginationUserDeleteAt{
			Status:  "success",
			Message: "Successfully fetched users",
			Data:    []*pb.UserResponseWithDeleteAt{},
			Pagination: &pb.PaginationMeta{
				CurrentPage:  1,
				PageSize:     10,
				TotalPages:   0,
				TotalRecords: 1,
			},
		}
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		res, err := s.userClient.FindByTrashed(s.ctx, findAllUserRequest)

		s.NoError(err)
		s.NotNil(res)

		s.Equal(expectedResponse.Status, res.Status)
		s.Equal(expectedResponse.Message, res.Message)
		s.Empty(res.Data)
	})

	s.Run("Failure Trashed User", func() {
		findAllUserRequest := &pb.FindAllUserRequest{
			Page:     1,
			PageSize: 10,
			Search:   "",
		}

		expectedError := status.Error(codes.Internal, "internal server error")

		res, err := s.userClient.FindByTrashed(s.ctx, findAllUserRequest)

		s.Error(err)
		s.NotNil(res)
		s.Equal(expectedError, err)
	})
}

func (s *ServerTestSuite) TestCreate() {
	s.Run("Success Create User", func() {
		createUserRequest := &pb.CreateUserRequest{
			Firstname:       "Jane",
			Lastname:        "Doe",
			Email:           "testt@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		createUserResponse := &pb.ApiResponseUser{
			Status:  "success",
			Message: "Successfully created user",
			Data: &pb.UserResponse{
				Id:        2,
				Firstname: "Jane",
				Lastname:  "Doe",
				Email:     "testt@example.com",
			},
		}

		res, err := s.userClient.Create(s.ctx, createUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(createUserResponse.Status, res.Status)
		s.Equal(createUserResponse.Message, res.Message)
		s.Equal(createUserResponse.Data.Id, res.Data.Id)
		s.Equal(createUserResponse.Data.Firstname, res.Data.Firstname)
		s.Equal(createUserResponse.Data.Lastname, res.Data.Lastname)
		s.Equal(createUserResponse.Data.Email, res.Data.Email)

	})

	s.Run("Failure Create User", func() {
		createUserRequest := &pb.CreateUserRequest{
			Firstname:       "Jane",
			Lastname:        "Doe",
			Email:           "testt@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		res, err := s.userClient.Create(s.ctx, createUserRequest)

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "rpc error: code = Internal desc = internal server error") // Validasi pesan error
	})

	s.Run("Validation Error Create User", func() {
		createUserRequest := &pb.CreateUserRequest{
			Firstname:       "Jane",
			Lastname:        "Doe",
			Email:           "invalid-email",
			Password:        "password123",
			ConfirmPassword: "differentPassword",
		}

		res, err := s.userClient.Create(s.ctx, createUserRequest)

		s.Error(err)
		s.Nil(res)
		s.EqualError(err, "rpc error: code = InvalidArgument desc = validation error") // Validasi pesan error
	})

}

func (s *ServerTestSuite) TestUpdate() {
	s.Run("Success Update User", func() {
		updateUserRequest := &pb.UpdateUserRequest{
			Id:              2,
			Firstname:       "Jane",
			Lastname:        "Doe",
			Email:           "testt@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		updateUserResponse := &pb.ApiResponseUser{
			Status:  "success",
			Message: "Successfully updated user",
			Data: &pb.UserResponse{
				Id:        2,
				Firstname: "Jane",
				Lastname:  "Doe",
				Email:     "testt@example.com",
			},
		}

		res, err := s.userClient.Update(s.ctx, updateUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(updateUserResponse.Status, res.Status)
		s.Equal(updateUserResponse.Message, res.Message)
		s.Equal(updateUserResponse.Data.Id, res.Data.Id)
		s.Equal(updateUserResponse.Data.Firstname, res.Data.Firstname)
		s.Equal(updateUserResponse.Data.Lastname, res.Data.Lastname)
		s.Equal(updateUserResponse.Data.Email, res.Data.Email)
	})

	s.Run("Failure Update User", func() {
		updateUserRequest := &pb.UpdateUserRequest{
			Id:              2,
			Firstname:       "Jane",
			Lastname:        "Doe",
			Email:           "testt@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		res, err := s.userClient.Update(s.ctx, updateUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = Internal")
	})

	s.Run("Validation Error Update User", func() {
		updateUserRequest := &pb.UpdateUserRequest{
			Id:              2,
			Firstname:       "Jane",
			Lastname:        "Doe",
			Email:           "invalid-email",
			Password:        "password123",
			ConfirmPassword: "differentPassword",
		}

		res, err := s.userClient.Update(s.ctx, updateUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = InvalidArgument")
	})

}

func (s *ServerTestSuite) TestTrashUser() {
	s.Run("Success Trash User", func() {
		trashUserRequest := &pb.FindByIdUserRequest{
			Id: 2,
		}

		trashUserResponse := &pb.ApiResponseUser{
			Status:  "success",
			Message: "Successfully trashed user",
			Data: &pb.UserResponse{
				Id:        2,
				Firstname: "Jane",
				Lastname:  "Doe",
				Email:     "testt@example.com",
			},
		}

		res, err := s.userClient.TrashedUser(s.ctx, trashUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(trashUserResponse.Status, res.Status)
		s.Equal(trashUserResponse.Message, res.Message)
		s.Equal(trashUserResponse.Data.Id, res.Data.Id)
		s.Equal(trashUserResponse.Data.Firstname, res.Data.Firstname)
		s.Equal(trashUserResponse.Data.Lastname, res.Data.Lastname)
		s.Equal(trashUserResponse.Data.Email, res.Data.Email)
	})

	s.Run("Failure Trash User", func() {
		trashUserRequest := &pb.FindByIdUserRequest{
			Id: 2,
		}

		res, err := s.userClient.TrashedUser(s.ctx, trashUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = Internal")
	})

	s.Run("Invalid Trash User", func() {
		trashUserRequest := &pb.FindByIdUserRequest{
			Id: 0,
		}

		res, err := s.userClient.TrashedUser(s.ctx, trashUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = InvalidArgument")
	})

}

func (s *ServerTestSuite) TestRestoreUser() {
	s.Run("Success Restore User", func() {
		restoreUserRequest := &pb.FindByIdUserRequest{
			Id: 2,
		}

		restoreUserResponse := &pb.ApiResponseUser{
			Status:  "success",
			Message: "Successfully restored user",
			Data: &pb.UserResponse{
				Id:        2,
				Firstname: "Jane",
				Lastname:  "Doe",
				Email:     "testt@example.com",
			},
		}

		res, err := s.userClient.RestoreUser(s.ctx, restoreUserRequest)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(restoreUserResponse.Status, res.Status)
		s.Equal(restoreUserResponse.Message, res.Message)
		s.Equal(restoreUserResponse.Data.Id, res.Data.Id)
		s.Equal(restoreUserResponse.Data.Firstname, res.Data.Firstname)
	})

	s.Run("Failure Restore User", func() {
		restoreUserRequest := &pb.FindByIdUserRequest{
			Id: 2,
		}

		res, err := s.userClient.RestoreUser(s.ctx, restoreUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = Internal")
	})

	s.Run("Invalid ID Restore User", func() {
		restoreUserRequest := &pb.FindByIdUserRequest{
			Id: -1,
		}

		res, err := s.userClient.RestoreUser(s.ctx, restoreUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = InvalidArgument")
	})

}

func (s *ServerTestSuite) TestDeletePermanentUser() {
	s.Run("Success Delete Permanent User", func() {
		deleteUserRequest := &pb.FindByIdUserRequest{
			Id: 1,
		}

		deleteUserResponse := &pb.ApiResponseUserDelete{
			Status:  "success",
			Message: "Successfully deleted user record permanently",
		}

		res, err := s.userClient.DeleteUserPermanent(s.ctx, deleteUserRequest)

		fmt.Println("res", res)
		fmt.Println("err", err)

		s.NoError(err)
		s.NotNil(res)
		s.Equal(deleteUserResponse.Status, res.Status)
		s.Equal(deleteUserResponse.Message, res.Message)
	})

	s.Run("Failure Delete Permanent User", func() {
		deleteUserRequest := &pb.FindByIdUserRequest{
			Id: 1,
		}

		res, err := s.userClient.DeleteUserPermanent(s.ctx, deleteUserRequest)

		s.Error(err)
		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = Internal")
	})

	s.Run("Invalid ID Delete Permanent User", func() {
		deleteUserRequest := &pb.FindByIdUserRequest{
			Id: -1,
		}

		res, err := s.userClient.DeleteUserPermanent(s.ctx, deleteUserRequest)

		s.Error(err)

		s.Nil(res)

		s.Contains(err.Error(), "rpc error: code = InvalidArgument")
	})
}
